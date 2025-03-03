package ui

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/rivo/tview"

	"github.com/gdamore/tcell/v2"
)

// TreeView é uma visualização em árvore de diretórios
type TreeView struct {
	*tview.TreeView
	app            *App
	rootDir        string
	showHiddenDirs bool
}

// NewTreeView cria uma nova visualização em árvore
func NewTreeView(app *App) *TreeView {
	// Criar árvore
	rootNode := tview.NewTreeNode("").SetReference("/")
	treeView := tview.NewTreeView().SetRoot(rootNode).SetCurrentNode(rootNode)

	// Configurar cores
	treeView.SetGraphicsColor(tcell.ColorGreen)

	// Configurar borda
	treeView.SetBorder(true).
		SetTitle(" Árvore ").
		SetTitleAlign(tview.AlignLeft)

	// Criar TreeView
	t := &TreeView{
		TreeView:       treeView,
		app:            app,
		showHiddenDirs: false,
	}

	// Configurar manipulador de seleção
	treeView.SetSelectedFunc(func(node *tview.TreeNode) {
		// Obter caminho do nó
		reference := node.GetReference()
		if reference == nil {
			return
		}
		path := reference.(string)

		// Navegar para o diretório
		t.app.navigateTo(path)
	})

	return t
}

// SetRootDir define o diretório raiz
func (t *TreeView) SetRootDir(rootDir string) {
	t.rootDir = rootDir
	t.Refresh()
}

// SetShowHiddenDirs define se diretórios ocultos são exibidos
func (t *TreeView) SetShowHiddenDirs(show bool) {
	t.showHiddenDirs = show
	t.Refresh()
}

// Refresh atualiza a visualização em árvore
func (t *TreeView) Refresh() {
	// Limpar árvore
	t.SetRoot(nil)

	// Verificar se o diretório raiz existe
	if t.rootDir == "" {
		return
	}

	// Criar nó raiz
	root := tview.NewTreeNode(filepath.Base(t.rootDir)).
		SetReference(t.rootDir).
		SetSelectable(true).
		SetExpanded(true).
		SetColor(tcell.ColorYellow)

	// Adicionar nó raiz
	t.SetRoot(root)

	// Adicionar nó para o diretório pai (..)
	parentDir := filepath.Dir(t.rootDir)
	if parentDir != t.rootDir { // Não adicionar ".." se estiver na raiz
		parentNode := tview.NewTreeNode("..").
			SetReference(parentDir).
			SetSelectable(true).
			SetColor(tcell.ColorBlue)
		root.AddChild(parentNode)
	}

	// Adicionar subdiretórios
	t.addSubDirectories(root, t.rootDir)
}

// addSubDirectories adiciona subdiretórios a um nó
func (t *TreeView) addSubDirectories(node *tview.TreeNode, path string) {
	// Listar diretórios
	entries, err := os.ReadDir(path)
	if err != nil {
		return
	}

	// Adicionar subdiretórios
	for _, entry := range entries {
		// Verificar se é um diretório
		if !entry.IsDir() {
			continue
		}

		// Verificar se é um diretório oculto
		name := entry.Name()
		if !t.showHiddenDirs && strings.HasPrefix(name, ".") {
			continue
		}

		// Criar caminho completo
		childPath := filepath.Join(path, name)

		// Criar nó
		child := tview.NewTreeNode(name).
			SetReference(childPath).
			SetSelectable(true).
			SetColor(tcell.ColorYellow)

		// Adicionar nó ao pai
		node.AddChild(child)
	}
}

// ExpandSelected expande o nó selecionado
func (t *TreeView) ExpandSelected() {
	// Obter nó selecionado
	node := t.GetCurrentNode()
	if node == nil {
		return
	}

	// Verificar se o nó já está expandido
	if node.IsExpanded() {
		return
	}

	// Expandir nó
	node.SetExpanded(true)

	// Obter caminho do nó
	reference := node.GetReference()
	if reference == nil {
		return
	}
	path := reference.(string)

	// Adicionar subdiretórios
	t.addSubDirectories(node, path)
}

// CollapseSelected colapsa o nó selecionado
func (t *TreeView) CollapseSelected() {
	// Obter nó selecionado
	node := t.GetCurrentNode()
	if node == nil {
		return
	}

	// Colapsar nó
	node.SetExpanded(false)
	node.ClearChildren()
}

// GetCurrentNode retorna o nó atualmente selecionado
func (t *TreeView) GetCurrentNode() *tview.TreeNode {
	return t.TreeView.GetCurrentNode()
}

// LoadTree carrega a árvore de diretórios
func (t *TreeView) LoadTree(rootDir string) {
	t.SetRootDir(rootDir)
}

// GetSelectedDirectory retorna o diretório selecionado
func (t *TreeView) GetSelectedDirectory() string {
	// Obter nó selecionado
	node := t.GetCurrentNode()
	if node == nil {
		return ""
	}

	// Obter referência do nó
	ref := node.GetReference()
	if ref == nil {
		return ""
	}

	// Converter para string
	path, ok := ref.(string)
	if !ok {
		return ""
	}

	return path
}

// UpdateTreeView atualiza a árvore de diretórios
func (t *TreeView) UpdateTreeView(currentDir string) {
	// Atualizar o diretório raiz
	t.rootDir = currentDir

	// Limpar a árvore
	t.TreeView.SetRoot(tview.NewTreeNode(""))

	// Criar nó raiz
	root := tview.NewTreeNode(filepath.Base(currentDir)).
		SetReference(currentDir).
		SetSelectable(true).
		SetExpanded(true).
		SetColor(tcell.ColorYellow)

	// Adicionar nó raiz
	t.SetRoot(root)

	// Adicionar nó para o diretório pai (..)
	parentDir := filepath.Dir(currentDir)
	if parentDir != currentDir { // Não adicionar ".." se estiver na raiz
		parentNode := tview.NewTreeNode("..").
			SetReference(parentDir).
			SetSelectable(true).
			SetColor(tcell.ColorBlue)
		root.AddChild(parentNode)
	}

	// Adicionar subdiretórios
	t.addSubDirectories(root, currentDir)

	// Expandir o nó raiz
	root.SetExpanded(true)
}

// addNodes adiciona nós filhos a um nó pai
func (t *TreeView) addNodes(node *tview.TreeNode, path string, level int) {
	// Limitar a profundidade para evitar sobrecarga
	if level > 3 {
		return
	}

	// Ler diretório
	files, err := os.ReadDir(path)
	if err != nil {
		return
	}

	// Adicionar diretórios
	for _, file := range files {
		if file.IsDir() {
			// Ignorar diretórios ocultos
			if strings.HasPrefix(file.Name(), ".") {
				continue
			}

			// Caminho completo
			childPath := filepath.Join(path, file.Name())

			// Criar nó filho
			child := tview.NewTreeNode(file.Name()).
				SetReference(childPath).
				SetSelectable(true).
				SetColor(tcell.ColorGreen)

			// Adicionar ao nó pai
			node.AddChild(child)

			// Adicionar nós filhos recursivamente (com limite de profundidade)
			if level < 2 {
				t.addNodes(child, childPath, level+1)
			}
		}
	}
}

// expandToPath expande a árvore até o caminho especificado
func (t *TreeView) expandToPath(node *tview.TreeNode, path string) {
	// Verificar se o nó tem referência
	ref := node.GetReference()
	if ref == nil {
		return
	}

	// Obter caminho do nó
	nodePath := ref.(string)

	// Verificar se o caminho contém o caminho do nó
	if strings.HasPrefix(path, nodePath) {
		// Expandir nó
		node.SetExpanded(true)

		// Expandir nós filhos
		for _, child := range node.GetChildren() {
			t.expandToPath(child, path)
		}
	}
}

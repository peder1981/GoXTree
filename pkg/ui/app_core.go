package ui

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// App representa a aplicação principal
type App struct {
	app            *tview.Application
	pages          *tview.Pages
	mainLayout     *tview.Flex
	fileView       *FileView
	treeView       *TreeView
	statusBar      *StatusBar
	menuBar        *MenuBar
	currentDir     string
	history        []string
	historyPos     int
	showHidden     bool
	selectedFiles  map[string]bool
	clipboard      string
	clipboardIsDir bool
}

// Clipboard representa a área de transferência
type Clipboard struct {
	Files     []string
	Operation string // "copy" ou "move"
}

// ViewMode representa o modo de visualização
type ViewMode int

const (
	ViewModeTree ViewMode = iota
	ViewModeFlat
	ViewModeDetails
)

// Constantes para os tipos de foco
const (
	focusTree   = iota // Foco na árvore de diretórios
	focusFiles         // Foco na lista de arquivos
	focusSearch        // Foco na busca
)

// NewApp cria uma nova instância da aplicação
func NewApp() *App {
	// Criar aplicação
	app := &App{
		app:        tview.NewApplication(),
		pages:      tview.NewPages(),
		mainLayout: tview.NewFlex(),
		history:    make([]string, 0),
		historyPos: -1,
		showHidden: false,
	}

	// Obter diretório atual
	var err error
	app.currentDir, err = os.Getwd()
	if err != nil {
		app.currentDir = "/"
	}

	// Criar componentes
	app.treeView = NewTreeView(app)
	app.fileView = NewFileView(app)
	app.statusBar = NewStatusBar()
	app.menuBar = NewMenuBar(app)

	// Configurar layout principal como vertical (Row)
	app.mainLayout.SetDirection(tview.FlexRow)

	// Criar barra de função para exibir as opções
	functionBar := tview.NewTextView().
		SetDynamicColors(true).
		SetText("[blue]F1[white]-Ajuda [blue]F2[white]-Renomear [blue]F7[white]-Criar Dir [blue]F8[white]-Excluir [blue]F10[white]-Sair [blue]Tab[white]-Alternar Painel [blue]ESC[white]-Voltar").
		SetTextAlign(tview.AlignCenter)

	// Adicionar menu bar ao topo
	app.mainLayout.AddItem(app.menuBar.menuBar, 1, 0, false)
	
	// Criar layout horizontal para árvore e lista de arquivos
	horizontalLayout := tview.NewFlex().SetDirection(tview.FlexColumn)
	horizontalLayout.AddItem(app.treeView.TreeView, 0, 1, true)
	horizontalLayout.AddItem(app.fileView.fileList, 0, 2, false)
	
	// Adicionar o layout horizontal ao layout principal
	app.mainLayout.AddItem(horizontalLayout, 0, 1, true)
	
	// Adicionar barra de status
	app.mainLayout.AddItem(app.statusBar.statusBar, 2, 0, false)
	
	// Adicionar barra de função na parte inferior
	app.mainLayout.AddItem(functionBar, 1, 0, false)

	// Configurar páginas
	app.pages.AddPage("main", app.mainLayout, true, true)

	// Configurar aplicação
	app.app.SetRoot(app.pages, true)

	// Configurar manipuladores de eventos
	app.setupKeyBindings()

	// Adicionar diretório inicial ao histórico
	app.addToHistory(app.currentDir)
	app.historyPos = 0

	// Carregar diretório inicial
	app.treeView.LoadTree(app.currentDir)
	if err := app.fileView.SetCurrentDir(app.currentDir); err != nil {
		fmt.Println(err)
	}
	app.refreshStatus()

	return app
}

// Run inicia a aplicação
func (a *App) Run() error {
	// Definir diretório inicial
	var err error
	a.currentDir, err = os.Getwd()
	if err != nil {
		a.currentDir = "/"
	}

	// Carregar visualizações
	a.treeView.LoadTree(a.currentDir)
	a.fileView.SetCurrentDir(a.currentDir)

	// Atualizar barra de status
	a.statusBar.UpdateStatus(a.currentDir)

	// Definir foco inicial
	a.app.SetFocus(a.treeView.TreeView)

	// Iniciar aplicação
	return a.app.SetRoot(a.pages, true).Run()
}

// setupKeyBindings configura os atalhos de teclado
func (a *App) setupKeyBindings() {
	a.app.SetInputCapture(a.handleKeyEvents)
}

// refreshStatus atualiza a barra de status
func (a *App) refreshStatus() {
	a.statusBar.UpdateStatus(a.currentDir)
}

// refreshAll atualiza todas as visualizações
func (a *App) refreshAll() {
	a.refreshView()
	a.refreshStatus()
}

// refreshViews atualiza as visualizações
func (a *App) refreshViews() {
	a.refreshView()
}

// setFocus define o foco em um componente
func (a *App) setFocus(p tview.Primitive) {
	a.app.SetFocus(p)
}

// confirmExit confirma a saída da aplicação
func (a *App) confirmExit() {
	a.showConfirmDialog("Sair", "Deseja realmente sair da aplicação?", func(confirmed bool) {
		if confirmed {
			a.app.Stop()
		}
	})
}

// showSystemInfo exibe informações do sistema
func (a *App) showSystemInfo() {
	// Obter informações do sistema
	hostname, _ := os.Hostname()
	wd, _ := os.Getwd()
	homeDir, _ := os.UserHomeDir()

	// Obter informações do diretório atual
	var dirSize int64
	var fileCount, dirCount int

	err := filepath.Walk(a.currentDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			dirCount++
		} else {
			fileCount++
			dirSize += info.Size()
		}

		return nil
	})

	// Formatar tamanho do diretório
	sizeStr := "Erro ao calcular tamanho"
	if err == nil {
		sizeStr = fmt.Sprintf("%d bytes", dirSize)
	}

	// Criar conteúdo
	content := fmt.Sprintf(
		"Informações do Sistema:\n"+
			"-------------------\n"+
			"Sistema Operacional: %s\n"+
			"Arquitetura: %s\n"+
			"Hostname: %s\n"+
			"Diretório Atual: %s\n"+
			"Diretório Home: %s\n"+
			"Versão Go: %s\n\n"+
			"Informações do Diretório:\n"+
			"----------------------\n"+
			"Tamanho Total: %s\n"+
			"Arquivos: %d\n"+
			"Diretórios: %d\n",
		"linux",
		"amd64",
		hostname,
		wd,
		homeDir,
		"1.17.2",
		sizeStr,
		fileCount,
		dirCount,
	)

	// Exibir diálogo
	textView := tview.NewTextView().
		SetText(content).
		SetDynamicColors(true).
		SetBorder(true).
		SetTitle("Informações do Sistema").
		SetTitleAlign(tview.AlignCenter)

	textView.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEscape || event.Key() == tcell.KeyEnter {
			a.pages.RemovePage("sysinfo")
			return nil
		}
		return event
	})

	a.pages.AddPage("sysinfo", textView, true, true)
	a.app.SetFocus(textView)
}

// viewFile visualiza o arquivo selecionado
func (a *App) viewFile() {
	selectedFile := a.fileView.GetSelectedFile()
	if selectedFile == "" {
		a.showError("Nenhum arquivo selecionado")
		return
	}

	// Verificar se é um diretório
	fileInfo, err := os.Stat(filepath.Join(a.currentDir, selectedFile))
	if err != nil {
		a.showError(fmt.Sprintf("Erro ao acessar arquivo: %v", err))
		return
	}

	if fileInfo.IsDir() {
		a.navigateTo(filepath.Join(a.currentDir, selectedFile))
		return
	}

	// Ler conteúdo do arquivo
	content, err := os.ReadFile(filepath.Join(a.currentDir, selectedFile))
	if err != nil {
		a.showError(fmt.Sprintf("Erro ao ler arquivo: %v", err))
		return
	}

	// Exibir conteúdo
	textView := tview.NewTextView().
		SetText(string(content)).
		SetScrollable(true).
		SetBorder(true).
		SetTitle(filepath.Base(selectedFile)).
		SetTitleAlign(tview.AlignCenter)

	textView.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEscape {
			a.pages.RemovePage("fileView")
			return nil
		}
		return event
	})

	a.pages.AddPage("fileView", textView, true, true)
	a.app.SetFocus(textView)
}

// addToHistory adiciona um diretório ao histórico de navegação
func (a *App) addToHistory(dir string) {
	// Se estamos no meio do histórico, remover tudo após a posição atual
	if a.historyPos >= 0 && a.historyPos < len(a.history)-1 {
		a.history = a.history[:a.historyPos+1]
	}

	// Verificar se o diretório já é o último no histórico
	if len(a.history) > 0 && a.history[len(a.history)-1] == dir {
		return
	}

	// Adicionar diretório ao histórico
	a.history = append(a.history, dir)
	a.historyPos = len(a.history) - 1
}

// handleKeyEvents manipula eventos de teclado
func (a *App) handleKeyEvents(event *tcell.EventKey) *tcell.EventKey {
	// Verificar teclas globais
	switch event.Key() {
	case tcell.KeyF1:
		// Mostrar ajuda
		a.showHelpDialog()
		return nil
	case tcell.KeyF10:
		// Sair da aplicação
		a.confirmExit()
		return nil
	case tcell.KeyEscape:
		// Comportamento universal do ESC
		if len(a.history) > 1 && a.historyPos > 0 {
			// Se há histórico, voltar para o diretório anterior
			a.goBack()
		} else {
			// Se não há histórico, perguntar se deseja sair
			a.confirmExit()
		}
		return nil
	}

	// Verificar teclas específicas para a visualização atual
	switch a.app.GetFocus() {
	case a.fileView.fileList:
		return a.handleFileViewKeys(event)
	case a.treeView.TreeView:
		return a.handleTreeViewKeys(event)
	}

	return event
}

// handleFileViewKeys manipula teclas na visualização de arquivos
func (a *App) handleFileViewKeys(event *tcell.EventKey) *tcell.EventKey {
	switch event.Key() {
	case tcell.KeyEnter:
		// Abrir arquivo ou diretório
		a.openFile()
		return nil
	case tcell.KeyF2:
		// Renomear arquivo
		a.renameFile()
		return nil
	case tcell.KeyF7:
		// Criar diretório
		a.createDirectory()
		return nil
	case tcell.KeyF8:
		// Excluir arquivo
		a.deleteFile()
		return nil
	case tcell.KeyTab:
		// Alternar foco para a árvore
		a.app.SetFocus(a.treeView.TreeView)
		return nil
	}

	// Verificar teclas de letra
	switch event.Rune() {
	case 'h':
		// Alternar exibição de arquivos ocultos
		a.toggleHidden()
		return nil
	case 'r':
		// Atualizar visualização
		a.refreshView()
		return nil
	}

	return event
}

// handleTreeViewKeys manipula teclas na visualização de árvore
func (a *App) handleTreeViewKeys(event *tcell.EventKey) *tcell.EventKey {
	switch event.Key() {
	case tcell.KeyEnter:
		// Navegar para o diretório selecionado
		a.navigateToSelectedDirectory()
		return nil
	case tcell.KeyTab:
		// Alternar foco para a lista de arquivos
		a.app.SetFocus(a.fileView.fileList)
		return nil
	}

	return event
}

// toggleHidden alterna a exibição de arquivos ocultos
func (a *App) toggleHidden() {
	a.showHidden = !a.showHidden
	a.fileView.SetShowHidden(a.showHidden)
}

// navigateToSelectedDirectory navega para o diretório selecionado na árvore
func (a *App) navigateToSelectedDirectory() {
	selectedDir := a.treeView.GetSelectedDirectory()
	if selectedDir != "" {
		a.navigateTo(selectedDir)
	}
}

// showError exibe uma mensagem de erro
func (a *App) showError(message string) {
	modal := tview.NewModal().
		SetText(message).
		AddButtons([]string{"OK"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			a.pages.RemovePage("error")
		})

	a.pages.AddPage("error", modal, true, true)
	a.app.SetFocus(modal)
}

// showMessage exibe uma mensagem informativa
func (a *App) showMessage(message string) {
	modal := tview.NewModal().
		SetText(message).
		AddButtons([]string{"OK"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			a.pages.RemovePage("message")
		})

	a.pages.AddPage("message", modal, true, true)
	a.app.SetFocus(modal)
}

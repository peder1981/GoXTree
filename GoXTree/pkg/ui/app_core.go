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

	// Aplicar tema retrô
	ApplyRetroThemeToApp(app)

	// Configurar layout principal como vertical (Row)
	app.mainLayout.SetDirection(tview.FlexRow)

	// Criar barra de função para exibir as opções
	functionBar := tview.NewTextView().
		SetDynamicColors(true).
		SetText("[yellow]F1[white]-Ajuda [yellow]F2[white]-Renomear [yellow]F3[white]-Busca [yellow]F4[white]-Busca Av. [yellow]F7[white]-Criar Dir [yellow]F8[white]-Excluir [yellow]F9[white]-Sincronizar [yellow]F10[white]-Sair").
		SetTextAlign(tview.AlignCenter).
		SetBackgroundColor(ColorBackground)

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
	app.SetupKeyHandlers()

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

// SetupKeyHandlers configura os manipuladores de teclas
func (a *App) SetupKeyHandlers() {
	a.app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		// Verificar se alguma tecla de função foi pressionada
		switch event.Key() {
		case tcell.KeyF1:
			a.showHelp()
			return nil
		case tcell.KeyF2:
			a.renameFile()
			return nil
		case tcell.KeyF3:
			a.searchFiles()
			return nil
		case tcell.KeyF4:
			a.advancedSearch()
			return nil
		case tcell.KeyF7:
			a.createDirectory()
			return nil
		case tcell.KeyF8:
			a.deleteFile()
			return nil
		case tcell.KeyF9:
			a.syncDirectories()
			return nil
		case tcell.KeyF10:
			a.app.Stop()
			return nil
		case tcell.KeyTab:
			a.toggleFocus()
			return nil
		case tcell.KeyEscape:
			// Comportamento do ESC depende do contexto
			if a.pages.HasPage("help") {
				a.pages.RemovePage("help")
				return nil
			} else if a.pages.HasPage("search") {
				a.pages.RemovePage("search")
				return nil
			} else if a.pages.HasPage("advancedSearch") {
				a.pages.RemovePage("advancedSearch")
				return nil
			} else if a.pages.HasPage("createDir") {
				a.pages.RemovePage("createDir")
				return nil
			} else if a.pages.HasPage("rename") {
				a.pages.RemovePage("rename")
				return nil
			} else if a.pages.HasPage("delete") {
				a.pages.RemovePage("delete")
				return nil
			} else if a.pages.HasPage("sync") {
				a.pages.RemovePage("sync")
				return nil
			} else if a.pages.HasPage("view") {
				a.pages.RemovePage("view")
				return nil
			} else if a.pages.HasPage("edit") {
				a.pages.RemovePage("edit")
				return nil
			} else if a.pages.HasPage("compare") {
				a.pages.RemovePage("compare")
				return nil
			} else if a.pages.HasPage("sysinfo") {
				a.pages.RemovePage("sysinfo")
				return nil
			}
		}

		// Verificar combinações de teclas Ctrl+letra
		if event.Modifiers() == tcell.ModCtrl {
			switch event.Rune() {
			case 'a', 'A': // Ctrl+A: Selecionar todos os arquivos
				a.selectAllFiles()
				return nil
			case 'd', 'D': // Ctrl+D: Desmarcar todos os arquivos
				a.unselectAllFiles()
				return nil
			case 'f', 'F': // Ctrl+F: Buscar arquivo
				a.searchFiles()
				return nil
			case 'g', 'G': // Ctrl+G: Ir para diretório
				a.goToDirectory()
				return nil
			case 'h', 'H': // Ctrl+H: Alternar arquivos ocultos
				a.toggleHiddenFiles()
				return nil
			case 'r', 'R': // Ctrl+R: Atualizar visualização
				a.refreshView()
				return nil
			case 's', 'S': // Ctrl+S: Selecionar/desmarcar arquivo atual
				a.toggleSelection()
				return nil
			case 'c', 'C': // Ctrl+C: Comparar arquivos selecionados
				a.compareSelectedFiles()
				return nil
			case 'v', 'V': // Ctrl+V: Visualizar arquivo
				a.viewCurrentFile()
				return nil
			case 'e', 'E': // Ctrl+E: Editar arquivo
				a.editCurrentFile()
				return nil
			case 'y', 'Y': // Ctrl+Y: Sincronizar diretórios
				a.syncDirectories()
				return nil
			case 'i', 'I': // Ctrl+I: Informações do sistema
				a.showSystemInfo()
				return nil
			}
		}

		return event
	})
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
		// Verificar se estamos na tela principal ou em uma tela de diálogo
		if a.pages.HasPage("textDialog") || a.pages.HasPage("inputDialog") || 
		   a.pages.HasPage("helpDialog") || a.pages.HasPage("menuDialog") || 
		   a.pages.HasPage("gotoDialog") || a.pages.HasPage("input") {
			// Se estamos em uma tela de diálogo, apenas fechar o diálogo
			return event // Deixar o manipulador específico da página tratar o evento
		} else if len(a.history) > 1 && a.historyPos > 0 {
			// Se estamos na tela principal e há histórico, voltar para o diretório anterior
			a.goBack()
		} else {
			// Se estamos na tela principal e não há histórico, perguntar se deseja sair
			a.confirmExit()
		}
		return nil
	case tcell.KeyCtrlA:
		// Selecionar todos os arquivos
		a.selectAllFiles()
		return nil
	case tcell.KeyCtrlD:
		// Desmarcar todos os arquivos
		a.unselectAllFiles()
		return nil
	case tcell.KeyCtrlF:
		// Buscar arquivo
		a.showSimpleSearchDialog()
		return nil
	case tcell.KeyCtrlG:
		// Ir para diretório
		a.showGotoDialog()
		return nil
	case tcell.KeyCtrlH:
		// Mostrar/ocultar arquivos ocultos
		a.toggleHidden()
		return nil
	case tcell.KeyCtrlR:
		// Atualizar visualização
		a.refreshView()
		return nil
	case tcell.KeyCtrlC:
		// Comparar arquivos selecionados
		a.compareSelectedFiles()
		return nil
	case tcell.KeyCtrlV:
		// Visualizar arquivo
		a.viewCurrentFile()
		return nil
	case tcell.KeyCtrlE:
		// Editar arquivo
		a.editCurrentFile()
		return nil
	case tcell.KeyCtrlY:
		// Sincronizar diretórios
		a.syncDirectoriesDialog()
		return nil
	case tcell.KeyCtrlS:
		// Selecionar arquivo
		a.toggleSelection()
		return nil
	case tcell.KeyCtrlI:
		// Informações do sistema
		a.showSystemInfo()
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

// refreshView atualiza a visualização
func (a *App) refreshView() {
	a.treeView.Refresh()
	a.updateFileList()
	a.statusBar.UpdateStatus(a.currentDir)
}

// updateFileList atualiza a lista de arquivos
func (a *App) updateFileList() {
	files, err := utils.ListFiles(a.currentDir)
	if err != nil {
		a.showError("Erro ao listar arquivos: " + err.Error())
		return
	}
	a.fileView.UpdateFileList(files, a.showHidden)
}

// searchFiles abre o diálogo de busca de arquivos
func (a *App) searchFiles() {
	a.showMessage("Busca de arquivos não implementada")
}

// advancedSearch abre o diálogo de busca avançada
func (a *App) advancedSearch() {
	a.showMessage("Busca avançada não implementada")
}

// syncDirectories sincroniza dois diretórios
func (a *App) syncDirectories() {
	a.showMessage("Sincronização de diretórios não implementada")
}

// goToDirectory abre o diálogo para ir para um diretório específico
func (a *App) goToDirectory() {
	a.showMessage("Ir para diretório não implementado")
}

// toggleHiddenFiles alterna a exibição de arquivos ocultos
func (a *App) toggleHiddenFiles() {
	a.showHidden = !a.showHidden
	a.updateFileList()
	if a.showHidden {
		a.statusBar.SetStatus("Arquivos ocultos: Visíveis")
	} else {
		a.statusBar.SetStatus("Arquivos ocultos: Ocultos")
	}
}

// showSystemInfo exibe informações do sistema
func (a *App) showSystemInfo() {
	a.showMessage("Informações do sistema não implementadas")
}

// renameFile abre o diálogo para renomear um arquivo
func (a *App) renameFile() {
	a.showMessage("Renomear arquivo não implementado")
}

// createDirectory abre o diálogo para criar um diretório
func (a *App) createDirectory() {
	a.showMessage("Criar diretório não implementado")
}

// deleteFile abre o diálogo para excluir um arquivo
func (a *App) deleteFile() {
	a.showMessage("Excluir arquivo não implementado")
}

// toggleFocus alterna o foco entre a árvore e a lista de arquivos
func (a *App) toggleFocus() {
	if a.app.GetFocus() == a.treeView.TreeView {
		a.app.SetFocus(a.fileView.fileList)
	} else {
		a.app.SetFocus(a.treeView.TreeView)
	}
}

// showMessage exibe uma mensagem na barra de status
func (a *App) showMessage(msg string) {
	a.statusBar.SetStatus(msg)
}

// showError exibe uma mensagem de erro na barra de status
func (a *App) showError(msg string) {
	a.statusBar.SetError(msg)
}

// showHelp exibe a tela de ajuda
func (a *App) showHelp() {
	helpView := NewHelpView(a)
	a.pages.AddAndSwitchToPage("help", helpView.helpView, true)
}

// addToHistory adiciona um diretório ao histórico
func (a *App) addToHistory(dir string) {
	// Se já estamos no final do histórico, adicionar novo item
	if a.historyPos == len(a.history)-1 {
		a.history = append(a.history, dir)
		a.historyPos = len(a.history) - 1
	} else {
		// Se não estamos no final, truncar o histórico e adicionar novo item
		a.history = a.history[:a.historyPos+1]
		a.history = append(a.history, dir)
		a.historyPos = len(a.history) - 1
	}
}

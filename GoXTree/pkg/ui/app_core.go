package ui

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/peder1981/GoXTree/pkg/utils"
	"github.com/rivo/tview"
	"github.com/sergi/go-diff/diffmatchpatch"
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
		app:           tview.NewApplication(),
		pages:         tview.NewPages(),
		mainLayout:    tview.NewFlex(),
		history:       make([]string, 0),
		historyPos:    -1,
		showHidden:    false,
		selectedFiles: make(map[string]bool),
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

	// Carregar configuração
	config, err := LoadConfig()
	if err != nil {
		// Se houver erro, usar tema retrô padrão
		ApplyRetroThemeToApp(app)
		app.showError(fmt.Sprintf("Erro ao carregar configuração: %s. Usando tema padrão.", err))
	} else {
		// Aplicar tema da configuração
		if err := ApplyTheme(app, config.Theme); err != nil {
			// Se houver erro no tema, usar tema retrô padrão
			ApplyRetroThemeToApp(app)
			app.showError(fmt.Sprintf("Erro ao aplicar tema: %s. Usando tema padrão.", err))
		}
		// Aplicar outras configurações
		app.showHidden = config.ShowHidden
	}

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
			// Sincronizar diretórios
			a.syncDirectories()
			return nil
		case tcell.KeyF10:
			a.app.Stop()
			return nil
		case tcell.KeyTab:
			a.toggleFocus()
			return nil
		case tcell.KeyEscape:
			// Verificar se estamos na tela principal ou em uma tela de diálogo
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
			} else if a.pages.HasPage("themeSelector") {
				a.pages.RemovePage("themeSelector")
				return nil
			}
		}

		// Verificar combinações de teclas Ctrl+letra
		if event.Modifiers()&tcell.ModCtrl != 0 {
			switch event.Rune() {
			case 'a', 'A': // Ctrl+A: Selecionar todos os arquivos
				a.selectAll()
				return nil
			case 'd', 'D': // Ctrl+D: Desmarcar todos os arquivos
				a.unselectAll()
				return nil
			case 'h', 'H': // Ctrl+H: Alternar exibição de arquivos ocultos
				a.toggleHiddenFiles()
				return nil
			case 'q', 'Q': // Ctrl+Q: Sair
				a.confirmExit()
				return nil
			case 'c', 'C': // Ctrl+C: Copiar
				a.copySelectedFiles()
				return nil
			case 'x', 'X': // Ctrl+X: Recortar
				a.cutSelectedFiles()
				return nil
			case 'v', 'V': // Ctrl+V: Colar
				a.pasteFiles()
				return nil
			case 'f', 'F': // Ctrl+F: Buscar
				a.searchFiles()
				return nil
			case 'g', 'G': // Ctrl+G: Ir para diretório
				a.goToDirectory()
				return nil
			case 'r', 'R': // Ctrl+R: Atualizar
				a.refreshAll()
				return nil
			case 't', 'T': // Ctrl+T: Alternar tema
				a.showThemeSelector()
				return nil
			}
		}

		// Verificar combinações de teclas Alt+letra
		if event.Modifiers()&tcell.ModAlt != 0 {
			switch event.Rune() {
			case 's', 'S': // Alt+S: Ordenar
				a.toggleSortOrder()
				return nil
			case 'v', 'V': // Alt+V: Visualizar
				a.viewFile()
				return nil
			case 'e', 'E': // Alt+E: Editar
				a.editFile()
				return nil
			case 'i', 'I': // Alt+I: Informações
				a.showSystemInfo()
				return nil
			case 'c', 'C': // Alt+C: Comparar arquivos selecionados
				a.compareSelectedFiles()
				return nil
			}
		}

		// Passar o evento para o manipulador específico do componente com foco
		focusedPrimitive := a.app.GetFocus()
		switch focusedPrimitive {
		case a.fileView.fileList:
			return a.handleFileViewKeys(event)
		case a.treeView.TreeView:
			return a.handleTreeViewKeys(event)
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
	info := fmt.Sprintf(`[yellow]GoXTree - Gerenciador de Arquivos[white]

[yellow]Informações do Sistema:[white]
  [green]Sistema Operacional:[white] %s
  [green]Arquitetura:[white] %s
  [green]Número de CPUs:[white] %d
  [green]Diretório Atual:[white] %s
  [green]Diretório Home:[white] %s
  [green]Usuário:[white] %s

[yellow]Informações do GoXTree:[white]
  [green]Versão:[white] 1.0.0
  [green]Tema:[white] %s
  [green]Arquivos Ocultos:[white] %t
  [green]Arquivos Selecionados:[white] %d

[yellow]Estatísticas:[white]
  [green]Memória Utilizada:[white] %s
  [green]Tempo de Execução:[white] %s

Pressione [green]ESC[white] para fechar esta janela.`,
		runtime.GOOS,
		runtime.GOARCH,
		runtime.NumCPU(),
		a.currentDir,
		os.Getenv("HOME"),
		os.Getenv("USER"),
		"Retrô", // TODO: Obter tema atual
		a.showHidden,
		len(a.selectedFiles),
		"N/A", // TODO: Obter memória utilizada
		"N/A", // TODO: Obter tempo de execução
	)
	
	// Criar visualização de texto
	textView := tview.NewTextView().
		SetText(info).
		SetDynamicColors(true).
		SetRegions(true).
		SetScrollable(true).
		SetTitle(" Informações do Sistema ").
		SetTitleAlign(tview.AlignLeft).
		SetBorderColor(tcell.ColorBlue).
		SetBackgroundColor(tcell.ColorBlack)
	
	// Configurar cores
	textView.SetBorder(true)
	
	// Adicionar manipulador de teclas para sair
	textView.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEscape {
			a.pages.SwitchToPage("main")
			return nil
		}
		return event
	})
	
	// Adicionar à página e mostrar
	a.pages.AddAndSwitchToPage("systemInfo", textView, true)
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

// navigateTo navega para um diretório específico
func (a *App) navigateTo(dir string) {
	// Verificar se o diretório existe
	fileInfo, err := os.Stat(dir)
	if err != nil {
		a.showError(fmt.Sprintf("Erro ao acessar diretório: %s", err))
		return
	}
	
	// Verificar se é um diretório
	if !fileInfo.IsDir() {
		a.showError(fmt.Sprintf("%s não é um diretório", dir))
		return
	}
	
	// Navegar para o diretório
	a.currentDir = dir
	a.treeView.LoadTree(dir)
	a.fileView.SetCurrentDir(dir)
	a.addToHistory(dir)
	a.statusBar.UpdateStatus(dir)
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
		a.selectAll()
		return nil
	case tcell.KeyCtrlD:
		// Desmarcar todos os arquivos
		a.unselectAll()
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
		a.toggleHiddenFiles()
		return nil
	case tcell.KeyCtrlR:
		// Atualizar visualização
		a.refreshAll()
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
		a.syncDirectories()
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
	// Verificar teclas de navegação
	switch event.Key() {
	case tcell.KeyUp:
		// Mover para cima
		return event
	case tcell.KeyDown:
		// Mover para baixo
		return event
	case tcell.KeyPgUp:
		// Página acima
		return event
	case tcell.KeyPgDn:
		// Página abaixo
		return event
	case tcell.KeyHome:
		// Ir para o início
		a.fileView.fileList.Select(1, 0)
		return nil
	case tcell.KeyEnd:
		// Ir para o fim
		a.fileView.fileList.Select(a.fileView.fileList.GetRowCount()-1, 0)
		return nil
	}

	// Verificar teclas de letra
	if event.Key() == tcell.KeyRune {
		switch event.Rune() {
		case 'h':
			// Alternar exibição de arquivos ocultos
			a.toggleHidden()
			return nil
		case 'r':
			// Atualizar visualização
			a.refreshView()
			return nil
		case ' ':
			// Selecionar arquivo atual e mover para o próximo (para teclados que não enviam KeySpace)
			a.toggleSelectionWithSpace()
			return nil
		}
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
	// Criar modal de erro
	modal := tview.NewModal().
		SetText(message).
		AddButtons([]string{"OK"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			a.pages.RemovePage("error")
		})
	
	// Configurar cores
	modal.SetBackgroundColor(ColorBackground)
	modal.SetTextColor(ColorError)
	modal.SetButtonBackgroundColor(ColorBorder)
	modal.SetButtonTextColor(ColorText)
	
	// Adicionar à página e mostrar
	a.pages.AddPage("error", modal, true, true)
	a.app.SetFocus(modal)
}

// showMessage exibe uma mensagem informativa
func (a *App) showMessage(message string) {
	// Criar modal de mensagem
	modal := tview.NewModal().
		SetText(message).
		AddButtons([]string{"OK"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			a.pages.RemovePage("message")
		})
	
	// Configurar cores
	modal.SetBackgroundColor(ColorBackground)
	modal.SetTextColor(ColorText)
	modal.SetButtonBackgroundColor(ColorBorder)
	modal.SetButtonTextColor(ColorText)
	
	// Adicionar à página e mostrar
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
	files, err := utils.ListFiles(a.currentDir, a.showHidden)
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
	a.showInputDialog("Ir para Diretório", "", func(path string) {
		// Verificar se o diretório existe
		fileInfo, err := os.Stat(path)
		if err != nil {
			a.showError(fmt.Sprintf("Erro ao acessar diretório: %s", err))
			return
		}
		
		// Verificar se é um diretório
		if !fileInfo.IsDir() {
			a.showError(fmt.Sprintf("%s não é um diretório", path))
			return
		}
		
		// Navegar para o diretório
		a.navigateTo(path)
	})
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

// renameFile abre o diálogo para renomear um arquivo
func (a *App) renameFile() {
	// Obter arquivo selecionado
	selectedFile := a.fileView.GetSelectedFile()
	if selectedFile == "" {
		a.showError("Nenhum arquivo selecionado")
		return
	}
	
	// Obter caminho completo
	oldPath := filepath.Join(a.currentDir, selectedFile)
	
	// Verificar se o arquivo existe
	if _, err := os.Stat(oldPath); err != nil {
		a.showError(fmt.Sprintf("Erro ao acessar arquivo: %v", err))
		return
	}
	
	// Exibir diálogo para novo nome
	a.showInputDialogWithValue("Renomear", selectedFile, func(newName string) {
		if newName == "" || newName == selectedFile {
			return
		}
		
		// Verificar se o nome contém caracteres inválidos
		if strings.ContainsAny(newName, "\\/:*?\"<>|") {
			a.showError("O nome contém caracteres inválidos")
			return
		}
		
		// Criar caminho completo para o novo nome
		newPath := filepath.Join(a.currentDir, newName)
		
		// Verificar se já existe um arquivo com o novo nome
		if _, err := os.Stat(newPath); err == nil {
			a.showConfirmDialog("Confirmar substituição", fmt.Sprintf("Já existe um arquivo ou diretório com o nome '%s'. Deseja substituí-lo?", newName), func(confirmed bool) {
				if confirmed {
					// Remover arquivo existente
					if err := os.RemoveAll(newPath); err != nil {
						a.showError(fmt.Sprintf("Erro ao remover arquivo existente: %v", err))
						return
					}
					
					// Renomear arquivo
					if err := os.Rename(oldPath, newPath); err != nil {
						a.showError(fmt.Sprintf("Erro ao renomear: %v", err))
						return
					}
					
					a.refreshCurrentDir()
					a.showMessage(fmt.Sprintf("'%s' renomeado para '%s'", selectedFile, newName))
				}
			})
		} else {
			// Renomear arquivo
			if err := os.Rename(oldPath, newPath); err != nil {
				a.showError(fmt.Sprintf("Erro ao renomear: %v", err))
				return
			}
			
			a.refreshCurrentDir()
			a.showMessage(fmt.Sprintf("'%s' renomeado para '%s'", selectedFile, newName))
		}
	})
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

// showHelp exibe a tela de ajuda
func (a *App) showHelp() {
	helpView := NewHelpView(a)
	helpView.Show()
}

// showThemeSelector exibe o seletor de temas
func (a *App) showThemeSelector() {
	// Criar lista de temas
	list := tview.NewList().
		AddItem("Retrô", "Tema clássico estilo DOS", 'r', func() {
			a.changeTheme("retro")
		}).
		AddItem("Moderno", "Tema moderno com ícones Unicode", 'm', func() {
			a.changeTheme("modern")
		}).
		AddItem("Escuro", "Tema escuro para ambientes com pouca luz", 'e', func() {
			a.changeTheme("dark")
		}).
		AddItem("Claro", "Tema claro para ambientes bem iluminados", 'c', func() {
			a.changeTheme("light")
		}).
		AddItem("Cancelar", "Voltar sem alterar o tema", 'x', func() {
			a.pages.RemovePage("themeSelector")
		})

	// Configurar lista
	list.SetBorder(true).
		SetTitle(" Selecionar Tema ").
		SetTitleAlign(tview.AlignCenter).
		SetBorderPadding(1, 1, 1, 1)

	// Criar modal
	modal := tview.NewFlex().
		AddItem(nil, 0, 1, false).
		AddItem(tview.NewFlex().
			SetDirection(tview.FlexRow).
			AddItem(nil, 0, 1, false).
			AddItem(list, 10, 1, true).
			AddItem(nil, 0, 1, false), 40, 1, true).
		AddItem(nil, 0, 1, false)

	// Adicionar à página
	a.pages.AddPage("themeSelector", modal, true, true)
	a.app.SetFocus(list)
}

// changeTheme altera o tema da aplicação
func (a *App) changeTheme(themeName string) {
	// Aplicar tema
	if err := ApplyTheme(a, themeName); err != nil {
		a.showError(fmt.Sprintf("Erro ao aplicar tema: %s", err))
		return
	}

	// Atualizar configuração
	config, err := LoadConfig()
	if err != nil {
		a.showError(fmt.Sprintf("Erro ao carregar configuração: %s", err))
	} else {
		config.Theme = themeName
		if err := SaveConfig(config); err != nil {
			a.showError(fmt.Sprintf("Erro ao salvar configuração: %s", err))
		}
	}

	// Atualizar visualizações
	a.refreshAll()
	a.pages.RemovePage("themeSelector")
	a.showMessage(fmt.Sprintf("Tema '%s' aplicado com sucesso", themeName))
}

// toggleSortOrder alterna a ordem de classificação dos arquivos
func (a *App) toggleSortOrder() {
	// Implementar alternância de ordem de classificação
	a.showMessage("Alternando ordem de classificação (não implementado)")
}

// showFileInfo exibe informações detalhadas sobre o arquivo selecionado
func (a *App) showFileInfo() {
	// Implementar exibição de informações do arquivo
	a.showMessage("Exibindo informações do arquivo (não implementado)")
}

// editFile abre o arquivo selecionado para edição
func (a *App) editFile() {
	// Implementar edição de arquivo
	a.showMessage("Editando arquivo (não implementado)")
}

// compareSelectedFiles compara os arquivos selecionados
func (a *App) compareSelectedFiles() {
	// Verificar se há exatamente dois arquivos selecionados
	if len(a.selectedFiles) != 2 {
		a.showError("Selecione exatamente dois arquivos para comparar")
		return
	}

	// Obter os nomes dos arquivos selecionados
	var files []string
	for file := range a.selectedFiles {
		files = append(files, file)
	}

	// Verificar se ambos são arquivos (não diretórios)
	file1Info, err := os.Stat(files[0])
	if err != nil {
		a.showError(fmt.Sprintf("Erro ao acessar arquivo: %s", err))
		return
	}
	file2Info, err := os.Stat(files[1])
	if err != nil {
		a.showError(fmt.Sprintf("Erro ao acessar arquivo: %s", err))
		return
	}
	if file1Info.IsDir() || file2Info.IsDir() {
		a.showError("Não é possível comparar diretórios")
		return
	}

	// Implementar comparação de arquivos
	a.compareFiles(files[0], files[1])
}

// compareFiles compara dois arquivos e exibe as diferenças
func (a *App) compareFiles(file1, file2 string) {
	// Ler conteúdo dos arquivos
	content1, err := os.ReadFile(file1)
	if err != nil {
		a.showError(fmt.Sprintf("Erro ao ler arquivo %s: %v", file1, err))
		return
	}
	
	content2, err := os.ReadFile(file2)
	if err != nil {
		a.showError(fmt.Sprintf("Erro ao ler arquivo %s: %v", file2, err))
		return
	}

	// Criar diff
	dmp := diffmatchpatch.New()
	diffs := dmp.DiffMain(string(content1), string(content2), false)
	diffText := dmp.DiffPrettyText(diffs)

	// Criar visualizador de texto
	textView := tview.NewTextView().
		SetDynamicColors(true).
		SetRegions(true).
		SetScrollable(true).
		SetText(diffText)

	// Configurar visualizador
	textView.SetBorder(true).
		SetTitle(fmt.Sprintf(" Comparando: %s <-> %s ", filepath.Base(file1), filepath.Base(file2))).
		SetTitleAlign(tview.AlignCenter)

	// Criar modal
	modal := tview.NewFlex().
		AddItem(nil, 0, 1, false).
		AddItem(tview.NewFlex().
			SetDirection(tview.FlexRow).
			AddItem(nil, 0, 1, false).
			AddItem(textView, 0, 1, true).
			AddItem(nil, 0, 1, false), 0, 1, true).
		AddItem(nil, 0, 1, false)

	// Adicionar à página
	a.pages.AddPage("compare", modal, true, true)
	a.app.SetFocus(textView)

	// Adicionar manipulador de teclas para fechar
	textView.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEscape || event.Key() == tcell.KeyEnter {
			a.pages.RemovePage("compare")
			return nil
		}
		return event
	})
}

// addToHistory adiciona um diretório ao histórico de navegação
func (a *App) addToHistory(dir string) {
	// Verificar se o diretório já está no final do histórico
	if len(a.history) > 0 && a.history[len(a.history)-1] == dir {
		return
	}

	// Adicionar diretório ao histórico
	a.history = append(a.history, dir)
	a.historyPos = len(a.history) - 1
}

// sortByName ordena os arquivos por nome
func (a *App) sortByName() {
	a.fileView.SetSortBy("name")
	a.refreshCurrentDir()
	a.showMessage("Arquivos ordenados por nome")
}

// sortByDate ordena os arquivos por data
func (a *App) sortByDate() {
	a.fileView.SetSortBy("date")
	a.refreshCurrentDir()
	a.showMessage("Arquivos ordenados por data")
}

// sortBySize ordena os arquivos por tamanho
func (a *App) sortBySize() {
	a.fileView.SetSortBy("size")
	a.refreshCurrentDir()
	a.showMessage("Arquivos ordenados por tamanho")
}

// showPreferences exibe as preferências do aplicativo
func (a *App) showPreferences() {
	form := tview.NewForm()
	form.SetTitle("Preferências").
		SetTitleAlign(tview.AlignCenter).
		SetBorder(true)
	
	showHidden := a.showHidden
	
	form.AddCheckbox("Mostrar arquivos ocultos", showHidden, func(checked bool) {
		showHidden = checked
	})
	
	form.AddButton("Salvar", func() {
		a.showHidden = showHidden
		a.refreshCurrentDir()
		a.pages.RemovePage("preferences")
		a.showMessage("Preferências salvas")
	})
	
	form.AddButton("Cancelar", func() {
		a.pages.RemovePage("preferences")
	})
	
	form.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEscape {
			a.pages.RemovePage("preferences")
			return nil
		}
		return event
	})
	
	a.pages.AddPage("preferences", tview.NewFlex().
		AddItem(nil, 0, 1, false).
		AddItem(tview.NewFlex().
			SetDirection(tview.FlexRow).
			AddItem(nil, 0, 1, false).
			AddItem(form, 10, 1, true).
			AddItem(nil, 0, 1, false), 50, 1, true).
		AddItem(nil, 0, 1, false), true, true)
	a.app.SetFocus(form)
}

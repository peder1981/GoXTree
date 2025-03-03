package ui

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// MenuBar representa a barra de menu
type MenuBar struct {
	app     *App
	menuBar *tview.TextView
}

// NewMenuBar cria uma nova barra de menu
func NewMenuBar(app *App) *MenuBar {
	menuBar := &MenuBar{
		app: app,
	}

	// Criar barra de menu
	menuBar.menuBar = tview.NewTextView()
	menuBar.menuBar.SetTextColor(tcell.ColorWhite)
	menuBar.menuBar.SetBackgroundColor(tcell.ColorBlue)

	// Atualizar menu
	menuBar.updateMenu()

	return menuBar
}

// updateMenu atualiza o texto da barra de menu
func (m *MenuBar) updateMenu() {
	// Definir texto do menu
	menuText := " F1:Ajuda | F2:Menu | F3:Buscar Simples | F4:Buscar Avançada | F5:Copiar | F6:Mover | F7:Criar Dir | F8:Excluir | F9:Sincronizar | F10:Sair "

	// Configurar texto
	m.menuBar.SetText(menuText)
}

// handleMenuKey manipula teclas de menu
func (m *MenuBar) handleMenuKey(key tcell.Key) bool {
	switch key {
	case tcell.KeyF1:
		m.app.showHelp()
		return true
	case tcell.KeyF2:
		m.app.showMainMenu()
		return true
	case tcell.KeyF3:
		m.app.showSimpleSearchDialog()
		return true
	case tcell.KeyF4:
		m.app.showAdvancedSearchDialog()
		return true
	case tcell.KeyF5:
		m.app.copyFile()
		return true
	case tcell.KeyF6:
		m.app.moveFile()
		return true
	case tcell.KeyF7:
		m.app.showCreateMenu()
		return true
	case tcell.KeyF8:
		m.app.deleteFile()
		return true
	case tcell.KeyF9:
		m.app.syncDirectories()
		return true
	case tcell.KeyF10:
		m.app.confirmExit()
		return true
	}
	return false
}

// showMainMenu exibe o menu principal
func (a *App) showMainMenu() {
	menu := tview.NewList().
		AddItem("Arquivo", "Operações de arquivo", 'a', nil).
		AddItem("Editar", "Operações de edição", 'e', nil).
		AddItem("Visualizar", "Opções de visualização", 'v', nil).
		AddItem("Ferramentas", "Ferramentas úteis", 't', nil).
		AddItem("Configurações", "Configurar aplicação", 'c', nil).
		AddItem("Ajuda", "Exibir ajuda", 'h', func() {
			a.pages.RemovePage("mainMenu")
			a.showHelpDialog()
		}).
		AddItem("Sair", "Sair da aplicação", 's', func() {
			a.pages.RemovePage("mainMenu")
			a.app.Stop()
		})

	menu.SetBorder(true).
		SetTitle("Menu Principal").
		SetTitleAlign(tview.AlignCenter)

	menu.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEscape {
			a.pages.RemovePage("mainMenu")
			return nil
		}
		return event
	})

	a.pages.AddPage("mainMenu", menu, true, true)
	a.app.SetFocus(menu)
}

// showCreateMenu exibe o menu de criação
func (a *App) showCreateMenu() {
	menu := tview.NewList().
		AddItem("Arquivo", "Criar novo arquivo", 'a', func() {
			a.pages.RemovePage("createMenu")
			a.createFile()
		}).
		AddItem("Diretório", "Criar novo diretório", 'd', func() {
			a.pages.RemovePage("createMenu")
			a.createDirectory()
		})

	menu.SetBorder(true).
		SetTitle("Criar").
		SetTitleAlign(tview.AlignCenter)

	menu.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEscape {
			a.pages.RemovePage("createMenu")
			return nil
		}
		return event
	})

	a.pages.AddPage("createMenu", menu, true, true)
	a.app.SetFocus(menu)
}

// moveFile move um arquivo ou diretório
func (a *App) moveFile() {
	selectedFile := a.getSelectedFile()
	if selectedFile == "" {
		a.showError("Nenhum arquivo selecionado")
		return
	}

	a.showInputDialog("Mover para", "Diretório de destino:", func(destDir string) {
		if destDir == "" {
			return
		}

		// Expandir caminho
		if strings.HasPrefix(destDir, "~") {
			homeDir, err := os.UserHomeDir()
			if err == nil {
				destDir = filepath.Join(homeDir, destDir[1:])
			}
		}

		// Verificar se é um caminho relativo
		if !filepath.IsAbs(destDir) {
			destDir = filepath.Join(a.currentDir, destDir)
		}

		// Verificar se o diretório de destino existe
		destInfo, err := os.Stat(destDir)
		if err != nil {
			a.showError(fmt.Sprintf("Diretório de destino não encontrado: %s", destDir))
			return
		}

		// Verificar se o destino é um diretório
		if !destInfo.IsDir() {
			a.showError(fmt.Sprintf("O destino não é um diretório: %s", destDir))
			return
		}

		// Caminho de destino completo
		fileName := filepath.Base(selectedFile)
		destPath := filepath.Join(destDir, fileName)

		// Verificar se o destino já existe
		if _, err := os.Stat(destPath); err == nil {
			a.showConfirmDialog("Substituir", fmt.Sprintf("'%s' já existe. Deseja substituir?", fileName), func(confirmed bool) {
				if !confirmed {
					return
				}

				// Mover arquivo ou diretório
				err := os.Rename(selectedFile, destPath)
				if err != nil {
					a.showError(fmt.Sprintf("Erro ao mover: %v", err))
					return
				}

				a.refreshCurrentDir()
				a.showMessage(fmt.Sprintf("'%s' movido com sucesso", fileName))
			})
		} else {
			// Mover arquivo ou diretório
			err := os.Rename(selectedFile, destPath)
			if err != nil {
				a.showError(fmt.Sprintf("Erro ao mover: %v", err))
				return
			}

			a.refreshCurrentDir()
			a.showMessage(fmt.Sprintf("'%s' movido com sucesso", fileName))
		}
	})
}

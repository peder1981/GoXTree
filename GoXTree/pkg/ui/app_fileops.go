package ui

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/peder1981/GoXTree/pkg/utils"
	"github.com/rivo/tview"
)

// showToolsMenu exibe o menu de ferramentas
func (a *App) showToolsMenu() {
	// Criar menu
	menu := tview.NewList()
	menu.SetTitle("Ferramentas").
		SetTitleAlign(tview.AlignCenter).
		SetBorder(true)
	
	// Adicionar itens ao menu
	menu.AddItem("Buscar Arquivo", "Busca por nome de arquivo", 'b', func() {
		a.pages.RemovePage("toolsMenu")
		a.showSearchDialog()
	})
	
	menu.AddItem("Comparar Arquivos", "Compara dois arquivos", 'c', func() {
		a.pages.RemovePage("toolsMenu")
		a.showCompareDialog()
	})
	
	menu.AddItem("Sincronizar Diretórios", "Sincroniza dois diretórios", 's', func() {
		a.pages.RemovePage("toolsMenu")
		a.syncDirectories()
	})
	
	menu.AddItem("Voltar", "Volta ao gerenciador de arquivos", 'v', func() {
		a.pages.RemovePage("toolsMenu")
	})
	
	// Configurar manipulador de eventos
	menu.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEscape {
			a.pages.RemovePage("toolsMenu")
			return nil
		}
		return event
	})
	
	// Exibir menu
	a.pages.AddPage("toolsMenu", menu, true, true)
	a.app.SetFocus(menu)
}

// getHomeDir retorna o diretório home do usuário
func (a *App) getHomeDir() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return homeDir, nil
}

// createFile cria um novo arquivo
func (a *App) createFile() {
	a.showInputDialog("Novo Arquivo", "Nome:", func(fileName string) {
		if fileName == "" {
			return
		}
		
		// Verificar se o nome contém caracteres inválidos
		if strings.ContainsAny(fileName, "\\/:*?\"<>|") {
			a.showError("O nome contém caracteres inválidos")
			return
		}
		
		// Criar caminho completo
		filePath := filepath.Join(a.currentDir, fileName)
		
		// Verificar se o arquivo já existe
		if _, err := os.Stat(filePath); err == nil {
			a.showError(fmt.Sprintf("Já existe um arquivo ou diretório com o nome '%s'", fileName))
			return
		}
		
		// Criar arquivo
		file, err := os.Create(filePath)
		if err != nil {
			a.showError(fmt.Sprintf("Erro ao criar arquivo: %v", err))
			return
		}
		file.Close()
		
		a.refreshCurrentDir()
		a.showMessage(fmt.Sprintf("Arquivo '%s' criado com sucesso", fileName))
	})
}

// copyFile copia um arquivo ou diretório
func (a *App) copyFile() {
	// Obter arquivo selecionado
	selectedFile := a.fileView.GetSelectedFile()
	if selectedFile == "" {
		a.showMessage("Nenhum arquivo selecionado")
		return
	}
	
	// Obter caminho completo
	srcPath := filepath.Join(a.currentDir, selectedFile)
	
	// Verificar se é um diretório
	fileInfo, err := os.Stat(srcPath)
	if err != nil {
		a.showMessage(fmt.Sprintf("Erro ao acessar arquivo: %v", err))
		return
	}
	
	isDir := fileInfo.IsDir()
	
	// Perguntar para onde copiar
	a.showInputDialogWithValue("Copiar para", a.currentDir, func(destDir string) {
		if destDir == "" {
			return
		}
		
		// Verificar se o destino existe
		destInfo, err := os.Stat(destDir)
		if err != nil {
			a.showError(fmt.Sprintf("Destino inválido: %v", err))
			return
		}
		
		// Verificar se o destino é um diretório
		if !destInfo.IsDir() {
			a.showError("O destino deve ser um diretório")
			return
		}
		
		// Construir caminho de destino
		destPath := filepath.Join(destDir, selectedFile)
		
		// Verificar se o destino já existe
		if _, err := os.Stat(destPath); err == nil {
			// Perguntar se deseja sobrescrever
			a.showConfirmDialog(fmt.Sprintf("Confirmar sobrescrita"), fmt.Sprintf("'%s' já existe. Sobrescrever?", selectedFile), func(confirmed bool) {
				if confirmed {
					a.doCopy(srcPath, destPath, isDir)
				}
			})
		} else {
			// Copiar diretamente
			a.doCopy(srcPath, destPath, isDir)
		}
	})
}

// doCopy realiza a cópia de um arquivo ou diretório
func (a *App) doCopy(src, dest string, isDir bool) {
	var err error
	
	if isDir {
		// Copiar diretório
		err = utils.CopyDir(src, dest)
	} else {
		// Copiar arquivo
		err = utils.CopyFile(src, dest)
	}
	
	if err != nil {
		a.showError(fmt.Sprintf("Erro ao copiar: %v", err))
		return
	}
	
	a.refreshCurrentDir()
	a.showMessage("Cópia concluída com sucesso")
}

// getSelectedFile retorna o caminho completo do arquivo selecionado
func (a *App) getSelectedFile() string {
	selectedFile := a.fileView.GetSelectedFile()
	if selectedFile == "" {
		return ""
	}
	
	return filepath.Join(a.currentDir, selectedFile)
}

// refreshCurrentDir atualiza a visualização do diretório atual
func (a *App) refreshCurrentDir() {
	a.refreshView()
}

// showSyncResults exibe os resultados da sincronização
func (a *App) showSyncResults(actions []utils.SyncAction) {
	// Criar texto com resultados
	var text string
	var copied, deleted, skipped int
	
	for _, action := range actions {
		switch action.Action {
		case "copy":
			copied++
			text += fmt.Sprintf("[green]+ Copiado:[white] %s -> %s\n", action.SourcePath, action.DestPath)
		case "delete":
			deleted++
			text += fmt.Sprintf("[red]- Excluído:[white] %s\n", action.DestPath)
		case "skip":
			skipped++
			text += fmt.Sprintf("[yellow]~ Ignorado:[white] %s (%s)\n", action.DestPath, action.Reason)
		}
	}
	
	// Adicionar resumo
	text += fmt.Sprintf("\n[blue]Resumo:[white]\n")
	text += fmt.Sprintf("  [green]Arquivos copiados:[white] %d\n", copied)
	text += fmt.Sprintf("  [red]Arquivos excluídos:[white] %d\n", deleted)
	text += fmt.Sprintf("  [yellow]Arquivos ignorados:[white] %d\n", skipped)
	text += fmt.Sprintf("  [blue]Total de operações:[white] %d\n", len(actions))
	
	// Criar visualizador de texto
	textView := tview.NewTextView().
		SetDynamicColors(true).
		SetText(text).
		SetScrollable(true).
		SetWordWrap(true)
	
	// Configurar borda
	textView.SetBorder(true).
		SetTitle(" Resultados da Sincronização ").
		SetTitleAlign(tview.AlignLeft)
	
	// Configurar manipulador de teclas
	textView.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEscape || event.Key() == tcell.KeyEnter {
			a.pages.RemovePage("syncResults")
			return nil
		}
		return event
	})
	
	// Exibir resultados
	a.pages.AddPage("syncResults", a.modal(textView, 70, 20), true, true)
}

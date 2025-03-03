package ui

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/peder1981/GoXTree/pkg/utils"
	"github.com/rivo/tview"

	"github.com/gdamore/tcell/v2"
)

// selectAll seleciona todos os arquivos no diretório atual
func (a *App) selectAll() {
	// Inicializar mapa de seleção se não existir
	if a.selectedFiles == nil {
		a.selectedFiles = make(map[string]bool)
	}

	// Obter lista de arquivos no diretório atual
	files, err := os.ReadDir(a.currentDir)
	if err != nil {
		a.showError("Erro ao listar arquivos: " + err.Error())
		return
	}

	// Adicionar todos os arquivos à seleção
	for _, file := range files {
		if !file.IsDir() {
			filePath := filepath.Join(a.currentDir, file.Name())
			a.selectedFiles[filePath] = true
		}
	}

	// Atualizar visualização
	a.refreshFileView()

	// Atualizar barra de status
	a.statusBar.SetStatus(fmt.Sprintf("%d arquivos selecionados", len(a.selectedFiles)))
}

// unselectAll remove a seleção de todos os arquivos
func (a *App) unselectAll() {
	// Limpar seleção atual
	a.selectedFiles = make(map[string]bool)

	// Atualizar visualização
	a.refreshFileView()

	// Atualizar barra de status
	a.statusBar.SetStatus("Seleção removida")
}

// toggleSelection alterna a seleção do arquivo atual
func (a *App) toggleSelection() {
	// Obter arquivo selecionado
	row, _ := a.fileView.fileList.GetSelection()
	if row <= 1 { // Cabeçalho ou diretório pai
		return
	}

	// Obter nome do arquivo
	fileName := a.fileView.fileList.GetCell(row, 0).Text
	filePath := filepath.Join(a.currentDir, fileName)

	// Verificar se o arquivo já está selecionado
	if _, ok := a.selectedFiles[filePath]; ok {
		// Remover da seleção
		delete(a.selectedFiles, filePath)
	} else {
		// Adicionar à seleção
		a.selectedFiles[filePath] = true
	}

	// Atualizar visualização
	a.refreshFileView()

	// Atualizar barra de status
	if len(a.selectedFiles) > 0 {
		a.statusBar.SetStatus(fmt.Sprintf("%d arquivos selecionados", len(a.selectedFiles)))
	} else {
		a.statusBar.SetStatus("Nenhum arquivo selecionado")
	}
}

// toggleSelectionWithSpace alterna a seleção do arquivo atual e move para o próximo
func (a *App) toggleSelectionWithSpace() {
	// Obter arquivo selecionado
	row, _ := a.fileView.fileList.GetSelection()
	if row <= 1 { // Cabeçalho ou diretório pai
		return
	}

	// Obter nome do arquivo
	fileName := a.fileView.fileList.GetCell(row, 0).Text
	filePath := filepath.Join(a.currentDir, fileName)

	// Verificar se o arquivo já está selecionado
	if _, ok := a.selectedFiles[filePath]; ok {
		// Remover da seleção
		delete(a.selectedFiles, filePath)
	} else {
		// Adicionar à seleção
		a.selectedFiles[filePath] = true
	}

	// Mover para o próximo item
	a.fileView.fileList.Select(row+1, 0)

	// Atualizar visualização
	a.refreshFileView()

	// Atualizar barra de status
	if len(a.selectedFiles) > 0 {
		a.statusBar.SetStatus(fmt.Sprintf("%d arquivos selecionados", len(a.selectedFiles)))
	} else {
		a.statusBar.SetStatus("Nenhum arquivo selecionado")
	}
}

// viewCurrentFile visualiza o conteúdo do arquivo atual
func (a *App) viewCurrentFile() {
	// Obter arquivo selecionado
	row, _ := a.fileView.fileList.GetSelection()
	if row <= 1 { // Cabeçalho ou diretório pai
		return
	}

	// Obter nome do arquivo
	fileName := a.fileView.fileList.GetCell(row, 0).Text
	filePath := filepath.Join(a.currentDir, fileName)

	// Verificar se é um diretório
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		a.showError("Erro ao acessar arquivo: " + err.Error())
		return
	}

	if fileInfo.IsDir() {
		a.showMessage("Não é possível visualizar diretórios")
		return
	}

	// Ler conteúdo do arquivo
	content, err := os.ReadFile(filePath)
	if err != nil {
		a.showError("Erro ao ler arquivo: " + err.Error())
		return
	}

	// Criar visualizador de texto
	textView := tview.NewTextView().
		SetDynamicColors(true).
		SetRegions(true).
		SetScrollable(true).
		SetText(string(content)).
		SetTitle(fmt.Sprintf(" Visualizando: %s ", fileName)).
		SetTitleColor(ColorTitle).
		SetBorderColor(ColorBorder).
		SetBackgroundColor(ColorBackground)

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
	a.pages.AddAndSwitchToPage("view", textView, true)
}

// editCurrentFile abre o arquivo atual para edição
func (a *App) editCurrentFile() {
	// Obter arquivo selecionado
	selectedFile := a.fileView.GetSelectedFile()
	if selectedFile == "" {
		a.showMessage("Nenhum arquivo selecionado")
		return
	}

	// Obter caminho completo
	filePath := filepath.Join(a.currentDir, selectedFile)

	// Verificar se é um arquivo
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		a.showMessage(fmt.Sprintf("Erro ao acessar arquivo: %v", err))
		return
	}

	if fileInfo.IsDir() {
		a.showMessage("Não é possível editar um diretório")
		return
	}

	// Ler conteúdo do arquivo
	content, err := os.ReadFile(filePath)
	if err != nil {
		a.showMessage(fmt.Sprintf("Erro ao ler arquivo: %v", err))
		return
	}

	// Criar página de edição
	editorPage := tview.NewFlex().SetDirection(tview.FlexRow)

	// Criar área de texto
	textArea := tview.NewTextView().
		SetText(string(content)).
		SetScrollable(true).
		SetWordWrap(true).
		SetChangedFunc(func() {
			a.app.Draw()
		})

	// Configurar borda
	textArea.SetBorder(true).
		SetTitle(fmt.Sprintf(" Editando: %s ", selectedFile)).
		SetTitleAlign(tview.AlignLeft)

	// Criar botões
	saveButton := tview.NewButton("Salvar").SetSelectedFunc(func() {
		// Obter conteúdo editado
		newContent := textArea.GetText(true)

		// Salvar arquivo
		err := os.WriteFile(filePath, []byte(newContent), 0644)
		if err != nil {
			a.showError(fmt.Sprintf("Erro ao salvar arquivo: %v", err))
			return
		}

		// Voltar à visualização principal
		a.pages.RemovePage("editor")
		a.refreshView()
		a.showMessage("Arquivo salvo com sucesso")
	})

	cancelButton := tview.NewButton("Cancelar").SetSelectedFunc(func() {
		// Voltar à visualização principal
		a.pages.RemovePage("editor")
	})

	// Criar layout de botões
	buttons := tview.NewFlex().
		AddItem(saveButton, 0, 1, false).
		AddItem(nil, 1, 0, false).
		AddItem(cancelButton, 0, 1, false)

	// Adicionar componentes à página
	editorPage.AddItem(textArea, 0, 1, true)
	editorPage.AddItem(buttons, 1, 0, false)

	// Adicionar página ao aplicativo
	a.pages.AddPage("editor", editorPage, true, true)

	// Configurar manipulador de teclas
	editorPage.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEscape {
			a.pages.RemovePage("editor")
			return nil
		}
		return event
	})
}

// copySelectedFiles copia os arquivos selecionados para a área de transferência
func (a *App) copySelectedFiles() {
	// Verificar se há arquivos selecionados
	if len(a.selectedFiles) == 0 {
		a.showMessage("Nenhum arquivo selecionado")
		return
	}

	// Definir operação como cópia
	a.clipboard = "copy"

	// Atualizar barra de status
	a.statusBar.SetStatus(fmt.Sprintf("%d arquivos copiados para a área de transferência", len(a.selectedFiles)))
}

// cutSelectedFiles recorta os arquivos selecionados para a área de transferência
func (a *App) cutSelectedFiles() {
	// Verificar se há arquivos selecionados
	if len(a.selectedFiles) == 0 {
		a.showMessage("Nenhum arquivo selecionado")
		return
	}

	// Definir operação como recorte
	a.clipboard = "cut"

	// Atualizar barra de status
	a.statusBar.SetStatus(fmt.Sprintf("%d arquivos recortados para a área de transferência", len(a.selectedFiles)))
}

// pasteFiles cola os arquivos da área de transferência
func (a *App) pasteFiles() {
	// Verificar se há arquivos na área de transferência
	if len(a.selectedFiles) == 0 {
		a.showMessage("Nenhum arquivo na área de transferência")
		return
	}

	// Contar arquivos processados
	processed := 0

	// Processar cada arquivo selecionado
	for filePath := range a.selectedFiles {
		// Obter nome do arquivo
		fileName := filepath.Base(filePath)

		// Definir caminho de destino
		destPath := filepath.Join(a.currentDir, fileName)

		// Verificar se o arquivo já existe no destino
		if _, err := os.Stat(destPath); err == nil {
			// Perguntar se deseja sobrescrever
			a.showConfirmDialog("Confirmação", fmt.Sprintf("O arquivo %s já existe. Sobrescrever?", fileName), func(confirmed bool) {
				if !confirmed {
					return
				}

				// Continuar com a cópia
				a.doCopy(filePath, destPath, false)

				a.refreshView()
			})
		} else {
			// Copiar arquivo
			a.doCopy(filePath, destPath, false)

			a.refreshView()
		}

		processed++
	}

	// Limpar seleção
	a.selectedFiles = make(map[string]bool)

	// Atualizar visualização
	a.refreshFileView()

	// Atualizar barra de status
	a.statusBar.SetStatus(fmt.Sprintf("%d arquivos processados", processed))
}

// invertSelection inverte a seleção de arquivos
func (a *App) invertSelection() {
	// Inicializar mapa de seleção se não existir
	if a.selectedFiles == nil {
		a.selectedFiles = make(map[string]bool)
	}

	// Obter lista de arquivos no diretório atual
	files, err := os.ReadDir(a.currentDir)
	if err != nil {
		a.showError("Erro ao listar arquivos: " + err.Error())
		return
	}

	// Inverter seleção
	for _, file := range files {
		if !file.IsDir() {
			filePath := filepath.Join(a.currentDir, file.Name())
			a.selectedFiles[filePath] = !a.selectedFiles[filePath]
		}
	}

	// Atualizar visualização
	a.refreshFileView()
	a.statusBar.SetText(fmt.Sprintf("Selecionados %d arquivos", len(a.selectedFiles)))
}

// selectByPattern seleciona arquivos por padrão
func (a *App) selectByPattern() {
	a.showInputDialog("Selecionar por Padrão", "Padrão:", func(pattern string) {
		if pattern == "" {
			return
		}

		// Compilar expressão regular
		regex, err := regexp.Compile(pattern)
		if err != nil {
			a.showError(fmt.Sprintf("Padrão inválido: %v", err))
			return
		}

		// Obter lista de arquivos
		files, err := os.ReadDir(a.currentDir)
		if err != nil {
			a.showError("Erro ao ler diretório: " + err.Error())
			return
		}

		// Contar arquivos selecionados
		count := 0

		// Selecionar arquivos que correspondem ao padrão
		for _, file := range files {
			// Ignorar diretório pai
			if file.Name() == ".." {
				continue
			}

			// Verificar se corresponde ao padrão
			if regex.MatchString(file.Name()) {
				// Adicionar à seleção
				filePath := filepath.Join(a.currentDir, file.Name())
				a.selectedFiles[filePath] = true
				count++
			}
		}

		// Atualizar visualização
		a.refreshFileView()
		a.statusBar.SetText(fmt.Sprintf("Selecionados %d arquivos", count))
	})
}

// showSelectionMenu exibe o menu de seleção
func (a *App) showSelectionMenu() {
	// Criar menu
	menu := tview.NewList()
	menu.SetTitle("Seleção").
		SetTitleAlign(tview.AlignCenter).
		SetBorder(true)

	// Adicionar itens ao menu
	menu.AddItem("Selecionar Todos", "Seleciona todos os arquivos", 't', func() {
		a.pages.RemovePage("selectionMenu")
		a.selectAll()
	})

	menu.AddItem("Desmarcar Todos", "Remove todas as seleções", 'd', func() {
		a.pages.RemovePage("selectionMenu")
		a.unselectAll()
	})

	menu.AddItem("Inverter Seleção", "Inverte a seleção atual", 'i', func() {
		a.pages.RemovePage("selectionMenu")
		a.invertSelection()
	})

	menu.AddItem("Selecionar por Padrão", "Seleciona arquivos por padrão", 'p', func() {
		a.pages.RemovePage("selectionMenu")
		a.selectByPattern()
	})

	menu.AddItem("Voltar", "Volta ao gerenciador de arquivos", 'v', func() {
		a.pages.RemovePage("selectionMenu")
	})

	// Configurar manipulador de eventos
	menu.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEscape {
			a.pages.RemovePage("selectionMenu")
			return nil
		}
		return event
	})

	// Adicionar página
	a.pages.AddPage("selectionMenu", menu, true, true)
	a.app.SetFocus(menu)
}

// getSelectedFiles retorna a lista de arquivos selecionados
func (a *App) getSelectedFiles() []string {
	var files []string
	for file := range a.selectedFiles {
		files = append(files, file)
	}
	return files
}

// getSelectedFilesCount retorna o número de arquivos selecionados
func (a *App) getSelectedFilesCount() int {
	return len(a.selectedFiles)
}

// isSelected verifica se um arquivo está selecionado
func (a *App) isSelected(file string) bool {
	return a.selectedFiles[file]
}

// showSelectedFilesInfo exibe informações sobre os arquivos selecionados
func (a *App) showSelectedFilesInfo() {
	// Verificar se há arquivos selecionados
	if len(a.selectedFiles) == 0 {
		a.showMessage("Nenhum arquivo selecionado")
		return
	}

	// Calcular estatísticas
	var (
		totalSize  int64
		fileCount  int
		dirCount   int
		newestTime time.Time
		oldestTime time.Time
		newestFile string
		oldestFile string
	)

	// Processar arquivos selecionados
	for file := range a.selectedFiles {
		// Obter informações do arquivo
		info, err := os.Stat(file)
		if err != nil {
			continue
		}

		// Atualizar estatísticas
		if info.IsDir() {
			dirCount++
		} else {
			fileCount++
			totalSize += info.Size()
		}

		// Verificar data de modificação
		modTime := info.ModTime()
		if newestTime.IsZero() || modTime.After(newestTime) {
			newestTime = modTime
			newestFile = filepath.Base(file)
		}

		if oldestTime.IsZero() || modTime.Before(oldestTime) {
			oldestTime = modTime
			oldestFile = filepath.Base(file)
		}
	}

	// Criar texto de informações
	var text strings.Builder
	text.WriteString(fmt.Sprintf("Arquivos selecionados: %d\n", len(a.selectedFiles)))
	text.WriteString(fmt.Sprintf("Arquivos: %d\n", fileCount))
	text.WriteString(fmt.Sprintf("Diretórios: %d\n", dirCount))
	text.WriteString(fmt.Sprintf("Tamanho total: %s\n", utils.FormatFileSize(totalSize)))

	if newestFile != "" {
		text.WriteString(fmt.Sprintf("\nArquivo mais recente: %s\n", newestFile))
		text.WriteString(fmt.Sprintf("Data: %s\n", newestTime.Format("02/01/2006 15:04:05")))
	}

	if oldestFile != "" {
		text.WriteString(fmt.Sprintf("\nArquivo mais antigo: %s\n", oldestFile))
		text.WriteString(fmt.Sprintf("Data: %s\n", oldestTime.Format("02/01/2006 15:04:05")))
	}

	// Exibir informações
	a.showTextDialog("Informações da Seleção", text.String())
}

// showTextDialog exibe um diálogo com texto
func (a *App) showTextDialog(title, text string) {
	// Criar visualização de texto
	textView := tview.NewTextView()
	textView.SetText(text).
		SetScrollable(true).
		SetTitle(title).
		SetTitleAlign(tview.AlignCenter).
		SetBorder(true)

	// Configurar manipulador de eventos
	textView.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEscape || event.Key() == tcell.KeyEnter {
			a.pages.RemovePage("textDialog")
			return nil
		}
		return event
	})

	// Adicionar página
	a.pages.AddPage("textDialog", textView, true, true)
	a.app.SetFocus(textView)
}

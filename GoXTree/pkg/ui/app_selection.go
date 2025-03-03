package ui

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/sergi/go-diff/diffmatchpatch"
)

// selectAllFiles seleciona todos os arquivos no diretório atual
func (a *App) selectAllFiles() {
	// Limpar seleção atual
	a.fileView.selectedFiles = make([]string, 0)
	
	// Obter lista de arquivos no diretório atual
	files, err := os.ReadDir(a.currentDir)
	if err != nil {
		a.showError("Erro ao listar arquivos: " + err.Error())
		return
	}
	
	// Adicionar todos os arquivos à seleção
	for _, file := range files {
		if !file.IsDir() {
			a.fileView.selectedFiles = append(a.fileView.selectedFiles, file.Name())
		}
	}
	
	// Atualizar visualização
	a.updateFileList()
	
	// Atualizar barra de status
	a.statusBar.SetStatus(fmt.Sprintf("%d arquivos selecionados", len(a.fileView.selectedFiles)))
}

// unselectAllFiles remove a seleção de todos os arquivos
func (a *App) unselectAllFiles() {
	// Limpar seleção atual
	a.fileView.selectedFiles = make([]string, 0)
	
	// Atualizar visualização
	a.updateFileList()
	
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
	
	// Verificar se o arquivo já está selecionado
	isSelected := false
	for i, selected := range a.fileView.selectedFiles {
		if selected == fileName {
			// Remover da seleção
			a.fileView.selectedFiles = append(a.fileView.selectedFiles[:i], a.fileView.selectedFiles[i+1:]...)
			isSelected = true
			break
		}
	}
	
	// Se não estava selecionado, adicionar à seleção
	if !isSelected {
		a.fileView.selectedFiles = append(a.fileView.selectedFiles, fileName)
	}
	
	// Atualizar visualização
	a.updateFileList()
	
	// Atualizar barra de status
	if len(a.fileView.selectedFiles) > 0 {
		a.statusBar.SetStatus(fmt.Sprintf("%d arquivos selecionados", len(a.fileView.selectedFiles)))
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
	
	textView.SetTextColor(ColorText)
	
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
		a.showMessage("Não é possível editar diretórios")
		return
	}
	
	// Ler conteúdo do arquivo
	content, err := os.ReadFile(filePath)
	if err != nil {
		a.showError("Erro ao ler arquivo: " + err.Error())
		return
	}
	
	// Criar editor de texto
	textArea := tview.NewTextArea().
		SetText(string(content), true).
		SetTitle(fmt.Sprintf(" Editando: %s ", fileName)).
		SetTitleColor(ColorTitle).
		SetBorderColor(ColorBorder).
		SetBackgroundColor(ColorBackground)
	
	textArea.SetTextStyle(tcell.StyleDefault.Foreground(ColorText))
	
	textArea.SetBorder(true)
	
	// Criar botões
	saveButton := tview.NewButton("Salvar").SetSelectedFunc(func() {
		// Obter conteúdo editado
		newContent := textArea.GetText()
		
		// Salvar arquivo
		err := os.WriteFile(filePath, []byte(newContent), 0644)
		if err != nil {
			a.showError("Erro ao salvar arquivo: " + err.Error())
			return
		}
		
		a.showMessage("Arquivo salvo com sucesso")
		a.pages.SwitchToPage("main")
	})
	
	cancelButton := tview.NewButton("Cancelar").SetSelectedFunc(func() {
		a.pages.SwitchToPage("main")
	})
	
	// Configurar cores dos botões
	saveButton.SetLabelColor(ColorText)
	saveButton.SetBackgroundColor(ColorBackground)
	cancelButton.SetLabelColor(ColorText)
	cancelButton.SetBackgroundColor(ColorBackground)
	
	// Criar layout
	buttons := tview.NewFlex().
		AddItem(saveButton, 0, 1, true).
		AddItem(nil, 1, 0, false).
		AddItem(cancelButton, 0, 1, true)
	
	layout := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(textArea, 0, 1, true).
		AddItem(buttons, 1, 0, false)
	
	// Adicionar manipulador de teclas para sair
	layout.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEscape {
			a.pages.SwitchToPage("main")
			return nil
		}
		return event
	})
	
	// Adicionar à página e mostrar
	a.pages.AddAndSwitchToPage("edit", layout, true)
}

// compareSelectedFiles compara dois arquivos selecionados
func (a *App) compareSelectedFiles() {
	// Verificar se há exatamente dois arquivos selecionados
	if len(a.fileView.selectedFiles) != 2 {
		a.showMessage("Selecione exatamente dois arquivos para comparar")
		return
	}
	
	// Obter caminhos dos arquivos
	file1Path := filepath.Join(a.currentDir, a.fileView.selectedFiles[0])
	file2Path := filepath.Join(a.currentDir, a.fileView.selectedFiles[1])
	
	// Verificar se são diretórios
	file1Info, err := os.Stat(file1Path)
	if err != nil {
		a.showError("Erro ao acessar arquivo: " + err.Error())
		return
	}
	
	file2Info, err := os.Stat(file2Path)
	if err != nil {
		a.showError("Erro ao acessar arquivo: " + err.Error())
		return
	}
	
	if file1Info.IsDir() || file2Info.IsDir() {
		a.showMessage("Não é possível comparar diretórios")
		return
	}
	
	// Ler conteúdo dos arquivos
	content1, err := os.ReadFile(file1Path)
	if err != nil {
		a.showError("Erro ao ler arquivo: " + err.Error())
		return
	}
	
	content2, err := os.ReadFile(file2Path)
	if err != nil {
		a.showError("Erro ao ler arquivo: " + err.Error())
		return
	}
	
	// Comparar arquivos
	dmp := diffmatchpatch.New()
	diffs := dmp.DiffMain(string(content1), string(content2), false)
	
	// Criar visualizador de texto
	textView := tview.NewTextView().
		SetDynamicColors(true).
		SetRegions(true).
		SetScrollable(true).
		SetTitle(fmt.Sprintf(" Comparando: %s e %s ", a.fileView.selectedFiles[0], a.fileView.selectedFiles[1])).
		SetTitleColor(ColorTitle).
		SetBorderColor(ColorBorder).
		SetBackgroundColor(ColorBackground)
	
	textView.SetTextColor(ColorText)
	
	textView.SetBorder(true)
	
	// Formatar diferenças
	var formattedText strings.Builder
	for _, diff := range diffs {
		text := diff.Text
		switch diff.Type {
		case diffmatchpatch.DiffInsert:
			formattedText.WriteString("[green]+ " + text + "[white]")
		case diffmatchpatch.DiffDelete:
			formattedText.WriteString("[red]- " + text + "[white]")
		case diffmatchpatch.DiffEqual:
			formattedText.WriteString("  " + text)
		}
	}
	
	textView.SetText(formattedText.String())
	
	// Adicionar manipulador de teclas para sair
	textView.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEscape {
			a.pages.SwitchToPage("main")
			return nil
		}
		return event
	})
	
	// Adicionar à página e mostrar
	a.pages.AddAndSwitchToPage("compare", textView, true)
}

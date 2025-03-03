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

// editFile edita o arquivo selecionado
func (a *App) editFile() {
	// Obter arquivo selecionado
	selectedFile := a.fileView.GetSelectedFile()
	if selectedFile == "" {
		a.showMessage("Nenhum arquivo selecionado")
		return
	}

	// Verificar se é um diretório
	filePath := filepath.Join(a.currentDir, selectedFile)
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		a.showMessage(fmt.Sprintf("Erro ao acessar arquivo: %v", err))
		return
	}

	if fileInfo.IsDir() {
		a.showMessage("Não é possível editar um diretório")
		return
	}

	// Verificar tamanho do arquivo
	if fileInfo.Size() > 1*1024*1024 { // 1MB
		a.showMessage("Arquivo muito grande para edição")
		return
	}

	// Ler conteúdo do arquivo
	content, err := utils.GetFileContent(filePath, 1*1024*1024)
	if err != nil {
		a.showMessage(fmt.Sprintf("Erro ao ler arquivo: %v", err))
		return
	}

	// Criar campo de texto para edição
	textArea := tview.NewTextArea().
		SetText(string(content), true)
	textArea.SetBorder(true).
		SetTitle(fmt.Sprintf(" Editando: %s ", selectedFile)).
		SetTitleAlign(tview.AlignLeft)

	// Criar botões
	saveButton := tview.NewButton("Salvar")
	cancelButton := tview.NewButton("Cancelar")

	// Configurar botões
	saveButton.SetSelectedFunc(func() {
		// Obter conteúdo editado
		newContent := textArea.GetText()

		// Salvar arquivo
		err := os.WriteFile(filePath, []byte(newContent), fileInfo.Mode())
		if err != nil {
			a.showMessage(fmt.Sprintf("Erro ao salvar arquivo: %v", err))
			return
		}

		// Voltar para a visualização principal
		a.app.SetRoot(a.mainLayout, true)
		a.app.SetInputCapture(a.handleKeyEvents)
		a.refreshView()
	})

	cancelButton.SetSelectedFunc(func() {
		// Voltar para a visualização principal
		a.app.SetRoot(a.mainLayout, true)
		a.app.SetInputCapture(a.handleKeyEvents)
	})

	// Criar layout
	buttons := tview.NewFlex().
		AddItem(nil, 0, 1, false).
		AddItem(saveButton, 10, 0, true).
		AddItem(nil, 1, 0, false).
		AddItem(cancelButton, 10, 0, true).
		AddItem(nil, 0, 1, false)

	flex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(textArea, 0, 1, true).
		AddItem(buttons, 1, 0, false)

	// Configurar manipulador de teclas
	flex.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEscape {
			a.app.SetRoot(a.mainLayout, true)
			a.app.SetInputCapture(a.handleKeyEvents)
			return nil
		}
		return event
	})

	// Exibir editor
	a.app.SetRoot(flex, true)
	a.app.SetFocus(textArea)
}

// createDirectory cria um novo diretório
func (a *App) createDirectory() {
	a.showInputDialog("Novo Diretório", "Nome:", func(dirName string) {
		if dirName == "" {
			return
		}
		
		// Verificar se o nome contém caracteres inválidos
		if strings.ContainsAny(dirName, "\\/:*?\"<>|") {
			a.showError("O nome contém caracteres inválidos")
			return
		}
		
		// Criar caminho completo
		dirPath := filepath.Join(a.currentDir, dirName)
		
		// Verificar se o diretório já existe
		if _, err := os.Stat(dirPath); err == nil {
			a.showError(fmt.Sprintf("Já existe um arquivo ou diretório com o nome '%s'", dirName))
			return
		}
		
		// Criar diretório
		err := os.MkdirAll(dirPath, 0755)
		if err != nil {
			a.showError(fmt.Sprintf("Erro ao criar diretório: %v", err))
			return
		}
		
		a.refreshCurrentDir()
		a.showMessage(fmt.Sprintf("Diretório '%s' criado com sucesso", dirName))
	})
}

// toggleHiddenFiles alterna a exibição de arquivos ocultos
func (a *App) toggleHiddenFiles() {
	a.showHidden = !a.showHidden
	
	if a.showHidden {
		a.statusBar.SetText("Arquivos ocultos: Visíveis")
	} else {
		a.statusBar.SetText("Arquivos ocultos: Ocultos")
	}
	
	a.refreshCurrentDir()
}

// toggleFocus alterna o foco entre as visualizações
func (a *App) toggleFocus() {
	// Verificar foco atual
	if a.app.GetFocus() == a.fileView.fileList {
		a.app.SetFocus(a.treeView.TreeView)
	} else {
		a.app.SetFocus(a.fileView.fileList)
	}
}

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
		a.syncDirectoriesDialog()
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
	
	// Adicionar página
	a.pages.AddPage("toolsMenu", menu, true, true)
	a.app.SetFocus(menu)
}

// getHomeDir retorna o diretório home do usuário
func (a *App) getHomeDir() (string, error) {
	// Obter diretório home
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	
	return homeDir, nil
}

// deleteFile exclui um arquivo ou diretório
func (a *App) deleteFile() {
	selectedFile := a.getSelectedFile()
	if selectedFile == "" {
		a.showError("Nenhum arquivo selecionado")
		return
	}
	
	fileInfo, err := os.Stat(selectedFile)
	if err != nil {
		a.showError(fmt.Sprintf("Erro ao acessar arquivo: %v", err))
		return
	}
	
	fileType := "arquivo"
	if fileInfo.IsDir() {
		fileType = "diretório"
	}
	
	message := fmt.Sprintf("Deseja realmente excluir o %s '%s'?", fileType, filepath.Base(selectedFile))
	a.showConfirmDialog("Excluir", message, func(confirmed bool) {
		if !confirmed {
			return
		}
		
		var err error
		if fileInfo.IsDir() {
			err = os.RemoveAll(selectedFile)
		} else {
			err = os.Remove(selectedFile)
		}
		
		if err != nil {
			a.showError(fmt.Sprintf("Erro ao excluir: %v", err))
			return
		}
		
		a.refreshCurrentDir()
		a.showMessage(fmt.Sprintf("%s excluído com sucesso", strings.Title(fileType)))
	})
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
		defer file.Close()
		
		a.refreshCurrentDir()
		a.showMessage(fmt.Sprintf("Arquivo '%s' criado com sucesso", fileName))
	})
}

// renameFile renomeia um arquivo ou diretório
func (a *App) renameFile() {
	// Obter arquivo selecionado
	file := a.fileView.GetSelectedFile()
	if file == "" || file == ".." {
		return
	}
	
	// Obter caminho completo
	filePath := filepath.Join(a.currentDir, file)
	
	// Verificar se o arquivo existe
	_, err := os.Stat(filePath)
	if err != nil {
		a.showError(fmt.Sprintf("Erro ao acessar arquivo: %v", err))
		return
	}
	
	// Exibir diálogo para renomear
	a.showInputDialogWithValue("Renomear para:", file, func(newName string) {
		if newName == "" || newName == file {
			return
		}
		
		// Obter caminho completo do novo nome
		newPath := filepath.Join(a.currentDir, newName)
		
		// Verificar se o novo nome já existe
		_, err := os.Stat(newPath)
		if err == nil {
			a.showConfirmDialog("Confirmação", fmt.Sprintf("O arquivo '%s' já existe. Sobrescrever?", newName), func(confirmed bool) {
				if confirmed {
					// Renomear arquivo
					err := os.Rename(filePath, newPath)
					if err != nil {
						a.showError(fmt.Sprintf("Erro ao renomear arquivo: %v", err))
						return
					}
					
					// Atualizar visualização
					a.refreshView()
				}
			})
			return
		}
		
		// Renomear arquivo
		err = os.Rename(filePath, newPath)
		if err != nil {
			a.showError(fmt.Sprintf("Erro ao renomear arquivo: %v", err))
			return
		}
		
		// Atualizar visualização
		a.refreshView()
	})
}

// copyFile copia um arquivo ou diretório
func (a *App) copyFile() {
	selectedFile := a.getSelectedFile()
	if selectedFile == "" {
		a.showError("Nenhum arquivo selecionado")
		return
	}
	
	fileInfo, err := os.Stat(selectedFile)
	if err != nil {
		a.showError(fmt.Sprintf("Erro ao acessar arquivo: %v", err))
		return
	}
	
	// Obter diretório de destino
	a.showInputDialog("Copiar para", "Diretório de destino:", func(destDir string) {
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
				if confirmed {
					// Copiar arquivo ou diretório
					a.doCopy(selectedFile, destPath, fileInfo.IsDir())
				}
			})
		} else {
			// Copiar arquivo ou diretório
			a.doCopy(selectedFile, destPath, fileInfo.IsDir())
		}
	})
}

// doCopy realiza a cópia de um arquivo ou diretório
func (a *App) doCopy(src, dest string, isDir bool) {
	// Verificar se é um diretório
	if isDir {
		// Copiar diretório
		err := utils.CopyDir(src, dest)
		if err != nil {
			a.showError(fmt.Sprintf("Erro ao copiar diretório: %v", err))
		} else {
			a.refreshCurrentDir()
		}
	} else {
		// Copiar arquivo
		err := utils.CopyFile(src, dest)
		if err != nil {
			a.showError(fmt.Sprintf("Erro ao copiar arquivo: %v", err))
		} else {
			a.refreshCurrentDir()
		}
	}
}

// getSelectedFile retorna o caminho completo do arquivo selecionado
func (a *App) getSelectedFile() string {
	// Obter arquivo selecionado
	fileName := a.fileView.GetSelectedFile()
	if fileName == "" {
		return ""
	}
	
	// Verificar se é o diretório pai
	if fileName == ".." {
		return filepath.Dir(a.currentDir)
	}
	
	// Retornar caminho completo
	return filepath.Join(a.currentDir, fileName)
}

// refreshCurrentDir atualiza a visualização do diretório atual
func (a *App) refreshCurrentDir() {
	a.navigateTo(a.currentDir)
}

// syncDirectoriesDialog exibe o diálogo para sincronizar diretórios
func (a *App) syncDirectoriesDialog() {
	// Criar formulário
	form := tview.NewForm()
	
	// Adicionar campos
	form.AddInputField("Diretório de origem:", a.currentDir, 40, nil, nil)
	form.AddInputField("Diretório de destino:", "", 40, nil, nil)
	
	// Adicionar opções
	form.AddCheckbox("Excluir arquivos órfãos no destino", false, nil)
	form.AddCheckbox("Sobrescrever arquivos mais novos", true, nil)
	form.AddCheckbox("Pular arquivos existentes", false, nil)
	form.AddCheckbox("Incluir arquivos ocultos", a.showHidden, nil)
	form.AddCheckbox("Apenas visualizar (não realizar alterações)", false, nil)
	
	// Adicionar botões
	form.AddButton("Sincronizar", func() {
		// Obter valores dos campos
		sourceDir := form.GetFormItem(0).(*tview.InputField).GetText()
		targetDir := form.GetFormItem(1).(*tview.InputField).GetText()
		deleteOrphans := form.GetFormItem(2).(*tview.Checkbox).IsChecked()
		overwriteNewer := form.GetFormItem(3).(*tview.Checkbox).IsChecked()
		skipExisting := form.GetFormItem(4).(*tview.Checkbox).IsChecked()
		includeHidden := form.GetFormItem(5).(*tview.Checkbox).IsChecked()
		previewOnly := form.GetFormItem(6).(*tview.Checkbox).IsChecked()
		
		// Verificar se os diretórios são válidos
		if sourceDir == "" || targetDir == "" {
			a.showError("Os diretórios de origem e destino são obrigatórios")
			return
		}
		
		// Expandir caminhos
		if strings.HasPrefix(sourceDir, "~") {
			homeDir, err := a.getHomeDir()
			if err == nil {
				sourceDir = filepath.Join(homeDir, sourceDir[1:])
			}
		}
		
		if strings.HasPrefix(targetDir, "~") {
			homeDir, err := a.getHomeDir()
			if err == nil {
				targetDir = filepath.Join(homeDir, targetDir[1:])
			}
		}
		
		// Verificar se os diretórios existem
		sourceStat, err := os.Stat(sourceDir)
		if err != nil || !sourceStat.IsDir() {
			a.showError(fmt.Sprintf("O diretório de origem '%s' não existe", sourceDir))
			return
		}
		
		targetStat, err := os.Stat(targetDir)
		if err != nil || !targetStat.IsDir() {
			a.showError(fmt.Sprintf("O diretório de destino '%s' não existe", targetDir))
			return
		}
		
		// Fechar o diálogo
		a.pages.RemovePage("syncDialog")
		
		// Iniciar sincronização
		a.syncDirectories(sourceDir, targetDir, deleteOrphans, overwriteNewer, skipExisting, includeHidden, previewOnly)
	})
	
	form.AddButton("Cancelar", func() {
		a.pages.RemovePage("syncDialog")
	})
	
	// Configurar borda
	form.SetBorder(true).
		SetTitle(" Sincronizar Diretórios ").
		SetTitleAlign(tview.AlignLeft)
	
	// Configurar manipulador de teclas
	form.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEscape {
			a.pages.RemovePage("syncDialog")
			return nil
		}
		return event
	})
	
	// Exibir diálogo
	a.pages.AddPage("syncDialog", a.modal(form, 60, 15), true, true)
}

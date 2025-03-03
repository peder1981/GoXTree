package ui

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/peder1981/GoXTree/pkg/utils"
	"github.com/rivo/tview"

	"github.com/gdamore/tcell/v2"
)

// showSearchDialog exibe o diálogo de busca
func (a *App) showSearchDialog() {
	// Criar formulário
	form := tview.NewForm()
	form.SetTitle("Busca de Arquivos").
		SetTitleAlign(tview.AlignCenter).
		SetBorder(true)

	// Variáveis para os campos
	var (
		pattern       string
		searchDir     string = a.currentDir
		recursive     bool   = true
		matchCase     bool   = false
		searchContent bool   = false
		fileType      string
	)

	// Adicionar campos
	form.AddInputField("Padrão de busca:", "", 40, nil, func(text string) {
		pattern = text
	})

	form.AddInputField("Diretório de busca:", searchDir, 40, nil, func(text string) {
		searchDir = text
	})

	form.AddCheckbox("Busca recursiva", recursive, func(checked bool) {
		recursive = checked
	})

	form.AddCheckbox("Diferenciar maiúsculas/minúsculas", matchCase, func(checked bool) {
		matchCase = checked
	})

	form.AddCheckbox("Buscar no conteúdo dos arquivos", searchContent, func(checked bool) {
		searchContent = checked
	})

	form.AddInputField("Tipo de arquivo (ex: .txt, .go):", "", 40, nil, func(text string) {
		fileType = text
	})

	// Adicionar botões
	form.AddButton("Buscar", func() {
		a.pages.RemovePage("searchDialog")
		a.performSearch(pattern, searchDir, recursive, matchCase, searchContent, fileType)
	})

	form.AddButton("Cancelar", func() {
		a.pages.RemovePage("searchDialog")
	})

	// Configurar manipulador de eventos
	form.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEscape {
			a.pages.RemovePage("searchDialog")
			return nil
		}
		return event
	})

	// Adicionar página
	a.pages.AddPage("searchDialog", form, true, true)
	a.app.SetFocus(form)
}

// performSearch realiza a busca de arquivos
func (a *App) performSearch(pattern, searchDir string, recursive, matchCase, searchContent bool, fileType string) {
	// Verificar padrão de busca
	if pattern == "" {
		a.showError("Padrão de busca não pode ser vazio")
		return
	}

	// Verificar diretório de busca
	if searchDir == "" {
		searchDir = a.currentDir
	}

	// Expandir caminho
	if strings.HasPrefix(searchDir, "~") {
		homeDir, err := a.getHomeDir()
		if err != nil {
			a.showError(fmt.Sprintf("Erro ao obter diretório home: %v", err))
			return
		}
		searchDir = filepath.Join(homeDir, searchDir[1:])
	}

	// Verificar se o diretório existe
	fileInfo, err := os.Stat(searchDir)
	if err != nil || !fileInfo.IsDir() {
		a.showError(fmt.Sprintf("Diretório de busca inválido: %s", searchDir))
		return
	}

	// Criar lista de resultados
	resultList := tview.NewList()
	resultList.SetTitle(fmt.Sprintf("Resultados da busca: %s", pattern)).
		SetTitleAlign(tview.AlignCenter).
		SetBorder(true)

	// Adicionar manipulador de eventos
	resultList.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEscape {
			a.pages.RemovePage("searchResults")
			return nil
		}
		return event
	})

	// Adicionar página
	a.pages.AddPage("searchResults", resultList, true, true)
	a.app.SetFocus(resultList)

	// Realizar busca em goroutine
	go func() {
		// Resultados
		var results []string

		// Buscar arquivos
		err := filepath.Walk(searchDir, func(path string, info os.FileInfo, err error) error {
			// Verificar erro
			if err != nil {
				return nil
			}

			// Verificar se é diretório e não é recursivo
			if info.IsDir() && path != searchDir && !recursive {
				return filepath.SkipDir
			}

			// Verificar tipo de arquivo
			if fileType != "" && !info.IsDir() {
				if !strings.HasSuffix(strings.ToLower(path), strings.ToLower(fileType)) {
					return nil
				}
			}

			// Verificar nome do arquivo
			fileName := filepath.Base(path)
			match := false

			if matchCase {
				match = strings.Contains(fileName, pattern)
			} else {
				match = strings.Contains(strings.ToLower(fileName), strings.ToLower(pattern))
			}

			// Verificar conteúdo do arquivo
			if !match && searchContent && !info.IsDir() {
				// Verificar tamanho do arquivo
				if info.Size() > 10*1024*1024 {
					return nil // Ignorar arquivos muito grandes
				}

				// Ler conteúdo do arquivo
				content, err := os.ReadFile(path)
				if err != nil {
					return nil
				}

				// Verificar conteúdo
				if matchCase {
					match = strings.Contains(string(content), pattern)
				} else {
					match = strings.Contains(strings.ToLower(string(content)), strings.ToLower(pattern))
				}
			}

			// Adicionar ao resultado
			if match {
				results = append(results, path)
			}

			return nil
		})

		// Verificar erro
		if err != nil {
			a.app.QueueUpdateDraw(func() {
				a.showError(fmt.Sprintf("Erro ao buscar arquivos: %v", err))
				a.pages.RemovePage("searchResults")
			})
			return
		}

		// Atualizar lista de resultados
		a.app.QueueUpdateDraw(func() {
			// Verificar se há resultados
			if len(results) == 0 {
				resultList.AddItem("Nenhum resultado encontrado", "", 0, nil)
				return
			}

			// Adicionar resultados
			for _, path := range results {
				resultList.AddItem(path, "", 0, func() {
					// Navegar para o diretório do arquivo
					dir := filepath.Dir(path)
					file := filepath.Base(path)

					a.navigateTo(dir)
					a.fileView.SelectFile(file)
					a.pages.RemovePage("searchResults")
				})
			}
		})
	}()
}

// showCompareDialog exibe o diálogo de comparação de arquivos
func (a *App) showCompareDialog() {
	// Criar formulário
	form := tview.NewForm()
	form.SetTitle("Comparar Arquivos").
		SetTitleAlign(tview.AlignCenter).
		SetBorder(true)

	// Variáveis para os campos
	var (
		file1 string
		file2 string
	)

	// Adicionar campos
	form.AddInputField("Arquivo 1:", "", 40, nil, func(text string) {
		file1 = text
	})

	form.AddInputField("Arquivo 2:", "", 40, nil, func(text string) {
		file2 = text
	})

	// Adicionar botões
	form.AddButton("Comparar", func() {
		a.pages.RemovePage("compareDialog")

		// Verificar arquivos
		if file1 == "" || file2 == "" {
			a.showError("Ambos os arquivos devem ser especificados")
			return
		}

		// Expandir caminhos
		if strings.HasPrefix(file1, "~") {
			homeDir, err := a.getHomeDir()
			if err == nil {
				file1 = filepath.Join(homeDir, file1[1:])
			}
		}

		if strings.HasPrefix(file2, "~") {
			homeDir, err := a.getHomeDir()
			if err == nil {
				file2 = filepath.Join(homeDir, file2[1:])
			}
		}

		// Verificar se os arquivos existem
		info1, err1 := os.Stat(file1)
		info2, err2 := os.Stat(file2)

		if err1 != nil || err2 != nil {
			a.showError("Um ou ambos os arquivos não existem")
			return
		}

		if info1.IsDir() || info2.IsDir() {
			a.showError("Não é possível comparar diretórios")
			return
		}

		// Comparar arquivos
		a.compareFiles(file1, file2)
	})

	form.AddButton("Cancelar", func() {
		a.pages.RemovePage("compareDialog")
	})

	// Configurar manipulador de eventos
	form.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEscape {
			a.pages.RemovePage("compareDialog")
			return nil
		}
		return event
	})

	// Adicionar página
	a.pages.AddPage("compareDialog", form, true, true)
	a.app.SetFocus(form)
}

// showSyncDialog exibe o diálogo de sincronização de diretórios
func (a *App) showSyncDialog() {
	// Criar formulário
	form := tview.NewForm()
	form.SetTitle("Sincronizar Diretórios").
		SetTitleAlign(tview.AlignCenter).
		SetBorder(true)

	// Variáveis para os campos
	var (
		sourceDir      string = a.currentDir
		targetDir      string
		deleteOrphans  bool = false
		overwriteNewer bool = false
		skipExisting   bool = false
		includeHidden  bool = false
		previewOnly    bool = true
	)

	// Adicionar campos
	form.AddInputField("Diretório de origem:", sourceDir, 40, nil, func(text string) {
		sourceDir = text
	})

	form.AddInputField("Diretório de destino:", "", 40, nil, func(text string) {
		targetDir = text
	})

	form.AddCheckbox("Excluir arquivos órfãos no destino", deleteOrphans, func(checked bool) {
		deleteOrphans = checked
	})

	form.AddCheckbox("Sobrescrever arquivos mais novos", overwriteNewer, func(checked bool) {
		overwriteNewer = checked
	})

	form.AddCheckbox("Ignorar arquivos existentes", skipExisting, func(checked bool) {
		skipExisting = checked
	})

	form.AddCheckbox("Incluir arquivos ocultos", includeHidden, func(checked bool) {
		includeHidden = checked
	})

	form.AddCheckbox("Apenas visualizar (sem executar)", previewOnly, func(checked bool) {
		previewOnly = checked
	})

	// Adicionar botões
	form.AddButton("Sincronizar", func() {
		a.pages.RemovePage("syncDialog")

		// Verificar diretórios
		if sourceDir == "" || targetDir == "" {
			a.showError("Ambos os diretórios devem ser especificados")
			return
		}

		// Expandir caminhos
		if strings.HasPrefix(sourceDir, "~") {
			homeDir, err := os.UserHomeDir()
			if err == nil {
				sourceDir = filepath.Join(homeDir, sourceDir[1:])
			}
		}

		if strings.HasPrefix(targetDir, "~") {
			homeDir, err := os.UserHomeDir()
			if err == nil {
				targetDir = filepath.Join(homeDir, targetDir[1:])
			}
		}

		// Verificar se os diretórios existem
		sourceInfo, err1 := os.Stat(sourceDir)
		targetInfo, err2 := os.Stat(targetDir)

		if err1 != nil {
			a.showError(fmt.Sprintf("Diretório de origem não existe: %s", sourceDir))
			return
		}

		if !sourceInfo.IsDir() {
			a.showError("Origem não é um diretório")
			return
		}

		if err2 == nil && !targetInfo.IsDir() {
			a.showError("Destino não é um diretório")
			return
		}

		// Sincronizar diretórios
		options := utils.SyncOptions{
			SourceDir:      sourceDir,
			DestDir:        targetDir,
			DeleteOrphaned: deleteOrphans,
			PreviewOnly:    previewOnly,
			SkipNewer:      !overwriteNewer,
			SkipExisting:   skipExisting,
			IncludeHidden:  includeHidden,
		}

		actions, err := utils.SyncDirectories(options)
		if err != nil {
			a.showError(fmt.Sprintf("Erro ao sincronizar: %v", err))
			return
		}

		// Mostrar resultado
		a.showSyncResults(actions)

		// Atualizar visualização
		a.refreshView()
	})

	form.AddButton("Cancelar", func() {
		a.pages.RemovePage("syncDialog")
	})

	// Configurar manipulador de eventos
	form.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEscape {
			a.pages.RemovePage("syncDialog")
			return nil
		}
		return event
	})

	// Adicionar página
	a.pages.AddPage("syncDialog", form, true, true)
	a.app.SetFocus(form)
}

// showConfirmDialog exibe o diálogo de confirmação
func (a *App) showConfirmDialog(title, message string, callback func(confirmed bool)) {
	// Criar formulário
	form := tview.NewForm()
	form.SetTitle(title).
		SetTitleAlign(tview.AlignCenter).
		SetBorder(true)

	// Adicionar mensagem
	form.AddFormItem(tview.NewTextView().SetText(message))

	// Adicionar botões
	form.AddButton("Sim", func() {
		a.pages.RemovePage("confirmDialog")
		callback(true)
	})

	form.AddButton("Não", func() {
		a.pages.RemovePage("confirmDialog")
		callback(false)
	})

	// Configurar manipulador de eventos
	form.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEscape {
			a.pages.RemovePage("confirmDialog")
			callback(false)
			return nil
		}
		return event
	})

	// Adicionar página
	a.pages.AddPage("confirmDialog", form, true, true)
	a.app.SetFocus(form)
}

// showSimpleSearchDialog exibe o diálogo de busca simples
func (a *App) showSimpleSearchDialog() {
	// Criar formulário
	form := tview.NewForm()
	form.SetTitle("Busca Simples").
		SetTitleAlign(tview.AlignCenter).
		SetBorder(true)

	// Adicionar campo de busca
	form.AddInputField("Buscar por:", "", 40, nil, nil)

	// Adicionar botões
	form.AddButton("Buscar", func() {
		// Obter padrão de busca
		pattern := form.GetFormItem(0).(*tview.InputField).GetText()
		if pattern == "" {
			a.showError("Informe um padrão de busca")
			return
		}

		// Fechar diálogo
		a.pages.RemovePage("simpleSearchDialog")

		// Realizar busca
		a.performSearch(pattern, a.currentDir, true, false, false, "")
	})

	form.AddButton("Cancelar", func() {
		a.pages.RemovePage("simpleSearchDialog")
	})

	// Configurar manipulador de teclas
	form.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEscape {
			a.pages.RemovePage("simpleSearchDialog")
			return nil
		}
		return event
	})

	// Exibir diálogo
	a.pages.AddPage("simpleSearchDialog", a.modal(form, 50, 7), true, true)
}

// showAdvancedSearchDialog exibe o diálogo de busca avançada
func (a *App) showAdvancedSearchDialog() {
	// Criar formulário
	form := tview.NewForm()
	form.SetTitle("Busca Avançada").
		SetTitleAlign(tview.AlignCenter).
		SetBorder(true)

	// Adicionar campos
	form.AddInputField("Buscar por:", "", 40, nil, nil)
	form.AddInputField("Diretório:", a.currentDir, 40, nil, nil)
	form.AddDropDown("Tipo de arquivo:", []string{"Todos", "Documentos", "Imagens", "Vídeos", "Áudio", "Compactados", "Código", "Personalizado"}, 0, nil)
	form.AddInputField("Extensão personalizada:", "", 20, nil, nil)

	// Adicionar opções
	form.AddCheckbox("Buscar em subdiretórios", true, nil)
	form.AddCheckbox("Diferenciar maiúsculas/minúsculas", false, nil)
	form.AddCheckbox("Buscar no conteúdo dos arquivos", false, nil)

	// Adicionar botões
	form.AddButton("Buscar", func() {
		// Obter valores dos campos
		pattern := form.GetFormItem(0).(*tview.InputField).GetText()
		searchDir := form.GetFormItem(1).(*tview.InputField).GetText()
		fileTypeIndex, _ := form.GetFormItem(2).(*tview.DropDown).GetCurrentOption()
		customExt := form.GetFormItem(3).(*tview.InputField).GetText()
		recursive := form.GetFormItem(4).(*tview.Checkbox).IsChecked()
		matchCase := form.GetFormItem(5).(*tview.Checkbox).IsChecked()
		searchContent := form.GetFormItem(6).(*tview.Checkbox).IsChecked()

		// Validar campos
		if pattern == "" {
			a.showError("Informe um padrão de busca")
			return
		}

		if searchDir == "" {
			a.showError("Informe um diretório de busca")
			return
		}

		// Expandir caminho
		if strings.HasPrefix(searchDir, "~") {
			homeDir, err := os.UserHomeDir()
			if err == nil {
				searchDir = filepath.Join(homeDir, searchDir[1:])
			}
		}

		// Verificar se o diretório existe
		fileInfo, err := os.Stat(searchDir)
		if err != nil || !fileInfo.IsDir() {
			a.showError(fmt.Sprintf("O diretório '%s' não existe", searchDir))
			return
		}

		// Determinar tipo de arquivo
		var fileType string
		switch fileTypeIndex {
		case 1: // Documentos
			fileType = ".doc,.docx,.pdf,.txt,.rtf,.odt"
		case 2: // Imagens
			fileType = ".jpg,.jpeg,.png,.gif,.bmp,.tiff,.svg"
		case 3: // Vídeos
			fileType = ".mp4,.avi,.mkv,.mov,.wmv,.flv"
		case 4: // Áudio
			fileType = ".mp3,.wav,.ogg,.flac,.aac"
		case 5: // Compactados
			fileType = ".zip,.rar,.7z,.tar,.gz"
		case 6: // Código
			fileType = ".go,.c,.cpp,.h,.py,.js,.html,.css,.java"
		case 7: // Personalizado
			fileType = customExt
		}

		// Fechar diálogo
		a.pages.RemovePage("advancedSearchDialog")

		// Realizar busca
		a.performSearch(pattern, searchDir, recursive, matchCase, searchContent, fileType)
	})

	form.AddButton("Cancelar", func() {
		a.pages.RemovePage("advancedSearchDialog")
	})

	// Configurar manipulador de teclas
	form.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEscape {
			a.pages.RemovePage("advancedSearchDialog")
			return nil
		}
		return event
	})

	// Exibir diálogo
	a.pages.AddPage("advancedSearchDialog", a.modal(form, 60, 15), true, true)
}

package ui

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/peder1981/GoXTree/pkg/utils"
	"github.com/rivo/tview"
	"github.com/sergi/go-diff/diffmatchpatch"
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
		pattern     string
		searchDir   string = a.currentDir
		recursive   bool   = true
		matchCase   bool   = false
		searchContent bool = false
		fileType    string
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

// compareFiles compara dois arquivos
func (a *App) compareFiles(file1, file2 string) {
	// Ler conteúdo dos arquivos
	content1, err1 := os.ReadFile(file1)
	content2, err2 := os.ReadFile(file2)
	
	if err1 != nil || err2 != nil {
		a.showError("Erro ao ler arquivos")
		return
	}
	
	// Comparar conteúdo
	if bytes.Equal(content1, content2) {
		a.showConfirmDialog("Comparação", "Os arquivos são idênticos", func(confirmed bool) {
			// Não fazer nada
		})
		return
	}
	
	// Criar visualização de diferenças
	diffView := tview.NewTextView()
	diffView.SetDynamicColors(true).
		SetScrollable(true).
		SetTitle(fmt.Sprintf("Comparação: %s <-> %s", filepath.Base(file1), filepath.Base(file2))).
		SetTitleAlign(tview.AlignCenter).
		SetBorder(true)
	
	// Configurar manipulador de eventos
	diffView.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEscape {
			a.pages.RemovePage("diffView")
			return nil
		}
		return event
	})
	
	// Adicionar página
	a.pages.AddPage("diffView", diffView, true, true)
	a.app.SetFocus(diffView)
	
	// Calcular diferenças em goroutine
	go func() {
		// Converter conteúdo para strings
		text1 := string(content1)
		text2 := string(content2)
		
		// Calcular diferenças usando diffmatchpatch
		dmp := diffmatchpatch.New()
		diffs := dmp.DiffMain(text1, text2, false)
		
		// Formatar saída
		var output strings.Builder
		
		for _, diff := range diffs {
			switch diff.Type {
			case diffmatchpatch.DiffEqual:
				output.WriteString(fmt.Sprintf("  %s\n", diff.Text))
			case diffmatchpatch.DiffDelete:
				output.WriteString(fmt.Sprintf("[red]- %s[white]\n", diff.Text))
			case diffmatchpatch.DiffInsert:
				output.WriteString(fmt.Sprintf("[green]+ %s[white]\n", diff.Text))
			}
		}
		
		// Atualizar visualização
		a.app.QueueUpdateDraw(func() {
			fmt.Fprintf(diffView, "%s", output.String())
		})
	}()
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
		deleteOrphans  bool   = false
		overwriteNewer bool   = false
		skipExisting   bool   = false
		includeHidden  bool   = false
		previewOnly    bool   = true
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
		a.syncDirectories(sourceDir, targetDir, deleteOrphans, overwriteNewer, skipExisting, includeHidden, previewOnly)
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

// syncDirectories sincroniza dois diretórios
func (a *App) syncDirectories(sourceDir, targetDir string, deleteOrphans, overwriteNewer, skipExisting, includeHidden, previewOnly bool) {
	// Criar visualização de texto
	textView := tview.NewTextView()
	textView.SetDynamicColors(true).
		SetScrollable(true).
		SetTitle(fmt.Sprintf("Sincronização: %s -> %s", sourceDir, targetDir)).
		SetTitleAlign(tview.AlignCenter).
		SetBorder(true)
	
	// Configurar manipulador de eventos
	textView.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEscape {
			a.pages.RemovePage("syncView")
			return nil
		}
		return event
	})
	
	// Adicionar página
	a.pages.AddPage("syncView", textView, true, true)
	a.app.SetFocus(textView)
	
	// Realizar sincronização em goroutine
	go func() {
		// Verificar se o diretório de destino existe
		_, err := os.Stat(targetDir)
		if os.IsNotExist(err) {
			// Criar diretório de destino
			err = os.MkdirAll(targetDir, 0755)
			if err != nil {
				a.app.QueueUpdateDraw(func() {
					fmt.Fprintf(textView, "[red]Erro ao criar diretório de destino: %v[white]", err)
				})
				return
			}
		}
		
		// Estatísticas
		var (
			totalFiles    int
			copiedFiles   int
			skippedFiles  int
			deletedFiles  int
			errorFiles    int
			totalSize     int64
			copiedSize    int64
		)
		
		// Função para adicionar linha ao log
		addLog := func(format string, args ...interface{}) {
			a.app.QueueUpdateDraw(func() {
				fmt.Fprintf(textView, format+"\n", args...)
			})
		}
		
		// Iniciar sincronização
		addLog("[yellow]Iniciando sincronização...[white]")
		
		// Mapear arquivos de destino
		destFiles := make(map[string]os.FileInfo)
		if deleteOrphans {
			addLog("Mapeando arquivos de destino...")
			
			err := filepath.Walk(targetDir, func(path string, info os.FileInfo, err error) error {
				if err != nil {
					return nil
				}
				
				// Ignorar o diretório raiz
				if path == targetDir {
					return nil
				}
				
				// Verificar se é arquivo oculto
				if !includeHidden && strings.HasPrefix(filepath.Base(path), ".") {
					if info.IsDir() {
						return filepath.SkipDir
					}
					return nil
				}
				
				// Adicionar ao mapa
				relPath, _ := filepath.Rel(targetDir, path)
				destFiles[relPath] = info
				
				return nil
			})
			
			if err != nil {
				addLog("[red]Erro ao mapear arquivos de destino: %v[white]", err)
			}
		}
		
		// Sincronizar arquivos
		err = filepath.Walk(sourceDir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				addLog("[red]Erro ao acessar %s: %v[white]", path, err)
				errorFiles++
				return nil
			}
			
			// Ignorar o diretório raiz
			if path == sourceDir {
				return nil
			}
			
			// Verificar se é arquivo oculto
			if !includeHidden && strings.HasPrefix(filepath.Base(path), ".") {
				if info.IsDir() {
					return filepath.SkipDir
				}
				return nil
			}
			
			// Calcular caminho relativo
			relPath, _ := filepath.Rel(sourceDir, path)
			destPath := filepath.Join(targetDir, relPath)
			
			// Remover do mapa de destino
			if deleteOrphans {
				delete(destFiles, relPath)
			}
			
			// Verificar se é diretório
			if info.IsDir() {
				// Criar diretório de destino
				if !previewOnly {
					err = os.MkdirAll(destPath, info.Mode())
					if err != nil {
						addLog("[red]Erro ao criar diretório %s: %v[white]", destPath, err)
						errorFiles++
					} else {
						addLog("Criado diretório: %s", relPath)
					}
				} else {
					addLog("[cyan]Seria criado diretório: %s[white]", relPath)
				}
				
				return nil
			}
			
			// Verificar se o arquivo já existe
			destInfo, err := os.Stat(destPath)
			if err == nil {
				// Arquivo existe
				if skipExisting {
					addLog("Ignorado arquivo existente: %s", relPath)
					skippedFiles++
					return nil
				}
				
				// Verificar data de modificação
				if !overwriteNewer && destInfo.ModTime().After(info.ModTime()) {
					addLog("Ignorado arquivo mais novo no destino: %s", relPath)
					skippedFiles++
					return nil
				}
				
				// Verificar tamanho
				if destInfo.Size() == info.Size() && destInfo.ModTime() == info.ModTime() {
					addLog("Ignorado arquivo idêntico: %s", relPath)
					skippedFiles++
					return nil
				}
			}
			
			// Copiar arquivo
			totalFiles++
			totalSize += info.Size()
			
			if !previewOnly {
				// Criar diretório pai
				err = os.MkdirAll(filepath.Dir(destPath), 0755)
				if err != nil {
					addLog("[red]Erro ao criar diretório pai para %s: %v[white]", destPath, err)
					errorFiles++
					return nil
				}
				
				// Copiar arquivo
				err = utils.CopyFile(path, destPath)
				if err != nil {
					addLog("[red]Erro ao copiar %s: %v[white]", relPath, err)
					errorFiles++
				} else {
					addLog("Copiado: %s (%s)", relPath, utils.FormatFileSize(info.Size()))
					copiedFiles++
					copiedSize += info.Size()
				}
			} else {
				addLog("[cyan]Seria copiado: %s (%s)[white]", relPath, utils.FormatFileSize(info.Size()))
				copiedFiles++
				copiedSize += info.Size()
			}
			
			return nil
		})
		
		if err != nil {
			addLog("[red]Erro ao sincronizar: %v[white]", err)
		}
		
		// Excluir arquivos órfãos
		if deleteOrphans && len(destFiles) > 0 {
			addLog("\n[yellow]Excluindo arquivos órfãos...[white]")
			
			for relPath, info := range destFiles {
				destPath := filepath.Join(targetDir, relPath)
				
				if !previewOnly {
					var err error
					if info.IsDir() {
						err = os.RemoveAll(destPath)
					} else {
						err = os.Remove(destPath)
					}
					
					if err != nil {
						addLog("[red]Erro ao excluir %s: %v[white]", relPath, err)
						errorFiles++
					} else {
						if info.IsDir() {
							addLog("Excluído diretório: %s", relPath)
						} else {
							addLog("Excluído arquivo: %s", relPath)
							deletedFiles++
						}
					}
				} else {
					if info.IsDir() {
						addLog("[cyan]Seria excluído diretório: %s[white]", relPath)
					} else {
						addLog("[cyan]Seria excluído arquivo: %s[white]", relPath)
						deletedFiles++
					}
				}
			}
		}
		
		// Exibir resumo
		addLog("\n[yellow]Resumo da sincronização:[white]")
		addLog("Total de arquivos: %d", totalFiles)
		addLog("Arquivos copiados: %d", copiedFiles)
		addLog("Arquivos ignorados: %d", skippedFiles)
		addLog("Arquivos excluídos: %d", deletedFiles)
		addLog("Erros: %d", errorFiles)
		addLog("Tamanho total: %s", utils.FormatFileSize(totalSize))
		addLog("Tamanho copiado: %s", utils.FormatFileSize(copiedSize))
		
		if previewOnly {
			addLog("\n[yellow]Esta foi apenas uma visualização. Nenhuma alteração foi realizada.[white]")
		}
	}()
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

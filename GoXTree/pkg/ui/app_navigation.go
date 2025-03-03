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

// refreshView atualiza as visualizações
func (a *App) refreshView() {
	a.treeView.LoadTree(a.currentDir)
	if err := a.fileView.SetCurrentDir(a.currentDir); err != nil {
		a.showError(fmt.Sprintf("Erro ao atualizar visualização: %s", err))
	}
	a.statusBar.UpdateStatus(a.currentDir)
}

// refreshTreeView atualiza a árvore de diretórios
func (a *App) refreshTreeView() {
	a.treeView.LoadTree(a.currentDir)
}

// refreshFileView atualiza a visualização de arquivos
func (a *App) refreshFileView() {
	if err := a.fileView.SetCurrentDir(a.currentDir); err != nil {
		a.showError(fmt.Sprintf("Erro ao atualizar visualização de arquivos: %s", err))
	}
}

// goBack navega para o diretório anterior no histórico
func (a *App) goBack() {
	if a.historyPos <= 0 {
		a.showMessage("Não há diretório anterior no histórico")
		return
	}

	a.historyPos--
	a.navigateTo(a.history[a.historyPos])
}

// goForward navega para o próximo diretório no histórico
func (a *App) goForward() {
	if a.historyPos >= len(a.history)-1 {
		a.showMessage("Não há diretório seguinte no histórico")
		return
	}

	a.historyPos++
	a.navigateTo(a.history[a.historyPos])
}

// goToParent navega para o diretório pai
func (a *App) goToParent() {
	// Verificar se já estamos na raiz
	if a.currentDir == "/" {
		return
	}

	// Obter diretório pai
	parentDir := filepath.Dir(a.currentDir)

	// Navegar para o diretório pai
	a.navigateTo(parentDir)
}

// goToHome navega para o diretório home do usuário
func (a *App) goToHome() {
	// Obter diretório home
	homeDir, err := os.UserHomeDir()
	if err != nil {
		a.showMessage(fmt.Sprintf("Erro ao obter diretório home: %v", err))
		return
	}

	// Navegar para o diretório home
	a.navigateTo(homeDir)
}

// goToRoot navega para o diretório raiz
func (a *App) goToRoot() {
	a.navigateTo("/")
}

// goToDirectory navega para um diretório específico
func (a *App) goToDirectory() {
	a.showInputDialog("Ir para Diretório", "Caminho:", func(dirPath string) {
		if dirPath == "" {
			return
		}
		
		// Expandir caminho
		if strings.HasPrefix(dirPath, "~") {
			homeDir, err := os.UserHomeDir()
			if err == nil {
				dirPath = filepath.Join(homeDir, dirPath[1:])
			}
		}
		
		// Verificar se é um caminho relativo
		if !filepath.IsAbs(dirPath) {
			dirPath = filepath.Join(a.currentDir, dirPath)
		}
		
		// Verificar se o diretório existe
		fileInfo, err := os.Stat(dirPath)
		if err != nil {
			a.showError(fmt.Sprintf("Diretório não encontrado: %s", dirPath))
			return
		}
		
		// Verificar se é um diretório
		if !fileInfo.IsDir() {
			a.showError(fmt.Sprintf("O caminho especificado não é um diretório: %s", dirPath))
			return
		}
		
		// Navegar para o diretório
		a.navigateTo(dirPath)
	})
}

// NavigateToParentDirectory navega para o diretório pai do diretório atual
func (a *App) NavigateToParentDirectory() {
	parent := filepath.Dir(a.currentDir)
	if parent != a.currentDir {
		// Adicionar diretório atual ao histórico
		a.addToHistory(a.currentDir)

		// Navegar para o diretório pai
		a.NavigateToDirectory(parent)
	}
}

// NavigateToDirectory navega para um diretório específico
func (a *App) NavigateToDirectory(path string) {
	// Verificar se o diretório existe
	fileInfo, err := os.Stat(path)
	if err != nil {
		a.showError(fmt.Sprintf("Erro ao acessar diretório: %v", err))
		return
	}

	// Verificar se é um diretório
	if !fileInfo.IsDir() {
		a.showError(fmt.Sprintf("'%s' não é um diretório", path))
		return
	}

	// Atualizar diretório atual
	a.currentDir = path

	// Atualizar visualizações
	a.refreshTreeView()
	a.refreshFileView()
}

// navigateTo navega para o diretório especificado
func (a *App) navigateTo(dir string) {
	// Verificar se o diretório existe
	fileInfo, err := os.Stat(dir)
	if err != nil {
		a.showError(fmt.Sprintf("Erro ao acessar diretório: %v", err))
		return
	}
	
	// Verificar se é um diretório
	if !fileInfo.IsDir() {
		a.showError(fmt.Sprintf("'%s' não é um diretório", dir))
		return
	}
	
	// Atualizar diretório atual
	a.currentDir = dir
	
	// Carregar arquivos
	files, err := os.ReadDir(dir)
	if err != nil {
		a.showError(fmt.Sprintf("Erro ao ler diretório: %v", err))
		return
	}
	
	// Filtrar arquivos ocultos se necessário
	var filteredFiles []os.DirEntry
	for _, file := range files {
		if !a.showHidden && strings.HasPrefix(file.Name(), ".") {
			continue
		}
		filteredFiles = append(filteredFiles, file)
	}
	
	// Converter para o formato utils.FileInfo
	var fileInfos []utils.FileInfo
	for _, entry := range filteredFiles {
		info, err := entry.Info()
		if err != nil {
			continue
		}
		
		fileInfo := utils.FileInfo{
			Name:     entry.Name(),
			Path:     filepath.Join(dir, entry.Name()),
			Size:     info.Size(),
			ModTime:  info.ModTime(),
			IsDir:    entry.IsDir(),
			IsHidden: strings.HasPrefix(entry.Name(), "."),
		}
		
		if !entry.IsDir() {
			fileInfo.Extension = strings.ToLower(filepath.Ext(entry.Name()))
		}
		
		fileInfos = append(fileInfos, fileInfo)
	}
	
	// Atualizar visualizações
	a.fileView.UpdateFileList(fileInfos, a.showHidden)
	a.treeView.UpdateTreeView(dir)
	
	// Atualizar barra de status
	a.statusBar.SetText(fmt.Sprintf("Diretório: %s", dir))
	
	// Adicionar ao histórico se não for navegação no histórico
	if a.historyPos < 0 || a.history[a.historyPos] != dir {
		a.addToHistory(dir)
	}
}

// openFile abre o arquivo selecionado
func (a *App) openFile() {
	// Obter arquivo selecionado
	selectedFile := a.fileView.GetSelectedFile()
	if selectedFile == "" {
		return
	}

	// Verificar se é um diretório
	filePath := filepath.Join(a.currentDir, selectedFile)
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		a.showMessage(fmt.Sprintf("Erro ao acessar arquivo: %v", err))
		return
	}

	// Se for um diretório, navegar para ele
	if fileInfo.IsDir() {
		a.navigateTo(filePath)
		return
	}

	// Verificar tipo de arquivo
	isText, err := utils.IsTextFile(filePath)
	if err != nil {
		a.showMessage(fmt.Sprintf("Erro ao verificar tipo de arquivo: %v", err))
		return
	}

	// Se for um arquivo de texto, abrir no visualizador
	if isText {
		a.viewFile()
		return
	}

	// Tentar abrir com o aplicativo padrão
	err = utils.OpenFileWithDefaultApp(filePath)
	if err != nil {
		a.showMessage(fmt.Sprintf("Erro ao abrir arquivo: %v", err))
	}
}

// selectNextFile seleciona o próximo arquivo
func (a *App) selectNextFile() {
	if a.fileView.GetItemCount() == 0 {
		return
	}

	currentItem := a.fileView.GetCurrentItem()
	if currentItem < a.fileView.GetItemCount()-1 {
		a.fileView.Select(currentItem + 1)
	}
}

// enterDirectory entra no diretório selecionado
func (a *App) enterDirectory() {
	// Obter arquivo selecionado
	selectedFile := a.fileView.GetSelectedFile()
	if selectedFile == "" {
		return
	}

	// Verificar se é um diretório
	filePath := filepath.Join(a.currentDir, selectedFile)
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		a.showMessage(fmt.Sprintf("Erro ao acessar arquivo: %v", err))
		return
	}

	// Se for um diretório, navegar para ele
	if fileInfo.IsDir() {
		a.navigateTo(filePath)
	}
}

// navigateBack navega para o diretório anterior no histórico
func (a *App) navigateBack() {
	if a.historyPos <= 0 {
		a.showMessage("Não há histórico anterior")
		return
	}
	
	a.historyPos--
	a.navigateTo(a.history[a.historyPos])
}

// navigateForward navega para o próximo diretório no histórico
func (a *App) navigateForward() {
	if a.historyPos >= len(a.history)-1 {
		a.showMessage("Não há histórico posterior")
		return
	}
	
	a.historyPos++
	a.navigateTo(a.history[a.historyPos])
}

// showDirectoryHistory exibe o histórico de navegação
func (a *App) showDirectoryHistory() {
	if len(a.history) == 0 {
		a.showMessage("Histórico de navegação vazio")
		return
	}
	
	historyList := tview.NewList()
	historyList.SetTitle("Histórico de Navegação").
		SetTitleAlign(tview.AlignCenter).
		SetBorder(true)
	
	// Adicionar diretórios ao histórico (do mais recente para o mais antigo)
	for i := len(a.history) - 1; i >= 0; i-- {
		dir := a.history[i]
		index := i // Capturar índice para o closure
		historyList.AddItem(dir, "", 0, func() {
			a.pages.RemovePage("history")
			a.historyPos = index
			a.navigateTo(a.history[index])
		})
	}
	
	historyList.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEscape {
			a.pages.RemovePage("history")
			return nil
		}
		return event
	})
	
	a.pages.AddPage("history", historyList, true, true)
	a.app.SetFocus(historyList)
}

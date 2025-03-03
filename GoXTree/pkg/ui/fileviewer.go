package ui

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/peder1981/GoXTree/pkg/viewer"
	"github.com/rivo/tview"

	"github.com/gdamore/tcell/v2"
)

// FileViewer representa o componente de visualização de arquivos
type FileViewer struct {
	app      *App
	pages    *tview.Pages
	filePath string
}

// NewFileViewer cria um novo visualizador de arquivos
func NewFileViewer(app *App) *FileViewer {
	return &FileViewer{
		app:   app,
		pages: tview.NewPages(),
	}
}

// ViewFile abre um arquivo no visualizador apropriado
func (fv *FileViewer) ViewFile(filePath string) error {
	// Verificar se o arquivo existe
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return fmt.Errorf("erro ao acessar arquivo: %w", err)
	}

	// Verificar se é um diretório
	if fileInfo.IsDir() {
		return fmt.Errorf("não é possível visualizar diretórios")
	}

	// Armazenar o caminho do arquivo
	fv.filePath = filePath

	// Determinar o tipo de visualizador com base na extensão do arquivo
	ext := strings.ToLower(filepath.Ext(filePath))

	// Criar o visualizador apropriado
	var viewerFlex *tview.Flex

	// Verificar se é uma imagem
	if isImageFile(ext) {
		imageViewer := viewer.NewImageViewer(fv.app.app)
		err := imageViewer.LoadFile(filePath)
		if err != nil {
			return fmt.Errorf("erro ao carregar imagem: %w", err)
		}
		viewerFlex = imageViewer.Show()
	} else if isTextFile(ext) {
		textViewer := viewer.NewTextViewer(fv.app.app)
		err := textViewer.LoadFile(filePath)
		if err != nil {
			return fmt.Errorf("erro ao carregar arquivo de texto: %w", err)
		}
		viewerFlex = textViewer.Show()
	} else {
		// Para outros tipos de arquivo, usar o visualizador hexadecimal
		hexViewer := viewer.NewHexViewer(fv.app.app)
		err := hexViewer.LoadFile(filePath)
		if err != nil {
			return fmt.Errorf("erro ao carregar arquivo para visualização hexadecimal: %w", err)
		}
		viewerFlex = hexViewer.Show()
	}

	// Adicionar o visualizador às páginas
	fv.pages.AddPage("viewer", viewerFlex, true, true)

	// Configurar manipulador de teclas para fechar o visualizador
	fv.app.app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEscape {
			fv.Close()
			return nil
		}
		return event
	})

	// Exibir o visualizador
	fv.app.app.SetRoot(fv.pages, true)

	return nil
}

// showViewer exibe o visualizador
func (fv *FileViewer) showViewer(filePath string) {
	// Ler conteúdo do arquivo
	content, err := os.ReadFile(filePath)
	if err != nil {
		fv.app.showError(fmt.Sprintf("Erro ao ler arquivo: %v", err))
		return
	}

	// Criar visualizador de texto
	textView := tview.NewTextView().
		SetDynamicColors(true).
		SetRegions(true).
		SetWordWrap(true).
		SetChangedFunc(func() {
			fv.app.app.Draw()
		})

	// Configurar texto
	textView.SetText(string(content))

	// Criar layout
	flex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(textView, 0, 1, true).
		AddItem(tview.NewTextView().
			SetTextAlign(tview.AlignCenter).
			SetText("F10/ESC: Sair | ↑/↓: Navegar"), 1, 0, false)

	// Configurar manipulador de eventos
	flex.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyF10, tcell.KeyEscape:
			// Sair do visualizador
			fv.app.pages.RemovePage("viewer")
			return nil
		}
		return event
	})

	// Adicionar página
	fv.app.pages.AddPage("viewer", flex, true, true)
	fv.app.app.SetFocus(flex)
}

// Close fecha o visualizador
func (fv *FileViewer) Close() {
	fv.app.pages.RemovePage("viewer")
}

// isTextFile verifica se um arquivo é de texto com base na extensão
func isTextFile(ext string) bool {
	textExtensions := []string{
		".txt", ".log", ".md", ".json", ".xml", ".html", ".htm", ".css", ".js",
		".go", ".c", ".cpp", ".h", ".hpp", ".py", ".rb", ".pl", ".php", ".java",
		".sh", ".bat", ".cmd", ".ini", ".cfg", ".conf", ".yaml", ".yml", ".toml",
		".csv", ".tsv", ".sql", ".prg", ".tlpp", ".ch",
	}

	for _, textExt := range textExtensions {
		if ext == textExt {
			return true
		}
	}

	return false
}

// isImageFile verifica se um arquivo é uma imagem com base na extensão
func isImageFile(ext string) bool {
	imageExtensions := []string{
		".jpg", ".jpeg", ".png", ".gif", ".bmp", ".tiff", ".webp",
	}

	for _, imgExt := range imageExtensions {
		if ext == imgExt {
			return true
		}
	}

	return false
}

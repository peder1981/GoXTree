package viewer

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/rivo/tview"
)

// FileViewer é um visualizador de arquivos
type FileViewer struct {
	app      *tview.Application
	textView *tview.TextView
	filePath string
}

// NewFileViewer cria um novo visualizador de arquivos
func NewFileViewer(app *tview.Application) *FileViewer {
	// Criar visualizador de texto
	textView := tview.NewTextView().
		SetDynamicColors(true).
		SetRegions(true).
		SetScrollable(true).
		SetWordWrap(true)

	// Configurar borda
	textView.SetBorder(true).
		SetTitle(" Visualizador de Arquivos ").
		SetTitleAlign(tview.AlignLeft)

	// Criar FileViewer
	f := &FileViewer{
		app:      app,
		textView: textView,
		filePath: "",
	}

	return f
}

// LoadFile carrega um arquivo para visualização
func (f *FileViewer) LoadFile(filePath string) error {
	// Verificar se o arquivo existe
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return err
	}

	// Verificar se é um arquivo
	if fileInfo.IsDir() {
		return fmt.Errorf("não é possível visualizar um diretório")
	}

	// Verificar tamanho do arquivo
	if fileInfo.Size() > 10*1024*1024 { // 10MB
		return fmt.Errorf("arquivo muito grande para visualização")
	}

	// Abrir arquivo
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Ler conteúdo
	content, err := io.ReadAll(file)
	if err != nil {
		return err
	}

	// Verificar se é um arquivo binário
	if isBinary(content) {
		return fmt.Errorf("não é possível visualizar um arquivo binário")
	}

	// Atualizar visualizador
	f.filePath = filePath
	f.textView.SetTitle(fmt.Sprintf(" Visualizando: %s ", filepath.Base(filePath)))
	f.textView.SetText(string(content))

	return nil
}

// GetView retorna a visualização
func (f *FileViewer) GetView() *tview.TextView {
	return f.textView
}

// isBinary verifica se um conteúdo é binário
func isBinary(content []byte) bool {
	// Verificar se há caracteres nulos
	for _, b := range content {
		if b == 0 {
			return true
		}
	}

	// Verificar se há muitos caracteres não imprimíveis
	nonPrintable := 0
	for _, b := range content {
		if b < 32 && b != 9 && b != 10 && b != 13 {
			nonPrintable++
		}
	}

	// Se mais de 10% dos caracteres são não imprimíveis, é binário
	return nonPrintable > len(content)/10
}

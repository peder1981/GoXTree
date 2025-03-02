package viewer

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/peder1981/GoXTree/pkg/utils"
	"github.com/rivo/tview"
)

// TextViewer representa um visualizador de texto
type TextViewer struct {
	app       *tview.Application
	textView  *tview.TextView
	statusBar *tview.TextView
	layout    *tview.Flex
	filePath  string
	fileInfo  os.FileInfo
	lines     []string
	lineCount int
}

// NewTextViewer cria um novo visualizador de texto
func NewTextViewer(app *tview.Application) *TextViewer {
	tv := &TextViewer{
		app:       app,
		textView:  tview.NewTextView(),
		statusBar: tview.NewTextView(),
		layout:    tview.NewFlex(),
	}

	// Configurar o visualizador de texto
	tv.textView.SetDynamicColors(true)
	tv.textView.SetRegions(true)
	tv.textView.SetScrollable(true)
	tv.textView.SetBorder(true)
	tv.textView.SetTitle(" Visualizador de Texto ")
	tv.textView.SetTitleAlign(tview.AlignLeft)

	// Configurar a barra de status
	tv.statusBar.SetTextColor(utils.ColorStatusText)
	tv.statusBar.SetBackgroundColor(utils.ColorStatusBar)
	tv.statusBar.SetDynamicColors(true)
	tv.statusBar.SetText("[::b]ESC[-:-:-] Fechar  [::b]↑/↓[-:-:-] Rolar  [::b]PgUp/PgDn[-:-:-] Página  [::b]Home/End[-:-:-] Início/Fim")

	// Configurar o layout
	tv.layout.SetDirection(tview.FlexRow).
		AddItem(tv.textView, 0, 1, true).
		AddItem(tv.statusBar, 1, 1, false)

	// Configurar manipuladores de eventos
	tv.textView.SetInputCapture(tv.handleKeyEvents)

	return tv
}

// LoadFile carrega um arquivo de texto
func (tv *TextViewer) LoadFile(filePath string) error {
	// Obter informações do arquivo
	info, err := os.Stat(filePath)
	if err != nil {
		return err
	}

	// Verificar se é um arquivo regular
	if !info.Mode().IsRegular() {
		return fmt.Errorf("não é um arquivo regular")
	}

	// Verificar o tamanho do arquivo
	if info.Size() > 10*1024*1024 { // 10MB
		return fmt.Errorf("arquivo muito grande para visualização (> 10MB)")
	}

	// Abrir o arquivo
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Ler o conteúdo do arquivo
	content, err := io.ReadAll(file)
	if err != nil {
		return err
	}

	// Armazenar informações do arquivo
	tv.filePath = filePath
	tv.fileInfo = info

	// Processar o conteúdo
	tv.lines = strings.Split(string(content), "\n")
	tv.lineCount = len(tv.lines)

	// Atualizar o título
	fileName := filepath.Base(filePath)
	tv.textView.SetTitle(fmt.Sprintf(" %s (%s) ", fileName, utils.FormatFileSize(info.Size())))

	// Exibir o conteúdo
	tv.displayContent(string(content))

	return nil
}

// displayContent exibe o conteúdo no visualizador
func (tv *TextViewer) displayContent(content string) {
	tv.textView.SetText(content)
	tv.textView.ScrollToBeginning()
}

// handleKeyEvents manipula eventos de teclado
func (tv *TextViewer) handleKeyEvents(event *tcell.EventKey) *tcell.EventKey {
	switch event.Key() {
	case tcell.KeyEscape:
		// Fechar o visualizador
		return nil
	}

	return event
}

// Show exibe o visualizador
func (tv *TextViewer) Show() *tview.Flex {
	return tv.layout
}

// HexViewer representa um visualizador hexadecimal
type HexViewer struct {
	app       *tview.Application
	textView  *tview.TextView
	statusBar *tview.TextView
	layout    *tview.Flex
	filePath  string
	fileInfo  os.FileInfo
}

// NewHexViewer cria um novo visualizador hexadecimal
func NewHexViewer(app *tview.Application) *HexViewer {
	hv := &HexViewer{
		app:       app,
		textView:  tview.NewTextView(),
		statusBar: tview.NewTextView(),
		layout:    tview.NewFlex(),
	}

	// Configurar o visualizador hexadecimal
	hv.textView.SetDynamicColors(true)
	hv.textView.SetRegions(true)
	hv.textView.SetScrollable(true)
	hv.textView.SetBorder(true)
	hv.textView.SetTitle(" Visualizador Hexadecimal ")
	hv.textView.SetTitleAlign(tview.AlignLeft)

	// Configurar a barra de status
	hv.statusBar.SetTextColor(utils.ColorStatusText)
	hv.statusBar.SetBackgroundColor(utils.ColorStatusBar)
	hv.statusBar.SetDynamicColors(true)
	hv.statusBar.SetText("[::b]ESC[-:-:-] Fechar  [::b]↑/↓[-:-:-] Rolar  [::b]PgUp/PgDn[-:-:-] Página  [::b]Home/End[-:-:-] Início/Fim")

	// Configurar o layout
	hv.layout.SetDirection(tview.FlexRow).
		AddItem(hv.textView, 0, 1, true).
		AddItem(hv.statusBar, 1, 1, false)

	// Configurar manipuladores de eventos
	hv.textView.SetInputCapture(hv.handleKeyEvents)

	return hv
}

// LoadFile carrega um arquivo para visualização hexadecimal
func (hv *HexViewer) LoadFile(filePath string) error {
	// Obter informações do arquivo
	info, err := os.Stat(filePath)
	if err != nil {
		return err
	}

	// Verificar se é um arquivo regular
	if !info.Mode().IsRegular() {
		return fmt.Errorf("não é um arquivo regular")
	}

	// Verificar o tamanho do arquivo
	if info.Size() > 1*1024*1024 { // 1MB
		return fmt.Errorf("arquivo muito grande para visualização hexadecimal (> 1MB)")
	}

	// Abrir o arquivo
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Ler o conteúdo do arquivo
	content, err := io.ReadAll(file)
	if err != nil {
		return err
	}

	// Armazenar informações do arquivo
	hv.filePath = filePath
	hv.fileInfo = info

	// Atualizar o título
	fileName := filepath.Base(filePath)
	hv.textView.SetTitle(fmt.Sprintf(" %s (%s) - Hex ", fileName, utils.FormatFileSize(info.Size())))

	// Exibir o conteúdo hexadecimal
	hv.displayHexContent(content)

	return nil
}

// displayHexContent exibe o conteúdo em formato hexadecimal
func (hv *HexViewer) displayHexContent(content []byte) {
	var hexOutput strings.Builder

	// Exibir 16 bytes por linha
	for i := 0; i < len(content); i += 16 {
		// Endereço
		hexOutput.WriteString(fmt.Sprintf("[yellow]%08x[white]  ", i))

		// Bytes em hexadecimal
		for j := 0; j < 16; j++ {
			if i+j < len(content) {
				hexOutput.WriteString(fmt.Sprintf("%02x ", content[i+j]))
			} else {
				hexOutput.WriteString("   ")
			}

			// Espaço extra no meio
			if j == 7 {
				hexOutput.WriteString(" ")
			}
		}

		// Separador
		hexOutput.WriteString(" |")

		// Caracteres ASCII
		for j := 0; j < 16; j++ {
			if i+j < len(content) {
				b := content[i+j]
				if b >= 32 && b <= 126 { // Caracteres imprimíveis
					hexOutput.WriteString(string(b))
				} else {
					hexOutput.WriteString(".")
				}
			} else {
				hexOutput.WriteString(" ")
			}
		}

		hexOutput.WriteString("|\n")
	}

	hv.textView.SetText(hexOutput.String())
	hv.textView.ScrollToBeginning()
}

// handleKeyEvents manipula eventos de teclado
func (hv *HexViewer) handleKeyEvents(event *tcell.EventKey) *tcell.EventKey {
	switch event.Key() {
	case tcell.KeyEscape:
		// Fechar o visualizador
		return nil
	}

	return event
}

// Show exibe o visualizador
func (hv *HexViewer) Show() *tview.Flex {
	return hv.layout
}

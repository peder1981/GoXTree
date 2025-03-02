package viewer

import (
	"fmt"
	"image"
	"image/color"
	_ "image/gif"  // Registro de decodificadores de imagem
	_ "image/jpeg" // Registro de decodificadores de imagem
	_ "image/png"  // Registro de decodificadores de imagem
	"os"
	"path/filepath"

	"github.com/gdamore/tcell/v2"
	"github.com/peder1981/GoXTree/pkg/utils"
	"github.com/rivo/tview"
)

// ImageViewer representa um visualizador de imagens em ASCII art
type ImageViewer struct {
	app       *tview.Application
	textView  *tview.TextView
	statusBar *tview.TextView
	layout    *tview.Flex
	filePath  string
	fileInfo  os.FileInfo
	image     image.Image
	width     int
	height    int
}

// NewImageViewer cria um novo visualizador de imagens
func NewImageViewer(app *tview.Application) *ImageViewer {
	iv := &ImageViewer{
		app:       app,
		textView:  tview.NewTextView(),
		statusBar: tview.NewTextView(),
		layout:    tview.NewFlex(),
	}

	// Configurar o visualizador de texto
	iv.textView.SetDynamicColors(true)
	iv.textView.SetRegions(true)
	iv.textView.SetScrollable(true)
	iv.textView.SetBorder(true)
	iv.textView.SetTitle(" Visualizador de Imagem ")
	iv.textView.SetTitleAlign(tview.AlignLeft)
	iv.textView.SetBorderColor(utils.ColorViewerBorder)
	iv.textView.SetTitleColor(utils.ColorViewerTitle)
	iv.textView.SetTextColor(utils.ColorViewerText)
	
	// Configurar a barra de status
	iv.statusBar.SetTextColor(utils.ColorStatusText)
	iv.statusBar.SetBackgroundColor(utils.ColorStatusBar)
	iv.statusBar.SetDynamicColors(true)
	iv.statusBar.SetText("[::b]ESC[-:-:-] Fechar  [::b]↑/↓[-:-:-] Rolar  [::b]PgUp/PgDn[-:-:-] Página  [::b]Home/End[-:-:-] Início/Fim")
	
	// Configurar o layout
	iv.layout.SetDirection(tview.FlexRow).
		AddItem(iv.textView, 0, 1, true).
		AddItem(iv.statusBar, 1, 1, false)
	
	// Configurar manipuladores de eventos
	iv.textView.SetInputCapture(iv.handleKeyEvents)

	return iv
}

// LoadFile carrega uma imagem
func (iv *ImageViewer) LoadFile(filePath string) error {
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
	
	// Decodificar a imagem
	img, _, err := image.Decode(file)
	if err != nil {
		return fmt.Errorf("erro ao decodificar imagem: %w", err)
	}
	
	// Armazenar informações
	iv.filePath = filePath
	iv.fileInfo = info
	iv.image = img
	iv.width = img.Bounds().Dx()
	iv.height = img.Bounds().Dy()
	
	// Atualizar o título
	fileName := filepath.Base(filePath)
	iv.textView.SetTitle(fmt.Sprintf(" %s (%dx%d) - %s ", fileName, iv.width, iv.height, utils.FormatFileSize(info.Size())))
	
	// Exibir a imagem
	iv.displayImageAsAscii()
	
	return nil
}

// displayImageAsAscii exibe a imagem como ASCII art
func (iv *ImageViewer) displayImageAsAscii() {
	if iv.image == nil {
		return
	}
	
	// Obter dimensões da imagem
	bounds := iv.image.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()
	
	// Calcular fator de escala para ajustar ao terminal
	// Assumindo que cada caractere tem proporção 2:1 (altura:largura)
	maxWidth := 80
	maxHeight := 40
	
	scaleX := float64(width) / float64(maxWidth)
	scaleY := float64(height) / float64(maxHeight) * 2 // Compensar pela proporção do caractere
	
	scale := scaleX
	if scaleY > scale {
		scale = scaleY
	}
	
	if scale < 1 {
		scale = 1
	}
	
	// Calcular dimensões escaladas
	scaledWidth := int(float64(width) / scale)
	scaledHeight := int(float64(height) / scale)
	
	// Caracteres para representar diferentes níveis de brilho
	// Do mais escuro para o mais claro
	asciiChars := " .:-=+*#%@"
	
	// Construir a representação ASCII
	var content string
	
	for y := 0; y < scaledHeight; y++ {
		for x := 0; x < scaledWidth; x++ {
			// Calcular coordenadas na imagem original
			origX := int(float64(x) * scale)
			origY := int(float64(y) * scale)
			
			// Obter cor do pixel
			pixel := iv.image.At(origX+bounds.Min.X, origY+bounds.Min.Y)
			
			// Converter para escala de cinza
			r, g, b, _ := color.RGBAModel.Convert(pixel).RGBA()
			gray := (r + g + b) / 3
			
			// Normalizar para 0-255
			gray = gray / 257
			
			// Mapear para um caractere ASCII
			charIndex := int(float64(gray) * float64(len(asciiChars)-1) / 255.0)
			if charIndex >= len(asciiChars) {
				charIndex = len(asciiChars) - 1
			}
			
			content += string(asciiChars[charIndex])
		}
		content += "\n"
	}
	
	// Exibir no visualizador
	iv.textView.SetText(content)
	iv.textView.ScrollToBeginning()
}

// handleKeyEvents manipula eventos de teclado
func (iv *ImageViewer) handleKeyEvents(event *tcell.EventKey) *tcell.EventKey {
	switch event.Key() {
	case tcell.KeyEscape:
		// O fechamento é tratado pelo FileViewer
		return event
	}
	
	return event
}

// Show exibe o visualizador
func (iv *ImageViewer) Show() *tview.Flex {
	return iv.layout
}

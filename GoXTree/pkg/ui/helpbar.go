package ui

import (
	"github.com/peder1981/GoXTree/pkg/utils"
	"github.com/rivo/tview"
)

// HelpBar representa a barra de ajuda
type HelpBar struct {
	app     *App
	helpBar *tview.TextView
}

// NewHelpBar cria uma nova barra de ajuda
func NewHelpBar(app *App) *HelpBar {
	hb := &HelpBar{
		app:     app,
		helpBar: tview.NewTextView(),
	}

	// Configurar a barra de ajuda
	hb.helpBar.SetTextColor(utils.ColorStatusText)
	hb.helpBar.SetBackgroundColor(utils.ColorStatusBar)
	hb.helpBar.SetDynamicColors(true)
	hb.helpBar.SetText("[::b]TAB[-:-:-] Alternar painel  [::b]ENTER[-:-:-] Selecionar  [::b]BACKSPACE[-:-:-] Diretório pai  [::b]ESPAÇO[-:-:-] Marcar  [::b]Ctrl+H[-:-:-] Ocultos")

	return hb
}

// SetHelp define o texto da barra de ajuda
func (hb *HelpBar) SetHelp(text string) {
	hb.helpBar.SetText(text)
}

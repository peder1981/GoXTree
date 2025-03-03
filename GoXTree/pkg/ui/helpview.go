package ui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// HelpView representa a visualização de ajuda
type HelpView struct {
	app      *App
	helpView *tview.TextView
	pages    *tview.Pages
}

// NewHelpView cria uma nova visualização de ajuda
func NewHelpView(app *App) *HelpView {
	hv := &HelpView{
		app:      app,
		helpView: tview.NewTextView(),
		pages:    tview.NewPages(),
	}

	// Configurar aparência
	hv.helpView.SetBorder(true)
	hv.helpView.SetTitle(" Ajuda do GoXTree ")
	hv.helpView.SetTitleAlign(tview.AlignCenter)
	hv.helpView.SetDynamicColors(true)
	hv.helpView.SetScrollable(true)
	hv.helpView.SetWordWrap(true)

	// Configurar texto de ajuda
	hv.helpView.SetText(hv.getHelpText())

	// Adicionar manipulador de teclas
	hv.helpView.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEscape || event.Key() == tcell.KeyF1 {
			hv.Close()
			return nil
		}
		return event
	})

	// Adicionar à página
	hv.pages.AddPage("help", hv.helpView, true, true)

	return hv
}

// Show exibe a visualização de ajuda
func (hv *HelpView) Show() {
	hv.app.app.SetRoot(hv.pages, true)
}

// Close fecha a visualização de ajuda
func (hv *HelpView) Close() {
	hv.app.pages.RemovePage("help")
}

// getHelpText retorna o texto de ajuda formatado
func (hv *HelpView) getHelpText() string {
	return `[yellow]GoXTree v1.0.0 - Ajuda[white]

[green]== Navegação ==[white]
[yellow]↑/↓[white] - Mover seleção para cima/baixo
[yellow]←/→[white] - Alternar entre painéis
[yellow]Enter[white] - Abrir diretório ou arquivo
[yellow]Backspace[white] - Ir para o diretório pai
[yellow]Tab[white] - Alternar entre árvore e lista de arquivos

[green]== Teclas de Função ==[white]
[yellow]F1[white] - Exibir esta ajuda
[yellow]F2[white] - Abrir menu
[yellow]F3[white] - Visualizar arquivo
[yellow]F4[white] - Editar arquivo
[yellow]F5[white] - Copiar arquivo(s)
[yellow]F6[white] - Mover arquivo(s)
[yellow]F7[white] - Criar diretório
[yellow]F8[white] - Excluir arquivo(s)
[yellow]F9[white] - Comprimir arquivo(s)
[yellow]F10[white] - Sair

[green]== Atalhos ==[white]
[yellow]Ctrl+A[white] - Selecionar todos os arquivos
[yellow]Ctrl+F[white] - Buscar arquivos
[yellow]Ctrl+G[white] - Ir para diretório
[yellow]Ctrl+H[white] - Mostrar/ocultar arquivos ocultos
[yellow]Ctrl+R[white] - Atualizar visualização
[yellow]Espaço[white] - Selecionar/desselecionar arquivo
[yellow]Esc[white] - Cancelar operação atual

[green]== Modos de Visualização ==[white]
- Árvore: Exibe a estrutura de diretórios em formato de árvore
- Lista: Exibe arquivos e diretórios em formato de lista
- Detalhes: Exibe arquivos com informações detalhadas

[green]== Sobre o GoXTree ==[white]
GoXTree é uma reimplementação moderna do XTree Gold, um popular gerenciador de arquivos da era DOS. Desenvolvido em Go, o GoXTree combina a simplicidade e eficiência do XTree original com recursos modernos.

[yellow]Pressione ESC ou F1 para fechar esta ajuda[white]`
}

// showHelp exibe a ajuda
func (hv *HelpView) showHelp() {
	// Criar visualizador de texto
	textView := tview.NewTextView().
		SetDynamicColors(true).
		SetRegions(true).
		SetWordWrap(true).
		SetChangedFunc(func() {
			hv.app.app.Draw()
		})
	
	// Configurar texto
	textView.SetText(hv.getHelpText())
	
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
			// Sair da ajuda
			hv.app.pages.RemovePage("help")
			return nil
		}
		return event
	})
	
	// Adicionar página
	hv.app.pages.AddPage("help", flex, true, true)
	hv.app.app.SetFocus(flex)
}

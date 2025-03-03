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
	return `[yellow]GoXTree - Gerenciador de Arquivos Retrô[white]

[yellow]Navegação:[white]
  [green]Setas[white]        - Mover cursor
  [green]Enter[white]        - Entrar no diretório / Abrir arquivo
  [green]Backspace[white]    - Voltar ao diretório pai
  [green]Tab[white]          - Alternar entre árvore e lista de arquivos
  [green]ESC[white]          - Voltar / Fechar janela atual

[yellow]Teclas de Função:[white]
  [green]F1[white]           - Mostrar ajuda
  [green]F2[white]           - Renomear arquivo/diretório
  [green]F3[white]           - Buscar arquivos
  [green]F4[white]           - Busca avançada
  [green]F7[white]           - Criar diretório
  [green]F8[white]           - Excluir arquivo/diretório
  [green]F9[white]           - Sincronizar diretórios
  [green]F10[white]          - Sair

[yellow]Atalhos Ctrl+Letra:[white]
  [green]Ctrl+A[white]       - Selecionar todos os arquivos
  [green]Ctrl+D[white]       - Desmarcar todos os arquivos
  [green]Ctrl+F[white]       - Buscar arquivo
  [green]Ctrl+G[white]       - Ir para diretório
  [green]Ctrl+H[white]       - Alternar arquivos ocultos
  [green]Ctrl+R[white]       - Atualizar visualização
  [green]Ctrl+S[white]       - Selecionar/desmarcar arquivo atual
  [green]Ctrl+C[white]       - Comparar arquivos selecionados
  [green]Ctrl+V[white]       - Visualizar arquivo
  [green]Ctrl+E[white]       - Editar arquivo
  [green]Ctrl+Y[white]       - Sincronizar diretórios
  [green]Ctrl+I[white]       - Informações do sistema

[yellow]Seleção de Arquivos:[white]
  - Use [green]Ctrl+S[white] para selecionar/desmarcar o arquivo atual
  - Use [green]Ctrl+A[white] para selecionar todos os arquivos
  - Use [green]Ctrl+D[white] para desmarcar todos os arquivos
  - Os arquivos selecionados são destacados em [cyan]ciano[white] com fundo azul escuro
  - Use [green]Ctrl+C[white] para comparar dois arquivos selecionados

[yellow]Tema Retrô:[white]
  - O GoXTree usa um tema retrô inspirado nos gerenciadores de arquivos DOS
  - Diferentes tipos de arquivos são destacados com cores distintas:
    * [blue]Diretórios[white] - Azul
    * [green]Executáveis[white] - Verde
    * [magenta]Arquivos compactados[white] - Magenta
    * [cyan]Arquivos de código[white] - Ciano
    * [red]Imagens e PDFs[white] - Vermelho
    * [yellow]Apresentações[white] - Amarelo
    * [white]Arquivos de texto[white] - Branco
    * [gray]Arquivos ocultos[white] - Cinza

[yellow]Visualização e Edição:[white]
  - Use [green]Ctrl+V[white] para visualizar o conteúdo do arquivo atual
  - Use [green]Ctrl+E[white] para editar o arquivo atual no editor interno
  - Use [green]ESC[white] para sair do visualizador/editor

[yellow]Comparação de Arquivos:[white]
  - Selecione exatamente dois arquivos usando [green]Ctrl+S[white]
  - Pressione [green]Ctrl+C[white] para comparar os arquivos
  - As diferenças são destacadas em cores:
    * [green]+ Texto adicionado[white]
    * [red]- Texto removido[white]
    * Texto sem alteração

[yellow]Sobre o GoXTree:[white]
  GoXTree é um gerenciador de arquivos moderno inspirado no XTree Gold,
  implementado em Go com interface de terminal usando as bibliotecas
  tcell e tview. Ele combina a simplicidade e eficiência dos gerenciadores
  de arquivos clássicos com recursos modernos.

Pressione [green]ESC[white] para fechar esta ajuda.`
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

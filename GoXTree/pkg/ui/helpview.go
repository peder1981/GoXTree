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

// GetRect implementa a interface tview.Primitive
func (hv *HelpView) GetRect() (int, int, int, int) {
	return hv.helpView.GetRect()
}

// SetRect implementa a interface tview.Primitive
func (hv *HelpView) SetRect(x, y, width, height int) {
	hv.helpView.SetRect(x, y, width, height)
}

// Draw implementa a interface tview.Primitive
func (hv *HelpView) Draw(screen tcell.Screen) {
	hv.helpView.Draw(screen)
}

// InputHandler implementa a interface tview.Primitive
func (hv *HelpView) InputHandler() func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
	return hv.helpView.InputHandler()
}

// Focus implementa a interface tview.Primitive
func (hv *HelpView) Focus(delegate func(p tview.Primitive)) {
	hv.helpView.Focus(delegate)
}

// Blur implementa a interface tview.Primitive
func (hv *HelpView) Blur() {
	hv.helpView.Blur()
}

// HasFocus implementa a interface tview.Primitive
func (hv *HelpView) HasFocus() bool {
	return hv.helpView.HasFocus()
}

// MouseHandler implementa a interface tview.Primitive
func (hv *HelpView) MouseHandler() func(action tview.MouseAction, event *tcell.EventMouse, setFocus func(p tview.Primitive)) (consumed bool, capture tview.Primitive) {
	return hv.helpView.MouseHandler()
}

// PasteHandler implementa a interface tview.Primitive
func (hv *HelpView) PasteHandler() func(text string) {
	// Retornar uma função que encapsula o PasteHandler do helpView
	return func(text string) {
		if handler := hv.helpView.PasteHandler(); handler != nil {
			handler(text, nil)
		}
	}
}

// Close fecha a visualização de ajuda
func (hv *HelpView) Close() {
	hv.app.app.SetRoot(hv.app.pages, true)
}

// getHelpText retorna o texto de ajuda formatado
func (hv *HelpView) getHelpText() string {
	return `[yellow]GoXTree - Ajuda[white]

[yellow]Navegação:[white]
  - Use as [green]setas[white] para navegar pela lista de arquivos
  - [green]Enter[white] para entrar em um diretório ou abrir um arquivo
  - [green]Backspace[white] para voltar ao diretório pai
  - [green]Home/End[white] para ir para o início/fim da lista
  - [green]PgUp/PgDn[white] para navegar páginas

[yellow]Seleção:[white]
  - [green]Espaço[white] para selecionar/deselecionar um arquivo
  - [green]Ins[white] para selecionar um arquivo e mover para o próximo
  - [green]Ctrl+A[white] para selecionar todos os arquivos
  - [green]Ctrl+N[white] para deselecionar todos os arquivos
  - [green]Ctrl+I[white] para inverter a seleção

[yellow]Operações de Arquivo:[white]
  - [green]F5[white] para copiar arquivos selecionados
  - [green]F6[white] para mover arquivos selecionados
  - [green]F7[white] para criar um novo diretório
  - [green]F8/Del[white] para excluir arquivos selecionados
  - [green]F9[white] para criar um novo arquivo

[yellow]Visualização:[white]
  - [green]F3[white] para alternar entre visualização em árvore e lista
  - [green]F4[white] para alternar entre visualização detalhada e simples
  - [green]F2[white] para alternar entre ordenação por nome, tamanho, data
  - [green]h[white] para alternar exibição de arquivos ocultos
  - [green]r[white] para atualizar a visualização atual

[yellow]Busca:[white]
  - [green]Ctrl+F[white] para buscar arquivos por nome
  - [green]Ctrl+G[white] para buscar arquivos por conteúdo
  - [green]F[white] para busca rápida na lista atual
  - [green]n[white] para ir para o próximo resultado da busca
  - [green]N[white] para ir para o resultado anterior da busca

[yellow]Esquema de Cores:[white]
  * [blue]Diretórios[white] - Azul
  * [green]Executáveis[white] - Verde
  * [magenta]Arquivos compactados[white] - Magenta
  * [cyan]Arquivos de código[white] - Ciano
  * [red]Imagens e PDFs[white] - Vermelho
  * [yellow]Apresentações[white] - Amarelo
  * [white]Arquivos de texto[white] - Branco
  * [gray]Arquivos ocultos[white] - Cinza

[yellow]Visualização e Edição:[white]
  - Use [green]Alt+V[white] para visualizar o conteúdo do arquivo atual
  - Use [green]Alt+E[white] para editar o arquivo atual no editor interno
  - Use [green]ESC[white] para sair do visualizador/editor

[yellow]Comparação de Arquivos:[white]
  - Selecione exatamente dois arquivos
  - Pressione [green]Alt+C[white] para comparar os arquivos
  - As diferenças são destacadas em cores:
    * [green]+ Texto adicionado[white]
    * [red]- Texto removido[white]
    * Texto sem alteração

[yellow]Configuração:[white]
  - As configurações são salvas em ~/.gxtree/config.json
  - Você pode personalizar:
    * Tema (retro, modern, dark, light)
    * Exibição de arquivos ocultos
    * Atalhos de teclado personalizados
    * Esquema de cores

[yellow]Sobre o GoXTree:[white]
  GoXTree é um gerenciador de arquivos moderno inspirado no XTree Gold,
  implementado em Go com interface de terminal usando as bibliotecas
  tcell e tview. Ele combina a simplicidade e eficiência dos gerenciadores
  de arquivos clássicos com recursos modernos.

Pressione [green]ESC[white] para fechar esta ajuda.`
}

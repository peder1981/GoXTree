package ui

import (
	"path/filepath"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// showHelpDialog exibe o diálogo de ajuda
func (a *App) showHelpDialog() {
	// Criar texto de ajuda
	helpText := `
GoXTree - Gerenciador de Arquivos

Teclas de Navegação:
  Tab         - Alternar entre visualizações
  Enter       - Abrir arquivo/diretório
  Backspace   - Voltar ao diretório anterior
  Setas       - Navegar na lista de arquivos

Teclas de Função:
  F1  - Ajuda
  F2  - Menu principal
  F3  - Busca simples
  F4  - Busca avançada
  F5  - Copiar arquivo
  F6  - Mover arquivo
  F7  - Criar diretório
  F8  - Excluir arquivo
  F9  - Sincronizar diretórios
  F10 - Sair

Teclas de Controle:
  Ctrl+F - Buscar arquivo
  Ctrl+G - Ir para diretório
  Ctrl+R - Atualizar visualizações
  Ctrl+H - Mostrar/ocultar arquivos ocultos
  Ctrl+S - Selecionar arquivo
  Ctrl+A - Selecionar todos os arquivos
  Ctrl+D - Desmarcar todos os arquivos
  Ctrl+C - Comparar arquivos selecionados
  Ctrl+V - Visualizar arquivo
  Ctrl+E - Editar arquivo
  Ctrl+Y - Sincronizar diretórios

Visualização em Árvore:
  + ou > - Expandir nó
  - ou < - Colapsar nó

Comportamento da tecla ESC:
  - Na tela principal: Volta ao diretório anterior ou pergunta se deseja sair
  - Em diálogos: Fecha o diálogo atual e volta à tela anterior
`

	// Criar texto
	textView := tview.NewTextView().
		SetText(helpText).
		SetDynamicColors(true).
		SetRegions(true).
		SetWordWrap(true).
		SetScrollable(true)

	// Configurar borda
	textView.SetBorder(true).
		SetTitle(" Ajuda ").
		SetTitleAlign(tview.AlignLeft)

	// Configurar manipulador de teclas
	textView.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEscape || event.Key() == tcell.KeyEnter {
			a.pages.RemovePage("helpDialog")
			return nil
		}
		return event
	})

	// Exibir diálogo
	a.pages.AddPage("helpDialog", a.modal(textView, 60, 20), true, true)
}

// showGotoDialog exibe o diálogo para ir para um diretório
func (a *App) showGotoDialog() {
	// Criar formulário
	form := tview.NewForm()
	form.AddInputField("Diretório:", "", 40, nil, nil)
	form.AddButton("OK", func() {
		// Obter diretório
		dir := form.GetFormItem(0).(*tview.InputField).GetText()
		if dir == "" {
			return
		}

		// Expandir caminho
		if strings.HasPrefix(dir, "~") {
			homeDir, err := a.getHomeDir()
			if err == nil {
				dir = filepath.Join(homeDir, dir[1:])
			}
		}

		// Verificar se é um caminho relativo
		if !filepath.IsAbs(dir) {
			dir = filepath.Join(a.currentDir, dir)
		}

		// Navegar para o diretório
		a.navigateTo(dir)
		a.pages.RemovePage("gotoDialog")
	})
	form.AddButton("Cancelar", func() {
		a.pages.RemovePage("gotoDialog")
	})

	// Configurar borda
	form.SetBorder(true).
		SetTitle(" Ir para Diretório ").
		SetTitleAlign(tview.AlignLeft)

	// Configurar manipulador de teclas
	form.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEscape {
			a.pages.RemovePage("gotoDialog")
			return nil
		}
		return event
	})

	// Exibir diálogo
	a.pages.AddPage("gotoDialog", a.modal(form, 50, 7), true, true)
}

// showMenu exibe o menu principal
func (a *App) showMenu(title string, options []string, callback func(int, string)) {
	// Criar lista
	list := tview.NewList()
	for i, option := range options {
		list.AddItem(option, "", rune('1'+i), nil)
	}

	// Configurar seleção
	list.SetSelectedFunc(func(index int, mainText, secondaryText string, shortcut rune) {
		a.pages.RemovePage("menuDialog")
		callback(index, mainText)
	})

	// Configurar borda
	list.SetBorder(true).
		SetTitle(" " + title + " ").
		SetTitleAlign(tview.AlignLeft)

	// Configurar manipulador de teclas
	list.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEscape {
			a.pages.RemovePage("menuDialog")
			return nil
		}
		return event
	})

	// Exibir diálogo
	a.pages.AddPage("menuDialog", a.modal(list, 30, len(options)+2), true, true)
}

// modal centraliza um primitivo na tela
func (a *App) modal(p tview.Primitive, width, height int) tview.Primitive {
	return tview.NewFlex().
		AddItem(nil, 0, 1, false).
		AddItem(
			tview.NewFlex().
				SetDirection(tview.FlexRow).
				AddItem(nil, 0, 1, false).
				AddItem(p, height, 1, true).
				AddItem(nil, 0, 1, false),
			width, 1, true).
		AddItem(nil, 0, 1, false)
}

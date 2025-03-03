package ui

import (
	"fmt"
	"strings"
	"path/filepath"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// showTextDialog exibe um diálogo com texto
func (a *App) showTextDialog(title, content string) {
	textView := tview.NewTextView().
		SetText(content).
		SetDynamicColors(true).
		SetScrollable(true).
		SetWordWrap(true)
	textView.SetBorder(true).
		SetTitle(fmt.Sprintf(" %s ", title)).
		SetTitleAlign(tview.AlignLeft)

	// Configurar manipulador de teclas
	textView.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEscape {
			a.pages.RemovePage("textDialog")
			return nil
		}
		return event
	})

	a.pages.AddPage("textDialog", textView, true, true)
}

// showInputDialog exibe um diálogo de entrada de texto
func (a *App) showInputDialog(title, defaultText string, callback func(string)) {
	// Criar formulário
	form := tview.NewForm()
	form.SetTitle(title).
		SetTitleAlign(tview.AlignCenter).
		SetBorder(true)
	
	// Adicionar campo de entrada
	form.AddInputField("", defaultText, 40, nil, nil)
	
	// Adicionar botões
	form.AddButton("OK", func() {
		a.pages.RemovePage("inputDialog")
		callback(form.GetFormItem(0).(*tview.InputField).GetText())
	})
	
	form.AddButton("Cancelar", func() {
		a.pages.RemovePage("inputDialog")
	})
	
	// Configurar manipulador de eventos
	form.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEscape {
			a.pages.RemovePage("inputDialog")
			return nil
		} else if event.Key() == tcell.KeyEnter {
			a.pages.RemovePage("inputDialog")
			callback(form.GetFormItem(0).(*tview.InputField).GetText())
			return nil
		}
		return event
	})
	
	// Adicionar página
	a.pages.AddPage("inputDialog", form, true, true)
	a.app.SetFocus(form)
}

// showInputDialogWithValue exibe um diálogo de entrada com valor inicial
func (a *App) showInputDialogWithValue(title, initialValue string, callback func(string)) {
	form := tview.NewForm()
	form.AddInputField(title, initialValue, 0, nil, nil)
	form.AddButton("OK", func() {
		value := form.GetFormItem(0).(*tview.InputField).GetText()
		a.pages.RemovePage("input")
		callback(value)
	})
	form.AddButton("Cancelar", func() {
		a.pages.RemovePage("input")
	})
	form.SetBorder(true).SetTitle(" Entrada ").SetTitleAlign(tview.AlignCenter)
	form.SetCancelFunc(func() {
		a.pages.RemovePage("input")
	})
	
	// Adicionar página
	a.pages.AddPage("input", form, true, true)
	a.app.SetFocus(form)
}

// showGoToDialog exibe o diálogo para ir para um diretório
func (a *App) showGoToDialog() {
	a.showInputDialog("Ir para diretório:", a.currentDir, func(path string) {
		if path == "" {
			return
		}
		
		// Expandir caminho
		if strings.HasPrefix(path, "~") {
			homeDir, err := a.getHomeDir()
			if err != nil {
				a.showError(fmt.Sprintf("Erro ao obter diretório home: %v", err))
				return
			}
			path = filepath.Join(homeDir, path[1:])
		}
		
		// Navegar para o diretório
		a.NavigateToDirectory(path)
	})
}

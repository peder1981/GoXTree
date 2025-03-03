package editor

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/gdamore/tcell/v2"
	"github.com/peder1981/GoXTree/pkg/utils"
	"github.com/rivo/tview"
)

// TextEditor representa um editor de texto simples
type TextEditor struct {
	app       *tview.Application
	textArea  *tview.TextArea
	statusBar *tview.TextView
	layout    *tview.Flex
	filePath  string
	fileInfo  os.FileInfo
	modified  bool
	onSave    func()
	onClose   func()
}

// NewTextEditor cria um novo editor de texto
func NewTextEditor(app *tview.Application) *TextEditor {
	te := &TextEditor{
		app:       app,
		textArea:  tview.NewTextArea(),
		statusBar: tview.NewTextView(),
		layout:    tview.NewFlex(),
		modified:  false,
	}

	// Configurar a área de texto
	te.textArea.SetBorder(true)
	te.textArea.SetTitle(" Editor de Texto ")
	te.textArea.SetTitleAlign(tview.AlignLeft)
	te.textArea.SetChangedFunc(func() {
		te.modified = true
		te.updateStatusBar()
	})

	// Configurar a barra de status
	te.statusBar.SetTextColor(utils.ColorStatusText)
	te.statusBar.SetBackgroundColor(utils.ColorStatusBar)
	te.statusBar.SetDynamicColors(true)
	te.updateStatusBar()

	// Configurar o layout
	te.layout.SetDirection(tview.FlexRow).
		AddItem(te.textArea, 0, 1, true).
		AddItem(te.statusBar, 1, 1, false)

	// Configurar manipuladores de eventos
	te.textArea.SetInputCapture(te.handleKeyEvents)

	return te
}

// LoadFile carrega um arquivo para edição
func (te *TextEditor) LoadFile(filePath string) error {
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
		return fmt.Errorf("arquivo muito grande para edição (> 10MB)")
	}

	// Ler o conteúdo do arquivo
	content, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	// Armazenar informações do arquivo
	te.filePath = filePath
	te.fileInfo = info
	te.modified = false

	// Definir o texto
	te.textArea.SetText(string(content), true)

	// Atualizar o título
	fileName := filepath.Base(filePath)
	te.textArea.SetTitle(fmt.Sprintf(" Editor de Texto - %s ", fileName))

	// Atualizar a barra de status
	te.updateStatusBar()

	return nil
}

// SaveFile salva o arquivo
func (te *TextEditor) SaveFile() error {
	// Verificar se há um arquivo aberto
	if te.filePath == "" {
		return fmt.Errorf("nenhum arquivo aberto")
	}

	// Obter o texto
	text := te.textArea.GetText()

	// Salvar o arquivo
	err := os.WriteFile(te.filePath, []byte(text), 0644)
	if err != nil {
		return err
	}

	// Atualizar o estado
	te.modified = false
	te.updateStatusBar()

	// Executar callback de salvamento
	if te.onSave != nil {
		te.onSave()
	}

	return nil
}

// SetOnSave define o callback a ser executado quando o arquivo for salvo
func (te *TextEditor) SetOnSave(callback func()) {
	te.onSave = callback
}

// SetOnClose define o callback a ser executado quando o editor for fechado
func (te *TextEditor) SetOnClose(callback func()) {
	te.onClose = callback
}

// updateStatusBar atualiza a barra de status
func (te *TextEditor) updateStatusBar() {
	var status string

	// Verificar se há um arquivo aberto
	if te.filePath != "" {
		fileName := filepath.Base(te.filePath)
		status = fmt.Sprintf(" %s ", fileName)

		// Indicar se o arquivo foi modificado
		if te.modified {
			status += " [::b]*[-:-:-] "
		}
	} else {
		status = " Novo Arquivo "
	}

	// Adicionar instruções
	status += " | [::b]Ctrl+S[-:-:-] Salvar | [::b]Ctrl+Q[-:-:-] Sair"

	// Atualizar a barra de status
	te.statusBar.SetText(status)
}

// handleKeyEvents manipula eventos de teclado
func (te *TextEditor) handleKeyEvents(event *tcell.EventKey) *tcell.EventKey {
	// Verificar combinações de teclas
	if event.Key() == tcell.KeyCtrlS {
		// Salvar o arquivo
		err := te.SaveFile()
		if err != nil {
			// Exibir mensagem de erro na barra de status
			te.statusBar.SetText(fmt.Sprintf(" Erro ao salvar: %s | [::b]Ctrl+S[-:-:-] Salvar | [::b]Ctrl+Q[-:-:-] Sair", err))
		}
		return nil
	} else if event.Key() == tcell.KeyCtrlQ || event.Key() == tcell.KeyEscape {
		// Verificar se há modificações não salvas
		if te.modified {
			// Exibir diálogo de confirmação
			te.showConfirmDialog("Arquivo modificado. Deseja salvar antes de sair?", func(save bool) {
				if save {
					// Salvar o arquivo
					err := te.SaveFile()
					if err != nil {
						// Exibir mensagem de erro na barra de status
						te.statusBar.SetText(fmt.Sprintf(" Erro ao salvar: %s | [::b]Ctrl+S[-:-:-] Salvar | [::b]Ctrl+Q[-:-:-] Sair", err))
						return
					}
				}

				// Fechar o editor
				if te.onClose != nil {
					te.onClose()
				}
			})
		} else {
			// Fechar o editor
			if te.onClose != nil {
				te.onClose()
			}
		}
		return nil
	}

	return event
}

// showConfirmDialog exibe um diálogo de confirmação
func (te *TextEditor) showConfirmDialog(message string, callback func(bool)) {
	modal := tview.NewModal().
		SetText(message).
		AddButtons([]string{"Sim", "Não", "Cancelar"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			if buttonLabel == "Sim" {
				callback(true)
			} else if buttonLabel == "Não" {
				callback(false)
			}
			// Se for "Cancelar", não faz nada
			te.app.SetRoot(te.layout, true)
		})

	te.app.SetRoot(modal, true)
}

// Show exibe o editor
func (te *TextEditor) Show() *tview.Flex {
	return te.layout
}

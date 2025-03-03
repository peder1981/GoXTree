package ui

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/gdamore/tcell/v2"
	"github.com/peder1981/GoXTree/pkg/editor"
	"github.com/rivo/tview"
)

// FileEditor representa o componente de edição de arquivos
type FileEditor struct {
	app      *App
	pages    *tview.Pages
	filePath string
}

// NewFileEditor cria um novo editor de arquivos
func NewFileEditor(app *App) *FileEditor {
	return &FileEditor{
		app:   app,
		pages: tview.NewPages(),
	}
}

// EditFile abre um arquivo no editor apropriado
func (fe *FileEditor) EditFile(filePath string) error {
	// Verificar se o arquivo existe
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		// Se o arquivo não existe, verificar se o diretório existe
		dir := filepath.Dir(filePath)
		_, err := os.Stat(dir)
		if err != nil {
			return fmt.Errorf("diretório inválido: %w", err)
		}
	} else if fileInfo.IsDir() {
		return fmt.Errorf("não é possível editar diretórios")
	}

	// Armazenar o caminho do arquivo
	fe.filePath = filePath

	// Por enquanto, só temos um editor de texto
	// No futuro, podemos adicionar editores específicos para outros tipos de arquivo
	textEditor := editor.NewTextEditor(fe.app.app)

	// Configurar callbacks
	textEditor.SetOnClose(func() {
		fe.Close()
	})

	textEditor.SetOnSave(func() {
		// Atualizar a visualização após salvar
		fe.app.refreshView()
	})

	// Carregar o arquivo, se existir
	if fileInfo != nil {
		err = textEditor.LoadFile(filePath)
		if err != nil {
			return fmt.Errorf("erro ao carregar arquivo: %w", err)
		}
	}

	// Adicionar o editor às páginas
	fe.pages.AddPage("editor", textEditor.Show(), true, true)

	// Exibir o editor
	fe.app.app.SetRoot(fe.pages, true)

	return nil
}

// Close fecha o editor
func (fe *FileEditor) Close() {
	fe.app.pages.RemovePage("editor")
}

// showEditor exibe o editor
func (fe *FileEditor) showEditor(filePath string) {
	// Ler conteúdo do arquivo
	content, err := os.ReadFile(filePath)
	if err != nil {
		fe.app.showError(fmt.Sprintf("Erro ao ler arquivo: %v", err))
		return
	}

	// Criar editor de texto
	textView := tview.NewTextView().
		SetDynamicColors(true).
		SetRegions(true).
		SetWordWrap(true).
		SetChangedFunc(func() {
			fe.app.app.Draw()
		})

	// Configurar texto
	textView.SetText(string(content))

	// Criar layout
	flex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(textView, 0, 1, true).
		AddItem(tview.NewTextView().
			SetTextAlign(tview.AlignCenter).
			SetText("F2: Salvar | F10: Sair"), 1, 0, false)

	// Configurar manipulador de eventos
	flex.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyF2:
			// Salvar arquivo
			fe.saveFile(filePath, textView.GetText(false))
			return nil
		case tcell.KeyF10, tcell.KeyEscape:
			// Sair do editor
			fe.app.pages.RemovePage("editor")
			return nil
		}
		return event
	})

	// Adicionar página
	fe.app.pages.AddPage("editor", flex, true, true)
	fe.app.app.SetFocus(flex)
}

// saveFile salva o conteúdo em um arquivo
func (fe *FileEditor) saveFile(filePath, content string) {
	// Criar diretório pai se necessário
	dir := filepath.Dir(filePath)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.MkdirAll(dir, 0755)
		if err != nil {
			fe.app.showError(fmt.Sprintf("Erro ao criar diretório: %v", err))
			return
		}
	}

	// Salvar conteúdo no arquivo
	err := os.WriteFile(filePath, []byte(content), 0644)
	if err != nil {
		fe.app.showError(fmt.Sprintf("Erro ao salvar arquivo: %v", err))
		return
	}

	// Exibir mensagem de sucesso
	fe.app.showMessage(fmt.Sprintf("Arquivo salvo: %s", filePath))
}

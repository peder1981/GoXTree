package ui

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/peder1981/GoXTree/pkg/utils"
	"github.com/rivo/tview"
)

// FileView representa a visualização de arquivos
type FileView struct {
	app        *App
	fileList   *tview.Table
	currentDir string
	files      []string
	showHidden bool
	itemCount  int
}

// NewFileView cria uma nova visualização de arquivos
func NewFileView(app *App) *FileView {
	// Criar tabela
	table := tview.NewTable().
		SetBorders(false).
		SetSelectable(true, false)

	// Configurar cabeçalho
	table.SetCell(0, 0, tview.NewTableCell("Nome").SetTextColor(tcell.ColorYellow).SetAlign(tview.AlignLeft).SetSelectable(false))
	table.SetCell(0, 1, tview.NewTableCell("Tamanho").SetTextColor(tcell.ColorYellow).SetAlign(tview.AlignRight).SetSelectable(false))
	table.SetCell(0, 2, tview.NewTableCell("Data").SetTextColor(tcell.ColorYellow).SetAlign(tview.AlignLeft).SetSelectable(false))
	table.SetCell(0, 3, tview.NewTableCell("Permissões").SetTextColor(tcell.ColorYellow).SetAlign(tview.AlignLeft).SetSelectable(false))

	// Configurar borda
	table.SetBorder(true).
		SetTitle(" Arquivos ").
		SetTitleAlign(tview.AlignLeft)

	// Criar FileView
	f := &FileView{
		app:        app,
		fileList:   table,
		currentDir: "",
		files:      make([]string, 0),
		showHidden: false,
		itemCount:  0,
	}

	// Configurar manipulador de seleção
	table.SetSelectedFunc(func(row, column int) {
		if row == 0 {
			return
		}
		f.app.openFile()
	})

	return f
}

// SetCurrentDir define o diretório atual
func (f *FileView) SetCurrentDir(dir string) error {
	// Verificar se o diretório existe
	fileInfo, err := os.Stat(dir)
	if err != nil {
		return err
	}

	// Verificar se é um diretório
	if !fileInfo.IsDir() {
		return fmt.Errorf("%s não é um diretório", dir)
	}

	f.currentDir = dir
	f.Refresh()
	return nil
}

// SetShowHidden define se arquivos ocultos são exibidos
func (f *FileView) SetShowHidden(show bool) {
	f.showHidden = show
	f.Refresh()
}

// Refresh atualiza a visualização de arquivos
func (f *FileView) Refresh() {
	// Limpar tabela
	f.fileList.Clear()

	// Configurar cabeçalho
	f.fileList.SetCell(0, 0, tview.NewTableCell("Nome").SetTextColor(tcell.ColorYellow).SetAlign(tview.AlignLeft).SetSelectable(false))
	f.fileList.SetCell(0, 1, tview.NewTableCell("Tamanho").SetTextColor(tcell.ColorYellow).SetAlign(tview.AlignRight).SetSelectable(false))
	f.fileList.SetCell(0, 2, tview.NewTableCell("Data").SetTextColor(tcell.ColorYellow).SetAlign(tview.AlignLeft).SetSelectable(false))
	f.fileList.SetCell(0, 3, tview.NewTableCell("Permissões").SetTextColor(tcell.ColorYellow).SetAlign(tview.AlignLeft).SetSelectable(false))

	// Verificar se o diretório atual existe
	if f.currentDir == "" {
		return
	}

	// Inicializar a lista de arquivos
	f.files = make([]string, 0)
	row := 1

	// Adicionar entrada para o diretório pai (..) se não estiver na raiz
	parentDir := strings.TrimSpace(filepath.Dir(f.currentDir))
	if parentDir != f.currentDir && parentDir != "" {
		// Adicionar ".." à lista de arquivos
		f.files = append(f.files, "..")

		// Adicionar linha à tabela
		f.fileList.SetCell(row, 0, tview.NewTableCell("📁 ..").SetTextColor(tcell.ColorBlue).SetAlign(tview.AlignLeft))
		f.fileList.SetCell(row, 1, tview.NewTableCell("<DIR>").SetTextColor(tcell.ColorBlue).SetAlign(tview.AlignRight))
		f.fileList.SetCell(row, 2, tview.NewTableCell("").SetTextColor(tcell.ColorBlue).SetAlign(tview.AlignLeft))
		f.fileList.SetCell(row, 3, tview.NewTableCell("rwx------").SetTextColor(tcell.ColorBlue).SetAlign(tview.AlignLeft))

		row++
	}

	// Listar arquivos
	files, err := utils.ListFiles(f.currentDir, f.showHidden)
	if err != nil {
		return
	}

	// Ordenar arquivos (diretórios primeiro, depois por nome)
	sort.Slice(files, func(i, j int) bool {
		if files[i].IsDir != files[j].IsDir {
			return files[i].IsDir
		}
		return strings.ToLower(files[i].Name) < strings.ToLower(files[j].Name)
	})

	// Adicionar arquivos à tabela
	for _, file := range files {
		// Adicionar nome do arquivo à lista
		f.files = append(f.files, file.Name)

		// Configurar cor
		color := tcell.ColorWhite
		if file.IsDir {
			color = tcell.ColorLightBlue
		}

		// Adicionar ícone
		icon := utils.GetFileIcon(file)
		name := fmt.Sprintf("%s %s", icon, file.Name)

		// Adicionar linha à tabela
		f.fileList.SetCell(row, 0, tview.NewTableCell(name).SetTextColor(color).SetAlign(tview.AlignLeft))
		if file.IsDir {
			f.fileList.SetCell(row, 1, tview.NewTableCell("<DIR>").SetTextColor(color).SetAlign(tview.AlignRight))
		} else {
			f.fileList.SetCell(row, 1, tview.NewTableCell(utils.FormatFileSize(file.Size)).SetTextColor(color).SetAlign(tview.AlignRight))
		}
		f.fileList.SetCell(row, 2, tview.NewTableCell(file.ModTime.Format("02/01/2006 15:04:05")).SetTextColor(color).SetAlign(tview.AlignLeft))
		f.fileList.SetCell(row, 3, tview.NewTableCell("rwx------").SetTextColor(color).SetAlign(tview.AlignLeft))

		row++
	}

	// Selecionar primeira linha
	if row > 1 {
		f.fileList.Select(1, 0)
	}

	// Atualizar contador de itens
	f.itemCount = row - 1
}

// UpdateFileList atualiza a lista de arquivos
func (f *FileView) UpdateFileList(files []utils.FileInfo, showHidden bool) {
	// Limpar tabela
	f.fileList.Clear()

	// Adicionar cabeçalho
	f.fileList.SetCell(0, 0, tview.NewTableCell("Nome").SetTextColor(tcell.ColorYellow).SetSelectable(false))
	f.fileList.SetCell(0, 1, tview.NewTableCell("Tamanho").SetTextColor(tcell.ColorYellow).SetSelectable(false))
	f.fileList.SetCell(0, 2, tview.NewTableCell("Data").SetTextColor(tcell.ColorYellow).SetSelectable(false))
	f.fileList.SetCell(0, 3, tview.NewTableCell("Permissões").SetTextColor(tcell.ColorYellow).SetSelectable(false))

	// Adicionar diretório pai
	f.fileList.SetCell(1, 0, tview.NewTableCell("..").SetTextColor(tcell.ColorBlue))
	f.fileList.SetCell(1, 1, tview.NewTableCell(""))
	f.fileList.SetCell(1, 2, tview.NewTableCell(""))
	f.fileList.SetCell(1, 3, tview.NewTableCell(""))

	// Filtrar arquivos ocultos
	var visibleFiles []utils.FileInfo
	for _, file := range files {
		if !showHidden && strings.HasPrefix(file.Name, ".") {
			continue
		}
		visibleFiles = append(visibleFiles, file)
	}

	// Ordenar arquivos (diretórios primeiro, depois por nome)
	sort.Slice(visibleFiles, func(i, j int) bool {
		if visibleFiles[i].IsDir != visibleFiles[j].IsDir {
			return visibleFiles[i].IsDir
		}
		return strings.ToLower(visibleFiles[i].Name) < strings.ToLower(visibleFiles[j].Name)
	})

	// Adicionar arquivos
	for i, file := range visibleFiles {
		row := i + 2 // +2 para o cabeçalho e o diretório pai

		// Nome
		nameCell := tview.NewTableCell(file.Name)
		if file.IsDir {
			nameCell.SetTextColor(tcell.ColorBlue)
		}
		f.fileList.SetCell(row, 0, nameCell)

		// Tamanho
		sizeText := ""
		if !file.IsDir {
			sizeText = utils.FormatFileSize(file.Size)
		}
		f.fileList.SetCell(row, 1, tview.NewTableCell(sizeText))

		// Data
		dateText := file.ModTime.Format("02/01/2006 15:04:05")
		f.fileList.SetCell(row, 2, tview.NewTableCell(dateText))

		// Permissões
		permText := "rwx------" // Valor padrão simplificado
		f.fileList.SetCell(row, 3, tview.NewTableCell(permText))
	}

	// Atualizar contagem de itens
	f.itemCount = len(visibleFiles) + 1 // +1 para o diretório pai

	// Selecionar primeiro item
	if f.itemCount > 0 {
		f.fileList.Select(1, 0)
	}

	// Atualizar título
	f.fileList.SetTitle(fmt.Sprintf(" Arquivos (%d) ", f.itemCount))
}

// GetSelectedFile retorna o arquivo selecionado
func (f *FileView) GetSelectedFile() string {
	row, _ := f.fileList.GetSelection()
	if row <= 0 || row >= len(f.files) {
		return ""
	}
	return f.files[row-1]
}

// GetItemCount retorna o número de itens na lista
func (f *FileView) GetItemCount() int {
	return len(f.files)
}

// GetCurrentItem retorna o índice do item atual
func (f *FileView) GetCurrentItem() int {
	row, _ := f.fileList.GetSelection()
	return row - 1
}

// Select seleciona um item pelo índice
func (f *FileView) Select(index int) {
	if index >= 0 && index < len(f.files) {
		f.fileList.Select(index+1, 0)
	}
}

// SelectFile seleciona um arquivo pelo nome
func (f *FileView) SelectFile(fileName string) bool {
	for i, file := range f.files {
		if file == fileName {
			f.fileList.Select(i+1, 0)
			return true
		}
	}
	return false
}

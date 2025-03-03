package ui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// Cores retrô inspiradas em terminais DOS
var (
	ColorBackground = tcell.ColorBlack
	ColorText       = tcell.ColorWhite
	ColorBorder     = tcell.ColorBlue
	ColorTitle      = tcell.ColorYellow
	ColorHighlight  = tcell.ColorGreen
	ColorSelected   = tcell.ColorWhite
	ColorHeader     = tcell.ColorYellow
	ColorError      = tcell.ColorRed
	ColorHelp       = tcell.ColorGreen
)

// Caracteres ASCII para bordas
const (
	BorderHorizontal      = '-'
	BorderVertical        = '|'
	BorderTopLeft         = '+'
	BorderTopRight        = '+'
	BorderBottomLeft      = '+'
	BorderBottomRight     = '+'
	BorderLeft            = '|'
	BorderRight           = '|'
	BorderTop             = '-'
	BorderBottom          = '-'
	TreeVerticalLine      = '|'
	TreeHorizontalLine    = '-'
	TreeCorner            = '+'
	TreeContinueCorner    = '+'
	TreeEndCorner         = '+'
	TreeDirectory         = '['
	TreeDirectoryOpen     = '['
	TreeDirectoryClosed   = ']'
	TreeFile              = ' '
)

// ApplyRetroTheme aplica o tema retrô à aplicação
func ApplyRetroTheme(app *tview.Application) {
	// Definir tema global
	tview.Styles = tview.Theme{
		PrimitiveBackgroundColor:    ColorBackground,
		ContrastBackgroundColor:     ColorBackground,
		MoreContrastBackgroundColor: ColorBackground,
		BorderColor:                 ColorBorder,
		TitleColor:                  ColorTitle,
		GraphicsColor:               ColorBorder,
		PrimaryTextColor:            ColorText,
		SecondaryTextColor:          ColorText,
		TertiaryTextColor:           ColorText,
		InverseTextColor:            ColorBackground,
		ContrastSecondaryTextColor:  ColorHighlight,
	}

	// Definir caracteres ASCII para bordas
	tview.Borders.Horizontal = BorderHorizontal
	tview.Borders.Vertical = BorderVertical
	tview.Borders.TopLeft = BorderTopLeft
	tview.Borders.TopRight = BorderTopRight
	tview.Borders.BottomLeft = BorderBottomLeft
	tview.Borders.BottomRight = BorderBottomRight
	tview.Borders.Left = BorderLeft
	tview.Borders.Right = BorderRight
	tview.Borders.Top = BorderTop
	tview.Borders.Bottom = BorderBottom

	// Definir caracteres ASCII para a árvore
	tview.TreeGraphics.VerticalLine = TreeVerticalLine
	tview.TreeGraphics.HorizontalLine = TreeHorizontalLine
	tview.TreeGraphics.CornerLine = TreeCorner
	tview.TreeGraphics.ContinueCorner = TreeContinueCorner
	tview.TreeGraphics.EndCorner = TreeEndCorner
	tview.TreeGraphics.Directory = TreeDirectory
	tview.TreeGraphics.DirectoryOpen = TreeDirectoryOpen
	tview.TreeGraphics.DirectoryClosed = TreeDirectoryClosed
	tview.TreeGraphics.File = TreeFile
}

// ApplyRetroThemeToApp aplica o tema retrô a todos os componentes da aplicação
func ApplyRetroThemeToApp(a *App) {
	// Aplicar tema global
	ApplyRetroTheme(a.app)

	// Personalizar TreeView
	a.treeView.TreeView.SetBackgroundColor(ColorBackground)
	a.treeView.TreeView.SetBorderColor(ColorBorder)
	a.treeView.TreeView.SetTitleColor(ColorTitle)
	a.treeView.TreeView.SetTitle(" Diretórios ")
	a.treeView.TreeView.SetBorder(true)

	// Personalizar FileView
	a.fileView.fileList.SetBackgroundColor(ColorBackground)
	a.fileView.fileList.SetBorderColor(ColorBorder)
	a.fileView.fileList.SetTitleColor(ColorTitle)
	a.fileView.fileList.SetTitle(" Arquivos ")
	a.fileView.fileList.SetBorder(true)
	
	// Personalizar cabeçalhos da tabela de arquivos
	a.fileView.fileList.SetCell(0, 0, tview.NewTableCell("Nome").
		SetTextColor(ColorHeader).SetAlign(tview.AlignLeft).SetSelectable(false))
	a.fileView.fileList.SetCell(0, 1, tview.NewTableCell("Tamanho").
		SetTextColor(ColorHeader).SetAlign(tview.AlignRight).SetSelectable(false))
	a.fileView.fileList.SetCell(0, 2, tview.NewTableCell("Data").
		SetTextColor(ColorHeader).SetAlign(tview.AlignLeft).SetSelectable(false))
	a.fileView.fileList.SetCell(0, 3, tview.NewTableCell("Permissões").
		SetTextColor(ColorHeader).SetAlign(tview.AlignLeft).SetSelectable(false))

	// Personalizar StatusBar
	a.statusBar.statusBar.SetBackgroundColor(ColorBackground)
	a.statusBar.statusBar.SetTextColor(ColorText)
	a.statusBar.statusBar.SetBorderColor(ColorBorder)
	a.statusBar.statusBar.SetBorder(true)

	// Personalizar MenuBar
	a.menuBar.menuBar.SetBackgroundColor(ColorBackground)
	a.menuBar.menuBar.SetTextColor(ColorText)
}

// ApplyModernThemeToApp aplica o tema moderno a todos os componentes da aplicação
func ApplyModernThemeToApp(a *App) {
	// Definir cores modernas
	modernBackground := tcell.ColorReset
	modernText := tcell.ColorBlack
	modernBorder := tcell.ColorBlue
	modernTitle := tcell.ColorBlue
	modernHighlight := tcell.ColorGreen
	modernSelected := tcell.ColorWhite
	modernHeader := tcell.ColorBlue

	// Aplicar tema global
	tview.Styles = tview.Theme{
		PrimitiveBackgroundColor:    modernBackground,
		ContrastBackgroundColor:     modernBackground,
		MoreContrastBackgroundColor: modernBackground,
		BorderColor:                 modernBorder,
		TitleColor:                  modernTitle,
		GraphicsColor:               modernBorder,
		PrimaryTextColor:            modernText,
		SecondaryTextColor:          modernText,
		TertiaryTextColor:           modernText,
		InverseTextColor:            modernBackground,
		ContrastSecondaryTextColor:  modernHighlight,
	}

	// Usar caracteres Unicode para bordas
	tview.Borders.Horizontal = '─'
	tview.Borders.Vertical = '│'
	tview.Borders.TopLeft = '┌'
	tview.Borders.TopRight = '┐'
	tview.Borders.BottomLeft = '└'
	tview.Borders.BottomRight = '┘'
	tview.Borders.Left = '│'
	tview.Borders.Right = '│'
	tview.Borders.Top = '─'
	tview.Borders.Bottom = '─'

	// Usar caracteres Unicode para a árvore
	tview.TreeGraphics.VerticalLine = '│'
	tview.TreeGraphics.HorizontalLine = '─'
	tview.TreeGraphics.CornerLine = '├'
	tview.TreeGraphics.ContinueCorner = '├'
	tview.TreeGraphics.EndCorner = '└'
	tview.TreeGraphics.Directory = '📁'
	tview.TreeGraphics.DirectoryOpen = '📂'
	tview.TreeGraphics.DirectoryClosed = '📁'
	tview.TreeGraphics.File = '📄'

	// Personalizar TreeView
	a.treeView.TreeView.SetBackgroundColor(modernBackground)
	a.treeView.TreeView.SetBorderColor(modernBorder)
	a.treeView.TreeView.SetTitleColor(modernTitle)
	a.treeView.TreeView.SetTitle(" Diretórios ")
	a.treeView.TreeView.SetBorder(true)

	// Personalizar FileView
	a.fileView.fileList.SetBackgroundColor(modernBackground)
	a.fileView.fileList.SetBorderColor(modernBorder)
	a.fileView.fileList.SetTitleColor(modernTitle)
	a.fileView.fileList.SetTitle(" Arquivos ")
	a.fileView.fileList.SetBorder(true)
	
	// Personalizar cabeçalhos da tabela de arquivos
	a.fileView.fileList.SetCell(0, 0, tview.NewTableCell("Nome").
		SetTextColor(modernHeader).SetAlign(tview.AlignLeft).SetSelectable(false))
	a.fileView.fileList.SetCell(0, 1, tview.NewTableCell("Tamanho").
		SetTextColor(modernHeader).SetAlign(tview.AlignRight).SetSelectable(false))
	a.fileView.fileList.SetCell(0, 2, tview.NewTableCell("Data").
		SetTextColor(modernHeader).SetAlign(tview.AlignLeft).SetSelectable(false))
	a.fileView.fileList.SetCell(0, 3, tview.NewTableCell("Permissões").
		SetTextColor(modernHeader).SetAlign(tview.AlignLeft).SetSelectable(false))

	// Personalizar StatusBar
	a.statusBar.statusBar.SetBackgroundColor(modernBackground)
	a.statusBar.statusBar.SetTextColor(modernText)
	a.statusBar.statusBar.SetBorderColor(modernBorder)
	a.statusBar.statusBar.SetBorder(true)

	// Personalizar MenuBar
	a.menuBar.menuBar.SetBackgroundColor(modernBackground)
	a.menuBar.menuBar.SetTextColor(modernText)
}

// ApplyDarkThemeToApp aplica o tema escuro a todos os componentes da aplicação
func ApplyDarkThemeToApp(a *App) {
	// Definir cores escuras
	darkBackground := tcell.ColorBlack
	darkText := tcell.ColorWhite
	darkBorder := tcell.ColorDarkBlue
	darkTitle := tcell.ColorLightBlue
	darkHighlight := tcell.ColorLightGreen
	darkSelected := tcell.ColorWhite
	darkHeader := tcell.ColorLightBlue

	// Aplicar tema global
	tview.Styles = tview.Theme{
		PrimitiveBackgroundColor:    darkBackground,
		ContrastBackgroundColor:     darkBackground,
		MoreContrastBackgroundColor: darkBackground,
		BorderColor:                 darkBorder,
		TitleColor:                  darkTitle,
		GraphicsColor:               darkBorder,
		PrimaryTextColor:            darkText,
		SecondaryTextColor:          darkText,
		TertiaryTextColor:           darkText,
		InverseTextColor:            darkBackground,
		ContrastSecondaryTextColor:  darkHighlight,
	}

	// Usar caracteres Unicode para bordas
	tview.Borders.Horizontal = '─'
	tview.Borders.Vertical = '│'
	tview.Borders.TopLeft = '┌'
	tview.Borders.TopRight = '┐'
	tview.Borders.BottomLeft = '└'
	tview.Borders.BottomRight = '┘'
	tview.Borders.Left = '│'
	tview.Borders.Right = '│'
	tview.Borders.Top = '─'
	tview.Borders.Bottom = '─'

	// Usar caracteres Unicode para a árvore
	tview.TreeGraphics.VerticalLine = '│'
	tview.TreeGraphics.HorizontalLine = '─'
	tview.TreeGraphics.CornerLine = '├'
	tview.TreeGraphics.ContinueCorner = '├'
	tview.TreeGraphics.EndCorner = '└'
	tview.TreeGraphics.Directory = '📁'
	tview.TreeGraphics.DirectoryOpen = '📂'
	tview.TreeGraphics.DirectoryClosed = '📁'
	tview.TreeGraphics.File = '📄'

	// Personalizar TreeView
	a.treeView.TreeView.SetBackgroundColor(darkBackground)
	a.treeView.TreeView.SetBorderColor(darkBorder)
	a.treeView.TreeView.SetTitleColor(darkTitle)
	a.treeView.TreeView.SetTitle(" Diretórios ")
	a.treeView.TreeView.SetBorder(true)

	// Personalizar FileView
	a.fileView.fileList.SetBackgroundColor(darkBackground)
	a.fileView.fileList.SetBorderColor(darkBorder)
	a.fileView.fileList.SetTitleColor(darkTitle)
	a.fileView.fileList.SetTitle(" Arquivos ")
	a.fileView.fileList.SetBorder(true)
	
	// Personalizar cabeçalhos da tabela de arquivos
	a.fileView.fileList.SetCell(0, 0, tview.NewTableCell("Nome").
		SetTextColor(darkHeader).SetAlign(tview.AlignLeft).SetSelectable(false))
	a.fileView.fileList.SetCell(0, 1, tview.NewTableCell("Tamanho").
		SetTextColor(darkHeader).SetAlign(tview.AlignRight).SetSelectable(false))
	a.fileView.fileList.SetCell(0, 2, tview.NewTableCell("Data").
		SetTextColor(darkHeader).SetAlign(tview.AlignLeft).SetSelectable(false))
	a.fileView.fileList.SetCell(0, 3, tview.NewTableCell("Permissões").
		SetTextColor(darkHeader).SetAlign(tview.AlignLeft).SetSelectable(false))

	// Personalizar StatusBar
	a.statusBar.statusBar.SetBackgroundColor(darkBackground)
	a.statusBar.statusBar.SetTextColor(darkText)
	a.statusBar.statusBar.SetBorderColor(darkBorder)
	a.statusBar.statusBar.SetBorder(true)

	// Personalizar MenuBar
	a.menuBar.menuBar.SetBackgroundColor(darkBackground)
	a.menuBar.menuBar.SetTextColor(darkText)
}

// ApplyLightThemeToApp aplica o tema claro a todos os componentes da aplicação
func ApplyLightThemeToApp(a *App) {
	// Definir cores claras
	lightBackground := tcell.ColorWhite
	lightText := tcell.ColorBlack
	lightBorder := tcell.ColorBlue
	lightTitle := tcell.ColorBlue
	lightHighlight := tcell.ColorGreen
	lightSelected := tcell.ColorBlack
	lightHeader := tcell.ColorBlue

	// Aplicar tema global
	tview.Styles = tview.Theme{
		PrimitiveBackgroundColor:    lightBackground,
		ContrastBackgroundColor:     lightBackground,
		MoreContrastBackgroundColor: lightBackground,
		BorderColor:                 lightBorder,
		TitleColor:                  lightTitle,
		GraphicsColor:               lightBorder,
		PrimaryTextColor:            lightText,
		SecondaryTextColor:          lightText,
		TertiaryTextColor:           lightText,
		InverseTextColor:            lightBackground,
		ContrastSecondaryTextColor:  lightHighlight,
	}

	// Usar caracteres Unicode para bordas
	tview.Borders.Horizontal = '─'
	tview.Borders.Vertical = '│'
	tview.Borders.TopLeft = '┌'
	tview.Borders.TopRight = '┐'
	tview.Borders.BottomLeft = '└'
	tview.Borders.BottomRight = '┘'
	tview.Borders.Left = '│'
	tview.Borders.Right = '│'
	tview.Borders.Top = '─'
	tview.Borders.Bottom = '─'

	// Usar caracteres Unicode para a árvore
	tview.TreeGraphics.VerticalLine = '│'
	tview.TreeGraphics.HorizontalLine = '─'
	tview.TreeGraphics.CornerLine = '├'
	tview.TreeGraphics.ContinueCorner = '├'
	tview.TreeGraphics.EndCorner = '└'
	tview.TreeGraphics.Directory = '📁'
	tview.TreeGraphics.DirectoryOpen = '📂'
	tview.TreeGraphics.DirectoryClosed = '📁'
	tview.TreeGraphics.File = '📄'

	// Personalizar TreeView
	a.treeView.TreeView.SetBackgroundColor(lightBackground)
	a.treeView.TreeView.SetBorderColor(lightBorder)
	a.treeView.TreeView.SetTitleColor(lightTitle)
	a.treeView.TreeView.SetTitle(" Diretórios ")
	a.treeView.TreeView.SetBorder(true)

	// Personalizar FileView
	a.fileView.fileList.SetBackgroundColor(lightBackground)
	a.fileView.fileList.SetBorderColor(lightBorder)
	a.fileView.fileList.SetTitleColor(lightTitle)
	a.fileView.fileList.SetTitle(" Arquivos ")
	a.fileView.fileList.SetBorder(true)
	
	// Personalizar cabeçalhos da tabela de arquivos
	a.fileView.fileList.SetCell(0, 0, tview.NewTableCell("Nome").
		SetTextColor(lightHeader).SetAlign(tview.AlignLeft).SetSelectable(false))
	a.fileView.fileList.SetCell(0, 1, tview.NewTableCell("Tamanho").
		SetTextColor(lightHeader).SetAlign(tview.AlignRight).SetSelectable(false))
	a.fileView.fileList.SetCell(0, 2, tview.NewTableCell("Data").
		SetTextColor(lightHeader).SetAlign(tview.AlignLeft).SetSelectable(false))
	a.fileView.fileList.SetCell(0, 3, tview.NewTableCell("Permissões").
		SetTextColor(lightHeader).SetAlign(tview.AlignLeft).SetSelectable(false))

	// Personalizar StatusBar
	a.statusBar.statusBar.SetBackgroundColor(lightBackground)
	a.statusBar.statusBar.SetTextColor(lightText)
	a.statusBar.statusBar.SetBorderColor(lightBorder)
	a.statusBar.statusBar.SetBorder(true)

	// Personalizar MenuBar
	a.menuBar.menuBar.SetBackgroundColor(lightBackground)
	a.menuBar.menuBar.SetTextColor(lightText)
}

// GetRetroColorScheme retorna um esquema de cores retrô para diferentes tipos de arquivos
func GetRetroColorScheme() map[string]tcell.Color {
	return map[string]tcell.Color{
		"dir":      tcell.ColorYellow,
		"exe":      tcell.ColorGreen,
		"zip":      tcell.ColorMagenta,
		"tar":      tcell.ColorMagenta,
		"gz":       tcell.ColorMagenta,
		"rar":      tcell.ColorMagenta,
		"7z":       tcell.ColorMagenta,
		"txt":      tcell.ColorWhite,
		"md":       tcell.ColorWhite,
		"go":       tcell.ColorCyan,
		"c":        tcell.ColorCyan,
		"cpp":      tcell.ColorCyan,
		"h":        tcell.ColorCyan,
		"py":       tcell.ColorCyan,
		"js":       tcell.ColorCyan,
		"html":     tcell.ColorCyan,
		"css":      tcell.ColorCyan,
		"json":     tcell.ColorCyan,
		"xml":      tcell.ColorCyan,
		"yaml":     tcell.ColorCyan,
		"yml":      tcell.ColorCyan,
		"toml":     tcell.ColorCyan,
		"ini":      tcell.ColorCyan,
		"conf":     tcell.ColorCyan,
		"sh":       tcell.ColorGreen,
		"bat":      tcell.ColorGreen,
		"cmd":      tcell.ColorGreen,
		"ps1":      tcell.ColorGreen,
		"jpg":      tcell.ColorRed,
		"jpeg":     tcell.ColorRed,
		"png":      tcell.ColorRed,
		"gif":      tcell.ColorRed,
		"bmp":      tcell.ColorRed,
		"svg":      tcell.ColorRed,
		"mp3":      tcell.ColorMagenta,
		"wav":      tcell.ColorMagenta,
		"ogg":      tcell.ColorMagenta,
		"mp4":      tcell.ColorMagenta,
		"avi":      tcell.ColorMagenta,
		"mkv":      tcell.ColorMagenta,
		"pdf":      tcell.ColorRed,
		"doc":      tcell.ColorCyan,
		"docx":     tcell.ColorCyan,
		"xls":      tcell.ColorGreen,
		"xlsx":     tcell.ColorGreen,
		"ppt":      tcell.ColorYellow,
		"pptx":     tcell.ColorYellow,
		"hidden":   tcell.ColorGray,
		"default":  tcell.ColorWhite,
	}
}

// GetModernColorScheme retorna um esquema de cores moderno para diferentes tipos de arquivos
func GetModernColorScheme() map[string]tcell.Color {
	return map[string]tcell.Color{
		"dir":      tcell.ColorBlue,
		"exe":      tcell.ColorGreen,
		"zip":      tcell.ColorPurple,
		"tar":      tcell.ColorPurple,
		"gz":       tcell.ColorPurple,
		"rar":      tcell.ColorPurple,
		"7z":       tcell.ColorPurple,
		"txt":      tcell.ColorBlack,
		"md":       tcell.ColorBlack,
		"go":       tcell.ColorTeal,
		"c":        tcell.ColorTeal,
		"cpp":      tcell.ColorTeal,
		"h":        tcell.ColorTeal,
		"py":       tcell.ColorTeal,
		"js":       tcell.ColorTeal,
		"html":     tcell.ColorTeal,
		"css":      tcell.ColorTeal,
		"json":     tcell.ColorTeal,
		"xml":      tcell.ColorTeal,
		"yaml":     tcell.ColorTeal,
		"yml":      tcell.ColorTeal,
		"toml":     tcell.ColorTeal,
		"ini":      tcell.ColorTeal,
		"conf":     tcell.ColorTeal,
		"sh":       tcell.ColorGreen,
		"bat":      tcell.ColorGreen,
		"cmd":      tcell.ColorGreen,
		"ps1":      tcell.ColorGreen,
		"jpg":      tcell.ColorRed,
		"jpeg":     tcell.ColorRed,
		"png":      tcell.ColorRed,
		"gif":      tcell.ColorRed,
		"bmp":      tcell.ColorRed,
		"svg":      tcell.ColorRed,
		"mp3":      tcell.ColorPurple,
		"wav":      tcell.ColorPurple,
		"ogg":      tcell.ColorPurple,
		"mp4":      tcell.ColorPurple,
		"avi":      tcell.ColorPurple,
		"mkv":      tcell.ColorPurple,
		"pdf":      tcell.ColorRed,
		"doc":      tcell.ColorTeal,
		"docx":     tcell.ColorTeal,
		"xls":      tcell.ColorGreen,
		"xlsx":     tcell.ColorGreen,
		"ppt":      tcell.ColorYellow,
		"pptx":     tcell.ColorYellow,
		"hidden":   tcell.ColorGray,
		"default":  tcell.ColorBlack,
	}
}

// GetDarkColorScheme retorna um esquema de cores escuro para diferentes tipos de arquivos
func GetDarkColorScheme() map[string]tcell.Color {
	return map[string]tcell.Color{
		"dir":      tcell.ColorLightBlue,
		"exe":      tcell.ColorLightGreen,
		"zip":      tcell.ColorPurple,
		"tar":      tcell.ColorPurple,
		"gz":       tcell.ColorPurple,
		"rar":      tcell.ColorPurple,
		"7z":       tcell.ColorPurple,
		"txt":      tcell.ColorWhite,
		"md":       tcell.ColorWhite,
		"go":       tcell.ColorTeal,
		"c":        tcell.ColorTeal,
		"cpp":      tcell.ColorTeal,
		"h":        tcell.ColorTeal,
		"py":       tcell.ColorTeal,
		"js":       tcell.ColorTeal,
		"html":     tcell.ColorTeal,
		"css":      tcell.ColorTeal,
		"json":     tcell.ColorTeal,
		"xml":      tcell.ColorTeal,
		"yaml":     tcell.ColorTeal,
		"yml":      tcell.ColorTeal,
		"toml":     tcell.ColorTeal,
		"ini":      tcell.ColorTeal,
		"conf":     tcell.ColorTeal,
		"sh":       tcell.ColorLightGreen,
		"bat":      tcell.ColorLightGreen,
		"cmd":      tcell.ColorLightGreen,
		"ps1":      tcell.ColorLightGreen,
		"jpg":      tcell.ColorRed,
		"jpeg":     tcell.ColorRed,
		"png":      tcell.ColorRed,
		"gif":      tcell.ColorRed,
		"bmp":      tcell.ColorRed,
		"svg":      tcell.ColorRed,
		"mp3":      tcell.ColorPurple,
		"wav":      tcell.ColorPurple,
		"ogg":      tcell.ColorPurple,
		"mp4":      tcell.ColorPurple,
		"avi":      tcell.ColorPurple,
		"mkv":      tcell.ColorPurple,
		"pdf":      tcell.ColorRed,
		"doc":      tcell.ColorTeal,
		"docx":     tcell.ColorTeal,
		"xls":      tcell.ColorLightGreen,
		"xlsx":     tcell.ColorLightGreen,
		"ppt":      tcell.ColorYellow,
		"pptx":     tcell.ColorYellow,
		"hidden":   tcell.ColorGray,
		"default":  tcell.ColorWhite,
	}
}

// GetLightColorScheme retorna um esquema de cores claro para diferentes tipos de arquivos
func GetLightColorScheme() map[string]tcell.Color {
	return map[string]tcell.Color{
		"dir":      tcell.ColorBlue,
		"exe":      tcell.ColorDarkGreen,
		"zip":      tcell.ColorPurple,
		"tar":      tcell.ColorPurple,
		"gz":       tcell.ColorPurple,
		"rar":      tcell.ColorPurple,
		"7z":       tcell.ColorPurple,
		"txt":      tcell.ColorBlack,
		"md":       tcell.ColorBlack,
		"go":       tcell.ColorTeal,
		"c":        tcell.ColorTeal,
		"cpp":      tcell.ColorTeal,
		"h":        tcell.ColorTeal,
		"py":       tcell.ColorTeal,
		"js":       tcell.ColorTeal,
		"html":     tcell.ColorTeal,
		"css":      tcell.ColorTeal,
		"json":     tcell.ColorTeal,
		"xml":      tcell.ColorTeal,
		"yaml":     tcell.ColorTeal,
		"yml":      tcell.ColorTeal,
		"toml":     tcell.ColorTeal,
		"ini":      tcell.ColorTeal,
		"conf":     tcell.ColorTeal,
		"sh":       tcell.ColorDarkGreen,
		"bat":      tcell.ColorDarkGreen,
		"cmd":      tcell.ColorDarkGreen,
		"ps1":      tcell.ColorDarkGreen,
		"jpg":      tcell.ColorRed,
		"jpeg":     tcell.ColorRed,
		"png":      tcell.ColorRed,
		"gif":      tcell.ColorRed,
		"bmp":      tcell.ColorRed,
		"svg":      tcell.ColorRed,
		"mp3":      tcell.ColorPurple,
		"wav":      tcell.ColorPurple,
		"ogg":      tcell.ColorPurple,
		"mp4":      tcell.ColorPurple,
		"avi":      tcell.ColorPurple,
		"mkv":      tcell.ColorPurple,
		"pdf":      tcell.ColorRed,
		"doc":      tcell.ColorTeal,
		"docx":     tcell.ColorTeal,
		"xls":      tcell.ColorDarkGreen,
		"xlsx":     tcell.ColorDarkGreen,
		"ppt":      tcell.ColorYellow,
		"pptx":     tcell.ColorYellow,
		"hidden":   tcell.ColorGray,
		"default":  tcell.ColorBlack,
	}
}

// GetFileColor retorna a cor para um determinado arquivo com base em sua extensão
func GetFileColor(filename string, isDir bool, isHidden bool) tcell.Color {
	if isHidden {
		return GetRetroColorScheme()["hidden"]
	}
	
	if isDir {
		return GetRetroColorScheme()["dir"]
	}
	
	// Obter extensão
	ext := ""
	for i := len(filename) - 1; i >= 0; i-- {
		if filename[i] == '.' {
			ext = filename[i+1:]
			break
		}
	}
	
	// Verificar se a extensão existe no esquema de cores
	if color, ok := GetRetroColorScheme()[ext]; ok {
		return color
	}
	
	return GetRetroColorScheme()["default"]
}

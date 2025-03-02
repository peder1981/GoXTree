package utils

import (
	"github.com/gdamore/tcell/v2"
)

// Constantes para teclas de função
const (
	KeyHelp        = tcell.KeyF1
	KeyMenu        = tcell.KeyF2
	KeyView        = tcell.KeyF3
	KeyEdit        = tcell.KeyF4
	KeyCopy        = tcell.KeyF5
	KeyMove        = tcell.KeyF6
	KeyMkdir       = tcell.KeyF7
	KeyDelete      = tcell.KeyF8
	KeyCompress    = tcell.KeyF9
	KeyQuit        = tcell.KeyF10
	KeyTogglePanel = tcell.KeyTab
)

// Constantes para teclas de atalho (runa)
const (
	KeyMarkFile     = ' '
	KeyMarkAll      = 'a'
	KeySearch       = 'f'
	KeyGoto         = 'g'
	KeyToggleHidden = 'h'
	KeyRefresh      = 'r'
)

// Constantes para tamanhos e limites
const (
	MaxHistoryItems   = 50
	MaxFileViewSize   = 10 * 1024 * 1024  // 10MB
	MaxCompressSize   = 100 * 1024 * 1024 // 100MB
	DefaultDateFormat = "02/01/2006 15:04:05"
)

// Constantes para ordenação
const (
	SortByName = iota
	SortBySize
	SortByDate
	SortByType
)

// Esquema de cores
var (
	// Cores principais
	ColorBackground = tcell.ColorBlack
	ColorText       = tcell.ColorWhite
	ColorTitle      = tcell.ColorWhite
	ColorBorder     = tcell.ColorBlue
	ColorSelected   = tcell.ColorNavy
	ColorHighlight  = tcell.ColorYellow

	// Cores específicas
	ColorDirectoryText   = tcell.ColorYellow
	ColorDirectoryBorder = tcell.ColorBlue
	ColorFileText        = tcell.ColorWhite
	ColorHiddenText      = tcell.ColorGray
	ColorStatusText      = tcell.ColorWhite
	ColorStatusBg        = tcell.ColorBlue
	ColorStatusBar       = tcell.ColorBlue
	ColorHelpText        = tcell.ColorWhite
	ColorHelpBg          = tcell.ColorBlue
	ColorErrorText       = tcell.ColorRed
	ColorMenuText        = tcell.ColorWhite
	ColorMenuBg          = tcell.ColorBlue
	ColorMenuSelected    = tcell.ColorNavy

	// Cores para visualizadores e editores
	ColorViewerBorder = tcell.ColorBlue
	ColorViewerText   = tcell.ColorWhite
	ColorViewerTitle  = tcell.ColorYellow
	ColorEditorBorder = tcell.ColorBlue
	ColorEditorText   = tcell.ColorWhite
	ColorEditorTitle  = tcell.ColorYellow
	ColorHexByte      = tcell.ColorGreen
	ColorHexOffset    = tcell.ColorYellow
	ColorHexAscii     = tcell.ColorLightBlue
)

// Mensagens
const (
	MsgConfirmDelete    = "Tem certeza que deseja excluir os itens selecionados?"
	MsgConfirmQuit      = "Tem certeza que deseja sair?"
	MsgNoSelection      = "Nenhum arquivo selecionado"
	MsgCopySuccess      = "Arquivos copiados com sucesso"
	MsgMoveSuccess      = "Arquivos movidos com sucesso"
	MsgDeleteSuccess    = "Arquivos excluídos com sucesso"
	MsgCreateDirSuccess = "Diretório criado com sucesso"
	MsgFileViewerTitle  = "Visualizador de Arquivos"
	MsgHelpTitle        = "Ajuda do GoXTree"
)

// Versão do aplicativo
const (
	AppName    = "GoXTree"
	AppVersion = "1.0.0"
	AppAuthor  = "Peder Munksgaard"
	AppYear    = "2025"
	AppLicense = "MIT License"
)

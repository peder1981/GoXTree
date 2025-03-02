package ui

import (
	"fmt"
	"path/filepath"
	"runtime"
	"time"

	"github.com/peder1981/GoXTree/pkg/utils"
	"github.com/rivo/tview"
)

// StatusBar representa a barra de status
type StatusBar struct {
	statusBar *tview.TextView
}

// NewStatusBar cria uma nova barra de status
func NewStatusBar() *StatusBar {
	// Criar texto
	statusBar := tview.NewTextView().
		SetDynamicColors(true).
		SetTextAlign(tview.AlignLeft).
		SetWrap(false)

	// Configurar borda
	statusBar.SetBorder(true).
		SetTitle(" Status ").
		SetTitleAlign(tview.AlignLeft)

	// Criar StatusBar
	s := &StatusBar{
		statusBar: statusBar,
	}

	return s
}

// Update atualiza a barra de status
func (s *StatusBar) Update(currentDir string, numFiles, numDirs int, dirSize int64) {
	// Obter informações do sistema
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	// Formatar texto das informações
	infoText := fmt.Sprintf(
		"[yellow]Diretório:[white] %s | [yellow]Arquivos:[white] %d | [yellow]Diretórios:[white] %d | [yellow]Tamanho:[white] %s | [yellow]Memória:[white] %s | [yellow]Hora:[white] %s",
		filepath.Base(currentDir),
		numFiles,
		numDirs,
		formatSize(dirSize),
		formatSize(int64(memStats.Alloc)),
		time.Now().Format("15:04:05"),
	)

	// Formatar texto das opções de teclas
	optionsText := "[blue]F1[white]-Ajuda [blue]F2[white]-Renomear [blue]F7[white]-Criar Dir [blue]F8[white]-Excluir [blue]F10[white]-Sair [blue]Tab[white]-Alternar Painel [blue]ESC[white]-Voltar"

	// Atualizar texto
	s.statusBar.Clear()
	fmt.Fprintf(s.statusBar, "%s\n%s", infoText, optionsText)
}

// UpdateStatus atualiza as informações da barra de status
func (s *StatusBar) UpdateStatus(currentDir string) {
	// Obter informações do diretório
	numFiles, numDirs, dirSize := utils.GetDirectoryStats(currentDir)
	
	// Definir um texto informativo
	infoText := fmt.Sprintf("[yellow]Diretório:[white] %s | [yellow]Arquivos:[white] %d | [yellow]Diretórios:[white] %d | [yellow]Tamanho:[white] %s",
		currentDir, numFiles, numDirs, formatSize(dirSize))
	
	// Atualizar texto
	s.statusBar.Clear()
	fmt.Fprintf(s.statusBar, "%s", infoText)
}

// SetStatus define uma mensagem de status personalizada
func (s *StatusBar) SetStatus(message string) {
	s.statusBar.SetText(message)
}

// SetText define o texto da barra de status
func (s *StatusBar) SetText(message string) {
	s.statusBar.SetText(message)
}

// formatSize formata o tamanho em bytes
func formatSize(size int64) string {
	const (
		KB = 1024
		MB = 1024 * KB
		GB = 1024 * MB
	)

	switch {
	case size < KB:
		return fmt.Sprintf("%d B", size)
	case size < MB:
		return fmt.Sprintf("%.1f KB", float64(size)/KB)
	case size < GB:
		return fmt.Sprintf("%.1f MB", float64(size)/MB)
	default:
		return fmt.Sprintf("%.1f GB", float64(size)/GB)
	}
}

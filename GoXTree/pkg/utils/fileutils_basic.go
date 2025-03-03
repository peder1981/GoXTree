package utils

import (
	"crypto/md5"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// FileInfo representa informações de um arquivo
type FileInfo struct {
	Name      string
	Path      string
	Size      int64
	ModTime   time.Time
	IsDir     bool
	IsHidden  bool
	Extension string
}

// ListFiles lista arquivos em um diretório
func ListFiles(dirPath string, showHidden bool) ([]FileInfo, error) {
	var files []FileInfo

	// Ler diretório
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return nil, err
	}

	// Processar entradas
	for _, entry := range entries {
		// Verificar se é arquivo oculto
		name := entry.Name()
		isHidden := strings.HasPrefix(name, ".")

		// Pular arquivos ocultos se não estiver mostrando
		if isHidden && !showHidden {
			continue
		}

		// Obter informações do arquivo
		info, err := entry.Info()
		if err != nil {
			continue
		}

		// Criar FileInfo
		fileInfo := FileInfo{
			Name:      name,
			Path:      filepath.Join(dirPath, name),
			Size:      info.Size(),
			ModTime:   info.ModTime(),
			IsDir:     entry.IsDir(),
			IsHidden:  isHidden,
			Extension: strings.ToLower(filepath.Ext(name)),
		}

		files = append(files, fileInfo)
	}

	return files, nil
}

// GetDirectoryTree obtém a árvore de diretórios
func GetDirectoryTree(rootDir string, maxDepth int) ([]FileInfo, error) {
	var result []FileInfo

	err := filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Calcular a profundidade relativa ao diretório raiz
		relPath, err := filepath.Rel(rootDir, path)
		if err != nil {
			return err
		}

		// Não incluir o diretório raiz
		if relPath == "." {
			return nil
		}

		// Verificar a profundidade
		depth := len(strings.Split(relPath, string(os.PathSeparator)))
		if maxDepth > 0 && depth > maxDepth {
			if info.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		name := filepath.Base(path)
		isHidden := strings.HasPrefix(name, ".")

		extension := ""
		if !info.IsDir() {
			extension = strings.TrimPrefix(filepath.Ext(name), ".")
		}

		result = append(result, FileInfo{
			Name:      name,
			Path:      path,
			Size:      info.Size(),
			ModTime:   info.ModTime(),
			IsDir:     info.IsDir(),
			IsHidden:  isHidden,
			Extension: extension,
		})

		return nil
	})

	return result, err
}

// CountItems conta arquivos e diretórios em um diretório
func CountItems(dirPath string) (int, int, error) {
	files, err := ListFiles(dirPath, true)
	if err != nil {
		return 0, 0, err
	}

	fileCount := 0
	dirCount := 0

	for _, file := range files {
		if file.IsDir {
			dirCount++
		} else {
			fileCount++
		}
	}

	return fileCount, dirCount, nil
}

// CopyFile copia um arquivo de origem para destino
func CopyFile(src, dst string) error {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return fmt.Errorf("%s não é um arquivo regular", src)
	}

	source, err := os.Open(src)
	if err != nil {
		return err
	}
	defer source.Close()

	// Criar diretório de destino se não existir
	dstDir := filepath.Dir(dst)
	if err := os.MkdirAll(dstDir, os.ModePerm); err != nil {
		return err
	}

	destination, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destination.Close()

	_, err = io.Copy(destination, source)
	if err != nil {
		return err
	}

	// Preservar permissões e data de modificação
	return os.Chmod(dst, sourceFileStat.Mode())
}

// CopyDirectory copia um diretório recursivamente
func CopyDirectory(src, dst string) error {
	srcInfo, err := os.Stat(src)
	if err != nil {
		return err
	}

	if err := os.MkdirAll(dst, srcInfo.Mode()); err != nil {
		return err
	}

	entries, err := os.ReadDir(src)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())

		info, err := entry.Info()
		if err != nil {
			continue
		}

		if info.IsDir() {
			if err := CopyDirectory(srcPath, dstPath); err != nil {
				return err
			}
		} else {
			if err := CopyFile(srcPath, dstPath); err != nil {
				return err
			}
		}
	}

	return nil
}

// MoveFile move um arquivo de origem para destino
func MoveFile(src, dst string) error {
	// Tentar renomear primeiro (mais eficiente, mas só funciona no mesmo sistema de arquivos)
	err := os.Rename(src, dst)
	if err == nil {
		return nil
	}

	// Se falhar, copiar e depois excluir
	if err := CopyFile(src, dst); err != nil {
		return err
	}
	return os.Remove(src)
}

// DeleteFile exclui um arquivo ou diretório
func DeleteFile(path string) error {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return err
	}

	if fileInfo.IsDir() {
		return os.RemoveAll(path)
	}
	return os.Remove(path)
}

// CreateDirectory cria um diretório
func CreateDirectory(path string) error {
	return os.MkdirAll(path, os.ModePerm)
}

// GetFileType retorna o tipo de arquivo com base na extensão
func GetFileType(extension string) string {
	extension = strings.ToLower(extension)

	textExtensions := []string{"txt", "md", "log", "json", "xml", "html", "css", "js", "go", "py", "java", "c", "cpp", "h", "sh", "bat", "ini", "cfg", "conf", "yaml", "yml"}
	imageExtensions := []string{"jpg", "jpeg", "png", "gif", "bmp", "tiff", "webp", "svg", "ico"}
	audioExtensions := []string{"mp3", "wav", "ogg", "flac", "aac", "wma", "m4a"}
	videoExtensions := []string{"mp4", "avi", "mkv", "mov", "wmv", "flv", "webm", "m4v"}
	archiveExtensions := []string{"zip", "rar", "tar", "gz", "7z", "bz2", "xz"}
	documentExtensions := []string{"pdf", "doc", "docx", "xls", "xlsx", "ppt", "pptx", "odt", "ods", "odp"}

	for _, ext := range textExtensions {
		if extension == ext {
			return "text"
		}
	}
	for _, ext := range imageExtensions {
		if extension == ext {
			return "image"
		}
	}
	for _, ext := range audioExtensions {
		if extension == ext {
			return "audio"
		}
	}
	for _, ext := range videoExtensions {
		if extension == ext {
			return "video"
		}
	}
	for _, ext := range archiveExtensions {
		if extension == ext {
			return "archive"
		}
	}
	for _, ext := range documentExtensions {
		if extension == ext {
			return "document"
		}
	}

	return "unknown"
}

// FormatFileSize formata o tamanho do arquivo
func FormatFileSize(size int64) string {
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

// FormatSize formata o tamanho do arquivo
func FormatSize(size int64) string {
	return FormatFileSize(size)
}

// GetFileIcon retorna um ícone para o tipo de arquivo
func GetFileIcon(fileInfo FileInfo) string {
	if fileInfo.IsDir {
		if fileInfo.IsHidden {
			return "📁"
		}
		return "📂"
	}

	fileType := GetFileType(fileInfo.Extension)
	switch fileType {
	case "text":
		return "📄"
	case "image":
		return "🖼️"
	case "audio":
		return "🎵"
	case "video":
		return "🎬"
	case "archive":
		return "📦"
	case "document":
		return "📑"
	default:
		if fileInfo.IsHidden {
			return "🔒"
		}
		return "📄"
	}
}

// GetFileContent obtém o conteúdo de um arquivo
func GetFileContent(path string, maxSize int64) ([]byte, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	if maxSize > 0 {
		stat, err := file.Stat()
		if err != nil {
			return nil, err
		}
		if stat.Size() > maxSize {
			return nil, fmt.Errorf("arquivo muito grande (tamanho: %d, máximo: %d)", stat.Size(), maxSize)
		}
	}

	return io.ReadAll(file)
}

// IsTextFile verifica se um arquivo é de texto
func IsTextFile(filePath string) (bool, error) {
	// Abrir arquivo
	file, err := os.Open(filePath)
	if err != nil {
		return false, err
	}
	defer file.Close()

	// Ler os primeiros 512 bytes
	buffer := make([]byte, 512)
	n, err := file.Read(buffer)
	if err != nil && err != io.EOF {
		return false, err
	}

	// Verificar se há caracteres nulos, que indicam arquivo binário
	for i := 0; i < n; i++ {
		if buffer[i] == 0 {
			return false, nil
		}
	}

	// Verificar se o arquivo contém apenas caracteres ASCII
	for i := 0; i < n; i++ {
		if buffer[i] > 127 {
			// Arquivo pode ser texto, mas não ASCII
			// Vamos considerar como texto mesmo assim
			return true, nil
		}
	}

	return true, nil
}

// GetParentDirectory retorna o diretório pai
func GetParentDirectory(path string) string {
	return filepath.Dir(path)
}

// GetDirectorySize calcula o tamanho total de um diretório
func GetDirectorySize(path string) (int64, error) {
	var size int64

	err := filepath.Walk(path, func(_ string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			size += info.Size()
		}
		return nil
	})

	return size, err
}

// GetDirectoryStats retorna estatísticas do diretório (número de arquivos, número de diretórios e tamanho total)
func GetDirectoryStats(dirPath string) (int, int, int64) {
	var fileCount, dirCount int
	var totalSize int64

	filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil // Ignorar erros
		}

		// Não contar o diretório raiz
		if path == dirPath {
			return nil
		}

		if info.IsDir() {
			dirCount++
		} else {
			fileCount++
			totalSize += info.Size()
		}
		return nil
	})

	return fileCount, dirCount, totalSize
}

// CalculateFileMD5 calcula o hash MD5 de um arquivo
func CalculateFileMD5(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", hash.Sum(nil)), nil
}

// GetFileMIMEType retorna o tipo MIME de um arquivo
func GetFileMIMEType(filePath string) string {
	file, err := os.Open(filePath)
	if err != nil {
		return "application/octet-stream"
	}
	defer file.Close()

	// Ler os primeiros 512 bytes para detectar o tipo MIME
	buffer := make([]byte, 512)
	n, err := file.Read(buffer)
	if err != nil && err != io.EOF {
		return "application/octet-stream"
	}
	buffer = buffer[:n]

	// Detectar o tipo MIME
	mimeType := http.DetectContentType(buffer)

	return mimeType
}

// isHidden verifica se um arquivo/diretório é oculto
func isHidden(path string) bool {
	name := filepath.Base(path)
	return strings.HasPrefix(name, ".")
}

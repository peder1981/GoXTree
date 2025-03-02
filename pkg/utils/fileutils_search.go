package utils

import (
	"bufio"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

// SearchOptions define as opções para busca avançada
type SearchOptions struct {
	Pattern        string
	Directory      string
	Recursive      bool
	CaseSensitive  bool
	MatchContent   bool
	MinSize        int64
	MaxSize        int64
	ModifiedAfter  time.Time
	ModifiedBefore time.Time
	FileTypes      []string
}

// SearchResult representa um resultado de busca
type SearchResult struct {
	Path      string
	Name      string
	Size      int64
	ModTime   time.Time
	IsDir     bool
	MatchLine string
	LineNum   int
}

// AdvancedSearchFiles realiza uma busca avançada de arquivos
func AdvancedSearchFiles(options SearchOptions) ([]SearchResult, error) {
	var results []SearchResult

	// Compilar a expressão regular para o padrão de busca
	var pattern *regexp.Regexp
	var err error

	if options.CaseSensitive {
		pattern, err = regexp.Compile(options.Pattern)
	} else {
		pattern, err = regexp.Compile("(?i)" + options.Pattern)
	}

	if err != nil {
		return nil, err
	}

	// Função para verificar se um arquivo corresponde aos critérios
	matchesFileType := func(ext string) bool {
		if len(options.FileTypes) == 0 {
			return true
		}
		ext = strings.ToLower(strings.TrimPrefix(ext, "."))
		for _, t := range options.FileTypes {
			if t == ext {
				return true
			}
		}
		return false
	}

	// Função para verificar se um arquivo corresponde aos critérios de tamanho e data
	matchesCriteria := func(info os.FileInfo) bool {
		// Verificar tamanho
		if options.MinSize > 0 && info.Size() < options.MinSize {
			return false
		}
		if options.MaxSize > 0 && info.Size() > options.MaxSize {
			return false
		}

		// Verificar data de modificação
		if !options.ModifiedAfter.IsZero() && info.ModTime().Before(options.ModifiedAfter) {
			return false
		}
		if !options.ModifiedBefore.IsZero() && info.ModTime().After(options.ModifiedBefore) {
			return false
		}

		return true
	}

	// Função para processar um arquivo
	processFile := func(path string, info os.FileInfo) error {
		// Verificar se o arquivo corresponde aos critérios
		if !matchesCriteria(info) {
			return nil
		}

		// Verificar se o tipo de arquivo corresponde
		if !info.IsDir() && !matchesFileType(filepath.Ext(path)) {
			return nil
		}

		// Verificar se o nome corresponde ao padrão
		name := filepath.Base(path)
		if pattern.MatchString(name) {
			results = append(results, SearchResult{
				Path:    path,
				Name:    name,
				Size:    info.Size(),
				ModTime: info.ModTime(),
				IsDir:   info.IsDir(),
			})
			return nil
		}

		// Se não for para verificar o conteúdo ou for um diretório, retornar
		if !options.MatchContent || info.IsDir() {
			return nil
		}

		// Verificar se é um arquivo de texto antes de procurar no conteúdo
		isText, err := IsTextFile(path)
		if err != nil || !isText {
			return nil
		}

		// Procurar no conteúdo do arquivo
		file, err := os.Open(path)
		if err != nil {
			return nil
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		lineNum := 0
		for scanner.Scan() {
			lineNum++
			line := scanner.Text()
			if pattern.MatchString(line) {
				results = append(results, SearchResult{
					Path:      path,
					Name:      name,
					Size:      info.Size(),
					ModTime:   info.ModTime(),
					IsDir:     info.IsDir(),
					MatchLine: line,
					LineNum:   lineNum,
				})
				break // Encontrou uma correspondência, não precisa continuar procurando neste arquivo
			}
		}

		return scanner.Err()
	}

	// Percorrer o diretório
	if options.Recursive {
		err = filepath.Walk(options.Directory, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return nil // Ignorar erros e continuar
			}
			return processFile(path, info)
		})
	} else {
		// Apenas listar arquivos no diretório atual
		entries, err := os.ReadDir(options.Directory)
		if err != nil {
			return nil, err
		}

		for _, entry := range entries {
			info, err := entry.Info()
			if err != nil {
				continue
			}
			path := filepath.Join(options.Directory, entry.Name())
			if err := processFile(path, info); err != nil {
				return nil, err
			}
		}
	}

	return results, err
}

// FindDuplicateFiles encontra arquivos duplicados em um diretório
func FindDuplicateFiles(rootDir string, recursive bool) (map[int64][]FileInfo, error) {
	// Mapa de tamanho -> lista de arquivos
	sizeMap := make(map[int64][]FileInfo)

	// Função para processar um diretório
	var processDir func(dirPath string) error
	processDir = func(dirPath string) error {
		entries, err := os.ReadDir(dirPath)
		if err != nil {
			return err
		}

		for _, entry := range entries {
			info, err := entry.Info()
			if err != nil {
				continue
			}

			path := filepath.Join(dirPath, entry.Name())

			if info.IsDir() {
				if recursive {
					if err := processDir(path); err != nil {
						return err
					}
				}
				continue
			}

			// Ignorar arquivos vazios
			if info.Size() == 0 {
				continue
			}

			name := entry.Name()
			isHidden := strings.HasPrefix(name, ".")
			extension := strings.TrimPrefix(filepath.Ext(name), ".")

			fileInfo := FileInfo{
				Name:      name,
				Path:      path,
				Size:      info.Size(),
				ModTime:   info.ModTime(),
				IsDir:     false,
				IsHidden:  isHidden,
				Extension: extension,
			}

			sizeMap[info.Size()] = append(sizeMap[info.Size()], fileInfo)
		}

		return nil
	}

	// Processar o diretório raiz
	if err := processDir(rootDir); err != nil {
		return nil, err
	}

	// Filtrar para manter apenas tamanhos com mais de um arquivo
	result := make(map[int64][]FileInfo)
	for size, files := range sizeMap {
		if len(files) > 1 {
			result[size] = files
		}
	}

	return result, nil
}

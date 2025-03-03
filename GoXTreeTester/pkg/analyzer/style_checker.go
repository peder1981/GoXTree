package analyzer

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/peder1981/GoXTreeTester/pkg/reporter"
)

// StyleChecker verifica a consistência do estilo de código
type StyleChecker struct {
	projectPath    string
	reporter       *reporter.Reporter
	maxLineLength  int
}

// NewStyleChecker cria um novo verificador de estilo
func NewStyleChecker(projectPath string, reporter *reporter.Reporter) *StyleChecker {
	return &StyleChecker{
		projectPath:   projectPath,
		reporter:      reporter,
		maxLineLength: 100, // Usa o valor global definido em analyzer.go
	}
}

// SetMaxLineLength define o comprimento máximo de linha para este verificador de estilo
func (sc *StyleChecker) SetMaxLineLength(length int) {
	sc.maxLineLength = length
}

// StyleIssue representa um problema de estilo de código
type StyleIssue struct {
	File     string
	Line     int
	Message  string
	Severity string
}

// CheckStyle verifica o estilo de código em todo o projeto
func (sc *StyleChecker) CheckStyle() ([]StyleIssue, error) {
	var issues []StyleIssue

	// Encontrar todos os arquivos Go no projeto
	var goFiles []string
	err := filepath.Walk(sc.projectPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasSuffix(path, ".go") {
			goFiles = append(goFiles, path)
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("erro ao percorrer diretório do projeto: %w", err)
	}

	// Verificar cada arquivo
	for _, filePath := range goFiles {
		fileIssues, err := sc.checkFile(filePath)
		if err != nil {
			sc.reporter.AddWarning(fmt.Sprintf("Erro ao verificar estilo em %s: %v", filePath, err))
			continue
		}
		issues = append(issues, fileIssues...)
	}

	return issues, nil
}

// checkFile verifica o estilo de código em um arquivo
func (sc *StyleChecker) checkFile(filePath string) ([]StyleIssue, error) {
	var issues []StyleIssue

	// Ler o conteúdo do arquivo
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("erro ao ler arquivo %s: %w", filePath, err)
	}

	// Verificar linhas longas
	issues = append(issues, sc.checkLineLengths(filePath, string(content))...)

	// Verificar espaços em branco no final das linhas
	issues = append(issues, sc.checkTrailingWhitespace(filePath, string(content))...)

	// Verificar nomenclatura de funções e variáveis
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, filePath, content, parser.ParseComments)
	if err != nil {
		return issues, fmt.Errorf("erro ao analisar arquivo %s: %w", filePath, err)
	}

	issues = append(issues, sc.checkNaming(filePath, fset, file)...)
	issues = append(issues, sc.checkComments(filePath, fset, file)...)

	return issues, nil
}

// checkLineLengths verifica se há linhas muito longas
func (sc *StyleChecker) checkLineLengths(filePath, content string) []StyleIssue {
	var issues []StyleIssue
	lines := strings.Split(content, "\n")

	for i, line := range lines {
		if len(line) > sc.maxLineLength {
			issues = append(issues, StyleIssue{
				File:     filePath,
				Line:     i + 1,
				Message:  fmt.Sprintf("Linha muito longa (%d caracteres)", len(line)),
				Severity: "warning",
			})
		}
	}

	return issues
}

// checkTrailingWhitespace verifica se há espaços em branco no final das linhas
func (sc *StyleChecker) checkTrailingWhitespace(filePath, content string) []StyleIssue {
	var issues []StyleIssue
	lines := strings.Split(content, "\n")
	trailingSpaceRegex := regexp.MustCompile(`\s+$`)

	for i, line := range lines {
		if trailingSpaceRegex.MatchString(line) {
			issues = append(issues, StyleIssue{
				File:     filePath,
				Line:     i + 1,
				Message:  "Espaços em branco no final da linha",
				Severity: "warning",
			})
		}
	}

	return issues
}

// checkNaming verifica a nomenclatura de funções e variáveis
func (sc *StyleChecker) checkNaming(filePath string, fset *token.FileSet, file *ast.File) []StyleIssue {
	var issues []StyleIssue
	
	// Verificar nomes de funções
	ast.Inspect(file, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.FuncDecl:
			// Funções exportadas devem começar com letra maiúscula
			if ast.IsExported(x.Name.Name) {
				// Verificar se a função tem comentário
				if x.Doc == nil || len(x.Doc.List) == 0 {
					issues = append(issues, StyleIssue{
						File:     filePath,
						Line:     fset.Position(x.Pos()).Line,
						Message:  fmt.Sprintf("Função exportada %s não tem comentário", x.Name.Name),
						Severity: "warning",
					})
				}
			}
		case *ast.TypeSpec:
			// Tipos exportados devem começar com letra maiúscula
			if ast.IsExported(x.Name.Name) {
				// Verificar se o tipo tem comentário
				if x.Doc == nil || len(x.Doc.List) == 0 {
					issues = append(issues, StyleIssue{
						File:     filePath,
						Line:     fset.Position(x.Pos()).Line,
						Message:  fmt.Sprintf("Tipo exportado %s não tem comentário", x.Name.Name),
						Severity: "warning",
					})
				}
			}
		}
		return true
	})

	return issues
}

// checkComments verifica se os comentários seguem o padrão do Go
func (sc *StyleChecker) checkComments(filePath string, fset *token.FileSet, file *ast.File) []StyleIssue {
	var issues []StyleIssue

	for _, comment := range file.Comments {
		for _, c := range comment.List {
			text := c.Text
			if strings.HasPrefix(text, "//") {
				// Verificar se há espaço após //
				if len(text) > 2 && text[2] != ' ' && !strings.HasPrefix(text, "///") {
					issues = append(issues, StyleIssue{
						File:     filePath,
						Line:     fset.Position(c.Pos()).Line,
						Message:  "Comentário deve ter espaço após //",
						Severity: "warning",
					})
				}
			}
		}
	}

	return issues
}

// FixStyleIssues corrige problemas de estilo automaticamente
func (sc *StyleChecker) FixStyleIssues(issues []StyleIssue) (int, error) {
	// Agrupar problemas por arquivo
	fileIssues := make(map[string][]StyleIssue)
	for _, issue := range issues {
		fileIssues[issue.File] = append(fileIssues[issue.File], issue)
	}

	fixedCount := 0

	// Corrigir cada arquivo
	for filePath, issues := range fileIssues {
		// Ler o conteúdo do arquivo
		content, err := os.ReadFile(filePath)
		if err != nil {
			return fixedCount, fmt.Errorf("erro ao ler arquivo %s: %w", filePath, err)
		}

		// Converter para string
		contentStr := string(content)
		lines := strings.Split(contentStr, "\n")

		// Corrigir problemas
		modified := false
		for _, issue := range issues {
			switch {
			case strings.Contains(issue.Message, "Espaços em branco no final da linha"):
				if issue.Line-1 < len(lines) {
					lines[issue.Line-1] = strings.TrimRight(lines[issue.Line-1], " \t")
					modified = true
					fixedCount++
				}
			case strings.Contains(issue.Message, "Comentário deve ter espaço após //"):
				if issue.Line-1 < len(lines) {
					line := lines[issue.Line-1]
					if strings.Contains(line, "//") && !strings.Contains(line, "// ") && !strings.Contains(line, "///") {
						lines[issue.Line-1] = strings.Replace(line, "//", "// ", 1)
						modified = true
						fixedCount++
					}
				}
			}
		}

		// Salvar o arquivo modificado
		if modified {
			newContent := strings.Join(lines, "\n")
			err = os.WriteFile(filePath, []byte(newContent), 0644)
			if err != nil {
				return fixedCount, fmt.Errorf("erro ao escrever arquivo %s: %w", filePath, err)
			}
			sc.reporter.AddFix(filePath, "Corrigidos problemas de estilo")
		}
	}

	return fixedCount, nil
}

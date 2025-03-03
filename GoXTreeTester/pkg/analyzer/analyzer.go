package analyzer

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"

	"github.com/peder1981/GoXTreeTester/pkg/reporter"
	"golang.org/x/tools/go/packages"
)

// Severity representa a gravidade de um problema
type Severity string

const (
	SeverityError   Severity = "ERROR"
	SeverityWarning Severity = "WARNING"
	SeverityInfo    Severity = "INFO"
)

// Issue representa um problema encontrado durante a análise
type Issue struct {
	File      string
	Line      int
	Column    int
	Message   string
	Severity  Severity
	Code      string
	Fixable   bool
	FixMethod string
}

// Analyzer é responsável por analisar o código-fonte
type Analyzer struct {
	projectPath  string
	ignoreErrors []string
	reporter     *reporter.Reporter
	fset         *token.FileSet
}

// NewAnalyzer cria um novo analisador
func NewAnalyzer(projectPath string, ignoreErrors []string, reporter *reporter.Reporter) *Analyzer {
	return &Analyzer{
		projectPath:  projectPath,
		ignoreErrors: ignoreErrors,
		reporter:     reporter,
		fset:         token.NewFileSet(),
	}
}

// Analyze analisa o código-fonte e retorna uma lista de problemas
func (a *Analyzer) Analyze() ([]Issue, error) {
	var issues []Issue

	// Configurar o carregamento de pacotes
	cfg := &packages.Config{
		Mode:  packages.NeedSyntax | packages.NeedTypes | packages.NeedTypesInfo | packages.NeedDeps,
		Tests: true,
		Dir:   a.projectPath,
		Fset:  a.fset,
	}

	// Carregar pacotes
	pkgs, err := packages.Load(cfg, "./...")
	if err != nil {
		return nil, fmt.Errorf("erro ao carregar pacotes: %w", err)
	}

	// Verificar erros de pacotes
	for _, pkg := range pkgs {
		for _, err := range pkg.Errors {
			issues = append(issues, Issue{
				File:     err.Pos,
				Line:     0, // Não podemos obter a linha exata sem um token.Pos
				Column:   0, // Não podemos obter a coluna exata sem um token.Pos
				Message:  err.Msg,
				Severity: SeverityError,
				Fixable:  false,
			})
		}
	}

	// Analisar cada arquivo Go no projeto
	err = filepath.Walk(a.projectPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Ignorar diretórios
		if info.IsDir() {
			// Ignorar diretórios ocultos e vendor
			if strings.HasPrefix(info.Name(), ".") || info.Name() == "vendor" {
				return filepath.SkipDir
			}
			return nil
		}

		// Processar apenas arquivos Go
		if !strings.HasSuffix(info.Name(), ".go") {
			return nil
		}

		// Analisar o arquivo
		fileIssues, err := a.analyzeFile(path)
		if err != nil {
			return err
		}

		issues = append(issues, fileIssues...)
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("erro ao percorrer diretórios: %w", err)
	}

	// Filtrar problemas ignorados
	if len(a.ignoreErrors) > 0 {
		var filteredIssues []Issue
		for _, issue := range issues {
			ignored := false
			for _, ignorePattern := range a.ignoreErrors {
				if strings.Contains(issue.Message, ignorePattern) {
					ignored = true
					break
				}
			}
			if !ignored {
				filteredIssues = append(filteredIssues, issue)
			}
		}
		issues = filteredIssues
	}

	// Adicionar problemas ao relatório
	for _, issue := range issues {
		a.reporter.AddIssue(issue.File, issue.Line, issue.Message, string(issue.Severity))
	}

	return issues, nil
}

// analyzeFile analisa um único arquivo Go
func (a *Analyzer) analyzeFile(filePath string) ([]Issue, error) {
	var issues []Issue

	// Analisar o arquivo
	file, err := parser.ParseFile(a.fset, filePath, nil, parser.AllErrors)
	if err != nil {
		return nil, fmt.Errorf("erro ao analisar arquivo %s: %w", filePath, err)
	}

	// Verificar problemas comuns
	issues = append(issues, a.checkDuplicateFunctions(file, filePath)...)
	issues = append(issues, a.checkUnusedImports(file, filePath)...)
	issues = append(issues, a.checkUndefinedMethods(file, filePath)...)
	issues = append(issues, a.checkNamingConsistency(file, filePath)...)

	return issues, nil
}

// checkDuplicateFunctions verifica funções duplicadas
func (a *Analyzer) checkDuplicateFunctions(file *ast.File, filePath string) []Issue {
	var issues []Issue
	functions := make(map[string]token.Pos)

	// Percorrer declarações
	for _, decl := range file.Decls {
		if funcDecl, ok := decl.(*ast.FuncDecl); ok {
			var funcName string
			if funcDecl.Recv != nil {
				// Método
				if len(funcDecl.Recv.List) > 0 {
					if starExpr, ok := funcDecl.Recv.List[0].Type.(*ast.StarExpr); ok {
						if ident, ok := starExpr.X.(*ast.Ident); ok {
							funcName = ident.Name + "." + funcDecl.Name.Name
						}
					} else if ident, ok := funcDecl.Recv.List[0].Type.(*ast.Ident); ok {
						funcName = ident.Name + "." + funcDecl.Name.Name
					}
				}
			} else {
				// Função
				funcName = funcDecl.Name.Name
			}

			if funcName != "" {
				if pos, exists := functions[funcName]; exists {
					issues = append(issues, Issue{
						File:      filePath,
						Line:      a.fset.Position(funcDecl.Pos()).Line,
						Column:    a.fset.Position(funcDecl.Pos()).Column,
						Message:   fmt.Sprintf("Função duplicada: %s (primeira declaração na linha %d)", funcName, a.fset.Position(pos).Line),
						Severity:  SeverityError,
						Fixable:   true,
						FixMethod: "RemoveDuplicateFunction",
					})
				} else {
					functions[funcName] = funcDecl.Pos()
				}
			}
		}
	}

	return issues
}

// checkUnusedImports verifica imports não utilizados
func (a *Analyzer) checkUnusedImports(file *ast.File, filePath string) []Issue {
	var issues []Issue

	// Coletar todos os identificadores usados no arquivo
	var usedIdents []*ast.Ident
	ast.Inspect(file, func(n ast.Node) bool {
		if ident, ok := n.(*ast.Ident); ok {
			usedIdents = append(usedIdents, ident)
		}
		return true
	})

	// Verificar cada import
	for _, imp := range file.Imports {
		if imp.Name != nil && imp.Name.Name == "_" {
			// Ignorar imports com underscore
			continue
		}

		importPath := strings.Trim(imp.Path.Value, "\"")
		importName := filepath.Base(importPath)
		if imp.Name != nil {
			importName = imp.Name.Name
		}

		used := false
		for _, ident := range usedIdents {
			if ident.Name == importName {
				used = true
				break
			}
		}

		if !used {
			issues = append(issues, Issue{
				File:      filePath,
				Line:      a.fset.Position(imp.Pos()).Line,
				Column:    a.fset.Position(imp.Pos()).Column,
				Message:   fmt.Sprintf("Import não utilizado: %s", importPath),
				Severity:  SeverityWarning,
				Fixable:   true,
				FixMethod: "RemoveUnusedImport",
			})
		}
	}

	return issues
}

// checkUndefinedMethods verifica chamadas para métodos não definidos
func (a *Analyzer) checkUndefinedMethods(file *ast.File, filePath string) []Issue {
	var issues []Issue

	// Coletar todos os métodos definidos
	definedMethods := make(map[string]bool)
	for _, decl := range file.Decls {
		if funcDecl, ok := decl.(*ast.FuncDecl); ok {
			if funcDecl.Recv != nil {
				// É um método
				if len(funcDecl.Recv.List) > 0 {
					if starExpr, ok := funcDecl.Recv.List[0].Type.(*ast.StarExpr); ok {
						if ident, ok := starExpr.X.(*ast.Ident); ok {
							methodName := ident.Name + "." + funcDecl.Name.Name
							definedMethods[methodName] = true
						}
					} else if ident, ok := funcDecl.Recv.List[0].Type.(*ast.Ident); ok {
						methodName := ident.Name + "." + funcDecl.Name.Name
						definedMethods[methodName] = true
					}
				}
			}
		}
	}

	// Lista de tipos padrão da biblioteca Go que não precisamos verificar
	standardTypes := map[string]bool{
		"string":         true,
		"[]byte":         true,
		"int":            true,
		"float64":        true,
		"bool":           true,
		"time.Time":      true,
		"strings.Builder": true,
		"bytes.Buffer":   true,
		"fmt.Stringer":   true,
		"io.Reader":      true,
		"io.Writer":      true,
		"os.FileInfo":    true,
		"os.File":        true,
		"image.Image":    true,
		"image.Rectangle": true,
		"tcell.Event":    true,
		"tcell.EventKey": true,
		"tview.Box":      true,
		"tview.TextView": true,
		"tview.Flex":     true,
		"tview.Grid":     true,
		"tview.Form":     true,
		"tview.Modal":    true,
		"tview.List":     true,
		"tview.Table":    true,
		"tview.TreeView": true,
		"tview.InputField": true,
		"tview.DropDown": true,
		"tview.Button":   true,
		"tview.Application": true,
		"tview.Pages":    true,
		"tview.Primitive": true,
	}

	// Verificar chamadas de métodos
	ast.Inspect(file, func(n ast.Node) bool {
		if callExpr, ok := n.(*ast.CallExpr); ok {
			if selExpr, ok := callExpr.Fun.(*ast.SelectorExpr); ok {
				if ident, ok := selExpr.X.(*ast.Ident); ok {
					methodName := ident.Name + "." + selExpr.Sel.Name
					
					// Verificar se o método é chamado em um objeto que não é um pacote
					if ident.Obj != nil && ident.Obj.Kind == ast.Var {
						// Ignorar tipos padrão da biblioteca Go
						if ident.Name == "hexOutput" || standardTypes[ident.Name] {
							return true
						}
						
						// Verificar se o nome do objeto termina com palavras comuns que indicam tipos externos
						if strings.HasSuffix(ident.Name, "View") || 
						   strings.HasSuffix(ident.Name, "Layout") ||
						   strings.HasSuffix(ident.Name, "Widget") ||
						   strings.HasSuffix(ident.Name, "Component") ||
						   strings.HasSuffix(ident.Name, "Info") ||
						   strings.HasSuffix(ident.Name, "Event") ||
						   strings.HasSuffix(ident.Name, "File") ||
						   strings.HasSuffix(ident.Name, "Reader") ||
						   strings.HasSuffix(ident.Name, "Writer") {
							return true
						}
						
						// Verificar se o método está definido
						if !definedMethods[methodName] {
							issues = append(issues, Issue{
								File:      filePath,
								Line:      a.fset.Position(callExpr.Pos()).Line,
								Column:    a.fset.Position(callExpr.Pos()).Column,
								Message:   fmt.Sprintf("Chamada para método não definido: %s", methodName),
								Severity:  SeverityError,
								Fixable:   false,
							})
						}
					}
				}
			}
		}
		return true
	})

	return issues
}

// checkNamingConsistency verifica consistência na nomenclatura
func (a *Analyzer) checkNamingConsistency(file *ast.File, filePath string) []Issue {
	var issues []Issue

	// Verificar nomes de funções e métodos
	for _, decl := range file.Decls {
		if funcDecl, ok := decl.(*ast.FuncDecl); ok {
			funcName := funcDecl.Name.Name
			
			// Verificar padrão camelCase
			if len(funcName) > 0 && funcName[0] >= 'A' && funcName[0] <= 'Z' {
				// Funções exportadas devem começar com maiúscula
				continue
			}
			
			// Verificar se o nome contém underscore (não é padrão Go)
			if strings.Contains(funcName, "_") {
				issues = append(issues, Issue{
					File:      filePath,
					Line:      a.fset.Position(funcDecl.Pos()).Line,
					Column:    a.fset.Position(funcDecl.Pos()).Column,
					Message:   fmt.Sprintf("Nome de função não segue o padrão camelCase: %s", funcName),
					Severity:  SeverityWarning,
					Fixable:   true,
					FixMethod: "FixNaming",
				})
			}
		}
	}

	return issues
}

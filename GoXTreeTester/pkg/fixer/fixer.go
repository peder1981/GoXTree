package fixer

import (
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"os"
	"strings"

	"github.com/peder1981/GoXTreeTester/pkg/analyzer"
	"github.com/peder1981/GoXTreeTester/pkg/reporter"
)

// Fixer é responsável por corrigir problemas no código
type Fixer struct {
	projectPath string
	reporter    *reporter.Reporter
}

// NewFixer cria um novo corretor
func NewFixer(projectPath string, reporter *reporter.Reporter) *Fixer {
	return &Fixer{
		projectPath: projectPath,
		reporter:    reporter,
	}
}

// FixIssues corrige problemas identificados
func (f *Fixer) FixIssues(issues []analyzer.Issue) (int, error) {
	fixedCount := 0

	// Agrupar problemas por arquivo
	fileIssues := make(map[string][]analyzer.Issue)
	for _, issue := range issues {
		fileIssues[issue.File] = append(fileIssues[issue.File], issue)
	}

	// Corrigir cada arquivo
	for filePath, issues := range fileIssues {
		fixed, err := f.fixFile(filePath, issues)
		if err != nil {
			return fixedCount, fmt.Errorf("erro ao corrigir arquivo %s: %w", filePath, err)
		}
		fixedCount += fixed
	}

	return fixedCount, nil
}

// fixFile corrige problemas em um único arquivo
func (f *Fixer) fixFile(filePath string, issues []analyzer.Issue) (int, error) {
	fixedCount := 0

	// Verificar se o arquivo existe
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return 0, fmt.Errorf("arquivo não encontrado: %s", filePath)
	}

	// Analisar o arquivo
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, filePath, nil, parser.ParseComments)
	if err != nil {
		return 0, fmt.Errorf("erro ao analisar arquivo: %w", err)
	}

	// Aplicar correções
	modified := false

	// Remover funções duplicadas
	if fixed, m := f.removeDuplicateFunctions(file, fset, issues); m {
		fixedCount += fixed
		modified = true
	}

	// Remover imports não utilizados
	if fixed, m := f.removeUnusedImports(file, issues); m {
		fixedCount += fixed
		modified = true
	}

	// Corrigir nomes de funções
	if fixed, m := f.fixNaming(file, fset, issues); m {
		fixedCount += fixed
		modified = true
	}

	// Corrigir chamadas para métodos não definidos
	if fixed, m := f.fixUndefinedMethods(file, fset, issues); m {
		fixedCount += fixed
		modified = true
	}

	// Salvar o arquivo modificado
	if modified {
		// Criar backup do arquivo original
		backupPath := filePath + ".bak"
		if err := copyFile(filePath, backupPath); err != nil {
			return fixedCount, fmt.Errorf("erro ao criar backup: %w", err)
		}

		// Formatar e salvar o arquivo modificado
		var buf strings.Builder
		if err := format.Node(&buf, fset, file); err != nil {
			return fixedCount, fmt.Errorf("erro ao formatar arquivo: %w", err)
		}

		if err := os.WriteFile(filePath, []byte(buf.String()), 0644); err != nil {
			return fixedCount, fmt.Errorf("erro ao salvar arquivo: %w", err)
		}

		// Adicionar informação ao relatório
		f.reporter.AddFix(filePath, fmt.Sprintf("Corrigidos %d problemas", fixedCount))
	}

	return fixedCount, nil
}

// removeDuplicateFunctions remove funções duplicadas
func (f *Fixer) removeDuplicateFunctions(file *ast.File, fset *token.FileSet, issues []analyzer.Issue) (int, bool) {
	fixedCount := 0
	modified := false

	// Coletar funções duplicadas
	duplicateFuncs := make(map[int]bool)
	for _, issue := range issues {
		if issue.FixMethod == "RemoveDuplicateFunction" {
			duplicateFuncs[issue.Line] = true
		}
	}

	if len(duplicateFuncs) == 0 {
		return 0, false
	}

	// Filtrar declarações, removendo funções duplicadas
	var newDecls []ast.Decl
	for _, decl := range file.Decls {
		if funcDecl, ok := decl.(*ast.FuncDecl); ok {
			line := fset.Position(funcDecl.Pos()).Line
			if duplicateFuncs[line] {
				// Pular esta função (removê-la)
				fixedCount++
				modified = true
				continue
			}
		}
		newDecls = append(newDecls, decl)
	}

	file.Decls = newDecls
	return fixedCount, modified
}

// removeUnusedImports remove imports não utilizados
func (f *Fixer) removeUnusedImports(file *ast.File, issues []analyzer.Issue) (int, bool) {
	fixedCount := 0
	modified := false

	// Coletar imports não utilizados
	unusedImports := make(map[string]bool)
	for _, issue := range issues {
		if issue.FixMethod == "RemoveUnusedImport" {
			// Extrair o nome do import da mensagem
			parts := strings.Split(issue.Message, ": ")
			if len(parts) > 1 {
				importPath := parts[1]
				unusedImports[importPath] = true
			}
		}
	}

	if len(unusedImports) == 0 {
		return 0, false
	}

	// Filtrar imports
	var specs []ast.Spec
	for _, spec := range file.Imports {
		importPath := strings.Trim(spec.Path.Value, "\"")
		if unusedImports[importPath] {
			// Pular este import (removê-lo)
			fixedCount++
			modified = true
			continue
		}
		specs = append(specs, spec)
	}

	// Atualizar a lista de imports
	if len(file.Imports) > 0 {
		for i, decl := range file.Decls {
			if genDecl, ok := decl.(*ast.GenDecl); ok && genDecl.Tok == token.IMPORT {
				if len(specs) > 0 {
					genDecl.Specs = specs
				} else {
					// Remover a declaração de import completamente se não houver mais imports
					file.Decls = append(file.Decls[:i], file.Decls[i+1:]...)
				}
				break
			}
		}
	}

	return fixedCount, modified
}

// fixNaming corrige problemas de nomenclatura
func (f *Fixer) fixNaming(file *ast.File, fset *token.FileSet, issues []analyzer.Issue) (int, bool) {
	fixedCount := 0
	modified := false

	// Coletar funções com problemas de nomenclatura
	namingIssues := make(map[string]string)
	for _, issue := range issues {
		if issue.FixMethod == "FixNaming" {
			// Extrair o nome da função da mensagem
			parts := strings.Split(issue.Message, ": ")
			if len(parts) > 1 {
				funcName := parts[1]
				// Converter nome_funcao para nomeFuncao (camelCase)
				camelCase := toCamelCase(funcName)
				namingIssues[funcName] = camelCase
			}
		}
	}

	if len(namingIssues) == 0 {
		return 0, false
	}

	// Corrigir nomes de funções
	for _, decl := range file.Decls {
		if funcDecl, ok := decl.(*ast.FuncDecl); ok {
			funcName := funcDecl.Name.Name
			if newName, exists := namingIssues[funcName]; exists {
				funcDecl.Name.Name = newName
				fixedCount++
				modified = true
			}
		}
	}

	return fixedCount, modified
}

// fixUndefinedMethods corrige chamadas para métodos não definidos
func (f *Fixer) fixUndefinedMethods(file *ast.File, fset *token.FileSet, issues []analyzer.Issue) (int, bool) {
	fixedCount := 0
	modified := false

	// Coletar métodos não definidos
	undefinedMethods := make(map[string]bool)
	for _, issue := range issues {
		if strings.Contains(issue.Message, "Chamada para método não definido") {
			// Extrair o nome do método da mensagem
			parts := strings.Split(issue.Message, ": ")
			if len(parts) > 1 {
				methodName := parts[1]
				undefinedMethods[methodName] = true
			}
		}
	}

	if len(undefinedMethods) == 0 {
		return 0, false
	}

	// Verificar se podemos encontrar métodos similares para substituir
	methodReplacements := f.findMethodReplacements(file, undefinedMethods)

	// Corrigir chamadas para métodos não definidos
	ast.Inspect(file, func(n ast.Node) bool {
		if callExpr, ok := n.(*ast.CallExpr); ok {
			if selExpr, ok := callExpr.Fun.(*ast.SelectorExpr); ok {
				if ident, ok := selExpr.X.(*ast.Ident); ok {
					methodName := ident.Name + "." + selExpr.Sel.Name
					
					if undefinedMethods[methodName] {
						if replacement, exists := methodReplacements[methodName]; exists {
							// Substituir pelo método correto
							selExpr.Sel.Name = replacement
							fixedCount++
							modified = true
						}
					}
				}
			}
		}
		return true
	})

	return fixedCount, modified
}

// findMethodReplacements encontra métodos similares para substituir métodos não definidos
func (f *Fixer) findMethodReplacements(file *ast.File, undefinedMethods map[string]bool) map[string]string {
	replacements := make(map[string]string)

	// Coletar todos os métodos definidos
	definedMethods := make(map[string]string)
	for _, decl := range file.Decls {
		if funcDecl, ok := decl.(*ast.FuncDecl); ok {
			if funcDecl.Recv != nil {
				// É um método
				if len(funcDecl.Recv.List) > 0 {
					var typeName string
					if starExpr, ok := funcDecl.Recv.List[0].Type.(*ast.StarExpr); ok {
						if ident, ok := starExpr.X.(*ast.Ident); ok {
							typeName = ident.Name
						}
					} else if ident, ok := funcDecl.Recv.List[0].Type.(*ast.Ident); ok {
						typeName = ident.Name
					}
					
					if typeName != "" {
						methodName := typeName + "." + funcDecl.Name.Name
						definedMethods[strings.ToLower(methodName)] = funcDecl.Name.Name
					}
				}
			}
		}
	}

	// Para cada método não definido, encontrar um método similar
	for methodName := range undefinedMethods {
		parts := strings.Split(methodName, ".")
		if len(parts) == 2 {
			typeName := parts[0]
			methodNameLower := strings.ToLower(methodName)
			
			// Verificar métodos similares
			for definedMethodLower, definedMethodName := range definedMethods {
				definedParts := strings.Split(definedMethodLower, ".")
				if len(definedParts) == 2 && definedParts[0] == strings.ToLower(typeName) {
					// Verificar similaridade
					similarity := calculateSimilarity(methodNameLower, definedMethodLower)
					if similarity > 0.7 { // 70% de similaridade
						replacements[methodName] = definedMethodName
						break
					}
				}
			}
		}
	}

	return replacements
}

// calculateSimilarity calcula a similaridade entre duas strings (algoritmo simples)
func calculateSimilarity(s1, s2 string) float64 {
	if s1 == s2 {
		return 1.0
	}
	
	// Implementação simples do algoritmo de Levenshtein
	m := len(s1)
	n := len(s2)
	
	// Criar matriz
	d := make([][]int, m+1)
	for i := range d {
		d[i] = make([]int, n+1)
	}
	
	// Inicializar primeira linha e coluna
	for i := 0; i <= m; i++ {
		d[i][0] = i
	}
	for j := 0; j <= n; j++ {
		d[0][j] = j
	}
	
	// Preencher matriz
	for j := 1; j <= n; j++ {
		for i := 1; i <= m; i++ {
			cost := 1
			if s1[i-1] == s2[j-1] {
				cost = 0
			}
			d[i][j] = min(
				d[i-1][j] + 1,     // deleção
				d[i][j-1] + 1,     // inserção
				d[i-1][j-1] + cost, // substituição
			)
		}
	}
	
	// Calcular similaridade
	maxLen := max(m, n)
	if maxLen == 0 {
		return 1.0
	}
	return 1.0 - float64(d[m][n])/float64(maxLen)
}

// min retorna o menor de três inteiros
func min(a, b, c int) int {
	if a < b {
		if a < c {
			return a
		}
		return c
	}
	if b < c {
		return b
	}
	return c
}

// max retorna o maior de dois inteiros
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// toCamelCase converte uma string com underscores para camelCase
func toCamelCase(s string) string {
	parts := strings.Split(s, "_")
	for i := 1; i < len(parts); i++ {
		if len(parts[i]) > 0 {
			parts[i] = strings.ToUpper(parts[i][:1]) + parts[i][1:]
		}
	}
	return strings.Join(parts, "")
}

// copyFile copia um arquivo de origem para destino
func copyFile(src, dst string) error {
	data, err := os.ReadFile(src)
	if err != nil {
		return err
	}
	return os.WriteFile(dst, data, 0644)
}

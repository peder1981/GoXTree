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
	"path/filepath"
	"bytes"
	"os/exec"
)

// Fixer é responsável por corrigir problemas no código
type Fixer struct {
	projectPath     string
	reporter        *reporter.Reporter
	specialPackages []string
}

// NewFixer cria um novo corretor
func NewFixer(projectPath string, reporter *reporter.Reporter) *Fixer {
	return &Fixer{
		projectPath:     projectPath,
		reporter:        reporter,
		specialPackages: []string{"github.com/gdamore/tcell/v2"}, // Pacote tcell é especial por padrão
	}
}

// SetSpecialPackages define pacotes especiais que devem ser tratados com cuidado
// durante a remoção de imports não utilizados
func (f *Fixer) SetSpecialPackages(packages []string) {
	f.specialPackages = packages
}

// FixIssues corrige problemas identificados pelo analisador
func (f *Fixer) FixIssues(issues []analyzer.Issue) (int, error) {
	fixedCount := 0

	// Agrupar problemas por arquivo
	fileIssues := make(map[string][]analyzer.Issue)
	for _, issue := range issues {
		fileIssues[issue.File] = append(fileIssues[issue.File], issue)
	}

	// Corrigir cada arquivo
	for file, issues := range fileIssues {
		fixed, err := f.fixFile(file, issues)
		if err != nil {
			return fixedCount, fmt.Errorf("erro ao corrigir arquivo %s: %w", file, err)
		}
		fixedCount += fixed
	}

	// Corrigir problemas de estilo
	styleIssues := []analyzer.StyleIssue{}
	for _, issue := range issues {
		if strings.Contains(issue.Message, "Espaços em branco no final da linha") ||
			strings.Contains(issue.Message, "Comentário deve ter espaço após //") {
			styleIssues = append(styleIssues, analyzer.StyleIssue{
				File:     issue.File,
				Line:     issue.Line,
				Message:  issue.Message,
				Severity: string(issue.Severity),
			})
		}
	}

	if len(styleIssues) > 0 {
		styleChecker := analyzer.NewStyleChecker(f.projectPath, f.reporter)
		fixed, err := styleChecker.FixStyleIssues(styleIssues)
		if err != nil {
			return fixedCount, fmt.Errorf("erro ao corrigir problemas de estilo: %w", err)
		}
		fixedCount += fixed
	}

	return fixedCount, nil
}

// FixUnusedImportsInProject remove imports não utilizados em todo o projeto
func (f *Fixer) FixUnusedImportsInProject(projectPath string) error {
	// Encontrar todos os arquivos Go no projeto
	var goFiles []string
	err := filepath.Walk(projectPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasSuffix(path, ".go") {
			goFiles = append(goFiles, path)
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("erro ao percorrer diretório do projeto: %w", err)
	}

	// Verificar quais arquivos usam tcell diretamente
	filesThatUseTcell := make(map[string]bool)
	for _, filePath := range goFiles {
		content, err := os.ReadFile(filePath)
		if err != nil {
			return fmt.Errorf("erro ao ler arquivo %s: %w", filePath, err)
		}
		
		if strings.Contains(string(content), "tcell.") {
			filesThatUseTcell[filePath] = true
		}
	}

	// Corrigir imports não utilizados em cada arquivo
	for _, filePath := range goFiles {
		// Se o arquivo usa tcell diretamente, não remover o import tcell
		if filesThatUseTcell[filePath] {
			// Verificar se o arquivo importa tcell
			fset := token.NewFileSet()
			file, err := parser.ParseFile(fset, filePath, nil, parser.ParseComments)
			if err != nil {
				return fmt.Errorf("erro ao analisar arquivo %s: %w", filePath, err)
			}
			
			// Verificar se o arquivo importa tcell
			importsTcell := false
			for _, imp := range file.Imports {
				path := strings.Trim(imp.Path.Value, "\"")
				if path == "github.com/gdamore/tcell/v2" {
					importsTcell = true
					break
				}
			}
			
			// Se o arquivo não importa tcell, adicionar o import
			if !importsTcell {
				err = f.addImport(filePath, "github.com/gdamore/tcell/v2")
				if err != nil {
					return fmt.Errorf("erro ao adicionar import tcell em %s: %w", filePath, err)
				}
				fmt.Printf("Adicionado import tcell em %s\n", filePath)
				continue
			}
		}
		
		err := f.FixUnusedImports(filePath)
		if err != nil {
			return fmt.Errorf("erro ao corrigir imports não utilizados em %s: %w", filePath, err)
		}
	}

	return nil
}

// addImport adiciona um import a um arquivo
func (f *Fixer) addImport(filePath, importPath string) error {
	// Ler o conteúdo do arquivo
	content, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("erro ao ler arquivo %s: %w", filePath, err)
	}
	
	// Converter para string
	contentStr := string(content)
	
	// Verificar se já existe um bloco de import
	importIndex := strings.Index(contentStr, "import (")
	if importIndex != -1 {
		// Encontrar o final do bloco de import
		endIndex := strings.Index(contentStr[importIndex:], ")")
		if endIndex != -1 {
			endIndex += importIndex
			// Inserir o novo import antes do fechamento do bloco
			newContent := contentStr[:endIndex] + "\n\t\"" + importPath + "\"" + contentStr[endIndex:]
			return os.WriteFile(filePath, []byte(newContent), 0644)
		}
	}
	
	// Se não encontrou um bloco de import, procurar por import único
	importIndex = strings.Index(contentStr, "import ")
	if importIndex != -1 {
		// Encontrar o final da linha
		endIndex := strings.Index(contentStr[importIndex:], "\n")
		if endIndex != -1 {
			endIndex += importIndex
			// Converter para bloco de import
			newContent := contentStr[:importIndex] + "import (\n\t" + contentStr[importIndex+7:endIndex] + "\n\t\"" + importPath + "\"\n)" + contentStr[endIndex:]
			return os.WriteFile(filePath, []byte(newContent), 0644)
		}
	}
	
	// Se não encontrou nenhum import, adicionar após o package
	packageIndex := strings.Index(contentStr, "package ")
	if packageIndex != -1 {
		// Encontrar o final da linha
		endIndex := strings.Index(contentStr[packageIndex:], "\n")
		if endIndex != -1 {
			endIndex += packageIndex
			// Adicionar o import após o package
			newContent := contentStr[:endIndex+1] + "\nimport (\n\t\"" + importPath + "\"\n)\n" + contentStr[endIndex+1:]
			return os.WriteFile(filePath, []byte(newContent), 0644)
		}
	}
	
	return fmt.Errorf("não foi possível adicionar import em %s", filePath)
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

// FixUnusedImports remove imports não utilizados do arquivo
func (f *Fixer) FixUnusedImports(filePath string) error {
	// Ler o arquivo
	content, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("erro ao ler arquivo %s: %w", filePath, err)
	}

	// Criar um arquivo temporário
	tempFile, err := os.CreateTemp(filepath.Dir(filePath), "gxtester-*.go")
	if err != nil {
		return fmt.Errorf("erro ao criar arquivo temporário: %w", err)
	}
	defer os.Remove(tempFile.Name())

	// Executar o comando goimports para remover imports não utilizados
	cmd := exec.Command("goimports", "-w", filePath)
	output, err := cmd.CombinedOutput()
	if err != nil {
		// Se goimports não estiver disponível, tentar usar go fmt
		cmd = exec.Command("go", "fmt", filePath)
		output, err = cmd.CombinedOutput()
		if err != nil {
			return fmt.Errorf("erro ao executar go fmt: %s: %w", string(output), err)
		}
	}

	// Verificar se o arquivo foi modificado
	newContent, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("erro ao ler arquivo após formatação %s: %w", filePath, err)
	}

	if bytes.Equal(content, newContent) {
		// Se o arquivo não foi modificado, tentar remover manualmente
		return f.manuallyRemoveUnusedImports(filePath)
	}

	return nil
}

// manuallyRemoveUnusedImports remove manualmente imports não utilizados
func (f *Fixer) manuallyRemoveUnusedImports(filePath string) error {
	// Analisar o arquivo
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, filePath, nil, parser.ParseComments)
	if err != nil {
		return fmt.Errorf("erro ao analisar arquivo %s: %w", filePath, err)
	}

	// Verificar imports não utilizados
	unusedImports := make(map[string]bool)
	for _, imp := range file.Imports {
		path := strings.Trim(imp.Path.Value, "\"")
		
		// Verificar se o import é um pacote especial
		isSpecialPackage := false
		for _, specialPkg := range f.specialPackages {
			if path == specialPkg {
				isSpecialPackage = true
				
				// Verificar se o arquivo contém referências ao pacote
				content, err := os.ReadFile(filePath)
				if err != nil {
					return fmt.Errorf("erro ao ler arquivo %s: %w", filePath, err)
				}
				
				// Extrair o nome do pacote do caminho completo
				pkgParts := strings.Split(specialPkg, "/")
				pkgName := pkgParts[len(pkgParts)-1]
				
				// Se o pacote tem versão (v2, v3, etc.), remover
				if strings.HasPrefix(pkgName, "v") && len(pkgName) <= 3 {
					pkgName = pkgParts[len(pkgParts)-2]
				}
				
				// Se o arquivo contém referências ao pacote, não remover o import
				if strings.Contains(string(content), pkgName+".") {
					isSpecialPackage = false
					break
				}
				
				unusedImports[path] = true
				break
			}
		}
		
		// Se não for um pacote especial, adicionar à lista de imports não utilizados
		if !isSpecialPackage {
			// Verificar se o import é usado
			// Aqui você pode adicionar lógica adicional para verificar se o import é usado
		}
	}

	// Se não houver imports não utilizados, retornar
	if len(unusedImports) == 0 {
		return nil
	}

	// Ler o conteúdo do arquivo
	content, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("erro ao ler arquivo %s: %w", filePath, err)
	}
	
	// Converter para string
	contentStr := string(content)
	
	// Remover imports não utilizados
	var newLines []string
	inImportBlock := false
	for _, line := range strings.Split(contentStr, "\n") {
		if strings.HasPrefix(strings.TrimSpace(line), "import") {
			if strings.Contains(line, "(") {
				inImportBlock = true
				newLines = append(newLines, line)
				continue
			} else if strings.Contains(line, "github.com/gdamore/tcell/v2") {
				// Ignorar linha de import único não utilizado
				continue
			}
		}

		if inImportBlock {
			if strings.Contains(line, ")") {
				inImportBlock = false
				newLines = append(newLines, line)
				continue
			}

			// Verificar se a linha contém um import não utilizado
			isUnused := false
			for imp := range unusedImports {
				if strings.Contains(line, imp) {
					isUnused = true
					break
				}
			}

			if !isUnused {
				newLines = append(newLines, line)
			}
		} else {
			newLines = append(newLines, line)
		}
	}

	// Escrever o conteúdo modificado de volta ao arquivo
	err = os.WriteFile(filePath, []byte(strings.Join(newLines, "\n")), 0644)
	if err != nil {
		return fmt.Errorf("erro ao escrever arquivo %s: %w", filePath, err)
	}

	return nil
}

package compiler

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/peder1981/advpl-tlpp-compiler/pkg/ast"
)

// SymbolKind representa o tipo de símbolo
type SymbolKind int

const (
	// SymbolKindFunction representa uma função
	SymbolKindFunction SymbolKind = iota
	// SymbolKindClass representa uma classe
	SymbolKindClass
	// SymbolKindMethod representa um método
	SymbolKindMethod
	// SymbolKindVariable representa uma variável
	SymbolKindVariable
	// SymbolKindParameter representa um parâmetro
	SymbolKindParameter
	// SymbolKindData representa um atributo de classe
	SymbolKindData
)

// SymbolInfo representa informações sobre um símbolo
type SymbolInfo struct {
	Name        string     `json:"name"`
	Kind        SymbolKind `json:"kind"`
	Line        int        `json:"line"`
	Column      int        `json:"column"`
	EndLine     int        `json:"endLine"`
	EndColumn   int        `json:"endColumn"`
	FilePath    string     `json:"filePath"`
	Description string     `json:"description"`
	Parent      string     `json:"parent,omitempty"`
	Children    []string   `json:"children,omitempty"`
	Signature   string     `json:"signature,omitempty"`
	ReturnType  string     `json:"returnType,omitempty"`
	IsStatic    bool       `json:"isStatic,omitempty"`
	IsPublic    bool       `json:"isPublic,omitempty"`
}

// DiagnosticSeverity representa a severidade de um diagnóstico
type DiagnosticSeverity int

const (
	// DiagnosticSeverityError representa um erro
	DiagnosticSeverityError DiagnosticSeverity = iota
	// DiagnosticSeverityWarning representa um aviso
	DiagnosticSeverityWarning
	// DiagnosticSeverityInformation representa uma informação
	DiagnosticSeverityInformation
	// DiagnosticSeverityHint representa uma dica
	DiagnosticSeverityHint
)

// Diagnostic representa um diagnóstico
type Diagnostic struct {
	Message   string             `json:"message"`
	Line      int                `json:"line"`
	Column    int                `json:"column"`
	EndLine   int                `json:"endLine"`
	EndColumn int                `json:"endColumn"`
	FilePath  string             `json:"filePath"`
	Severity  DiagnosticSeverity `json:"severity"`
	Code      string             `json:"code,omitempty"`
	Source    string             `json:"source,omitempty"`
}

// CompletionItem representa um item de completação
type CompletionItem struct {
	Label         string `json:"label"`
	Kind          int    `json:"kind"`
	Detail        string `json:"detail,omitempty"`
	Documentation string `json:"documentation,omitempty"`
	InsertText    string `json:"insertText,omitempty"`
	SortText      string `json:"sortText,omitempty"`
}

// IDEIntegration é responsável pela integração com o IDE
type IDEIntegration struct {
	symbols      map[string]*SymbolInfo
	diagnostics  []*Diagnostic
	completions  map[string][]*CompletionItem
	filePath     string
	currentClass string
}

// NewIDEIntegration cria uma nova instância de IDEIntegration
func NewIDEIntegration(filePath string) *IDEIntegration {
	return &IDEIntegration{
		symbols:     make(map[string]*SymbolInfo),
		diagnostics: []*Diagnostic{},
		completions: make(map[string][]*CompletionItem),
		filePath:    filePath,
	}
}

// ProcessProgram processa um programa AST para extrair informações para o IDE
func (ide *IDEIntegration) ProcessProgram(program *ast.Program) {
	for _, stmt := range program.Statements {
		ide.processStatement(stmt)
	}
}

// processStatement processa um statement para extrair informações para o IDE
func (ide *IDEIntegration) processStatement(stmt ast.Statement) {
	switch s := stmt.(type) {
	case *ast.FunctionStatement:
		ide.processFunctionStatement(s)
	case *ast.ClassStatement:
		ide.processClassStatement(s)
	case *ast.LocalStatement:
		ide.processLocalStatement(s)
	case *ast.PublicStatement:
		ide.processPublicStatement(s)
	case *ast.PrivateStatement:
		ide.processPrivateStatement(s)
	case *ast.BlockStatement:
		for _, blockStmt := range s.Statements {
			ide.processStatement(blockStmt)
		}
	}
}

// processFunctionStatement processa uma declaração de função
func (ide *IDEIntegration) processFunctionStatement(stmt *ast.FunctionStatement) {
	name := stmt.Name.Value
	line := stmt.Token.Line
	column := stmt.Token.Column
	
	// Criar informações do símbolo
	symbol := &SymbolInfo{
		Name:        name,
		Kind:        SymbolKindFunction,
		Line:        line,
		Column:      column,
		EndLine:     stmt.EndToken.Line,
		EndColumn:   stmt.EndToken.Column,
		FilePath:    ide.filePath,
		Description: fmt.Sprintf("%s função %s", if stmt.Static { "Static" } else { "" }, name),
		IsStatic:    stmt.Static,
		IsPublic:    !stmt.Static,
	}
	
	// Construir assinatura
	params := make([]string, len(stmt.Parameters))
	for i, param := range stmt.Parameters {
		params[i] = param.Value
	}
	symbol.Signature = fmt.Sprintf("%s(%s)", name, strings.Join(params, ", "))
	
	// Adicionar símbolo
	ide.symbols[name] = symbol
	
	// Processar corpo da função
	ide.processStatement(stmt.Body)
}

// processClassStatement processa uma declaração de classe
func (ide *IDEIntegration) processClassStatement(stmt *ast.ClassStatement) {
	name := stmt.Name.Value
	line := stmt.Token.Line
	column := stmt.Token.Column
	
	// Salvar classe atual
	previousClass := ide.currentClass
	ide.currentClass = name
	
	// Criar informações do símbolo
	symbol := &SymbolInfo{
		Name:        name,
		Kind:        SymbolKindClass,
		Line:        line,
		Column:      column,
		EndLine:     stmt.EndToken.Line,
		EndColumn:   stmt.EndToken.Column,
		FilePath:    ide.filePath,
		Description: fmt.Sprintf("Classe %s", name),
		IsPublic:    true,
	}
	
	// Adicionar informações de herança
	if stmt.Parent != nil {
		symbol.Parent = stmt.Parent.Value
		symbol.Description = fmt.Sprintf("Classe %s (herda de %s)", name, stmt.Parent.Value)
	}
	
	// Adicionar símbolo
	ide.symbols[name] = symbol
	
	// Processar atributos
	for _, data := range stmt.Data {
		ide.processDataStatement(data)
	}
	
	// Processar métodos
	for _, method := range stmt.Methods {
		ide.processMethodStatement(method)
		symbol.Children = append(symbol.Children, fmt.Sprintf("%s.%s", name, method.Name.Value))
	}
	
	// Restaurar classe anterior
	ide.currentClass = previousClass
}

// processMethodStatement processa uma declaração de método
func (ide *IDEIntegration) processMethodStatement(stmt *ast.MethodStatement) {
	name := stmt.Name.Value
	line := stmt.Token.Line
	column := stmt.Token.Column
	className := ide.currentClass
	
	// Criar informações do símbolo
	symbol := &SymbolInfo{
		Name:        name,
		Kind:        SymbolKindMethod,
		Line:        line,
		Column:      column,
		EndLine:     stmt.EndToken.Line,
		EndColumn:   stmt.EndToken.Column,
		FilePath:    ide.filePath,
		Description: fmt.Sprintf("Método %s da classe %s", name, className),
		Parent:      className,
		IsPublic:    true,
	}
	
	// Construir assinatura
	params := make([]string, len(stmt.Parameters))
	for i, param := range stmt.Parameters {
		params[i] = param.Value
	}
	symbol.Signature = fmt.Sprintf("%s(%s) Class %s", name, strings.Join(params, ", "), className)
	
	// Adicionar símbolo
	ide.symbols[fmt.Sprintf("%s.%s", className, name)] = symbol
	
	// Processar corpo do método
	ide.processStatement(stmt.Body)
}

// processDataStatement processa uma declaração de atributo de classe
func (ide *IDEIntegration) processDataStatement(stmt *ast.DataStatement) {
	name := stmt.Name.Value
	line := stmt.Token.Line
	column := stmt.Token.Column
	className := ide.currentClass
	
	// Criar informações do símbolo
	symbol := &SymbolInfo{
		Name:        name,
		Kind:        SymbolKindData,
		Line:        line,
		Column:      column,
		EndLine:     line,
		EndColumn:   column + len(name),
		FilePath:    ide.filePath,
		Description: fmt.Sprintf("Atributo %s da classe %s", name, className),
		Parent:      className,
		IsPublic:    true,
	}
	
	// Adicionar símbolo
	ide.symbols[fmt.Sprintf("%s.%s", className, name)] = symbol
}

// processLocalStatement processa uma declaração de variável local
func (ide *IDEIntegration) processLocalStatement(stmt *ast.LocalStatement) {
	name := stmt.Name.Value
	line := stmt.Token.Line
	column := stmt.Token.Column
	
	// Criar informações do símbolo
	symbol := &SymbolInfo{
		Name:        name,
		Kind:        SymbolKindVariable,
		Line:        line,
		Column:      column,
		EndLine:     line,
		EndColumn:   column + len(name),
		FilePath:    ide.filePath,
		Description: fmt.Sprintf("Variável local %s", name),
		IsPublic:    false,
	}
	
	// Adicionar símbolo
	ide.symbols[name] = symbol
}

// processPublicStatement processa uma declaração de variável pública
func (ide *IDEIntegration) processPublicStatement(stmt *ast.PublicStatement) {
	name := stmt.Name.Value
	line := stmt.Token.Line
	column := stmt.Token.Column
	
	// Criar informações do símbolo
	symbol := &SymbolInfo{
		Name:        name,
		Kind:        SymbolKindVariable,
		Line:        line,
		Column:      column,
		EndLine:     line,
		EndColumn:   column + len(name),
		FilePath:    ide.filePath,
		Description: fmt.Sprintf("Variável pública %s", name),
		IsPublic:    true,
	}
	
	// Adicionar símbolo
	ide.symbols[name] = symbol
}

// processPrivateStatement processa uma declaração de variável privada
func (ide *IDEIntegration) processPrivateStatement(stmt *ast.PrivateStatement) {
	name := stmt.Name.Value
	line := stmt.Token.Line
	column := stmt.Token.Column
	
	// Criar informações do símbolo
	symbol := &SymbolInfo{
		Name:        name,
		Kind:        SymbolKindVariable,
		Line:        line,
		Column:      column,
		EndLine:     line,
		EndColumn:   column + len(name),
		FilePath:    ide.filePath,
		Description: fmt.Sprintf("Variável privada %s", name),
		IsPublic:    false,
	}
	
	// Adicionar símbolo
	ide.symbols[name] = symbol
}

// AddDiagnostic adiciona um diagnóstico
func (ide *IDEIntegration) AddDiagnostic(message string, line, column, endLine, endColumn int, severity DiagnosticSeverity, code, source string) {
	diagnostic := &Diagnostic{
		Message:   message,
		Line:      line,
		Column:    column,
		EndLine:   endLine,
		EndColumn: endColumn,
		FilePath:  ide.filePath,
		Severity:  severity,
		Code:      code,
		Source:    source,
	}
	ide.diagnostics = append(ide.diagnostics, diagnostic)
}

// AddCompletionItem adiciona um item de completação
func (ide *IDEIntegration) AddCompletionItem(scope, label string, kind int, detail, documentation, insertText string) {
	item := &CompletionItem{
		Label:         label,
		Kind:          kind,
		Detail:        detail,
		Documentation: documentation,
		InsertText:    insertText,
		SortText:      label,
	}
	
	if _, ok := ide.completions[scope]; !ok {
		ide.completions[scope] = []*CompletionItem{}
	}
	ide.completions[scope] = append(ide.completions[scope], item)
}

// GetSymbols retorna os símbolos como JSON
func (ide *IDEIntegration) GetSymbols() (string, error) {
	symbols := make([]*SymbolInfo, 0, len(ide.symbols))
	for _, symbol := range ide.symbols {
		symbols = append(symbols, symbol)
	}
	
	data, err := json.Marshal(symbols)
	if err != nil {
		return "", err
	}
	
	return string(data), nil
}

// GetDiagnostics retorna os diagnósticos como JSON
func (ide *IDEIntegration) GetDiagnostics() (string, error) {
	data, err := json.Marshal(ide.diagnostics)
	if err != nil {
		return "", err
	}
	
	return string(data), nil
}

// GetCompletions retorna os itens de completação para um escopo como JSON
func (ide *IDEIntegration) GetCompletions(scope string) (string, error) {
	items, ok := ide.completions[scope]
	if !ok {
		items = []*CompletionItem{}
	}
	
	data, err := json.Marshal(items)
	if err != nil {
		return "", err
	}
	
	return string(data), nil
}

// GenerateCompletionItems gera itens de completação a partir dos símbolos
func (ide *IDEIntegration) GenerateCompletionItems() {
	// Adicionar símbolos globais ao escopo global
	for _, symbol := range ide.symbols {
		if symbol.Parent == "" {
			ide.addSymbolToCompletions("global", symbol)
		}
	}
	
	// Adicionar métodos e atributos ao escopo da classe
	for _, symbol := range ide.symbols {
		if symbol.Parent != "" {
			ide.addSymbolToCompletions(symbol.Parent, symbol)
		}
	}
}

// addSymbolToCompletions adiciona um símbolo aos itens de completação
func (ide *IDEIntegration) addSymbolToCompletions(scope string, symbol *SymbolInfo) {
	var kind int
	var detail, documentation, insertText string
	
	switch symbol.Kind {
	case SymbolKindFunction:
		kind = 3 // Função
		detail = symbol.Signature
		documentation = symbol.Description
		insertText = fmt.Sprintf("%s($0)", symbol.Name)
	case SymbolKindClass:
		kind = 5 // Classe
		detail = "Classe"
		documentation = symbol.Description
		insertText = symbol.Name
	case SymbolKindMethod:
		kind = 6 // Método
		detail = symbol.Signature
		documentation = symbol.Description
		insertText = fmt.Sprintf("%s($0)", symbol.Name)
	case SymbolKindVariable:
		kind = 7 // Variável
		detail = "Variável"
		documentation = symbol.Description
		insertText = symbol.Name
	case SymbolKindParameter:
		kind = 7 // Parâmetro
		detail = "Parâmetro"
		documentation = symbol.Description
		insertText = symbol.Name
	case SymbolKindData:
		kind = 8 // Atributo
		detail = "Atributo"
		documentation = symbol.Description
		insertText = symbol.Name
	}
	
	ide.AddCompletionItem(scope, symbol.Name, kind, detail, documentation, insertText)
}

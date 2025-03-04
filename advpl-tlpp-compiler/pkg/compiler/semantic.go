package compiler

import (
	"fmt"
	"strings"

	"advpl-tlpp-compiler/pkg/ast"
)

// SemanticError representa um erro semântico
type SemanticError struct {
	Message  string
	Line     int
	Column   int
	FileName string
}

// Error implementa a interface error
func (e *SemanticError) Error() string {
	return fmt.Sprintf("%s:%d:%d: %s", e.FileName, e.Line, e.Column, e.Message)
}

// Symbol representa um símbolo na tabela de símbolos
type Symbol struct {
	Name       string
	Type       string
	Scope      string
	IsFunction bool
	IsClass    bool
	IsMethod   bool
	IsVariable bool
	IsStatic   bool
	Line       int
	Column     int
	FileName   string
	Parent     string
}

// SymbolTable representa a tabela de símbolos
type SymbolTable struct {
	symbols       map[string]*Symbol
	scopes        []string
	currentScope  string
	currentClass  string
	currentMethod string
	fileName      string
}

// NewSymbolTable cria uma nova tabela de símbolos
func NewSymbolTable(fileName string) *SymbolTable {
	st := &SymbolTable{
		symbols:      make(map[string]*Symbol),
		scopes:       []string{"global"},
		currentScope: "global",
		fileName:     fileName,
	}
	return st
}

// EnterScope entra em um novo escopo
func (st *SymbolTable) EnterScope(scope string) {
	st.scopes = append(st.scopes, scope)
	st.currentScope = scope
}

// ExitScope sai do escopo atual
func (st *SymbolTable) ExitScope() {
	if len(st.scopes) > 1 {
		st.scopes = st.scopes[:len(st.scopes)-1]
		st.currentScope = st.scopes[len(st.scopes)-1]
	}
}

// AddSymbol adiciona um símbolo à tabela
func (st *SymbolTable) AddSymbol(name, symbolType, scope string, line, column int) *Symbol {
	fullName := st.getFullName(name)
	symbol := &Symbol{
		Name:     name,
		Type:     symbolType,
		Scope:    scope,
		Line:     line,
		Column:   column,
		FileName: st.fileName,
	}
	st.symbols[fullName] = symbol
	return symbol
}

// GetSymbol obtém um símbolo da tabela
func (st *SymbolTable) GetSymbol(name string) *Symbol {
	fullName := st.getFullName(name)
	return st.symbols[fullName]
}

// SymbolExists verifica se um símbolo existe na tabela
func (st *SymbolTable) SymbolExists(name string) bool {
	fullName := st.getFullName(name)
	_, exists := st.symbols[fullName]
	return exists
}

// getFullName obtém o nome completo de um símbolo
func (st *SymbolTable) getFullName(name string) string {
	if st.currentClass != "" && st.currentMethod != "" {
		return fmt.Sprintf("%s.%s.%s", st.currentClass, st.currentMethod, name)
	} else if st.currentClass != "" {
		return fmt.Sprintf("%s.%s", st.currentClass, name)
	} else if st.currentScope != "global" {
		return fmt.Sprintf("%s.%s", st.currentScope, name)
	}
	return name
}

// SemanticAnalyzer é responsável pela análise semântica
type SemanticAnalyzer struct {
	symbolTable *SymbolTable
	errors      []*SemanticError
	fileName    string
}

// NewSemanticAnalyzer cria um novo analisador semântico
func NewSemanticAnalyzer(fileName string) *SemanticAnalyzer {
	return &SemanticAnalyzer{
		symbolTable: NewSymbolTable(fileName),
		errors:      []*SemanticError{},
		fileName:    fileName,
	}
}

// Analyze realiza a análise semântica do programa
func (sa *SemanticAnalyzer) Analyze(program *ast.Program) []*SemanticError {
	// Primeira passagem: coletar declarações
	sa.collectDeclarations(program)

	// Segunda passagem: verificar referências
	sa.checkReferences(program)

	return sa.errors
}

// collectDeclarations coleta todas as declarações no programa
func (sa *SemanticAnalyzer) collectDeclarations(program *ast.Program) {
	for _, stmt := range program.Statements {
		sa.collectDeclarationFromStatement(stmt)
	}
}

// collectDeclarationFromStatement coleta declarações de um statement
func (sa *SemanticAnalyzer) collectDeclarationFromStatement(stmt ast.Statement) {
	switch s := stmt.(type) {
	case *ast.FunctionStatement:
		sa.collectFunctionDeclaration(s)
	case *ast.ClassStatement:
		sa.collectClassDeclaration(s)
	case *ast.LocalStatement:
		sa.collectVariableDeclaration(s, "local")
	case *ast.PublicStatement:
		sa.collectVariableDeclaration(s, "public")
	case *ast.PrivateStatement:
		sa.collectVariableDeclaration(s, "private")
	case *ast.BlockStatement:
		for _, blockStmt := range s.Statements {
			sa.collectDeclarationFromStatement(blockStmt)
		}
	}
}

// collectFunctionDeclaration coleta uma declaração de função
func (sa *SemanticAnalyzer) collectFunctionDeclaration(stmt *ast.FunctionStatement) {
	name := stmt.Name.Value
	line := stmt.Token.Line
	column := stmt.Token.Column

	// Verificar se a função já foi declarada
	if sa.symbolTable.SymbolExists(name) {
		sa.addError(fmt.Sprintf("função '%s' já declarada", name), line, column)
		return
	}

	// Adicionar função à tabela de símbolos
	symbol := sa.symbolTable.AddSymbol(name, "function", sa.symbolTable.currentScope, line, column)
	symbol.IsFunction = true
	symbol.IsStatic = stmt.Static

	// Entrar no escopo da função
	sa.symbolTable.EnterScope(name)

	// Coletar parâmetros
	for _, param := range stmt.Parameters {
		paramSymbol := sa.symbolTable.AddSymbol(param.Value, "parameter", name, param.Token.Line, param.Token.Column)
		paramSymbol.IsVariable = true
	}

	// Coletar declarações no corpo da função
	sa.collectDeclarationFromStatement(stmt.Body)

	// Sair do escopo da função
	sa.symbolTable.ExitScope()
}

// collectClassDeclaration coleta uma declaração de classe
func (sa *SemanticAnalyzer) collectClassDeclaration(stmt *ast.ClassStatement) {
	name := stmt.Name.Value
	line := stmt.Token.Line
	column := stmt.Token.Column

	// Verificar se a classe já foi declarada
	if sa.symbolTable.SymbolExists(name) {
		sa.addError(fmt.Sprintf("classe '%s' já declarada", name), line, column)
		return
	}

	// Adicionar classe à tabela de símbolos
	symbol := sa.symbolTable.AddSymbol(name, "class", sa.symbolTable.currentScope, line, column)
	symbol.IsClass = true

	// Se a classe tem um pai, verificar se ele existe
	if stmt.Parent != nil {
		parentName := stmt.Parent.Value
		symbol.Parent = parentName
	}

	// Entrar no escopo da classe
	sa.symbolTable.EnterScope(name)
	sa.symbolTable.currentClass = name

	// Coletar atributos
	for _, data := range stmt.Data {
		dataSymbol := sa.symbolTable.AddSymbol(data.Name.Value, "data", name, data.Token.Line, data.Token.Column)
		dataSymbol.IsVariable = true
	}

	// Coletar métodos
	for _, method := range stmt.Methods {
		sa.collectMethodDeclaration(method)
	}

	// Sair do escopo da classe
	sa.symbolTable.currentClass = ""
	sa.symbolTable.ExitScope()
}

// collectMethodDeclaration coleta uma declaração de método
func (sa *SemanticAnalyzer) collectMethodDeclaration(stmt *ast.MethodStatement) {
	name := stmt.Name.Value
	line := stmt.Token.Line
	column := stmt.Token.Column
	className := sa.symbolTable.currentClass

	// Verificar se o método já foi declarado nesta classe
	fullName := fmt.Sprintf("%s.%s", className, name)
	if sa.symbolTable.SymbolExists(fullName) {
		sa.addError(fmt.Sprintf("método '%s' já declarado na classe '%s'", name, className), line, column)
		return
	}

	// Adicionar método à tabela de símbolos
	symbol := sa.symbolTable.AddSymbol(name, "method", className, line, column)
	symbol.IsMethod = true

	// Entrar no escopo do método
	sa.symbolTable.currentMethod = name

	// Coletar parâmetros
	for _, param := range stmt.Parameters {
		paramSymbol := sa.symbolTable.AddSymbol(param.Value, "parameter", fullName, param.Token.Line, param.Token.Column)
		paramSymbol.IsVariable = true
	}

	// Coletar declarações no corpo do método
	sa.collectDeclarationFromStatement(stmt.Body)

	// Sair do escopo do método
	sa.symbolTable.currentMethod = ""
}

// collectVariableDeclaration coleta uma declaração de variável
func (sa *SemanticAnalyzer) collectVariableDeclaration(stmt ast.Statement, varType string) {
	var name string
	var line, column int

	switch s := stmt.(type) {
	case *ast.LocalStatement:
		name = s.Name.Value
		line = s.Token.Line
		column = s.Token.Column
	case *ast.PublicStatement:
		name = s.Name.Value
		line = s.Token.Line
		column = s.Token.Column
	case *ast.PrivateStatement:
		name = s.Name.Value
		line = s.Token.Line
		column = s.Token.Column
	default:
		return
	}

	// Verificar se a variável já foi declarada no escopo atual
	if sa.symbolTable.SymbolExists(name) {
		sa.addError(fmt.Sprintf("variável '%s' já declarada neste escopo", name), line, column)
		return
	}

	// Adicionar variável à tabela de símbolos
	symbol := sa.symbolTable.AddSymbol(name, varType, sa.symbolTable.currentScope, line, column)
	symbol.IsVariable = true
}

// checkReferences verifica as referências no programa
func (sa *SemanticAnalyzer) checkReferences(program *ast.Program) {
	for _, stmt := range program.Statements {
		sa.checkReferencesInStatement(stmt)
	}
}

// checkReferencesInStatement verifica as referências em um statement
func (sa *SemanticAnalyzer) checkReferencesInStatement(stmt ast.Statement) {
	switch s := stmt.(type) {
	case *ast.ExpressionStatement:
		sa.checkReferencesInExpression(s.Expression)
	case *ast.ReturnStatement:
		if s.Value != nil {
			sa.checkReferencesInExpression(s.Value)
		}
	case *ast.BlockStatement:
		for _, blockStmt := range s.Statements {
			sa.checkReferencesInStatement(blockStmt)
		}
	case *ast.FunctionStatement:
		sa.symbolTable.EnterScope(s.Name.Value)
		sa.checkReferencesInStatement(s.Body)
		sa.symbolTable.ExitScope()
	case *ast.ClassStatement:
		sa.symbolTable.EnterScope(s.Name.Value)
		sa.symbolTable.currentClass = s.Name.Value
		for _, method := range s.Methods {
			sa.symbolTable.currentMethod = method.Name.Value
			sa.checkReferencesInStatement(method.Body)
			sa.symbolTable.currentMethod = ""
		}
		sa.symbolTable.currentClass = ""
		sa.symbolTable.ExitScope()
	case *ast.MethodStatement:
		sa.symbolTable.currentMethod = s.Name.Value
		sa.checkReferencesInStatement(s.Body)
		sa.symbolTable.currentMethod = ""
	case *ast.LocalStatement:
		if s.Value != nil {
			sa.checkReferencesInExpression(s.Value)
		}
	case *ast.PublicStatement:
		if s.Value != nil {
			sa.checkReferencesInExpression(s.Value)
		}
	case *ast.PrivateStatement:
		if s.Value != nil {
			sa.checkReferencesInExpression(s.Value)
		}
	}
}

// checkReferencesInExpression verifica as referências em uma expressão
func (sa *SemanticAnalyzer) checkReferencesInExpression(expr ast.Expression) {
	switch e := expr.(type) {
	case *ast.Identifier:
		// Verificar se o identificador existe
		if !sa.symbolTable.SymbolExists(e.Value) {
			// Não reportar erro para "Self" em métodos
			if strings.ToLower(e.Value) == "self" && sa.symbolTable.currentMethod != "" {
				return
			}
			sa.addError(fmt.Sprintf("identificador '%s' não declarado", e.Value), e.Token.Line, e.Token.Column)
		}
	case *ast.PrefixExpression:
		sa.checkReferencesInExpression(e.Right)
	case *ast.InfixExpression:
		sa.checkReferencesInExpression(e.Left)
		sa.checkReferencesInExpression(e.Right)
	case *ast.CallExpression:
		sa.checkReferencesInExpression(e.Function)
		for _, arg := range e.Arguments {
			sa.checkReferencesInExpression(arg)
		}
	case *ast.ArrayLiteral:
		for _, elem := range e.Elements {
			sa.checkReferencesInExpression(elem)
		}
	case *ast.IndexExpression:
		sa.checkReferencesInExpression(e.Left)
		sa.checkReferencesInExpression(e.Index)
	case *ast.IfExpression:
		sa.checkReferencesInExpression(e.Condition)
		sa.checkReferencesInStatement(e.Consequence)
		for _, elseIf := range e.ElseIfs {
			sa.checkReferencesInExpression(elseIf.Condition)
			sa.checkReferencesInStatement(elseIf.Body)
		}
		if e.Alternative != nil {
			sa.checkReferencesInStatement(e.Alternative)
		}
	case *ast.WhileExpression:
		sa.checkReferencesInExpression(e.Condition)
		sa.checkReferencesInStatement(e.Body)
	case *ast.ForExpression:
		sa.checkReferencesInExpression(e.Start)
		sa.checkReferencesInExpression(e.End)
		if e.Step != nil {
			sa.checkReferencesInExpression(e.Step)
		}
		sa.checkReferencesInStatement(e.Body)
	}
}

// addError adiciona um erro à lista de erros
func (sa *SemanticAnalyzer) addError(message string, line, column int) {
	sa.errors = append(sa.errors, &SemanticError{
		Message:  message,
		Line:     line,
		Column:   column,
		FileName: sa.fileName,
	})
}

// GetErrors retorna a lista de erros
func (sa *SemanticAnalyzer) GetErrors() []*SemanticError {
	return sa.errors
}

// HasErrors verifica se há erros
func (sa *SemanticAnalyzer) HasErrors() bool {
	return len(sa.errors) > 0
}

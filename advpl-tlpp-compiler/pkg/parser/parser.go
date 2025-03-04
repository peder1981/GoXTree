package parser

import (
	"fmt"
	"strconv"

	"github.com/peder1981/advpl-tlpp-compiler/pkg/ast"
	"github.com/peder1981/advpl-tlpp-compiler/pkg/lexer"
)

// Precedências dos operadores
const (
	_ int = iota
	LOWEST
	EQUALS      // ==
	LESSGREATER // > ou <
	SUM         // +
	PRODUCT     // *
	PREFIX      // -X ou !X
	CALL        // myFunction(X)
	INDEX       // array[index]
)

// Tabela de precedências
var precedences = map[lexer.TokenType]int{
	lexer.TOKEN_EQ:    EQUALS,
	lexer.TOKEN_NOT_EQ: EQUALS,
	lexer.TOKEN_LT:     LESSGREATER,
	lexer.TOKEN_GT:     LESSGREATER,
	lexer.TOKEN_LT_EQ:  LESSGREATER,
	lexer.TOKEN_GT_EQ:  LESSGREATER,
	lexer.TOKEN_PLUS:   SUM,
	lexer.TOKEN_MINUS:  SUM,
	lexer.TOKEN_MUL:    PRODUCT,
	lexer.TOKEN_DIV:    PRODUCT,
	lexer.TOKEN_MOD:    PRODUCT,
	lexer.TOKEN_LPAREN: CALL,
	lexer.TOKEN_LBRACKET: INDEX,
}

// Tipos de funções de parsing
type (
	prefixParseFn func() ast.Expression
	infixParseFn  func(ast.Expression) ast.Expression
)

// Parser representa o analisador sintático
type Parser struct {
	l      *lexer.Lexer
	errors []string

	curToken  lexer.Token
	peekToken lexer.Token

	prefixParseFns map[lexer.TokenType]prefixParseFn
	infixParseFns  map[lexer.TokenType]infixParseFn
}

// New cria um novo analisador sintático
func New(l *lexer.Lexer) *Parser {
	p := &Parser{
		l:      l,
		errors: []string{},
	}

	// Inicializa os mapas de funções de parsing
	p.prefixParseFns = make(map[lexer.TokenType]prefixParseFn)
	p.infixParseFns = make(map[lexer.TokenType]infixParseFn)

	// Registra funções de parsing de prefixo
	p.registerPrefix(lexer.TOKEN_IDENT, p.parseIdentifier)
	p.registerPrefix(lexer.TOKEN_INT, p.parseIntegerLiteral)
	p.registerPrefix(lexer.TOKEN_FLOAT, p.parseFloatLiteral)
	p.registerPrefix(lexer.TOKEN_STRING, p.parseStringLiteral)
	p.registerPrefix(lexer.TOKEN_DATE, p.parseDateLiteral)
	p.registerPrefix(lexer.TOKEN_TRUE, p.parseBooleanLiteral)
	p.registerPrefix(lexer.TOKEN_FALSE, p.parseBooleanLiteral)
	p.registerPrefix(lexer.TOKEN_NIL, p.parseNilLiteral)
	p.registerPrefix(lexer.TOKEN_LPAREN, p.parseGroupedExpression)
	p.registerPrefix(lexer.TOKEN_BANG, p.parsePrefixExpression)
	p.registerPrefix(lexer.TOKEN_MINUS, p.parsePrefixExpression)
	p.registerPrefix(lexer.TOKEN_IF, p.parseIfExpression)
	p.registerPrefix(lexer.TOKEN_WHILE, p.parseWhileExpression)
	p.registerPrefix(lexer.TOKEN_FOR, p.parseForExpression)
	p.registerPrefix(lexer.TOKEN_LBRACKET, p.parseArrayLiteral)

	// Registra funções de parsing de infix
	p.registerInfix(lexer.TOKEN_PLUS, p.parseInfixExpression)
	p.registerInfix(lexer.TOKEN_MINUS, p.parseInfixExpression)
	p.registerInfix(lexer.TOKEN_MUL, p.parseInfixExpression)
	p.registerInfix(lexer.TOKEN_DIV, p.parseInfixExpression)
	p.registerInfix(lexer.TOKEN_MOD, p.parseInfixExpression)
	p.registerInfix(lexer.TOKEN_EQ, p.parseInfixExpression)
	p.registerInfix(lexer.TOKEN_NOT_EQ, p.parseInfixExpression)
	p.registerInfix(lexer.TOKEN_LT, p.parseInfixExpression)
	p.registerInfix(lexer.TOKEN_GT, p.parseInfixExpression)
	p.registerInfix(lexer.TOKEN_LT_EQ, p.parseInfixExpression)
	p.registerInfix(lexer.TOKEN_GT_EQ, p.parseInfixExpression)
	p.registerInfix(lexer.TOKEN_LPAREN, p.parseCallExpression)
	p.registerInfix(lexer.TOKEN_LBRACKET, p.parseIndexExpression)

	// Lê dois tokens para inicializar curToken e peekToken
	p.nextToken()
	p.nextToken()

	return p
}

// Errors retorna os erros de parsing
func (p *Parser) Errors() []string {
	return p.errors
}

// nextToken avança para o próximo token
func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

// curTokenIs verifica se o token atual é do tipo esperado
func (p *Parser) curTokenIs(t lexer.TokenType) bool {
	return p.curToken.Type == t
}

// peekTokenIs verifica se o próximo token é do tipo esperado
func (p *Parser) peekTokenIs(t lexer.TokenType) bool {
	return p.peekToken.Type == t
}

// expectPeek verifica se o próximo token é do tipo esperado e avança
func (p *Parser) expectPeek(t lexer.TokenType) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	}
	p.peekError(t)
	return false
}

// peekError adiciona um erro quando o próximo token não é o esperado
func (p *Parser) peekError(t lexer.TokenType) {
	msg := fmt.Sprintf("esperava próximo token como %s, obteve %s", t, p.peekToken.Type)
	p.errors = append(p.errors, msg)
}

// registerPrefix registra uma função de parsing de prefixo
func (p *Parser) registerPrefix(tokenType lexer.TokenType, fn prefixParseFn) {
	p.prefixParseFns[tokenType] = fn
}

// registerInfix registra uma função de parsing de infix
func (p *Parser) registerInfix(tokenType lexer.TokenType, fn infixParseFn) {
	p.infixParseFns[tokenType] = fn
}

// ParseProgram analisa o programa completo
func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}
	program.Statements = []ast.Statement{}

	for !p.curTokenIs(lexer.TOKEN_EOF) {
		stmt := p.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		p.nextToken()
	}

	return program
}

// parseStatement analisa um statement
func (p *Parser) parseStatement() ast.Statement {
	switch p.curToken.Type {
	case lexer.TOKEN_FUNCTION:
		return p.parseFunctionStatement()
	case lexer.TOKEN_CLASS:
		return p.parseClassStatement()
	case lexer.TOKEN_LOCAL:
		return p.parseLocalStatement()
	case lexer.TOKEN_PUBLIC:
		return p.parsePublicStatement()
	case lexer.TOKEN_PRIVATE:
		return p.parsePrivateStatement()
	case lexer.TOKEN_RETURN:
		return p.parseReturnStatement()
	default:
		return p.parseExpressionStatement()
	}
}

// parseIdentifier analisa um identificador
func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
}

// parseIntegerLiteral analisa um literal inteiro
func (p *Parser) parseIntegerLiteral() ast.Expression {
	lit := &ast.IntegerLiteral{Token: p.curToken}

	value, err := strconv.ParseInt(p.curToken.Literal, 0, 64)
	if err != nil {
		msg := fmt.Sprintf("não foi possível analisar %q como inteiro", p.curToken.Literal)
		p.errors = append(p.errors, msg)
		return nil
	}

	lit.Value = value
	return lit
}

// parseFloatLiteral analisa um literal de ponto flutuante
func (p *Parser) parseFloatLiteral() ast.Expression {
	lit := &ast.FloatLiteral{Token: p.curToken}

	value, err := strconv.ParseFloat(p.curToken.Literal, 64)
	if err != nil {
		msg := fmt.Sprintf("não foi possível analisar %q como float", p.curToken.Literal)
		p.errors = append(p.errors, msg)
		return nil
	}

	lit.Value = value
	return lit
}

// parseStringLiteral analisa um literal de string
func (p *Parser) parseStringLiteral() ast.Expression {
	return &ast.StringLiteral{Token: p.curToken, Value: p.curToken.Literal}
}

// parseDateLiteral analisa um literal de data
func (p *Parser) parseDateLiteral() ast.Expression {
	return &ast.DateLiteral{Token: p.curToken, Value: p.curToken.Literal}
}

// parseBooleanLiteral analisa um literal booleano
func (p *Parser) parseBooleanLiteral() ast.Expression {
	return &ast.BooleanLiteral{Token: p.curToken, Value: p.curTokenIs(lexer.TOKEN_TRUE)}
}

// parseNilLiteral analisa um literal nil
func (p *Parser) parseNilLiteral() ast.Expression {
	return &ast.NilLiteral{Token: p.curToken}
}

// parsePrefixExpression analisa uma expressão de prefixo
func (p *Parser) parsePrefixExpression() ast.Expression {
	expression := &ast.PrefixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
	}

	p.nextToken()

	expression.Right = p.parseExpression(PREFIX)

	return expression
}

// parseInfixExpression analisa uma expressão infixa
func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	expression := &ast.InfixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
		Left:     left,
	}

	precedence := p.curPrecedence()
	p.nextToken()
	expression.Right = p.parseExpression(precedence)

	return expression
}

// parseGroupedExpression analisa uma expressão agrupada
func (p *Parser) parseGroupedExpression() ast.Expression {
	p.nextToken()

	exp := p.parseExpression(LOWEST)

	if !p.expectPeek(lexer.TOKEN_RPAREN) {
		return nil
	}

	return exp
}

// parseIfExpression analisa uma expressão if
func (p *Parser) parseIfExpression() ast.Expression {
	expression := &ast.IfExpression{Token: p.curToken}

	if !p.expectPeek(lexer.TOKEN_LPAREN) {
		return nil
	}

	p.nextToken()
	expression.Condition = p.parseExpression(LOWEST)

	if !p.expectPeek(lexer.TOKEN_RPAREN) {
		return nil
	}

	expression.Consequence = p.parseBlockStatement()

	// Parse ELSEIF blocks
	for p.peekTokenIs(lexer.TOKEN_ELSEIF) {
		p.nextToken() // consume ELSEIF

		elseIf := &ast.ElseIfExpression{Token: p.curToken}

		if !p.expectPeek(lexer.TOKEN_LPAREN) {
			return nil
		}

		p.nextToken()
		elseIf.Condition = p.parseExpression(LOWEST)

		if !p.expectPeek(lexer.TOKEN_RPAREN) {
			return nil
		}

		elseIf.Body = p.parseBlockStatement()

		expression.ElseIfs = append(expression.ElseIfs, elseIf)
	}

	if p.peekTokenIs(lexer.TOKEN_ELSE) {
		p.nextToken()
		expression.Alternative = p.parseBlockStatement()
	}

	if !p.expectPeek(lexer.TOKEN_ENDIF) {
		return nil
	}

	return expression
}

// parseWhileExpression analisa uma expressão while
func (p *Parser) parseWhileExpression() ast.Expression {
	expression := &ast.WhileExpression{Token: p.curToken}

	if !p.expectPeek(lexer.TOKEN_LPAREN) {
		return nil
	}

	p.nextToken()
	expression.Condition = p.parseExpression(LOWEST)

	if !p.expectPeek(lexer.TOKEN_RPAREN) {
		return nil
	}

	expression.Body = p.parseBlockStatement()

	if !p.expectPeek(lexer.TOKEN_ENDDO) {
		return nil
	}

	return expression
}

// parseForExpression analisa uma expressão for
func (p *Parser) parseForExpression() ast.Expression {
	expression := &ast.ForExpression{Token: p.curToken}

	if !p.expectPeek(lexer.TOKEN_IDENT) {
		return nil
	}

	expression.Counter = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	if !p.expectPeek(lexer.TOKEN_FROM) {
		return nil
	}

	p.nextToken()
	expression.Start = p.parseExpression(LOWEST)

	if !p.expectPeek(lexer.TOKEN_TO) {
		return nil
	}

	p.nextToken()
	expression.End = p.parseExpression(LOWEST)

	if p.peekTokenIs(lexer.TOKEN_STEP) {
		p.nextToken()
		p.nextToken()
		expression.Step = p.parseExpression(LOWEST)
	}

	expression.Body = p.parseBlockStatement()

	if !p.expectPeek(lexer.TOKEN_NEXT) {
		return nil
	}

	return expression
}

// parseBlockStatement analisa um bloco de código
func (p *Parser) parseBlockStatement() *ast.BlockStatement {
	block := &ast.BlockStatement{Token: p.curToken}
	block.Statements = []ast.Statement{}

	p.nextToken()

	for !p.curTokenIs(lexer.TOKEN_EOF) &&
		!p.curTokenIs(lexer.TOKEN_ENDIF) &&
		!p.curTokenIs(lexer.TOKEN_ELSE) &&
		!p.curTokenIs(lexer.TOKEN_ELSEIF) &&
		!p.curTokenIs(lexer.TOKEN_ENDDO) &&
		!p.curTokenIs(lexer.TOKEN_NEXT) &&
		!p.curTokenIs(lexer.TOKEN_ENDCLASS) {
		stmt := p.parseStatement()
		if stmt != nil {
			block.Statements = append(block.Statements, stmt)
		}
		p.nextToken()
	}

	return block
}

// parseArrayLiteral analisa um literal de array
func (p *Parser) parseArrayLiteral() ast.Expression {
	array := &ast.ArrayLiteral{Token: p.curToken}

	array.Elements = p.parseExpressionList(lexer.TOKEN_RBRACKET)

	return array
}

// parseExpressionList analisa uma lista de expressões
func (p *Parser) parseExpressionList(end lexer.TokenType) []ast.Expression {
	list := []ast.Expression{}

	if p.peekTokenIs(end) {
		p.nextToken()
		return list
	}

	p.nextToken()
	list = append(list, p.parseExpression(LOWEST))

	for p.peekTokenIs(lexer.TOKEN_COMMA) {
		p.nextToken()
		p.nextToken()
		list = append(list, p.parseExpression(LOWEST))
	}

	if !p.expectPeek(end) {
		return nil
	}

	return list
}

// parseIndexExpression analisa uma expressão de acesso a índice
func (p *Parser) parseIndexExpression(left ast.Expression) ast.Expression {
	exp := &ast.IndexExpression{Token: p.curToken, Left: left}

	p.nextToken()
	exp.Index = p.parseExpression(LOWEST)

	if !p.expectPeek(lexer.TOKEN_RBRACKET) {
		return nil
	}

	return exp
}

// parseCallExpression analisa uma chamada de função
func (p *Parser) parseCallExpression(function ast.Expression) ast.Expression {
	exp := &ast.CallExpression{Token: p.curToken, Function: function}
	exp.Arguments = p.parseExpressionList(lexer.TOKEN_RPAREN)
	return exp
}

// parseFunctionStatement analisa uma declaração de função
func (p *Parser) parseFunctionStatement() *ast.FunctionStatement {
	stmt := &ast.FunctionStatement{Token: p.curToken}

	// Verificar se é uma função estática
	stmt.Static = p.curTokenIs(lexer.TOKEN_STATIC)
	if stmt.Static {
		if !p.expectPeek(lexer.TOKEN_FUNCTION) {
			return nil
		}
	}

	if !p.expectPeek(lexer.TOKEN_IDENT) {
		return nil
	}

	stmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	if !p.expectPeek(lexer.TOKEN_LPAREN) {
		return nil
	}

	stmt.Parameters = p.parseFunctionParameters()

	stmt.Body = p.parseBlockStatement()

	return stmt
}

// parseFunctionParameters analisa os parâmetros de uma função
func (p *Parser) parseFunctionParameters() []*ast.Identifier {
	identifiers := []*ast.Identifier{}

	if p.peekTokenIs(lexer.TOKEN_RPAREN) {
		p.nextToken()
		return identifiers
	}

	p.nextToken()

	ident := &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
	identifiers = append(identifiers, ident)

	for p.peekTokenIs(lexer.TOKEN_COMMA) {
		p.nextToken()
		p.nextToken()
		ident := &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
		identifiers = append(identifiers, ident)
	}

	if !p.expectPeek(lexer.TOKEN_RPAREN) {
		return nil
	}

	return identifiers
}

// parseClassStatement analisa uma declaração de classe
func (p *Parser) parseClassStatement() *ast.ClassStatement {
	stmt := &ast.ClassStatement{Token: p.curToken}

	if !p.expectPeek(lexer.TOKEN_IDENT) {
		return nil
	}

	stmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	// Verificar herança
	if p.peekTokenIs(lexer.TOKEN_FROM) {
		p.nextToken()
		if !p.expectPeek(lexer.TOKEN_IDENT) {
			return nil
		}
		stmt.Parent = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
	}

	// Parse class members
	for !p.curTokenIs(lexer.TOKEN_ENDCLASS) && !p.curTokenIs(lexer.TOKEN_EOF) {
		p.nextToken()

		switch p.curToken.Type {
		case lexer.TOKEN_METHOD:
			method := p.parseMethodStatement()
			if method != nil {
				stmt.Methods = append(stmt.Methods, method)
			}
		case lexer.TOKEN_DATA:
			data := p.parseDataStatement()
			if data != nil {
				stmt.Data = append(stmt.Data, data)
			}
		}
	}

	if !p.expectPeek(lexer.TOKEN_ENDCLASS) {
		return nil
	}

	return stmt
}

// parseMethodStatement analisa uma declaração de método
func (p *Parser) parseMethodStatement() *ast.MethodStatement {
	stmt := &ast.MethodStatement{Token: p.curToken}

	if !p.expectPeek(lexer.TOKEN_IDENT) {
		return nil
	}

	stmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	if !p.expectPeek(lexer.TOKEN_LPAREN) {
		return nil
	}

	stmt.Parameters = p.parseFunctionParameters()
	stmt.Body = p.parseBlockStatement()

	return stmt
}

// parseDataStatement analisa uma declaração de atributo de classe
func (p *Parser) parseDataStatement() *ast.DataStatement {
	stmt := &ast.DataStatement{Token: p.curToken}

	if !p.expectPeek(lexer.TOKEN_IDENT) {
		return nil
	}

	stmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	return stmt
}

// parseLocalStatement analisa uma declaração de variável local
func (p *Parser) parseLocalStatement() *ast.LocalStatement {
	stmt := &ast.LocalStatement{Token: p.curToken}

	if !p.expectPeek(lexer.TOKEN_IDENT) {
		return nil
	}

	stmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	if p.peekTokenIs(lexer.TOKEN_ASSIGN) {
		p.nextToken()
		p.nextToken()
		stmt.Value = p.parseExpression(LOWEST)
	}

	return stmt
}

// parsePublicStatement analisa uma declaração de variável pública
func (p *Parser) parsePublicStatement() *ast.PublicStatement {
	stmt := &ast.PublicStatement{Token: p.curToken}

	if !p.expectPeek(lexer.TOKEN_IDENT) {
		return nil
	}

	stmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	if p.peekTokenIs(lexer.TOKEN_ASSIGN) {
		p.nextToken()
		p.nextToken()
		stmt.Value = p.parseExpression(LOWEST)
	}

	return stmt
}

// parsePrivateStatement analisa uma declaração de variável privada
func (p *Parser) parsePrivateStatement() *ast.PrivateStatement {
	stmt := &ast.PrivateStatement{Token: p.curToken}

	if !p.expectPeek(lexer.TOKEN_IDENT) {
		return nil
	}

	stmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	if p.peekTokenIs(lexer.TOKEN_ASSIGN) {
		p.nextToken()
		p.nextToken()
		stmt.Value = p.parseExpression(LOWEST)
	}

	return stmt
}

// parseReturnStatement analisa uma declaração de retorno
func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	stmt := &ast.ReturnStatement{Token: p.curToken}

	p.nextToken()

	stmt.Value = p.parseExpression(LOWEST)

	return stmt
}

// parseExpressionStatement analisa um statement de expressão
func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	stmt := &ast.ExpressionStatement{Token: p.curToken}
	stmt.Expression = p.parseExpression(LOWEST)
	return stmt
}

// parseExpression analisa uma expressão
func (p *Parser) parseExpression(precedence int) ast.Expression {
	prefix := p.prefixParseFns[p.curToken.Type]
	if prefix == nil {
		p.noPrefixParseFnError(p.curToken.Type)
		return nil
	}
	leftExp := prefix()

	for !p.peekTokenIs(lexer.TOKEN_SEMICOLON) && precedence < p.peekPrecedence() {
		infix := p.infixParseFns[p.peekToken.Type]
		if infix == nil {
			return leftExp
		}

		p.nextToken()
		leftExp = infix(leftExp)
	}

	return leftExp
}

func (p *Parser) peekPrecedence() int {
	if p, ok := precedences[p.peekToken.Type]; ok {
		return p
	}
	return LOWEST
}

func (p *Parser) curPrecedence() int {
	if p, ok := precedences[p.curToken.Type]; ok {
		return p
	}
	return LOWEST
}

func (p *Parser) noPrefixParseFnError(t lexer.TokenType) {
	msg := fmt.Sprintf("nenhuma função de parse de prefixo para %s encontrada", t)
	p.errors = append(p.errors, msg)
}
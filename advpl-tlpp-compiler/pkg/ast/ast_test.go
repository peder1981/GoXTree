package ast

import (
	"testing"
)

func TestNode(t *testing.T) {
	// Testa a interface Node
	var nodes []Node = []Node{
		&Program{},
		&FunctionStatement{},
		&ClassStatement{},
		&MethodStatement{},
		&DataStatement{},
		&VariableStatement{},
		&ReturnStatement{},
		&ExpressionStatement{},
		&BlockStatement{},
		&IfExpression{},
		&WhileExpression{},
		&ForExpression{},
		&CallExpression{},
		&IndexExpression{},
		&PrefixExpression{},
		&InfixExpression{},
		&AssignExpression{},
		&Identifier{},
		&IntegerLiteral{},
		&FloatLiteral{},
		&StringLiteral{},
		&BooleanLiteral{},
		&NilLiteral{},
		&ArrayLiteral{},
	}

	// Verifica se todos os nós implementam a interface Node
	for i, node := range nodes {
		if node.TokenLiteral() == "" {
			// TokenLiteral deve retornar uma string vazia para nós vazios
			continue
		}
		
		nodeStr := node.String()
		if nodeStr == "" {
			t.Errorf("Node[%d] (%T).String() retornou string vazia", i, node)
		}
	}
}

func TestProgram(t *testing.T) {
	program := &Program{
		Statements: []Statement{
			&FunctionStatement{
				Token: Token{Type: FUNCTION, Literal: "Function"},
				Name: &Identifier{
					Token: Token{Type: IDENT, Literal: "Teste"},
					Value: "Teste",
				},
				Parameters: []*Identifier{},
				Body: &BlockStatement{
					Token: Token{Type: LBRACE, Literal: "{"},
					Statements: []Statement{
						&ReturnStatement{
							Token: Token{Type: RETURN, Literal: "Return"},
							ReturnValue: &IntegerLiteral{
								Token: Token{Type: INT, Literal: "5"},
								Value: 5,
							},
						},
					},
				},
			},
		},
	}

	if program.TokenLiteral() != "Function" {
		t.Errorf("program.TokenLiteral() incorreto. Esperado='Function', obtido='%s'", 
			program.TokenLiteral())
	}

	programStr := program.String()
	if len(programStr) == 0 {
		t.Errorf("program.String() retornou string vazia")
	}
}

func TestFunctionStatement(t *testing.T) {
	function := &FunctionStatement{
		Token: Token{Type: FUNCTION, Literal: "Function"},
		Name: &Identifier{
			Token: Token{Type: IDENT, Literal: "Soma"},
			Value: "Soma",
		},
		Parameters: []*Identifier{
			{
				Token: Token{Type: IDENT, Literal: "a"},
				Value: "a",
			},
			{
				Token: Token{Type: IDENT, Literal: "b"},
				Value: "b",
			},
		},
		Body: &BlockStatement{
			Token: Token{Type: LBRACE, Literal: "{"},
			Statements: []Statement{
				&ReturnStatement{
					Token: Token{Type: RETURN, Literal: "Return"},
					ReturnValue: &InfixExpression{
						Token: Token{Type: PLUS, Literal: "+"},
						Left: &Identifier{
							Token: Token{Type: IDENT, Literal: "a"},
							Value: "a",
						},
						Operator: "+",
						Right: &Identifier{
							Token: Token{Type: IDENT, Literal: "b"},
							Value: "b",
						},
					},
				},
			},
		},
	}

	if function.TokenLiteral() != "Function" {
		t.Errorf("function.TokenLiteral() incorreto. Esperado='Function', obtido='%s'", 
			function.TokenLiteral())
	}

	if function.Name.Value != "Soma" {
		t.Errorf("function.Name.Value incorreto. Esperado='Soma', obtido='%s'", 
			function.Name.Value)
	}

	if len(function.Parameters) != 2 {
		t.Errorf("function.Parameters tem tamanho incorreto. Esperado=2, obtido=%d", 
			len(function.Parameters))
	}

	params := []string{"a", "b"}
	for i, param := range function.Parameters {
		if param.Value != params[i] {
			t.Errorf("function.Parameters[%d].Value incorreto. Esperado='%s', obtido='%s'", 
				i, params[i], param.Value)
		}
	}

	functionStr := function.String()
	if len(functionStr) == 0 {
		t.Errorf("function.String() retornou string vazia")
	}
}

func TestClassStatement(t *testing.T) {
	class := &ClassStatement{
		Token: Token{Type: CLASS, Literal: "Class"},
		Name: &Identifier{
			Token: Token{Type: IDENT, Literal: "Pessoa"},
			Value: "Pessoa",
		},
		Parent: &Identifier{
			Token: Token{Type: IDENT, Literal: "Base"},
			Value: "Base",
		},
		Body: &BlockStatement{
			Token: Token{Type: LBRACE, Literal: "{"},
			Statements: []Statement{
				&DataStatement{
					Token: Token{Type: DATA, Literal: "Data"},
					Name: &Identifier{
						Token: Token{Type: IDENT, Literal: "Nome"},
						Value: "Nome",
					},
				},
				&MethodStatement{
					Token: Token{Type: METHOD, Literal: "Method"},
					Name: &Identifier{
						Token: Token{Type: IDENT, Literal: "New"},
						Value: "New",
					},
					Parameters: []*Identifier{
						{
							Token: Token{Type: IDENT, Literal: "cNome"},
							Value: "cNome",
						},
					},
					Body: &BlockStatement{
						Token: Token{Type: LBRACE, Literal: "{"},
						Statements: []Statement{
							&ExpressionStatement{
								Token: Token{Type: IDENT, Literal: "::Nome"},
								Expression: &AssignExpression{
									Token: Token{Type: ASSIGN, Literal: ":="},
									Left: &Identifier{
										Token: Token{Type: IDENT, Literal: "::Nome"},
										Value: "::Nome",
									},
									Value: &Identifier{
										Token: Token{Type: IDENT, Literal: "cNome"},
										Value: "cNome",
									},
								},
							},
						},
					},
					IsConstructor: true,
				},
			},
		},
	}

	if class.TokenLiteral() != "Class" {
		t.Errorf("class.TokenLiteral() incorreto. Esperado='Class', obtido='%s'", 
			class.TokenLiteral())
	}

	if class.Name.Value != "Pessoa" {
		t.Errorf("class.Name.Value incorreto. Esperado='Pessoa', obtido='%s'", 
			class.Name.Value)
	}

	if class.Parent.Value != "Base" {
		t.Errorf("class.Parent.Value incorreto. Esperado='Base', obtido='%s'", 
			class.Parent.Value)
	}

	classStr := class.String()
	if len(classStr) == 0 {
		t.Errorf("class.String() retornou string vazia")
	}
}

func TestExpressions(t *testing.T) {
	// Testa InfixExpression
	infix := &InfixExpression{
		Token: Token{Type: PLUS, Literal: "+"},
		Left: &IntegerLiteral{
			Token: Token{Type: INT, Literal: "5"},
			Value: 5,
		},
		Operator: "+",
		Right: &IntegerLiteral{
			Token: Token{Type: INT, Literal: "10"},
			Value: 10,
		},
	}

	if infix.String() != "(5 + 10)" {
		t.Errorf("infix.String() incorreto. Esperado='(5 + 10)', obtido='%s'", 
			infix.String())
	}

	// Testa PrefixExpression
	prefix := &PrefixExpression{
		Token: Token{Type: BANG, Literal: "!"},
		Operator: "!",
		Right: &BooleanLiteral{
			Token: Token{Type: TRUE, Literal: ".T."},
			Value: true,
		},
	}

	if prefix.String() != "(!.T.)" {
		t.Errorf("prefix.String() incorreto. Esperado='(!.T.)', obtido='%s'", 
			prefix.String())
	}

	// Testa IfExpression
	ifExpr := &IfExpression{
		Token: Token{Type: IF, Literal: "If"},
		Condition: &InfixExpression{
			Token: Token{Type: GT, Literal: ">"},
			Left: &Identifier{
				Token: Token{Type: IDENT, Literal: "x"},
				Value: "x",
			},
			Operator: ">",
			Right: &IntegerLiteral{
				Token: Token{Type: INT, Literal: "5"},
				Value: 5,
			},
		},
		Consequence: &BlockStatement{
			Token: Token{Type: LBRACE, Literal: "{"},
			Statements: []Statement{
				&ExpressionStatement{
					Token: Token{Type: IDENT, Literal: "y"},
					Expression: &AssignExpression{
						Token: Token{Type: ASSIGN, Literal: ":="},
						Left: &Identifier{
							Token: Token{Type: IDENT, Literal: "y"},
							Value: "y",
						},
						Value: &IntegerLiteral{
							Token: Token{Type: INT, Literal: "10"},
							Value: 10,
						},
					},
				},
			},
		},
		Alternative: &BlockStatement{
			Token: Token{Type: LBRACE, Literal: "{"},
			Statements: []Statement{
				&ExpressionStatement{
					Token: Token{Type: IDENT, Literal: "y"},
					Expression: &AssignExpression{
						Token: Token{Type: ASSIGN, Literal: ":="},
						Left: &Identifier{
							Token: Token{Type: IDENT, Literal: "y"},
							Value: "y",
						},
						Value: &IntegerLiteral{
							Token: Token{Type: INT, Literal: "5"},
							Value: 5,
						},
					},
				},
			},
		},
	}

	ifStr := ifExpr.String()
	if len(ifStr) == 0 {
		t.Errorf("ifExpr.String() retornou string vazia")
	}
}

func TestLiterals(t *testing.T) {
	// Testa IntegerLiteral
	intLit := &IntegerLiteral{
		Token: Token{Type: INT, Literal: "42"},
		Value: 42,
	}

	if intLit.String() != "42" {
		t.Errorf("intLit.String() incorreto. Esperado='42', obtido='%s'", 
			intLit.String())
	}

	// Testa FloatLiteral
	floatLit := &FloatLiteral{
		Token: Token{Type: FLOAT, Literal: "3.14"},
		Value: 3.14,
	}

	if floatLit.String() != "3.14" {
		t.Errorf("floatLit.String() incorreto. Esperado='3.14', obtido='%s'", 
			floatLit.String())
	}

	// Testa StringLiteral
	strLit := &StringLiteral{
		Token: Token{Type: STRING, Literal: "\"hello\""},
		Value: "hello",
	}

	if strLit.String() != "\"hello\"" {
		t.Errorf("strLit.String() incorreto. Esperado='\"hello\"', obtido='%s'", 
			strLit.String())
	}

	// Testa BooleanLiteral
	boolLit := &BooleanLiteral{
		Token: Token{Type: TRUE, Literal: ".T."},
		Value: true,
	}

	if boolLit.String() != ".T." {
		t.Errorf("boolLit.String() incorreto. Esperado='.T.', obtido='%s'", 
			boolLit.String())
	}

	// Testa NilLiteral
	nilLit := &NilLiteral{
		Token: Token{Type: NIL, Literal: "Nil"},
	}

	if nilLit.String() != "Nil" {
		t.Errorf("nilLit.String() incorreto. Esperado='Nil', obtido='%s'", 
			nilLit.String())
	}

	// Testa ArrayLiteral
	arrayLit := &ArrayLiteral{
		Token: Token{Type: LBRACKET, Literal: "{"},
		Elements: []Expression{
			&IntegerLiteral{
				Token: Token{Type: INT, Literal: "1"},
				Value: 1,
			},
			&IntegerLiteral{
				Token: Token{Type: INT, Literal: "2"},
				Value: 2,
			},
			&IntegerLiteral{
				Token: Token{Type: INT, Literal: "3"},
				Value: 3,
			},
		},
	}

	if arrayLit.String() != "{1, 2, 3}" {
		t.Errorf("arrayLit.String() incorreto. Esperado='{1, 2, 3}', obtido='%s'", 
			arrayLit.String())
	}
}

// Estruturas necessárias para os testes

type Token struct {
	Type    string
	Literal string
	Line    int
	Column  int
}

// Constantes para tipos de tokens
const (
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"
	
	// Identificadores + literais
	IDENT   = "IDENT"
	INT     = "INT"
	FLOAT   = "FLOAT"
	STRING  = "STRING"
	
	// Operadores
	ASSIGN   = ":="
	PLUS     = "+"
	MINUS    = "-"
	BANG     = "!"
	ASTERISK = "*"
	SLASH    = "/"
	
	LT = "<"
	GT = ">"
	
	EQ     = "=="
	NOT_EQ = "!="
	
	// Delimitadores
	COMMA     = ","
	SEMICOLON = ";"
	COLON     = ":"
	
	LPAREN   = "("
	RPAREN   = ")"
	LBRACE   = "{"
	RBRACE   = "}"
	LBRACKET = "["
	RBRACKET = "]"
	
	// Palavras-chave
	FUNCTION   = "FUNCTION"
	STATIC     = "STATIC"
	CLASS      = "CLASS"
	METHOD     = "METHOD"
	DATA       = "DATA"
	FROM       = "FROM"
	LOCAL      = "LOCAL"
	PUBLIC     = "PUBLIC"
	PRIVATE    = "PRIVATE"
	RETURN     = "RETURN"
	TRUE       = "TRUE"
	FALSE      = "FALSE"
	IF         = "IF"
	ELSE       = "ELSE"
	ELSEIF     = "ELSEIF"
	WHILE      = "WHILE"
	FOR        = "FOR"
	TO         = "TO"
	STEP       = "STEP"
	NIL        = "NIL"
	CONSTRUCTOR = "CONSTRUCTOR"
)

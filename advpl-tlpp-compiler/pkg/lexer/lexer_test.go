package lexer

import (
	"testing"
)

func TestNextToken(t *testing.T) {
	input := `Function Teste(cParam1, nParam2)
	Local cVar := "Teste"
	Local nVar := 10
	Local lVar := .T.
	Local dVar := Date()
	
	If (nVar > 5)
		cVar := "Maior que 5"
	ElseIf (nVar == 5)
		cVar := "Igual a 5"
	Else
		cVar := "Menor que 5"
	EndIf
	
	While (nVar > 0)
		nVar--
	EndDo
	
	For nI := 1 To 10 Step 2
		cVar += Str(nI)
	Next
	
	Return cVar
EndFunction

Class Exemplo From Base
	Data cAtributo
	Data nAtributo
	
	Method New() Constructor
	Method Metodo1()
EndClass

Method New() Class Exemplo
	::cAtributo := "Valor Inicial"
	::nAtributo := 0
	Return Self
	
Method Metodo1() Class Exemplo
	Local xRet := Nil
	
	If (::nAtributo > 0)
		xRet := ::cAtributo + Str(::nAtributo)
	Else
		xRet := ::cAtributo
	EndIf
	
	Return xRet`

	tests := []struct {
		expectedType    TokenType
		expectedLiteral string
	}{
		{TOKEN_FUNCTION, "Function"},
		{TOKEN_IDENT, "Teste"},
		{TOKEN_LPAREN, "("},
		{TOKEN_IDENT, "cParam1"},
		{TOKEN_COMMA, ","},
		{TOKEN_IDENT, "nParam2"},
		{TOKEN_RPAREN, ")"},
		{TOKEN_LOCAL, "Local"},
		{TOKEN_IDENT, "cVar"},
		{TOKEN_ASSIGN, ":="},
		{TOKEN_STRING, "\"Teste\""},
		{TOKEN_LOCAL, "Local"},
		{TOKEN_IDENT, "nVar"},
		{TOKEN_ASSIGN, ":="},
		{TOKEN_INT, "10"},
		{TOKEN_LOCAL, "Local"},
		{TOKEN_IDENT, "lVar"},
		{TOKEN_ASSIGN, ":="},
		{TOKEN_TRUE, ".T."},
		{TOKEN_LOCAL, "Local"},
		{TOKEN_IDENT, "dVar"},
		{TOKEN_ASSIGN, ":="},
		{TOKEN_IDENT, "Date"},
		{TOKEN_LPAREN, "("},
		{TOKEN_RPAREN, ")"},
		{TOKEN_IF, "If"},
		{TOKEN_LPAREN, "("},
		{TOKEN_IDENT, "nVar"},
		{TOKEN_GT, ">"},
		{TOKEN_INT, "5"},
		{TOKEN_RPAREN, ")"},
		{TOKEN_IDENT, "cVar"},
		{TOKEN_ASSIGN, ":="},
		{TOKEN_STRING, "\"Maior que 5\""},
		{TOKEN_ELSEIF, "ElseIf"},
		{TOKEN_LPAREN, "("},
		{TOKEN_IDENT, "nVar"},
		{TOKEN_EQ, "=="},
		{TOKEN_INT, "5"},
		{TOKEN_RPAREN, ")"},
		{TOKEN_IDENT, "cVar"},
		{TOKEN_ASSIGN, ":="},
		{TOKEN_STRING, "\"Igual a 5\""},
		{TOKEN_ELSE, "Else"},
		{TOKEN_IDENT, "cVar"},
		{TOKEN_ASSIGN, ":="},
		{TOKEN_STRING, "\"Menor que 5\""},
		{TOKEN_ENDIF, "EndIf"},
		{TOKEN_WHILE, "While"},
		{TOKEN_LPAREN, "("},
		{TOKEN_IDENT, "nVar"},
		{TOKEN_GT, ">"},
		{TOKEN_INT, "0"},
		{TOKEN_RPAREN, ")"},
		{TOKEN_IDENT, "nVar"},
		{TOKEN_MINUS, "-"},
		{TOKEN_MINUS, "-"},
		{TOKEN_ENDDO, "EndDo"},
		{TOKEN_FOR, "For"},
		{TOKEN_IDENT, "nI"},
		{TOKEN_ASSIGN, ":="},
		{TOKEN_INT, "1"},
		{TOKEN_TO, "To"},
		{TOKEN_INT, "10"},
		{TOKEN_STEP, "Step"},
		{TOKEN_INT, "2"},
		{TOKEN_IDENT, "cVar"},
		{TOKEN_PLUS, "+"},
		{TOKEN_ASSIGN, "="},
		{TOKEN_IDENT, "Str"},
		{TOKEN_LPAREN, "("},
		{TOKEN_IDENT, "nI"},
		{TOKEN_RPAREN, ")"},
		{TOKEN_NEXT, "Next"},
		{TOKEN_RETURN, "Return"},
		{TOKEN_IDENT, "cVar"},
		{TOKEN_ENDFUNCTION, "EndFunction"},
		{TOKEN_CLASS, "Class"},
		{TOKEN_IDENT, "Exemplo"},
		{TOKEN_FROM, "From"},
		{TOKEN_IDENT, "Base"},
		{TOKEN_DATA, "Data"},
		{TOKEN_IDENT, "cAtributo"},
		{TOKEN_DATA, "Data"},
		{TOKEN_IDENT, "nAtributo"},
		{TOKEN_METHOD, "Method"},
		{TOKEN_IDENT, "New"},
		{TOKEN_LPAREN, "("},
		{TOKEN_RPAREN, ")"},
		{TOKEN_IDENT, "Constructor"},
		{TOKEN_METHOD, "Method"},
		{TOKEN_IDENT, "Metodo1"},
		{TOKEN_LPAREN, "("},
		{TOKEN_RPAREN, ")"},
		{TOKEN_ENDCLASS, "EndClass"},
		{TOKEN_METHOD, "Method"},
		{TOKEN_IDENT, "New"},
		{TOKEN_LPAREN, "("},
		{TOKEN_RPAREN, ")"},
		{TOKEN_IDENT, "Class"},
		{TOKEN_IDENT, "Exemplo"},
		{TOKEN_IDENT, "::"},
		{TOKEN_IDENT, "cAtributo"},
		{TOKEN_ASSIGN, ":="},
		{TOKEN_STRING, "\"Valor Inicial\""},
		{TOKEN_IDENT, "::"},
		{TOKEN_IDENT, "nAtributo"},
		{TOKEN_ASSIGN, ":="},
		{TOKEN_INT, "0"},
		{TOKEN_RETURN, "Return"},
		{TOKEN_SELF, "Self"},
		{TOKEN_METHOD, "Method"},
		{TOKEN_IDENT, "Metodo1"},
		{TOKEN_LPAREN, "("},
		{TOKEN_RPAREN, ")"},
		{TOKEN_IDENT, "Class"},
		{TOKEN_IDENT, "Exemplo"},
		{TOKEN_LOCAL, "Local"},
		{TOKEN_IDENT, "xRet"},
		{TOKEN_ASSIGN, ":="},
		{TOKEN_NIL, "Nil"},
		{TOKEN_IF, "If"},
		{TOKEN_LPAREN, "("},
		{TOKEN_IDENT, "::"},
		{TOKEN_IDENT, "nAtributo"},
		{TOKEN_GT, ">"},
		{TOKEN_INT, "0"},
		{TOKEN_RPAREN, ")"},
		{TOKEN_IDENT, "xRet"},
		{TOKEN_ASSIGN, ":="},
		{TOKEN_IDENT, "::"},
		{TOKEN_IDENT, "cAtributo"},
		{TOKEN_PLUS, "+"},
		{TOKEN_IDENT, "Str"},
		{TOKEN_LPAREN, "("},
		{TOKEN_IDENT, "::"},
		{TOKEN_IDENT, "nAtributo"},
		{TOKEN_RPAREN, ")"},
		{TOKEN_ELSE, "Else"},
		{TOKEN_IDENT, "xRet"},
		{TOKEN_ASSIGN, ":="},
		{TOKEN_IDENT, "::"},
		{TOKEN_IDENT, "cAtributo"},
		{TOKEN_ENDIF, "EndIf"},
		{TOKEN_RETURN, "Return"},
		{TOKEN_IDENT, "xRet"},
		{TOKEN_EOF, ""},
	}

	l := New(input, "test.prw")

	for i, tt := range tests {
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("teste[%d] - tipo de token incorreto. esperado=%d, obtido=%d (%s)",
				i, tt.expectedType, tok.Type, tok.Literal)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("teste[%d] - literal incorreto. esperado=%q, obtido=%q",
				i, tt.expectedLiteral, tok.Literal)
		}
	}
}

func TestTokenPosition(t *testing.T) {
	input := `Function Teste()
	Local x := 10
	Return x
EndFunction`

	tests := []struct {
		expectedType   TokenType
		expectedLine   int
		expectedColumn int
	}{
		{TOKEN_FUNCTION, 1, 1},
		{TOKEN_IDENT, 1, 10},
		{TOKEN_LPAREN, 1, 15},
		{TOKEN_RPAREN, 1, 16},
		{TOKEN_LOCAL, 2, 2},
		{TOKEN_IDENT, 2, 8},
		{TOKEN_ASSIGN, 2, 10},
		{TOKEN_INT, 2, 13},
		{TOKEN_RETURN, 3, 2},
		{TOKEN_IDENT, 3, 9},
		{TOKEN_ENDFUNCTION, 4, 1},
		{TOKEN_EOF, 4, 12},
	}

	l := New(input, "test.prw")

	for i, tt := range tests {
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("teste[%d] - tipo de token incorreto. esperado=%d, obtido=%d",
				i, tt.expectedType, tok.Type)
		}

		if tok.Line != tt.expectedLine {
			t.Fatalf("teste[%d] - linha incorreta. esperado=%d, obtido=%d",
				i, tt.expectedLine, tok.Line)
		}

		if tok.Column != tt.expectedColumn {
			t.Fatalf("teste[%d] - coluna incorreta. esperado=%d, obtido=%d",
				i, tt.expectedColumn, tok.Column)
		}
	}
}

func TestComments(t *testing.T) {
	input := `// Comentário de linha
Function Teste() // Outro comentário
	/* Comentário
	   de bloco */
	Local x := 10
	Return x
EndFunction`

	expected := []TokenType{
		TOKEN_FUNCTION,
		TOKEN_IDENT,
		TOKEN_LPAREN,
		TOKEN_RPAREN,
		TOKEN_LOCAL,
		TOKEN_IDENT,
		TOKEN_ASSIGN,
		TOKEN_INT,
		TOKEN_RETURN,
		TOKEN_IDENT,
		TOKEN_ENDFUNCTION,
		TOKEN_EOF,
	}

	l := New(input, "test.prw")

	for i, exp := range expected {
		tok := l.NextToken()
		if tok.Type != exp {
			t.Fatalf("teste[%d] - tipo de token incorreto. esperado=%d, obtido=%d (%s)",
				i, exp, tok.Type, tok.Literal)
		}
	}
}

func TestStringLiterals(t *testing.T) {
	input := `"String simples"
"String com \"escape\""
'String com aspas simples'`

	tests := []struct {
		expectedType    TokenType
		expectedLiteral string
	}{
		{TOKEN_STRING, "\"String simples\""},
		{TOKEN_STRING, "\"String com \\\"escape\\\"\""},
		{TOKEN_STRING, "'String com aspas simples'"},
		{TOKEN_EOF, ""},
	}

	l := New(input, "test.prw")

	for i, tt := range tests {
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("teste[%d] - tipo de token incorreto. esperado=%d, obtido=%d",
				i, tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("teste[%d] - literal incorreto. esperado=%q, obtido=%q",
				i, tt.expectedLiteral, tok.Literal)
		}
	}
}

func TestNumbers(t *testing.T) {
	input := `123
12.34
0.5
.5`

	tests := []struct {
		expectedType    TokenType
		expectedLiteral string
	}{
		{TOKEN_INT, "123"},
		{TOKEN_FLOAT, "12.34"},
		{TOKEN_FLOAT, "0.5"},
		{TOKEN_FLOAT, ".5"},
		{TOKEN_EOF, ""},
	}

	l := New(input, "test.prw")

	for i, tt := range tests {
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("teste[%d] - tipo de token incorreto. esperado=%d, obtido=%d",
				i, tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("teste[%d] - literal incorreto. esperado=%q, obtido=%q",
				i, tt.expectedLiteral, tok.Literal)
		}
	}
}

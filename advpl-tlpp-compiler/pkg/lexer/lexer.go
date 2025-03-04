package lexer

import (
	"fmt"
	"unicode"
	"unicode/utf8"
)

// TokenType representa o tipo de um token
type TokenType int

// Constantes para os tipos de tokens
const (
	TOKEN_ILLEGAL TokenType = iota
	TOKEN_EOF
	TOKEN_IDENT
	TOKEN_INT
	TOKEN_FLOAT
	TOKEN_STRING
	TOKEN_DATE
	
	// Operadores
	TOKEN_ASSIGN   // :=
	TOKEN_PLUS     // +
	TOKEN_MINUS    // -
	TOKEN_BANG     // !
	TOKEN_ASTERISK // *
	TOKEN_SLASH    // /
	TOKEN_PERCENT  // %
	
	TOKEN_EQ     // ==
	TOKEN_NOT_EQ // !=
	TOKEN_LT     // <
	TOKEN_GT     // >
	TOKEN_LE     // <=
	TOKEN_GE     // >=
	
	// Delimitadores
	TOKEN_COMMA     // ,
	TOKEN_SEMICOLON // ;
	TOKEN_COLON     // :
	TOKEN_LPAREN    // (
	TOKEN_RPAREN    // )
	TOKEN_LBRACE    // {
	TOKEN_RBRACE    // }
	TOKEN_LBRACKET  // [
	TOKEN_RBRACKET  // ]
	
	// Palavras-chave
	TOKEN_FUNCTION
	TOKEN_STATIC
	TOKEN_RETURN
	TOKEN_IF
	TOKEN_ELSE
	TOKEN_ELSEIF
	TOKEN_ENDIF
	TOKEN_WHILE
	TOKEN_DO
	TOKEN_ENDDO
	TOKEN_FOR
	TOKEN_NEXT
	TOKEN_LOCAL
	TOKEN_PUBLIC
	TOKEN_PRIVATE
	TOKEN_CLASS
	TOKEN_ENDCLASS
	TOKEN_METHOD
	TOKEN_DATA
	TOKEN_FROM
	TOKEN_TO
	TOKEN_STEP
	TOKEN_NIL
	TOKEN_TRUE
	TOKEN_FALSE
	TOKEN_SELF
	TOKEN_NEW
)

// Token representa um token léxico
type Token struct {
	Type    TokenType
	Literal string
	Line    int
	Column  int
	File    string
}

// String retorna uma representação em string do token
func (t Token) String() string {
	return fmt.Sprintf("Token{Type: %v, Literal: %q, Line: %d, Column: %d, File: %s}",
		t.Type, t.Literal, t.Line, t.Column, t.File)
}

// Lexer é responsável pela análise léxica
type Lexer struct {
	input        string
	position     int  // posição atual no input (aponta para o caractere atual)
	readPosition int  // posição atual de leitura (após o caractere atual)
	ch           rune // caractere atual sendo examinado
	
	line   int // linha atual
	column int // coluna atual
	
	file string // nome do arquivo
}

// New cria um novo Lexer
func New(input string, file string) *Lexer {
	l := &Lexer{
		input:  input,
		file:   file,
		line:   1,
		column: 0,
	}
	l.readChar()
	return l
}

// readChar lê o próximo caractere e avança a posição no input
func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0 // EOF
	} else {
		r, width := utf8.DecodeRuneInString(l.input[l.readPosition:])
		l.ch = r
		l.position = l.readPosition
		l.readPosition += width
	}
	l.column++
}

// peekChar retorna o próximo caractere sem avançar a posição
func (l *Lexer) peekChar() rune {
	if l.readPosition >= len(l.input) {
		return 0 // EOF
	}
	r, _ := utf8.DecodeRuneInString(l.input[l.readPosition:])
	return r
}

// NextToken retorna o próximo token
func (l *Lexer) NextToken() Token {
	var tok Token
	
	l.skipWhitespace()
	
	// Registrar a posição do token
	tok.Line = l.line
	tok.Column = l.column
	tok.File = l.file
	
	switch l.ch {
	case '=':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = Token{Type: TOKEN_EQ, Literal: string(ch) + string(l.ch)}
		} else {
			tok = newToken(TOKEN_ASSIGN, l.ch)
		}
	case '+':
		tok = newToken(TOKEN_PLUS, l.ch)
	case '-':
		tok = newToken(TOKEN_MINUS, l.ch)
	case '!':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = Token{Type: TOKEN_NOT_EQ, Literal: string(ch) + string(l.ch)}
		} else {
			tok = newToken(TOKEN_BANG, l.ch)
		}
	case '*':
		tok = newToken(TOKEN_ASTERISK, l.ch)
	case '/':
		// Verificar se é um comentário
		if l.peekChar() == '/' {
			l.skipLineComment()
			return l.NextToken()
		} else if l.peekChar() == '*' {
			l.skipBlockComment()
			return l.NextToken()
		} else {
			tok = newToken(TOKEN_SLASH, l.ch)
		}
	case '%':
		tok = newToken(TOKEN_PERCENT, l.ch)
	case '<':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = Token{Type: TOKEN_LE, Literal: string(ch) + string(l.ch)}
		} else {
			tok = newToken(TOKEN_LT, l.ch)
		}
	case '>':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = Token{Type: TOKEN_GE, Literal: string(ch) + string(l.ch)}
		} else {
			tok = newToken(TOKEN_GT, l.ch)
		}
	case ',':
		tok = newToken(TOKEN_COMMA, l.ch)
	case ';':
		tok = newToken(TOKEN_SEMICOLON, l.ch)
	case ':':
		tok = newToken(TOKEN_COLON, l.ch)
	case '(':
		tok = newToken(TOKEN_LPAREN, l.ch)
	case ')':
		tok = newToken(TOKEN_RPAREN, l.ch)
	case '{':
		tok = newToken(TOKEN_LBRACE, l.ch)
	case '}':
		tok = newToken(TOKEN_RBRACE, l.ch)
	case '[':
		tok = newToken(TOKEN_LBRACKET, l.ch)
	case ']':
		tok = newToken(TOKEN_RBRACKET, l.ch)
	case '"', '\'':
		tok.Type = TOKEN_STRING
		tok.Literal = l.readString(l.ch)
	case 0:
		tok.Literal = ""
		tok.Type = TOKEN_EOF
	default:
		if isLetter(l.ch) {
			tok.Literal = l.readIdentifier()
			tok.Type = lookupIdent(tok.Literal)
			return tok
		} else if isDigit(l.ch) {
			return l.readNumber()
		} else {
			tok = newToken(TOKEN_ILLEGAL, l.ch)
		}
	}
	
	l.readChar()
	return tok
}

// newToken cria um novo token
func newToken(tokenType TokenType, ch rune) Token {
	return Token{Type: tokenType, Literal: string(ch)}
}

// skipWhitespace pula espaços em branco
func (l *Lexer) skipWhitespace() {
	for unicode.IsSpace(l.ch) {
		if l.ch == '\n' {
			l.line++
			l.column = 0
		}
		l.readChar()
	}
}

// skipLineComment pula comentários de linha
func (l *Lexer) skipLineComment() {
	l.readChar() // pular o primeiro '/'
	l.readChar() // pular o segundo '/'
	
	for l.ch != '\n' && l.ch != 0 {
		l.readChar()
	}
}

// skipBlockComment pula comentários de bloco
func (l *Lexer) skipBlockComment() {
	l.readChar() // pular o '/'
	l.readChar() // pular o '*'
	
	for {
		if l.ch == 0 {
			// EOF antes do fim do comentário
			break
		}
		
		if l.ch == '*' && l.peekChar() == '/' {
			l.readChar() // pular o '*'
			l.readChar() // pular o '/'
			break
		}
		
		if l.ch == '\n' {
			l.line++
			l.column = 0
		}
		
		l.readChar()
	}
}

// readIdentifier lê um identificador
func (l *Lexer) readIdentifier() string {
	position := l.position
	for isLetter(l.ch) || isDigit(l.ch) || l.ch == '_' {
		l.readChar()
	}
	return l.input[position:l.position]
}

// readNumber lê um número (inteiro ou decimal)
func (l *Lexer) readNumber() Token {
	position := l.position
	isFloat := false
	
	// Ler a parte inteira
	for isDigit(l.ch) {
		l.readChar()
	}
	
	// Verificar se é um número decimal
	if l.ch == '.' && isDigit(l.peekChar()) {
		isFloat = true
		l.readChar() // consumir o ponto
		
		// Ler a parte decimal
		for isDigit(l.ch) {
			l.readChar()
		}
	}
	
	// Verificar se é uma data (DD/MM/YYYY ou YYYY/MM/DD)
	if l.ch == '/' && isDigit(l.peekChar()) {
		// Potencial data
		dateStart := position
		dateContent := l.input[dateStart:l.position]
		
		// Ler o primeiro separador
		l.readChar() // consumir a primeira '/'
		
		// Ler o segundo grupo de dígitos
		digitCount := 0
		for isDigit(l.ch) {
			digitCount++
			l.readChar()
		}
		
		// Verificar o segundo separador
		if l.ch == '/' && isDigit(l.peekChar()) && digitCount > 0 {
			l.readChar() // consumir a segunda '/'
			
			// Ler o terceiro grupo de dígitos
			digitCount = 0
			for isDigit(l.ch) {
				digitCount++
				l.readChar()
			}
			
			if digitCount > 0 {
				// É uma data válida
				return Token{
					Type:    TOKEN_DATE,
					Literal: l.input[dateStart:l.position],
				}
			}
		}
		
		// Se chegou aqui, não é uma data válida
		// Retornar como número normal
		return Token{
			Type:    TOKEN_FLOAT,
			Literal: dateContent,
		}
	}
	
	if isFloat {
		return Token{
			Type:    TOKEN_FLOAT,
			Literal: l.input[position:l.position],
		}
	}
	
	return Token{
		Type:    TOKEN_INT,
		Literal: l.input[position:l.position],
	}
}

// readString lê uma string
func (l *Lexer) readString(delimiter rune) string {
	position := l.position + 1 // pular o delimitador inicial
	
	for {
		l.readChar()
		if l.ch == delimiter || l.ch == 0 {
			break
		}
		
		// Lidar com caracteres de escape
		if l.ch == '\\' {
			l.readChar() // pular o '\'
			// Aqui podemos adicionar lógica para interpretar caracteres de escape
		}
		
		if l.ch == '\n' {
			l.line++
			l.column = 0
		}
	}
	
	return l.input[position:l.position]
}

// isLetter verifica se um caractere é uma letra
func isLetter(ch rune) bool {
	return unicode.IsLetter(ch)
}

// isDigit verifica se um caractere é um dígito
func isDigit(ch rune) bool {
	return unicode.IsDigit(ch)
}

// lookupIdent verifica se um identificador é uma palavra-chave
func lookupIdent(ident string) TokenType {
	keywords := map[string]TokenType{
		"function": TOKEN_FUNCTION,
		"static":   TOKEN_STATIC,
		"return":   TOKEN_RETURN,
		"if":       TOKEN_IF,
		"else":     TOKEN_ELSE,
		"elseif":   TOKEN_ELSEIF,
		"endif":    TOKEN_ENDIF,
		"while":    TOKEN_WHILE,
		"do":       TOKEN_DO,
		"enddo":    TOKEN_ENDDO,
		"for":      TOKEN_FOR,
		"next":     TOKEN_NEXT,
		"local":    TOKEN_LOCAL,
		"public":   TOKEN_PUBLIC,
		"private":  TOKEN_PRIVATE,
		"class":    TOKEN_CLASS,
		"endclass": TOKEN_ENDCLASS,
		"method":   TOKEN_METHOD,
		"data":     TOKEN_DATA,
		"from":     TOKEN_FROM,
		"to":       TOKEN_TO,
		"step":     TOKEN_STEP,
		"nil":      TOKEN_NIL,
		"true":     TOKEN_TRUE,
		"false":    TOKEN_FALSE,
		"self":     TOKEN_SELF,
		"new":      TOKEN_NEW,
	}
	
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	
	return TOKEN_IDENT
}

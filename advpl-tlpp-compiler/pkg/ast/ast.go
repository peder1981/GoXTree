package ast

import (
	"bytes"
	"strings"

	"github.com/peder1981/advpl-tlpp-compiler/pkg/lexer"
)

// Node representa um nó na AST
type Node interface {
	TokenLiteral() string
	String() string
}

// Statement representa uma declaração
type Statement interface {
	Node
	statementNode()
}

// Expression representa uma expressão
type Expression interface {
	Node
	expressionNode()
}

// Program representa um programa AdvPL/TLPP
type Program struct {
	Statements []Statement
}

// TokenLiteral retorna o literal do token do primeiro statement
func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	}
	return ""
}

// String retorna uma representação em string do programa
func (p *Program) String() string {
	var out bytes.Buffer

	for _, s := range p.Statements {
		out.WriteString(s.String())
	}

	return out.String()
}

// Identifier representa um identificador
type Identifier struct {
	Token lexer.Token // token.TOKEN_IDENT
	Value string
}

func (i *Identifier) expressionNode() {}

// TokenLiteral retorna o literal do token
func (i *Identifier) TokenLiteral() string { return i.Token.Literal }

// String retorna uma representação em string do identificador
func (i *Identifier) String() string { return i.Value }

// FunctionStatement representa uma declaração de função
type FunctionStatement struct {
	Token      lexer.Token // token.TOKEN_FUNCTION
	Name       *Identifier
	Parameters []*Identifier
	Body       *BlockStatement
	Static     bool
}

func (fs *FunctionStatement) statementNode() {}

// TokenLiteral retorna o literal do token
func (fs *FunctionStatement) TokenLiteral() string { return fs.Token.Literal }

// String retorna uma representação em string da declaração de função
func (fs *FunctionStatement) String() string {
	var out bytes.Buffer

	params := []string{}
	for _, p := range fs.Parameters {
		params = append(params, p.String())
	}

	staticStr := ""
	if fs.Static {
		staticStr = "STATIC "
	}

	out.WriteString(staticStr)
	out.WriteString("FUNCTION ")
	out.WriteString(fs.Name.String())
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") ")
	out.WriteString(fs.Body.String())

	return out.String()
}

// BlockStatement representa um bloco de código
type BlockStatement struct {
	Token      lexer.Token // token.TOKEN_LBRACE ou token que inicia o bloco
	Statements []Statement
}

func (bs *BlockStatement) statementNode() {}

// TokenLiteral retorna o literal do token
func (bs *BlockStatement) TokenLiteral() string { return bs.Token.Literal }

// String retorna uma representação em string do bloco
func (bs *BlockStatement) String() string {
	var out bytes.Buffer

	for _, s := range bs.Statements {
		out.WriteString(s.String())
	}

	return out.String()
}

// ReturnStatement representa uma declaração de retorno
type ReturnStatement struct {
	Token lexer.Token // token.TOKEN_RETURN
	Value Expression
}

func (rs *ReturnStatement) statementNode() {}

// TokenLiteral retorna o literal do token
func (rs *ReturnStatement) TokenLiteral() string { return rs.Token.Literal }

// String retorna uma representação em string da declaração de retorno
func (rs *ReturnStatement) String() string {
	var out bytes.Buffer

	out.WriteString(rs.TokenLiteral() + " ")

	if rs.Value != nil {
		out.WriteString(rs.Value.String())
	}

	out.WriteString("\n")

	return out.String()
}

// ExpressionStatement representa uma expressão usada como statement
type ExpressionStatement struct {
	Token      lexer.Token // o primeiro token da expressão
	Expression Expression
}

func (es *ExpressionStatement) statementNode() {}

// TokenLiteral retorna o literal do token
func (es *ExpressionStatement) TokenLiteral() string { return es.Token.Literal }

// String retorna uma representação em string do statement de expressão
func (es *ExpressionStatement) String() string {
	if es.Expression != nil {
		return es.Expression.String()
	}
	return ""
}

// IntegerLiteral representa um literal inteiro
type IntegerLiteral struct {
	Token lexer.Token
	Value int64
}

func (il *IntegerLiteral) expressionNode() {}

// TokenLiteral retorna o literal do token
func (il *IntegerLiteral) TokenLiteral() string { return il.Token.Literal }

// String retorna uma representação em string do literal inteiro
func (il *IntegerLiteral) String() string { return il.Token.Literal }

// FloatLiteral representa um literal de ponto flutuante
type FloatLiteral struct {
	Token lexer.Token
	Value float64
}

func (fl *FloatLiteral) expressionNode() {}

// TokenLiteral retorna o literal do token
func (fl *FloatLiteral) TokenLiteral() string { return fl.Token.Literal }

// String retorna uma representação em string do literal de ponto flutuante
func (fl *FloatLiteral) String() string { return fl.Token.Literal }

// StringLiteral representa um literal de string
type StringLiteral struct {
	Token lexer.Token
	Value string
}

func (sl *StringLiteral) expressionNode() {}

// TokenLiteral retorna o literal do token
func (sl *StringLiteral) TokenLiteral() string { return sl.Token.Literal }

// String retorna uma representação em string do literal de string
func (sl *StringLiteral) String() string { return "\"" + sl.Value + "\"" }

// DateLiteral representa um literal de data
type DateLiteral struct {
	Token lexer.Token
	Value string // formato AAAAMMDD
}

func (dl *DateLiteral) expressionNode() {}

// TokenLiteral retorna o literal do token
func (dl *DateLiteral) TokenLiteral() string { return dl.Token.Literal }

// String retorna uma representação em string do literal de data
func (dl *DateLiteral) String() string { return dl.Token.Literal }

// BooleanLiteral representa um literal booleano
type BooleanLiteral struct {
	Token lexer.Token
	Value bool
}

func (bl *BooleanLiteral) expressionNode() {}

// TokenLiteral retorna o literal do token
func (bl *BooleanLiteral) TokenLiteral() string { return bl.Token.Literal }

// String retorna uma representação em string do literal booleano
func (bl *BooleanLiteral) String() string {
	if bl.Value {
		return ".T."
	}
	return ".F."
}

// NilLiteral representa o valor nil
type NilLiteral struct {
	Token lexer.Token
}

func (nl *NilLiteral) expressionNode() {}

// TokenLiteral retorna o literal do token
func (nl *NilLiteral) TokenLiteral() string { return nl.Token.Literal }

// String retorna uma representação em string do nil
func (nl *NilLiteral) String() string { return "NIL" }

// PrefixExpression representa uma expressão de prefixo
type PrefixExpression struct {
	Token    lexer.Token // O token de prefixo, ex: !
	Operator string
	Right    Expression
}

func (pe *PrefixExpression) expressionNode() {}

// TokenLiteral retorna o literal do token
func (pe *PrefixExpression) TokenLiteral() string { return pe.Token.Literal }

// String retorna uma representação em string da expressão de prefixo
func (pe *PrefixExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(pe.Operator)
	out.WriteString(pe.Right.String())
	out.WriteString(")")

	return out.String()
}

// InfixExpression representa uma expressão infixa
type InfixExpression struct {
	Token    lexer.Token // O token de operação, ex: +
	Left     Expression
	Operator string
	Right    Expression
}

func (ie *InfixExpression) expressionNode() {}

// TokenLiteral retorna o literal do token
func (ie *InfixExpression) TokenLiteral() string { return ie.Token.Literal }

// String retorna uma representação em string da expressão infixa
func (ie *InfixExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(ie.Left.String())
	out.WriteString(" " + ie.Operator + " ")
	out.WriteString(ie.Right.String())
	out.WriteString(")")

	return out.String()
}

// CallExpression representa uma chamada de função
type CallExpression struct {
	Token     lexer.Token // O token '('
	Function  Expression  // Identificador ou expressão de função
	Arguments []Expression
}

func (ce *CallExpression) expressionNode() {}

// TokenLiteral retorna o literal do token
func (ce *CallExpression) TokenLiteral() string { return ce.Token.Literal }

// String retorna uma representação em string da chamada de função
func (ce *CallExpression) String() string {
	var out bytes.Buffer

	args := []string{}
	for _, a := range ce.Arguments {
		args = append(args, a.String())
	}

	out.WriteString(ce.Function.String())
	out.WriteString("(")
	out.WriteString(strings.Join(args, ", "))
	out.WriteString(")")

	return out.String()
}

// IfExpression representa uma expressão if
type IfExpression struct {
	Token       lexer.Token // O token 'if'
	Condition   Expression
	Consequence *BlockStatement
	Alternative *BlockStatement
	ElseIfs     []*ElseIfExpression
}

func (ie *IfExpression) expressionNode() {}

// TokenLiteral retorna o literal do token
func (ie *IfExpression) TokenLiteral() string { return ie.Token.Literal }

// String retorna uma representação em string da expressão if
func (ie *IfExpression) String() string {
	var out bytes.Buffer

	out.WriteString("IF ")
	out.WriteString(ie.Condition.String())
	out.WriteString("\n")
	out.WriteString(ie.Consequence.String())

	for _, elseif := range ie.ElseIfs {
		out.WriteString(elseif.String())
	}

	if ie.Alternative != nil {
		out.WriteString("ELSE\n")
		out.WriteString(ie.Alternative.String())
	}

	out.WriteString("ENDIF")

	return out.String()
}

// ElseIfExpression representa uma expressão elseif
type ElseIfExpression struct {
	Token     lexer.Token // O token 'elseif'
	Condition Expression
	Body      *BlockStatement
}

func (eie *ElseIfExpression) expressionNode() {}

// TokenLiteral retorna o literal do token
func (eie *ElseIfExpression) TokenLiteral() string { return eie.Token.Literal }

// String retorna uma representação em string da expressão elseif
func (eie *ElseIfExpression) String() string {
	var out bytes.Buffer

	out.WriteString("ELSEIF ")
	out.WriteString(eie.Condition.String())
	out.WriteString("\n")
	out.WriteString(eie.Body.String())

	return out.String()
}

// WhileExpression representa uma expressão while
type WhileExpression struct {
	Token     lexer.Token // O token 'while'
	Condition Expression
	Body      *BlockStatement
}

func (we *WhileExpression) expressionNode() {}

// TokenLiteral retorna o literal do token
func (we *WhileExpression) TokenLiteral() string { return we.Token.Literal }

// String retorna uma representação em string da expressão while
func (we *WhileExpression) String() string {
	var out bytes.Buffer

	out.WriteString("WHILE ")
	out.WriteString(we.Condition.String())
	out.WriteString("\n")
	out.WriteString(we.Body.String())
	out.WriteString("ENDDO")

	return out.String()
}

// ForExpression representa uma expressão for
type ForExpression struct {
	Token    lexer.Token // O token 'for'
	Counter  *Identifier
	Start    Expression
	End      Expression
	Step     Expression
	Body     *BlockStatement
}

func (fe *ForExpression) expressionNode() {}

// TokenLiteral retorna o literal do token
func (fe *ForExpression) TokenLiteral() string { return fe.Token.Literal }

// String retorna uma representação em string da expressão for
func (fe *ForExpression) String() string {
	var out bytes.Buffer

	out.WriteString("FOR ")
	out.WriteString(fe.Counter.String())
	out.WriteString(" := ")
	out.WriteString(fe.Start.String())
	out.WriteString(" TO ")
	out.WriteString(fe.End.String())
	
	if fe.Step != nil {
		out.WriteString(" STEP ")
		out.WriteString(fe.Step.String())
	}
	
	out.WriteString("\n")
	out.WriteString(fe.Body.String())
	out.WriteString("NEXT")

	return out.String()
}

// ClassStatement representa uma declaração de classe
type ClassStatement struct {
	Token    lexer.Token // O token 'class'
	Name     *Identifier
	Parent   *Identifier // Classe pai, se houver
	Methods  []*MethodStatement
	Data     []*DataStatement
}

func (cs *ClassStatement) statementNode() {}

// TokenLiteral retorna o literal do token
func (cs *ClassStatement) TokenLiteral() string { return cs.Token.Literal }

// String retorna uma representação em string da declaração de classe
func (cs *ClassStatement) String() string {
	var out bytes.Buffer

	out.WriteString("CLASS ")
	out.WriteString(cs.Name.String())
	
	if cs.Parent != nil {
		out.WriteString(" FROM ")
		out.WriteString(cs.Parent.String())
	}
	
	out.WriteString("\n")
	
	for _, data := range cs.Data {
		out.WriteString(data.String())
		out.WriteString("\n")
	}
	
	for _, method := range cs.Methods {
		out.WriteString(method.String())
		out.WriteString("\n")
	}
	
	out.WriteString("ENDCLASS")

	return out.String()
}

// MethodStatement representa uma declaração de método
type MethodStatement struct {
	Token      lexer.Token // O token 'method'
	Name       *Identifier
	Parameters []*Identifier
	Body       *BlockStatement
}

func (ms *MethodStatement) statementNode() {}

// TokenLiteral retorna o literal do token
func (ms *MethodStatement) TokenLiteral() string { return ms.Token.Literal }

// String retorna uma representação em string da declaração de método
func (ms *MethodStatement) String() string {
	var out bytes.Buffer

	params := []string{}
	for _, p := range ms.Parameters {
		params = append(params, p.String())
	}

	out.WriteString("METHOD ")
	out.WriteString(ms.Name.String())
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") ")
	out.WriteString(ms.Body.String())

	return out.String()
}

// DataStatement representa uma declaração de atributo de classe
type DataStatement struct {
	Token lexer.Token // O token 'data'
	Name  *Identifier
}

func (ds *DataStatement) statementNode() {}

// TokenLiteral retorna o literal do token
func (ds *DataStatement) TokenLiteral() string { return ds.Token.Literal }

// String retorna uma representação em string da declaração de atributo
func (ds *DataStatement) String() string {
	var out bytes.Buffer

	out.WriteString("DATA ")
	out.WriteString(ds.Name.String())

	return out.String()
}

// LocalStatement representa uma declaração de variável local
type LocalStatement struct {
	Token lexer.Token // O token 'local'
	Name  *Identifier
	Value Expression
}

func (ls *LocalStatement) statementNode() {}

// TokenLiteral retorna o literal do token
func (ls *LocalStatement) TokenLiteral() string { return ls.Token.Literal }

// String retorna uma representação em string da declaração local
func (ls *LocalStatement) String() string {
	var out bytes.Buffer

	out.WriteString("LOCAL ")
	out.WriteString(ls.Name.String())
	
	if ls.Value != nil {
		out.WriteString(" := ")
		out.WriteString(ls.Value.String())
	}

	return out.String()
}

// PublicStatement representa uma declaração de variável pública
type PublicStatement struct {
	Token lexer.Token // O token 'public'
	Name  *Identifier
	Value Expression
}

func (ps *PublicStatement) statementNode() {}

// TokenLiteral retorna o literal do token
func (ps *PublicStatement) TokenLiteral() string { return ps.Token.Literal }

// String retorna uma representação em string da declaração pública
func (ps *PublicStatement) String() string {
	var out bytes.Buffer

	out.WriteString("PUBLIC ")
	out.WriteString(ps.Name.String())
	
	if ps.Value != nil {
		out.WriteString(" := ")
		out.WriteString(ps.Value.String())
	}

	return out.String()
}

// PrivateStatement representa uma declaração de variável privada
type PrivateStatement struct {
	Token lexer.Token // O token 'private'
	Name  *Identifier
	Value Expression
}

func (ps *PrivateStatement) statementNode() {}

// TokenLiteral retorna o literal do token
func (ps *PrivateStatement) TokenLiteral() string { return ps.Token.Literal }

// String retorna uma representação em string da declaração privada
func (ps *PrivateStatement) String() string {
	var out bytes.Buffer

	out.WriteString("PRIVATE ")
	out.WriteString(ps.Name.String())
	
	if ps.Value != nil {
		out.WriteString(" := ")
		out.WriteString(ps.Value.String())
	}

	return out.String()
}

// ArrayLiteral representa um literal de array
type ArrayLiteral struct {
	Token    lexer.Token // O token '['
	Elements []Expression
}

func (al *ArrayLiteral) expressionNode() {}

// TokenLiteral retorna o literal do token
func (al *ArrayLiteral) TokenLiteral() string { return al.Token.Literal }

// String retorna uma representação em string do literal de array
func (al *ArrayLiteral) String() string {
	var out bytes.Buffer

	elements := []string{}
	for _, el := range al.Elements {
		elements = append(elements, el.String())
	}

	out.WriteString("[")
	out.WriteString(strings.Join(elements, ", "))
	out.WriteString("]")

	return out.String()
}

// IndexExpression representa uma expressão de acesso a índice
type IndexExpression struct {
	Token lexer.Token // O token '['
	Left  Expression
	Index Expression
}

func (ie *IndexExpression) expressionNode() {}

// TokenLiteral retorna o literal do token
func (ie *IndexExpression) TokenLiteral() string { return ie.Token.Literal }

// String retorna uma representação em string da expressão de índice
func (ie *IndexExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(ie.Left.String())
	out.WriteString("[")
	out.WriteString(ie.Index.String())
	out.WriteString("])")

	return out.String()
}

// AssignmentExpression representa uma expressão de atribuição
type AssignmentExpression struct {
	Token lexer.Token // O token ':='
	Left  Expression
	Value Expression
}

func (ae *AssignmentExpression) expressionNode() {}

// TokenLiteral retorna o literal do token
func (ae *AssignmentExpression) TokenLiteral() string { return ae.Token.Literal }

// String retorna uma representação em string da expressão de atribuição
func (ae *AssignmentExpression) String() string {
	var out bytes.Buffer

	out.WriteString(ae.Left.String())
	out.WriteString(" := ")
	out.WriteString(ae.Value.String())

	return out.String()
}

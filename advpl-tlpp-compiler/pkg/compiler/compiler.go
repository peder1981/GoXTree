package compiler

import (
	"fmt"
	"strings"
	"time"

	"advpl-tlpp-compiler/pkg/ast"
)

// Options representa as opções de compilação
type Options struct {
	Verbose      bool
	Optimize     bool
	Dialect      string
	IncludeDirs  []string
	CheckSyntax  bool
	GenerateDocs bool
}

// Stats mantém estatísticas sobre o código compilado
type Stats struct {
	FunctionCount int
	VariableCount int
	LineCount     int
}

// Compiler representa o compilador
type Compiler struct {
	program  *ast.Program
	options  Options
	stats    Stats
	output   strings.Builder
	indent   int
	includes map[string]bool
}

// New cria um novo compilador
func New(program *ast.Program, options Options) *Compiler {
	return &Compiler{
		program:  program,
		options:  options,
		stats:    Stats{},
		includes: make(map[string]bool),
	}
}

// GetStats retorna as estatísticas de compilação
func (c *Compiler) GetStats() Stats {
	return c.stats
}

// Compile compila o programa e retorna o código objeto
func (c *Compiler) Compile() (string, error) {
	// Processar cada statement do programa
	for _, stmt := range c.program.Statements {
		if err := c.compileStatement(stmt); err != nil {
			return "", err
		}
	}

	return c.output.String(), nil
}

// compileStatement compila um statement
func (c *Compiler) compileStatement(stmt ast.Statement) error {
	switch s := stmt.(type) {
	case *ast.FunctionStatement:
		return c.compileFunctionStatement(s)
	case *ast.ClassStatement:
		return c.compileClassStatement(s)
	case *ast.LocalStatement:
		return c.compileLocalStatement(s)
	case *ast.PublicStatement:
		return c.compilePublicStatement(s)
	case *ast.PrivateStatement:
		return c.compilePrivateStatement(s)
	case *ast.ReturnStatement:
		return c.compileReturnStatement(s)
	case *ast.ExpressionStatement:
		return c.compileExpressionStatement(s)
	default:
		return fmt.Errorf("tipo de statement não suportado: %T", stmt)
	}
}

// compileFunctionStatement compila uma declaração de função
func (c *Compiler) compileFunctionStatement(stmt *ast.FunctionStatement) error {
	c.stats.FunctionCount++

	// Gerar cabeçalho da função
	if stmt.Static {
		c.writeLine("Static Function %s(", stmt.Name.Value)
	} else {
		c.writeLine("Function %s(", stmt.Name.Value)
	}

	// Gerar parâmetros
	params := make([]string, len(stmt.Parameters))
	for i, param := range stmt.Parameters {
		params[i] = param.Value
	}
	c.writeLine("%s)", strings.Join(params, ", "))

	// Gerar corpo da função
	c.indent++
	if err := c.compileBlockStatement(stmt.Body); err != nil {
		return err
	}
	c.indent--

	c.writeLine("Return")
	c.writeLine("")

	return nil
}

// compileClassStatement compila uma declaração de classe
func (c *Compiler) compileClassStatement(stmt *ast.ClassStatement) error {
	// Gerar cabeçalho da classe
	if stmt.Parent != nil {
		c.writeLine("Class %s From %s", stmt.Name.Value, stmt.Parent.Value)
	} else {
		c.writeLine("Class %s", stmt.Name.Value)
	}

	c.indent++

	// Gerar atributos
	for _, data := range stmt.Data {
		c.writeLine("Data %s", data.Name.Value)
	}

	// Gerar métodos
	for _, method := range stmt.Methods {
		if err := c.compileMethodStatement(method); err != nil {
			return err
		}
	}

	c.indent--
	c.writeLine("EndClass")
	c.writeLine("")

	return nil
}

// compileMethodStatement compila uma declaração de método
func (c *Compiler) compileMethodStatement(stmt *ast.MethodStatement) error {
	c.stats.FunctionCount++

	// Gerar cabeçalho do método
	c.writeLine("Method %s(", stmt.Name.Value)

	// Gerar parâmetros
	params := make([]string, len(stmt.Parameters))
	for i, param := range stmt.Parameters {
		params[i] = param.Value
	}
	c.writeLine("%s) Class %s", strings.Join(params, ", "), "self")

	// Gerar corpo do método
	c.indent++
	if err := c.compileBlockStatement(stmt.Body); err != nil {
		return err
	}
	c.indent--

	c.writeLine("Return")
	c.writeLine("")

	return nil
}

// compileLocalStatement compila uma declaração de variável local
func (c *Compiler) compileLocalStatement(stmt *ast.LocalStatement) error {
	c.stats.VariableCount++

	if stmt.Value != nil {
		value, err := c.compileExpression(stmt.Value)
		if err != nil {
			return err
		}
		c.writeLine("Local %s := %s", stmt.Name.Value, value)
	} else {
		c.writeLine("Local %s", stmt.Name.Value)
	}

	return nil
}

// compilePublicStatement compila uma declaração de variável pública
func (c *Compiler) compilePublicStatement(stmt *ast.PublicStatement) error {
	c.stats.VariableCount++

	if stmt.Value != nil {
		value, err := c.compileExpression(stmt.Value)
		if err != nil {
			return err
		}
		c.writeLine("Public %s := %s", stmt.Name.Value, value)
	} else {
		c.writeLine("Public %s", stmt.Name.Value)
	}

	return nil
}

// compilePrivateStatement compila uma declaração de variável privada
func (c *Compiler) compilePrivateStatement(stmt *ast.PrivateStatement) error {
	c.stats.VariableCount++

	if stmt.Value != nil {
		value, err := c.compileExpression(stmt.Value)
		if err != nil {
			return err
		}
		c.writeLine("Private %s := %s", stmt.Name.Value, value)
	} else {
		c.writeLine("Private %s", stmt.Name.Value)
	}

	return nil
}

// compileReturnStatement compila uma declaração de retorno
func (c *Compiler) compileReturnStatement(stmt *ast.ReturnStatement) error {
	if stmt.Value != nil {
		value, err := c.compileExpression(stmt.Value)
		if err != nil {
			return err
		}
		c.writeLine("Return %s", value)
	} else {
		c.writeLine("Return")
	}

	return nil
}

// compileExpressionStatement compila um statement de expressão
func (c *Compiler) compileExpressionStatement(stmt *ast.ExpressionStatement) error {
	expr, err := c.compileExpression(stmt.Expression)
	if err != nil {
		return err
	}
	c.writeLine(expr)
	return nil
}

// compileBlockStatement compila um bloco de código
func (c *Compiler) compileBlockStatement(block *ast.BlockStatement) error {
	for _, stmt := range block.Statements {
		if err := c.compileStatement(stmt); err != nil {
			return err
		}
	}
	return nil
}

// compileExpression compila uma expressão
func (c *Compiler) compileExpression(expr ast.Expression) (string, error) {
	switch e := expr.(type) {
	case *ast.Identifier:
		return e.Value, nil
	case *ast.IntegerLiteral:
		return fmt.Sprintf("%d", e.Value), nil
	case *ast.FloatLiteral:
		return fmt.Sprintf("%g", e.Value), nil
	case *ast.StringLiteral:
		return fmt.Sprintf(`"%s"`, e.Value), nil
	case *ast.DateLiteral:
		return fmt.Sprintf("CTOD('%s')", e.Value), nil
	case *ast.BooleanLiteral:
		if e.Value {
			return ".T.", nil
		}
		return ".F.", nil
	case *ast.NilLiteral:
		return "NIL", nil
	case *ast.PrefixExpression:
		right, err := c.compileExpression(e.Right)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("%s%s", e.Operator, right), nil
	case *ast.InfixExpression:
		left, err := c.compileExpression(e.Left)
		if err != nil {
			return "", err
		}
		right, err := c.compileExpression(e.Right)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("%s %s %s", left, e.Operator, right), nil
	case *ast.CallExpression:
		function, err := c.compileExpression(e.Function)
		if err != nil {
			return "", err
		}

		args := make([]string, len(e.Arguments))
		for i, arg := range e.Arguments {
			compiled, err := c.compileExpression(arg)
			if err != nil {
				return "", err
			}
			args[i] = compiled
		}

		return fmt.Sprintf("%s(%s)", function, strings.Join(args, ", ")), nil
	case *ast.ArrayLiteral:
		elements := make([]string, len(e.Elements))
		for i, elem := range e.Elements {
			compiled, err := c.compileExpression(elem)
			if err != nil {
				return "", err
			}
			elements[i] = compiled
		}
		return fmt.Sprintf("{%s}", strings.Join(elements, ", ")), nil
	case *ast.IndexExpression:
		left, err := c.compileExpression(e.Left)
		if err != nil {
			return "", err
		}
		index, err := c.compileExpression(e.Index)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("%s[%s]", left, index), nil
	case *ast.IfExpression:
		return c.compileIfExpression(e)
	case *ast.WhileExpression:
		return c.compileWhileExpression(e)
	case *ast.ForExpression:
		return c.compileForExpression(e)
	default:
		return "", fmt.Errorf("tipo de expressão não suportado: %T", expr)
	}
}

// compileIfExpression compila uma expressão if
func (c *Compiler) compileIfExpression(expr *ast.IfExpression) (string, error) {
	var result strings.Builder

	condition, err := c.compileExpression(expr.Condition)
	if err != nil {
		return "", err
	}

	result.WriteString(fmt.Sprintf("If %s\n", condition))
	c.indent++
	if err := c.compileBlockStatement(expr.Consequence); err != nil {
		return "", err
	}
	c.indent--

	for _, elseIf := range expr.ElseIfs {
		elseIfCond, err := c.compileExpression(elseIf.Condition)
		if err != nil {
			return "", err
		}

		result.WriteString(fmt.Sprintf("ElseIf %s\n", elseIfCond))
		c.indent++
		if err := c.compileBlockStatement(elseIf.Body); err != nil {
			return "", err
		}
		c.indent--
	}

	if expr.Alternative != nil {
		result.WriteString("Else\n")
		c.indent++
		if err := c.compileBlockStatement(expr.Alternative); err != nil {
			return "", err
		}
		c.indent--
	}

	result.WriteString("EndIf")
	return result.String(), nil
}

// compileWhileExpression compila uma expressão while
func (c *Compiler) compileWhileExpression(expr *ast.WhileExpression) (string, error) {
	var result strings.Builder

	condition, err := c.compileExpression(expr.Condition)
	if err != nil {
		return "", err
	}

	result.WriteString(fmt.Sprintf("While %s\n", condition))
	c.indent++
	if err := c.compileBlockStatement(expr.Body); err != nil {
		return "", err
	}
	c.indent--
	result.WriteString("EndDo")

	return result.String(), nil
}

// compileForExpression compila uma expressão for
func (c *Compiler) compileForExpression(expr *ast.ForExpression) (string, error) {
	var result strings.Builder

	start, err := c.compileExpression(expr.Start)
	if err != nil {
		return "", err
	}

	end, err := c.compileExpression(expr.End)
	if err != nil {
		return "", err
	}

	result.WriteString(fmt.Sprintf("For %s := %s To %s", expr.Counter.Value, start, end))

	if expr.Step != nil {
		step, err := c.compileExpression(expr.Step)
		if err != nil {
			return "", err
		}
		result.WriteString(fmt.Sprintf(" Step %s", step))
	}

	result.WriteString("\n")
	c.indent++
	if err := c.compileBlockStatement(expr.Body); err != nil {
		return "", err
	}
	c.indent--
	result.WriteString("Next")

	return result.String(), nil
}

// writeLine escreve uma linha de código com indentação
func (c *Compiler) writeLine(format string, args ...interface{}) {
	c.stats.LineCount++
	indent := strings.Repeat("\t", c.indent)
	line := fmt.Sprintf(format, args...)
	c.output.WriteString(indent + line + "\n")
}

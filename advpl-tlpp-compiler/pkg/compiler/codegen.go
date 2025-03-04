package compiler

import (
	"fmt"
	"strings"
	"time"

	"github.com/peder1981/advpl-tlpp-compiler/pkg/ast"
)

// CodeGenerator é responsável por gerar código objeto a partir da AST
type CodeGenerator struct {
	program      *ast.Program
	options      Options
	output       strings.Builder
	indent       int
	currentClass string
	sourceFile   string
	timestamp    time.Time
}

// NewCodeGenerator cria um novo gerador de código
func NewCodeGenerator(program *ast.Program, sourceFile string, options Options) *CodeGenerator {
	return &CodeGenerator{
		program:    program,
		options:    options,
		sourceFile: sourceFile,
		timestamp:  time.Now(),
	}
}

// Generate gera o código objeto
func (g *CodeGenerator) Generate() (string, error) {
	// Adicionar cabeçalho
	g.generateHeader()

	// Processar cada statement do programa
	for _, stmt := range g.program.Statements {
		if err := g.generateStatement(stmt); err != nil {
			return "", err
		}
	}

	// Adicionar rodapé
	g.generateFooter()

	return g.output.String(), nil
}

// generateHeader gera o cabeçalho do código objeto
func (g *CodeGenerator) generateHeader() {
	g.writeLine("//----------------------------------------------------------")
	g.writeLine("// Código gerado pelo compilador AdvPL/TLPP")
	g.writeLine("// Arquivo fonte: %s", g.sourceFile)
	g.writeLine("// Data: %s", g.timestamp.Format("2006-01-02 15:04:05"))
	g.writeLine("//----------------------------------------------------------")
	g.writeLine("")

	// Incluir bibliotecas padrão
	g.writeLine("#INCLUDE 'PROTHEUS.CH'")
	g.writeLine("#INCLUDE 'TOTVS.CH'")
	g.writeLine("")
}

// generateFooter gera o rodapé do código objeto
func (g *CodeGenerator) generateFooter() {
	g.writeLine("")
	g.writeLine("//----------------------------------------------------------")
	g.writeLine("// Fim do código gerado")
	g.writeLine("//----------------------------------------------------------")
}

// generateStatement gera código para um statement
func (g *CodeGenerator) generateStatement(stmt ast.Statement) error {
	switch s := stmt.(type) {
	case *ast.FunctionStatement:
		return g.generateFunctionStatement(s)
	case *ast.ClassStatement:
		return g.generateClassStatement(s)
	case *ast.LocalStatement:
		return g.generateLocalStatement(s)
	case *ast.PublicStatement:
		return g.generatePublicStatement(s)
	case *ast.PrivateStatement:
		return g.generatePrivateStatement(s)
	case *ast.ReturnStatement:
		return g.generateReturnStatement(s)
	case *ast.ExpressionStatement:
		return g.generateExpressionStatement(s)
	default:
		return fmt.Errorf("tipo de statement não suportado para geração de código: %T", stmt)
	}
}

// generateFunctionStatement gera código para uma declaração de função
func (g *CodeGenerator) generateFunctionStatement(stmt *ast.FunctionStatement) error {
	// Gerar comentário de documentação
	g.writeLine("/*/{Protheus.doc} %s", stmt.Name.Value)
	g.writeLine("    Função %s", stmt.Name.Value)
	g.writeLine("    @type function")
	g.writeLine("    @author Compilador AdvPL/TLPP")
	g.writeLine("    @since %s", g.timestamp.Format("2006-01-02"))
	
	// Documentar parâmetros
	for _, param := range stmt.Parameters {
		g.writeLine("    @param %s, Variável, Parâmetro da função", param.Value)
	}
	
	g.writeLine("/*/")

	// Gerar cabeçalho da função
	if stmt.Static {
		g.writeLine("Static Function %s(", stmt.Name.Value)
	} else {
		g.writeLine("Function %s(", stmt.Name.Value)
	}

	// Gerar parâmetros
	params := make([]string, len(stmt.Parameters))
	for i, param := range stmt.Parameters {
		params[i] = param.Value
	}
	g.writeLine("%s)", strings.Join(params, ", "))

	// Gerar corpo da função
	g.indent++
	if err := g.generateBlockStatement(stmt.Body); err != nil {
		return err
	}
	g.indent--

	g.writeLine("Return")
	g.writeLine("")

	return nil
}

// generateClassStatement gera código para uma declaração de classe
func (g *CodeGenerator) generateClassStatement(stmt *ast.ClassStatement) error {
	g.currentClass = stmt.Name.Value

	// Gerar comentário de documentação
	g.writeLine("/*/{Protheus.doc} %s", stmt.Name.Value)
	g.writeLine("    Classe %s", stmt.Name.Value)
	if stmt.Parent != nil {
		g.writeLine("    Herda de %s", stmt.Parent.Value)
	}
	g.writeLine("    @type class")
	g.writeLine("    @author Compilador AdvPL/TLPP")
	g.writeLine("    @since %s", g.timestamp.Format("2006-01-02"))
	g.writeLine("/*/")

	// Gerar cabeçalho da classe
	if stmt.Parent != nil {
		g.writeLine("Class %s From %s", stmt.Name.Value, stmt.Parent.Value)
	} else {
		g.writeLine("Class %s", stmt.Name.Value)
	}

	g.indent++

	// Gerar atributos
	for _, data := range stmt.Data {
		g.writeLine("Data %s", data.Name.Value)
	}

	// Gerar métodos
	for _, method := range stmt.Methods {
		if err := g.generateMethodStatement(method); err != nil {
			return err
		}
	}

	g.indent--
	g.writeLine("EndClass")
	g.writeLine("")

	g.currentClass = ""
	return nil
}

// generateMethodStatement gera código para uma declaração de método
func (g *CodeGenerator) generateMethodStatement(stmt *ast.MethodStatement) error {
	// Gerar comentário de documentação
	g.writeLine("/*/{Protheus.doc} %s", stmt.Name.Value)
	g.writeLine("    Método %s da classe %s", stmt.Name.Value, g.currentClass)
	g.writeLine("    @type method")
	g.writeLine("    @author Compilador AdvPL/TLPP")
	g.writeLine("    @since %s", g.timestamp.Format("2006-01-02"))
	
	// Documentar parâmetros
	for _, param := range stmt.Parameters {
		g.writeLine("    @param %s, Variável, Parâmetro do método", param.Value)
	}
	
	g.writeLine("/*/")

	// Gerar cabeçalho do método
	g.writeLine("Method %s(", stmt.Name.Value)

	// Gerar parâmetros
	params := make([]string, len(stmt.Parameters))
	for i, param := range stmt.Parameters {
		params[i] = param.Value
	}
	g.writeLine("%s) Class %s", strings.Join(params, ", "), g.currentClass)

	// Gerar corpo do método
	g.indent++
	if err := g.generateBlockStatement(stmt.Body); err != nil {
		return err
	}
	g.indent--

	g.writeLine("Return Self")
	g.writeLine("")

	return nil
}

// generateLocalStatement gera código para uma declaração de variável local
func (g *CodeGenerator) generateLocalStatement(stmt *ast.LocalStatement) error {
	if stmt.Value != nil {
		value, err := g.generateExpression(stmt.Value)
		if err != nil {
			return err
		}
		g.writeLine("Local %s := %s", stmt.Name.Value, value)
	} else {
		g.writeLine("Local %s", stmt.Name.Value)
	}

	return nil
}

// generatePublicStatement gera código para uma declaração de variável pública
func (g *CodeGenerator) generatePublicStatement(stmt *ast.PublicStatement) error {
	if stmt.Value != nil {
		value, err := g.generateExpression(stmt.Value)
		if err != nil {
			return err
		}
		g.writeLine("Public %s := %s", stmt.Name.Value, value)
	} else {
		g.writeLine("Public %s", stmt.Name.Value)
	}

	return nil
}

// generatePrivateStatement gera código para uma declaração de variável privada
func (g *CodeGenerator) generatePrivateStatement(stmt *ast.PrivateStatement) error {
	if stmt.Value != nil {
		value, err := g.generateExpression(stmt.Value)
		if err != nil {
			return err
		}
		g.writeLine("Private %s := %s", stmt.Name.Value, value)
	} else {
		g.writeLine("Private %s", stmt.Name.Value)
	}

	return nil
}

// generateReturnStatement gera código para uma declaração de retorno
func (g *CodeGenerator) generateReturnStatement(stmt *ast.ReturnStatement) error {
	if stmt.Value != nil {
		value, err := g.generateExpression(stmt.Value)
		if err != nil {
			return err
		}
		g.writeLine("Return %s", value)
	} else {
		g.writeLine("Return")
	}

	return nil
}

// generateExpressionStatement gera código para um statement de expressão
func (g *CodeGenerator) generateExpressionStatement(stmt *ast.ExpressionStatement) error {
	expr, err := g.generateExpression(stmt.Expression)
	if err != nil {
		return err
	}
	g.writeLine(expr)
	return nil
}

// generateBlockStatement gera código para um bloco de código
func (g *CodeGenerator) generateBlockStatement(block *ast.BlockStatement) error {
	for _, stmt := range block.Statements {
		if err := g.generateStatement(stmt); err != nil {
			return err
		}
	}
	return nil
}

// generateExpression gera código para uma expressão
func (g *CodeGenerator) generateExpression(expr ast.Expression) (string, error) {
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
		right, err := g.generateExpression(e.Right)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("%s%s", e.Operator, right), nil
	case *ast.InfixExpression:
		left, err := g.generateExpression(e.Left)
		if err != nil {
			return "", err
		}
		right, err := g.generateExpression(e.Right)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("%s %s %s", left, e.Operator, right), nil
	case *ast.CallExpression:
		function, err := g.generateExpression(e.Function)
		if err != nil {
			return "", err
		}

		args := make([]string, len(e.Arguments))
		for i, arg := range e.Arguments {
			compiled, err := g.generateExpression(arg)
			if err != nil {
				return "", err
			}
			args[i] = compiled
		}

		return fmt.Sprintf("%s(%s)", function, strings.Join(args, ", ")), nil
	case *ast.ArrayLiteral:
		elements := make([]string, len(e.Elements))
		for i, elem := range e.Elements {
			compiled, err := g.generateExpression(elem)
			if err != nil {
				return "", err
			}
			elements[i] = compiled
		}
		return fmt.Sprintf("{%s}", strings.Join(elements, ", ")), nil
	case *ast.IndexExpression:
		left, err := g.generateExpression(e.Left)
		if err != nil {
			return "", err
		}
		index, err := g.generateExpression(e.Index)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("%s[%s]", left, index), nil
	case *ast.IfExpression:
		return g.generateIfExpression(e)
	case *ast.WhileExpression:
		return g.generateWhileExpression(e)
	case *ast.ForExpression:
		return g.generateForExpression(e)
	default:
		return "", fmt.Errorf("tipo de expressão não suportado para geração de código: %T", expr)
	}
}

// generateIfExpression gera código para uma expressão if
func (g *CodeGenerator) generateIfExpression(expr *ast.IfExpression) (string, error) {
	var result strings.Builder

	condition, err := g.generateExpression(expr.Condition)
	if err != nil {
		return "", err
	}

	result.WriteString(fmt.Sprintf("If %s\n", condition))
	g.indent++
	if err := g.generateBlockStatement(expr.Consequence); err != nil {
		return "", err
	}
	g.indent--

	for _, elseIf := range expr.ElseIfs {
		elseIfCond, err := g.generateExpression(elseIf.Condition)
		if err != nil {
			return "", err
		}

		result.WriteString(fmt.Sprintf("ElseIf %s\n", elseIfCond))
		g.indent++
		if err := g.generateBlockStatement(elseIf.Body); err != nil {
			return "", err
		}
		g.indent--
	}

	if expr.Alternative != nil {
		result.WriteString("Else\n")
		g.indent++
		if err := g.generateBlockStatement(expr.Alternative); err != nil {
			return "", err
		}
		g.indent--
	}

	result.WriteString("EndIf")
	return result.String(), nil
}

// generateWhileExpression gera código para uma expressão while
func (g *CodeGenerator) generateWhileExpression(expr *ast.WhileExpression) (string, error) {
	var result strings.Builder

	condition, err := g.generateExpression(expr.Condition)
	if err != nil {
		return "", err
	}

	result.WriteString(fmt.Sprintf("While %s\n", condition))
	g.indent++
	if err := g.generateBlockStatement(expr.Body); err != nil {
		return "", err
	}
	g.indent--
	result.WriteString("EndDo")

	return result.String(), nil
}

// generateForExpression gera código para uma expressão for
func (g *CodeGenerator) generateForExpression(expr *ast.ForExpression) (string, error) {
	var result strings.Builder

	start, err := g.generateExpression(expr.Start)
	if err != nil {
		return "", err
	}

	end, err := g.generateExpression(expr.End)
	if err != nil {
		return "", err
	}

	result.WriteString(fmt.Sprintf("For %s := %s To %s", expr.Counter.Value, start, end))

	if expr.Step != nil {
		step, err := g.generateExpression(expr.Step)
		if err != nil {
			return "", err
		}
		result.WriteString(fmt.Sprintf(" Step %s", step))
	}

	result.WriteString("\n")
	g.indent++
	if err := g.generateBlockStatement(expr.Body); err != nil {
		return "", err
	}
	g.indent--
	result.WriteString("Next")

	return result.String(), nil
}

// writeLine escreve uma linha de código com indentação
func (g *CodeGenerator) writeLine(format string, args ...interface{}) {
	indent := strings.Repeat("\t", g.indent)
	line := fmt.Sprintf(format, args...)
	g.output.WriteString(indent + line + "\n")
}

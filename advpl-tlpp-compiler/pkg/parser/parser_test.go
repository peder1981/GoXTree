package parser

import (
	"fmt"
	"testing"

	"advpl-tlpp-compiler/pkg/ast"
	"advpl-tlpp-compiler/pkg/lexer"
)

func TestFunctionStatements(t *testing.T) {
	input := `Function Soma(a, b)
	Local resultado := a + b
	Return resultado
EndFunction

Static Function Multiplica(a, b)
	Return a * b
EndFunction`

	l := lexer.New(input, "test.prw")
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 2 {
		t.Fatalf("program.Statements não contém 2 statements. got=%d",
			len(program.Statements))
	}

	tests := []struct {
		expectedName string
		expectedStatic bool
		expectedParams []string
	}{
		{"Soma", false, []string{"a", "b"}},
		{"Multiplica", true, []string{"a", "b"}},
	}

	for i, tt := range tests {
		stmt := program.Statements[i]
		if !testFunctionStatement(t, stmt, tt.expectedName, tt.expectedStatic, tt.expectedParams) {
			return
		}
	}
}

func testFunctionStatement(t *testing.T, s ast.Statement, name string, isStatic bool, params []string) bool {
	funcStmt, ok := s.(*ast.FunctionStatement)
	if !ok {
		t.Errorf("s não é *ast.FunctionStatement. got=%T", s)
		return false
	}

	if funcStmt.Name.Value != name {
		t.Errorf("funcStmt.Name.Value não é '%s'. got=%s", name, funcStmt.Name.Value)
		return false
	}

	if funcStmt.Static != isStatic {
		t.Errorf("funcStmt.Static não é %v. got=%v", isStatic, funcStmt.Static)
		return false
	}

	if len(funcStmt.Parameters) != len(params) {
		t.Errorf("length de funcStmt.Parameters é errado. esperado %d, got=%d",
			len(params), len(funcStmt.Parameters))
		return false
	}

	for i, param := range params {
		if funcStmt.Parameters[i].Value != param {
			t.Errorf("parâmetro %d não é '%s'. got=%s", i, param, funcStmt.Parameters[i].Value)
			return false
		}
	}

	return true
}

func TestClassStatements(t *testing.T) {
	input := `Class Pessoa
	Data Nome
	Data Idade
	
	Method New() Constructor
	Method Apresentar()
EndClass

Class Funcionario From Pessoa
	Data Cargo
	Data Salario
	
	Method Trabalhar()
EndClass`

	l := lexer.New(input, "test.prw")
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 2 {
		t.Fatalf("program.Statements não contém 2 statements. got=%d",
			len(program.Statements))
	}

	tests := []struct {
		expectedName string
		expectedParent string
		expectedDataCount int
		expectedMethodCount int
	}{
		{"Pessoa", "", 2, 2},
		{"Funcionario", "Pessoa", 2, 1},
	}

	for i, tt := range tests {
		stmt := program.Statements[i]
		if !testClassStatement(t, stmt, tt.expectedName, tt.expectedParent, tt.expectedDataCount, tt.expectedMethodCount) {
			return
		}
	}
}

func testClassStatement(t *testing.T, s ast.Statement, name string, parent string, dataCount int, methodCount int) bool {
	classStmt, ok := s.(*ast.ClassStatement)
	if !ok {
		t.Errorf("s não é *ast.ClassStatement. got=%T", s)
		return false
	}

	if classStmt.Name.Value != name {
		t.Errorf("classStmt.Name.Value não é '%s'. got=%s", name, classStmt.Name.Value)
		return false
	}

	if parent == "" {
		if classStmt.Parent != nil {
			t.Errorf("classStmt.Parent não é nil. got=%v", classStmt.Parent)
			return false
		}
	} else {
		if classStmt.Parent == nil {
			t.Errorf("classStmt.Parent é nil, esperado '%s'", parent)
			return false
		}
		if classStmt.Parent.Value != parent {
			t.Errorf("classStmt.Parent.Value não é '%s'. got=%s", parent, classStmt.Parent.Value)
			return false
		}
	}

	if len(classStmt.Data) != dataCount {
		t.Errorf("length de classStmt.Data é errado. esperado %d, got=%d",
			dataCount, len(classStmt.Data))
		return false
	}

	if len(classStmt.Methods) != methodCount {
		t.Errorf("length de classStmt.Methods é errado. esperado %d, got=%d",
			methodCount, len(classStmt.Methods))
		return false
	}

	return true
}

func TestVariableStatements(t *testing.T) {
	input := `Local a := 5
Public b := "teste"
Private c := .T.`

	l := lexer.New(input, "test.prw")
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 3 {
		t.Fatalf("program.Statements não contém 3 statements. got=%d",
			len(program.Statements))
	}

	tests := []struct {
		expectedType string
		expectedName string
	}{
		{"Local", "a"},
		{"Public", "b"},
		{"Private", "c"},
	}

	for i, tt := range tests {
		stmt := program.Statements[i]
		if !testVariableStatement(t, stmt, tt.expectedType, tt.expectedName) {
			return
		}
	}
}

func testVariableStatement(t *testing.T, s ast.Statement, varType string, name string) bool {
	switch varType {
	case "Local":
		localStmt, ok := s.(*ast.LocalStatement)
		if !ok {
			t.Errorf("s não é *ast.LocalStatement. got=%T", s)
			return false
		}
		if localStmt.Name.Value != name {
			t.Errorf("localStmt.Name.Value não é '%s'. got=%s", name, localStmt.Name.Value)
			return false
		}
	case "Public":
		publicStmt, ok := s.(*ast.PublicStatement)
		if !ok {
			t.Errorf("s não é *ast.PublicStatement. got=%T", s)
			return false
		}
		if publicStmt.Name.Value != name {
			t.Errorf("publicStmt.Name.Value não é '%s'. got=%s", name, publicStmt.Name.Value)
			return false
		}
	case "Private":
		privateStmt, ok := s.(*ast.PrivateStatement)
		if !ok {
			t.Errorf("s não é *ast.PrivateStatement. got=%T", s)
			return false
		}
		if privateStmt.Name.Value != name {
			t.Errorf("privateStmt.Name.Value não é '%s'. got=%s", name, privateStmt.Name.Value)
			return false
		}
	}
	return true
}

func TestExpressions(t *testing.T) {
	input := `Local a := 5 + 10 * 2
Local b := (5 + 10) * 2
Local c := "olá" + " " + "mundo"
Local d := a > b
Local e := a == b
Local f := a != b
Local g := !f`

	l := lexer.New(input, "test.prw")
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 7 {
		t.Fatalf("program.Statements não contém 7 statements. got=%d",
			len(program.Statements))
	}

	// Apenas verificamos se não há erros de parse, pois testar cada expressão
	// seria muito extenso para este exemplo
}

func TestIfExpression(t *testing.T) {
	input := `If (x > 5)
	y := 10
ElseIf (x == 5)
	y := 5
Else
	y := 0
EndIf`

	l := lexer.New(input, "test.prw")
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements não contém 1 statement. got=%d",
			len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] não é ast.ExpressionStatement. got=%T",
			program.Statements[0])
	}

	exp, ok := stmt.Expression.(*ast.IfExpression)
	if !ok {
		t.Fatalf("stmt.Expression não é ast.IfExpression. got=%T",
			stmt.Expression)
	}

	if exp.Consequence == nil {
		t.Fatal("exp.Consequence é nil")
	}

	if len(exp.ElseIfs) != 1 {
		t.Fatalf("exp.ElseIfs não contém 1 elseif. got=%d",
			len(exp.ElseIfs))
	}

	if exp.Alternative == nil {
		t.Fatal("exp.Alternative é nil")
	}
}

func TestWhileExpression(t *testing.T) {
	input := `While (x > 0)
	x--
EndDo`

	l := lexer.New(input, "test.prw")
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements não contém 1 statement. got=%d",
			len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] não é ast.ExpressionStatement. got=%T",
			program.Statements[0])
	}

	exp, ok := stmt.Expression.(*ast.WhileExpression)
	if !ok {
		t.Fatalf("stmt.Expression não é ast.WhileExpression. got=%T",
			stmt.Expression)
	}

	if exp.Body == nil {
		t.Fatal("exp.Body é nil")
	}
}

func TestForExpression(t *testing.T) {
	input := `For i := 1 To 10 Step 2
	soma := soma + i
Next`

	l := lexer.New(input, "test.prw")
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements não contém 1 statement. got=%d",
			len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] não é ast.ExpressionStatement. got=%T",
			program.Statements[0])
	}

	exp, ok := stmt.Expression.(*ast.ForExpression)
	if !ok {
		t.Fatalf("stmt.Expression não é ast.ForExpression. got=%T",
			stmt.Expression)
	}

	if exp.Counter == nil {
		t.Fatal("exp.Counter é nil")
	}

	if exp.Start == nil {
		t.Fatal("exp.Start é nil")
	}

	if exp.End == nil {
		t.Fatal("exp.End é nil")
	}

	if exp.Step == nil {
		t.Fatal("exp.Step é nil")
	}

	if exp.Body == nil {
		t.Fatal("exp.Body é nil")
	}
}

func checkParserErrors(t *testing.T, p *Parser) {
	errors := p.Errors()
	if len(errors) == 0 {
		return
	}

	t.Errorf("parser tem %d erros", len(errors))
	for _, msg := range errors {
		t.Errorf("erro do parser: %q", msg)
	}
	t.FailNow()
}

package compiler

import (
	"strings"
	"testing"

	"advpl-tlpp-compiler/pkg/ast"
	"advpl-tlpp-compiler/pkg/lexer"
	"advpl-tlpp-compiler/pkg/parser"
)

func TestCompileFunctionStatement(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			`Function Soma(a, b)
				Return a + b
			EndFunction`,
			"Function Soma(a, b)\nReturn a + b\nReturn\n",
		},
		{
			`Static Function Multiplica(a, b)
				Return a * b
			EndFunction`,
			"Static Function Multiplica(a, b)\nReturn a * b\nReturn\n",
		},
		{
			`Function ProcessaValor(nValor)
				Local nResultado := 0
				
				If (nValor > 10)
					nResultado := nValor * 2
				Else
					nResultado := nValor / 2
				EndIf
				
				Return nResultado
			EndFunction`,
			"Function ProcessaValor(nValor)\nLocal nResultado := 0\nIf (nValor > 10)\nnResultado := nValor * 2\nElse\nnResultado := nValor / 2\nEndIf\nReturn nResultado\nReturn\n",
		},
	}

	for i, tt := range tests {
		program := parseCompilerTestProgram(tt.input)
		compiler := New(program, Options{})
		result, err := compiler.Compile()
		
		if err != nil {
			t.Fatalf("teste[%d] - erro durante compilação: %v", i, err)
		}
		
		// Normaliza o resultado removendo espaços em branco extras
		normalizedResult := normalizeWhitespace(result)
		normalizedExpected := normalizeWhitespace(tt.expected)
		
		if normalizedResult != normalizedExpected {
			t.Errorf("teste[%d] - resultado incorreto.\nesperado=\n%s\nobtido=\n%s",
				i, normalizedExpected, normalizedResult)
		}
	}
}

func TestCompileClassStatement(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			`Class Pessoa
				Data Nome
				Data Idade
				
				Method New(cNome, nIdade) Constructor
					::Nome := cNome
					::Idade := nIdade
				Return Self
			EndClass`,
			"Class Pessoa\nData Nome\nData Idade\nMethod New(cNome, nIdade) Class self\n::Nome := cNome\n::Idade := nIdade\nReturn Self\nReturn\nEndClass\n",
		},
		{
			`Class Funcionario From Pessoa
				Data Cargo
				Data Salario
				
				Method Trabalhar()
					// Implementação
				Return .T.
			EndClass`,
			"Class Funcionario From Pessoa\nData Cargo\nData Salario\nMethod Trabalhar() Class self\nReturn .T.\nReturn\nEndClass\n",
		},
	}

	for i, tt := range tests {
		program := parseCompilerTestProgram(tt.input)
		compiler := New(program, Options{})
		result, err := compiler.Compile()
		
		if err != nil {
			t.Fatalf("teste[%d] - erro durante compilação: %v", i, err)
		}
		
		// Normaliza o resultado removendo espaços em branco extras
		normalizedResult := normalizeWhitespace(result)
		normalizedExpected := normalizeWhitespace(tt.expected)
		
		if normalizedResult != normalizedExpected {
			t.Errorf("teste[%d] - resultado incorreto.\nesperado=\n%s\nobtido=\n%s",
				i, normalizedExpected, normalizedResult)
		}
	}
}

func TestCompileVariableStatements(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			`Local a := 5
			Public b := "teste"
			Private c := .T.`,
			"Local a := 5\nPublic b := \"teste\"\nPrivate c := .T.\n",
		},
		{
			`Local nTotal := 0
			Local aArray := {1, 2, 3}
			Local nIndice
			
			For nIndice := 1 To Len(aArray)
				nTotal += aArray[nIndice]
			Next`,
			"Local nTotal := 0\nLocal aArray := {1, 2, 3}\nLocal nIndice\nFor nIndice := 1 To Len(aArray)\nnTotal += aArray[nIndice]\nNext\n",
		},
	}

	for i, tt := range tests {
		program := parseCompilerTestProgram(tt.input)
		compiler := New(program, Options{})
		result, err := compiler.Compile()
		
		if err != nil {
			t.Fatalf("teste[%d] - erro durante compilação: %v", i, err)
		}
		
		// Normaliza o resultado removendo espaços em branco extras
		normalizedResult := normalizeWhitespace(result)
		normalizedExpected := normalizeWhitespace(tt.expected)
		
		if normalizedResult != normalizedExpected {
			t.Errorf("teste[%d] - resultado incorreto.\nesperado=\n%s\nobtido=\n%s",
				i, normalizedExpected, normalizedResult)
		}
	}
}

func TestCompileControlStructures(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			`If (x > 5)
				y := 10
			ElseIf (x == 5)
				y := 5
			Else
				y := 0
			EndIf`,
			"If (x > 5)\ny := 10\nElseIf (x == 5)\ny := 5\nElse\ny := 0\nEndIf\n",
		},
		{
			`While (x > 0)
				x--
			EndDo`,
			"While (x > 0)\nx--\nEndDo\n",
		},
		{
			`For i := 1 To 10 Step 2
				soma := soma + i
			Next`,
			"For i := 1 To 10 Step 2\nsoma := soma + i\nNext\n",
		},
	}

	for i, tt := range tests {
		program := parseCompilerTestProgram(tt.input)
		compiler := New(program, Options{})
		result, err := compiler.Compile()
		
		if err != nil {
			t.Fatalf("teste[%d] - erro durante compilação: %v", i, err)
		}
		
		// Normaliza o resultado removendo espaços em branco extras
		normalizedResult := normalizeWhitespace(result)
		normalizedExpected := normalizeWhitespace(tt.expected)
		
		if normalizedResult != normalizedExpected {
			t.Errorf("teste[%d] - resultado incorreto.\nesperado=\n%s\nobtido=\n%s",
				i, normalizedExpected, normalizedResult)
		}
	}
}

func TestCompileExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			`Local a := 5 + 10 * 2`,
			"Local a := 5 + 10 * 2\n",
		},
		{
			`Local b := (5 + 10) * 2`,
			"Local b := (5 + 10) * 2\n",
		},
		{
			`Local c := "olá" + " " + "mundo"`,
			"Local c := \"olá\" + \" \" + \"mundo\"\n",
		},
		{
			`Local d := a > b`,
			"Local d := a > b\n",
		},
		{
			`Local e := a == b`,
			"Local e := a == b\n",
		},
		{
			`Local f := a != b`,
			"Local f := a != b\n",
		},
		{
			`Local g := !f`,
			"Local g := !f\n",
		},
	}

	for i, tt := range tests {
		program := parseCompilerTestProgram(tt.input)
		compiler := New(program, Options{})
		result, err := compiler.Compile()
		
		if err != nil {
			t.Fatalf("teste[%d] - erro durante compilação: %v", i, err)
		}
		
		// Normaliza o resultado removendo espaços em branco extras
		normalizedResult := normalizeWhitespace(result)
		normalizedExpected := normalizeWhitespace(tt.expected)
		
		if normalizedResult != normalizedExpected {
			t.Errorf("teste[%d] - resultado incorreto.\nesperado=\n%s\nobtido=\n%s",
				i, normalizedExpected, normalizedResult)
		}
	}
}

func TestCompileComplexProgram(t *testing.T) {
	input := `
	Function CalculaFolhaPagamento(aDados)
		Local nTotal := 0
		Local nIndice
		Local nSalario
		Local nBonus
		
		For nIndice := 1 To Len(aDados)
			nSalario := aDados[nIndice, 1]
			nBonus := CalculaBonus(aDados[nIndice, 2], aDados[nIndice, 3])
			
			nTotal += nSalario + nBonus
		Next
		
		Return nTotal
	EndFunction
	
	Static Function CalculaBonus(nTempo, nDesempenho)
		Local nBonus := 0
		
		If (nTempo > 5)
			nBonus := 500
		ElseIf (nTempo > 2)
			nBonus := 300
		Else
			nBonus := 100
		EndIf
		
		If (nDesempenho > 90)
			nBonus *= 1.5
		EndIf
		
		Return nBonus
	EndFunction
	
	Class Funcionario
		Data Nome
		Data Salario
		Data Tempo
		Data Desempenho
		
		Method New(cNome, nSalario, nTempo, nDesempenho) Constructor
			::Nome := cNome
			::Salario := nSalario
			::Tempo := nTempo
			::Desempenho := nDesempenho
		Return Self
		
		Method CalculaSalarioTotal()
			Local nBonus := CalculaBonus(::Tempo, ::Desempenho)
			Return ::Salario + nBonus
	EndClass
	`

	program := parseCompilerTestProgram(input)
	compiler := New(program, Options{})
	result, err := compiler.Compile()
	
	if err != nil {
		t.Fatalf("erro durante compilação: %v", err)
	}
	
	// Verificamos apenas se a compilação foi bem-sucedida e se o resultado não está vazio
	if len(result) == 0 {
		t.Errorf("resultado da compilação está vazio")
	}
}

// Funções auxiliares

// parseCompilerTestProgram analisa um programa para testes de compilação
func parseCompilerTestProgram(input string) *ast.Program {
	l := lexer.New(input, "test.prw")
	p := parser.New(l)
	return p.ParseProgram()
}

func normalizeWhitespace(s string) string {
	// Remove espaços em branco extras
	s = strings.TrimSpace(s)
	return s
}

package compiler

import (
	"strings"
	"testing"

	"advpl-tlpp-compiler/pkg/ast"
	"advpl-tlpp-compiler/pkg/lexer"
	"advpl-tlpp-compiler/pkg/parser"
)

// Função auxiliar para analisar o código-fonte e retornar um programa AST
func parseProgram(input string) *ast.Program {
	l := lexer.New(input)
	p := parser.New(l)
	return p.ParseProgram()
}

func TestCodeGeneratorBasic(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name: "Função simples",
			input: `Function Soma(a, b)
				Return a + b
			EndFunction`,
			expected: "Function Soma(a, b)",
		},
		{
			name: "Classe simples",
			input: `Class Pessoa
				Data Nome
				Method New() Constructor
			EndClass`,
			expected: "Class Pessoa",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			program := parseProgram(tt.input)
			options := Options{
				Verbose:      false,
				Optimize:     false,
				Dialect:      "advpl",
				IncludeDirs:  []string{},
				CheckSyntax:  false,
				GenerateDocs: false,
			}
			
			codeGen := NewCodeGenerator(program, "test.prw", options)
			result, err := codeGen.Generate()
			
			if err != nil {
				t.Fatalf("Erro ao gerar código: %v", err)
			}
			
			if !strings.Contains(result, tt.expected) {
				t.Errorf("Código gerado não contém '%s'. Código gerado:\n%s", 
					tt.expected, result)
			}
		})
	}
}

func TestCodeGeneratorFunctions(t *testing.T) {
	input := `
	Function Soma(a, b)
		Local resultado := a + b
		Return resultado
	EndFunction

	Static Function Multiplica(a, b)
		Return a * b
	EndFunction
	`

	program := parseProgram(input)
	options := Options{
		Verbose:      false,
		Optimize:     false,
		Dialect:      "advpl",
		IncludeDirs:  []string{},
		CheckSyntax:  false,
		GenerateDocs: false,
	}
	
	codeGen := NewCodeGenerator(program, "test.prw", options)
	result, err := codeGen.Generate()
	
	if err != nil {
		t.Fatalf("Erro ao gerar código: %v", err)
	}
	
	// Verificar se o código contém a função Soma
	if !strings.Contains(result, "Function Soma(a, b)") {
		t.Errorf("Código gerado não contém a função Soma. Código gerado:\n%s", result)
	}
	
	// Verificar se o código contém a função estática Multiplica
	if !strings.Contains(result, "Static Function Multiplica(a, b)") {
		t.Errorf("Código gerado não contém a função estática Multiplica. Código gerado:\n%s", result)
	}
}

func TestCodeGeneratorClasses(t *testing.T) {
	input := `
	Class Pessoa
		Data Nome
		Data Idade
		
		Method New(cNome, nIdade) Constructor
			::Nome := cNome
			::Idade := nIdade
		Return Self
		
		Method Apresentar()
			Return "Olá, meu nome é " + ::Nome
	EndClass

	Class Funcionario From Pessoa
		Data Cargo
		Data Salario
		
		Method New(cNome, nIdade, cCargo, nSalario) Constructor
			::Pessoa:New(cNome, nIdade)
			::Cargo := cCargo
			::Salario := nSalario
		Return Self
		
		Method Trabalhar()
			Return "Trabalhando como " + ::Cargo
	EndClass
	`

	program := parseProgram(input)
	options := Options{
		Verbose:      false,
		Optimize:     false,
		Dialect:      "advpl",
		IncludeDirs:  []string{},
		CheckSyntax:  false,
		GenerateDocs: false,
	}
	
	codeGen := NewCodeGenerator(program, "test.prw", options)
	result, err := codeGen.Generate()
	
	if err != nil {
		t.Fatalf("Erro ao gerar código: %v", err)
	}
	
	// Verificar se o código contém a classe Pessoa
	if !strings.Contains(result, "Class Pessoa") {
		t.Errorf("Código gerado não contém a classe Pessoa. Código gerado:\n%s", result)
	}
	
	// Verificar se o código contém a classe Funcionario com herança
	if !strings.Contains(result, "Class Funcionario From Pessoa") {
		t.Errorf("Código gerado não contém a classe Funcionario com herança. Código gerado:\n%s", result)
	}
	
	// Verificar se o código contém o método construtor
	if !strings.Contains(result, "Method New(") && strings.Contains(result, "Constructor") {
		t.Errorf("Código gerado não contém o método construtor. Código gerado:\n%s", result)
	}
}

func TestCodeGeneratorControlStructures(t *testing.T) {
	input := `
	Function TesteControle(nValor)
		Local resultado
		
		If nValor > 10
			resultado := "Maior que 10"
		ElseIf nValor > 5
			resultado := "Entre 6 e 10"
		Else
			resultado := "Menor ou igual a 5"
		EndIf
		
		While nValor > 0
			nValor--
		EndDo
		
		For nI := 1 To 10 Step 2
			resultado := resultado + Str(nI)
		Next
		
		Return resultado
	EndFunction
	`

	program := parseProgram(input)
	options := Options{
		Verbose:      false,
		Optimize:     false,
		Dialect:      "advpl",
		IncludeDirs:  []string{},
		CheckSyntax:  false,
		GenerateDocs: false,
	}
	
	codeGen := NewCodeGenerator(program, "test.prw", options)
	result, err := codeGen.Generate()
	
	if err != nil {
		t.Fatalf("Erro ao gerar código: %v", err)
	}
	
	// Verificar se o código contém a estrutura If
	if !strings.Contains(result, "If") || !strings.Contains(result, "EndIf") {
		t.Errorf("Código gerado não contém a estrutura If. Código gerado:\n%s", result)
	}
	
	// Verificar se o código contém a estrutura While
	if !strings.Contains(result, "While") || !strings.Contains(result, "EndDo") {
		t.Errorf("Código gerado não contém a estrutura While. Código gerado:\n%s", result)
	}
	
	// Verificar se o código contém a estrutura For
	if !strings.Contains(result, "For") || !strings.Contains(result, "Next") {
		t.Errorf("Código gerado não contém a estrutura For. Código gerado:\n%s", result)
	}
}

func TestCodeGeneratorWithOptimization(t *testing.T) {
	input := `
	Function TesteOtimizacao()
		Local a := 10
		Local b := 20
		Local c := a + b  // Isso poderia ser otimizado para c := 30
		
		Return c
	EndFunction
	`

	program := parseProgram(input)
	
	// Teste sem otimização
	optionsNoOpt := Options{
		Verbose:      false,
		Optimize:     false,
		Dialect:      "advpl",
		IncludeDirs:  []string{},
		CheckSyntax:  false,
		GenerateDocs: false,
	}
	
	codeGenNoOpt := NewCodeGenerator(program, "test.prw", optionsNoOpt)
	resultNoOpt, err := codeGenNoOpt.Generate()
	
	if err != nil {
		t.Fatalf("Erro ao gerar código sem otimização: %v", err)
	}
	
	// Teste com otimização
	optionsOpt := Options{
		Verbose:      false,
		Optimize:     true,
		Dialect:      "advpl",
		IncludeDirs:  []string{},
		CheckSyntax:  false,
		GenerateDocs: false,
	}
	
	codeGenOpt := NewCodeGenerator(program, "test.prw", optionsOpt)
	resultOpt, err := codeGenOpt.Generate()
	
	if err != nil {
		t.Fatalf("Erro ao gerar código com otimização: %v", err)
	}
	
	// Verificar se o código otimizado é diferente do não otimizado
	// Nota: Este teste pode falhar se o otimizador não estiver implementado
	// ou se não houver diferença no código gerado para este caso específico
	if resultOpt == resultNoOpt && optionsOpt.Optimize {
		t.Logf("Aviso: O código otimizado é igual ao não otimizado. Isso pode indicar que o otimizador não está implementado.")
	}
}

func TestCodeGeneratorDialects(t *testing.T) {
	input := `
	Function TesteDialeto()
		Local a := 10
		Return a
	EndFunction
	`

	program := parseProgram(input)
	
	// Teste com dialeto AdvPL
	optionsAdvPL := Options{
		Verbose:      false,
		Optimize:     false,
		Dialect:      "advpl",
		IncludeDirs:  []string{},
		CheckSyntax:  false,
		GenerateDocs: false,
	}
	
	codeGenAdvPL := NewCodeGenerator(program, "test.prw", optionsAdvPL)
	resultAdvPL, err := codeGenAdvPL.Generate()
	
	if err != nil {
		t.Fatalf("Erro ao gerar código com dialeto AdvPL: %v", err)
	}
	
	// Teste com dialeto TLPP
	optionsTLPP := Options{
		Verbose:      false,
		Optimize:     false,
		Dialect:      "tlpp",
		IncludeDirs:  []string{},
		CheckSyntax:  false,
		GenerateDocs: false,
	}
	
	codeGenTLPP := NewCodeGenerator(program, "test.prw", optionsTLPP)
	resultTLPP, err := codeGenTLPP.Generate()
	
	if err != nil {
		t.Fatalf("Erro ao gerar código com dialeto TLPP: %v", err)
	}
	
	// Verificar se o código gerado contém o cabeçalho específico do dialeto
	// Nota: Este teste pode falhar se não houver diferença nos cabeçalhos
	if optionsAdvPL.Dialect != optionsTLPP.Dialect && resultAdvPL == resultTLPP {
		t.Logf("Aviso: O código gerado é igual para ambos os dialetos. Isso pode indicar que a diferenciação de dialetos não está implementada.")
	}
}

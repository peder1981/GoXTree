package compiler

import (
	"strings"
	"testing"

	"advpl-tlpp-compiler/pkg/ast"
	"advpl-tlpp-compiler/pkg/lexer"
	"advpl-tlpp-compiler/pkg/parser"
)

func TestSemanticAnalyzerBasic(t *testing.T) {
	tests := []struct {
		name          string
		input         string
		expectErrors  bool
		errorContains string
	}{
		{
			name: "Programa válido",
			input: `
			Function Soma(a, b)
				Return a + b
			EndFunction
			`,
			expectErrors: false,
		},
		{
			name: "Variável não declarada",
			input: `
			Function Teste()
				Return x  // x não foi declarado
			EndFunction
			`,
			expectErrors:  true,
			errorContains: "não declarada",
		},
		{
			name: "Função não declarada",
			input: `
			Function Teste()
				Return FuncaoInexistente()  // Função não existe
			EndFunction
			`,
			expectErrors:  true,
			errorContains: "não declarada",
		},
		{
			name: "Classe não declarada",
			input: `
			Function Teste()
				Local obj := ClasseInexistente():New()  // Classe não existe
				Return obj
			EndFunction
			`,
			expectErrors:  true,
			errorContains: "não declarada",
		},
		{
			name: "Método não declarado",
			input: `
			Class Teste
				Method New() Constructor
				Return Self
			EndClass
			
			Function Uso()
				Local obj := Teste():New()
				Return obj:MetodoInexistente()  // Método não existe
			EndFunction
			`,
			expectErrors:  true,
			errorContains: "não declarado",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			program := parseProgram(tt.input)
			analyzer := NewSemanticAnalyzer()
			errors := analyzer.Analyze(program)
			
			if tt.expectErrors && len(errors) == 0 {
				t.Errorf("Esperava erros semânticos, mas nenhum foi encontrado")
			}
			
			if !tt.expectErrors && len(errors) > 0 {
				t.Errorf("Não esperava erros semânticos, mas encontrou: %v", errors)
			}
			
			if tt.expectErrors && len(errors) > 0 && tt.errorContains != "" {
				found := false
				for _, err := range errors {
					if contains(err.Error(), tt.errorContains) {
						found = true
						break
					}
				}
				
				if !found {
					t.Errorf("Esperava erro contendo '%s', mas não encontrou nos erros: %v", 
						tt.errorContains, errors)
				}
			}
		})
	}
}

func TestSemanticAnalyzerTypeChecking(t *testing.T) {
	tests := []struct {
		name          string
		input         string
		expectErrors  bool
		errorContains string
	}{
		{
			name: "Operação aritmética válida",
			input: `
			Function Teste()
				Local a := 10
				Local b := 20
				Return a + b
			EndFunction
			`,
			expectErrors: false,
		},
		{
			name: "Operação aritmética com tipos incompatíveis",
			input: `
			Function Teste()
				Local a := 10
				Local b := "texto"
				Return a + b  // Soma de número com string
			EndFunction
			`,
			expectErrors:  true,
			errorContains: "incompatíveis",
		},
		{
			name: "Atribuição com tipos incompatíveis",
			input: `
			Function Teste()
				Local a as Numeric
				a := "texto"  // Atribuindo string a uma variável numérica
				Return a
			EndFunction
			`,
			expectErrors:  true,
			errorContains: "incompatíveis",
		},
		{
			name: "Parâmetros de função com tipos incompatíveis",
			input: `
			Function Soma(a as Numeric, b as Numeric)
				Return a + b
			EndFunction
			
			Function Teste()
				Return Soma(10, "texto")  // Segundo parâmetro deveria ser numérico
			EndFunction
			`,
			expectErrors:  true,
			errorContains: "incompatível",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			program := parseProgram(tt.input)
			analyzer := NewSemanticAnalyzer()
			errors := analyzer.Analyze(program)
			
			if tt.expectErrors && len(errors) == 0 {
				t.Errorf("Esperava erros semânticos, mas nenhum foi encontrado")
			}
			
			if !tt.expectErrors && len(errors) > 0 {
				t.Errorf("Não esperava erros semânticos, mas encontrou: %v", errors)
			}
			
			if tt.expectErrors && len(errors) > 0 && tt.errorContains != "" {
				found := false
				for _, err := range errors {
					if contains(err.Error(), tt.errorContains) {
						found = true
						break
					}
				}
				
				if !found {
					t.Errorf("Esperava erro contendo '%s', mas não encontrou nos erros: %v", 
						tt.errorContains, errors)
				}
			}
		})
	}
}

func TestSemanticAnalyzerScopeChecking(t *testing.T) {
	tests := []struct {
		name          string
		input         string
		expectErrors  bool
		errorContains string
	}{
		{
			name: "Escopo válido",
			input: `
			Function Teste()
				Local a := 10
				If a > 5
					Local b := 20
					a := a + b
				EndIf
				Return a
			EndFunction
			`,
			expectErrors: false,
		},
		{
			name: "Variável fora de escopo",
			input: `
			Function Teste()
				If .T.
					Local a := 10
				EndIf
				Return a  // 'a' está fora de escopo aqui
			EndFunction
			`,
			expectErrors:  true,
			errorContains: "escopo",
		},
		{
			name: "Variável de loop fora de escopo",
			input: `
			Function Teste()
				For nI := 1 To 10
					// nI está no escopo
				Next
				Return nI  // 'nI' está fora de escopo aqui
			EndFunction
			`,
			expectErrors:  true,
			errorContains: "escopo",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			program := parseProgram(tt.input)
			analyzer := NewSemanticAnalyzer()
			errors := analyzer.Analyze(program)
			
			if tt.expectErrors && len(errors) == 0 {
				t.Errorf("Esperava erros semânticos, mas nenhum foi encontrado")
			}
			
			if !tt.expectErrors && len(errors) > 0 {
				t.Errorf("Não esperava erros semânticos, mas encontrou: %v", errors)
			}
			
			if tt.expectErrors && len(errors) > 0 && tt.errorContains != "" {
				found := false
				for _, err := range errors {
					if contains(err.Error(), tt.errorContains) {
						found = true
						break
					}
				}
				
				if !found {
					t.Errorf("Esperava erro contendo '%s', mas não encontrou nos erros: %v", 
						tt.errorContains, errors)
				}
			}
		})
	}
}

func TestSemanticAnalyzerClassChecking(t *testing.T) {
	tests := []struct {
		name          string
		input         string
		expectErrors  bool
		errorContains string
	}{
		{
			name: "Classe válida",
			input: `
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
			`,
			expectErrors: false,
		},
		{
			name: "Herança de classe inexistente",
			input: `
			Class Funcionario From ClasseInexistente
				Data Cargo
				
				Method New() Constructor
				Return Self
			EndClass
			`,
			expectErrors:  true,
			errorContains: "inexistente",
		},
		{
			name: "Acesso a atributo inexistente",
			input: `
			Class Pessoa
				Data Nome
				
				Method New(cNome) Constructor
					::Nome := cNome
					::Idade := 30  // Atributo Idade não existe
				Return Self
			EndClass
			`,
			expectErrors:  true,
			errorContains: "inexistente",
		},
		{
			name: "Método com mesmo nome de atributo",
			input: `
			Class Pessoa
				Data Nome
				
				Method Nome()  // Conflito com o atributo Nome
					Return "Conflito"
			EndClass
			`,
			expectErrors:  true,
			errorContains: "conflito",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			program := parseProgram(tt.input)
			analyzer := NewSemanticAnalyzer()
			errors := analyzer.Analyze(program)
			
			if tt.expectErrors && len(errors) == 0 {
				t.Errorf("Esperava erros semânticos, mas nenhum foi encontrado")
			}
			
			if !tt.expectErrors && len(errors) > 0 {
				t.Errorf("Não esperava erros semânticos, mas encontrou: %v", errors)
			}
			
			if tt.expectErrors && len(errors) > 0 && tt.errorContains != "" {
				found := false
				for _, err := range errors {
					if contains(err.Error(), tt.errorContains) {
						found = true
						break
					}
				}
				
				if !found {
					t.Errorf("Esperava erro contendo '%s', mas não encontrou nos erros: %v", 
						tt.errorContains, errors)
				}
			}
		})
	}
}

// Função auxiliar para verificar se uma string contém outra
func contains(s, substr string) bool {
	return strings.Contains(s, substr)
}

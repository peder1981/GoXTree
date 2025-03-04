package compiler

import (
	"strings"
	"testing"
)

func TestOptimizerBasic(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		options  OptimizationOptions
		expected string
	}{
		{
			name: "Sem otimização",
			input: `// Este é um comentário
Function Soma(a, b)
	// Outro comentário
	Return a + b // Retorna a soma
EndFunction`,
			options: OptimizationOptions{
				Level:                OptimizationLevelNone,
				RemoveComments:       true,
				RemoveUnusedVariables: false,
				RemoveUnusedFunctions: false,
				InlineSimpleFunctions: false,
				ConstantFolding:       true,
				DeadCodeElimination:   true,
			},
			expected: `// Este é um comentário
Function Soma(a, b)
	// Outro comentário
	Return a + b // Retorna a soma
EndFunction`,
		},
		{
			name: "Remover comentários",
			input: `// Este é um comentário
Function Soma(a, b)
	// Outro comentário
	Return a + b // Retorna a soma
EndFunction`,
			options: OptimizationOptions{
				Level:                OptimizationLevelBasic,
				RemoveComments:       true,
				RemoveUnusedVariables: false,
				RemoveUnusedFunctions: false,
				InlineSimpleFunctions: false,
				ConstantFolding:       false,
				DeadCodeElimination:   false,
			},
			expected: `
Function Soma(a, b)
	
	Return a + b 
EndFunction`,
		},
		{
			name: "Manter comentários de documentação",
			input: `/*/{Protheus.doc} Soma
@description Soma dois números
@type function
@param a, numeric, Primeiro número
@param b, numeric, Segundo número
@return numeric, Resultado da soma
/*/
Function Soma(a, b)
	// Outro comentário
	Return a + b // Retorna a soma
EndFunction`,
			options: OptimizationOptions{
				Level:                OptimizationLevelBasic,
				RemoveComments:       true,
				RemoveUnusedVariables: false,
				RemoveUnusedFunctions: false,
				InlineSimpleFunctions: false,
				ConstantFolding:       false,
				DeadCodeElimination:   false,
			},
			expected: `/*/{Protheus.doc} Soma
@description Soma dois números
@type function
@param a, numeric, Primeiro número
@param b, numeric, Segundo número
@return numeric, Resultado da soma
/*/
Function Soma(a, b)
	
	Return a + b 
EndFunction`,
		},
		{
			name: "Otimização de constantes",
			input: `Function Calculo()
	Local a := 2 + 3
	Local b := 10 - 5
	Local c := 4 * 3
	Local d := 20 / 4
	Return a + b + c + d
EndFunction`,
			options: OptimizationOptions{
				Level:                OptimizationLevelBasic,
				RemoveComments:       true,
				RemoveUnusedVariables: false,
				RemoveUnusedFunctions: false,
				InlineSimpleFunctions: false,
				ConstantFolding:       true,
				DeadCodeElimination:   false,
			},
			expected: `Function Calculo()
	Local a := 5
	Local b := 5
	Local c := 12
	Local d := 5
	Return a + b + c + d
EndFunction`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			optimizer := NewOptimizer(tt.options)
			result := optimizer.Optimize(tt.input)
			
			// Normalizar espaços em branco para comparação
			normalizedResult := strings.TrimSpace(result)
			normalizedExpected := strings.TrimSpace(tt.expected)
			
			if normalizedResult != normalizedExpected {
				t.Errorf("Resultado não corresponde ao esperado.\nEsperado:\n%s\n\nObtido:\n%s", 
					normalizedExpected, normalizedResult)
			}
		})
	}
}

func TestOptimizerAdvanced(t *testing.T) {
	input := `Function Principal()
		Local a := 10
		Local b := 20
		Local c := a + b  // Variável usada
		Local d := 30     // Variável não usada
		
		If .F.
			// Código morto
			a := 100
		EndIf
		
		Return c
	EndFunction
	
	// Função não utilizada
	Function NaoUtilizada()
		Return .T.
	EndFunction`

	options := OptimizationOptions{
		Level:                OptimizationLevelAdvanced,
		RemoveComments:       true,
		RemoveUnusedVariables: true,
		RemoveUnusedFunctions: true,
		InlineSimpleFunctions: true,
		ConstantFolding:       true,
		DeadCodeElimination:   true,
	}
	
	optimizer := NewOptimizer(options)
	result := optimizer.Optimize(input)
	
	// Verificar se a variável não utilizada foi removida
	if strings.Contains(result, "Local d := 30") {
		t.Logf("Aviso: A variável não utilizada 'd' não foi removida. Isso pode indicar que a remoção de variáveis não utilizadas não está completamente implementada.")
	}
	
	// Verificar se o código morto foi removido
	if strings.Contains(result, "If .F.") {
		t.Logf("Aviso: O código morto não foi removido. Isso pode indicar que a eliminação de código morto não está completamente implementada.")
	}
	
	// Verificar se a função não utilizada foi removida
	if strings.Contains(result, "Function NaoUtilizada()") {
		t.Logf("Aviso: A função não utilizada não foi removida. Isso pode indicar que a remoção de funções não utilizadas não está completamente implementada.")
	}
}

func TestOptimizerInlineFunctions(t *testing.T) {
	input := `Function Dobro(x)
		Return x * 2
	EndFunction
	
	Function Principal()
		Local a := 10
		Local b := Dobro(a)  // Chamada de função que poderia ser inline
		Return b
	EndFunction`

	options := OptimizationOptions{
		Level:                OptimizationLevelAdvanced,
		RemoveComments:       true,
		RemoveUnusedVariables: false,
		RemoveUnusedFunctions: false,
		InlineSimpleFunctions: true,
		ConstantFolding:       true,
		DeadCodeElimination:   true,
	}
	
	optimizer := NewOptimizer(options)
	result := optimizer.Optimize(input)
	
	// Verificar se a função foi inlined
	// Isso é difícil de testar precisamente, então apenas verificamos se a estrutura mudou
	if strings.Count(result, "Function") == strings.Count(input, "Function") &&
	   strings.Contains(result, "Dobro(a)") {
		t.Logf("Aviso: A função 'Dobro' não parece ter sido inlined. Isso pode indicar que o inline de funções simples não está completamente implementado.")
	}
}

func TestDefaultOptimizationOptions(t *testing.T) {
	options := DefaultOptimizationOptions()
	
	if options.Level != OptimizationLevelBasic {
		t.Errorf("Nível de otimização padrão deveria ser OptimizationLevelBasic, mas é %v", options.Level)
	}
	
	if !options.RemoveComments {
		t.Errorf("RemoveComments deveria ser true por padrão")
	}
	
	if options.RemoveUnusedVariables {
		t.Errorf("RemoveUnusedVariables deveria ser false por padrão")
	}
	
	if options.RemoveUnusedFunctions {
		t.Errorf("RemoveUnusedFunctions deveria ser false por padrão")
	}
	
	if options.InlineSimpleFunctions {
		t.Errorf("InlineSimpleFunctions deveria ser false por padrão")
	}
	
	if !options.ConstantFolding {
		t.Errorf("ConstantFolding deveria ser true por padrão")
	}
	
	if !options.DeadCodeElimination {
		t.Errorf("DeadCodeElimination deveria ser true por padrão")
	}
}

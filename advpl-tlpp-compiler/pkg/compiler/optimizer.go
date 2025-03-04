package compiler

import (
	"fmt"
	"regexp"
	"strings"
)

// OptimizationLevel define o nível de otimização
type OptimizationLevel int

const (
	// OptimizationLevelNone sem otimização
	OptimizationLevelNone OptimizationLevel = iota
	// OptimizationLevelBasic otimização básica
	OptimizationLevelBasic
	// OptimizationLevelAdvanced otimização avançada
	OptimizationLevelAdvanced
)

// OptimizationOptions define as opções de otimização
type OptimizationOptions struct {
	Level                OptimizationLevel
	RemoveComments       bool
	RemoveUnusedVariables bool
	RemoveUnusedFunctions bool
	InlineSimpleFunctions bool
	ConstantFolding       bool
	DeadCodeElimination   bool
}

// DefaultOptimizationOptions retorna as opções padrão
func DefaultOptimizationOptions() OptimizationOptions {
	return OptimizationOptions{
		Level:                OptimizationLevelBasic,
		RemoveComments:       true,
		RemoveUnusedVariables: false,
		RemoveUnusedFunctions: false,
		InlineSimpleFunctions: false,
		ConstantFolding:       true,
		DeadCodeElimination:   true,
	}
}

// Optimizer é responsável por otimizar o código gerado
type Optimizer struct {
	options OptimizationOptions
}

// NewOptimizer cria um novo otimizador
func NewOptimizer(options OptimizationOptions) *Optimizer {
	return &Optimizer{
		options: options,
	}
}

// Optimize otimiza o código gerado
func (o *Optimizer) Optimize(code string) string {
	if o.options.Level == OptimizationLevelNone {
		return code
	}

	// Aplicar otimizações básicas
	if o.options.Level >= OptimizationLevelBasic {
		// Remover comentários
		if o.options.RemoveComments {
			code = o.removeComments(code)
		}

		// Remover linhas em branco duplicadas
		code = o.removeExtraBlankLines(code)

		// Otimizar expressões constantes
		if o.options.ConstantFolding {
			code = o.foldConstants(code)
		}
	}

	// Aplicar otimizações avançadas
	if o.options.Level >= OptimizationLevelAdvanced {
		// Eliminar código morto
		if o.options.DeadCodeElimination {
			code = o.eliminateDeadCode(code)
		}

		// Remover variáveis não utilizadas
		if o.options.RemoveUnusedVariables {
			code = o.removeUnusedVariables(code)
		}

		// Remover funções não utilizadas
		if o.options.RemoveUnusedFunctions {
			code = o.removeUnusedFunctions(code)
		}

		// Inline de funções simples
		if o.options.InlineSimpleFunctions {
			code = o.inlineSimpleFunctions(code)
		}
	}

	return code
}

// removeComments remove comentários do código
func (o *Optimizer) removeComments(code string) string {
	// Se a opção de remover comentários não estiver ativada, retornar o código original
	if !o.options.RemoveComments {
		return code
	}

	// Manter comentários de documentação (/*/{Protheus.doc} ... /*/)
	docCommentRegex := regexp.MustCompile(`(?s)/\*/\{Protheus\.doc\}.*?/\*/`)
	docComments := docCommentRegex.FindAllString(code, -1)
	
	// Substituir temporariamente os comentários de documentação
	for i, comment := range docComments {
		placeholder := fmt.Sprintf("__DOC_COMMENT_%d__", i)
		code = strings.Replace(code, comment, placeholder, 1)
	}

	// Remover comentários de linha (//)
	lineCommentRegex := regexp.MustCompile(`//.*`)
	lines := strings.Split(code, "\n")
	for i, line := range lines {
		lines[i] = lineCommentRegex.ReplaceAllString(line, "")
	}
	code = strings.Join(lines, "\n")

	// Remover comentários de bloco (/* ... */)
	blockCommentRegex := regexp.MustCompile(`(?s)/\*.*?\*/`)
	code = blockCommentRegex.ReplaceAllString(code, "")

	// Restaurar comentários de documentação
	for i, comment := range docComments {
		placeholder := fmt.Sprintf("__DOC_COMMENT_%d__", i)
		code = strings.Replace(code, placeholder, comment, 1)
	}

	return code
}

// removeExtraBlankLines remove linhas em branco duplicadas
func (o *Optimizer) removeExtraBlankLines(code string) string {
	// Substituir múltiplas linhas em branco por uma única linha
	blankLinesRegex := regexp.MustCompile(`\n\s*\n\s*\n`)
	for blankLinesRegex.MatchString(code) {
		code = blankLinesRegex.ReplaceAllString(code, "\n\n")
	}
	return code
}

// foldConstants realiza a otimização de constantes
func (o *Optimizer) foldConstants(code string) string {
	// Substituir expressões constantes por seus valores
	// Exemplo: 2 + 3 -> 5
	constExprRegex := regexp.MustCompile(`(\d+)\s*\+\s*(\d+)`)
	code = constExprRegex.ReplaceAllStringFunc(code, func(match string) string {
		parts := constExprRegex.FindStringSubmatch(match)
		if len(parts) == 3 {
			var a, b int
			fmt.Sscanf(parts[1], "%d", &a)
			fmt.Sscanf(parts[2], "%d", &b)
			return fmt.Sprintf("%d", a+b)
		}
		return match
	})

	// Substituir expressões de subtração
	constSubRegex := regexp.MustCompile(`(\d+)\s*\-\s*(\d+)`)
	code = constSubRegex.ReplaceAllStringFunc(code, func(match string) string {
		parts := constSubRegex.FindStringSubmatch(match)
		if len(parts) == 3 {
			var a, b int
			fmt.Sscanf(parts[1], "%d", &a)
			fmt.Sscanf(parts[2], "%d", &b)
			return fmt.Sprintf("%d", a-b)
		}
		return match
	})

	// Substituir expressões de multiplicação
	constMulRegex := regexp.MustCompile(`(\d+)\s*\*\s*(\d+)`)
	code = constMulRegex.ReplaceAllStringFunc(code, func(match string) string {
		parts := constMulRegex.FindStringSubmatch(match)
		if len(parts) == 3 {
			var a, b int
			fmt.Sscanf(parts[1], "%d", &a)
			fmt.Sscanf(parts[2], "%d", &b)
			return fmt.Sprintf("%d", a*b)
		}
		return match
	})

	// Substituir expressões de divisão (apenas se o divisor não for zero)
	constDivRegex := regexp.MustCompile(`(\d+)\s*\/\s*(\d+)`)
	code = constDivRegex.ReplaceAllStringFunc(code, func(match string) string {
		parts := constDivRegex.FindStringSubmatch(match)
		if len(parts) == 3 {
			var a, b int
			fmt.Sscanf(parts[1], "%d", &a)
			fmt.Sscanf(parts[2], "%d", &b)
			if b != 0 {
				return fmt.Sprintf("%d", a/b)
			}
		}
		return match
	})

	// Otimizar expressões booleanas
	code = strings.ReplaceAll(code, ".T. .AND. .T.", ".T.")
	code = strings.ReplaceAll(code, ".F. .AND. .T.", ".F.")
	code = strings.ReplaceAll(code, ".T. .AND. .F.", ".F.")
	code = strings.ReplaceAll(code, ".F. .AND. .F.", ".F.")
	code = strings.ReplaceAll(code, ".T. .OR. .T.", ".T.")
	code = strings.ReplaceAll(code, ".F. .OR. .T.", ".T.")
	code = strings.ReplaceAll(code, ".T. .OR. .F.", ".T.")
	code = strings.ReplaceAll(code, ".F. .OR. .F.", ".F.")
	code = strings.ReplaceAll(code, ".NOT. .T.", ".F.")
	code = strings.ReplaceAll(code, ".NOT. .F.", ".T.")

	return code
}

// eliminateDeadCode elimina código morto
func (o *Optimizer) eliminateDeadCode(code string) string {
	// Eliminar blocos if com condição constante
	ifTrueRegex := regexp.MustCompile(`(?s)If\s+\.T\.\s*\n(.*?)\nElse\s*\n.*?\nEndIf`)
	code = ifTrueRegex.ReplaceAllString(code, "$1")

	ifFalseRegex := regexp.MustCompile(`(?s)If\s+\.F\.\s*\n.*?\nElse\s*\n(.*?)\nEndIf`)
	code = ifFalseRegex.ReplaceAllString(code, "$1")

	// Eliminar loops while com condição falsa
	whileFalseRegex := regexp.MustCompile(`(?s)While\s+\.F\.\s*\n.*?\nEndDo`)
	code = whileFalseRegex.ReplaceAllString(code, "")

	return code
}

// removeUnusedVariables remove variáveis não utilizadas
func (o *Optimizer) removeUnusedVariables(code string) string {
	// Esta é uma implementação simplificada
	// Uma implementação real exigiria análise de fluxo de dados
	
	// Encontrar todas as declarações de variáveis locais
	localVarRegex := regexp.MustCompile(`(?m)^(\s*)Local\s+([a-zA-Z_][a-zA-Z0-9_]*)\s*(?::=.*)?$`)
	declarations := localVarRegex.FindAllStringSubmatch(code, -1)
	
	for _, decl := range declarations {
		indent := decl[1]
		varName := decl[2]
		
		// Verificar se a variável é usada no código
		varUsageRegex := regexp.MustCompile(fmt.Sprintf(`[^a-zA-Z0-9_]%s[^a-zA-Z0-9_]`, regexp.QuoteMeta(varName)))
		usageCount := len(varUsageRegex.FindAllString(code, -1))
		
		// Se a variável é usada apenas uma vez (na declaração), removê-la
		if usageCount <= 1 {
			varDeclRegex := regexp.MustCompile(fmt.Sprintf(`(?m)^%sLocal\s+%s\s*(?::=.*)?$\n`, regexp.QuoteMeta(indent), regexp.QuoteMeta(varName)))
			code = varDeclRegex.ReplaceAllString(code, "")
		}
	}
	
	return code
}

// removeUnusedFunctions remove funções não utilizadas
func (o *Optimizer) removeUnusedFunctions(code string) string {
	// Esta é uma implementação simplificada
	// Uma implementação real exigiria análise de fluxo de dados
	
	// Encontrar todas as declarações de funções
	funcDeclRegex := regexp.MustCompile(`(?s)(Function|Static Function)\s+([a-zA-Z_][a-zA-Z0-9_]*)\s*\(.*?\).*?Return.*?$`)
	declarations := funcDeclRegex.FindAllStringSubmatch(code, -1)
	
	for _, decl := range declarations {
		funcType := decl[1]
		funcName := decl[2]
		
		// Não remover funções de usuário (User Function)
		if strings.Contains(funcType, "User") {
			continue
		}
		
		// Verificar se a função é usada no código
		funcUsageRegex := regexp.MustCompile(fmt.Sprintf(`[^a-zA-Z0-9_]%s\s*\(`, regexp.QuoteMeta(funcName)))
		usageCount := len(funcUsageRegex.FindAllString(code, -1))
		
		// Se a função não é usada, removê-la
		if usageCount == 0 {
			funcDefRegex := regexp.MustCompile(fmt.Sprintf(`(?s)%s\s+%s\s*\(.*?Return.*?$\n`, regexp.QuoteMeta(funcType), regexp.QuoteMeta(funcName)))
			code = funcDefRegex.ReplaceAllString(code, "")
		}
	}
	
	return code
}

// inlineSimpleFunctions realiza o inline de funções simples
func (o *Optimizer) inlineSimpleFunctions(code string) string {
	// Esta é uma implementação simplificada
	// Uma implementação real exigiria análise mais complexa
	
	// Encontrar funções simples (uma única linha de retorno)
	simpleFuncRegex := regexp.MustCompile(`(?s)(Static Function)\s+([a-zA-Z_][a-zA-Z0-9_]*)\s*\((.*?)\)\s*\n\s*Return\s+(.*?)\s*\n`)
	simpleFuncs := simpleFuncRegex.FindAllStringSubmatch(code, -1)
	
	for _, funcMatch := range simpleFuncs {
		funcType := funcMatch[1]
		funcName := funcMatch[2]
		funcParams := funcMatch[3]
		funcReturn := funcMatch[4]
		
		// Verificar se a função é usada no código
		funcUsageRegex := regexp.MustCompile(fmt.Sprintf(`%s\s*\((.*?)\)`, regexp.QuoteMeta(funcName)))
		funcUsages := funcUsageRegex.FindAllStringSubmatch(code, -1)
		
		// Se a função é usada e tem poucos parâmetros, fazer inline
		if len(funcUsages) > 0 && len(strings.Split(funcParams, ",")) <= 2 {
			// Remover a definição da função
			funcDefRegex := regexp.MustCompile(fmt.Sprintf(`(?s)%s\s+%s\s*\(%s\)\s*\n\s*Return\s+%s\s*\n`, 
				regexp.QuoteMeta(funcType), regexp.QuoteMeta(funcName), regexp.QuoteMeta(funcParams), regexp.QuoteMeta(funcReturn)))
			code = funcDefRegex.ReplaceAllString(code, "")
			
			// Substituir chamadas da função pelo seu valor de retorno
			for _ = range funcUsages {
				// Substituição simples para funções sem parâmetros
				if funcParams == "" {
					callRegex := regexp.MustCompile(fmt.Sprintf(`%s\s*\(\s*\)`, regexp.QuoteMeta(funcName)))
					code = callRegex.ReplaceAllString(code, funcReturn)
				} else if !strings.Contains(funcParams, ",") {
					// Para funções com um único parâmetro, fazer substituição direta
					paramName := strings.TrimSpace(funcParams)
					callRegex := regexp.MustCompile(fmt.Sprintf(`%s\s*\(\s*(.*?)\s*\)`, regexp.QuoteMeta(funcName)))
					code = callRegex.ReplaceAllStringFunc(code, func(match string) string {
						parts := callRegex.FindStringSubmatch(match)
						if len(parts) == 2 {
							paramValue := parts[1]
							return strings.ReplaceAll(funcReturn, paramName, paramValue)
						}
						return match
					})
				}
			}
		}
	}
	
	return code
}

package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/peder1981/advpl-tlpp-compiler/pkg/compiler"
	"github.com/peder1981/advpl-tlpp-compiler/pkg/lexer"
	"github.com/peder1981/advpl-tlpp-compiler/pkg/parser"
	"github.com/peder1981/advpl-tlpp-compiler/pkg/utils"
)

var (
	outputFile   string
	verbose      bool
	optimize     bool
	dialect      string
	includeDir   string
	includeDirs  []string
	showVersion  bool
	checkSyntax  bool
	generateDocs bool
)

const version = "0.1.0"

func init() {
	flag.StringVar(&outputFile, "o", "", "Arquivo de saída (padrão: nome do arquivo de entrada com extensão .ppo)")
	flag.BoolVar(&verbose, "v", false, "Modo verboso")
	flag.BoolVar(&optimize, "optimize", false, "Otimizar código gerado")
	flag.StringVar(&dialect, "dialect", "advpl", "Dialeto da linguagem (advpl ou tlpp)")
	flag.StringVar(&includeDir, "I", "", "Diretório de inclusão (pode ser especificado múltiplas vezes)")
	flag.BoolVar(&showVersion, "version", false, "Mostrar versão do compilador")
	flag.BoolVar(&checkSyntax, "check", false, "Apenas verificar sintaxe sem gerar código")
	flag.BoolVar(&generateDocs, "docs", false, "Gerar documentação a partir dos comentários")

	flag.Func("include", "Diretório de inclusão (pode ser especificado múltiplas vezes)", func(s string) error {
		includeDirs = append(includeDirs, s)
		return nil
	})
}

func main() {
	flag.Parse()

	if showVersion {
		fmt.Printf("AdvPL/TLPP Compiler versão %s\n", version)
		return
	}

	// Verificar se foi fornecido um arquivo de entrada
	args := flag.Args()
	if len(args) == 0 {
		fmt.Println("Erro: Nenhum arquivo de entrada especificado")
		fmt.Println("Uso: advpl-compiler [opções] arquivo.prw")
		flag.PrintDefaults()
		os.Exit(1)
	}

	inputFile := args[0]

	// Verificar se o arquivo existe
	if _, err := os.Stat(inputFile); os.IsNotExist(err) {
		fmt.Printf("Erro: Arquivo não encontrado: %s\n", inputFile)
		os.Exit(1)
	}

	// Determinar o arquivo de saída se não foi especificado
	if outputFile == "" {
		ext := filepath.Ext(inputFile)
		outputFile = strings.TrimSuffix(inputFile, ext) + ".ppo"
	}

	// Adicionar o diretório de inclusão se especificado
	if includeDir != "" {
		includeDirs = append(includeDirs, includeDir)
	}

	// Configurar o compilador
	compilerOptions := compiler.Options{
		Verbose:      verbose,
		Optimize:     optimize,
		Dialect:      dialect,
		IncludeDirs:  includeDirs,
		CheckSyntax:  checkSyntax,
		GenerateDocs: generateDocs,
	}

	// Iniciar o processo de compilação
	if verbose {
		fmt.Printf("Compilando %s para %s\n", inputFile, outputFile)
		fmt.Printf("Dialeto: %s\n", dialect)
		if len(includeDirs) > 0 {
			fmt.Printf("Diretórios de inclusão: %s\n", strings.Join(includeDirs, ", "))
		}
	}

	// Ler o arquivo de entrada
	source, err := os.ReadFile(inputFile)
	if err != nil {
		fmt.Printf("Erro ao ler arquivo: %v\n", err)
		os.Exit(1)
	}

	// Criar o lexer
	l := lexer.New(string(source), inputFile)

	// Criar o parser
	p := parser.New(l)

	// Analisar o código
	program, err := p.ParseProgram()
	if err != nil {
		fmt.Printf("Erro de análise: %v\n", err)
		os.Exit(1)
	}

	// Se estamos apenas verificando a sintaxe, terminar aqui
	if checkSyntax {
		if p.Errors() == nil || len(p.Errors()) == 0 {
			fmt.Println("Análise sintática concluída sem erros.")
			return
		}
		fmt.Println("Erros de sintaxe encontrados:")
		for _, err := range p.Errors() {
			fmt.Printf("  - %s\n", err)
		}
		os.Exit(1)
	}

	// Compilar o programa
	c := compiler.New(program, compilerOptions)
	result, err := c.Compile()
	if err != nil {
		fmt.Printf("Erro de compilação: %v\n", err)
		os.Exit(1)
	}

	// Salvar o resultado
	err = os.WriteFile(outputFile, []byte(result), 0644)
	if err != nil {
		fmt.Printf("Erro ao salvar arquivo de saída: %v\n", err)
		os.Exit(1)
	}

	if verbose {
		fmt.Printf("Compilação concluída com sucesso: %s\n", outputFile)
		stats := c.GetStats()
		fmt.Printf("Estatísticas:\n")
		fmt.Printf("  - Funções: %d\n", stats.FunctionCount)
		fmt.Printf("  - Variáveis: %d\n", stats.VariableCount)
		fmt.Printf("  - Linhas de código: %d\n", stats.LineCount)
	}

	// Gerar documentação se solicitado
	if generateDocs {
		docsFile := strings.TrimSuffix(inputFile, filepath.Ext(inputFile)) + ".md"
		docs := utils.GenerateDocs(program)
		err = os.WriteFile(docsFile, []byte(docs), 0644)
		if err != nil {
			fmt.Printf("Erro ao gerar documentação: %v\n", err)
		} else if verbose {
			fmt.Printf("Documentação gerada: %s\n", docsFile)
		}
	}
}

package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"advpl-tlpp-compiler/pkg/compiler"
	"advpl-tlpp-compiler/pkg/executor"
	"advpl-tlpp-compiler/pkg/lexer"
	"advpl-tlpp-compiler/pkg/parser"
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
	runCode      bool
	appServer    string
	environment  string
	keepTemp     bool
)

const version = "0.2.0"

func init() {
	flag.StringVar(&outputFile, "o", "", "Arquivo de saída (padrão: nome do arquivo de entrada com extensão .ppo)")
	flag.BoolVar(&verbose, "v", false, "Modo verboso")
	flag.BoolVar(&optimize, "optimize", false, "Otimizar código gerado")
	flag.StringVar(&dialect, "dialect", "advpl", "Dialeto da linguagem (advpl ou tlpp)")
	flag.StringVar(&includeDir, "I", "", "Diretório de inclusão (pode ser especificado múltiplas vezes)")
	flag.BoolVar(&showVersion, "version", false, "Mostrar versão do compilador")
	flag.BoolVar(&checkSyntax, "check", false, "Apenas verificar sintaxe sem gerar código")
	flag.BoolVar(&generateDocs, "docs", false, "Gerar documentação a partir dos comentários")
	flag.BoolVar(&runCode, "run", false, "Executar o código após compilação")
	flag.StringVar(&appServer, "server", "localhost", "Servidor do Protheus para execução")
	flag.StringVar(&environment, "env", "ENVIRONMENT", "Ambiente do Protheus para execução")
	flag.BoolVar(&keepTemp, "keep-temp", false, "Manter arquivos temporários após execução")

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
	program := p.ParseProgram()
	if program == nil {
		fmt.Printf("Erro de análise: programa nulo\n")
		os.Exit(1)
	}

	// Se estamos apenas verificando a sintaxe, terminar aqui
	if checkSyntax {
		if len(p.Errors()) > 0 {
			fmt.Println("Erros de sintaxe encontrados:")
			for _, err := range p.Errors() {
				fmt.Printf("  - %s\n", err)
			}
			os.Exit(1)
		}
		fmt.Println("Análise sintática concluída sem erros.")
		return
	}

	// Usar o novo gerador de código
	codeGen := compiler.NewCodeGenerator(program, inputFile, compilerOptions)
	result, err := codeGen.Generate()
	if err != nil {
		fmt.Printf("Erro na geração de código: %v\n", err)
		os.Exit(1)
	}

	// Salvar o resultado
	outputPath := filepath.Join(filepath.Dir(inputFile), outputFile)
	err = os.WriteFile(outputPath, []byte(result), 0644)
	if err != nil {
		fmt.Printf("Erro ao salvar arquivo de saída: %v\n", err)
		os.Exit(1)
	}

	if verbose {
		fmt.Printf("Compilação concluída com sucesso: %s\n", outputFile)
		fmt.Printf("Tamanho do código gerado: %d bytes\n", len(result))
	}

	// Gerar documentação se solicitado
	if generateDocs {
		docsFile := strings.TrimSuffix(inputFile, filepath.Ext(inputFile)) + ".md"
		docs := ""
		err = os.WriteFile(docsFile, []byte(docs), 0644)
		if err != nil {
			fmt.Printf("Erro ao gerar documentação: %v\n", err)
		} else if verbose {
			fmt.Printf("Documentação gerada: %s\n", docsFile)
		}
	}

	// Executar o código se solicitado
	if runCode {
		fmt.Println("Executando o código compilado...")
		
		// Configurar opções de execução
		execOptions := executor.ExecutionOptions{
			TempDir:       os.TempDir(),
			KeepTempFiles: keepTemp,
			Verbose:       verbose,
			Timeout:       30 * time.Second,
			Environment:   make(map[string]string),
		}
		
		// Criar executor
		exec := executor.New(execOptions)
		
		// Tentar executar com o Protheus real
		result, err := exec.ExecuteWithProtheus(outputPath, appServer, environment)
		if err != nil {
			fmt.Printf("Não foi possível executar com o Protheus: %v\n", err)
			fmt.Println("Executando em modo de simulação...")
			
			// Executar em modo de simulação
			result, err = exec.ExecuteFile(outputPath)
			if err != nil {
				fmt.Printf("Erro na execução: %v\n", err)
				os.Exit(1)
			}
		}
		
		// Mostrar resultado da execução
		if result.Success {
			fmt.Println("Execução concluída com sucesso!")
			fmt.Printf("Tempo de execução: %v\n", result.Duration)
			fmt.Println("Saída:")
			fmt.Println(result.Output)
		} else {
			fmt.Println("Erro na execução:")
			fmt.Println(result.ErrorMessage)
			if result.Output != "" {
				fmt.Println("Saída:")
				fmt.Println(result.Output)
			}
			os.Exit(1)
		}
	}
}

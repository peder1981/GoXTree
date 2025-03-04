package executor

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/peder1981/advpl-tlpp-compiler/pkg/ast"
	"github.com/peder1981/advpl-tlpp-compiler/pkg/compiler"
	"github.com/peder1981/advpl-tlpp-compiler/pkg/lexer"
	"github.com/peder1981/advpl-tlpp-compiler/pkg/parser"
)

// ExecutionOptions representa as opções para execução
type ExecutionOptions struct {
	TempDir       string
	KeepTempFiles bool
	Verbose       bool
	Timeout       time.Duration
	Environment   map[string]string
}

// DefaultExecutionOptions retorna as opções padrão
func DefaultExecutionOptions() ExecutionOptions {
	return ExecutionOptions{
		TempDir:       os.TempDir(),
		KeepTempFiles: false,
		Verbose:       false,
		Timeout:       30 * time.Second,
		Environment:   make(map[string]string),
	}
}

// ExecutionResult representa o resultado da execução
type ExecutionResult struct {
	Success      bool
	Output       string
	ErrorMessage string
	Duration     time.Duration
}

// Executor é responsável por executar código AdvPL/TLPP
type Executor struct {
	options ExecutionOptions
}

// New cria um novo executor
func New(options ExecutionOptions) *Executor {
	return &Executor{
		options: options,
	}
}

// ExecuteString compila e executa um código fonte AdvPL/TLPP
func (e *Executor) ExecuteString(source string) (*ExecutionResult, error) {
	// Criar lexer
	l := lexer.New(source)

	// Criar parser
	p := parser.New(l)

	// Parsear o programa
	program := p.ParseProgram()
	if len(p.Errors()) > 0 {
		return &ExecutionResult{
			Success:      false,
			ErrorMessage: fmt.Sprintf("Erros de parsing: %v", p.Errors()),
		}, nil
	}

	// Compilar o programa
	return e.executeProgram(program, "source.prw")
}

// ExecuteFile compila e executa um arquivo AdvPL/TLPP
func (e *Executor) ExecuteFile(filePath string) (*ExecutionResult, error) {
	// Ler o arquivo
	source, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("erro ao ler o arquivo: %v", err)
	}

	// Criar lexer
	l := lexer.New(string(source))

	// Criar parser
	p := parser.New(l)

	// Parsear o programa
	program := p.ParseProgram()
	if len(p.Errors()) > 0 {
		return &ExecutionResult{
			Success:      false,
			ErrorMessage: fmt.Sprintf("Erros de parsing: %v", p.Errors()),
		}, nil
	}

	// Compilar o programa
	return e.executeProgram(program, filepath.Base(filePath))
}

// executeProgram compila e executa um programa AST
func (e *Executor) executeProgram(program *ast.Program, sourceFileName string) (*ExecutionResult, error) {
	startTime := time.Now()

	// Criar opções de compilação
	options := compiler.Options{
		Verbose:      e.options.Verbose,
		Optimize:     true,
		Dialect:      "advpl",
		CheckSyntax:  true,
		GenerateDocs: true,
	}

	// Criar gerador de código
	codeGen := compiler.NewCodeGenerator(program, sourceFileName, options)

	// Gerar código
	code, err := codeGen.Generate()
	if err != nil {
		return nil, fmt.Errorf("erro ao gerar código: %v", err)
	}

	// Criar arquivo temporário
	tempDir := filepath.Join(e.options.TempDir, fmt.Sprintf("advpl_exec_%d", time.Now().UnixNano()))
	if err := os.MkdirAll(tempDir, 0755); err != nil {
		return nil, fmt.Errorf("erro ao criar diretório temporário: %v", err)
	}

	// Limpar arquivos temporários se necessário
	if !e.options.KeepTempFiles {
		defer os.RemoveAll(tempDir)
	}

	// Salvar código em arquivo temporário
	outputFile := filepath.Join(tempDir, sourceFileName)
	if err := ioutil.WriteFile(outputFile, []byte(code), 0644); err != nil {
		return nil, fmt.Errorf("erro ao escrever arquivo temporário: %v", err)
	}

	// Executar o código (simulado, já que não temos um runtime real)
	result, err := e.simulateExecution(outputFile)
	if err != nil {
		return nil, err
	}

	result.Duration = time.Since(startTime)
	return result, nil
}

// simulateExecution simula a execução do código AdvPL/TLPP
// Em um ambiente real, isso chamaria o runtime do Protheus
func (e *Executor) simulateExecution(filePath string) (*ExecutionResult, error) {
	// Em um ambiente real, isso executaria o código no Protheus
	// Por enquanto, apenas verificamos a sintaxe
	
	// Ler o arquivo
	code, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("erro ao ler o arquivo compilado: %v", err)
	}

	// Verificar se o código contém erros óbvios
	if strings.Contains(string(code), "SYNTAX ERROR") {
		return &ExecutionResult{
			Success:      false,
			ErrorMessage: "Erro de sintaxe detectado no código gerado",
			Output:       string(code),
		}, nil
	}

	// Simular execução bem-sucedida
	return &ExecutionResult{
		Success: true,
		Output:  fmt.Sprintf("Simulação de execução bem-sucedida para %s", filePath),
	}, nil
}

// ExecuteWithProtheus executa o código usando o Protheus real (se disponível)
func (e *Executor) ExecuteWithProtheus(filePath string, appServer, environment string) (*ExecutionResult, error) {
	// Verificar se o Protheus está disponível
	if _, err := exec.LookPath("appserver"); err != nil {
		return nil, fmt.Errorf("appserver do Protheus não encontrado no PATH")
	}

	startTime := time.Now()

	// Preparar comando para executar o Protheus
	cmd := exec.Command("appserver", 
		"-run", filePath,
		"-server", appServer,
		"-env", environment)

	// Configurar ambiente
	cmd.Env = os.Environ()
	for k, v := range e.options.Environment {
		cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", k, v))
	}

	// Capturar saída
	output, err := cmd.CombinedOutput()
	if err != nil {
		return &ExecutionResult{
			Success:      false,
			Output:       string(output),
			ErrorMessage: err.Error(),
			Duration:     time.Since(startTime),
		}, nil
	}

	return &ExecutionResult{
		Success:  true,
		Output:   string(output),
		Duration: time.Since(startTime),
	}, nil
}

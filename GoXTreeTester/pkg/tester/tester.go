package tester

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/peder1981/GoXTreeTester/pkg/reporter"
)

// Configurações globais para o testador
var (
	testTimeout = "30s"
)

// SetTestTimeout define o tempo limite para execução de testes
func SetTestTimeout(timeout string) {
	testTimeout = timeout
}

// TestResult representa o resultado de um teste
type TestResult struct {
	Name         string
	Passed       bool
	Duration     time.Duration
	ErrorMessage string
}

// Tester é responsável por executar testes
type Tester struct {
	projectPath string
	reporter    *reporter.Reporter
	timeout     string
}

// NewTester cria um novo testador
func NewTester(projectPath string, reporter *reporter.Reporter) *Tester {
	return &Tester{
		projectPath: projectPath,
		reporter:    reporter,
		timeout:     testTimeout,
	}
}

// SetTimeout define o tempo limite para execução de testes para este testador
func (t *Tester) SetTimeout(timeout string) {
	t.timeout = timeout
}

// RunTests executa os testes do projeto
func (t *Tester) RunTests() ([]TestResult, error) {
	var results []TestResult

	// Executar testes Go padrão
	goTestResults, err := t.runGoTests()
	if err != nil {
		return nil, err
	}
	results = append(results, goTestResults...)

	// Executar testes funcionais
	funcTestResults, err := t.runFunctionalTests()
	if err != nil {
		return nil, err
	}
	results = append(results, funcTestResults...)

	// Adicionar resultados ao relatório
	for _, result := range results {
		status := "PASSED"
		if !result.Passed {
			status = "FAILED"
		}
		t.reporter.AddTestResult(result.Name, status, result.Duration.String(), result.ErrorMessage)
	}

	return results, nil
}

// runGoTests executa os testes Go padrão
func (t *Tester) runGoTests() ([]TestResult, error) {
	var results []TestResult

	// Verificar se existem testes
	hasTests := false
	err := filepath.Walk(t.projectPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasSuffix(path, "_test.go") {
			hasTests = true
			return filepath.SkipAll
		}
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("erro ao procurar arquivos de teste: %w", err)
	}

	if !hasTests {
		// Criar testes básicos se não existirem
		err = t.generateBasicTests()
		if err != nil {
			return nil, fmt.Errorf("erro ao gerar testes básicos: %w", err)
		}
	}

	// Executar testes Go com timeout configurado
	cmd := exec.Command("go", "test", "./...", "-v", "-timeout", t.timeout)
	cmd.Dir = t.projectPath
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	startTime := time.Now()
	err = cmd.Run()
	duration := time.Since(startTime)

	output := stdout.String()
	if output == "" {
		output = stderr.String()
	}

	// Analisar saída dos testes
	lines := strings.Split(output, "\n")
	for _, line := range lines {
		if strings.Contains(line, "--- PASS") || strings.Contains(line, "--- FAIL") {
			parts := strings.Split(line, ": ")
			if len(parts) < 2 {
				continue
			}

			testInfo := strings.Split(parts[0], " ")
			if len(testInfo) < 3 {
				continue
			}

			testName := testInfo[2]
			passed := strings.Contains(line, "--- PASS")
			testDuration, _ := time.ParseDuration(strings.TrimSpace(parts[1]))

			var errorMsg string
			if !passed && len(parts) > 2 {
				errorMsg = strings.Join(parts[2:], ": ")
			}

			results = append(results, TestResult{
				Name:         testName,
				Passed:       passed,
				Duration:     testDuration,
				ErrorMessage: errorMsg,
			})
		}
	}

	// Se não conseguiu analisar a saída, adicionar um resultado geral
	if len(results) == 0 {
		passed := err == nil
		errorMsg := ""
		if !passed {
			errorMsg = stderr.String()
		}

		results = append(results, TestResult{
			Name:         "GoTests",
			Passed:       passed,
			Duration:     duration,
			ErrorMessage: errorMsg,
		})
	}

	return results, nil
}

// runFunctionalTests executa testes funcionais personalizados
func (t *Tester) runFunctionalTests() ([]TestResult, error) {
	var results []TestResult

	// Teste de compilação
	startTime := time.Now()
	cmd := exec.Command("go", "build", "./cmd/gxtree")
	cmd.Dir = t.projectPath
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	err := cmd.Run()
	duration := time.Since(startTime)

	results = append(results, TestResult{
		Name:         "CompilationTest",
		Passed:       err == nil,
		Duration:     duration,
		ErrorMessage: stderr.String(),
	})

	// Teste de execução
	if err == nil {
		startTime = time.Now()
		execPath := filepath.Join(t.projectPath, "gxtree")
		if runtime := os.Getenv("GOOS"); runtime == "windows" {
			execPath += ".exe"
		}

		// Verificar se o executável existe
		if _, err := os.Stat(execPath); err == nil {
			// Executar com timeout
			cmd = exec.Command(execPath, "--version")
			var stdout, stderr bytes.Buffer
			cmd.Stdout = &stdout
			cmd.Stderr = &stderr

			// Criar um canal para sinalizar conclusão
			done := make(chan error, 1)
			go func() {
				done <- cmd.Run()
			}()

			// Aguardar conclusão ou timeout
			var execErr error
			select {
			case execErr = <-done:
				// Comando concluído
			case <-time.After(5 * time.Second):
				// Timeout
				if cmd.Process != nil {
					cmd.Process.Kill()
				}
				execErr = fmt.Errorf("timeout após 5 segundos")
			}

			duration = time.Since(startTime)
			errorMsg := ""
			if execErr != nil {
				errorMsg = execErr.Error()
				if stderr.Len() > 0 {
					errorMsg += "\n" + stderr.String()
				}
			}

			results = append(results, TestResult{
				Name:         "ExecutionTest",
				Passed:       execErr == nil,
				Duration:     duration,
				ErrorMessage: errorMsg,
			})
		}
	}

	// Teste de interface
	results = append(results, t.testUIComponents()...)

	return results, nil
}

// testUIComponents testa componentes da interface do usuário
func (t *Tester) testUIComponents() []TestResult {
	var results []TestResult

	// Verificar arquivos de componentes UI
	uiFiles := []string{
		"app.go",
		"app_core.go",
		"file_view.go",
		"menu_bar.go",
		"status_bar.go",
		"help_view.go",
	}

	for _, file := range uiFiles {
		filePath := filepath.Join(t.projectPath, "pkg", "ui", file)
		startTime := time.Now()
		_, err := os.Stat(filePath)
		duration := time.Since(startTime)

		results = append(results, TestResult{
			Name:         "UIComponentTest_" + file,
			Passed:       err == nil,
			Duration:     duration,
			ErrorMessage: err.Error(),
		})
	}

	return results
}

// generateBasicTests gera testes básicos para o projeto
func (t *Tester) generateBasicTests() error {
	// Criar diretório de testes
	testsDir := filepath.Join(t.projectPath, "pkg", "ui", "tests")
	if err := os.MkdirAll(testsDir, 0755); err != nil {
		return err
	}

	// Criar arquivo de teste básico
	testFilePath := filepath.Join(testsDir, "app_test.go")
	testContent := `package ui_test

import (
	"testing"
	
	"github.com/peder1981/GoXTree/pkg/ui"
)

func TestNewApp(t *testing.T) {
	app := ui.NewApp()
	if app == nil {
		t.Error("NewApp() returned nil")
	}
}

func TestAppInitialization(t *testing.T) {
	app := ui.NewApp()
	if app == nil {
		t.Skip("NewApp() returned nil, skipping initialization test")
	}
	
	// Verificar se os componentes básicos foram inicializados
	// Nota: Este teste pode precisar ser adaptado com base na implementação real
}
`

	return os.WriteFile(testFilePath, []byte(testContent), 0644)
}

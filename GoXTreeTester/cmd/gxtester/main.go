package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/fatih/color"
	"github.com/peder1981/GoXTreeTester/pkg/analyzer"
	"github.com/peder1981/GoXTreeTester/pkg/fixer"
	"github.com/peder1981/GoXTreeTester/pkg/reporter"
	"github.com/peder1981/GoXTreeTester/pkg/tester"
)

var (
	projectPath  string
	autoFix      bool
	reportPath   string
	testOnly     bool
	analyzeOnly  bool
	verbose      bool
	ignoreErrors []string
)

func init() {
	flag.StringVar(&projectPath, "project", "", "Caminho para o projeto GoXTree (obrigatório)")
	flag.BoolVar(&autoFix, "autofix", false, "Corrigir problemas automaticamente")
	flag.StringVar(&reportPath, "report", "", "Caminho para salvar o relatório (padrão: ./gxtree_test_report.html)")
	flag.BoolVar(&testOnly, "test-only", false, "Executar apenas os testes")
	flag.BoolVar(&analyzeOnly, "analyze-only", false, "Executar apenas a análise")
	flag.BoolVar(&verbose, "verbose", false, "Exibir informações detalhadas")
	flag.Func("ignore", "Erros a serem ignorados (pode ser usado múltiplas vezes)", func(s string) error {
		ignoreErrors = append(ignoreErrors, s)
		return nil
	})
}

func main() {
	flag.Parse()

	if projectPath == "" {
		fmt.Println("Erro: O caminho do projeto é obrigatório")
		flag.Usage()
		os.Exit(1)
	}

	// Verificar se o caminho existe
	absPath, err := filepath.Abs(projectPath)
	if err != nil {
		fmt.Printf("Erro ao resolver o caminho do projeto: %v\n", err)
		os.Exit(1)
	}

	if _, err := os.Stat(absPath); os.IsNotExist(err) {
		fmt.Printf("Erro: O diretório do projeto não existe: %s\n", absPath)
		os.Exit(1)
	}

	// Definir caminho padrão para o relatório
	if reportPath == "" {
		reportPath = "gxtree_test_report.html"
	}

	// Inicializar componentes
	rep := reporter.NewReporter(reportPath, verbose)
	anl := analyzer.NewAnalyzer(absPath, ignoreErrors, rep)
	tst := tester.NewTester(absPath, rep)
	fix := fixer.NewFixer(absPath, rep)

	// Banner
	printBanner()

	// Executar análise
	if !testOnly {
		color.Green("Iniciando análise do código...")
		issues, err := anl.Analyze()
		if err != nil {
			fmt.Printf("Erro durante a análise: %v\n", err)
			os.Exit(1)
		}

		if len(issues) > 0 {
			color.Yellow("Encontrados %d problemas durante a análise", len(issues))
			for _, issue := range issues {
				fmt.Printf("- %s: %s em %s:%d\n", issue.Severity, issue.Message, issue.File, issue.Line)
			}

			// Corrigir problemas automaticamente se solicitado
			if autoFix {
				color.Green("Corrigindo problemas automaticamente...")
				fixed, err := fix.FixIssues(issues)
				if err != nil {
					fmt.Printf("Erro durante a correção: %v\n", err)
				} else {
					color.Green("Corrigidos %d de %d problemas", fixed, len(issues))
				}
			}
		} else {
			color.Green("Nenhum problema encontrado durante a análise!")
		}
	}

	// Executar testes
	if !analyzeOnly {
		color.Green("Iniciando testes...")
		results, err := tst.RunTests()
		if err != nil {
			fmt.Printf("Erro durante os testes: %v\n", err)
			os.Exit(1)
		}

		// Exibir resultados dos testes
		passedTests := 0
		for _, result := range results {
			if result.Passed {
				passedTests++
			}
		}

		if passedTests == len(results) {
			color.Green("Todos os %d testes passaram!", len(results))
		} else {
			color.Yellow("%d de %d testes passaram", passedTests, len(results))
			for _, result := range results {
				if !result.Passed {
					fmt.Printf("- Falha em %s: %s\n", result.Name, result.ErrorMessage)
				}
			}
		}
	}

	// Gerar relatório
	color.Green("Gerando relatório em %s...", reportPath)
	if err := rep.GenerateReport(); err != nil {
		fmt.Printf("Erro ao gerar relatório: %v\n", err)
		os.Exit(1)
	}

	color.Green("Concluído!")
}

func printBanner() {
	banner := `
  _____      __   ______               ______          __           
 / ___/___  / /_ /_  __/_______  ___  /_  __/__  _____/ /____  _____
/ (_ / _ \/ __/  / / / ___/ _ \/ _ \  / / / _ \/ ___/ __/ _ \/ ___/
\___/\___/\__/  /_/ /_/   \___/_//_/ /_/  \___/_/   \__/\___/_/    
                                                                   
`
	color.Cyan(banner)
	color.Cyan("Versão 1.0.0 - Testador Automático para GoXTree")
	fmt.Println()
}

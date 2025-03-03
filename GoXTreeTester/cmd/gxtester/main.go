package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/fatih/color"
	"github.com/peder1981/GoXTreeTester/pkg/analyzer"
	"github.com/peder1981/GoXTreeTester/pkg/fixer"
	"github.com/peder1981/GoXTreeTester/pkg/reporter"
	"github.com/peder1981/GoXTreeTester/pkg/tester"
)

// Config representa a configuração do GoXTreeTester
type Config struct {
	Project struct {
		Name          string   `json:"name"`
		Path          string   `json:"path"`
		IgnorePatterns []string `json:"ignore_patterns"`
	} `json:"project"`
	Style struct {
		MaxLineLength        int  `json:"max_line_length"`
		CheckTrailingWhitespace bool `json:"check_trailing_whitespace"`
		CheckCommentFormat   bool `json:"check_comment_format"`
		CheckExportedDocs    bool `json:"check_exported_docs"`
	} `json:"style"`
	Imports struct {
		SpecialPackages []string `json:"special_packages"`
	} `json:"imports"`
	Testing struct {
		Timeout     string   `json:"timeout"`
		IgnoreTests []string `json:"ignore_tests"`
	} `json:"testing"`
	Reporting struct {
		DefaultReportPath string `json:"default_report_path"`
		IncludeWarnings   bool   `json:"include_warnings"`
		IncludeFixes      bool   `json:"include_fixes"`
	} `json:"reporting"`
}

var (
	projectPath     string
	autoFix         bool
	fixImportsOnly  bool
	checkStyleOnly  bool
	reportPath      string
	testOnly        bool
	analyzeOnly     bool
	verbose         bool
	ignoreErrors    []string
	configPath      string
	config          Config
)

func init() {
	flag.StringVar(&projectPath, "project", "", "Caminho para o projeto GoXTree (obrigatório se não for especificado no config.json)")
	flag.BoolVar(&autoFix, "autofix", false, "Corrigir problemas automaticamente")
	flag.BoolVar(&fixImportsOnly, "fix-imports", false, "Corrigir apenas imports não utilizados")
	flag.BoolVar(&checkStyleOnly, "check-style", false, "Verificar apenas o estilo de código")
	flag.StringVar(&reportPath, "report", "", "Caminho para salvar o relatório (padrão: ./gxtree_test_report.html)")
	flag.BoolVar(&testOnly, "test-only", false, "Executar apenas os testes")
	flag.BoolVar(&analyzeOnly, "analyze-only", false, "Executar apenas a análise")
	flag.BoolVar(&verbose, "verbose", false, "Exibir informações detalhadas")
	flag.StringVar(&configPath, "config", "config.json", "Caminho para o arquivo de configuração")
	flag.Func("ignore", "Erros a serem ignorados (pode ser usado múltiplas vezes)", func(s string) error {
		ignoreErrors = append(ignoreErrors, s)
		return nil
	})
}

func loadConfig() error {
	// Verificar se o arquivo de configuração existe
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		if configPath != "config.json" {
			return fmt.Errorf("arquivo de configuração não encontrado: %s", configPath)
		}
		// Se o arquivo padrão não existe, apenas retorne sem erro
		return nil
	}

	// Ler o arquivo de configuração
	data, err := ioutil.ReadFile(configPath)
	if err != nil {
		return fmt.Errorf("erro ao ler arquivo de configuração: %v", err)
	}

	// Decodificar o JSON
	if err := json.Unmarshal(data, &config); err != nil {
		return fmt.Errorf("erro ao decodificar arquivo de configuração: %v", err)
	}

	// Aplicar configurações do arquivo se não foram especificadas na linha de comando
	if projectPath == "" && config.Project.Path != "" {
		projectPath = config.Project.Path
	}

	if reportPath == "" && config.Reporting.DefaultReportPath != "" {
		reportPath = config.Reporting.DefaultReportPath
	}

	return nil
}

func main() {
	flag.Parse()

	// Carregar configuração
	if err := loadConfig(); err != nil {
		fmt.Printf("Aviso: %v\n", err)
	}

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
	fix := fixer.NewFixer(absPath, rep)
	
	// Configurar o analisador de estilo com base no arquivo de configuração
	if config.Style.MaxLineLength > 0 {
		analyzer.SetMaxLineLength(config.Style.MaxLineLength)
	}
	
	// Inicializar o testador apenas se for necessário
	var tst *tester.Tester
	if !fixImportsOnly && !analyzeOnly && !checkStyleOnly {
		tst = tester.NewTester(absPath, rep)
		
		// Configurar o testador com base no arquivo de configuração
		if config.Testing.Timeout != "" {
			tester.SetTestTimeout(config.Testing.Timeout)
		}
	}

	// Banner
	printBanner()

	// Verificar apenas o estilo de código
	if checkStyleOnly {
		color.Green("Verificando estilo de código...")
		styleChecker := analyzer.NewStyleChecker(absPath, rep)
		
		// Configurar o verificador de estilo com base no arquivo de configuração
		if config.Style.MaxLineLength > 0 {
			styleChecker.SetMaxLineLength(config.Style.MaxLineLength)
		}
		
		styleIssues, err := styleChecker.CheckStyle()
		if err != nil {
			fmt.Printf("Erro ao verificar estilo de código: %v\n", err)
			os.Exit(1)
		}

		if len(styleIssues) > 0 {
			color.Yellow("Encontrados %d problemas de estilo", len(styleIssues))
			for _, issue := range styleIssues {
				fmt.Printf("- %s: %s em %s:%d\n", issue.Severity, issue.Message, issue.File, issue.Line)
			}

			// Corrigir problemas automaticamente se solicitado
			if autoFix {
				color.Green("Corrigindo problemas de estilo automaticamente...")
				fixed, err := styleChecker.FixStyleIssues(styleIssues)
				if err != nil {
					fmt.Printf("Erro durante a correção: %v\n", err)
				} else {
					color.Green("Corrigidos %d de %d problemas de estilo", fixed, len(styleIssues))
				}
			}
		} else {
			color.Green("Nenhum problema de estilo encontrado!")
		}

		color.Green("Concluído!")
		return
	}

	// Executar análise
	if !testOnly && !fixImportsOnly {
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

	// Corrigir apenas imports não utilizados se solicitado
	if fixImportsOnly {
		color.Green("Corrigindo imports não utilizados...")
		
		// Configurar o fixer com base no arquivo de configuração
		if len(config.Imports.SpecialPackages) > 0 {
			fix.SetSpecialPackages(config.Imports.SpecialPackages)
		}
		
		err := fix.FixUnusedImportsInProject(absPath)
		if err != nil {
			fmt.Printf("Erro durante a correção de imports: %v\n", err)
		} else {
			color.Green("Imports corrigidos com sucesso")
		}
		
		// Se estamos apenas corrigindo imports, não precisamos executar testes ou gerar relatório
		color.Green("Concluído!")
		return
	}

	// Executar testes
	if !analyzeOnly && !fixImportsOnly {
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
	if !fixImportsOnly {
		color.Green("Gerando relatório em %s...", reportPath)
		
		// Configurar o reporter com base no arquivo de configuração
		if config.Reporting.IncludeWarnings {
			rep.SetIncludeWarnings(true)
		}
		
		if config.Reporting.IncludeFixes {
			rep.SetIncludeFixes(true)
		}
		
		if err := rep.GenerateReport(); err != nil {
			fmt.Printf("Erro ao gerar relatório: %v\n", err)
			os.Exit(1)
		}
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

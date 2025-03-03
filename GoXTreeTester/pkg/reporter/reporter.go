package reporter

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// Reporter é responsável por gerar relatórios
type Reporter struct {
	reportPath string
	verbose    bool
	startTime  time.Time
	issues     []Issue
	testResults []TestResult
	fixes      []Fix
}

// Issue representa um problema encontrado
type Issue struct {
	File     string
	Line     int
	Message  string
	Severity string
}

// TestResult representa o resultado de um teste
type TestResult struct {
	Name         string
	Status       string
	Duration     string
	ErrorMessage string
}

// Fix representa uma correção aplicada
type Fix struct {
	File    string
	Message string
}

// NewReporter cria um novo gerador de relatórios
func NewReporter(reportPath string, verbose bool) *Reporter {
	return &Reporter{
		reportPath: reportPath,
		verbose:    verbose,
		startTime:  time.Now(),
	}
}

// AddIssue adiciona um problema ao relatório
func (r *Reporter) AddIssue(file string, line int, message, severity string) {
	r.issues = append(r.issues, Issue{
		File:     file,
		Line:     line,
		Message:  message,
		Severity: severity,
	})

	if r.verbose {
		fmt.Printf("[%s] %s:%d - %s\n", severity, file, line, message)
	}
}

// AddTestResult adiciona um resultado de teste ao relatório
func (r *Reporter) AddTestResult(name, status, duration, errorMessage string) {
	r.testResults = append(r.testResults, TestResult{
		Name:         name,
		Status:       status,
		Duration:     duration,
		ErrorMessage: errorMessage,
	})

	if r.verbose {
		fmt.Printf("[TEST] %s: %s (%s)\n", name, status, duration)
		if errorMessage != "" {
			fmt.Printf("       Error: %s\n", errorMessage)
		}
	}
}

// AddFix adiciona uma correção ao relatório
func (r *Reporter) AddFix(file, message string) {
	r.fixes = append(r.fixes, Fix{
		File:    file,
		Message: message,
	})

	if r.verbose {
		fmt.Printf("[FIX] %s: %s\n", file, message)
	}
}

// GenerateReport gera o relatório final
func (r *Reporter) GenerateReport() error {
	// Criar diretório para o relatório se necessário
	reportDir := filepath.Dir(r.reportPath)
	if reportDir != "" && reportDir != "." {
		if err := os.MkdirAll(reportDir, 0755); err != nil {
			return fmt.Errorf("erro ao criar diretório para relatório: %w", err)
		}
	}

	// Gerar conteúdo do relatório
	content := r.generateHTMLReport()

	// Salvar relatório
	err := os.WriteFile(r.reportPath, []byte(content), 0644)
	if err != nil {
		return fmt.Errorf("erro ao salvar relatório: %w", err)
	}

	return nil
}

// generateHTMLReport gera o relatório em formato HTML
func (r *Reporter) generateHTMLReport() string {
	// Calcular estatísticas
	duration := time.Since(r.startTime)
	issuesByType := make(map[string]int)
	for _, issue := range r.issues {
		issuesByType[issue.Severity]++
	}

	passedTests := 0
	failedTests := 0
	for _, test := range r.testResults {
		if test.Status == "PASSED" {
			passedTests++
		} else {
			failedTests++
		}
	}

	// Construir HTML
	var html strings.Builder

	html.WriteString(`<!DOCTYPE html>
<html lang="pt-br">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Relatório de Testes - GoXTree</title>
    <style>
        body {
            font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
            line-height: 1.6;
            color: #333;
            margin: 0;
            padding: 20px;
            background-color: #f5f5f5;
        }
        .container {
            max-width: 1200px;
            margin: 0 auto;
            background-color: #fff;
            padding: 20px;
            border-radius: 5px;
            box-shadow: 0 2px 5px rgba(0,0,0,0.1);
        }
        h1, h2, h3 {
            color: #2c3e50;
        }
        h1 {
            text-align: center;
            padding-bottom: 10px;
            border-bottom: 2px solid #eee;
        }
        .summary {
            display: flex;
            justify-content: space-around;
            margin: 20px 0;
            padding: 15px;
            background-color: #f8f9fa;
            border-radius: 5px;
        }
        .summary-item {
            text-align: center;
        }
        .summary-item .count {
            font-size: 24px;
            font-weight: bold;
        }
        .error { color: #e74c3c; }
        .warning { color: #f39c12; }
        .info { color: #3498db; }
        .success { color: #2ecc71; }
        table {
            width: 100%;
            border-collapse: collapse;
            margin: 20px 0;
        }
        th, td {
            padding: 12px 15px;
            text-align: left;
            border-bottom: 1px solid #ddd;
        }
        th {
            background-color: #f8f9fa;
            font-weight: bold;
        }
        tr:hover {
            background-color: #f5f5f5;
        }
        .badge {
            display: inline-block;
            padding: 3px 7px;
            border-radius: 3px;
            font-size: 12px;
            font-weight: bold;
            text-transform: uppercase;
        }
        .badge-error { background-color: #e74c3c; color: white; }
        .badge-warning { background-color: #f39c12; color: white; }
        .badge-info { background-color: #3498db; color: white; }
        .badge-success { background-color: #2ecc71; color: white; }
        .footer {
            text-align: center;
            margin-top: 30px;
            padding-top: 20px;
            border-top: 1px solid #eee;
            color: #7f8c8d;
            font-size: 14px;
        }
        .collapsible {
            background-color: #f8f9fa;
            color: #444;
            cursor: pointer;
            padding: 18px;
            width: 100%;
            border: none;
            text-align: left;
            outline: none;
            font-size: 15px;
            transition: 0.4s;
            border-radius: 5px;
            margin: 5px 0;
        }
        .active, .collapsible:hover {
            background-color: #eee;
        }
        .content {
            padding: 0 18px;
            display: none;
            overflow: hidden;
            background-color: #f1f1f1;
            border-radius: 0 0 5px 5px;
        }
    </style>
</head>
<body>
    <div class="container">
        <h1>Relatório de Testes - GoXTree</h1>
        
        <div class="summary">
            <div class="summary-item">
                <div class="count">` + fmt.Sprintf("%d", len(r.issues)) + `</div>
                <div>Problemas Encontrados</div>
            </div>
            <div class="summary-item">
                <div class="count">` + fmt.Sprintf("%d", len(r.fixes)) + `</div>
                <div>Correções Aplicadas</div>
            </div>
            <div class="summary-item">
                <div class="count">` + fmt.Sprintf("%d", passedTests) + `</div>
                <div class="success">Testes Passaram</div>
            </div>
            <div class="summary-item">
                <div class="count">` + fmt.Sprintf("%d", failedTests) + `</div>
                <div class="error">Testes Falharam</div>
            </div>
            <div class="summary-item">
                <div class="count">` + formatDuration(duration) + `</div>
                <div>Tempo Total</div>
            </div>
        </div>
        
        <h2>Problemas Encontrados</h2>`)

	if len(r.issues) > 0 {
		html.WriteString(`
        <table>
            <thead>
                <tr>
                    <th>Arquivo</th>
                    <th>Linha</th>
                    <th>Severidade</th>
                    <th>Mensagem</th>
                </tr>
            </thead>
            <tbody>`)

		for _, issue := range r.issues {
			severityClass := "badge-info"
			if issue.Severity == "ERROR" {
				severityClass = "badge-error"
			} else if issue.Severity == "WARNING" {
				severityClass = "badge-warning"
			}

			html.WriteString(fmt.Sprintf(`
                <tr>
                    <td>%s</td>
                    <td>%d</td>
                    <td><span class="badge %s">%s</span></td>
                    <td>%s</td>
                </tr>`, issue.File, issue.Line, severityClass, issue.Severity, issue.Message))
		}

		html.WriteString(`
            </tbody>
        </table>`)
	} else {
		html.WriteString(`
        <p class="success">Nenhum problema encontrado!</p>`)
	}

	html.WriteString(`
        <h2>Resultados dos Testes</h2>`)

	if len(r.testResults) > 0 {
		html.WriteString(`
        <table>
            <thead>
                <tr>
                    <th>Nome</th>
                    <th>Status</th>
                    <th>Duração</th>
                    <th>Detalhes</th>
                </tr>
            </thead>
            <tbody>`)

		for _, test := range r.testResults {
			statusClass := "badge-success"
			if test.Status != "PASSED" {
				statusClass = "badge-error"
			}

			html.WriteString(fmt.Sprintf(`
                <tr>
                    <td>%s</td>
                    <td><span class="badge %s">%s</span></td>
                    <td>%s</td>
                    <td>`, test.Name, statusClass, test.Status, test.Duration))

			if test.ErrorMessage != "" {
				html.WriteString(`<button class="collapsible">Ver detalhes</button>
                    <div class="content">
                        <pre>` + test.ErrorMessage + `</pre>
                    </div>`)
			} else {
				html.WriteString("-")
			}

			html.WriteString(`</td>
                </tr>`)
		}

		html.WriteString(`
            </tbody>
        </table>`)
	} else {
		html.WriteString(`
        <p>Nenhum teste executado!</p>`)
	}

	html.WriteString(`
        <h2>Correções Aplicadas</h2>`)

	if len(r.fixes) > 0 {
		html.WriteString(`
        <table>
            <thead>
                <tr>
                    <th>Arquivo</th>
                    <th>Descrição</th>
                </tr>
            </thead>
            <tbody>`)

		for _, fix := range r.fixes {
			html.WriteString(fmt.Sprintf(`
                <tr>
                    <td>%s</td>
                    <td>%s</td>
                </tr>`, fix.File, fix.Message))
		}

		html.WriteString(`
            </tbody>
        </table>`)
	} else {
		html.WriteString(`
        <p>Nenhuma correção aplicada!</p>`)
	}

	html.WriteString(`
        <div class="footer">
            <p>Relatório gerado em ` + time.Now().Format("02/01/2006 15:04:05") + `</p>
            <p>GoXTreeTester v1.0.0</p>
        </div>
    </div>

    <script>
        var coll = document.getElementsByClassName("collapsible");
        for (var i = 0; i < coll.length; i++) {
            coll[i].addEventListener("click", function() {
                this.classList.toggle("active");
                var content = this.nextElementSibling;
                if (content.style.display === "block") {
                    content.style.display = "none";
                } else {
                    content.style.display = "block";
                }
            });
        }
    </script>
</body>
</html>`)

	return html.String()
}

// formatDuration formata a duração de forma legível
func formatDuration(d time.Duration) string {
	d = d.Round(time.Second)
	h := d / time.Hour
	d -= h * time.Hour
	m := d / time.Minute
	d -= m * time.Minute
	s := d / time.Second
	
	if h > 0 {
		return fmt.Sprintf("%dh %dm %ds", h, m, s)
	}
	if m > 0 {
		return fmt.Sprintf("%dm %ds", m, s)
	}
	return fmt.Sprintf("%ds", s)
}

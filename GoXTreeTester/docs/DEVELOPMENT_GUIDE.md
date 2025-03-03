# Guia de Desenvolvimento do GoXTreeTester

Este guia é destinado a desenvolvedores que desejam contribuir com o GoXTreeTester. Ele explica a estrutura do código, como adicionar novos recursos e como testar suas alterações.

## Estrutura do Código

O GoXTreeTester é organizado em vários pacotes:

```
GoXTreeTester/
├── cmd/
│   └── gxtester/
│       └── main.go         # Ponto de entrada da aplicação
├── pkg/
│   ├── analyzer/
│   │   ├── analyzer.go     # Análise de código
│   │   └── style_checker.go # Verificação de estilo de código
│   ├── fixer/
│   │   └── fixer.go        # Correção automática
│   ├── tester/
│   │   └── tester.go       # Execução de testes
│   └── reporter/
│       └── reporter.go     # Geração de relatórios
└── docs/
    ├── ARCHITECTURE.md     # Documentação da arquitetura
    ├── USAGE_GUIDE.md      # Guia de uso
    └── DEVELOPMENT_GUIDE.md # Este guia
```

## Fluxo de Trabalho de Desenvolvimento

### 1. Configurar o Ambiente de Desenvolvimento

Antes de começar a desenvolver, você precisa configurar seu ambiente:

```bash
# Clonar o repositório
git clone https://github.com/peder1981/GoXTreeTester.git

# Entrar no diretório
cd GoXTreeTester

# Instalar dependências
go mod download
```

### 2. Entender o Código Existente

Antes de fazer alterações, é importante entender como o código existente funciona. Leia a documentação em `docs/ARCHITECTURE.md` e explore o código-fonte.

### 3. Fazer Alterações

Ao fazer alterações, siga estas diretrizes:

- **Mantenha a estrutura modular**: Cada pacote tem uma responsabilidade específica.
- **Escreva testes**: Adicione testes para qualquer novo recurso ou correção de bug.
- **Siga as convenções de estilo do Go**: Use `gofmt` e `golint` para garantir que seu código siga as convenções de estilo do Go.
- **Documente seu código**: Adicione comentários explicando o que seu código faz.

### 4. Testar Suas Alterações

Antes de enviar suas alterações, teste-as:

```bash
# Executar testes
go test ./...

# Compilar
go build -o gxtester ./cmd/gxtester

# Executar
./gxtester --project=/caminho/para/GoXTree
```

### 5. Enviar Suas Alterações

Quando estiver satisfeito com suas alterações, envie-as:

```bash
# Criar um branch
git checkout -b meu-recurso

# Adicionar suas alterações
git add .

# Fazer commit
git commit -m "Adicionado novo recurso: descrição"

# Enviar para o GitHub
git push origin meu-recurso
```

## Adicionar Novos Recursos

### Adicionar uma Nova Verificação de Estilo

Para adicionar uma nova verificação de estilo, você precisa modificar o arquivo `pkg/analyzer/style_checker.go`:

1. Adicione uma nova função para a verificação:

```go
// checkNewStyleIssue verifica se há um novo problema de estilo
func (sc *StyleChecker) checkNewStyleIssue(file *ast.File, fset *token.FileSet) []Issue {
    var issues []Issue
    
    // Implementar a verificação
    
    return issues
}
```

2. Chame a nova função em `CheckStyle`:

```go
func (sc *StyleChecker) CheckStyle(file *ast.File, fset *token.FileSet, filePath string) []Issue {
    var issues []Issue
    
    // Chamadas existentes
    issues = append(issues, sc.checkLineLengths(file, fset)...)
    issues = append(issues, sc.checkTrailingWhitespace(file, fset)...)
    
    // Nova chamada
    issues = append(issues, sc.checkNewStyleIssue(file, fset)...)
    
    return issues
}
```

### Adicionar uma Nova Correção Automática

Para adicionar uma nova correção automática, você precisa modificar o arquivo `pkg/fixer/fixer.go`:

1. Adicione uma nova função para a correção:

```go
// FixNewIssue corrige um novo tipo de problema
func (f *Fixer) FixNewIssue(filePath string) error {
    // Implementar a correção
    
    return nil
}
```

2. Chame a nova função em `FixIssues`:

```go
func (f *Fixer) FixIssues(issues []analyzer.Issue) []string {
    var fixed []string
    
    // Correções existentes
    
    // Nova correção
    for _, issue := range issues {
        if issue.Type == "Novo Tipo de Problema" {
            err := f.FixNewIssue(issue.FilePath)
            if err == nil {
                fixed = append(fixed, issue.FilePath)
            }
        }
    }
    
    return fixed
}
```

## Componentes Principais

### Analyzer

O Analyzer é responsável por analisar o código-fonte do GoXTree e identificar problemas. Ele usa o pacote `go/ast` para analisar a estrutura do código.

#### Adicionar uma Nova Análise

Para adicionar uma nova análise, você precisa modificar o arquivo `pkg/analyzer/analyzer.go`:

1. Adicione uma nova função para a análise:

```go
// checkNewIssue verifica se há um novo tipo de problema
func (a *Analyzer) checkNewIssue(file *ast.File, fset *token.FileSet, filePath string) []Issue {
    var issues []Issue
    
    // Implementar a análise
    
    return issues
}
```

2. Chame a nova função em `Analyze`:

```go
func (a *Analyzer) Analyze(projectPath string) []Issue {
    var issues []Issue
    
    // Análises existentes
    
    // Nova análise
    for _, filePath := range a.goFiles {
        fileIssues := a.checkNewIssue(file, fset, filePath)
        issues = append(issues, fileIssues...)
    }
    
    return issues
}
```

### StyleChecker

O StyleChecker é responsável por verificar o estilo do código. Ele verifica problemas como linhas muito longas, espaços em branco no final das linhas e falta de documentação.

#### Estrutura do StyleChecker

```go
// StyleChecker é responsável por verificar o estilo do código
type StyleChecker struct {
    MaxLineLength int
}

// NewStyleChecker cria um novo StyleChecker
func NewStyleChecker() *StyleChecker {
    return &StyleChecker{
        MaxLineLength: 100,
    }
}

// CheckStyle verifica o estilo do código
func (sc *StyleChecker) CheckStyle(file *ast.File, fset *token.FileSet, filePath string) []Issue {
    var issues []Issue
    
    issues = append(issues, sc.checkLineLengths(file, fset)...)
    issues = append(issues, sc.checkTrailingWhitespace(file, fset)...)
    issues = append(issues, sc.checkCommentFormat(file, fset)...)
    issues = append(issues, sc.checkExportedDocs(file, fset)...)
    
    return issues
}
```

### Fixer

O Fixer é responsável por corrigir automaticamente os problemas identificados pelo Analyzer. Ele implementa várias estratégias de correção.

#### Estrutura do Fixer

```go
// Fixer é responsável por corrigir problemas no código
type Fixer struct {
    ProjectPath string
}

// NewFixer cria um novo Fixer
func NewFixer(projectPath string) *Fixer {
    return &Fixer{
        ProjectPath: projectPath,
    }
}

// FixIssues corrige os problemas identificados
func (f *Fixer) FixIssues(issues []analyzer.Issue) []string {
    var fixed []string
    
    // Implementar correções
    
    return fixed
}
```

### Reporter

O Reporter é responsável por gerar relatórios sobre a análise e os testes. Ele coleta informações durante a execução e gera um relatório HTML detalhado.

#### Estrutura do Reporter

```go
// Reporter é responsável por gerar relatórios
type Reporter struct {
    Issues      []analyzer.Issue
    TestResults []tester.TestResult
    Fixes       []string
    Warnings    []string
}

// NewReporter cria um novo Reporter
func NewReporter() *Reporter {
    return &Reporter{
        Issues:      []analyzer.Issue{},
        TestResults: []tester.TestResult{},
        Fixes:       []string{},
        Warnings:    []string{},
    }
}

// GenerateReport gera um relatório HTML
func (r *Reporter) GenerateReport(outputPath string) error {
    // Implementar geração de relatório
    
    return nil
}
```

## Boas Práticas de Desenvolvimento

### Testes

Escreva testes para qualquer novo recurso ou correção de bug. O GoXTreeTester usa o pacote `testing` padrão do Go para testes.

Exemplo de teste:

```go
func TestStyleChecker_CheckLineLengths(t *testing.T) {
    // Configurar
    sc := NewStyleChecker()
    
    // Executar
    issues := sc.checkLineLengths(file, fset)
    
    // Verificar
    if len(issues) != 1 {
        t.Errorf("Esperado 1 problema, obtido %d", len(issues))
    }
}
```

### Documentação

Documente seu código com comentários. Siga as convenções de documentação do Go:

- Comentários de pacote: `// Package analyzer ...`
- Comentários de função: `// FunctionName faz algo ...`
- Comentários de tipo: `// TypeName é um tipo que ...`

### Tratamento de Erros

Trate erros adequadamente. Não ignore erros, e forneça mensagens de erro úteis.

Exemplo:

```go
file, err := os.Open(filePath)
if err != nil {
    return fmt.Errorf("falha ao abrir arquivo %s: %w", filePath, err)
}
defer file.Close()
```

### Logging

Use o pacote `log` para logging. Forneça mensagens de log úteis.

Exemplo:

```go
log.Printf("Analisando arquivo: %s", filePath)
```

## Conclusão

Este guia fornece uma visão geral de como contribuir com o GoXTreeTester. Se você tiver dúvidas ou precisar de ajuda, entre em contato com a equipe de desenvolvimento.

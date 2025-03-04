# Manual Técnico - Compilador AdvPL/TLPP

## Índice
1. [Arquitetura](#arquitetura)
2. [Componentes](#componentes)
3. [Fluxo de Compilação](#fluxo-de-compilação)
4. [Desenvolvimento](#desenvolvimento)
5. [Testes](#testes)
6. [Contribuição](#contribuição)

## Arquitetura

O compilador AdvPL/TLPP é estruturado em módulos independentes que seguem o padrão de compilação em fases:

1. **Análise Léxica** (pkg/lexer)
   - Tokenização do código fonte
   - Identificação de palavras-chave
   - Tratamento de comentários e strings

2. **Análise Sintática** (pkg/parser)
   - Construção da AST
   - Validação da estrutura do código
   - Tratamento de expressões e statements

3. **Análise Semântica** (pkg/compiler)
   - Verificação de tipos
   - Validação de escopo
   - Resolução de símbolos

4. **Geração de Código** (pkg/compiler)
   - Otimização
   - Geração do código objeto
   - Formatação do output

5. **Integração IDE** (pkg/lsp)
   - Servidor LSP
   - Diagnósticos
   - Completação de código

## Componentes

### Lexer (pkg/lexer)

O lexer implementa a análise léxica usando um autômato finito:

```go
type Lexer struct {
    input        string
    position     int
    readPosition int
    ch           rune
    line         int
    column       int
    file         string
}
```

Principais funções:
- `NextToken()`: Retorna o próximo token
- `readChar()`: Lê o próximo caractere
- `skipWhitespace()`: Ignora espaços em branco
- `readIdentifier()`: Lê um identificador
- `readNumber()`: Lê um número

### Parser (pkg/parser)

O parser implementa análise sintática descendente recursiva:

```go
type Parser struct {
    lexer        *lexer.Lexer
    currentToken lexer.Token
    peekToken    lexer.Token
    errors       []string
}
```

Principais funções:
- `ParseProgram()`: Ponto de entrada do parser
- `parseStatement()`: Parse de statements
- `parseExpression()`: Parse de expressões
- `parseFunction()`: Parse de funções
- `parseClass()`: Parse de classes

### Compilador (pkg/compiler)

O compilador realiza a análise semântica e geração de código:

```go
type Compiler struct {
    program  *ast.Program
    options  Options
    stats    Stats
    output   strings.Builder
    indent   int
    includes map[string]bool
}
```

Principais funções:
- `Compile()`: Compila o programa
- `compileStatement()`: Compila statements
- `compileExpression()`: Compila expressões
- `optimize()`: Otimiza o código

### Servidor LSP (pkg/lsp)

O servidor LSP implementa o protocolo LSP:

```go
type Server struct {
    documents     sync.Map
    capabilities  ServerCapabilities
    compiler      *compiler.Compiler
    ide           *compiler.IDEIntegration
    configuration Configuration
}
```

Principais funções:
- `Initialize()`: Inicializa o servidor
- `DidOpen()`: Abre um documento
- `DidChange()`: Processa mudanças
- `DocumentSymbol()`: Retorna símbolos
- `Completion()`: Fornece completação

## Fluxo de Compilação

1. **Entrada**
   - Leitura do arquivo fonte
   - Configuração de opções
   - Resolução de includes

2. **Análise Léxica**
   ```go
   lexer := lexer.New(input, filename)
   token := lexer.NextToken()
   ```

3. **Análise Sintática**
   ```go
   parser := parser.New(lexer)
   program := parser.ParseProgram()
   ```

4. **Análise Semântica**
   ```go
   compiler := compiler.New(program, options)
   compiler.Analyze()
   ```

5. **Geração de Código**
   ```go
   output := compiler.Generate()
   ```

## Desenvolvimento

### Ambiente de Desenvolvimento

1. **Configuração do Go**
   ```bash
   go mod init github.com/peder1981/advpl-tlpp-compiler
   go mod tidy
   ```

2. **Estrutura de Diretórios**
   ```
   .
   ├── cmd/
   │   ├── compiler/
   │   └── lsp/
   ├── pkg/
   │   ├── ast/
   │   ├── compiler/
   │   ├── lexer/
   │   ├── parser/
   │   └── lsp/
   ├── docs/
   └── examples/
   ```

3. **Ferramentas Recomendadas**
   - VSCode com extensão Go
   - Delve para debugging
   - golangci-lint para linting

### Convenções de Código

1. **Nomenclatura**
   - Prefixo "ITX" para identificadores
   - CamelCase para tipos e funções
   - snake_case para variáveis privadas

2. **Documentação**
   - Comentários em português
   - Documentação de funções públicas
   - Exemplos de uso

3. **Organização**
   - Máximo 500 linhas por arquivo
   - Um pacote por diretório
   - Testes junto aos arquivos

## Testes

### Testes Unitários

```go
func TestLexer(t *testing.T) {
    input := `Function Test()`
    lexer := New(input)
    
    expected := []struct {
        Type    TokenType
        Literal string
    }{
        {TOKEN_FUNCTION, "Function"},
        {TOKEN_IDENT, "Test"},
        {TOKEN_LPAREN, "("},
        {TOKEN_RPAREN, ")"},
    }
    
    for _, exp := range expected {
        token := lexer.NextToken()
        if token.Type != exp.Type {
            t.Errorf("wrong token type")
        }
    }
}
```

### Testes de Integração

```go
func TestCompilation(t *testing.T) {
    input := `Function Test()
        Local x := 10
        Return x
    EndFunc`
    
    program, err := parser.ParseSource(input)
    if err != nil {
        t.Fatal(err)
    }
    
    compiler := New(program, Options{})
    output, err := compiler.Generate()
    if err != nil {
        t.Fatal(err)
    }
    
    // Verificar output
}
```

### Testes do LSP

```go
func TestLSPServer(t *testing.T) {
    server := NewServer()
    
    // Teste Initialize
    result, err := server.Initialize(context.Background(), InitializeParams{})
    if err != nil {
        t.Fatal(err)
    }
    
    // Teste DidOpen
    err = server.DidOpen(context.Background(), DidOpenTextDocumentParams{})
    if err != nil {
        t.Fatal(err)
    }
}
```

## Contribuição

1. **Fork do Repositório**
   ```bash
   git clone https://github.com/peder1981/advpl-tlpp-compiler
   cd advpl-tlpp-compiler
   git checkout -b feature/nova-funcionalidade
   ```

2. **Desenvolvimento**
   - Escreva testes primeiro
   - Implemente a funcionalidade
   - Documente as mudanças

3. **Pull Request**
   - Descreva as mudanças
   - Referencie issues
   - Aguarde review

4. **Code Review**
   - Verificação de estilo
   - Testes passando
   - Documentação atualizada

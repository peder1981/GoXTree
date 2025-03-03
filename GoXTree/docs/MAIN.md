# Documentação do Ponto de Entrada

## Visão Geral

O arquivo `main.go` é o ponto de entrada da aplicação GoXTree. Ele é responsável por inicializar a aplicação, processar argumentos de linha de comando e iniciar a interface do usuário.

## Localização

O arquivo `main.go` está localizado no diretório `cmd/gxtree/`:

```
GoXTree/
└── cmd/
    └── gxtree/
        └── main.go
```

## Estrutura

O arquivo `main.go` contém a função `main()`, que é o ponto de entrada da aplicação. Ele também pode conter outras funções auxiliares para processamento de argumentos de linha de comando e inicialização da aplicação.

```go
package main

import (
    "flag"
    "fmt"
    "os"
    "path/filepath"

    "github.com/peder1981/goxree/pkg/ui"
)

func main() {
    // Processar argumentos de linha de comando
    var (
        showHelp    bool
        showVersion bool
    )

    flag.BoolVar(&showHelp, "help", false, "Exibe informações de ajuda")
    flag.BoolVar(&showHelp, "h", false, "Exibe informações de ajuda (atalho)")
    flag.BoolVar(&showVersion, "version", false, "Exibe a versão do programa")
    flag.BoolVar(&showVersion, "v", false, "Exibe a versão do programa (atalho)")

    flag.Parse()

    // Exibir ajuda se solicitado
    if showHelp {
        printHelp()
        return
    }

    // Exibir versão se solicitado
    if showVersion {
        printVersion()
        return
    }

    // Determinar o diretório inicial
    var startDir string
    if flag.NArg() > 0 {
        // Usar o diretório fornecido como argumento
        startDir = flag.Arg(0)
        // Verificar se o diretório existe
        if _, err := os.Stat(startDir); os.IsNotExist(err) {
            fmt.Printf("Erro: diretório não encontrado: %s\n", startDir)
            os.Exit(1)
        }
    } else {
        // Usar o diretório atual
        var err error
        startDir, err = os.Getwd()
        if err != nil {
            fmt.Printf("Erro: não foi possível obter o diretório atual: %s\n", err)
            os.Exit(1)
        }
    }

    // Converter para caminho absoluto
    startDir, err := filepath.Abs(startDir)
    if err != nil {
        fmt.Printf("Erro: não foi possível obter o caminho absoluto: %s\n", err)
        os.Exit(1)
    }

    // Inicializar a aplicação
    app := ui.NewApp()
    app.SetRootDirectory(startDir)

    // Executar a aplicação
    if err := app.Run(); err != nil {
        fmt.Printf("Erro: %s\n", err)
        os.Exit(1)
    }
}

// printHelp exibe informações de ajuda
func printHelp() {
    fmt.Println("GoXTree - Gerenciador de Arquivos")
    fmt.Println("Uso: goxree [opções] [diretório]")
    fmt.Println()
    fmt.Println("Opções:")
    flag.PrintDefaults()
    fmt.Println()
    fmt.Println("Argumentos:")
    fmt.Println("  diretório    Diretório inicial (opcional, padrão: diretório atual)")
    fmt.Println()
    fmt.Println("Exemplos:")
    fmt.Println("  goxree                     # Inicia no diretório atual")
    fmt.Println("  goxree /home/user/docs     # Inicia no diretório especificado")
    fmt.Println("  goxree -h                  # Exibe esta ajuda")
    fmt.Println("  goxree -v                  # Exibe a versão do programa")
}

// printVersion exibe a versão do programa
func printVersion() {
    fmt.Println("GoXTree v1.0.0")
    fmt.Println("Copyright (c) 2025 Peder Munksgaard")
    fmt.Println("Licença: MIT")
}
```

## Fluxo de Execução

1. A função `main()` é chamada quando o programa é iniciado
2. Os argumentos de linha de comando são processados usando o pacote `flag`
3. Se a opção `-help` ou `-h` for especificada, a função `printHelp()` é chamada e o programa é encerrado
4. Se a opção `-version` ou `-v` for especificada, a função `printVersion()` é chamada e o programa é encerrado
5. O diretório inicial é determinado:
   - Se um diretório for fornecido como argumento, ele é usado
   - Caso contrário, o diretório atual é usado
6. O diretório inicial é convertido para um caminho absoluto
7. A aplicação é inicializada chamando `ui.NewApp()`
8. O diretório raiz da aplicação é definido chamando `app.SetRootDirectory(startDir)`
9. A aplicação é executada chamando `app.Run()`
10. Se ocorrer algum erro durante a execução, ele é exibido e o programa é encerrado com código de saída 1

## Argumentos de Linha de Comando

O GoXTree aceita os seguintes argumentos de linha de comando:

| Argumento | Descrição |
|-----------|-----------|
| `-help`, `-h` | Exibe informações de ajuda |
| `-version`, `-v` | Exibe a versão do programa |
| `diretório` | Diretório inicial (opcional, padrão: diretório atual) |

## Códigos de Saída

| Código | Descrição |
|--------|-----------|
| 0 | Sucesso |
| 1 | Erro |

## Personalização

O arquivo `main.go` pode ser personalizado para adicionar novos argumentos de linha de comando ou modificar o comportamento de inicialização da aplicação.

### Adicionando Novos Argumentos

Para adicionar um novo argumento de linha de comando, use o pacote `flag`:

```go
var newOption string
flag.StringVar(&newOption, "option", "default", "Descrição da opção")
```

### Modificando o Comportamento de Inicialização

Para modificar o comportamento de inicialização da aplicação, altere o código na função `main()` após o processamento dos argumentos de linha de comando.

## Considerações de Desempenho

- O arquivo `main.go` é executado apenas uma vez durante a inicialização da aplicação
- O processamento de argumentos de linha de comando é rápido e não afeta significativamente o desempenho da aplicação
- A inicialização da aplicação pode ser lenta se o diretório inicial contiver muitos arquivos

## Extensibilidade

O arquivo `main.go` pode ser estendido para adicionar novas funcionalidades:

1. Adicionar novos argumentos de linha de comando
2. Adicionar novas opções de configuração
3. Adicionar suporte para arquivos de configuração
4. Adicionar suporte para plugins ou extensões

Por exemplo, para adicionar suporte para um arquivo de configuração:

```go
var configFile string
flag.StringVar(&configFile, "config", "", "Arquivo de configuração")

// Carregar configuração
if configFile != "" {
    config, err := loadConfig(configFile)
    if err != nil {
        fmt.Printf("Erro: não foi possível carregar o arquivo de configuração: %s\n", err)
        os.Exit(1)
    }
    
    // Aplicar configuração
    app.ApplyConfig(config)
}

// Função para carregar configuração
func loadConfig(path string) (ui.Config, error) {
    // Implementação para carregar configuração
    // ...
    return config, nil
}
```

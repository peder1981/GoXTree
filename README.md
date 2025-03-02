# GoXTree

GoXTree é uma recriação moderna do lendário gerenciador de arquivos XTree Gold, implementado em Go. Este projeto visa trazer de volta a eficiência e a simplicidade do XTree Gold original, mas com tecnologias modernas e compatibilidade com sistemas operacionais atuais.

## Características

- Interface de texto com elementos gráficos usando TUI (Terminal User Interface)
- Visualização em árvore hierárquica de diretórios
- Visualizador de arquivos para múltiplos formatos (texto, hexadecimal, imagens, etc.)
- Editor de texto integrado
- Funções de busca para localizar arquivos
- Operações básicas de gerenciamento de arquivos (copiar, mover, excluir)
- Baixo consumo de recursos
- Interface intuitiva inspirada no XTree Gold original
- Navegação por teclas de função e atalhos de teclado
- Comparação de arquivos e diretórios
- Suporte a múltiplas plataformas (Windows, Linux, macOS)

## Requisitos

- Go 1.16 ou superior
- Bibliotecas:
  - github.com/gdamore/tcell/v2
  - github.com/rivo/tview
  - github.com/sergi/go-diff/diffmatchpatch

## Instalação

### Binários pré-compilados

Você pode baixar os binários pré-compilados para seu sistema operacional na seção de [Releases](https://github.com/peder1981/GoXTree/releases).

### A partir do código fonte

```bash
# Clone o repositório
git clone https://github.com/peder1981/GoXTree.git

# Entre no diretório
cd GoXTree

# Compile o projeto
go build -o bin/goxTree ./cmd/gxtree/main.go

# Execute o programa
./bin/goxTree
```

## Compilação para múltiplas plataformas

O GoXTree pode ser compilado para diversas plataformas e arquiteturas. Para compilar para todas as plataformas principais, execute o script:

```bash
# Entre no diretório do projeto
cd GoXTree

# Execute o script de compilação
./build_all.sh
```

Ou compile manualmente para cada plataforma:

```bash
# Linux AMD64
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/goxTree-linux-amd64 ./cmd/gxtree/main.go

# Linux ARM64
CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o bin/goxTree-linux-arm64 ./cmd/gxtree/main.go

# Windows AMD64
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o bin/goxTree-windows-amd64.exe ./cmd/gxtree/main.go

# macOS AMD64
CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o bin/goxTree-darwin-amd64 ./cmd/gxtree/main.go

# macOS ARM64 (Apple Silicon)
CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -o bin/goxTree-darwin-arm64 ./cmd/gxtree/main.go
```

## Arquitetura do Projeto

O GoXTree segue uma arquitetura modular com os seguintes componentes principais:

### Estrutura de Diretórios
```
GoXTree/
├── bin/                  # Binários compilados
├── cmd/
│   └── gxtree/
│       └── main.go       # Ponto de entrada da aplicação
├── pkg/
    ├── ui/               # Componentes da interface do usuário
    │   ├── app_core.go   # Núcleo da aplicação
    │   ├── app_navigation.go # Navegação entre diretórios
    │   ├── dialog.go     # Diálogos e popups
    │   ├── file_view.go  # Visualização de arquivos
    │   ├── menu_bar.go   # Barra de menu
    │   ├── status_bar.go # Barra de status
    │   └── tree_view.go  # Visualização em árvore
    └── utils/            # Utilitários e funções auxiliares
        ├── file_utils.go # Utilitários para manipulação de arquivos
        └── ui_utils.go   # Utilitários para a interface
```

## Uso

### Navegação

O GoXTree oferece uma interface intuitiva para navegação de arquivos e diretórios:

- Use as setas para navegar entre arquivos e diretórios
- Pressione `Enter` para entrar em um diretório ou abrir um arquivo
- Pressione `ESC` para voltar ao diretório anterior
- Pressione `Tab` para alternar entre a visualização em árvore e a lista de arquivos

### Teclas de Função

| Tecla | Função |
|-------|--------|
| F1 | Ajuda |
| F2 | Renomear |
| F3 | Visualizar |
| F4 | Editar |
| F5 | Copiar |
| F6 | Mover |
| F7 | Criar diretório |
| F8 | Excluir |
| F9 | Comprimir |
| F10 | Sair |

### Barra de Status

A barra de status exibe informações sobre o diretório atual, incluindo:
- Caminho completo do diretório
- Número de arquivos
- Número de diretórios
- Tamanho total do diretório

### Barra de Função

A barra de função na parte inferior da tela exibe as teclas de função disponíveis e suas ações correspondentes.

## Recursos Avançados

### Comparação de Arquivos

O GoXTree permite comparar o conteúdo de dois arquivos, destacando as diferenças entre eles:

1. Selecione o primeiro arquivo
2. Pressione `Alt+C` para marcá-lo para comparação
3. Navegue até o segundo arquivo
4. Pressione `Alt+C` novamente para iniciar a comparação

### Busca de Arquivos

Para buscar arquivos ou diretórios:

1. Pressione `Alt+F` para abrir o diálogo de busca
2. Digite o padrão de busca (suporta expressões regulares)
3. Pressione `Enter` para iniciar a busca

## Contribuição

Contribuições são bem-vindas! Sinta-se à vontade para abrir issues ou enviar pull requests.

## Licença

Este projeto está licenciado sob a [MIT License](LICENSE).

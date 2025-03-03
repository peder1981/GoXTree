# Guia de Desenvolvimento

Este documento fornece informações para desenvolvedores que desejam contribuir para o projeto GoXTree.

## Ambiente de Desenvolvimento

### Requisitos

- Go 1.16 ou superior
- Git
- Editor de código (VSCode, GoLand, etc.)
- Terminal Linux ou compatível

### Configuração do Ambiente

1. Clone o repositório:
   ```bash
   git clone https://github.com/peder1981/GoXTree.git
   cd GoXTree
   ```

2. Instale as dependências:
   ```bash
   go mod tidy
   ```

3. Compile o projeto:
   ```bash
   go build -o goxree ./cmd/gxtree/main.go
   ```

4. Execute o programa:
   ```bash
   ./goxree
   ```

## Estrutura do Código

O código está organizado em pacotes seguindo as melhores práticas de Go:

```
GoXTree/
├── cmd/
│   └── gxtree/
│       └── main.go       # Ponto de entrada da aplicação
├── pkg/
│   ├── editor/           # Componentes de edição de arquivos
│   │   └── texteditor.go # Editor de texto
│   ├── ui/               # Componentes da interface do usuário
│   │   ├── app.go        # Aplicação principal
│   │   ├── fileeditor.go # Interface para edição de arquivos
│   │   ├── fileview.go   # Visualização de arquivos
│   │   ├── fileviewer.go # Interface para visualização de arquivos
│   │   ├── helpview.go   # Tela de ajuda
│   │   ├── menubar.go    # Barra de menu
│   │   ├── statusbar.go  # Barra de status
│   │   └── treeview.go   # Visualização em árvore de diretórios
│   ├── utils/            # Utilitários
│   │   ├── constants.go  # Constantes utilizadas no projeto
│   │   └── fileutils.go  # Funções para manipulação de arquivos
│   └── viewer/           # Componentes de visualização de arquivos
│       ├── imageviewer.go # Visualizador de imagens
│       └── textviewer.go  # Visualizador de texto e hexadecimal
├── docs/                 # Documentação do projeto
└── README.md             # Documentação principal
```

## Padrões de Código

### Estilo de Código

Seguimos o estilo de código padrão do Go, conforme definido pelo `gofmt`. Antes de enviar um pull request, certifique-se de que seu código está formatado corretamente:

```bash
gofmt -w .
```

### Comentários

Todos os pacotes, tipos, funções e métodos exportados devem ter comentários de documentação seguindo o padrão do Go:

```go
// Package example fornece exemplos de código.
package example

// ExampleType é um exemplo de tipo.
type ExampleType struct {
    // Field é um campo de exemplo.
    Field string
}

// ExampleFunction é um exemplo de função.
//
// Parâmetros:
//   - param: Um parâmetro de exemplo
//
// Retorno:
//   - string: Um valor de retorno de exemplo
func ExampleFunction(param string) string {
    return param
}
```

### Tratamento de Erros

Erros devem ser tratados e propagados adequadamente. Evite ignorar erros sem uma boa razão:

```go
// Bom
file, err := os.Open("file.txt")
if err != nil {
    return err
}

// Ruim
file, _ := os.Open("file.txt") // Ignora o erro
```

### Testes

Todos os pacotes devem ter testes unitários. Os testes devem ser colocados no mesmo pacote que o código que estão testando, em arquivos com o sufixo `_test.go`:

```go
// fileutils_test.go
package utils

import "testing"

func TestCopyFile(t *testing.T) {
    // Teste para a função CopyFile
}
```

## Fluxo de Trabalho de Desenvolvimento

### Branches

- `main`: Branch principal, sempre estável
- `develop`: Branch de desenvolvimento, onde as novas funcionalidades são integradas
- `feature/<nome>`: Branches para novas funcionalidades
- `bugfix/<nome>`: Branches para correções de bugs
- `release/<versão>`: Branches para preparação de releases

### Processo de Desenvolvimento

1. Crie uma branch a partir de `develop`:
   ```bash
   git checkout develop
   git pull
   git checkout -b feature/nova-funcionalidade
   ```

2. Desenvolva a funcionalidade e adicione testes

3. Formate o código:
   ```bash
   gofmt -w .
   ```

4. Execute os testes:
   ```bash
   go test ./...
   ```

5. Commit e push:
   ```bash
   git add .
   git commit -m "Adiciona nova funcionalidade"
   git push origin feature/nova-funcionalidade
   ```

6. Crie um Pull Request para a branch `develop`

### Revisão de Código

Todos os Pull Requests devem ser revisados por pelo menos um outro desenvolvedor antes de serem mesclados. A revisão deve verificar:

- Funcionalidade: O código faz o que deveria fazer?
- Qualidade: O código segue os padrões de qualidade do projeto?
- Testes: Existem testes adequados para a funcionalidade?
- Documentação: O código está adequadamente documentado?

## Geração de Documentação

A documentação do código pode ser gerada usando o `godoc`:

```bash
# Instalar godoc
go install golang.org/x/tools/cmd/godoc@latest

# Gerar documentação
./generate_docs.sh
```

A documentação gerada estará disponível em `docs/godoc/`.

## Releases

### Versionamento

Seguimos o [Versionamento Semântico](https://semver.org/):

- MAJOR.MINOR.PATCH
  - MAJOR: Mudanças incompatíveis com versões anteriores
  - MINOR: Adições de funcionalidades compatíveis com versões anteriores
  - PATCH: Correções de bugs compatíveis com versões anteriores

### Processo de Release

1. Crie uma branch de release a partir de `develop`:
   ```bash
   git checkout develop
   git pull
   git checkout -b release/1.0.0
   ```

2. Atualize a versão no código e na documentação

3. Commit as alterações:
   ```bash
   git add .
   git commit -m "Prepara release 1.0.0"
   ```

4. Mescle a branch de release em `main` e `develop`:
   ```bash
   git checkout main
   git pull
   git merge --no-ff release/1.0.0
   git tag -a v1.0.0 -m "Versão 1.0.0"
   git push origin main --tags
   
   git checkout develop
   git pull
   git merge --no-ff release/1.0.0
   git push origin develop
   ```

5. Exclua a branch de release:
   ```bash
   git branch -d release/1.0.0
   ```

## Recursos Adicionais

- [Documentação do Go](https://golang.org/doc/)
- [Effective Go](https://golang.org/doc/effective_go)
- [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- [Documentação do tcell](https://pkg.go.dev/github.com/gdamore/tcell/v2)
- [Documentação do tview](https://pkg.go.dev/github.com/rivo/tview)

# Guia para Desenvolvedores do GoXTree

Este documento fornece informações detalhadas para desenvolvedores que desejam contribuir com o projeto GoXTree.

## Estrutura do Código

O GoXTree segue uma arquitetura modular com separação clara de responsabilidades:

### Componentes Principais

1. **Núcleo da Aplicação** (`pkg/ui/app_core.go`)
   - Gerencia o ciclo de vida da aplicação
   - Integra todos os componentes da interface
   - Gerencia eventos globais de teclado

2. **Navegação** (`pkg/ui/app_navigation.go`)
   - Implementa a navegação entre diretórios
   - Gerencia o histórico de navegação
   - Fornece funções para navegar para diretórios específicos

3. **Visualização em Árvore** (`pkg/ui/tree_view.go`)
   - Exibe a estrutura hierárquica de diretórios
   - Permite navegação e seleção de diretórios

4. **Visualização de Arquivos** (`pkg/ui/file_view.go`)
   - Exibe os arquivos no diretório atual
   - Permite seleção e manipulação de arquivos

5. **Barra de Status** (`pkg/ui/status_bar.go`)
   - Exibe informações sobre o diretório atual
   - Mostra estatísticas como número de arquivos e tamanho total

6. **Barra de Menu** (`pkg/ui/menu_bar.go`)
   - Exibe o menu principal da aplicação
   - Fornece acesso às principais funcionalidades

7. **Diálogos** (`pkg/ui/dialog.go`)
   - Implementa diálogos modais para interação com o usuário
   - Inclui diálogos de confirmação, entrada de texto e mensagens

8. **Utilitários** (`pkg/utils/`)
   - Fornece funções auxiliares para manipulação de arquivos e interface

## Fluxo de Execução

1. O ponto de entrada da aplicação é `cmd/gxtree/main.go`
2. A função `main()` cria uma nova instância da aplicação (`app_core.go:NewApp()`)
3. A aplicação inicializa todos os componentes da interface
4. O método `Run()` inicia o loop principal da aplicação
5. A aplicação processa eventos de teclado e atualiza a interface conforme necessário

## Convenções de Código

### Nomenclatura

- **Arquivos**: Nomes em snake_case (ex: `file_view.go`)
- **Estruturas**: Nomes em PascalCase (ex: `FileView`)
- **Métodos**: Nomes em camelCase (ex: `loadTree()`)
- **Variáveis**: Nomes em camelCase (ex: `currentDir`)
- **Constantes**: Nomes em SCREAMING_SNAKE_CASE (ex: `MAX_DEPTH`)

### Organização de Código

- Cada componente da interface é definido em seu próprio arquivo
- Métodos relacionados são agrupados juntos
- Comentários de documentação seguem o padrão Go

## Adicionando Novos Recursos

### Novo Componente de Interface

1. Crie um novo arquivo em `pkg/ui/` para o componente
2. Defina uma estrutura para o componente
3. Implemente um construtor (`New[ComponentName]()`)
4. Integre o componente em `app_core.go`

### Nova Funcionalidade

1. Identifique o componente apropriado para a funcionalidade
2. Implemente a lógica necessária
3. Adicione manipuladores de eventos de teclado, se necessário
4. Atualize a documentação

## Tratamento de Erros

- Use o padrão Go de retornar erros como último valor de retorno
- Utilize o método `showError()` para exibir erros ao usuário
- Registre erros detalhados para depuração

## Compilação e Testes

### Compilação Local

```bash
go build -o bin/goxTree ./cmd/gxtree/main.go
```

### Compilação para Múltiplas Plataformas

```bash
./build_all.sh
```

### Testes

```bash
go test ./...
```

## Diretrizes de Contribuição

1. Faça um fork do repositório
2. Crie um branch para sua feature (`git checkout -b feature/nova-funcionalidade`)
3. Faça commit das suas alterações (`git commit -m 'Adiciona nova funcionalidade'`)
4. Envie para o branch (`git push origin feature/nova-funcionalidade`)
5. Abra um Pull Request

## Recursos Úteis

- [Documentação do Go](https://golang.org/doc/)
- [Documentação do tcell](https://github.com/gdamore/tcell)
- [Documentação do tview](https://github.com/rivo/tview)
- [XTree Gold - Wikipedia](https://en.wikipedia.org/wiki/XTree)

## Contato

Para dúvidas ou sugestões, entre em contato com o mantenedor do projeto:
- GitHub: [@peder1981](https://github.com/peder1981)

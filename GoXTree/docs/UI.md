# Documentação da Interface do Usuário

## Visão Geral

A interface do usuário do GoXTree é implementada usando as bibliotecas [tcell](https://github.com/gdamore/tcell) e [tview](https://github.com/rivo/tview), que fornecem uma base para criar interfaces de terminal ricas e interativas.

## Componentes da Interface

### App (app.go)

O componente App é o controlador principal da aplicação. Ele gerencia o layout da interface, processa eventos de teclado e coordena a interação entre os diferentes componentes.

#### Métodos Principais

- `NewApp()`: Cria uma nova instância do App
- `Run()`: Inicia a execução da aplicação
- `Stop()`: Interrompe a execução da aplicação
- `SetupLayout()`: Configura o layout da interface
- `ProcessKey(key tcell.Key, mod tcell.ModMask)`: Processa eventos de teclado
- `ShowInputDialog(title, label string, callback func(string))`: Exibe um diálogo de entrada de texto
- `ShowConfirmDialog(title, message string, callback func(bool))`: Exibe um diálogo de confirmação
- `ShowErrorDialog(title, message string)`: Exibe um diálogo de erro
- `ShowInfoDialog(title, message string)`: Exibe um diálogo de informação
- `NavigateToDirectory(path string)`: Navega para um diretório específico
- `RefreshStatus()`: Atualiza a barra de status

#### Estrutura

```go
type App struct {
    app            *tview.Application
    layout         *tview.Flex
    treeView       *TreeView
    fileView       *FileView
    menuBar        *MenuBar
    statusBar      *StatusBar
    currentDir     string
    currentPanel   int
    selectedFile   string
    markedFiles    map[string]bool
}
```

### TreeView (treeview.go)

O componente TreeView exibe a estrutura hierárquica de diretórios e permite a navegação entre eles.

#### Métodos Principais

- `NewTreeView(rootDir string)`: Cria uma nova instância do TreeView
- `SetRootDirectory(rootDir string)`: Define o diretório raiz
- `GetCurrentDirectory()`: Retorna o diretório atual
- `ExpandDirectory(path string)`: Expande um diretório
- `CollapseDirectory(path string)`: Colapsa um diretório
- `SelectDirectory(path string)`: Seleciona um diretório
- `ProcessKey(key tcell.Key, mod tcell.ModMask)`: Processa eventos de teclado
- `SetChangeDirCallback(callback func(string))`: Define o callback para mudança de diretório

#### Estrutura

```go
type TreeView struct {
    *tview.TreeView
    rootDir           string
    currentDir        string
    changeDirCallback func(string)
}
```

### FileView (fileview.go)

O componente FileView exibe os arquivos no diretório selecionado, com informações como nome, tamanho, data de modificação, etc.

#### Métodos Principais

- `NewFileView()`: Cria uma nova instância do FileView
- `SetDirectory(dir string)`: Define o diretório a ser exibido
- `GetCurrentDirectory()`: Retorna o diretório atual
- `GetSelectedFile()`: Retorna o arquivo selecionado
- `MarkFile(file string)`: Marca um arquivo para operações em lote
- `UnmarkFile(file string)`: Desmarca um arquivo
- `MarkAll()`: Marca todos os arquivos
- `UnmarkAll()`: Desmarca todos os arquivos
- `GetMarkedFiles()`: Retorna os arquivos marcados
- `ProcessKey(key tcell.Key, mod tcell.ModMask)`: Processa eventos de teclado
- `SetFileActionCallback(callback func(string, string))`: Define o callback para ações em arquivos
- `GetFileCount()`: Retorna o número de arquivos no diretório atual
- `GetDirCount()`: Retorna o número de diretórios no diretório atual

#### Estrutura

```go
type FileView struct {
    *tview.Table
    currentDir        string
    files             []os.FileInfo
    selectedRow       int
    markedFiles       map[string]bool
    fileActionCallback func(string, string)
}
```

### MenuBar (menubar.go)

O componente MenuBar exibe os menus disponíveis e processa as ações selecionadas.

#### Métodos Principais

- `NewMenuBar()`: Cria uma nova instância do MenuBar
- `AddMenu(name string, items []MenuItem)`: Adiciona um menu
- `ShowMenu(name string)`: Exibe um menu
- `HideMenu()`: Esconde o menu atual
- `ProcessKey(key tcell.Key, mod tcell.ModMask)`: Processa eventos de teclado

#### Estrutura

```go
type MenuBar struct {
    *tview.Flex
    menus       map[string]*tview.List
    currentMenu string
    visible     bool
}

type MenuItem struct {
    Label    string
    Shortcut string
    Action   func()
}
```

### StatusBar (statusbar.go)

O componente StatusBar exibe informações sobre o estado atual da aplicação, como diretório atual, número de arquivos, espaço livre, etc.

#### Métodos Principais

- `NewStatusBar()`: Cria uma nova instância do StatusBar
- `SetDirectory(dir string)`: Define o diretório atual
- `SetFileCount(count int)`: Define o número de arquivos
- `SetDirCount(count int)`: Define o número de diretórios
- `SetDiskInfo(free, total int64)`: Define informações sobre o disco
- `Update()`: Atualiza as informações exibidas

#### Estrutura

```go
type StatusBar struct {
    *tview.TextView
    currentDir string
    fileCount  int
    dirCount   int
    freeSpace  int64
    totalSpace int64
}
```

### HelpView (helpview.go)

O componente HelpView exibe informações sobre como usar a aplicação, atalhos de teclado, etc.

#### Métodos Principais

- `NewHelpView()`: Cria uma nova instância do HelpView
- `Show()`: Exibe a tela de ajuda
- `Hide()`: Esconde a tela de ajuda
- `ProcessKey(key tcell.Key, mod tcell.ModMask)`: Processa eventos de teclado

#### Estrutura

```go
type HelpView struct {
    *tview.TextView
    visible bool
}
```

### FileViewer (fileviewer.go)

O componente FileViewer coordena a visualização de arquivos, selecionando o visualizador apropriado com base no tipo de arquivo.

#### Métodos Principais

- `NewFileViewer()`: Cria uma nova instância do FileViewer
- `ViewFile(path string)`: Visualiza um arquivo
- `Close()`: Fecha o visualizador atual
- `ProcessKey(key tcell.Key, mod tcell.ModMask)`: Processa eventos de teclado

#### Estrutura

```go
type FileViewer struct {
    *tview.Pages
    currentViewer Viewer
    currentFile   string
}

type Viewer interface {
    View(path string) error
    Close()
    ProcessKey(key tcell.Key, mod tcell.ModMask) bool
}
```

### FileEditor (fileeditor.go)

O componente FileEditor coordena a edição de arquivos, selecionando o editor apropriado com base no tipo de arquivo.

#### Métodos Principais

- `NewFileEditor()`: Cria uma nova instância do FileEditor
- `EditFile(path string)`: Edita um arquivo
- `Close()`: Fecha o editor atual
- `ProcessKey(key tcell.Key, mod tcell.ModMask)`: Processa eventos de teclado

#### Estrutura

```go
type FileEditor struct {
    *tview.Pages
    currentEditor Editor
    currentFile   string
}

type Editor interface {
    Edit(path string) error
    Save() error
    Close()
    ProcessKey(key tcell.Key, mod tcell.ModMask) bool
}
```

## Fluxo de Interação

1. O usuário inicia a aplicação, que cria uma instância do App
2. O App configura o layout da interface, criando instâncias dos componentes TreeView, FileView, MenuBar e StatusBar
3. O App inicia a execução da aplicação, exibindo a interface e processando eventos de teclado
4. O usuário navega pelos diretórios usando o TreeView, que notifica o App sobre mudanças de diretório
5. O App atualiza o FileView para exibir os arquivos no diretório selecionado
6. O usuário seleciona arquivos no FileView e executa operações como visualização, edição, cópia, movimentação, etc.
7. O App coordena essas operações, atualizando a interface conforme necessário
8. O usuário pode acessar menus através do MenuBar para executar operações adicionais
9. O StatusBar exibe informações sobre o estado atual da aplicação
10. O usuário pode acessar a tela de ajuda através do HelpView para obter informações sobre como usar a aplicação

## Personalização da Interface

A interface do GoXTree pode ser personalizada através de configurações definidas no arquivo constants.go:

- Cores de fundo e texto
- Prefixos e sufixos para diretórios e arquivos
- Formatação de datas e tamanhos
- Atalhos de teclado

## Considerações de Acessibilidade

- A interface é projetada para ser acessível através do teclado, sem necessidade de mouse
- As cores são escolhidas para garantir contraste adequado
- Mensagens de erro e informação são claras e informativas

# Documentação dos Visualizadores e Editores

## Visão Geral

O GoXTree inclui componentes para visualização e edição de diferentes tipos de arquivos. Esses componentes são projetados para serem modulares e extensíveis, permitindo a adição de suporte para novos tipos de arquivos no futuro.

## Visualizadores (pkg/viewer)

Os visualizadores são componentes responsáveis por exibir o conteúdo de arquivos em diferentes formatos. O GoXTree inclui visualizadores para arquivos de texto, imagens e visualização hexadecimal.

### Interface Viewer

Todos os visualizadores implementam a interface `Viewer`:

```go
type Viewer interface {
    View(path string) error
    Close()
    ProcessKey(key tcell.Key, mod tcell.ModMask) bool
}
```

- `View(path string) error`: Abre e exibe o conteúdo do arquivo especificado
- `Close()`: Fecha o visualizador e libera recursos
- `ProcessKey(key tcell.Key, mod tcell.ModMask) bool`: Processa eventos de teclado específicos do visualizador

### TextViewer (textviewer.go)

O TextViewer é responsável por exibir o conteúdo de arquivos de texto em formato legível.

#### Métodos Principais

- `NewTextViewer()`: Cria uma nova instância do TextViewer
- `View(path string) error`: Abre e exibe o conteúdo do arquivo de texto
- `Close()`: Fecha o visualizador e libera recursos
- `ProcessKey(key tcell.Key, mod tcell.ModMask) bool`: Processa eventos de teclado
- `ScrollUp()`: Rola o conteúdo para cima
- `ScrollDown()`: Rola o conteúdo para baixo
- `ScrollPageUp()`: Rola uma página para cima
- `ScrollPageDown()`: Rola uma página para baixo
- `GoToTop()`: Vai para o início do arquivo
- `GoToBottom()`: Vai para o final do arquivo
- `Search(text string)`: Busca texto no arquivo

#### Estrutura

```go
type TextViewer struct {
    *tview.TextView
    filePath     string
    content      []string
    currentLine  int
    searchText   string
    searchIndex  int
}
```

#### Funcionalidades

- Exibição de arquivos de texto com quebra de linha
- Navegação pelo conteúdo (rolagem, paginação)
- Busca de texto
- Destaque de sintaxe para linguagens de programação comuns
- Exibição de números de linha

### HexViewer (hexviewer.go)

O HexViewer é responsável por exibir o conteúdo de arquivos em formato hexadecimal, útil para visualizar arquivos binários.

#### Métodos Principais

- `NewHexViewer()`: Cria uma nova instância do HexViewer
- `View(path string) error`: Abre e exibe o conteúdo do arquivo em formato hexadecimal
- `Close()`: Fecha o visualizador e libera recursos
- `ProcessKey(key tcell.Key, mod tcell.ModMask) bool`: Processa eventos de teclado
- `ScrollUp()`: Rola o conteúdo para cima
- `ScrollDown()`: Rola o conteúdo para baixo
- `ScrollPageUp()`: Rola uma página para cima
- `ScrollPageDown()`: Rola uma página para baixo
- `GoToOffset(offset int64)`: Vai para um deslocamento específico no arquivo

#### Estrutura

```go
type HexViewer struct {
    *tview.TextView
    filePath     string
    content      []byte
    offset       int64
    bytesPerLine int
}
```

#### Funcionalidades

- Exibição de arquivos em formato hexadecimal
- Navegação pelo conteúdo (rolagem, paginação)
- Ir para um deslocamento específico
- Exibição de valores ASCII correspondentes
- Exibição de deslocamentos

### ImageViewer (imageviewer.go)

O ImageViewer é responsável por exibir imagens em formato ASCII art no terminal.

#### Métodos Principais

- `NewImageViewer()`: Cria uma nova instância do ImageViewer
- `View(path string) error`: Abre e exibe a imagem em formato ASCII art
- `Close()`: Fecha o visualizador e libera recursos
- `ProcessKey(key tcell.Key, mod tcell.ModMask) bool`: Processa eventos de teclado
- `ZoomIn()`: Aumenta o zoom da imagem
- `ZoomOut()`: Diminui o zoom da imagem
- `ResetZoom()`: Restaura o zoom original

#### Estrutura

```go
type ImageViewer struct {
    *tview.TextView
    filePath     string
    image        image.Image
    zoomLevel    float64
}
```

#### Funcionalidades

- Exibição de imagens em formato ASCII art
- Suporte para diferentes formatos de imagem (JPEG, PNG, GIF, BMP)
- Zoom in/out
- Ajuste automático ao tamanho do terminal

## Editores (pkg/editor)

Os editores são componentes responsáveis por permitir a edição de diferentes tipos de arquivos. O GoXTree inclui um editor de texto básico.

### Interface Editor

Todos os editores implementam a interface `Editor`:

```go
type Editor interface {
    Edit(path string) error
    Save() error
    Close()
    ProcessKey(key tcell.Key, mod tcell.ModMask) bool
}
```

- `Edit(path string) error`: Abre o arquivo especificado para edição
- `Save() error`: Salva as alterações no arquivo
- `Close()`: Fecha o editor e libera recursos
- `ProcessKey(key tcell.Key, mod tcell.ModMask) bool`: Processa eventos de teclado específicos do editor

### TextEditor (texteditor.go)

O TextEditor é responsável por permitir a edição de arquivos de texto.

#### Métodos Principais

- `NewTextEditor()`: Cria uma nova instância do TextEditor
- `Edit(path string) error`: Abre o arquivo de texto para edição
- `Save() error`: Salva as alterações no arquivo
- `Close()`: Fecha o editor e libera recursos
- `ProcessKey(key tcell.Key, mod tcell.ModMask) bool`: Processa eventos de teclado
- `InsertChar(ch rune)`: Insere um caractere na posição atual do cursor
- `DeleteChar()`: Exclui o caractere na posição atual do cursor
- `DeleteLine()`: Exclui a linha atual
- `NewLine()`: Insere uma nova linha
- `MoveCursor(direction int)`: Move o cursor na direção especificada
- `GoToLine(line int)`: Vai para a linha especificada
- `Search(text string)`: Busca texto no arquivo
- `Replace(search, replace string)`: Substitui texto no arquivo

#### Estrutura

```go
type TextEditor struct {
    *tview.TextView
    filePath     string
    content      []string
    cursorX      int
    cursorY      int
    modified     bool
    searchText   string
    searchIndex  int
}
```

#### Funcionalidades

- Edição de arquivos de texto
- Navegação pelo conteúdo (cursor, rolagem, paginação)
- Busca e substituição de texto
- Destaque de sintaxe para linguagens de programação comuns
- Exibição de números de linha
- Indicação de modificações não salvas
- Desfazer/refazer operações

## Integração com a Interface do Usuário

Os visualizadores e editores são integrados à interface do usuário através dos componentes FileViewer e FileEditor, que selecionam o visualizador ou editor apropriado com base no tipo de arquivo.

### Seleção de Visualizador

O componente FileViewer seleciona o visualizador apropriado com base na extensão do arquivo:

```go
func (fv *FileViewer) ViewFile(path string) error {
    ext := strings.ToLower(filepath.Ext(path))
    var viewer Viewer
    
    switch ext {
    case ".jpg", ".jpeg", ".png", ".gif", ".bmp":
        viewer = NewImageViewer()
    case ".txt", ".go", ".c", ".cpp", ".h", ".py", ".js", ".html", ".css", ".md", ".json", ".xml", ".yaml", ".yml":
        viewer = NewTextViewer()
    default:
        viewer = NewHexViewer()
    }
    
    fv.currentViewer = viewer
    fv.currentFile = path
    
    return viewer.View(path)
}
```

### Seleção de Editor

O componente FileEditor seleciona o editor apropriado com base na extensão do arquivo:

```go
func (fe *FileEditor) EditFile(path string) error {
    ext := strings.ToLower(filepath.Ext(path))
    var editor Editor
    
    switch ext {
    case ".txt", ".go", ".c", ".cpp", ".h", ".py", ".js", ".html", ".css", ".md", ".json", ".xml", ".yaml", ".yml":
        editor = NewTextEditor()
    default:
        return fmt.Errorf("unsupported file type for editing: %s", ext)
    }
    
    fe.currentEditor = editor
    fe.currentFile = path
    
    return editor.Edit(path)
}
```

## Extensibilidade

O sistema de visualizadores e editores é projetado para ser facilmente extensível. Para adicionar suporte para um novo tipo de arquivo:

1. Implemente a interface `Viewer` ou `Editor` para o novo tipo de arquivo
2. Atualize a lógica de seleção no FileViewer ou FileEditor para usar o novo visualizador ou editor para as extensões de arquivo apropriadas

Por exemplo, para adicionar um visualizador de PDF:

```go
type PDFViewer struct {
    *tview.TextView
    filePath string
    document *pdf.Document
    page     int
}

func NewPDFViewer() *PDFViewer {
    return &PDFViewer{
        TextView: tview.NewTextView(),
        page:     1,
    }
}

func (pv *PDFViewer) View(path string) error {
    // Implementação para abrir e exibir um arquivo PDF
    // ...
    return nil
}

func (pv *PDFViewer) Close() {
    // Implementação para fechar o visualizador
    // ...
}

func (pv *PDFViewer) ProcessKey(key tcell.Key, mod tcell.ModMask) bool {
    // Implementação para processar eventos de teclado
    // ...
    return true
}
```

E então atualizar a lógica de seleção no FileViewer:

```go
func (fv *FileViewer) ViewFile(path string) error {
    ext := strings.ToLower(filepath.Ext(path))
    var viewer Viewer
    
    switch ext {
    case ".jpg", ".jpeg", ".png", ".gif", ".bmp":
        viewer = NewImageViewer()
    case ".txt", ".go", ".c", ".cpp", ".h", ".py", ".js", ".html", ".css", ".md", ".json", ".xml", ".yaml", ".yml":
        viewer = NewTextViewer()
    case ".pdf":
        viewer = NewPDFViewer()
    default:
        viewer = NewHexViewer()
    }
    
    fv.currentViewer = viewer
    fv.currentFile = path
    
    return viewer.View(path)
}
```

# Documentação dos Utilitários

## Visão Geral

O pacote `utils` contém funções e constantes utilitárias utilizadas em todo o projeto GoXTree. Esses utilitários fornecem funcionalidades comuns para manipulação de arquivos, formatação de dados e definição de constantes.

## FileUtils (fileutils.go)

O arquivo `fileutils.go` contém funções para manipulação de arquivos e diretórios.

### Funções Principais

#### Manipulação de Arquivos

```go
// CopyFile copia um arquivo de origem para destino
func CopyFile(src, dst string) error

// MoveFile move um arquivo de origem para destino
func MoveFile(src, dst string) error

// DeleteFile exclui um arquivo
func DeleteFile(path string) error

// CreateDirectory cria um diretório
func CreateDirectory(path string) error

// DeleteDirectory exclui um diretório e seu conteúdo
func DeleteDirectory(path string) error

// RenameFile renomeia um arquivo ou diretório
func RenameFile(oldPath, newPath string) error
```

#### Informações sobre Arquivos

```go
// GetFileSize retorna o tamanho de um arquivo em bytes
func GetFileSize(path string) (int64, error)

// GetFileModTime retorna a data de modificação de um arquivo
func GetFileModTime(path string) (time.Time, error)

// GetFilePermissions retorna as permissões de um arquivo
func GetFilePermissions(path string) (os.FileMode, error)

// IsDirectory verifica se o caminho é um diretório
func IsDirectory(path string) (bool, error)

// IsFile verifica se o caminho é um arquivo
func IsFile(path string) (bool, error)

// IsSymlink verifica se o caminho é um link simbólico
func IsSymlink(path string) (bool, error)

// IsHidden verifica se o arquivo ou diretório é oculto
func IsHidden(path string) (bool, error)
```

#### Listagem de Arquivos e Diretórios

```go
// ListFiles lista os arquivos em um diretório
func ListFiles(dir string) ([]os.FileInfo, error)

// ListDirectories lista os diretórios em um diretório
func ListDirectories(dir string) ([]os.FileInfo, error)

// ListAll lista todos os arquivos e diretórios em um diretório
func ListAll(dir string) ([]os.FileInfo, error)

// FindFiles encontra arquivos que correspondem a um padrão
func FindFiles(dir, pattern string) ([]string, error)

// CountFiles conta o número de arquivos em um diretório
func CountFiles(dir string) (int, error)

// CountDirectories conta o número de diretórios em um diretório
func CountDirectories(dir string) (int, error)
```

#### Informações sobre Disco

```go
// GetDiskUsage retorna o espaço total e livre em um disco
func GetDiskUsage(path string) (total, free int64, err error)

// GetDiskUsagePercent retorna a porcentagem de uso de um disco
func GetDiskUsagePercent(path string) (float64, error)
```

#### Formatação de Dados

```go
// FormatFileSize formata um tamanho de arquivo em uma string legível
func FormatFileSize(size int64) string

// FormatDate formata uma data em uma string legível
func FormatDate(date time.Time) string

// FormatPermissions formata permissões de arquivo em uma string legível
func FormatPermissions(mode os.FileMode) string
```

### Exemplo de Uso

```go
// Listar arquivos em um diretório
files, err := utils.ListFiles("/home/user/documents")
if err != nil {
    log.Fatal(err)
}

// Exibir informações sobre cada arquivo
for _, file := range files {
    size, _ := utils.GetFileSize(filepath.Join("/home/user/documents", file.Name()))
    modTime, _ := utils.GetFileModTime(filepath.Join("/home/user/documents", file.Name()))
    
    fmt.Printf("Nome: %s\n", file.Name())
    fmt.Printf("Tamanho: %s\n", utils.FormatFileSize(size))
    fmt.Printf("Modificado em: %s\n", utils.FormatDate(modTime))
    fmt.Println("---")
}

// Copiar um arquivo
err = utils.CopyFile("/home/user/documents/file.txt", "/home/user/backup/file.txt")
if err != nil {
    log.Fatal(err)
}

// Obter informações sobre o disco
total, free, err := utils.GetDiskUsage("/home")
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Espaço total: %s\n", utils.FormatFileSize(total))
fmt.Printf("Espaço livre: %s\n", utils.FormatFileSize(free))
fmt.Printf("Uso: %.2f%%\n", float64(total-free)/float64(total)*100)
```

## Constants (constants.go)

O arquivo `constants.go` contém constantes utilizadas em todo o projeto.

### Constantes Principais

#### Teclas de Atalho

```go
const (
    // Teclas de função
    KeyHelp       = tcell.KeyF1
    KeyMenu       = tcell.KeyF2
    KeyView       = tcell.KeyF3
    KeyEdit       = tcell.KeyF4
    KeyCopy       = tcell.KeyF5
    KeyMove       = tcell.KeyF6
    KeyMkdir      = tcell.KeyF7
    KeyDelete     = tcell.KeyF8
    KeyCompress   = tcell.KeyF9
    KeyQuit       = tcell.KeyF10
    
    // Teclas de controle
    KeyMarkAll    = tcell.KeyCtrlA
    KeySearch     = tcell.KeyCtrlF
    KeyGoToDir    = tcell.KeyCtrlG
    KeyToggleHidden = tcell.KeyCtrlH
    KeyRefresh    = tcell.KeyCtrlR
)
```

#### Cores e Estilos

```go
const (
    // Cores
    ColorBackground = tcell.ColorBlack
    ColorText       = tcell.ColorWhite
    ColorHighlight  = tcell.ColorYellow
    ColorSelected   = tcell.ColorBlue
    ColorMarked     = tcell.ColorGreen
    ColorError      = tcell.ColorRed
    ColorInfo       = tcell.ColorCyan
    
    // Estilos
    StyleNormal     = tcell.StyleDefault.Background(ColorBackground).Foreground(ColorText)
    StyleHighlight  = tcell.StyleDefault.Background(ColorBackground).Foreground(ColorHighlight)
    StyleSelected   = tcell.StyleDefault.Background(ColorSelected).Foreground(ColorText)
    StyleMarked     = tcell.StyleDefault.Background(ColorBackground).Foreground(ColorMarked)
    StyleError      = tcell.StyleDefault.Background(ColorBackground).Foreground(ColorError)
    StyleInfo       = tcell.StyleDefault.Background(ColorBackground).Foreground(ColorInfo)
)
```

#### Prefixos e Sufixos

```go
const (
    // Prefixos para diretórios e arquivos
    PrefixDirectory = "📁 "
    PrefixFile      = "📄 "
    PrefixSymlink   = "🔗 "
    PrefixHidden    = "🔒 "
    
    // Sufixos para diretórios e arquivos
    SuffixDirectory = "/"
    SuffixExecutable = "*"
    SuffixSymlink   = "@"
)
```

#### Painéis

```go
const (
    // Identificadores de painéis
    PanelTree = 0
    PanelFile = 1
)
```

#### Formatação

```go
const (
    // Formatação de datas
    DateFormat = "2006-01-02 15:04:05"
    
    // Formatação de tamanhos
    SizeUnitB  = "B"
    SizeUnitKB = "KB"
    SizeUnitMB = "MB"
    SizeUnitGB = "GB"
    SizeUnitTB = "TB"
)
```

### Exemplo de Uso

```go
// Definir o estilo de um texto
textView.SetTextStyle(constants.StyleNormal)

// Verificar uma tecla de atalho
if key == constants.KeyHelp {
    showHelp()
}

// Formatar um diretório
dirName := constants.PrefixDirectory + dir.Name() + constants.SuffixDirectory

// Identificar o painel atual
if currentPanel == constants.PanelTree {
    // Processar eventos no painel de árvore
} else if currentPanel == constants.PanelFile {
    // Processar eventos no painel de arquivos
}
```

## Considerações de Desempenho

- As funções de manipulação de arquivos são projetadas para serem eficientes, mas podem ser lentas para arquivos muito grandes ou operações em lote
- A função `ListAll` carrega todos os arquivos e diretórios em memória, o que pode ser um problema para diretórios com muitos arquivos
- As funções de formatação são otimizadas para exibição na interface do usuário, não para processamento em lote

## Extensibilidade

O pacote `utils` é projetado para ser facilmente extensível. Para adicionar novas funcionalidades:

1. Adicione novas funções ao arquivo `fileutils.go` ou crie novos arquivos para categorias específicas de utilitários
2. Adicione novas constantes ao arquivo `constants.go` ou crie novos arquivos para categorias específicas de constantes
3. Atualize a documentação para refletir as novas funcionalidades

Por exemplo, para adicionar suporte para compressão de arquivos:

```go
// CompressFiles comprime arquivos em um arquivo ZIP
func CompressFiles(files []string, dst string) error {
    // Implementação para comprimir arquivos
    // ...
    return nil
}

// DecompressFile descomprime um arquivo ZIP
func DecompressFile(src, dst string) error {
    // Implementação para descomprimir arquivos
    // ...
    return nil
}

// IsCompressed verifica se um arquivo é comprimido
func IsCompressed(path string) (bool, error) {
    // Implementação para verificar se um arquivo é comprimido
    // ...
    return false, nil
}
```

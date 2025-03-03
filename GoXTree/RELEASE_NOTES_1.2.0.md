# GoXTree v1.2.0 - Correções e Melhorias

## Visão Geral

Esta versão traz importantes correções de compatibilidade com as versões mais recentes das bibliotecas tcell e tview, além de melhorias gerais no código e desempenho.

## Principais Correções

### Compatibilidade com a versão mais recente da biblioteca tview
- Correção do método PasteHandler() em helpview.go
- Remoção de referências a campos inexistentes em tview.Borders (Left, Right, Top, Bottom)
- Remoção de referências a tview.TreeGraphics que não existe na versão atual

### Correção de chamadas de métodos
- Substituição de syncDirectoriesDialog() por syncDirectories()
- Correção da chamada para textArea.GetText() adicionando o parâmetro necessário

### Correção de cores
- Substituição de tcell.ColorMagenta e tcell.ColorCyan por tcell.NewRGBColor
- Correção de referências a ColorSelected

### Limpeza de código
- Remoção de imports não utilizados
- Remoção de variáveis declaradas e não utilizadas
- Padronização de métodos e funções

## Melhorias

- Estrutura geral do código
- Tratamento de erros
- Desempenho da aplicação

## Requisitos

- Go 1.16 ou superior
- Bibliotecas:
  - github.com/gdamore/tcell/v2
  - github.com/rivo/tview
  - github.com/sergi/go-diff/diffmatchpatch

## Downloads

- Linux AMD64: [goxTree-linux-amd64](https://github.com/peder1981/GoXTree/releases/download/v1.2.0/goxTree-linux-amd64)
- Linux ARM64: [goxTree-linux-arm64](https://github.com/peder1981/GoXTree/releases/download/v1.2.0/goxTree-linux-arm64)
- Windows AMD64: [goxTree-windows-amd64.exe](https://github.com/peder1981/GoXTree/releases/download/v1.2.0/goxTree-windows-amd64.exe)
- macOS AMD64: [goxTree-darwin-amd64](https://github.com/peder1981/GoXTree/releases/download/v1.2.0/goxTree-darwin-amd64)
- macOS ARM64: [goxTree-darwin-arm64](https://github.com/peder1981/GoXTree/releases/download/v1.2.0/goxTree-darwin-arm64)

## Instalação

1. Baixe o binário correspondente ao seu sistema operacional
2. Torne-o executável (Linux/macOS): `chmod +x goxTree-*`
3. Execute o binário: `./goxTree-*`

## Agradecimentos

Agradecemos a todos que contribuíram para esta versão com sugestões, relatórios de bugs e melhorias no código.

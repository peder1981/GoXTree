# Documentação do GoXTree

Bem-vindo à documentação do GoXTree, uma recriação moderna do lendário gerenciador de arquivos XTree Gold, implementado em Go.

## Índice

1. [Visão Geral](#visão-geral)
2. [Arquitetura](ARQUITETURA.md)
3. [Interface do Usuário](UI.md)
4. [Visualizadores e Editores](VISUALIZADORES_EDITORES.md)
5. [Utilitários](UTILITARIOS.md)
6. [Ponto de Entrada](MAIN.md)
7. [Instalação e Uso](#instalação-e-uso)
8. [Contribuição](#contribuição)
9. [Licença](#licença)

## Visão Geral

O GoXTree é uma recriação moderna do lendário gerenciador de arquivos XTree Gold, implementado em Go. Este projeto visa trazer de volta a eficiência e a simplicidade do XTree Gold original, mas com tecnologias modernas e compatibilidade com sistemas operacionais atuais.

### Características

- Interface de texto com elementos gráficos usando TUI (Terminal User Interface)
- Visualização em árvore hierárquica de diretórios
- Visualizador de arquivos para múltiplos formatos (texto, hexadecimal, imagens, etc.)
- Editor de texto integrado
- Funções de busca para localizar arquivos
- Operações básicas de gerenciamento de arquivos (copiar, mover, excluir)
- Baixo consumo de recursos
- Interface intuitiva inspirada no XTree Gold original

### Inspiração

O XTree Gold foi um popular gerenciador de arquivos para DOS lançado pela primeira vez em 1989 como uma versão aprimorada do XTree original (lançado em 1985). Foi desenvolvido pela Executive Systems, Inc. (ESI) e se tornou um dos gerenciadores de arquivos mais populares da era DOS.

Principais características do XTree Gold original:
1. Interface de texto com elementos gráficos, incluindo menus em cascata
2. Visualização em árvore hierárquica de diretórios (revolucionária para a época)
3. Capacidade de visualizar arquivos em diversos formatos (texto, hexadecimal, etc.)
4. Suporte para compressão de arquivos ZIP
5. Editor de texto integrado chamado "1Word"
6. Capacidade de visualizar arquivos gráficos em vários formatos
7. Funções de busca avançadas para localizar arquivos
8. Capacidade de comparar diretórios e encontrar arquivos duplicados
9. Baixo consumo de memória (importante na era DOS com limite de 640KB)
10. Funções para recuperação de arquivos perdidos

## Instalação e Uso

### Requisitos

- Go 1.16 ou superior
- Bibliotecas:
  - github.com/gdamore/tcell/v2
  - github.com/rivo/tview

### Instalação

```bash
# Clone o repositório
git clone https://github.com/peder1981/GoXTree.git

# Entre no diretório
cd GoXTree

# Compile o projeto
go build -o goxree ./cmd/gxtree/main.go

# Execute o programa
./goxree
```

### Uso

O GoXTree apresenta uma interface dividida em dois painéis principais:

1. **Painel de Árvore (esquerdo)**: Mostra a estrutura hierárquica de diretórios
2. **Painel de Arquivos (direito)**: Mostra os arquivos no diretório selecionado

#### Navegação Básica

- Use as teclas de seta para navegar pelos diretórios e arquivos
- Pressione `Enter` para expandir/recolher diretórios ou abrir arquivos
- Pressione `Tab` para alternar entre o painel de árvore e o painel de arquivos

#### Visualização de Arquivos

- Selecione um arquivo e pressione `F3` para visualizá-lo
- O tipo de visualizador será escolhido automaticamente com base na extensão do arquivo:
  - Arquivos de texto: Visualizador de texto
  - Imagens: Visualizador de imagens em ASCII art
  - Outros arquivos: Visualizador hexadecimal

#### Edição de Arquivos

- Selecione um arquivo de texto e pressione `F4` para editá-lo
- Use `Ctrl+S` para salvar as alterações
- Use `Ctrl+Q` ou `Esc` para sair do editor

#### Operações com Arquivos

- `F5`: Copiar arquivos selecionados
- `F6`: Mover arquivos selecionados
- `F7`: Criar novo diretório
- `F8`: Excluir arquivos selecionados
- `F9`: Comprimir arquivos selecionados

### Atalhos de Teclado

| Tecla | Função |
|-------|--------|
| F1 | Ajuda |
| F2 | Menu |
| F3 | Visualizar arquivo |
| F4 | Editar arquivo |
| F5 | Copiar |
| F6 | Mover |
| F7 | Criar diretório |
| F8 | Excluir |
| F9 | Comprimir |
| F10 | Sair |
| Tab | Alternar entre painéis |
| Space | Marcar/desmarcar arquivo |
| Ctrl+A | Marcar todos |
| Ctrl+F | Buscar arquivo |
| Ctrl+G | Ir para diretório |
| Ctrl+H | Mostrar arquivos ocultos |
| Ctrl+R | Atualizar |

## Contribuição

Contribuições são bem-vindas! Se você deseja contribuir para o GoXTree, siga estas etapas:

1. Faça um fork do repositório
2. Crie uma branch para sua feature (`git checkout -b feature/nova-feature`)
3. Faça commit das suas alterações (`git commit -m 'Adiciona nova feature'`)
4. Faça push para a branch (`git push origin feature/nova-feature`)
5. Abra um Pull Request

### Diretrizes de Contribuição

- Siga o estilo de código existente
- Adicione testes para novas funcionalidades
- Atualize a documentação conforme necessário
- Certifique-se de que todos os testes passam antes de enviar um Pull Request

## Licença

Este projeto está licenciado sob a [MIT License](../LICENSE).

## Agradecimentos

Este projeto é uma homenagem ao XTree Gold original, desenvolvido pela Executive Systems, Inc. (ESI) em 1989. GoXTree não é afiliado à ESI, Central Point Software, ou Symantec.

# Arquitetura do GoXTree

## Visão Geral

O GoXTree é uma recriação moderna do gerenciador de arquivos XTree Gold, implementado em Go. A arquitetura do projeto segue uma estrutura modular com componentes bem definidos e responsabilidades separadas.

## Estrutura de Diretórios

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

## Componentes Principais

### 1. Interface do Usuário (pkg/ui)

#### App (app.go)
O componente central que gerencia toda a aplicação. Ele é responsável por inicializar todos os outros componentes, gerenciar o layout da interface, processar eventos de teclado e coordenar a interação entre os diferentes componentes.

**Responsabilidades:**
- Inicialização da aplicação
- Gerenciamento do layout da interface
- Processamento de eventos de teclado
- Coordenação da interação entre componentes
- Gerenciamento de diálogos e menus

#### TreeView (treeview.go)
Componente responsável pela visualização em árvore de diretórios. Ele exibe a estrutura hierárquica de diretórios e permite a navegação entre eles.

**Responsabilidades:**
- Exibição da estrutura hierárquica de diretórios
- Navegação entre diretórios
- Expansão e colapso de diretórios
- Seleção de diretórios

#### FileView (fileview.go)
Componente responsável pela visualização de arquivos em formato de tabela. Ele exibe os arquivos no diretório selecionado, com informações como nome, tamanho, data de modificação, etc.

**Responsabilidades:**
- Exibição de arquivos em formato de tabela
- Seleção de arquivos
- Marcação de arquivos para operações em lote
- Ordenação de arquivos por diferentes critérios

#### MenuBar (menubar.go)
Componente responsável pela barra de menu da aplicação. Ele exibe os menus disponíveis e processa as ações selecionadas.

**Responsabilidades:**
- Exibição de menus
- Processamento de ações de menu
- Exibição de atalhos de teclado

#### StatusBar (statusbar.go)
Componente responsável pela barra de status da aplicação. Ele exibe informações sobre o estado atual da aplicação, como diretório atual, número de arquivos, espaço livre, etc.

**Responsabilidades:**
- Exibição de informações sobre o estado atual
- Atualização dinâmica de informações

#### HelpView (helpview.go)
Componente responsável pela tela de ajuda da aplicação. Ele exibe informações sobre como usar a aplicação, atalhos de teclado, etc.

**Responsabilidades:**
- Exibição de informações de ajuda
- Navegação na tela de ajuda

#### FileViewer (fileviewer.go)
Interface para visualização de arquivos. Ele coordena a seleção do visualizador apropriado com base no tipo de arquivo.

**Responsabilidades:**
- Seleção do visualizador apropriado
- Coordenação da visualização de arquivos

#### FileEditor (fileeditor.go)
Interface para edição de arquivos. Ele coordena a seleção do editor apropriado com base no tipo de arquivo.

**Responsabilidades:**
- Seleção do editor apropriado
- Coordenação da edição de arquivos

### 2. Visualizadores (pkg/viewer)

#### TextViewer (textviewer.go)
Componente responsável pela visualização de arquivos de texto. Ele exibe o conteúdo de arquivos de texto em formato legível.

**Responsabilidades:**
- Exibição de arquivos de texto
- Navegação no conteúdo do arquivo
- Busca de texto

#### ImageViewer (imageviewer.go)
Componente responsável pela visualização de imagens. Ele exibe imagens em formato ASCII art no terminal.

**Responsabilidades:**
- Exibição de imagens em ASCII art
- Suporte para diferentes formatos de imagem

### 3. Editores (pkg/editor)

#### TextEditor (texteditor.go)
Componente responsável pela edição de arquivos de texto. Ele permite a edição de arquivos de texto com funcionalidades básicas.

**Responsabilidades:**
- Edição de arquivos de texto
- Salvamento de alterações
- Busca e substituição de texto

### 4. Utilitários (pkg/utils)

#### FileUtils (fileutils.go)
Funções utilitárias para manipulação de arquivos. Ele fornece funções para operações comuns com arquivos, como cópia, movimentação, exclusão, etc.

**Responsabilidades:**
- Operações com arquivos (cópia, movimentação, exclusão)
- Obtenção de informações sobre arquivos
- Verificação de permissões

#### Constants (constants.go)
Constantes utilizadas no projeto. Ele define constantes como teclas de atalho, cores, etc.

**Responsabilidades:**
- Definição de constantes utilizadas no projeto

## Fluxo de Dados

1. O usuário interage com a interface através de eventos de teclado
2. O App processa esses eventos e os encaminha para o componente apropriado
3. O componente executa a ação solicitada e atualiza seu estado
4. O App atualiza a interface para refletir o novo estado

## Padrões de Design

### Padrão MVC (Model-View-Controller)
O GoXTree segue uma variação do padrão MVC:
- **Model**: Representado pelos utilitários e pelo sistema de arquivos
- **View**: Representado pelos componentes de interface do usuário
- **Controller**: Representado pelo App e pelos controladores específicos de cada componente

### Padrão Observer
Os componentes da interface do usuário observam mudanças no estado da aplicação e se atualizam automaticamente.

### Padrão Strategy
Os visualizadores e editores de arquivos são selecionados dinamicamente com base no tipo de arquivo, seguindo o padrão Strategy.

## Considerações de Desempenho

- O GoXTree é projetado para ser leve e eficiente, com baixo consumo de recursos
- A visualização de diretórios é carregada sob demanda, para evitar o carregamento de toda a árvore de diretórios de uma vez
- A visualização de arquivos grandes é otimizada para evitar o carregamento completo do arquivo na memória

## Extensibilidade

O GoXTree é projetado para ser facilmente extensível:
- Novos visualizadores podem ser adicionados implementando a interface apropriada
- Novos editores podem ser adicionados implementando a interface apropriada
- Novos comandos podem ser adicionados ao menu e ao processamento de eventos de teclado

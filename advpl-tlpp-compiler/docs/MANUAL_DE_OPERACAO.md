# Manual de Operação - Compilador AdvPL/TLPP

## Índice
1. [Introdução](#introdução)
2. [Instalação](#instalação)
3. [Estrutura do Projeto](#estrutura-do-projeto)
4. [Uso do Compilador](#uso-do-compilador)
5. [Servidor LSP](#servidor-lsp)
6. [Integração com IDEs](#integração-com-ides)
7. [Troubleshooting](#troubleshooting)
8. [Referências](#referências)

## Introdução

O Compilador AdvPL/TLPP é uma ferramenta moderna para desenvolvimento em AdvPL e TLPP, linguagens utilizadas na plataforma TOTVS Protheus. O compilador oferece:

- Análise sintática e semântica
- Geração de código otimizado
- Integração com IDEs via LSP
- Suporte a orientação a objetos
- Diagnósticos em tempo real

## Instalação

### Pré-requisitos
- Go 1.22 ou superior
- Git

### Passos de Instalação

1. Clone o repositório:
```bash
git clone https://github.com/peder1981/advpl-tlpp-compiler
cd advpl-tlpp-compiler
```

2. Compile o projeto:
```bash
go build -o bin/advpl-compiler ./cmd/compiler
go build -o bin/advpl-lsp ./cmd/lsp
```

3. Adicione ao PATH:
```bash
export PATH=$PATH:/caminho/para/advpl-tlpp-compiler/bin
```

## Estrutura do Projeto

```
advpl-tlpp-compiler/
├── cmd/
│   ├── compiler/     # Ponto de entrada do compilador
│   └── lsp/         # Servidor LSP
├── pkg/
│   ├── ast/         # Árvore sintática abstrata
│   ├── compiler/    # Núcleo do compilador
│   ├── lexer/       # Análise léxica
│   ├── parser/      # Análise sintática
│   └── lsp/         # Implementação LSP
├── docs/            # Documentação
└── examples/        # Exemplos de código
```

## Uso do Compilador

### Compilação Básica
```bash
advpl-compiler arquivo.prw
```

### Opções do Compilador
- `-v`: Modo verbose
- `-o`: Otimização de código
- `-d`: Dialeto (advpl ou tlpp)
- `-i`: Diretórios de include
- `-s`: Apenas verificação de sintaxe
- `-doc`: Gera documentação

### Exemplos de Uso

1. Compilar um arquivo:
```bash
advpl-compiler fonte.prw
```

2. Compilar com otimização:
```bash
advpl-compiler -o fonte.prw
```

3. Verificar sintaxe:
```bash
advpl-compiler -s fonte.prw
```

## Servidor LSP

O servidor LSP fornece recursos avançados de IDE para editores compatíveis.

### Iniciar o Servidor
```bash
advpl-lsp
```

### Recursos Disponíveis
- Realce de sintaxe
- Completação de código
- Navegação (ir para definição)
- Hover com informações
- Formatação de código
- Diagnósticos em tempo real

### Configuração do Log
O servidor gera logs em `advpl-lsp.log` para debug.

## Integração com IDEs

### VSCode

1. Instale a extensão "AdvPL/TLPP Language Support"

2. Configure o servidor LSP:
```json
{
  "advpl.lspPath": "/caminho/para/advpl-lsp",
  "advpl.dialect": "advpl",
  "advpl.maxNumberOfProblems": 100
}
```

### Neovim

1. Configure o LSP no init.lua:
```lua
require'lspconfig'.advpl.setup{
  cmd = { "advpl-lsp" },
  filetypes = { "advpl", "tlpp", "prw" },
  root_dir = function(fname)
    return require'lspconfig'.util.find_git_ancestor(fname)
  end,
}
```

### Emacs

1. Configure o eglot:
```elisp
(add-to-list 'eglot-server-programs
             '((advpl-mode tlpp-mode)
               . ("advpl-lsp")))
```

## Troubleshooting

### Problemas Comuns

1. **Erro de Compilação**
   - Verifique a sintaxe do código
   - Confirme se todos os includes estão acessíveis
   - Verifique o log de erros

2. **LSP Não Conecta**
   - Verifique se o servidor está no PATH
   - Confirme as configurações do editor
   - Verifique o arquivo de log

3. **Completação Não Funciona**
   - Verifique se o documento foi salvo
   - Confirme se o servidor LSP está rodando
   - Verifique a sintaxe do código

### Logs e Diagnóstico

- Logs do compilador: `stderr`
- Logs do LSP: `advpl-lsp.log`
- Diagnósticos: Exibidos no editor

## Referências

1. [Documentação AdvPL](https://tdn.totvs.com/display/tec/AdvPL)
2. [Language Server Protocol](https://microsoft.github.io/language-server-protocol/)
3. [TOTVS Developer Network](https://tdn.totvs.com/)
4. [Protheus Documentation](https://tdn.totvs.com/display/tec/Protheus)

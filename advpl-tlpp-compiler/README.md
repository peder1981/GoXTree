# AdvPL/TLPP Compiler e IDE

Um compilador e ambiente de desenvolvimento integrado (IDE) para as linguagens AdvPL e TLPP, implementado em Go.

## Sobre o Projeto

Este projeto consiste em dois componentes principais:

1. **Compilador AdvPL/TLPP**: Um compilador completo para as linguagens AdvPL (Advanced Protheus Language) e TLPP (Totvs Language for Protheus and Protheus), utilizadas no desenvolvimento de aplicações para o ERP Protheus da TOTVS.

2. **IDE baseado em ASCII**: Um ambiente de desenvolvimento integrado baseado em terminal, que oferece funcionalidades como destaque de sintaxe, completação de código, navegação de código, depuração e integração com o compilador.

## Características

### Compilador
- Análise léxica e sintática completa para AdvPL e TLPP
- Verificação de tipos e semântica
- Otimização de código
- Geração de código objeto compatível com o Protheus
- Suporte a diferentes dialetos e versões das linguagens

### IDE
- Interface baseada em terminal (ASCII)
- Destaque de sintaxe para AdvPL e TLPP
- Completação de código inteligente
- Navegação de código (definições, referências)
- Depuração integrada
- Gerenciamento de projetos
- Integração com o compilador
- Suporte a temas e personalização

## Requisitos

- Go 1.22 ou superior
- Terminal com suporte a cores e caracteres Unicode

## Instalação

```bash
go install github.com/peder1981/advpl-tlpp-compiler/cmd/compiler@latest
go install github.com/peder1981/advpl-tlpp-compiler/cmd/ide@latest
```

## Uso

### Compilador

```bash
advpl-compiler arquivo.prw [opções]
```

### IDE

```bash
advpl-ide [diretório_do_projeto]
```

## Estrutura do Projeto

```
advpl-tlpp-compiler/
├── cmd/
│   ├── compiler/     # Ponto de entrada do compilador
│   └── ide/          # Ponto de entrada do IDE
├── pkg/
│   ├── ast/          # Estruturas de árvore sintática abstrata
│   ├── compiler/     # Lógica de compilação
│   ├── ide/          # Componentes do IDE
│   ├── lexer/        # Analisador léxico
│   ├── parser/       # Analisador sintático
│   └── utils/        # Utilitários compartilhados
└── examples/         # Exemplos de código AdvPL e TLPP
```

## Licença

Este projeto está licenciado sob a [MIT License](LICENSE).

## Contribuição

Contribuições são bem-vindas! Por favor, sinta-se à vontade para abrir issues ou enviar pull requests.

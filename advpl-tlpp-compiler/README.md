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
- Execução de código compilado (simulação ou integração com Protheus)
- Geração automática de documentação

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
- Opcional: Protheus AppServer para execução de código real

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

Opções disponíveis:
- `-o arquivo.ppo`: Especifica o arquivo de saída
- `-v`: Modo verboso
- `-optimize`: Otimiza o código gerado
- `-dialect [advpl|tlpp]`: Especifica o dialeto da linguagem
- `-I diretório`: Adiciona um diretório de inclusão
- `-check`: Apenas verifica a sintaxe sem gerar código
- `-docs`: Gera documentação a partir dos comentários
- `-run`: Executa o código após compilação
- `-server servidor`: Especifica o servidor Protheus para execução
- `-env ambiente`: Especifica o ambiente Protheus para execução
- `-keep-temp`: Mantém arquivos temporários após execução

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
│   ├── compiler/     # Lógica de compilação e geração de código
│   ├── executor/     # Executor de código compilado
│   ├── ide/          # Componentes do IDE
│   ├── lexer/        # Analisador léxico
│   ├── parser/       # Analisador sintático
│   └── utils/        # Utilitários compartilhados
└── examples/         # Exemplos de código AdvPL e TLPP
```

## Pipeline de Compilação

O processo de compilação segue as seguintes etapas:

1. **Análise Léxica**: O código fonte é transformado em tokens pelo lexer.
2. **Análise Sintática**: Os tokens são analisados pelo parser para criar uma Árvore Sintática Abstrata (AST).
3. **Análise Semântica**: A AST é verificada quanto a erros semânticos e tipos.
4. **Geração de Código**: A AST é transformada em código objeto pelo gerador de código.
5. **Otimização** (opcional): O código gerado é otimizado para melhor desempenho.
6. **Execução** (opcional): O código compilado pode ser executado diretamente.

## Exemplos

O diretório `examples/` contém exemplos de código AdvPL e TLPP que podem ser usados para testar o compilador:

- `cliente.prw`: Implementação de uma classe de cliente com métodos para gerenciamento.
- Outros exemplos serão adicionados no futuro.

## Licença

Este projeto está licenciado sob a [MIT License](LICENSE).

## Contribuição

Contribuições são bem-vindas! Por favor, sinta-se à vontade para abrir issues ou enviar pull requests.

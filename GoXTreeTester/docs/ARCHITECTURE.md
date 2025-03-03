# Arquitetura do GoXTreeTester

Este documento descreve a arquitetura e o funcionamento interno do GoXTreeTester, uma ferramenta de análise e teste automático para o projeto GoXTree.

## Visão Geral

O GoXTreeTester é composto por quatro módulos principais:

1. **Analyzer**: Responsável pela análise estática do código e verificação de estilo
2. **Fixer**: Responsável pela correção automática de problemas
3. **Tester**: Responsável pela execução de testes
4. **Reporter**: Responsável pela geração de relatórios

## Módulo Analyzer

O módulo Analyzer é responsável por analisar o código-fonte do GoXTree e identificar problemas. Ele é composto por vários componentes:

### Analyzer

O componente principal que coordena a análise. Ele percorre o código-fonte do GoXTree e executa várias verificações:

- **Verificação de estrutura do projeto**: Verifica se a estrutura do projeto segue o padrão esperado.
- **Verificação de imports**: Verifica se há imports não utilizados ou faltando.
- **Verificação de tipos e métodos**: Verifica se há chamadas para métodos não definidos.
- **Verificação de estilo de código**: Verifica se o código segue as boas práticas de estilo.

### StyleChecker

O StyleChecker é responsável por verificar o estilo do código. Ele verifica:

- **Linhas muito longas**: Identifica linhas com mais de 100 caracteres.
- **Espaços em branco no final das linhas**: Identifica linhas com espaços em branco no final.
- **Comentários mal formatados**: Verifica se os comentários seguem o padrão do Go.
- **Documentação de tipos e funções exportados**: Verifica se tipos e funções exportados têm comentários de documentação.

## Módulo Fixer

O módulo Fixer é responsável por corrigir automaticamente os problemas identificados pelo Analyzer. Ele implementa várias estratégias de correção:

### FixIssues

Corrige vários tipos de problemas, como:

- **Imports não utilizados**: Remove imports que não são utilizados no código.
- **Problemas de estilo**: Corrige problemas de estilo, como espaços em branco no final das linhas e comentários mal formatados.

### FixUnusedImports e FixUnusedImportsInProject

Funções específicas para corrigir imports não utilizados. A função `FixUnusedImportsInProject` é especialmente importante, pois:

1. Verifica quais arquivos usam diretamente o pacote tcell
2. Garante que esses arquivos importem o pacote tcell, mesmo que o import pareça não ser utilizado
3. Remove imports não utilizados de outros arquivos

Isso é necessário porque o pacote tcell é usado de forma indireta em muitos arquivos, através de tipos como `tcell.Key`.

## Módulo Tester

O módulo Tester é responsável por executar testes no GoXTree. Ele implementa vários tipos de testes:

### Testes Unitários

Executa os testes unitários definidos no projeto GoXTree.

### Testes Funcionais

Executa testes funcionais que verificam o comportamento do GoXTree como um todo. Inclui testes de:

- **Inicialização da aplicação**: Verifica se a aplicação inicializa corretamente.
- **Navegação de diretórios**: Verifica se a navegação de diretórios funciona corretamente.
- **Operações de arquivo**: Verifica se operações como copiar, mover e excluir arquivos funcionam corretamente.
- **Componentes de UI**: Verifica se os componentes de interface do usuário funcionam corretamente.

## Módulo Reporter

O módulo Reporter é responsável por gerar relatórios sobre a análise e os testes. Ele coleta informações durante a execução e gera um relatório HTML detalhado.

### Reporter

O componente principal que coleta informações e gera o relatório. Ele mantém listas de:

- **Issues**: Problemas identificados durante a análise.
- **TestResults**: Resultados dos testes executados.
- **Fixes**: Correções aplicadas automaticamente.
- **Warnings**: Avisos gerados durante a execução.

### GenerateReport

Gera um relatório HTML detalhado com todas as informações coletadas. O relatório inclui:

- **Resumo**: Um resumo dos problemas, testes e correções.
- **Problemas**: Lista detalhada de problemas identificados.
- **Resultados dos Testes**: Lista detalhada de resultados dos testes.
- **Correções**: Lista de correções aplicadas automaticamente.
- **Avisos**: Lista de avisos gerados durante a execução.

## Fluxo de Execução

O fluxo de execução típico do GoXTreeTester é:

1. O usuário executa o comando `gxtester` com as opções desejadas.
2. O programa inicializa os componentes necessários (Analyzer, Fixer, Tester, Reporter).
3. Se a opção `--fix-imports` for especificada, o programa corrige apenas os imports não utilizados.
4. Se a opção `--check-style` for especificada, o programa verifica apenas o estilo do código.
5. Se nenhuma dessas opções for especificada, o programa executa a análise completa e os testes.
6. Se a opção `--autofix` for especificada, o programa corrige automaticamente os problemas identificados.
7. O programa gera um relatório HTML com os resultados.

## Conclusão

O GoXTreeTester é uma ferramenta poderosa para análise e teste do GoXTree. Ele ajuda a manter a qualidade do código, identificando problemas e propondo correções automaticamente.

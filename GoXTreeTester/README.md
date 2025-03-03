# GoXTreeTester

GoXTreeTester é uma ferramenta de análise e teste automático para o projeto GoXTree. Ela identifica problemas no código, executa testes e propõe correções automaticamente.

## Funcionalidades

- **Análise Estática de Código**: Identifica problemas como funções duplicadas, imports não utilizados, métodos não definidos e inconsistências de nomenclatura.
- **Execução de Testes**: Executa testes unitários e funcionais automaticamente.
- **Correção Automática**: Corrige problemas identificados automaticamente quando possível.
- **Geração de Relatórios**: Gera relatórios detalhados em formato HTML com os resultados da análise e dos testes.

## Requisitos

- Go 1.22 ou superior
- Acesso ao código-fonte do GoXTree

## Instalação

```bash
# Clonar o repositório
git clone https://github.com/peder1981/GoXTreeTester.git

# Entrar no diretório
cd GoXTreeTester

# Compilar
go build -o gxtester ./cmd/gxtester
```

## Uso

```bash
# Executar análise e testes
./gxtester --project=/caminho/para/GoXTree

# Executar análise e testes com correção automática
./gxtester --project=/caminho/para/GoXTree --autofix

# Executar apenas análise
./gxtester --project=/caminho/para/GoXTree --analyze-only

# Executar apenas testes
./gxtester --project=/caminho/para/GoXTree --test-only

# Exibir informações detalhadas
./gxtester --project=/caminho/para/GoXTree --verbose

# Ignorar erros específicos
./gxtester --project=/caminho/para/GoXTree --ignore="Import não utilizado"

# Especificar caminho para o relatório
./gxtester --project=/caminho/para/GoXTree --report=./relatorio.html
```

## Estrutura do Projeto

- **cmd/gxtester**: Ponto de entrada da aplicação
- **pkg/analyzer**: Módulo de análise estática de código
- **pkg/tester**: Módulo de execução de testes
- **pkg/fixer**: Módulo de correção automática
- **pkg/reporter**: Módulo de geração de relatórios

## Como Funciona

1. **Análise**: O módulo de análise percorre o código-fonte do GoXTree e identifica problemas como funções duplicadas, imports não utilizados, métodos não definidos e inconsistências de nomenclatura.

2. **Testes**: O módulo de testes executa testes unitários e funcionais para verificar se o código está funcionando corretamente.

3. **Correção**: O módulo de correção aplica correções automáticas para os problemas identificados, quando possível.

4. **Relatório**: O módulo de relatório gera um relatório detalhado em formato HTML com os resultados da análise e dos testes.

## Exemplos de Problemas Detectados

- Funções duplicadas
- Imports não utilizados
- Chamadas para métodos não definidos
- Inconsistências de nomenclatura

## Exemplos de Correções Automáticas

- Remoção de funções duplicadas
- Remoção de imports não utilizados
- Correção de nomes de funções
- Substituição de chamadas para métodos não definidos por métodos similares

## Licença

Este projeto está licenciado sob a licença MIT - veja o arquivo LICENSE para mais detalhes.

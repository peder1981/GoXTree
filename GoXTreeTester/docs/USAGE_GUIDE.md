# Guia de Uso do GoXTreeTester

Este guia explica como usar o GoXTreeTester para analisar, testar e melhorar o código do projeto GoXTree.

## Instalação

Antes de usar o GoXTreeTester, você precisa instalá-lo:

```bash
# Clonar o repositório
git clone https://github.com/peder1981/GoXTreeTester.git

# Entrar no diretório
cd GoXTreeTester

# Compilar
go build -o gxtester ./cmd/gxtester
```

## Uso Básico

O uso mais básico do GoXTreeTester é executá-lo sem opções adicionais, especificando apenas o caminho para o projeto GoXTree:

```bash
./gxtester --project=/caminho/para/GoXTree
```

Isso executará a análise completa, os testes e gerará um relatório HTML.

## Opções Disponíveis

O GoXTreeTester oferece várias opções para personalizar sua execução:

### Opções de Análise

- `--analyze-only`: Executa apenas a análise, sem os testes.
- `--check-style`: Verifica apenas o estilo do código.
- `--ignore="Padrão"`: Ignora problemas que correspondem ao padrão especificado.

### Opções de Teste

- `--test-only`: Executa apenas os testes, sem a análise.
- `--test-timeout=30s`: Define o tempo limite para a execução dos testes.

### Opções de Correção

- `--autofix`: Corrige automaticamente os problemas identificados.
- `--fix-imports`: Corrige apenas os imports não utilizados.

### Opções de Relatório

- `--report=./relatorio.html`: Especifica o caminho para o relatório HTML.
- `--verbose`: Exibe informações detalhadas durante a execução.

## Exemplos de Uso

### Análise Completa com Correção Automática

```bash
./gxtester --project=/caminho/para/GoXTree --autofix
```

Isso executará a análise completa, os testes e corrigirá automaticamente os problemas identificados.

### Verificação de Estilo

```bash
./gxtester --project=/caminho/para/GoXTree --check-style
```

Isso verificará apenas o estilo do código, sem executar os testes ou corrigir problemas.

### Verificação de Estilo com Correção Automática

```bash
./gxtester --project=/caminho/para/GoXTree --check-style --autofix
```

Isso verificará o estilo do código e corrigirá automaticamente os problemas identificados.

### Correção de Imports Não Utilizados

```bash
./gxtester --project=/caminho/para/GoXTree --fix-imports
```

Isso corrigirá apenas os imports não utilizados, sem executar a análise completa ou os testes.

### Execução de Testes

```bash
./gxtester --project=/caminho/para/GoXTree --test-only
```

Isso executará apenas os testes, sem a análise ou correção.

### Execução Detalhada

```bash
./gxtester --project=/caminho/para/GoXTree --verbose
```

Isso executará a análise completa e os testes, exibindo informações detalhadas durante a execução.

## Interpretando o Relatório

O GoXTreeTester gera um relatório HTML detalhado com os resultados da análise e dos testes. O relatório é dividido em várias seções:

### Resumo

O resumo mostra uma visão geral dos problemas, testes e correções. Ele inclui:

- Número total de problemas identificados
- Número de testes executados e seus resultados
- Número de correções aplicadas automaticamente
- Número de avisos gerados durante a execução

### Problemas

A seção de problemas lista todos os problemas identificados durante a análise. Para cada problema, o relatório mostra:

- Tipo de problema
- Arquivo e linha onde o problema foi encontrado
- Descrição detalhada do problema
- Sugestão de correção, quando disponível

### Resultados dos Testes

A seção de resultados dos testes lista todos os testes executados. Para cada teste, o relatório mostra:

- Nome do teste
- Resultado (sucesso ou falha)
- Tempo de execução
- Mensagem de erro, quando aplicável

### Correções

A seção de correções lista todas as correções aplicadas automaticamente. Para cada correção, o relatório mostra:

- Tipo de correção
- Arquivo e linha onde a correção foi aplicada
- Descrição detalhada da correção
- Código antes e depois da correção

### Avisos

A seção de avisos lista todos os avisos gerados durante a execução. Para cada aviso, o relatório mostra:

- Tipo de aviso
- Arquivo e linha onde o aviso foi gerado
- Descrição detalhada do aviso

## Dicas e Boas Práticas

### Executar Regularmente

Execute o GoXTreeTester regularmente para manter a qualidade do código. Idealmente, execute-o antes de cada commit.

### Corrigir Problemas Manualmente

Nem todos os problemas podem ser corrigidos automaticamente. Para problemas que requerem correção manual, o relatório fornece sugestões detalhadas.

### Revisar Correções Automáticas

Sempre revise as correções automáticas antes de confirmar as alterações. Embora o GoXTreeTester seja cuidadoso, é sempre bom verificar se as correções não introduziram novos problemas.

### Usar com Integração Contínua

Você pode integrar o GoXTreeTester ao seu pipeline de integração contínua para garantir que todos os commits mantenham a qualidade do código.

## Solução de Problemas

### O GoXTreeTester não encontra o projeto GoXTree

Verifique se o caminho para o projeto GoXTree está correto e se você tem permissão para acessá-lo.

### O GoXTreeTester não corrige alguns problemas

Nem todos os problemas podem ser corrigidos automaticamente. Alguns problemas requerem intervenção manual.

### O GoXTreeTester está muito lento

Se o GoXTreeTester estiver muito lento, tente executá-lo com opções mais específicas, como `--check-style` ou `--fix-imports`, em vez de executar a análise completa.

### O relatório HTML não é gerado

Verifique se você tem permissão para escrever no diretório onde o relatório deve ser gerado. Se não especificar um caminho para o relatório, ele será gerado no diretório atual.

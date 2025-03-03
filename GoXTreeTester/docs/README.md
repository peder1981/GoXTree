# Documentação do GoXTreeTester

Bem-vindo à documentação do GoXTreeTester, uma ferramenta de análise e teste automático para o projeto GoXTree.

## Índice

1. [Guia de Uso](USAGE_GUIDE.md) - Como usar o GoXTreeTester para analisar, testar e melhorar o código do GoXTree.
2. [Arquitetura](ARCHITECTURE.md) - Descrição detalhada da arquitetura e funcionamento interno do GoXTreeTester.
3. [Guia de Desenvolvimento](DEVELOPMENT_GUIDE.md) - Guia para desenvolvedores que desejam contribuir com o GoXTreeTester.

## Visão Geral

O GoXTreeTester é uma ferramenta de análise e teste automático para o projeto GoXTree. Ele identifica problemas no código, executa testes e propõe correções automaticamente.

### Funcionalidades Principais

- **Análise Estática de Código**: Identifica problemas como funções duplicadas, imports não utilizados, métodos não definidos e inconsistências de nomenclatura.
- **Verificação de Estilo de Código**: Analisa o estilo do código, identificando problemas como linhas muito longas, espaços em branco no final das linhas e falta de documentação.
- **Execução de Testes**: Executa testes unitários e funcionais automaticamente.
- **Correção Automática**: Corrige problemas identificados automaticamente quando possível.
- **Geração de Relatórios**: Gera relatórios detalhados em formato HTML com os resultados da análise e dos testes.

## Início Rápido

Para começar a usar o GoXTreeTester, siga estes passos:

1. Instale o GoXTreeTester:

```bash
git clone https://github.com/peder1981/GoXTreeTester.git
cd GoXTreeTester
go build -o gxtester ./cmd/gxtester
```

2. Execute o GoXTreeTester:

```bash
./gxtester --project=/caminho/para/GoXTree
```

3. Veja o relatório gerado e corrija os problemas identificados.

Para mais informações, consulte o [Guia de Uso](USAGE_GUIDE.md).

## Contribuindo

Se você deseja contribuir com o GoXTreeTester, consulte o [Guia de Desenvolvimento](DEVELOPMENT_GUIDE.md) para obter informações sobre como configurar seu ambiente de desenvolvimento, entender o código existente e enviar suas alterações.

## Licença

Este projeto está licenciado sob a licença MIT - veja o arquivo LICENSE para mais detalhes.

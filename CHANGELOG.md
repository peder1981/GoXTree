# Histórico de Alterações do GoXTree

## [1.1.0] - 2025-03-02

### Adicionado
- Barra de função dedicada na parte inferior da tela
- Suporte a compilação para múltiplas plataformas
- Script de compilação automatizado (build_all.sh)
- Documentação detalhada para desenvolvedores (DEVELOPERS.md)
- Verificação para evitar duplicação de entradas no histórico de navegação

### Alterado
- Separação do comportamento das teclas F10 e ESC
  - F10: Sempre pergunta se o usuário deseja sair da aplicação
  - ESC: Comportamento universal (voltar ao diretório anterior ou perguntar se deseja sair)
- Simplificação da barra de status para mostrar apenas informações essenciais
- Reorganização do layout principal para vertical (FlexRow)
- Melhoria na navegação de histórico com tratamento de casos de borda
- Atualização do README.md com informações detalhadas

### Corrigido
- Problema com a exibição da barra de status
- Inconsistência no comportamento das teclas F10 e ESC
- Problemas na navegação entre diretórios
- Tratamento incorreto do histórico de navegação

## [1.0.0] - 2025-02-15

### Adicionado
- Interface de texto com elementos gráficos usando TUI
- Visualização em árvore hierárquica de diretórios
- Visualizador de arquivos para múltiplos formatos
- Editor de texto integrado
- Funções de busca para localizar arquivos
- Operações básicas de gerenciamento de arquivos
- Comparação de arquivos usando a biblioteca go-diff/diffmatchpatch
- Navegação de diretórios com suporte a histórico
- Opção ".." para navegar para o diretório pai

### Alterado
- Padronização dos métodos de diálogo
- Melhoria na exibição de informações de arquivos
- Padronização da estrutura de métodos
- Padronização da nomenclatura de métodos e variáveis

### Corrigido
- Problemas com visualização de arquivos e diretórios
- Implementação incorreta dos manipuladores de eventos de teclado
- Problemas de tipos e conversões
- Imports não utilizados

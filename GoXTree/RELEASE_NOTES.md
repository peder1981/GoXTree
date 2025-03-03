# GoXTree v1.0.0 - Notas de Lançamento

É com grande satisfação que anunciamos o lançamento oficial do GoXTree v1.0.0, um gerenciador de arquivos moderno inspirado no lendário XTree Gold, agora reimaginado e implementado em Go!

## Sobre o GoXTree

O GoXTree é uma homenagem ao clássico XTree Gold, trazendo sua eficiência e simplicidade para os sistemas operacionais modernos. Desenvolvido inteiramente em Go, o GoXTree oferece uma interface de terminal intuitiva e poderosa para gerenciar seus arquivos com agilidade e precisão.

## Principais Características

- **Interface TUI Elegante**: Interface de texto com elementos gráficos coloridos usando a biblioteca tview
- **Visualização em Árvore**: Navegue facilmente pela estrutura hierárquica de diretórios
- **Operações de Arquivo Completas**: Copie, mova, renomeie e exclua arquivos com facilidade
- **Visualizador Integrado**: Visualize arquivos de texto, binários e imagens diretamente no aplicativo
- **Comparação de Arquivos**: Compare o conteúdo de dois arquivos e veja as diferenças destacadas
- **Navegação Eficiente**: Use atalhos de teclado e teclas de função para uma navegação rápida
- **Multiplataforma**: Funciona em Linux, Windows e macOS (Intel e Apple Silicon)
- **Baixo Consumo de Recursos**: Leve e eficiente, ideal para sistemas com recursos limitados

## Novidades nesta Versão

- Barra de função dedicada para fácil acesso às operações comuns
- Comportamento aprimorado das teclas F10 e ESC para navegação intuitiva
- Melhor exibição de informações na barra de status
- Navegação aprimorada com suporte a histórico de diretórios
- Opção ".." para navegar facilmente para o diretório pai
- Compilação otimizada para múltiplas plataformas e arquiteturas

## Compatibilidade

O GoXTree é compatível com:
- Linux (AMD64, ARM64)
- Windows (AMD64)
- macOS/Darwin (AMD64, ARM64)

## Instalação

Baixe o binário pré-compilado para o seu sistema ou compile a partir do código fonte:

```bash
# Clone o repositório
git clone https://github.com/peder1981/GoXTree.git

# Entre no diretório
cd GoXTree

# Compile o projeto
go build -o bin/goxTree ./cmd/gxtree/main.go

# Execute o programa
./bin/goxTree
```

## Agradecimentos

Agradecemos a todos que contribuíram para tornar este projeto possível. O GoXTree é um projeto de código aberto e todas as contribuições são bem-vindas!

## O Futuro do GoXTree

Estamos comprometidos com o desenvolvimento contínuo do GoXTree, com planos para adicionar:
- Suporte a compressão/descompressão de arquivos
- Integração com sistemas de controle de versão
- Personalização avançada de temas e cores
- Suporte a plugins
- E muito mais!

Junte-se a nós nesta jornada para criar o melhor gerenciador de arquivos baseado em terminal disponível hoje!

---

*"Simplicidade é a sofisticação suprema." - Leonardo da Vinci*

# GoXTree - Gerenciador de Arquivos Moderno

GoXTree é um gerenciador de arquivos moderno inspirado no clássico XTree Gold, implementado em Go com uma interface de terminal (TUI).

## Características

- Interface simples e intuitiva
- Navegação rápida entre diretórios
- Visualização de arquivos
- Operações básicas de arquivos (copiar, mover, excluir)
- Suporte a múltiplas plataformas (Linux, Windows, macOS)
- Baixo consumo de recursos
- Inspirado na interface clássica do XTree Gold

## Instalação

Baixe a versão mais recente do GoXTree na [página de releases](https://github.com/peder1981/GoXTree/releases).

### Linux e macOS

```bash
# Extrair o arquivo
tar -xzf goxTree-linux-amd64.tar.gz

# Tornar executável
chmod +x goxTree-linux-amd64

# Mover para um diretório no PATH (opcional)
sudo mv goxTree-linux-amd64 /usr/local/bin/goxTree
```

### Windows

1. Baixe o arquivo `goxTree-windows-amd64.exe`
2. Renomeie para `goxTree.exe` (opcional)
3. Execute o arquivo

## Uso

```bash
# Iniciar no diretório atual
goxTree

# Iniciar em um diretório específico
goxTree /caminho/para/diretorio
```

## Atalhos de Teclado

| Tecla | Função |
|-------|--------|
| Setas | Navegar na lista de arquivos |
| Enter | Entrar no diretório / Abrir arquivo |
| ESC | Voltar ao diretório anterior / Sair |
| F1 | Ajuda |
| F2 | Renomear arquivo/diretório |
| F7 | Criar diretório |
| F8 | Excluir arquivo/diretório |
| F10 | Sair |
| Tab | Alternar entre painéis |

## Contribuindo

Contribuições são bem-vindas! Veja [CONTRIBUTING.md](https://github.com/peder1981/GoXTree/blob/main/.github/CONTRIBUTING.md) para mais detalhes.

## Licença

Este projeto está licenciado sob a licença MIT - veja o arquivo [LICENSE](https://github.com/peder1981/GoXTree/blob/main/LICENSE) para detalhes.

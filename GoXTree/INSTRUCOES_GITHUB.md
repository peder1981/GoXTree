# Instruções para Publicação no GitHub

Para publicar o GoXTree no GitHub e criar uma release, siga estas etapas:

## 1. Criar o Repositório no GitHub

1. Acesse https://github.com/new
2. Nome do repositório: GoXTree
3. Descrição: Gerenciador de arquivos moderno inspirado no XTree Gold, implementado em Go
4. Visibilidade: Pública
5. Inicialize com README: Não
6. Clique em "Criar repositório"

## 2. Enviar o Código para o GitHub

```bash
# No diretório do projeto
cd /home/peder/my-advpl-project/GoXTree

# Configurar o remote (já feito)
# git remote add origin https://github.com/peder1981/GoXTree.git

# Enviar o código
git push -u origin main

# Enviar a tag
git push origin v1.0.0
```

## 3. Criar a Release no GitHub

1. Acesse https://github.com/peder1981/GoXTree/releases/new
2. Tag: v1.0.0
3. Título: GoXTree v1.0.0 - Lançamento Oficial
4. Descrição: Cole o conteúdo do arquivo GITHUB_RELEASE.md
5. Anexe os binários compilados:
   - bin/goxTree-linux-amd64
   - bin/goxTree-linux-arm64
   - bin/goxTree-windows-amd64.exe
   - bin/goxTree-darwin-amd64
   - bin/goxTree-darwin-arm64
6. Marque como "Latest release"
7. Clique em "Publish release"

## 4. Configurar o Repositório

1. Acesse as configurações do repositório
2. Habilite Issues e Pull Requests
3. Configure as GitHub Pages para exibir a documentação
4. Adicione tópicos relevantes: go, file-manager, terminal, tui, xtree-gold

## 5. Promover o Projeto

1. Compartilhe nas redes sociais
2. Poste em fóruns como Reddit r/golang, r/commandline
3. Adicione ao awesome-go (https://github.com/avelino/awesome-go)
4. Considere escrever um artigo sobre o projeto

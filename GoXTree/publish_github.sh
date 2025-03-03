#!/bin/bash

# Script para publicar o GoXTree no GitHub
# Autor: Peder
# Data: 2025-03-03
# Versão: 1.2.0

echo "Publicando GoXTree v1.2.0 no GitHub..."

# Verificar se o repositório já existe
REPO_EXISTS=$(curl -s -o /dev/null -w "%{http_code}" https://github.com/peder1981/GoXTree)

if [ "$REPO_EXISTS" == "404" ]; then
    echo "O repositório não existe no GitHub. Por favor, crie-o primeiro em https://github.com/new"
    echo "Nome: GoXTree"
    echo "Descrição: Gerenciador de arquivos moderno inspirado no XTree Gold, implementado em Go"
    echo "Visibilidade: Pública"
    echo "Pressione Enter quando o repositório estiver criado..."
    read
fi

# Configurar o remote (caso ainda não esteja configurado)
if ! git remote | grep -q "^origin$"; then
    echo "Configurando remote 'origin'..."
    git remote add origin https://github.com/peder1981/GoXTree.git
fi

# Verificar a branch atual
CURRENT_BRANCH=$(git rev-parse --abbrev-ref HEAD)
echo "Branch atual: $CURRENT_BRANCH"

# Criar tag para a nova versão
echo "Criando tag v1.2.0..."
git tag -a v1.2.0 -m "Versão 1.2.0 - Correções e melhorias de compatibilidade"

# Enviar o código
echo "Enviando código para o GitHub..."
git push -u origin $CURRENT_BRANCH

# Enviar a tag
echo "Enviando tag v1.2.0 para o GitHub..."
git push origin v1.2.0

echo ""
echo "Código enviado com sucesso!"
echo ""
echo "Para criar a release no GitHub:"
echo "1. Acesse https://github.com/peder1981/GoXTree/releases/new"
echo "2. Tag: v1.2.0"
echo "3. Título: GoXTree v1.2.0 - Correções e Melhorias"
echo "4. Descrição: Cole o conteúdo do arquivo CHANGELOG.md para a versão 1.2.0"
echo "5. Anexe os binários compilados do diretório bin/"
echo "6. Marque como 'Latest release'"
echo "7. Clique em 'Publish release'"
echo ""
echo "Obrigado por usar o GoXTree!"

#!/bin/bash

# Script para preparar o repositório GoXTree para o GitHub
# Autor: Peder
# Data: 2025-03-02

echo "Preparando o repositório GoXTree para o GitHub..."

# Verificar se estamos no diretório correto
if [ ! -f "go.mod" ]; then
    echo "Erro: Este script deve ser executado no diretório raiz do GoXTree."
    exit 1
fi

# Adicionar todos os arquivos ao git
echo "Adicionando arquivos ao git..."
git add -A

# Verificar se há arquivos para commit
if git diff-index --quiet HEAD --; then
    echo "Nenhuma alteração para commit."
else
    # Fazer commit das alterações
    echo "Fazendo commit das alterações..."
    git commit -m "Preparação para lançamento no GitHub"
fi

# Criar tag v1.0.0 se não existir
if ! git tag | grep -q "^v1.0.0$"; then
    echo "Criando tag v1.0.0..."
    git tag -a v1.0.0 -m "Versão 1.0.0 - Lançamento inicial"
fi

# Verificar se o remote origin existe
if ! git remote | grep -q "^origin$"; then
    echo "O remote 'origin' não está configurado."
    echo "Por favor, crie o repositório no GitHub e execute:"
    echo "git remote add origin https://github.com/peder1981/GoXTree.git"
else
    echo "Remote 'origin' já configurado."
fi

# Compilar para todas as plataformas
echo "Compilando para todas as plataformas..."
bash build_all.sh

echo ""
echo "Repositório preparado com sucesso!"
echo ""
echo "Para enviar para o GitHub, execute:"
echo "git push -u origin main"
echo "git push origin v1.0.0"
echo ""
echo "Para criar um pacote completo, execute:"
echo "./package.sh"

#!/bin/bash

# Script para criar um pacote completo do GoXTree para o GitHub
# Autor: Peder
# Data: 2025-03-02

echo "Criando pacote do GoXTree para o GitHub..."

# Verificar se o diretório bin existe
if [ ! -d "bin" ]; then
    mkdir -p bin
fi

# Compilar para todas as plataformas
echo "Compilando para todas as plataformas..."
bash build_all.sh

# Criar diretório para o pacote
PACKAGE_DIR="goxTree-package"
rm -rf $PACKAGE_DIR
mkdir -p $PACKAGE_DIR

# Copiar arquivos essenciais
echo "Copiando arquivos para o pacote..."
cp -r README.md README_pt-BR.md LICENSE CHANGELOG.md DEVELOPERS.md RELEASE_NOTES.md GITHUB_RELEASE.md $PACKAGE_DIR/
cp -r assets $PACKAGE_DIR/
cp -r bin $PACKAGE_DIR/

# Copiar arquivos do GitHub
echo "Copiando arquivos do GitHub..."
mkdir -p $PACKAGE_DIR/.github
cp -r .github/* $PACKAGE_DIR/.github/

# Copiar código-fonte
echo "Copiando código-fonte..."
mkdir -p $PACKAGE_DIR/cmd $PACKAGE_DIR/pkg
cp -r cmd/* $PACKAGE_DIR/cmd/
cp -r pkg/* $PACKAGE_DIR/pkg/

# Copiar arquivos de configuração
echo "Copiando arquivos de configuração..."
cp go.mod $PACKAGE_DIR/
cp .gitignore $PACKAGE_DIR/

# Criar arquivo ZIP
echo "Criando arquivo ZIP..."
zip -r goxTree-package.zip $PACKAGE_DIR

# Limpar
echo "Limpando arquivos temporários..."
rm -rf $PACKAGE_DIR

echo "Pacote criado com sucesso: goxTree-package.zip"
echo "Este pacote está pronto para ser enviado para o GitHub."

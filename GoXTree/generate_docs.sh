#!/bin/bash

# Script para gerar documentação do GoXTree

# Verificar se o godoc está instalado
if ! command -v godoc &> /dev/null; then
    echo "godoc não está instalado. Instalando..."
    go install golang.org/x/tools/cmd/godoc@latest
fi

# Criar diretório para documentação gerada
mkdir -p docs/godoc

# Gerar documentação HTML
echo "Gerando documentação HTML..."
godoc -html -url=/pkg/github.com/peder1981/GoXTree/ > docs/godoc/index.html
godoc -html -url=/pkg/github.com/peder1981/GoXTree/pkg/ui/ > docs/godoc/ui.html
godoc -html -url=/pkg/github.com/peder1981/GoXTree/pkg/editor/ > docs/godoc/editor.html
godoc -html -url=/pkg/github.com/peder1981/GoXTree/pkg/viewer/ > docs/godoc/viewer.html
godoc -html -url=/pkg/github.com/peder1981/GoXTree/pkg/utils/ > docs/godoc/utils.html

echo "Documentação gerada com sucesso em docs/godoc/"
echo "Para visualizar a documentação, abra o arquivo docs/godoc/index.html em um navegador"

#!/bin/bash

# Script para abrir a página de criação de release no GitHub
# Autor: Peder
# Data: 2025-03-03
# Versão: 1.2.0

echo "Abrindo página de criação de release no GitHub..."

# Verificar se o xdg-open está disponível
if command -v xdg-open &> /dev/null; then
    xdg-open "https://github.com/peder1981/GoXTree/releases/new?tag=v1.2.0" &
elif command -v open &> /dev/null; then
    open "https://github.com/peder1981/GoXTree/releases/new?tag=v1.2.0"
else
    echo "Não foi possível abrir o navegador automaticamente."
    echo "Por favor, acesse manualmente: https://github.com/peder1981/GoXTree/releases/new?tag=v1.2.0"
fi

echo ""
echo "Instruções para criar a release:"
echo "1. Título: GoXTree v1.2.0 - Correções e Melhorias"
echo "2. Descrição: Cole o conteúdo do arquivo RELEASE_NOTES_1.2.0.md"
echo "3. Anexe os binários compilados do diretório bin/:"
echo "   - bin/goxTree-linux-amd64"
echo "   - bin/goxTree-linux-arm64"
echo "   - bin/goxTree-windows-amd64.exe"
echo "   - bin/goxTree-darwin-amd64"
echo "   - bin/goxTree-darwin-arm64"
echo "4. Marque como 'Latest release'"
echo "5. Clique em 'Publish release'"
echo ""
echo "Obrigado por usar o GoXTree!"

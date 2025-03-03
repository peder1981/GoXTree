#!/bin/bash

# Script para compilar o GoXTree para múltiplas plataformas
# Autor: Peder
# Data: 2025-03-03
# Versão: 1.2.0

# Criar diretório de saída se não existir
mkdir -p bin

echo "Compilando GoXTree v1.2.0 para múltiplas plataformas..."

# Linux AMD64
echo "Compilando para Linux AMD64..."
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/goxTree-linux-amd64 ./cmd/gxtree/main.go
chmod +x bin/goxTree-linux-amd64

# Linux ARM64
echo "Compilando para Linux ARM64..."
CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o bin/goxTree-linux-arm64 ./cmd/gxtree/main.go
chmod +x bin/goxTree-linux-arm64

# Windows AMD64
echo "Compilando para Windows AMD64..."
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o bin/goxTree-windows-amd64.exe ./cmd/gxtree/main.go

# macOS AMD64
echo "Compilando para macOS AMD64..."
CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o bin/goxTree-darwin-amd64 ./cmd/gxtree/main.go
chmod +x bin/goxTree-darwin-amd64

# macOS ARM64 (Apple Silicon)
echo "Compilando para macOS ARM64 (Apple Silicon)..."
CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -o bin/goxTree-darwin-arm64 ./cmd/gxtree/main.go
chmod +x bin/goxTree-darwin-arm64

# Criar arquivo de versão
echo "1.2.0" > bin/VERSION

echo "Compilação concluída! Binários disponíveis no diretório bin/"
echo ""
echo "Binários gerados:"
ls -lh bin/

# Verificar tamanho total dos binários
echo ""
echo "Tamanho total dos binários:"
du -sh bin/

#!/bin/bash

# Script para facilitar o uso do GoXTreeTester

# Verificar se o GoXTree existe
if [ ! -d "/home/peder/my-go-projects/GoXTree" ]; then
  echo "Erro: O diretório do GoXTree não foi encontrado em /home/peder/my-go-projects/GoXTree"
  exit 1
fi

# Compilar o GoXTreeTester
echo "Compilando o GoXTreeTester..."
go build -o gxtester ./cmd/gxtester
if [ $? -ne 0 ]; then
  echo "Erro: Falha ao compilar o GoXTreeTester"
  exit 1
fi

# Definir opções padrão
PROJECT_PATH="/home/peder/my-go-projects/GoXTree"
REPORT_PATH="./report.html"
VERBOSE=false
AUTOFIX=false
CHECK_STYLE=false
FIX_IMPORTS=false
ANALYZE_ONLY=false
TEST_ONLY=false

# Processar argumentos
while [[ $# -gt 0 ]]; do
  case $1 in
    --report=*)
      REPORT_PATH="${1#*=}"
      shift
      ;;
    --verbose)
      VERBOSE=true
      shift
      ;;
    --autofix)
      AUTOFIX=true
      shift
      ;;
    --check-style)
      CHECK_STYLE=true
      shift
      ;;
    --fix-imports)
      FIX_IMPORTS=true
      shift
      ;;
    --analyze-only)
      ANALYZE_ONLY=true
      shift
      ;;
    --test-only)
      TEST_ONLY=true
      shift
      ;;
    *)
      echo "Opção desconhecida: $1"
      echo "Uso: $0 [--report=CAMINHO] [--verbose] [--autofix] [--check-style] [--fix-imports] [--analyze-only] [--test-only]"
      exit 1
      ;;
  esac
done

# Construir comando
CMD="./gxtester --project=$PROJECT_PATH --report=$REPORT_PATH"

if [ "$VERBOSE" = true ]; then
  CMD="$CMD --verbose"
fi

if [ "$AUTOFIX" = true ]; then
  CMD="$CMD --autofix"
fi

if [ "$CHECK_STYLE" = true ]; then
  CMD="$CMD --check-style"
fi

if [ "$FIX_IMPORTS" = true ]; then
  CMD="$CMD --fix-imports"
fi

if [ "$ANALYZE_ONLY" = true ]; then
  CMD="$CMD --analyze-only"
fi

if [ "$TEST_ONLY" = true ]; then
  CMD="$CMD --test-only"
fi

# Executar o GoXTreeTester
echo "Executando o GoXTreeTester..."
echo "Comando: $CMD"
eval $CMD

# Verificar resultado
if [ $? -eq 0 ]; then
  echo "GoXTreeTester executado com sucesso!"
  echo "Relatório gerado em: $REPORT_PATH"
else
  echo "Erro: Falha ao executar o GoXTreeTester"
  exit 1
fi

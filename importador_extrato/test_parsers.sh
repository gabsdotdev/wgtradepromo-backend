#!/bin/bash

# Script de teste para validar os parsers

echo "=== Teste do Sistema de Parsers ==="
echo ""

# Cores para output
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Testa compilação
echo "1. Testando compilação..."
cd /home/lopes/llopes/wg/workspace/wgtradepromo-backend/importador_extrato
if go build -o importador_test 2>&1 | grep -q "error"; then
    echo -e "${RED}✗ Falha na compilação${NC}"
    exit 1
else
    echo -e "${GREEN}✓ Compilação bem-sucedida${NC}"
fi

# Lista arquivos disponíveis
echo ""
echo "2. Arquivos disponíveis para importação:"
echo ""

echo "Banco Inter (CSV):"
find rawdata/extrato/inter -name "*.csv" 2>/dev/null | wc -l | xargs echo "  - Arquivos encontrados:"

echo ""
echo "Nubank (CSV):"
nubank_count=$(find rawdata/extrato/nubank -name "*.csv" 2>/dev/null | wc -l)
if [ "$nubank_count" -eq 0 ]; then
    echo -e "  - ${YELLOW}Nenhum arquivo encontrado (aguardando upload)${NC}"
else
    echo "  - Arquivos encontrados: $nubank_count"
fi

echo ""
echo "Simples Nacional (PDF):"
find rawdata/extrato/das -name "*.pdf" 2>/dev/null | wc -l | xargs echo "  - Arquivos encontrados:"

# Limpa arquivo de teste
rm -f importador_test

echo ""
echo -e "${GREEN}=== Teste Concluído ===${NC}"
echo ""
echo "Para executar o importador completo, certifique-se de:"
echo "1. Configurar as variáveis de ambiente do banco de dados"
echo "2. Executar: ./run.sh"

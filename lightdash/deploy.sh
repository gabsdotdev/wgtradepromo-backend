#!/bin/bash
# sudo apt install yamllint && chmod +x lightdash/deploy.sh
# Cores para output
GREEN='\033[0;32m'
RED='\033[0;31m'
NC='\033[0m'

echo "🔍 Validando arquivos YAML..."

# Validar sintaxe YAML
yamllint dash_financeiro_geral.yml
if [ $? -ne 0 ]; then
    echo -e "${RED}❌ Erro de sintaxe no dash_financeiro_geral.yml${NC}"
    exit 1
fi

yamllint charts_financeiro.yml
if [ $? -ne 0 ]; then
    echo -e "${RED}❌ Erro de sintaxe no charts_financeiro.yml${NC}"
    exit 1
fi

echo -e "${GREEN}✅ Sintaxe YAML validada com sucesso${NC}"

# Validar com Lightdash CLI
echo "🚀 Iniciando deploy no Lightdash..."

# Atualizar charts
echo "📊 Atualizando charts..."
lightdash deploy charts_financeiro.yml
if [ $? -ne 0 ]; then
    echo -e "${RED}❌ Erro ao fazer deploy dos charts${NC}"
    exit 1
fi

# Atualizar dashboard
echo "📈 Atualizando dashboard..."
lightdash deploy dash_financeiro_geral.yml
if [ $? -ne 0 ]; then
    echo -e "${RED}❌ Erro ao fazer deploy do dashboard${NC}"
    exit 1
fi

echo -e "${GREEN}✅ Deploy concluído com sucesso!${NC}"

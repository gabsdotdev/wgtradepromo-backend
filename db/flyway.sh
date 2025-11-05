#!/usr/bin/env bash
set -e

# ==============================
# 🐳 Script para rodar Flyway via Docker
# Uso:
#   ./flyway.sh [comando]
# Exemplo:
#   ./flyway.sh migrate
#   ./flyway.sh info
#   ./flyway.sh validate
# ==============================

# --- 1. Lê o comando (padrão: migrate)
FLYWAY_CMD=${1:-migrate}

# --- 2. Resolve IP do host Windows
HOST_IP=$(getent hosts windows-host | awk '{print $1}' || true)

# Se não conseguir resolver, tenta pegar IP do nameserver (WSL2)
if [ -z "$HOST_IP" ]; then
  HOST_IP=$(grep nameserver /etc/resolv.conf | awk '{print $2}')
fi

if [ -z "$HOST_IP" ]; then
  echo "❌ Não foi possível determinar o IP do host Windows."
  exit 1
fi

echo "➡️  Resolvendo 'windows-host' como $HOST_IP"
echo "➡️  Executando comando Flyway: $FLYWAY_CMD"
echo

# --- 3. Configurações do banco (ajuste conforme necessário)
DB_NAME="wgtrade"
DB_USER="postgres"
DB_PASS="postgres"
DB_PORT="5432"

# --- 4. Executa o Flyway no container
docker run --rm \
  -v "$(pwd)/migration:/flyway/sql" \
  -e FLYWAY_LOCATIONS=filesystem:/flyway/sql \
  -e FLYWAY_URL="jdbc:postgresql://$HOST_IP:$DB_PORT/$DB_NAME" \
  -e FLYWAY_USER="$DB_USER" \
  -e FLYWAY_PASSWORD="$DB_PASS" \
  flyway/flyway:11.15.0 $FLYWAY_CMD

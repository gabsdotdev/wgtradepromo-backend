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

# --- 2. Define DB_HOST apenas se ainda não estiver definido no ambiente
if [[ -z "${DB_HOST:-}" ]]; then
  # Tenta resolver IP do host Windows
  HOST_IP=$(getent hosts windows-host | awk '{print $1}' || true)

  # Se não conseguir resolver, tenta pegar IP do nameserver (WSL2)
  if [[ -z "$HOST_IP" ]]; then
    HOST_IP="$(grep -m1 nameserver /etc/resolv.conf | awk '{print $2}' || true)"
  fi

  if [[ -z "$HOST_IP" ]]; then
    echo "❌ Não foi possível determinar o IP do host Windows e DB_HOST não está definido."
    echo "   Defina DB_HOST ou garanta que 'windows-host' resolva corretamente."
    exit 1
  fi

  DB_HOST="$HOST_IP"
  echo "➡️  DB_HOST não definido; resolvido 'windows-host' como: $DB_HOST"
else
  echo "➡️  DB_HOST já definido no ambiente: $DB_HOST"
fi

echo "➡️  Executando comando Flyway: $FLYWAY_CMD"
echo

# --- 3. Configurações do banco (ajuste conforme necessário)
DB_NAME="${DB_NAME:-wgtrade}"
DB_USER="${DB_USER:-postgres}"
DB_PASS="${DB_PASS:-postgres}"
DB_PORT="${DB_PORT:-5432}"

# --- 5. Executa o Flyway no container
echo "Executando Flyway:"
echo "  Host: $DB_HOST"
echo "  Porta: $DB_PORT"
echo "  Banco: $DB_NAME"
echo "  Usuário: $DB_USER"
echo "  Comando: $FLYWAY_CMD"
echo

# --- 5. Executa o Flyway no container
docker run --rm \
  -v "$(pwd)/migration:/flyway/sql" \
  -e FLYWAY_LOCATIONS=filesystem:/flyway/sql \
  -e FLYWAY_URL="jdbc:postgresql://$DB_HOST:$DB_PORT/$DB_NAME" \
  -e FLYWAY_USER="$DB_USER" \
  -e FLYWAY_PASSWORD="$DB_PASS" \
  flyway/flyway:11.15.0 $FLYWAY_CMD

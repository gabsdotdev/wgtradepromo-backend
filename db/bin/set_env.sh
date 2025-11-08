#!/bin/bash
# Script para definir variáveis de ambiente do banco de dados PostgreSQL

if [[ "${BASH_SOURCE[0]}" == "${0}" ]]; then
  echo "⚠️  Execute este script com 'source set_env.sh' para manter as variáveis no shell atual."
  exit 1
fi

PROFILE="${1:-local}"
ENV_FILE=".env.${PROFILE}.db"
EXAMPLE_FILE=".env.example.db"
GITIGNORE_FILE=".gitignore"

if [[ ! -f "$ENV_FILE" ]]; then
  if [[ -f "$EXAMPLE_FILE" ]]; then
    cp -f "$EXAMPLE_FILE" "$ENV_FILE"
    echo "⚠️  Arquivo '$ENV_FILE' não existia e foi criado a partir de '$EXAMPLE_FILE'."
    echo "📝  Edite '$ENV_FILE' e preencha as variáveis antes de rodar novamente:"
    echo "    source set_env.sh ${PROFILE}"
    
    if [[ -f "$GITIGNORE_FILE" ]]; then
      if ! grep -Fxq "$ENV_FILE" "$GITIGNORE_FILE"; then
        echo "" >> "$GITIGNORE_FILE"
        echo "$ENV_FILE" >> "$GITIGNORE_FILE"
        echo "📁  Linha adicionada ao .gitignore: $ENV_FILE"
      fi
    else
      echo "$ENV_FILE" > "$GITIGNORE_FILE"
      echo "📁  Criado novo .gitignore com entrada: $ENV_FILE"
    fi

    return 1
  else
    echo "❌ Arquivo de exemplo '$EXAMPLE_FILE' não encontrado."
    echo "   Crie manualmente '$ENV_FILE' com este modelo:"
    cat <<'EOF'
DB_NAME="__DB_NAME__"
DB_USER="__DB_USER__"
DB_PASS="__DB_PASS__"
DB_PORT="5432"
DB_HOST="__DB_HOST__"
EOF
    return 1
  fi
fi

set -a
source "$ENV_FILE"
set +a

# Variáveis obrigatórias
req=(DB_NAME DB_USER DB_PASS DB_PORT)

# Validação compacta
fail=0
for v in "${req[@]}"; do
  [ -n "${!v:-}" ] || { echo "❌ Falta definir: $v"; fail=1; }
done

if [ "$fail" -ne 0 ]; then
  echo
  echo "💡 Corrija as variáveis acima e rode novamente."
  # se script foi chamado com 'source', usa return; senão, usa exit
  (return 0 2>/dev/null) && return 1 || exit 1
fi

# --- Máscara (mostra só 2 primeiros e 2 últimos caracteres)
if [ -n "$DB_PASS" ]; then
  pass_len=${#DB_PASS}
  if [ "$pass_len" -le 4 ]; then
    masked_pass="$DB_PASS"
  else
    start=${DB_PASS:0:2}
    end=${DB_PASS: -2}
    middle_len=$((pass_len - 4))
    masked_pass="${start}$(printf '%*s' "$middle_len" '' | tr ' ' '*')${end}"
  fi
else
  masked_pass="(vazio)"
fi

echo "Variáveis de ambiente configuradas:"
echo "DB_NAME=$DB_NAME"
echo "DB_USER=$DB_USER"
echo "DB_PASS=$masked_pass"
echo "DB_PORT=$DB_PORT"
echo "DB_HOST=$DB_HOST"

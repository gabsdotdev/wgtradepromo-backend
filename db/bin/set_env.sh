#!/bin/bash
# Script para definir variáveis de ambiente do banco de dados PostgreSQL

if [[ "${BASH_SOURCE[0]}" == "${0}" ]]; then
  echo "⚠️  Execute este script com 'source set_env.sh' para manter as variáveis no shell atual."
  exit 1
fi

export DB_NAME=""
export DB_USER=""
export DB_PASS=""
export DB_PORT=""
export DB_HOST=""

# Variáveis obrigatórias
req=(DB_NAME DB_USER DB_PASS DB_PORT DB_HOST)

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

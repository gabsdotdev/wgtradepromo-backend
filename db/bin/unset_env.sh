#!/bin/bash
# Script para remover variáveis de ambiente do banco de dados PostgreSQL

if [[ "${BASH_SOURCE[0]}" == "${0}" ]]; then
  echo "⚠️  Execute este script com 'source set_env.sh' para manter as variáveis no shell atual."
  exit 1
fi

unset DB_NAME
unset DB_USER
unset DB_PASS
unset DB_PORT
unset DB_HOST

echo "Variáveis de ambiente removidas:"
echo "DB_NAME, DB_USER, DB_PASS, DB_PORT e DB_HOST foram desfeitas."

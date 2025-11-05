-- =========================================================
-- UNDO: Remover tabelas e índices criados anteriormente
-- =========================================================

-- 1. Remover índices de transacoes
DROP INDEX IF EXISTS ix_transacoes_tipo_transacao;
DROP INDEX IF EXISTS ix_transacoes_tipo_operacao;
DROP INDEX IF EXISTS ix_transacoes_data;
DROP INDEX IF EXISTS ix_transacoes_conta;

-- 2. Remover tabela transacoes
DROP TABLE IF EXISTS transacoes CASCADE;

-- 3. Remover índices de contas
DROP INDEX IF EXISTS ix_contas_ativo;
DROP INDEX IF EXISTS ix_contas_empresa;

-- 4. Remover tabela contas
DROP TABLE IF EXISTS contas CASCADE;

-- 5. Remover índices de empresas
DROP INDEX IF EXISTS ix_empresas_ativa;

-- 6. Remover tabela empresas
DROP TABLE IF EXISTS empresas CASCADE;

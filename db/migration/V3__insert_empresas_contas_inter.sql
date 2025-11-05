-- =========================================================
-- Inserir CONTAS no BANCO INTER
-- =========================================================

-- ⚠️ Para associar corretamente as contas às empresas,
-- precisamos recuperar os UUIDs das empresas recém-criadas.

WITH empresa_wg_trade AS (
    SELECT id FROM empresas WHERE cnpj = '33722929000183'
),
empresa_wg_promo AS (
    SELECT id FROM empresas WHERE cnpj = '23756078000136'
)
INSERT INTO contas (
    id, empresa_id, banco, agencia, numero, nome,
    saldo_inicial, ativo, criado_em, atualizado_em
)
VALUES
    (
        '019a5259-946f-7ce5-86c3-877cae7f0d88',
        (SELECT id FROM empresa_wg_trade),
        'Inter',
        '0001',
        '264583213',
        'Conta Inter PJ — Matriz',
        0,
        TRUE,
        NOW(),
        NOW()
    ),
    (
        '019a5259-c655-7532-bb1b-d2635394e0e4',
        (SELECT id FROM empresa_wg_promo),
        'Inter',
        '0001',
        '344126161',
        'Conta Inter PJ — Filial',
        0,
        TRUE,
        NOW(),
        NOW()
    );

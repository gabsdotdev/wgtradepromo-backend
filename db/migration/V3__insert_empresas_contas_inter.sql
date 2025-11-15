-- =========================================================
-- Inserir CONTAS no BANCO INTER
-- =========================================================

-- ⚠️ Para associar corretamente as contas às empresas,
-- precisamos recuperar os UUIDs das empresas recém-criadas.

WITH empresa_wg_trade AS (
    SELECT id FROM cadastros.empresas WHERE cnpj = '33722929000183'
),
empresa_wg_promo AS (
    SELECT id FROM cadastros.empresas WHERE cnpj = '23756078000136'
)
INSERT INTO financeiro.contas (
    id, empresa_id, banco, agencia, numero, nome,
    saldo_inicial, ativo, criado_em, atualizado_em
)
VALUES
    (
        '019a5259-946f-7ce5-86c3-877cae7f0d88',
        (SELECT id FROM empresa_wg_trade),
        '077 - Inter',
        '0001 ',
        '264583213',
        'Conta Inter PJ — WG Trade',
        0,
        TRUE,
        NOW(),
        NOW()
    ),
    (
        '019a5259-c655-7532-bb1b-d2635394e0e4',
        (SELECT id FROM empresa_wg_promo),
        '077 - Inter',
        '0001',
        '344126161',
        'Conta Inter PJ — WG Promo',
        0,
        TRUE,
        NOW(),
        NOW()
    ),
    (
        '019a7d96-59d9-7a47-9e05-6e3c0920c08a',
        (SELECT id FROM empresa_wg_promo),
        '260 - Nubank',
        '0001',
        '6115932263',
        'Conta Nubank PJ — WG Promo',
        0,
        TRUE,
        NOW(),
        NOW()
    );

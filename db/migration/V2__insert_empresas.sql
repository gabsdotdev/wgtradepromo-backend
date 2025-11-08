-- =========================================================
-- Inserção de cadastros.empresas
-- =========================================================

INSERT INTO cadastros.empresas (
    id,
    nome,
    cnpj,
    ativa,
    criado_em,
    atualizado_em
) VALUES 
    (
        '019a5253-a0ee-73d0-9e68-3fa3676e048a',
        'WG TRADE PROMOCOES E EVENTOS',
        '33722929000183', -- CNPJ sem máscara
        TRUE,
        NOW(),
        NOW()
    ),
    (
        '019a5253-cc5d-7039-997c-d1d5c5831c23',
        'W & G PROMOCOES E EVENTOS',
        '23756078000136',
        TRUE,
        NOW(),
        NOW()
    );
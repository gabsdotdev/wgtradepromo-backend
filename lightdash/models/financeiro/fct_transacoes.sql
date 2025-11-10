{{ config(materialized='view') }}

WITH t AS (
  SELECT *
  FROM {{ source('financeiro', 'transacoes') }}
),
c AS (
  SELECT *
  FROM {{ source('financeiro', 'contas') }}
),
e AS (
  SELECT *
  FROM {{ source('cadastros', 'empresas') }}
)
SELECT
  t.id,
  t.conta_id,
  c.empresa_id,
  e.nome               AS empresa,
  c.nome               AS conta,
  t.data,
  t.titulo,
  t.descricao,
  t.tipo_operacao,           -- 'credito' | 'debito'
  t.tipo_transacao,
  t.valor,                   -- sempre positivo (constraint)
  CASE
    WHEN t.tipo_operacao = 'credito' THEN t.valor
    WHEN t.tipo_operacao = 'debito'  THEN -t.valor
    ELSE 0
  END AS valor_signed
FROM t
JOIN c ON c.id = t.conta_id
JOIN e ON e.id = c.empresa_id
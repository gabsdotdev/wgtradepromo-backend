-- =========================================================
-- TABELA: financeiro.das_documentos
-- =========================================================
CREATE TABLE financeiro.das_documentos (
  id                  UUID PRIMARY KEY,
  empresa_id          UUID NOT NULL REFERENCES cadastros.empresas(id) ON DELETE RESTRICT,

  -- Período de apuração: sempre o primeiro dia do mês
  periodo_apuracao    DATE NOT NULL,
  CONSTRAINT ck_das_periodo_dia1 CHECK (DATE_PART('day', periodo_apuracao) = 1),

  data_vencimento     DATE NOT NULL,

  numero_documento    VARCHAR(40) NOT NULL,
  valor_total         NUMERIC(14,2) NOT NULL,

  -- Status textual controlado
  status              VARCHAR(20) NOT NULL DEFAULT 'EMITIDO',
  CONSTRAINT ck_das_status CHECK (status IN ('EMITIDO', 'PAGO', 'VENCIDO', 'CANCELADO')),

  -- Caminho ou URL do arquivo PDF (opcional)
  arquivo_path        TEXT,

  criado_em           TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  atualizado_em       TIMESTAMPTZ NOT NULL DEFAULT NOW(),

  -- Regras de unicidade
  CONSTRAINT uq_das_empresa_periodo UNIQUE (empresa_id, periodo_apuracao),
  CONSTRAINT uq_das_numero_documento UNIQUE (numero_documento)
);

-- =========================================================
-- ÍNDICES recomendados
-- =========================================================
CREATE INDEX ix_das_periodo ON financeiro.das_documentos (periodo_apuracao);
CREATE INDEX ix_das_status ON financeiro.das_documentos (status);

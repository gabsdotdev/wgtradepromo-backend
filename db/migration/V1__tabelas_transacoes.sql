-- =========================================================
-- TABELA: empresas
-- =========================================================
CREATE TABLE empresas (
  id              UUID PRIMARY KEY,                               
  nome            VARCHAR(120) NOT NULL,
  cnpj            CHAR(14) NOT NULL,                              
  ativa           BOOLEAN NOT NULL DEFAULT TRUE,
  criado_em       TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  atualizado_em   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  CONSTRAINT uq_empresas_cnpj UNIQUE (cnpj)
);

CREATE INDEX ix_empresas_ativa ON empresas (ativa);

-- =========================================================
-- TABELA: contas
-- =========================================================
CREATE TABLE contas (
  id              UUID PRIMARY KEY,                               
  empresa_id      UUID NOT NULL REFERENCES empresas(id) ON DELETE RESTRICT,
  banco           VARCHAR(60) NOT NULL,                           
  agencia         VARCHAR(20),
  numero          VARCHAR(40),
  nome            VARCHAR(100) NOT NULL,                      
  saldo_inicial   NUMERIC(14,2) NOT NULL DEFAULT 0,
  ativo           BOOLEAN NOT NULL DEFAULT TRUE,
  criado_em       TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  atualizado_em   TIMESTAMPTZ NOT NULL DEFAULT NOW(),

  CONSTRAINT uq_conta_empresa_banco_agencia_numero
    UNIQUE (empresa_id, banco, agencia, numero)
);

CREATE INDEX ix_contas_empresa      ON contas (empresa_id);
CREATE INDEX ix_contas_ativo        ON contas (ativo);

-- =========================================================
-- TABELA: transacoes
-- =========================================================
CREATE TABLE transacoes (
  id                      UUID PRIMARY KEY,                    
  conta_id                UUID NOT NULL REFERENCES contas(id) ON DELETE RESTRICT,

  data                    DATE NOT NULL,
  titulo                  VARCHAR(150) NOT NULL,
  descricao               TEXT,

  tipo_operacao           VARCHAR(10) NOT NULL CHECK (tipo_operacao IN ('credito','debito')),
  tipo_transacao          VARCHAR(30)  NOT NULL,                 

  valor                   NUMERIC(14,2) NOT NULL CHECK (valor > 0),

  criado_em               TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  atualizado_em           TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX ix_transacoes_conta          ON transacoes (conta_id);
CREATE INDEX ix_transacoes_data           ON transacoes (data);
CREATE INDEX ix_transacoes_tipo_operacao  ON transacoes (tipo_operacao);
CREATE INDEX ix_transacoes_tipo_transacao ON transacoes (tipo_transacao);
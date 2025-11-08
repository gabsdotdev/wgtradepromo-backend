package models

import (
	"time"
)

// Transaction represents a financial transaction in the database
type Transaction struct {
	ID            string    `db:"id"`
	ContaID       string    `db:"conta_id"`
	Data          time.Time `db:"data"`
	Titulo        string    `db:"titulo"`
	Descricao     string    `db:"descricao"`
	TipoOperacao  string    `db:"tipo_operacao"`
	TipoTransacao string    `db:"tipo_transacao"`
	Valor         float64   `db:"valor"`
	CriadoEm      time.Time `db:"criado_em"`
	AtualizadoEm  time.Time `db:"atualizado_em"`
	Fingerprint   string    `db:"fingerprint"`
}

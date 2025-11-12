package models

import "time"

type DasDocumento struct {
	ID              string    `db:"id"`
	EmpresaID       string    `db:"empresa_id"`
	PeriodoApuracao time.Time `db:"periodo_apuracao"`
	DataVencimento  time.Time `db:"data_vencimento"`
	NumeroDocumento string    `db:"numero_documento"`
	ValorTotal      float64   `db:"valor_total"`
	Status          string    `db:"status"` // EMITIDO, PAGO, VENCIDO, CANCELADO
	ArquivoPath     *string   `db:"arquivo_path"`
	CriadoEm        time.Time `db:"criado_em"`
	AtualizadoEm    time.Time `db:"atualizado_em"`
}

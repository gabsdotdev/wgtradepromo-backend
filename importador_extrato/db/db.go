package db

import (
	"crypto/sha256"
	"fmt"
	"time"

	"github.com/gabsdotdev/wgtradepromo-backend/importador_extrato/models"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type DB struct {
	*sqlx.DB
}

func NewConnection(connectionString string) (*DB, error) {
	db, err := sqlx.Connect("postgres", connectionString)
	if err != nil {
		return nil, fmt.Errorf("error connecting to database: %v", err)
	}

	return &DB{db}, nil
}

func (db *DB) InsertTransaction(contaID string, date time.Time, description, details string, amount float64) error {
	// Generate UUID v7 for the transaction
	id := uuid.Must(uuid.NewV7())

	// Generate fingerprint to avoid duplicates
	fingerprint := generateFingerprint(contaID, date, description, details, amount)

	// Determine operation type based on amount
	tipoOperacao := "debito"
	if amount > 0 {
		tipoOperacao = "credito"
	}

	// Use absolute value for the amount as per database schema
	if amount < 0 {
		amount = -amount
	}

	now := time.Now()

	// Create transaction record
	tx := &models.Transaction{
		ID:            id.String(),
		ContaID:       contaID,
		Data:          date,
		Titulo:        description,
		Descricao:     details,
		TipoOperacao:  tipoOperacao,
		TipoTransacao: getTipoTransacao(description),
		Valor:         amount,
		CriadoEm:      now,
		AtualizadoEm:  now,
		Fingerprint:   fingerprint,
	}

	// Insert into database
	query := `
		INSERT INTO financeiro.transacoes (
			id, conta_id, data, titulo, descricao,
			tipo_operacao, tipo_transacao, valor,
			criado_em, atualizado_em, fingerprint
		) VALUES (
			:id, :conta_id, :data, :titulo, :descricao,
			:tipo_operacao, :tipo_transacao, :valor,
			:criado_em, :atualizado_em, :fingerprint
		)
	`

	_, err := db.NamedExec(query, tx)
	if err != nil {
		return fmt.Errorf("error inserting transaction: %v", err)
	}

	return nil
}

func generateFingerprint(contaID string, date time.Time, description, details string, amount float64) string {
	data := fmt.Sprintf("%s|%s|%s|%s|%.2f",
		contaID,
		date.Format("2006-01-02"),
		description,
		details,
		amount,
	)
	hash := sha256.Sum256([]byte(data))
	return fmt.Sprintf("%x", hash)
}

func getTipoTransacao(description string) string {
	// Map common transaction descriptions to types
	switch {
	case contains(description, "Pix"):
		return "pix"
	case contains(description, "Transferência"):
		return "transferencia"
	case contains(description, "Pagamento"):
		return "pagamento"
	default:
		return "outros"
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && s[:len(substr)] == substr
}

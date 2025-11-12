package parser

import "time"

// Statement representa um extrato bancário ou documento fiscal
type Statement struct {
	AccountNumber string
	Period        string
	Balance       float64
	Transactions  []Transaction
	DasDocumento  *DasDocumento
}

// Transaction representa uma transação financeira
type Transaction struct {
	Date        time.Time
	Description string
	Details     string
	Amount      float64
	Balance     float64
}

// DasDocumento representa um documento DAS do Simples Nacional
type DasDocumento struct {
	CNPJ            string
	PeriodoApuracao time.Time
	DataVencimento  time.Time
	NumeroDocumento string
	ValorTotal      float64
}

// Parser é a interface que todos os parsers devem implementar
type Parser interface {
	// Parse processa o arquivo e retorna um Statement
	Parse(filename string) (*Statement, error)

	// CanParse verifica se o parser pode processar o arquivo
	CanParse(filename string) bool

	// GetName retorna o nome do parser (ex: "Inter", "Nubank", "Simples Nacional")
	GetName() string
}

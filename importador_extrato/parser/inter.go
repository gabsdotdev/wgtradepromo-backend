package parser

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// InterParser é o parser para extratos do Banco Inter em formato CSV
type InterParser struct{}

// NewInterParser cria uma nova instância do parser do Inter
func NewInterParser() *InterParser {
	return &InterParser{}
}

// GetName retorna o nome do parser
func (p *InterParser) GetName() string {
	return "Banco Inter"
}

// CanParse verifica se o arquivo pode ser processado por este parser
func (p *InterParser) CanParse(filename string) bool {
	// Verifica se o arquivo está na pasta inter e é um CSV
	ext := strings.ToLower(filepath.Ext(filename))
	if ext != ".csv" {
		return false
	}

	// Verifica se está na pasta inter
	return strings.Contains(filepath.ToSlash(filename), "/extrato/inter/")
}

// Parse processa um arquivo CSV do Banco Inter
func (p *InterParser) Parse(filename string) (*Statement, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("error opening file: %v", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.Comma = ';'
	reader.FieldsPerRecord = -1
	reader.TrimLeadingSpace = true
	reader.LazyQuotes = true

	// Skip header line
	_, err = reader.Read()
	if err != nil {
		return nil, fmt.Errorf("error reading header: %v", err)
	}

	// Read account info
	stmt := &Statement{}
	accountLine, err := reader.Read()
	if err != nil {
		return nil, fmt.Errorf("error reading account info: %v", err)
	}
	stmt.AccountNumber = strings.TrimSpace(accountLine[1])

	// Read period
	periodLine, err := reader.Read()
	if err != nil {
		return nil, fmt.Errorf("error reading period: %v", err)
	}
	stmt.Period = strings.TrimSpace(periodLine[1])

	// Read balance
	balanceLine, err := reader.Read()
	if err != nil {
		return nil, fmt.Errorf("error reading balance: %v", err)
	}
	stmt.Balance = parseAmount(balanceLine[1])

	// Skip column headers
	_, err = reader.Read()
	if err != nil {
		return nil, fmt.Errorf("error reading column headers: %v", err)
	}

	// Read transactions
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("error reading transaction: %v", err)
		}

		date, err := parseDate(record[0])
		if err != nil {
			return nil, fmt.Errorf("error parsing date: %v", err)
		}

		transaction := Transaction{
			Date:        date,
			Description: strings.TrimSpace(record[1]),
			Details:     strings.TrimSpace(record[2]),
			Amount:      parseAmount(record[3]),
			Balance:     parseAmount(record[4]),
		}
		stmt.Transactions = append(stmt.Transactions, transaction)
	}

	return stmt, nil
}

func parseDate(dateStr string) (time.Time, error) {
	return time.Parse("02/01/2006", strings.TrimSpace(dateStr))
}

func parseAmount(amountStr string) float64 {
	// Remove R$ and spaces
	clean := strings.TrimSpace(strings.ReplaceAll(amountStr, "R$", ""))

	// Replace comma with dot for decimal point
	clean = strings.ReplaceAll(clean, ".", "")
	clean = strings.ReplaceAll(clean, ",", ".")

	// Convert to float
	amount := 0.0
	fmt.Sscanf(clean, "%f", &amount)
	return amount
}

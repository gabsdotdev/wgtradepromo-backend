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

// NubankParser é o parser para extratos do Nubank em formato CSV
type NubankParser struct{}

// NewNubankParser cria uma nova instância do parser do Nubank
func NewNubankParser() *NubankParser {
	return &NubankParser{}
}

// GetName retorna o nome do parser
func (p *NubankParser) GetName() string {
	return "Nubank"
}

// CanParse verifica se o arquivo pode ser processado por este parser
func (p *NubankParser) CanParse(filename string) bool {
	// Verifica se o arquivo está na pasta nubank e é um CSV
	ext := strings.ToLower(filepath.Ext(filename))
	if ext != ".csv" {
		return false
	}

	// Verifica se está na pasta nubank
	return strings.Contains(filepath.ToSlash(filename), "/extrato/nubank/")
}

// Parse processa um arquivo CSV do Nubank
func (p *NubankParser) Parse(filename string) (*Statement, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("error opening file: %v", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.Comma = ','
	reader.FieldsPerRecord = -1
	reader.TrimLeadingSpace = true

	stmt := &Statement{
		AccountNumber: extractAccountFromFilename(filename),
		Transactions:  []Transaction{},
	}

	// Read and skip header line
	header, err := reader.Read()
	if err != nil {
		return nil, fmt.Errorf("error reading header: %v", err)
	}

	// O formato do Nubank pode variar, esta é uma implementação base
	// Esperamos colunas como: Data, Descrição, Valor, etc.
	// Ajuste conforme o formato real do arquivo
	_ = header // Para quando tivermos o arquivo real, podemos validar as colunas

	// Read transactions
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("error reading transaction: %v", err)
		}

		// Este é um formato exemplo - ajustar quando tiver arquivo real
		// Formato esperado: Data, Descrição, Valor
		if len(record) < 3 {
			continue
		}

		date, err := parseNubankDate(record[0])
		if err != nil {
			fmt.Printf("Warning: error parsing date '%s': %v - skipping transaction\n", record[0], err)
			continue
		}

		amount := parseNubankAmount(record[2])

		transaction := Transaction{
			Date:        date,
			Description: strings.TrimSpace(record[1]),
			Details:     "", // Ajustar se houver campo de detalhes
			Amount:      amount,
		}
		stmt.Transactions = append(stmt.Transactions, transaction)
	}

	return stmt, nil
}

// parseNubankDate tenta diferentes formatos de data usados pelo Nubank
func parseNubankDate(dateStr string) (time.Time, error) {
	dateStr = strings.TrimSpace(dateStr)

	// Formatos comuns do Nubank
	formats := []string{
		"2006-01-02", // ISO format: 2024-01-15
		"02/01/2006", // Brazilian format: 15/01/2024
		"01/02/2006", // US format: 01/15/2024
	}

	for _, format := range formats {
		if date, err := time.Parse(format, dateStr); err == nil {
			return date, nil
		}
	}

	return time.Time{}, fmt.Errorf("unable to parse date: %s", dateStr)
}

// parseNubankAmount converte string de valor para float64
func parseNubankAmount(amountStr string) float64 {
	// Remove espaços
	clean := strings.TrimSpace(amountStr)

	// Remove R$ se presente
	clean = strings.ReplaceAll(clean, "R$", "")
	clean = strings.TrimSpace(clean)

	// Replace comma with dot for decimal point
	clean = strings.ReplaceAll(clean, ",", ".")

	// Convert to float
	amount := 0.0
	fmt.Sscanf(clean, "%f", &amount)
	return amount
}

// extractAccountFromFilename extrai o número da conta do nome do arquivo
func extractAccountFromFilename(filename string) string {
	// Tenta extrair um padrão de conta do nome do arquivo
	// Ex: nubank-123456-2024.csv -> 123456
	base := filepath.Base(filename)
	parts := strings.Split(base, "-")
	if len(parts) > 1 {
		return parts[1]
	}
	return "nubank-default"
}

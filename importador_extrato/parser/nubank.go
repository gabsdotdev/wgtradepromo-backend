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

		// Formato do CSV do Nubank: Data, Valor, Identificador, Descrição
		if len(record) < 4 {
			continue
		}

		date, err := parseNubankDate(record[0])
		if err != nil {
			fmt.Printf("Warning: error parsing date '%s': %v - skipping transaction\n", record[0], err)
			continue
		}

		amount := parseNubankAmount(record[1])
		identifier := strings.TrimSpace(record[2])
		fullDescription := strings.TrimSpace(record[3])

		// Extrai título e detalhes da descrição
		description, details := extractNubankDescriptionAndDetails(fullDescription)

		transaction := Transaction{
			Date:        date,
			Description: description,
			Details:     fmt.Sprintf("ID: %s | %s", identifier, details),
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
	// Ex: NU_6115932263_13JAN2025_12NOV2025.csv -> 6115932263
	base := filepath.Base(filename)
	parts := strings.Split(base, "_")
	if len(parts) >= 2 {
		return parts[1]
	}
	return "nubank-default"
}

// extractNubankDescriptionAndDetails extrai título e detalhes da descrição do Nubank
func extractNubankDescriptionAndDetails(fullDescription string) (string, string) {
	// Padrões comuns de descrição do Nubank:
	// "Transferência recebida pelo Pix - NOME - CPF/CNPJ - BANCO (CODIGO) Agência: X Conta: Y"
	// "Transferência enviada pelo Pix - NOME - CPF/CNPJ - BANCO (CODIGO) Agência: X Conta: Y"
	// "Pagamento - ESTABELECIMENTO"
	// etc.

	// Separa pelo primeiro hífen para extrair o tipo de transação
	parts := strings.SplitN(fullDescription, " - ", 2)

	if len(parts) < 2 {
		// Se não houver hífen, usa a descrição completa como título
		return fullDescription, ""
	}

	transactionType := strings.TrimSpace(parts[0])
	remainingInfo := strings.TrimSpace(parts[1])

	return transactionType, remainingInfo
}

package parser

import (
	"errors"
	"fmt"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/ledongthuc/pdf"
)

// DasSimplesNacionalParser é o parser para documentos DAS do Simples Nacional em formato PDF
type DasSimplesNacionalParser struct{}

// NewDasSimplesNacionalParser cria uma nova instância do parser do Simples Nacional
func NewDasSimplesNacionalParser() *DasSimplesNacionalParser {
	return &DasSimplesNacionalParser{}
}

// GetName retorna o nome do parser
func (p *DasSimplesNacionalParser) GetName() string {
	return "Simples Nacional"
}

// CanParse verifica se o arquivo pode ser processado por este parser
func (p *DasSimplesNacionalParser) CanParse(filename string) bool {
	// Verifica se o arquivo está na pasta simples_nacional e é um PDF
	ext := strings.ToLower(filepath.Ext(filename))
	if ext != ".pdf" {
		return false
	}

	// Verifica se está na pasta simples_nacional
	return strings.Contains(filepath.ToSlash(filename), "/extrato/das/")
}

// Parse processa um arquivo PDF do Simples Nacional (DAS)
func (p *DasSimplesNacionalParser) Parse(filename string) (*Statement, error) {
	// Abre o PDF
	f, r, err := pdf.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("error opening PDF: %v", err)
	}
	defer f.Close()

	stmt := &Statement{
		AccountNumber: "simples-nacional", // Identificador especial para Simples Nacional
		Transactions:  []Transaction{},
	}

	var fullText strings.Builder

	// Extrai texto de todas as páginas
	totalPages := r.NumPage()
	for pageNum := 1; pageNum <= totalPages; pageNum++ {
		page := r.Page(pageNum)
		if page.V.IsNull() {
			continue
		}

		text, err := page.GetPlainText(nil)
		if err != nil {
			fmt.Printf("Warning: error reading page %d: %v\n", pageNum, err)
			continue
		}

		fullText.WriteString(text)
		fullText.WriteString("\n")
	}

	text := fullText.String()

	// Parse das informações do DAS
	dasDocumento, err := p.extractDASTransactions(text)
	if err != nil {
		return nil, fmt.Errorf("error extracting DAS transactions: %v", err)
	}

	stmt.DasDocumento = &dasDocumento

	return stmt, nil
}

// extractDASTransactions extrai transações do texto do DAS
func (p *DasSimplesNacionalParser) extractDASTransactions(texto string) (DasDocumento, error) {
	var out DasDocumento

	// Padrões regex para extrair informações do DAS
	// Expressões regulares específicas para cada informação
	cnpjRegex := regexp.MustCompile(`\d{2}\.\d{3}\.\d{3}/\d{4}-\d{2}`)
	numeroDocRegex := regexp.MustCompile(`\d{2}\.\d{2}\.\d{5}\.\d{7}-\d{1}`)
	periodoRegex := regexp.MustCompile(`([A-ZÇ][a-zç]+/[0-9]{4})`) // ex: Outubro/2025
	vencimentoRegex := regexp.MustCompile(`(\d{2}/\d{2}/\d{4})`)
	valorRegex := regexp.MustCompile(`(\d{1,3}(?:\.\d{3})*,\d{2})`)

	// Extrair os valores
	cnpj := cnpjRegex.FindString(texto)
	numeroDocumento := numeroDocRegex.FindString(texto)
	periodoStr := periodoRegex.FindString(texto)
	vencimentoStr := vencimentoRegex.FindString(texto)
	valorStr := valorRegex.FindString(texto)

	periodo, err := parsePeriodoPT(periodoStr) // 1º dia do mês
	if err != nil {
		return out, fmt.Errorf("periodo_apuracao inválido (%q): %w", periodoStr, err)
	}

	vencimento, err := parseDateBR(vencimentoStr)
	if err != nil {
		return out, fmt.Errorf("data_vencimento inválida (%q): %w", vencimentoStr, err)
	}

	valor, err := parseValorBR(valorStr)
	if err != nil {
		return out, fmt.Errorf("valor_total inválido (%q): %w", valorStr, err)
	}

	out.CNPJ = regexp.MustCompile(`\D`).ReplaceAllString(cnpj, "")
	out.NumeroDocumento = regexp.MustCompile(`\D`).ReplaceAllString(numeroDocumento, "")
	out.PeriodoApuracao = periodo
	out.DataVencimento = vencimento
	out.ValorTotal = valor

	return out, nil
}

// parseBRLAmount converte string em formato brasileiro (1.234,56) para float64
func parseBRLAmount(amountStr string) float64 {
	// Remove espaços
	clean := strings.TrimSpace(amountStr)

	// Remove pontos (separador de milhar)
	clean = strings.ReplaceAll(clean, ".", "")

	// Substitui vírgula por ponto (separador decimal)
	clean = strings.ReplaceAll(clean, ",", ".")

	// Convert to float
	amount := 0.0
	fmt.Sscanf(clean, "%f", &amount)
	return amount
}

// parseDateBR converte "dd/mm/aaaa" para time.Time (UTC, sem hora de preocupação)
func parseDateBR(s string) (time.Time, error) {
	return time.Parse("02/01/2006", s)
}

// parsePeriodoPT converte "Outubro/2025" -> time.Date(2025, 10, 1, 0, 0, 0, 0, time.UTC)
func parsePeriodoPT(s string) (time.Time, error) {
	partes := strings.Split(s, "/")
	if len(partes) != 2 {
		return time.Time{}, errors.New("formato de período inválido")
	}
	mesNome := strings.ToLower(strings.TrimSpace(partes[0]))
	anoStr := strings.TrimSpace(partes[1])

	mesMap := map[string]time.Month{
		"janeiro":   time.January,
		"fevereiro": time.February,
		"março":     time.March, "marco": time.March,
		"abril":    time.April,
		"maio":     time.May,
		"junho":    time.June,
		"julho":    time.July,
		"agosto":   time.August,
		"setembro": time.September,
		"outubro":  time.October,
		"novembro": time.November,
		"dezembro": time.December,
	}
	mes, ok := mesMap[mesNome]
	if !ok {
		return time.Time{}, fmt.Errorf("mês inválido: %q", mesNome)
	}

	ano, err := strconv.Atoi(anoStr)
	if err != nil || ano < 1900 || ano > 3000 {
		return time.Time{}, fmt.Errorf("ano inválido: %q", anoStr)
	}

	return time.Date(ano, mes, 1, 0, 0, 0, 0, time.UTC), nil
}

// parseValorBR converte "4.803,14" -> 4803.14 (float64)
// Se preferir evitar float, troque o retorno para int64 (centavos).
func parseValorBR(s string) (float64, error) {
	s = strings.TrimSpace(s)
	s = strings.ReplaceAll(s, ".", "")
	s = strings.ReplaceAll(s, ",", ".")
	return strconv.ParseFloat(s, 64)
}

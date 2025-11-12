package parser

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/ledongthuc/pdf"
)

// ExtratoSimplesNacionalParser é o parser para documentos DAS do Simples Nacional em formato PDF
type ExtratoSimplesNacionalParser struct {
	parser *FiscalParser
}

// NewExtratoSimplesNacionalParser cria uma nova instância do parser do Simples Nacional
func NewExtratoSimplesNacionalParser() *ExtratoSimplesNacionalParser {
	return &ExtratoSimplesNacionalParser{parser: NewFiscalParser()}
}

// GetName retorna o nome do parser
func (p *ExtratoSimplesNacionalParser) GetName() string {
	return "Extrato Simples Nacional"
}

// CanParse verifica se o arquivo pode ser processado por este parser
func (p *ExtratoSimplesNacionalParser) CanParse(filename string) bool {
	// Verifica se o arquivo está na pasta simples_nacional e é um PDF
	ext := strings.ToLower(filepath.Ext(filename))
	if ext != ".pdf" {
		return false
	}

	// Verifica se está na pasta simples_nacional
	return strings.Contains(filepath.ToSlash(filename), "/extrato/extrato_simples_nacional/")
}

// Parse processa um arquivo PDF do Simples Nacional (DAS)
func (p *ExtratoSimplesNacionalParser) Parse(filename string) (*Statement, error) {
	// Abre o PDF
	f, r, err := pdf.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("error opening PDF: %v", err)
	}
	defer f.Close()

	stmt := &Statement{
		AccountNumber: "extrato-simples-nacional", // Identificador especial para Extrato Simples Nacional
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
func (p *ExtratoSimplesNacionalParser) extractDASTransactions(texto string) (DasDocumento, error) {
	var out DasDocumento

	os.WriteFile("conteudo.txt", []byte(texto), 0644)

	conteudo := cutAfterCaseInsensitive(texto, "principal")
	if conteudo == "" {
		conteudo = texto
	}

	// Padrões regex para extrair informações do DAS
	// Expressões regulares específicas para cada informação
	cnpjRegex := regexp.MustCompile(`(\d{2}\.\d{3}\.\d{3}/\d{4}-\d{2})`)
	numeroDocRegex := regexp.MustCompile(`(?i)Número\s*[:\s]*(\d{17})`)
	periodoRegex := regexp.MustCompile(`(?i)Período\s*de\s*Apuração\s*\(PA\)\s*[:\s]*([0-9]{2}/[0-9]{4})`)
	vencimentoRegex := regexp.MustCompile(`(?i)Data\s*de\s*Vencimento\s*[:\s]*([0-9]{2}/[0-9]{2}/[0-9]{4})`)
	valorRegex := regexp.MustCompile(`(?i)Total\s*[:\s]*([\d.,]+)`)

	// Extrair os valores
	cnpj := match1(cnpjRegex, texto)
	numeroDocumento := match1(numeroDocRegex, texto)
	periodoStr := match1(periodoRegex, texto)
	vencimentoStr := match1(vencimentoRegex, texto)
	valorStr := match1(valorRegex, conteudo)

	periodo, err := p.parser.ParsePeriodoNum(strings.TrimSpace(periodoStr)) // 1º dia do mês
	if err != nil {
		return out, fmt.Errorf("periodo_apuracao inválido (%q): %w", periodoStr, err)
	}

	vencimento, err := p.parser.ParseDateBR(vencimentoStr)
	if err != nil {
		return out, fmt.Errorf("data_vencimento inválida (%q): %w", vencimentoStr, err)
	}

	valor, err := p.parser.ParseValorBR(valorStr)
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

func cutAfterCaseInsensitive(haystack, needle string) string {
	hl := strings.ToLower(haystack)
	nl := strings.ToLower(needle)
	idx := strings.Index(hl, nl)
	if idx == -1 {
		return ""
	}
	// avança o comprimento do 'needle' na string original (mesmo offset)
	return haystack[idx+len(needle):]
}

func match1(re *regexp.Regexp, s string) string {
	m := re.FindStringSubmatch(s)
	if len(m) > 1 {
		return strings.TrimSpace(m[1])
	}
	return ""
}

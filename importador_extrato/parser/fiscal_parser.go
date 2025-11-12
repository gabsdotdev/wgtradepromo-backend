package parser

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// FiscalParser encapsula utilitários para parse de campos financeiros/fiscais brasileiros.
type FiscalParser struct {
	Location *time.Location // permite customizar timezone (ex.: time.Local, time.UTC)
}

// NewFiscalParser retorna uma instância padrão (UTC).
func NewFiscalParser() *FiscalParser {
	return &FiscalParser{Location: time.UTC}
}

// ParseDateBR converte "dd/mm/aaaa" para time.Time
func (p *FiscalParser) ParseDateBR(s string) (time.Time, error) {
	t, err := time.ParseInLocation("02/01/2006", strings.TrimSpace(s), p.Location)
	if err != nil {
		return time.Time{}, fmt.Errorf("data inválida (%q): %w", s, err)
	}
	return t, nil
}

// ParsePeriodoPT converte "Outubro/2025" → 2025-10-01
func (p *FiscalParser) ParsePeriodoPT(s string) (time.Time, error) {
	partes := strings.Split(s, "/")
	if len(partes) != 2 {
		return time.Time{}, errors.New("formato de período inválido (esperado: Mês/Ano)")
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

	return time.Date(ano, mes, 1, 0, 0, 0, 0, p.Location), nil
}

// ParsePeriodoNum converte "11/2025" → 2025-11-01
func (p *FiscalParser) ParsePeriodoNum(s string) (time.Time, error) {
	partes := strings.Split(s, "/")
	if len(partes) != 2 {
		return time.Time{}, errors.New("formato de período inválido (esperado: MM/AAAA)")
	}

	mesStr := strings.TrimSpace(partes[0])
	anoStr := strings.TrimSpace(partes[1])

	mes, err := strconv.Atoi(mesStr)
	if err != nil || mes < 1 || mes > 12 {
		return time.Time{}, fmt.Errorf("mês inválido: %q", mesStr)
	}

	ano, err := strconv.Atoi(anoStr)
	if err != nil || ano < 1900 || ano > 3000 {
		return time.Time{}, fmt.Errorf("ano inválido: %q", anoStr)
	}

	return time.Date(ano, time.Month(mes), 1, 0, 0, 0, 0, p.Location), nil
}

// ParseValorBR converte "4.803,14" → 4803.14
// Se quiser evitar float, retorne int64 (centavos).
func (p *FiscalParser) ParseValorBR(s string) (float64, error) {
	s = strings.TrimSpace(s)
	s = strings.ReplaceAll(s, ".", "")
	s = strings.ReplaceAll(s, ",", ".")
	v, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0, fmt.Errorf("valor inválido (%q): %w", s, err)
	}
	return v, nil
}

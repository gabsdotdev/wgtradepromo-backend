package parser

import (
	"fmt"
	"path/filepath"
)

// ParserFactory cria e retorna o parser apropriado para um arquivo
type ParserFactory struct {
	parsers []Parser
}

// NewParserFactory cria uma nova instância da factory com todos os parsers disponíveis
func NewParserFactory() *ParserFactory {
	return &ParserFactory{
		parsers: []Parser{
			NewInterParser(),
			NewNubankParser(),
			NewDasSimplesNacionalParser(),
			NewExtratoSimplesNacionalParser(),
		},
	}
}

// GetParser retorna o parser apropriado para o arquivo fornecido
func (f *ParserFactory) GetParser(filename string) (Parser, error) {
	for _, parser := range f.parsers {
		if parser.CanParse(filename) {
			return parser, nil
		}
	}

	return nil, fmt.Errorf("no parser found for file: %s", filepath.Base(filename))
}

// GetAllParsers retorna todos os parsers disponíveis
func (f *ParserFactory) GetAllParsers() []Parser {
	return f.parsers
}

// ListSupportedParsers retorna os nomes de todos os parsers suportados
func (f *ParserFactory) ListSupportedParsers() []string {
	names := make([]string, len(f.parsers))
	for i, parser := range f.parsers {
		names[i] = parser.GetName()
	}
	return names
}

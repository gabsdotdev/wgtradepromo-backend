package parser

// ParseCSV é mantido para compatibilidade com código existente
// Deprecated: Use InterParser.Parse() instead
func ParseCSV(filename string) (*Statement, error) {
	parser := NewInterParser()
	return parser.Parse(filename)
}

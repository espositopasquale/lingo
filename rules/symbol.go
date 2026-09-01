package rules

type SymbolRule struct{}

func (SymbolRule) Name() string {
	return "symbol"
}

func (SymbolRule) Priority() int {
	return 10
}

func (SymbolRule) Match(
	text string,
	start int,
) (Match, bool) {

	c := NewCursor(text, start)

	if !c.IsSymbol() {
		return Match{}, false
	}

	c = c.Advance()

	return Match{
		Start: start,
		End:   c.Pos(),
		Type:  Symbol,
	}, true
}

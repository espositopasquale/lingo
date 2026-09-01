package rules

type PunctuationRule struct{}

func (PunctuationRule) Name() string {
	return "punctuation"
}

func (PunctuationRule) Priority() int {
	return 50
}

func (PunctuationRule) Match(
	text string,
	start int,
) (Match, bool) {

	c := NewCursor(text, start)

	if !c.IsPunct() {
		return Match{}, false
	}

	r, _ := c.Rune()
	c = c.Advance()

	for !c.End() {
		next, _ := c.Rune()

		if next != r {
			break
		}

		c = c.Advance()
	}

	return Match{
		Start: start,
		End:   c.Pos(),
		Type:  Punctuation,
	}, true
}

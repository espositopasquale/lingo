package rules

type DateRule struct{}

func (DateRule) Name() string {
	return "date"
}

func (DateRule) Priority() int {
	return 500
}

func (DateRule) Match(text string, start int) (Match, bool) {
	c := NewCursor(text, start)

	var ok bool

	c, ok = readDigits(c, 1, 2)
	if !ok {
		return Match{}, false
	}

	separator, _ := c.Rune()

	switch separator {
	case '/', '-', '.':
	default:
		return Match{}, false
	}

	c = c.Advance()

	c, ok = readDigits(c, 1, 2)
	if !ok {
		return Match{}, false
	}

	if !c.MatchRune(separator) {
		return Match{}, false
	}

	c = c.Advance()

	c, ok = readDigits(c, 4, 4)
	if !ok {
		return Match{}, false
	}

	return Match{
		Start: start,
		End:   c.Pos(),
		Type:  Date,
	}, true
}

func readDigits(c Cursor, min, max int) (Cursor, bool) {
	count := 0

	for !c.End() &&
		count < max &&
		c.IsDigit() {

		c = c.Advance()
		count++
	}

	return c, count >= min
}

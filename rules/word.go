package rules

type WordRule struct{}

func (WordRule) Name() string {
	return "word"
}

func (WordRule) Priority() int {
	return 100
}

func (WordRule) Match(text string, start int) (Match, bool) {
	c := NewCursor(text, start)

	if !c.IsLetter() {
		return Match{}, false
	}

	c = c.Advance()

	for !c.End() {
		if c.IsLetter() {
			c = c.Advance()
			continue
		}

		if isApostrophe(c) {
			next := c.Advance()

			if !next.End() && next.IsLetter() {
				c = next.Advance()
				continue
			}
		}

		break
	}

	return Match{
		Start: start,
		End:   c.Pos(),
		Type:  Word,
	}, true
}

func isApostrophe(c Cursor) bool {
	if c.End() {
		return false
	}

	r, _ := c.Rune()

	switch r {
	case '\'', '’', 'ʼ':
		return true
	default:
		return false
	}
}

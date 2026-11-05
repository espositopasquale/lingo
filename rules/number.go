package rules

type NumberRule struct{}

func (NumberRule) Name() string {
	return "number"
}

func (NumberRule) Priority() int {
	return 300
}

func (NumberRule) Match(text string, start int) (Match, bool) {
	c := NewCursor(text, start)

	if c.MatchRune('+') || c.MatchRune('-') {
		c = c.Advance()
	}

	if !c.IsDigit() {
		return Match{}, false
	}

	for c.IsDigit() {
		c = c.Advance()
	}

	for c.MatchRune('.') {
		next := c.Advance()
		digits := 0

		for next.IsDigit() && digits < 3 {
			next = next.Advance()
			digits++
		}

		if digits != 3 || next.IsDigit() {
			break
		}

		c = next
	}

	if c.MatchRune(',') {
		next := c.Advance()

		if next.IsDigit() {
			c = next

			for c.IsDigit() {
				c = c.Advance()
			}
		}
	}

	return Match{
		Start: start,
		End:   c.Pos(),
		Type:  Number,
	}, true
}

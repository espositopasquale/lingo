package rules

type TimeRule struct{}

func (TimeRule) Name() string {
	return "time"
}

func (TimeRule) Priority() int {
	return 450
}

func (TimeRule) Match(text string, start int) (Match, bool) {
	c := NewCursor(text, start)

	var ok bool

	c, ok = readDigits(c, 1, 2)
	if !ok {
		return Match{}, false
	}

	if !c.MatchRune(':') {
		return Match{}, false
	}

	c = c.Advance()

	c, ok = readDigits(c, 2, 2)
	if !ok {
		return Match{}, false
	}

	return Match{
		Start: start,
		End:   c.Pos(),
		Type:  Time,
	}, true
}

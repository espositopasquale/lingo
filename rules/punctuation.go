package rules

type PunctuationRule struct {
	// Single keeps adjacent punctuation marks separate, as the UD profile expects.
	// An ellipsis (...) stays together.
	Single bool
}

func (PunctuationRule) Name() string {
	return "punctuation"
}

func (PunctuationRule) Priority() int {
	return 50
}

func (rule PunctuationRule) Match(
	text string,
	start int,
) (Match, bool) {

	c := NewCursor(text, start)

	if !c.IsPunct() {
		return Match{}, false
	}

	r, _ := c.Rune()
	c = c.Advance()
	if rule.Single {
		if r == '.' && start+3 <= len(text) && text[start:start+3] == "..." {
			c = NewCursor(text, start+3)
		}
		return Match{Start: start, End: c.Pos(), Type: Punctuation}, true
	}

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

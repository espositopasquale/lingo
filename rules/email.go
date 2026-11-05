package rules

import (
	"strings"
	"unicode"
)

type EmailRule struct{}

func (EmailRule) Name() string {
	return "email"
}

func (EmailRule) Priority() int {
	return 700
}

func (EmailRule) Match(text string, start int) (Match, bool) {
	if start < 0 || start >= len(text) {
		return Match{}, false
	}

	c := NewCursor(text, start)

	if !isEmailLocalChar(c) {
		return Match{}, false
	}

	hasLocal := false

	for !c.End() {
		r, _ := c.Rune()

		if isEmailLocalRune(r) {
			hasLocal = true
			c = c.Advance()
			continue
		}

		break
	}

	if !hasLocal || !c.MatchRune('@') {
		return Match{}, false
	}

	c = c.Advance()

	labels := 0
	lastLabel := ""
	for {
		labelStart := c.Pos()
		for !c.End() {
			r, _ := c.Rune()
			if !unicode.IsLetter(r) && !unicode.IsDigit(r) && r != '-' {
				break
			}
			c = c.Advance()
		}
		label := text[labelStart:c.Pos()]
		if label == "" || strings.HasPrefix(label, "-") || strings.HasSuffix(label, "-") {
			return Match{}, false
		}
		labels++
		lastLabel = label
		if !c.MatchRune('.') || !c.Advance().IsLetter() && !c.Advance().IsDigit() {
			break
		}
		c = c.Advance()
	}
	if labels < 2 || strings.IndexFunc(lastLabel, func(r rune) bool { return !unicode.IsLetter(r) }) >= 0 {
		return Match{}, false
	}

	return Match{
		Start: start,
		End:   c.Pos(),
		Type:  Email,
	}, true
}

func isEmailLocalChar(c Cursor) bool {
	if c.End() {
		return false
	}

	r, _ := c.Rune()

	return isEmailLocalRune(r)
}

func isEmailLocalRune(r rune) bool {
	return unicode.IsLetter(r) ||
		unicode.IsDigit(r) ||
		strings.ContainsRune(
			".-_+",
			r,
		)
}

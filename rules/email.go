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

	domainStart := c.Pos()
	domainChars := 0

	for !c.End() {
		r, _ := c.Rune()

		if unicode.IsLetter(r) ||
			unicode.IsDigit(r) ||
			r == '-' {

			domainChars++
			c = c.Advance()
			continue
		}

		break
	}

	if domainChars == 0 {
		return Match{}, false
	}

	if c.Pos() == domainStart {
		return Match{}, false
	}

	domain := text[domainStart:c.Pos()]

	if strings.HasPrefix(domain, "-") ||
		strings.HasSuffix(domain, "-") {
		return Match{}, false
	}

	if !c.MatchRune('.') {
		return Match{}, false
	}

	c = c.Advance()

	tldStart := c.Pos()

	for !c.End() {
		r, _ := c.Rune()

		if unicode.IsLetter(r) {
			c = c.Advance()
			continue
		}

		break
	}

	if c.Pos() == tldStart {
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

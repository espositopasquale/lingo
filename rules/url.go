package rules

import (
	"strings"
	"unicode"
	"unicode/utf8"
)

type URLRule struct{}

func (URLRule) Name() string {
	return "url"
}

func (URLRule) Priority() int {
	return 800
}

func (URLRule) Match(text string, start int) (Match, bool) {
	if start < 0 || start >= len(text) {
		return Match{}, false
	}

	remaining := text[start:]

	prefixes := []string{
		"https://",
		"http://",
		"www.",
	}

	var prefix string

	for _, candidate := range prefixes {
		if strings.HasPrefix(
			strings.ToLower(remaining),
			candidate,
		) {
			prefix = candidate
			break
		}
	}

	if prefix == "" {
		return Match{}, false
	}

	c := NewCursor(text, start+len(prefix))

	hostStart := c.Pos()

	for !c.End() {
		r, _ := c.Rune()

		if unicode.IsLetter(r) ||
			unicode.IsDigit(r) ||
			strings.ContainsRune(
				".-_",
				r,
			) {

			c = c.Advance()
			continue
		}

		break
	}

	if c.Pos() == hostStart {
		return Match{}, false
	}

	host := text[hostStart:c.Pos()]

	if strings.HasPrefix(host, ".") ||
		strings.HasSuffix(host, ".") {

		return Match{}, false
	}

	for !c.End() {
		r, _ := c.Rune()

		if unicode.IsSpace(r) ||
			unicode.IsPunct(r) &&
				!strings.ContainsRune(
					"/?&=#%.-_:",
					r,
				) {

			break
		}

		c = c.Advance()
	}

	end := trimURLPunctuation(text, start, c.Pos())

	if end <= start+len(prefix) {
		return Match{}, false
	}

	return Match{
		Start: start,
		End:   end,
		Type:  URL,
	}, true
}

func trimURLPunctuation(
	text string,
	start int,
	end int,
) int {
	for end > start {
		r, size := lastRune(text[start:end])

		switch r {
		case '.', ',', ';', ':', '!', '?', ')', ']', '}':
			end -= size
		default:
			return end
		}
	}

	return end
}

func lastRune(text string) (rune, int) {
	for len(text) > 0 {
		r, size := utf8.DecodeLastRuneInString(text)

		return r, size
	}

	return 0, 0
}

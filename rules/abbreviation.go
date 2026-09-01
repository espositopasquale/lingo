package rules

import "strings"

type AbbreviationRule struct {
	words map[string]struct{}
}

func NewAbbreviationRule(words []string) *AbbreviationRule {
	m := make(map[string]struct{}, len(words))

	for _, word := range words {
		word = strings.TrimSpace(word)

		if word == "" {
			continue
		}

		m[strings.ToLower(word)] = struct{}{}
	}

	return &AbbreviationRule{
		words: m,
	}
}

func (r *AbbreviationRule) Name() string {
	return "abbreviation"
}

func (r *AbbreviationRule) Priority() int {
	return 400
}

func (r *AbbreviationRule) Match(
	text string,
	start int,
) (Match, bool) {
	if start < 0 || start >= len(text) {
		return Match{}, false
	}

	remaining := text[start:]

	bestEnd := -1

	for abbreviation := range r.words {
		if len(abbreviation) > len(remaining) {
			continue
		}

		candidate := remaining[:len(abbreviation)]

		if !strings.EqualFold(candidate, abbreviation) {
			continue
		}

		if !validAbbreviationBoundary(
			text,
			start,
			start+len(abbreviation),
		) {
			continue
		}

		if start+len(abbreviation) > bestEnd {
			bestEnd = start + len(abbreviation)
		}
	}

	if bestEnd == -1 {
		return Match{}, false
	}

	return Match{
		Start: start,
		End:   bestEnd,
		Type:  Abbreviation,
	}, true
}

func validAbbreviationBoundary(
	text string,
	start int,
	end int,
) bool {
	if end >= len(text) {
		return true
	}

	c := NewCursor(text, end)

	if c.IsLetter() || c.IsDigit() {
		return false
	}

	_ = start

	return true
}

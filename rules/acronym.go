package rules

import "strings"

type AcronymRule struct {
	words map[string]struct{}
}

func NewAcronymRule(words []string) *AcronymRule {
	m := make(map[string]struct{}, len(words))

	for _, word := range words {
		word = strings.TrimSpace(word)

		if word == "" {
			continue
		}

		m[strings.ToUpper(word)] = struct{}{}
	}

	return &AcronymRule{
		words: m,
	}
}

func (r *AcronymRule) Name() string {
	return "acronym"
}

func (r *AcronymRule) Priority() int {
	return 350
}

func (r *AcronymRule) Match(
	text string,
	start int,
) (Match, bool) {
	if start < 0 || start >= len(text) {
		return Match{}, false
	}

	c := NewCursor(text, start)

	if !c.IsLetter() {
		return Match{}, false
	}

	for c.IsLetter() {
		c = c.Advance()
	}

	end := c.Pos()
	value := text[start:end]

	if _, ok := r.words[strings.ToUpper(value)]; !ok {
		return Match{}, false
	}

	if end < len(text) {
		next := NewCursor(text, end)

		if next.IsLetter() || next.IsDigit() {
			return Match{}, false
		}
	}

	return Match{
		Start: start,
		End:   end,
		Type:  Acronym,
	}, true
}

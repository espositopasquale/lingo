package rules

import (
	"strings"
	"unicode"
	"unicode/utf8"
)

type AcronymRule struct {
	words lexiconNode
	mixed map[string]bool
}

func NewAcronymRule(words []string) *AcronymRule {
	r := &AcronymRule{mixed: make(map[string]bool)}
	for _, word := range words {
		word = strings.TrimSpace(word)
		if word == "" {
			continue
		}
		r.words.add(word, false)
		if word != strings.ToUpper(word) {
			r.mixed[word] = true
		}
	}
	return r
}

func (*AcronymRule) Name() string  { return "acronym" }
func (*AcronymRule) Priority() int { return 400 }

func uppercaseLetters(value string) bool {
	letters := 0
	for _, r := range value {
		if unicode.IsLetter(r) {
			if !unicode.IsUpper(r) {
				return false
			}
			letters++
		}
	}
	return letters > 0
}

func (r *AcronymRule) Match(text string, start int) (Match, bool) {
	if start < 0 || start >= len(text) {
		return Match{}, false
	}
	best := r.words.match(text, start, func(end int) bool {
		value := text[start:end]
		return uppercaseLetters(value) || r.mixed[value]
	})

	// Allow the last dot to be omitted in "A.A.A." sequences. Ordinary uppercase words still need a match.
	pos, count := start, 0
	for pos < len(text) {
		letter, size := utf8.DecodeRuneInString(text[pos:])
		if !unicode.IsUpper(letter) {
			break
		}
		pos += size
		count++
		if pos >= len(text) || text[pos] != '.' {
			break
		}
		pos++
	}
	if count >= 2 && !wordContinues(text, pos) && pos > best {
		best = pos
	}
	if best < 0 {
		return Match{}, false
	}
	return Match{Start: start, End: best, Type: Acronym}, true
}

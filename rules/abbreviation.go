package rules

import (
	"strings"
	"unicode"
	"unicode/utf8"
)

type lexiconNode struct {
	children    map[rune]*lexiconNode
	terminal    bool
	numericOnly bool
}

func (n *lexiconNode) add(word string, numericOnly bool) {
	for _, r := range strings.TrimSpace(word) {
		if n.children == nil {
			n.children = make(map[rune]*lexiconNode)
		}
		r = foldRune(r)
		if n.children[r] == nil {
			n.children[r] = &lexiconNode{}
		}
		n = n.children[r]
	}
	if !n.terminal {
		n.numericOnly = numericOnly
	} else {
		n.numericOnly = n.numericOnly && numericOnly
	}
	n.terminal = true
}

// Map case-equivalent runes to the same key, including σ/ς/Σ and k/K/K.
func foldRune(r rune) rune {
	min := r
	for next := unicode.SimpleFold(r); next != r; next = unicode.SimpleFold(next) {
		if next < min {
			min = next
		}
	}
	return min
}

func (n *lexiconNode) match(text string, start int, accept func(int) bool) int {
	best := -1
	for pos := start; pos < len(text); {
		r, size := utf8.DecodeRuneInString(text[pos:])
		n = n.children[foldRune(r)]
		if n == nil {
			break
		}
		pos += size
		if n.terminal && !wordContinues(text, pos) && (!n.numericOnly || followingNumber(text, pos)) && (accept == nil || accept(pos)) {
			best = pos
		}
	}
	return best
}

func followingNumber(text string, pos int) bool {
	for _, r := range text[pos:] {
		if unicode.IsSpace(r) {
			continue
		}
		return unicode.IsDigit(r)
	}
	return false
}

type AbbreviationRule struct{ words lexiconNode }

func NewAbbreviationRule(words []string) *AbbreviationRule {
	return NewAbbreviationRuleWithNumericOnly(words, nil)
}

func NewAbbreviationRuleWithNumericOnly(words, numericOnly []string) *AbbreviationRule {
	r := &AbbreviationRule{}
	for _, word := range words {
		if strings.TrimSpace(word) != "" {
			r.words.add(word, false)
		}
	}
	for _, word := range numericOnly {
		if strings.TrimSpace(word) != "" {
			r.words.add(word, true)
		}
	}
	return r
}

func (*AbbreviationRule) Name() string  { return "abbreviation" }
func (*AbbreviationRule) Priority() int { return 400 }

func (r *AbbreviationRule) Match(text string, start int) (Match, bool) {
	if start < 0 || start >= len(text) {
		return Match{}, false
	}
	end := r.words.match(text, start, nil)
	if end < 0 {
		return Match{}, false
	}
	return Match{Start: start, End: end, Type: Abbreviation}, true
}

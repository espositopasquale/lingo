package tokenizer

import (
	"strings"

	"github.com/espositopasquale/lin-go/rules"
)

type italianNumberRule struct{ ud bool }

func (italianNumberRule) Name() string  { return "italian-number" }
func (italianNumberRule) Priority() int { return 300 }

func (rule italianNumberRule) Match(text string, start int) (rules.Match, bool) {
	if rule.ud {
		c := rules.NewCursor(text, start)
		r, _ := c.Rune()
		if apostrophe(r) {
			next := c.Advance()
			count := 0
			for next.IsDigit() {
				count++
				next = next.Advance()
			}
			if count == 2 && !wordContinuation(text, next.Pos()) {
				return rules.Match{Start: start, End: next.Pos(), Type: rules.Number}, true
			}
		}
	}
	if !rule.ud {
		return (rules.NumberRule{}).Match(text, start)
	}
	match, ok := udNumberAtom(text, start)
	if !ok {
		return match, false
	}
	c := rules.NewCursor(text, match.End)
	if c.MatchRune('-') || c.MatchRune('–') || c.MatchRune('x') || c.MatchRune('×') {
		next := c.Advance()
		if next.IsDigit() {
			if second, found := udNumberAtom(text, next.Pos()); found {
				match.End = second.End
			}
		}
	}
	c = rules.NewCursor(text, match.End)
	if c.MatchRune('%') || c.MatchRune('‰') {
		match.End = c.Advance().Pos()
	} else if r, _ := c.Rune(); apostrophe(r) || r == '′' {
		if !wordContinuation(text, c.Advance().Pos()) {
			match.End = c.Advance().Pos()
		}
	}
	return match, true
}

func udNumberAtom(text string, start int) (rules.Match, bool) {
	match, ok := (rules.NumberRule{}).Match(text, start)
	if !ok {
		return match, false
	}
	c := rules.NewCursor(text, match.End)
	// Accept decimal points too, for values and version numbers such as 1.5 or 2.0.
	if c.MatchRune('.') && !strings.ContainsAny(text[start:match.End], ".,") && c.Advance().IsDigit() {
		c = c.Advance()
		for c.IsDigit() {
			c = c.Advance()
		}
		match.End = c.Pos()
	}
	return match, true
}

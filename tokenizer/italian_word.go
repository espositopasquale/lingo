package tokenizer

import (
	"strings"
	"unicode/utf8"

	"github.com/espositopasquale/lingo/rules"
)

type italianWordRule struct{ ud bool }

func (italianWordRule) Name() string  { return "italian-word" }
func (italianWordRule) Priority() int { return 100 }

// Split after these prefixes when an apostrophe joins them to the next word.
var italianElisions = map[string]bool{
	"l": true, "d": true, "all": true, "dall": true, "dell": true,
	"nell": true, "sull": true, "coll": true, "pell": true,
	"un": true, "quest": true, "quell": true, "sant": true,
	"bell": true, "buon": true, "grand": true, "tutt": true,
	"anch": true, "senz": true, "com": true, "c": true,
	"m": true, "t": true, "s": true, "v": true,
	"gl": true, "agl": true, "dagl": true, "degl": true,
	"negl": true, "sugl": true,
	"dov": true, "cos": true, "quant": true, "mezz": true,
	"cent": true, "vent": true, "trent": true, "quarant": true,
	"cinquant": true, "sessant": true, "settant": true,
	"ottant": true, "novant": true,
}

var italianTruncations = map[string]bool{
	"po": true, "be": true, "mo": true, "va": true,
	"fa": true, "sta": true, "da": true, "di": true, "to": true,
	"e": true, // Also accept E' or e' for è/È.
}

func apostrophe(r rune) bool { return r == '\'' || r == '’' || r == 'ʼ' }

func (rule italianWordRule) Match(text string, start int) (rules.Match, bool) {
	if rule.ud {
		c := rules.NewCursor(text, start)
		r, _ := c.Rune()
		if apostrophe(r) {
			next := c.Advance()
			if next.MatchRune('s') && !wordContinuation(text, next.Advance().Pos()) {
				return rules.Match{Start: start, End: next.Advance().Pos(), Type: rules.Word}, true
			}
		}
	}
	match, ok := (rules.WordRule{}).Match(text, start)
	if !ok {
		return match, false
	}
	if rule.ud {
		// Keep hyphenated words together, as in e-mail or by-pass.
		// URL and email rules run at a higher priority.
		for {
			c := rules.NewCursor(text, match.End)
			if !c.MatchRune('-') && !c.MatchRune('‑') {
				break
			}
			part, found := (rules.WordRule{}).Match(text, c.Advance().Pos())
			if !found {
				break
			}
			match.End = part.End
		}
		if strings.EqualFold(text[start:match.End], "e") && strings.HasPrefix(text[match.End:], "/o") && !wordContinuation(text, match.End+2) {
			match.End += 2
		}
		for offset, r := range text[start:match.End] {
			if !apostrophe(r) {
				continue
			}
			if italianElisions[strings.ToLower(text[start:start+offset])] {
				match.End = start + offset + utf8.RuneLen(r)
				return match, true
			}
			if text[start+offset+utf8.RuneLen(r):match.End] == "s" {
				match.End = start + offset
				return match, true
			}
			break
		}
	}
	if match.End < len(text) && (italianTruncations[strings.ToLower(text[start:match.End])] || rule.ud && italianElisions[strings.ToLower(text[start:match.End])]) {
		r, size := utf8.DecodeRuneInString(text[match.End:])
		if apostrophe(r) {
			match.End += size
		}
	}
	if !rule.ud && strings.HasSuffix(strings.ToLower(text[start:match.End]), "'art") && strings.HasPrefix(text[match.End:], ".") && !wordContinuation(text, match.End+1) {
		match.End++
	}
	if match.End == start+1 && text[start] >= 'A' && text[start] <= 'Z' && strings.HasPrefix(text[match.End:], ".") && !wordContinuation(text, match.End+1) {
		match.End++
		match.Type = rules.Abbreviation
	}
	return match, true
}

func wordContinuation(text string, position int) bool {
	c := rules.NewCursor(text, position)
	return c.IsLetter() || c.IsDigit() || c.IsMark()
}

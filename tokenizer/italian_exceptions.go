package tokenizer

import "github.com/espositopasquale/lin-go/rules"

type italianExceptionRule struct{ *rules.AbbreviationRule }

func (italianExceptionRule) Name() string { return "italian-exception" }

func (r italianExceptionRule) Match(text string, start int) (rules.Match, bool) {
	match, ok := r.AbbreviationRule.Match(text, start)
	match.Type = rules.Word
	return match, ok
}

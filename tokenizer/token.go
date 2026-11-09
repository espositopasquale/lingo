package tokenizer

import "github.com/espositopasquale/lingo/rules"

type TokenType = rules.TokenType

const (
	Unknown      = rules.Unknown
	Word         = rules.Word
	Number       = rules.Number
	Date         = rules.Date
	Time         = rules.Time
	Abbreviation = rules.Abbreviation
	Acronym      = rules.Acronym
	Unit         = rules.Unit
	Email        = rules.Email
	URL          = rules.URL
	Punctuation  = rules.Punctuation
	Symbol       = rules.Symbol
)

type Token struct {
	Value string
	Type  TokenType
	Start int
	End   int
}

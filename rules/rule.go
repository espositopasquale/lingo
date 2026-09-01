package rules

type TokenType int

const (
	Unknown TokenType = iota
	Word
	Number
	Date
	Time
	Abbreviation
	Acronym
	Unit
	Email
	URL
	Punctuation
	Symbol
)

type Match struct {
	Start int
	End   int
	Type  TokenType
}

type Rule interface {
	Name() string
	Priority() int
	Match(text string, start int) (Match, bool)
}

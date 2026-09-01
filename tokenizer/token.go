package tokenizer

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

func (t TokenType) String() string {
	switch t {
	case Word:
		return "word"
	case Number:
		return "number"
	case Date:
		return "date"
	case Time:
		return "time"
	case Abbreviation:
		return "abbreviation"
	case Acronym:
		return "acronym"
	case Unit:
		return "unit"
	case Email:
		return "email"
	case URL:
		return "url"
	case Punctuation:
		return "punctuation"
	case Symbol:
		return "symbol"
	default:
		return "unknown"
	}
}

type Token struct {
	Value string
	Type  TokenType
	Start int
	End   int
}

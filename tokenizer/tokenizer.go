package tokenizer

import (
	"unicode"
	"unicode/utf8"

	"github.com/espositopasquale/lin-go/rules"
)

type Tokenizer struct {
	engine  *rules.Engine
	options Options
}

func Tokenize(text string) []Token {
	return NewItalian().Tokenize(text)
}

func TokenizeWithoutPunctuation(text string) []Token {
	options := DefaultOptions()
	options.KeepPunctuation = false

	return NewItalian(options).Tokenize(text)
}

func (t *Tokenizer) Tokenize(text string) []Token {
	tokens := make([]Token, 0)

	for pos := 0; pos < len(text); {
		r, size := utf8.DecodeRuneInString(text[pos:])

		if unicode.IsSpace(r) {
			pos += size
			continue
		}

		match, ok := t.engine.Match(text, pos)

		if !ok {
			tokens = append(tokens, Token{
				Value: text[pos : pos+size],
				Type:  Symbol,
				Start: pos,
				End:   pos + size,
			})

			pos += size
			continue
		}

		token := Token{
			Value: text[match.Start:match.End],
			Type:  match.Type,
			Start: match.Start,
			End:   match.End,
		}

		if token.Type == Punctuation &&
			!t.options.KeepPunctuation {

			pos = match.End
			continue
		}

		tokens = append(tokens, token)

		pos = match.End
	}

	return tokens
}

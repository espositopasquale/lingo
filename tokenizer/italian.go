package tokenizer

import "github.com/espositopasquale/lin-go/rules"

func NewItalian(options ...Options) *Tokenizer {
	return newItalian(false, options...)
}

// NewItalianUD tokenizes Italian text using Universal Dependencies conventions.
// It splits elisions like l'amico into l' and amico, but keeps nella and dammelo
// together. It returns surface tokens, without dependency or morphological analysis.
func NewItalianUD(options ...Options) *Tokenizer {
	return newItalian(true, options...)
}

func newItalian(ud bool, options ...Options) *Tokenizer {
	opts := DefaultOptions()

	if len(options) > 0 {
		opts = options[0]
	}

	engine := rules.NewEngine(
		rules.URLRule{},
		rules.EmailRule{},

		rules.DateRule{},
		rules.TimeRule{},

		italianAcronyms(),
		italianAbbreviations(),
		italianExceptionRule{rules.NewAbbreviationRule(italianWordExceptions)},
		italianUnits(),

		italianNumberRule{ud: ud},
		italianWordRule{ud: ud},
		rules.PunctuationRule{Single: ud},
	)

	return &Tokenizer{
		engine:  engine,
		options: opts,
	}
}

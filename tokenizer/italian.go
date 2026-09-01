package tokenizer

import "github.com/espositopasquale/lin-go/rules"

func NewItalian(options ...Options) *Tokenizer {
	opts := DefaultOptions()

	if len(options) > 0 {
		opts = options[0]
	}

	engine := rules.NewEngine(
		rules.URLRule{},
		rules.EmailRule{},

		rules.DateRule{},
		rules.TimeRule{},

		italianAbbreviations(),
		italianAcronyms(),
		italianUnits(),

		rules.NumberRule{},
		rules.WordRule{},
		rules.PunctuationRule{},
		rules.SymbolRule{},
	)

	return &Tokenizer{
		engine:  engine,
		options: opts,
	}
}

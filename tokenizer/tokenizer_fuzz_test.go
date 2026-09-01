package tokenizer

import (
	"testing"
)

func FuzzTokenizerDoesNotPanic(f *testing.F) {
	seeds := []string{
		"",
		"Ciao",
		"Ciao, mondo!",
		"L'amico è arrivato.",
		"31/08/2026 alle 14:30",
		"1.250,50 €",
		"Perché è così?",
		"🙂 🚀",
		"àèìòù",
		"日本語",
		"\x00",
		"\n\t\r",
	}

	for _, seed := range seeds {
		f.Add(seed)
	}

	f.Fuzz(func(t *testing.T, text string) {
		tokens := Tokenize(text)

		for _, token := range tokens {
			if token.Start < 0 {
				t.Fatalf(
					"Start negativo: %d",
					token.Start,
				)
			}

			if token.End < token.Start {
				t.Fatalf(
					"End < Start: %d < %d",
					token.End,
					token.Start,
				)
			}

			if token.End > len(text) {
				t.Fatalf(
					"End fuori limite: %d > %d",
					token.End,
					len(text),
				)
			}

			if text[token.Start:token.End] != token.Value {
				t.Fatalf(
					"token e offset non corrispondono: %q [%d:%d]",
					token.Value,
					token.Start,
					token.End,
				)
			}
		}
	})
}

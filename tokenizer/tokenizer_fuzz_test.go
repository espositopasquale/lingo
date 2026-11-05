package tokenizer

import (
	"strings"
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
		"Perche\u0301 e\u0300 cosi\u0300",
		"l’UE nell'«opera» e O’Connor",
		"’90, 60x60, 1.5-2.5 e 2,5%",
		"日本語",
		"\x00",
		"\n\t\r",
	}

	for _, seed := range seeds {
		f.Add(seed)
	}

	f.Fuzz(func(t *testing.T, text string) {
		for _, tk := range []*Tokenizer{NewItalian(), NewItalianUD()} {
			tokens := tk.Tokenize(text)
			previous := 0

			for _, token := range tokens {
				if token.Start < 0 {
					t.Fatalf(
						"Start negativo: %d",
						token.Start,
					)
				}

				if token.End <= token.Start {
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
				if token.Start < previous || strings.TrimSpace(text[previous:token.Start]) != "" {
					t.Fatalf("token sovrapposti o testo perso prima di %d", token.Start)
				}
				previous = token.End
			}
			if strings.TrimSpace(text[previous:]) != "" {
				t.Fatalf("testo perso dopo %d", previous)
			}
		}
	})
}

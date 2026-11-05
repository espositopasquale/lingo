package tokenizer

import "testing"

func TestItalianProfiles(t *testing.T) {
	input := "L'amica com’è? Un po' nella sala con O’Connor."
	cases := []struct {
		name string
		tk   *Tokenizer
		want []string
	}{
		{"default", NewItalian(), []string{"L'amica", "com’è", "?", "Un", "po'", "nella", "sala", "con", "O’Connor", "."}},
		{"ud", NewItalianUD(), []string{"L'", "amica", "com’", "è", "?", "Un", "po'", "nella", "sala", "con", "O’Connor", "."}},
	}
	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			tokens := test.tk.Tokenize(input)
			if len(tokens) != len(test.want) {
				t.Fatalf("got %v, want %q", tokens, test.want)
			}
			for index, token := range tokens {
				if token.Value != test.want[index] || input[token.Start:token.End] != token.Value {
					t.Fatalf("token %d: %+v, want %q", index, token, test.want[index])
				}
			}
		})
	}
}

func TestItalianClassificationAndUnicode(t *testing.T) {
	cases := []struct {
		text string
		word string
		kind TokenType
	}{
		{"usa il web ad aprile a casa", "usa", Word},
		{"usa il web ad aprile a casa", "web", Word},
		{"usa il web ad aprile a casa", "ad", Word},
		{"usa il web ad aprile a casa", "a", Word},
		{"L'amico è qui", "L'amico", Word},
		{"lʼamico è qui", "lʼamico", Word},
		{"Perche\u0301 e\u0300 qui", "Perche\u0301", Word},
		{"UE e M5S", "M5S", Acronym},
		{"UE e M5S", "UE", Acronym},
		{"12 cm e 3 m² a 20 °C", "cm", Unit},
		{"12 cm e 3 m² a 20 °C", "m²", Unit},
		{"12 cm e 3 m² a 20 °C", "°C", Unit},
		{"12 cm e 3 m² a 20 °C", "a", Word},
		{"Vai su https://esempio.it.", "https://esempio.it", URL},
		{"Scrivi a nome@mail.esempio.it.", "nome@mail.esempio.it", Email},
	}
	for _, test := range cases {
		t.Run(test.word, func(t *testing.T) {
			for _, token := range NewItalian().Tokenize(test.text) {
				if token.Value == test.word {
					if token.Type != test.kind {
						t.Fatalf("%q: got %s, want %s", test.word, token.Type, test.kind)
					}
					return
				}
			}
			t.Fatalf("token %q missing", test.word)
		})
	}
}

func TestUDWithoutPunctuation(t *testing.T) {
	options := DefaultOptions()
	options.KeepPunctuation = false
	tokens := NewItalianUD(options).Tokenize("L'amica?!")
	if len(tokens) != 2 || tokens[0].Value != "L'" || tokens[1].Value != "amica" {
		t.Fatalf("elision apostrophes must remain in the word token: %+v", tokens)
	}
}

func TestSymbolFallback(t *testing.T) {
	for _, tk := range []*Tokenizer{NewItalian(), NewItalianUD()} {
		for _, text := range []string{"€", "🙂", "\u200d", "\x00", "\xff"} {
			tokens := tk.Tokenize(text)
			if len(tokens) != 1 {
				t.Fatalf("fallback for %q: got %+v", text, tokens)
			}
			token := tokens[0]
			if token.Type != Symbol || token.Value != text || token.Start != 0 || token.End != len(text) {
				t.Fatalf("fallback must preserve symbol type, bytes and offsets for %q: %+v", text, token)
			}
		}
	}
}

func TestUDDevelopmentRegressions(t *testing.T) {
	cases := []struct {
		text string
		want []string
	}{
		{"Dov'è? Cos'è? mezz'ora e trent'anni", []string{"Dov'", "è", "?", "Cos'", "è", "?", "mezz'", "ora", "e", "trent'", "anni"}},
		{"sull'«opera» l'\"amica\"", []string{"sull'", "«", "opera", "»", "l'", "\"", "amica", "\""}},
		{"Negli anni '90 e ’50", []string{"Negli", "anni", "'90", "e", "’50"}},
		{"Al 43' costa L. 22.500", []string{"Al", "43'", "costa", "L.", "22.500"}},
		{"60x60, versione 2.0 e 1.5-2.5", []string{"60x60", ",", "versione", "2.0", "e", "1.5-2.5"}},
		{"by-pass, e-mail e/o king-size", []string{"by-pass", ",", "e-mail", "e/o", "king-size"}},
		{"Martha's Vineyard e King's", []string{"Martha", "'s", "Vineyard", "e", "King", "'s"}},
		{"Il sig. M. Rossi disse: E' qui.", []string{"Il", "sig.", "M.", "Rossi", "disse", ":", "E'", "qui", "."}},
	}
	for _, test := range cases {
		t.Run(test.text, func(t *testing.T) {
			tokens := NewItalianUD().Tokenize(test.text)
			if len(tokens) != len(test.want) {
				t.Fatalf("got %+v; want %q", tokens, test.want)
			}
			for index, token := range tokens {
				if token.Value != test.want[index] || test.text[token.Start:token.End] != token.Value {
					t.Fatalf("got %+v; want %q", token, test.want[index])
				}
			}
		})
	}
}

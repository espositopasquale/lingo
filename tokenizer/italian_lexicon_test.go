package tokenizer

import "testing"

func TestIntegratedItalianLexicons(t *testing.T) {
	cases := []struct {
		text   string
		values []string
		kinds  []TokenType
	}{
		{"CNR C.N.R. C.N.R R.E.M. CdA AD", []string{"CNR", "C.N.R.", "C.N.R", "R.E.M.", "CdA", "AD"}, []TokenType{Acronym, Acronym, Acronym, Acronym, Acronym, Acronym}},
		{"usa ad Ad Usa Cda SCONOSCIUTA", []string{"usa", "ad", "Ad", "Usa", "Cda", "SCONOSCIUTA"}, []TokenType{Word, Word, Word, Word, Word, Word}},
		{"No. 12", []string{"No.", "12"}, []TokenType{Abbreviation, Number}},
		{"No. grazie", []string{"No", ".", "grazie"}, []TokenType{Word, Punctuation, Word}},
		{"Mons. Cost. att. i.e. C.so", []string{"Mons.", "Cost.", "att.", "i.e.", "C.so"}, []TokenType{Abbreviation, Abbreviation, Abbreviation, Abbreviation, Abbreviation}},
		{"S.p.A. s.p.a. S.r.l.", []string{"S.p.A.", "s.p.a.", "S.r.l."}, []TokenType{Acronym, Abbreviation, Acronym}},
		{"by-pass e-mail c/c e/o n° de\"", []string{"by-pass", "e-mail", "c/c", "e/o", "n°", "de\""}, []TokenType{Word, Word, Word, Word, Word, Word}},
		{"70 cm.", []string{"70", "cm", "."}, []TokenType{Number, Unit, Punctuation}},
		{"A. Rossi a. casa", []string{"A.", "Rossi", "a", ".", "casa"}, []TokenType{Abbreviation, Word, Word, Punctuation, Word}},
		{"sett..", []string{"sett.", "."}, []TokenType{Abbreviation, Punctuation}},
	}
	for _, profile := range []struct {
		name string
		tk   *Tokenizer
	}{{"default", NewItalian()}, {"ud", NewItalianUD()}} {
		for _, tc := range cases {
			t.Run(profile.name+"/"+tc.text, func(t *testing.T) {
				got := profile.tk.Tokenize(tc.text)
				if len(got) != len(tc.values) {
					t.Fatalf("got %+v; want %q", got, tc.values)
				}
				for i, token := range got {
					if token.Value != tc.values[i] || token.Type != tc.kinds[i] || tc.text[token.Start:token.End] != token.Value {
						t.Errorf("token %d: got %+v; want %q (%s)", i, token, tc.values[i], tc.kinds[i])
					}
				}
			})
		}
	}
}

func TestImportedElisionsRetainProfileConventions(t *testing.T) {
	for _, prefix := range []string{"all'", "dall'", "dell'", "L'", "l'", "nell'"} {
		text := prefix + "art."
		standard := NewItalian().Tokenize(text)
		ud := NewItalianUD().Tokenize(text)
		if len(standard) != 1 || standard[0].Value != text || len(ud) != 2 || ud[0].Value != prefix || ud[1].Value != "art." || ud[1].Type != Abbreviation {
			t.Errorf("%q: default %+v; UD %+v", text, standard, ud)
		}
	}
}

func TestSigleWithoutPunctuation(t *testing.T) {
	opts := DefaultOptions()
	opts.KeepPunctuation = false
	for _, tk := range []*Tokenizer{NewItalian(opts), NewItalianUD(opts)} {
		got := tk.Tokenize("C.N.R., e-mail!")
		if len(got) != 2 || got[0].Value != "C.N.R." || got[0].Type != Acronym || got[1].Value != "e-mail" {
			t.Fatalf("got %+v", got)
		}
	}
}

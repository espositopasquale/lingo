package tokenizer

import "testing"

func TestUTF8Offsets(t *testing.T) {
	text := "Città, perché è già così."

	tokens := Tokenize(text)

	for _, token := range tokens {
		actual := text[token.Start:token.End]

		if actual != token.Value {
			t.Fatalf(
				"offset errato: token=%q actual=%q [%d:%d]",
				token.Value,
				actual,
				token.Start,
				token.End,
			)
		}
	}
}

func TestItalianWords(t *testing.T) {
	text := "L'amico dell'anno com'è?"

	tokens := TokenizeWithoutPunctuation(text)

	expected := []string{
		"L'amico",
		"dell'anno",
		"com'è",
	}

	if len(tokens) != len(expected) {
		t.Fatalf(
			"attesi %d token, ottenuti %d",
			len(expected),
			len(tokens),
		)
	}

	for i, token := range tokens {
		if token.Value != expected[i] {
			t.Errorf(
				"token %d: atteso %q, ottenuto %q",
				i,
				expected[i],
				token.Value,
			)
		}
	}
}

func TestNumber(t *testing.T) {
	text := "Ho pagato 1.250,50 euro."

	tokens := Tokenize(text)

	for _, token := range tokens {
		if token.Value == "1.250,50" {
			if token.Type != Number {
				t.Fatalf(
					"atteso Number, ottenuto %s",
					token.Type,
				)
			}

			return
		}
	}

	t.Fatal("numero non trovato")
}

func TestDate(t *testing.T) {
	text := "Oggi è il 31/08/2026."

	tokens := Tokenize(text)

	for _, token := range tokens {
		if token.Value == "31/08/2026" {
			if token.Type != Date {
				t.Fatalf(
					"atteso Date, ottenuto %s",
					token.Type,
				)
			}

			return
		}
	}

	t.Fatal("data non trovata")
}

func TestTime(t *testing.T) {
	text := "L'appuntamento è alle 14:30."

	tokens := Tokenize(text)

	for _, token := range tokens {
		if token.Value == "14:30" {
			if token.Type != Time {
				t.Fatalf(
					"atteso Time, ottenuto %s",
					token.Type,
				)
			}

			return
		}
	}

	t.Fatal("orario non trovato")
}

func TestWithoutPunctuation(t *testing.T) {
	text := "Ciao, come stai?!"

	tokens := TokenizeWithoutPunctuation(text)

	expected := []string{
		"Ciao",
		"come",
		"stai",
	}

	if len(tokens) != len(expected) {
		t.Fatalf(
			"attesi %d token, ottenuti %d",
			len(expected),
			len(tokens),
		)
	}

	for i, token := range tokens {
		if token.Value != expected[i] {
			t.Errorf(
				"atteso %q, ottenuto %q",
				expected[i],
				token.Value,
			)
		}
	}
}

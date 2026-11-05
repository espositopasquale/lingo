package conllu

import (
	"errors"
	"io"
	"strings"
	"testing"
)

func row(id, form, misc string) string {
	return strings.Join([]string{id, form, "_", "X", "_", "_", "_", "_", "_", misc}, "\t") + "\n"
}

func TestSurfaceTokensAndUTF8Offsets(t *testing.T) {
	input := "# sent_id = mwt\n# text = È nella città.  \n" +
		row("1", "È", "_") + row("2-3", "nella", "_") +
		row("2", "in", "_") + row("3", "la", "_") +
		row("3.1", "fantasma", "_") + row("4", "città", "SpaceAfter=No") + row("5", ".", "_")
	reader := NewReader(strings.NewReader(input))
	sentence, err := reader.Read()
	if err != nil {
		t.Fatal(err)
	}
	if sentence.ID != "mwt" || sentence.Words != 5 || len(sentence.Tokens) != 4 || sentence.Text != "È nella città.  " {
		t.Fatalf("unexpected sentence: %+v", sentence)
	}
	if sentence.Tokens[1].ID != "2-3" || sentence.Tokens[1].Start != 3 || sentence.Tokens[1].End != 8 {
		t.Fatalf("MWT must use the original surface span: %+v", sentence.Tokens[1])
	}
	for _, token := range sentence.Tokens {
		if sentence.Text[token.Start:token.End] != token.Form {
			t.Fatalf("incorrect UTF-8 offsets: %+v", token)
		}
	}
	if _, err := reader.Read(); err != io.EOF {
		t.Fatalf("expected EOF, got %v", err)
	}
}

func TestReconstructionAndSentenceBoundaries(t *testing.T) {
	input := "# sent_id = reconstructed\n" +
		row("1-2", "al", "SpacesAfter=\\s\\s") + row("1", "a", "_") + row("2", "il", "_") +
		row("3", "bar", "SpaceAfter=No") + row("4", ".", "_") + "\n\n" +
		"# text = Ciao\u00a0mondo!\n" + row("1", "Ciao", "_") + row("2", "mondo", "SpaceAfter=No") + row("3", "!", "_")
	reader := NewReader(strings.NewReader(input))
	first, err := reader.Read()
	if err != nil || first.Text != "al  bar." || len(first.Tokens) != 3 {
		t.Fatalf("reconstruction: %+v, %v", first, err)
	}
	second, err := reader.Read()
	if err != nil || second.Text != "Ciao\u00a0mondo!" || second.Tokens[1].Start != 6 {
		t.Fatalf("second sentence: %+v, %v", second, err)
	}
}

func TestRejectMalformedCorpus(t *testing.T) {
	cases := map[string]string{
		"columns":        "1\tciao\n",
		"missing-id":     row("2", "ciao", "_"),
		"zero-id":        row("0", "ciao", "_"),
		"leading-zero":   row("01", "ciao", "_"),
		"empty-node":     row("1.x", "ciao", "_"),
		"bad-range":      row("1-1", "ciao", "_"),
		"missing-parts":  row("1-2", "nella", "_") + row("1", "in", "_"),
		"overlap":        row("1-3", "nella", "_") + row("1", "in", "_") + row("2-3", "la", "_"),
		"form-mismatch":  "# text = Ciao!\n" + row("1", "ciao", "_"),
		"omitted-text":   "# text = ciao mondo\n" + row("1", "ciao", "_"),
		"duplicate-text": "# text = ciao\n# text = ciao\n" + row("1", "ciao", "_"),
		"bad-spacing":    row("1", "ciao", "SpacesAfter=\\x"),
		"invalid-utf8":   row("1", "\xff", "_"),
	}
	for name, input := range cases {
		t.Run(name, func(t *testing.T) {
			if _, err := NewReader(strings.NewReader(input)).Read(); err == nil {
				t.Fatal("malformed corpus must fail instead of being silently skipped")
			}
		})
	}
}

type failingReader struct{}

func (failingReader) Read([]byte) (int, error) { return 0, errors.New("read failed") }

func TestInputError(t *testing.T) {
	if _, err := NewReader(failingReader{}).Read(); err == nil || !strings.Contains(err.Error(), "read failed") {
		t.Fatalf("input errors must propagate: %v", err)
	}
}

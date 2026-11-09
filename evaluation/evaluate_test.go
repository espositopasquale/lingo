package evaluation

import (
	"crypto/sha256"
	"fmt"
	"io"
	"math"
	"os"
	"strings"
	"testing"

	"github.com/espositopasquale/lingo/conllu"
	"github.com/espositopasquale/lingo/tokenizer"
)

func fixture(t testing.TB) []byte {
	t.Helper()
	data, err := os.ReadFile("testdata/italian.conllu")
	if err != nil {
		t.Fatal(err)
	}
	return data
}

func TestItalianRegressionSuite(t *testing.T) {
	data := fixture(t)
	report, err := Evaluate(strings.NewReader(string(data)), tokenizer.NewItalianUD().Tokenize, 10)
	if err != nil {
		t.Fatal(err)
	}
	if report.Sentences != 22 || report.ExactSentences != 22 || report.Tokens.F1 != 1 {
		t.Fatalf("Italian regression suite failed: %+v", report)
	}
	if report.CorpusSHA256 != fmt.Sprintf("%x", sha256.Sum256(data)) {
		t.Fatal("report must identify the exact input bytes")
	}
}

func TestExactSpanScoresAndDiagnosticsLimit(t *testing.T) {
	// Merge the first two reference tokens to check how a boundary error is scored.
	// Only the final period matches: precision=1/2, recall=1/3, F1=2/5.
	input := "# sent_id = score\n# text = L'amica.\n" +
		"1\tL'\t_\tDET\t_\t_\t2\tdet\t_\tSpaceAfter=No\n" +
		"2\tamica\t_\tNOUN\t_\t_\t0\troot\t_\tSpaceAfter=No\n" +
		"3\t.\t_\tPUNCT\t_\t_\t2\tpunct\t_\t_\n\n"
	report, err := Evaluate(strings.NewReader(input+input), tokenizer.NewItalian().Tokenize, 1)
	if err != nil {
		t.Fatal(err)
	}
	if report.Sentences != 2 || len(report.Examples) != 1 || report.ExactSentences != 0 {
		t.Fatalf("diagnostic limit must not truncate evaluation: %+v", report)
	}
	if math.Abs(report.Tokens.Precision-.5) > 1e-12 || math.Abs(report.Tokens.Recall-1.0/3) > 1e-12 || math.Abs(report.Tokens.F1-.4) > 1e-12 {
		t.Fatalf("incorrect exact-span metrics: %+v", report.Tokens)
	}
}

func TestRejectInvalidPredictions(t *testing.T) {
	input := "# text = È.\n1\tÈ\t_\tX\t_\t_\t0\troot\t_\tSpaceAfter=No\n2\t.\t_\tPUNCT\t_\t_\t1\tpunct\t_\t_\n"
	cases := map[string][]tokenizer.Token{
		"missing":       {},
		"wrong-value":   {{Value: "E", Start: 0, End: 2}},
		"out-of-bounds": {{Value: "È.", Start: 0, End: 4}},
		"overlap":       {{Value: "È", Start: 0, End: 2}, {Value: "È.", Start: 0, End: 3}},
		"utf8-split":    {{Value: "\xc3", Start: 0, End: 1}, {Value: "\x88.", Start: 1, End: 3}},
	}
	for name, prediction := range cases {
		t.Run(name, func(t *testing.T) {
			_, err := Evaluate(strings.NewReader(input), func(string) []tokenizer.Token { return prediction }, 0)
			if err == nil {
				t.Fatal("invalid predictions must not produce a score")
			}
		})
	}
	if _, err := Evaluate(strings.NewReader(""), tokenizer.Tokenize, 0); err == nil {
		t.Fatal("empty corpus must fail")
	}
	if _, err := Evaluate(strings.NewReader(input), nil, 0); err == nil {
		t.Fatal("nil tokenizer must fail")
	}
	if _, err := Evaluate(strings.NewReader(input), tokenizer.Tokenize, -1); err == nil {
		t.Fatal("negative diagnostic limit must fail")
	}
}

func BenchmarkItalianProfiles(b *testing.B) {
	reader := conllu.NewReader(strings.NewReader(string(fixture(b))))
	var parts []string
	for {
		sentence, err := reader.Read()
		if err != nil {
			if err != io.EOF {
				b.Fatal(err)
			}
			break
		}
		parts = append(parts, sentence.Text)
	}
	text := strings.Join(parts, "\n")
	profiles := map[string]*tokenizer.Tokenizer{"default": tokenizer.NewItalian(), "ud": tokenizer.NewItalianUD()}
	for name, tk := range profiles {
		b.Run(name, func(b *testing.B) {
			b.SetBytes(int64(len(text)))
			b.ReportAllocs()
			for b.Loop() {
				tk.Tokenize(text)
			}
		})
	}
}

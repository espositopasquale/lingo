package main

import (
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"math"
	"os"
	"strings"

	"github.com/espositopasquale/lingo/evaluation"
	"github.com/espositopasquale/lingo/tokenizer"
)

func main() { os.Exit(run(os.Args[1:], os.Stdout, os.Stderr)) }

func run(args []string, stdout, stderr io.Writer) int {
	flags := flag.NewFlagSet("lingo-eval", flag.ContinueOnError)
	flags.SetOutput(stderr)
	corpus := flags.String("corpus", "", "local CoNLL-U file (required)")
	profile := flags.String("profile", "ud", "Italian tokenization profile: ud or default")
	format := flags.String("format", "text", "report format: text or json")
	maxErrors := flags.Int("max-errors", 5, "maximum diagnostic examples; all sentences are scored")
	minimum := flags.Float64("min-f1", 0, "fail below this token F1, in [0,1]")
	expected := flags.String("expect-sha256", "", "require this corpus SHA-256 for reproducibility")
	if err := flags.Parse(args); err != nil {
		if err == flag.ErrHelp {
			return 0
		}
		return 2
	}
	if *corpus == "" || flags.NArg() != 0 || *maxErrors < 0 || math.IsNaN(*minimum) || *minimum < 0 || *minimum > 1 {
		fmt.Fprintln(stderr, "error: supply -corpus, nonnegative -max-errors and -min-f1 in [0,1]")
		return 2
	}
	if *format != "text" && *format != "json" {
		fmt.Fprintln(stderr, "error: -format must be text or json")
		return 2
	}
	if *expected != "" {
		decoded, err := hex.DecodeString(*expected)
		if err != nil || len(decoded) != 32 {
			fmt.Fprintln(stderr, "error: -expect-sha256 must be a 64-character hexadecimal digest")
			return 2
		}
	}
	var tk *tokenizer.Tokenizer
	switch *profile {
	case "ud":
		tk = tokenizer.NewItalianUD()
	case "default":
		tk = tokenizer.NewItalian()
	default:
		fmt.Fprintln(stderr, "error: -profile must be ud or default")
		return 2
	}
	file, err := os.Open(*corpus)
	if err != nil {
		fmt.Fprintln(stderr, "error:", err)
		return 1
	}
	defer file.Close()
	report, err := evaluation.Evaluate(file, tk.Tokenize, *maxErrors)
	if err != nil {
		fmt.Fprintln(stderr, "error:", err)
		return 1
	}
	if *expected != "" && !strings.EqualFold(*expected, report.CorpusSHA256) {
		fmt.Fprintf(stderr, "error: corpus SHA-256 mismatch: got %s\n", report.CorpusSHA256)
		return 1
	}
	if *format == "json" {
		encoder := json.NewEncoder(stdout)
		encoder.SetIndent("", "  ")
		if err := encoder.Encode(struct {
			Profile string `json:"profile"`
			evaluation.Report
		}{*profile, report}); err != nil {
			fmt.Fprintln(stderr, "error writing report:", err)
			return 1
		}
	} else {
		fmt.Fprintf(stdout, "Profile: %s\nCorpus SHA-256: %s\n", *profile, report.CorpusSHA256)
		fmt.Fprintf(stdout, "Sentences: %d; exact: %d (%.2f%%)\n", report.Sentences, report.ExactSentences, report.SentenceAccuracy*100)
		fmt.Fprintf(stdout, "Tokens: gold=%d predicted=%d correct=%d\n", report.Tokens.Gold, report.Tokens.Predicted, report.Tokens.Correct)
		fmt.Fprintf(stdout, "Precision: %.4f  Recall: %.4f  F1: %.4f\n", report.Tokens.Precision, report.Tokens.Recall, report.Tokens.F1)
		for _, example := range report.Examples {
			fmt.Fprintf(stdout, "\n%s\n  gold: %q\n  got:  %q\n", example.SentenceID, example.Gold, example.Predicted)
		}
	}
	if report.Tokens.F1 < *minimum {
		fmt.Fprintf(stderr, "error: token F1 %.6f is below required %.6f\n", report.Tokens.F1, *minimum)
		return 1
	}
	return 0
}

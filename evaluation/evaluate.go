// Package evaluation compares token boundaries with the reference in a CoNLL-U file.
// It uses the corpus sentence boundaries and includes punctuation. Token types,
// dependency trees and the words inside multiword tokens are not scored.
package evaluation

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"unicode"
	"unicode/utf8"

	"github.com/espositopasquale/lingo/conllu"
	"github.com/espositopasquale/lingo/tokenizer"
)

type TokenizeFunc func(string) []tokenizer.Token

type Scores struct {
	Gold      int     `json:"gold"`
	Predicted int     `json:"predicted"`
	Correct   int     `json:"correct"`
	Precision float64 `json:"precision"`
	Recall    float64 `json:"recall"`
	F1        float64 `json:"f1"`
}

type Span struct {
	Start int `json:"start"`
	End   int `json:"end"`
}

type Mismatch struct {
	SentenceID string   `json:"sentence_id"`
	Text       string   `json:"text"`
	Gold       []string `json:"gold"`
	Predicted  []string `json:"predicted"`
}

type Report struct {
	CorpusSHA256     string     `json:"corpus_sha256"`
	Sentences        int        `json:"sentences"`
	ExactSentences   int        `json:"exact_sentences"`
	SentenceAccuracy float64    `json:"sentence_accuracy"`
	Tokens           Scores     `json:"tokens"`
	Examples         []Mismatch `json:"examples,omitempty"`
}

// Evaluate scores every sentence by comparing exact token boundaries.
// Precision, recall and F1 use totals across the corpus and range from 0 to 1;
// a zero denominator gives a score of zero. maxErrors limits the examples saved
// in the report, not the sentences scored. Tokens must preserve the original text
// and cover every byte except whitespace.
func Evaluate(input io.Reader, tokenize TokenizeFunc, maxErrors int) (Report, error) {
	if tokenize == nil || maxErrors < 0 {
		return Report{}, fmt.Errorf("tokenizer is required and maxErrors must be nonnegative")
	}
	var report Report
	hash := sha256.New()
	reader := conllu.NewReader(io.TeeReader(input, hash))
	for {
		sentence, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return Report{}, err
		}
		predicted := tokenize(sentence.Text)
		if err := validate(sentence.Text, predicted); err != nil {
			return Report{}, fmt.Errorf("sentence %q: %w", sentence.ID, err)
		}
		gold := make(map[Span]bool, len(sentence.Tokens))
		for _, token := range sentence.Tokens {
			gold[Span{token.Start, token.End}] = true
		}
		correct := 0
		for _, token := range predicted {
			if gold[Span{token.Start, token.End}] {
				correct++
			}
		}
		report.Sentences++
		report.Tokens.Gold += len(sentence.Tokens)
		report.Tokens.Predicted += len(predicted)
		report.Tokens.Correct += correct
		if correct == len(sentence.Tokens) && correct == len(predicted) {
			report.ExactSentences++
		} else if len(report.Examples) < maxErrors {
			mismatch := Mismatch{SentenceID: sentence.ID, Text: sentence.Text}
			for _, token := range sentence.Tokens {
				mismatch.Gold = append(mismatch.Gold, token.Form)
			}
			for _, token := range predicted {
				mismatch.Predicted = append(mismatch.Predicted, token.Value)
			}
			report.Examples = append(report.Examples, mismatch)
		}
	}
	if report.Sentences == 0 {
		return Report{}, fmt.Errorf("corpus contains no sentences")
	}
	report.CorpusSHA256 = hex.EncodeToString(hash.Sum(nil))
	report.SentenceAccuracy = float64(report.ExactSentences) / float64(report.Sentences)
	scores := &report.Tokens
	if scores.Predicted > 0 {
		scores.Precision = float64(scores.Correct) / float64(scores.Predicted)
	}
	if scores.Gold > 0 {
		scores.Recall = float64(scores.Correct) / float64(scores.Gold)
	}
	if scores.Gold+scores.Predicted > 0 {
		scores.F1 = 2 * float64(scores.Correct) / float64(scores.Gold+scores.Predicted)
	}
	return report, nil
}

func validate(text string, tokens []tokenizer.Token) error {
	previous := 0
	for _, token := range tokens {
		if token.Start < previous || token.End <= token.Start || token.End > len(text) {
			return fmt.Errorf("invalid or overlapping predicted span [%d:%d]", token.Start, token.End)
		}
		if text[token.Start:token.End] != token.Value || !utf8.ValidString(token.Value) {
			return fmt.Errorf("prediction %q does not preserve its UTF-8 text span", token.Value)
		}
		if !spaceOnly(text[previous:token.Start]) {
			return fmt.Errorf("prediction omitted non-whitespace text before byte %d", token.Start)
		}
		if spaceOnly(token.Value) {
			return fmt.Errorf("prediction contains a whitespace-only token")
		}
		previous = token.End
	}
	if !spaceOnly(text[previous:]) {
		return fmt.Errorf("prediction omitted non-whitespace text after byte %d", previous)
	}
	return nil
}

func spaceOnly(text string) bool {
	for _, r := range text {
		if !unicode.IsSpace(r) {
			return false
		}
	}
	return true
}

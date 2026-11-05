package main

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
)

const corpus = "../../evaluation/testdata/italian.conllu"

func TestJSONAndRegressionGate(t *testing.T) {
	var stdout, stderr bytes.Buffer
	args := []string{"-corpus", corpus, "-format", "json", "-min-f1", "1"}
	if code := run(args, &stdout, &stderr); code != 0 {
		t.Fatalf("exit %d: %s", code, stderr.String())
	}
	var report struct {
		Profile      string `json:"profile"`
		CorpusSHA256 string `json:"corpus_sha256"`
		Sentences    int    `json:"sentences"`
	}
	if err := json.Unmarshal(stdout.Bytes(), &report); err != nil || report.Profile != "ud" || report.Sentences != 22 {
		t.Fatalf("invalid JSON report: %s, %v", stdout.String(), err)
	}
	first := stdout.String()
	stdout.Reset()
	if code := run(append(args, "-expect-sha256", report.CorpusSHA256), &stdout, &stderr); code != 0 || stdout.String() != first {
		t.Fatalf("same input/configuration must produce identical JSON, exit %d: %s", code, stderr.String())
	}
	stdout.Reset()
	stderr.Reset()
	if code := run([]string{"-corpus", corpus, "-profile", "default", "-min-f1", "1"}, &stdout, &stderr); code != 1 || !strings.Contains(stderr.String(), "below required") {
		t.Fatalf("regression gate must fail: code=%d stderr=%s", code, stderr.String())
	}
}

func TestInvalidArgumentsAndCorpus(t *testing.T) {
	cases := [][]string{
		{}, {"-corpus", corpus, "-profile", "bad"}, {"-corpus", corpus, "-format", "xml"},
		{"-corpus", corpus, "-min-f1", "NaN"}, {"-corpus", corpus, "-min-f1", "2"},
		{"-corpus", corpus, "-max-errors", "-1"}, {"-corpus", corpus, "-expect-sha256", "bad"},
		{"-corpus", "does-not-exist"}, {"-corpus", corpus, "-expect-sha256", strings.Repeat("0", 64)},
	}
	for _, args := range cases {
		var stdout, stderr bytes.Buffer
		if code := run(args, &stdout, &stderr); code == 0 || stdout.Len() != 0 || stderr.Len() == 0 {
			t.Fatalf("must fail without a misleading report: %q, exit %d, stdout=%s stderr=%s", args, code, stdout.String(), stderr.String())
		}
	}
}

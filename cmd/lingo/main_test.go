package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"strings"
	"testing"

	"github.com/espositopasquale/lingo/tokenizer"
)

func TestRunTextAndJSON(t *testing.T) {
	var stdout, stderr bytes.Buffer
	code := run([]string{"-ud", "-format", "json", "-text", "L'amico!"}, strings.NewReader(""), &stdout, &stderr, false)
	if code != 0 || stderr.Len() != 0 {
		t.Fatalf("exit %d: %s", code, stderr.String())
	}
	var tokens []struct {
		Value string `json:"value"`
		Type  string `json:"type"`
	}
	if err := json.Unmarshal(stdout.Bytes(), &tokens); err != nil {
		t.Fatalf("invalid JSON: %v (%s)", err, stdout.String())
	}
	if len(tokens) != 3 || tokens[0].Value != "L'" || tokens[2].Type != tokenizer.Punctuation.String() {
		t.Fatalf("unexpected tokens: %#v", tokens)
	}
}

func TestInputPrecedenceAndNoPunctuation(t *testing.T) {
	var stdout, stderr bytes.Buffer
	args := []string{"-no-punct", "-text", "Scelto!", "ignorato"}
	if code := run(args, strings.NewReader("anche questo"), &stdout, &stderr, true); code != 0 {
		t.Fatalf("exit %d: %s", code, stderr.String())
	}
	if got := stdout.String(); !strings.Contains(got, `"Scelto"`) || strings.Contains(got, `"!"`) || strings.Contains(got, "ignorato") {
		t.Fatalf("unexpected output: %s", got)
	}
}

func TestInvalidArgumentsAndReadError(t *testing.T) {
	for _, args := range [][]string{{"-format", "xml"}, {"-unknown"}} {
		var stdout, stderr bytes.Buffer
		if code := run(args, strings.NewReader(""), &stdout, &stderr, false); code != 2 || stderr.Len() == 0 {
			t.Fatalf("args %q: exit=%d stderr=%q", args, code, stderr.String())
		}
	}
	var stdout, stderr bytes.Buffer
	if code := run(nil, errorReader{}, &stdout, &stderr, true); code != 1 || !strings.Contains(stderr.String(), "error reading stdin") {
		t.Fatalf("exit=%d stderr=%q", code, stderr.String())
	}
}

func TestVersion(t *testing.T) {
	old := version
	version = "v1.0.0-test"
	t.Cleanup(func() { version = old })
	var stdout, stderr bytes.Buffer
	if code := run([]string{"-version"}, strings.NewReader(""), &stdout, &stderr, false); code != 0 || stdout.String() != "v1.0.0-test\n" {
		t.Fatalf("exit=%d stdout=%q stderr=%q", code, stdout.String(), stderr.String())
	}
}

func TestREPLAcceptsLongLines(t *testing.T) {
	line := strings.Repeat("a", 70*1024)
	var stdout, stderr bytes.Buffer
	if code := run(nil, strings.NewReader(line+"\nquit\n"), &stdout, &stderr, false); code != 0 {
		t.Fatalf("exit=%d stderr=%q", code, stderr.String())
	}
}

type errorReader struct{}

func (errorReader) Read([]byte) (int, error) { return 0, errors.New("read failed") }

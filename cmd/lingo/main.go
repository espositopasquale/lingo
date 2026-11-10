package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"runtime/debug"
	"strings"

	"github.com/espositopasquale/lingo/tokenizer"
)

var version = "dev"

func main() { os.Exit(run(os.Args[1:], os.Stdin, os.Stdout, os.Stderr, stdinHasData(os.Stdin))) }

func run(args []string, stdin io.Reader, stdout, stderr io.Writer, hasStdin bool) int {
	flags := flag.NewFlagSet("lingo", flag.ContinueOnError)
	flags.SetOutput(stderr)
	noPunct := flags.Bool("no-punct", false, "strip punctuation tokens from the output")
	textFlag := flags.String("text", "", "text to tokenize (skips stdin/args)")
	ud := flags.Bool("ud", false, "use Italian Universal Dependencies surface-token conventions")
	format := flags.String("format", "text", "output format: text or json")
	showVersion := flags.Bool("version", false, "print version and exit")
	if err := flags.Parse(args); err != nil {
		if err == flag.ErrHelp {
			return 0
		}
		return 2
	}
	if *showVersion {
		fmt.Fprintln(stdout, buildVersion())
		return 0
	}
	if *format != "text" && *format != "json" {
		fmt.Fprintln(stderr, "error: -format must be text or json")
		return 2
	}

	options := tokenizer.DefaultOptions()
	options.KeepPunctuation = !*noPunct
	tk := tokenizer.NewItalian(options)
	if *ud {
		tk = tokenizer.NewItalianUD(options)
	}

	var text string
	switch {
	case *textFlag != "":
		text = *textFlag
	case flags.NArg() > 0:
		text = strings.Join(flags.Args(), " ")
	case hasStdin:
		contents, err := io.ReadAll(stdin)
		if err != nil {
			fmt.Fprintln(stderr, "error reading stdin:", err)
			return 1
		}
		text = string(contents)
	default:
		return repl(stdin, stdout, stderr, tk, *format)
	}

	if err := writeTokens(stdout, tk.Tokenize(text), len(text), *format); err != nil {
		fmt.Fprintln(stderr, "error writing output:", err)
		return 1
	}
	return 0
}

func stdinHasData(file *os.File) bool {
	stat, err := file.Stat()
	return err == nil && (stat.Mode()&os.ModeCharDevice) == 0
}

func repl(stdin io.Reader, stdout, stderr io.Writer, tk *tokenizer.Tokenizer, format string) int {
	scanner := bufio.NewScanner(stdin)
	// Scanner's 64 KiB default is surprisingly small for pasted documents.
	scanner.Buffer(make([]byte, 64*1024), 10*1024*1024)
	fmt.Fprintln(stdout, "LinGo tokenizer. Type a sentence and press Enter (or type 'exit' to quit).")
	for {
		fmt.Fprint(stdout, "> ")
		if !scanner.Scan() {
			break
		}
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		if line == "exit" || line == "quit" {
			break
		}
		if err := writeTokens(stdout, tk.Tokenize(line), len(line), format); err != nil {
			fmt.Fprintln(stderr, "error writing output:", err)
			return 1
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(stderr, "error reading input:", err)
		return 1
	}
	return 0
}

func writeTokens(writer io.Writer, tokens []tokenizer.Token, textLength int, format string) error {
	if format == "json" {
		type jsonToken struct {
			Value string `json:"value"`
			Type  string `json:"type"`
			Start int    `json:"start"`
			End   int    `json:"end"`
		}
		output := make([]jsonToken, len(tokens))
		for i, token := range tokens {
			output[i] = jsonToken{token.Value, token.Type.String(), token.Start, token.End}
		}
		encoder := json.NewEncoder(writer)
		encoder.SetEscapeHTML(false)
		return encoder.Encode(output)
	}
	offsetWidth := len(fmt.Sprintf("%d", textLength))
	valueWidth := 0
	for _, token := range tokens {
		if width := len(fmt.Sprintf("%q", token.Value)); width > valueWidth {
			valueWidth = width
		}
	}
	for _, token := range tokens {
		if _, err := fmt.Fprintf(writer, "%-*s | %-13s | [%*d:%*d]\n",
			valueWidth, fmt.Sprintf("%q", token.Value), token.Type,
			offsetWidth, token.Start, offsetWidth, token.End); err != nil {
			return err
		}
	}
	return nil
}

func buildVersion() string {
	if version != "dev" {
		return version
	}
	if info, ok := debug.ReadBuildInfo(); ok && info.Main.Version != "" && info.Main.Version != "(devel)" {
		return info.Main.Version
	}
	return version
}

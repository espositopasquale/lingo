package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/espositopasquale/lin-go/tokenizer"
)

func main() {
	noPunct := flag.Bool("no-punct", false, "strip punctuation tokens from the output")
	textFlag := flag.String("text", "", "text to tokenize (skips stdin/args)")
	flag.Parse()

	options := tokenizer.DefaultOptions()
	options.KeepPunctuation = !*noPunct

	switch {
	case *textFlag != "":
		tokenize(*textFlag, options)
	case flag.NArg() > 0:
		tokenize(strings.Join(flag.Args(), " "), options)
	case stdinHasData():
		text, err := io.ReadAll(os.Stdin)
		if err != nil {
			fmt.Fprintln(os.Stderr, "error reading stdin:", err)
			os.Exit(1)
		}
		tokenize(string(text), options)
	default:
		repl(options)
	}
}

func stdinHasData() bool {
	stat, err := os.Stdin.Stat()
	if err != nil {
		return false
	}

	return (stat.Mode() & os.ModeCharDevice) == 0
}

func repl(options tokenizer.Options) {
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Println("Lin-go tokenizer. Type a sentence and press Enter (or type 'exit' to quit).")

	for {
		fmt.Print("> ")

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

		tokenize(line, options)
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "error reading input:", err)
		os.Exit(1)
	}
}

func tokenize(text string, options tokenizer.Options) {
	tokens := tokenizer.NewItalian(options).Tokenize(text)

	offsetWidth := len(fmt.Sprintf("%d", len(text)))
	valueWidth := 0

	for _, token := range tokens {
		if width := len(fmt.Sprintf("%q", token.Value)); width > valueWidth {
			valueWidth = width
		}
	}

	for _, token := range tokens {
		fmt.Printf(
			"%-*s | %-13s | [%*d:%*d]\n",
			valueWidth, fmt.Sprintf("%q", token.Value),
			token.Type,
			offsetWidth, token.Start,
			offsetWidth, token.End,
		)
	}
}

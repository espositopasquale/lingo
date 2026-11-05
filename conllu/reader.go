// Package conllu reads surface tokens from Universal Dependencies CoNLL-U files.
// A multiword token uses the FORM from its range row; its component rows and empty
// nodes are skipped. Dependency graphs are not validated.
package conllu

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"
)

// Token holds a surface form and its byte offsets in the original UTF-8 sentence.
type Token struct {
	ID    string
	Form  string
	Start int
	End   int
}

type Sentence struct {
	ID     string
	Text   string
	Tokens []Token
	// Words counts rows with integer IDs, including words inside multiword tokens.
	Words int
}

type Reader struct {
	scanner *bufio.Scanner
	line    int
}

func NewReader(input io.Reader) *Reader {
	scanner := bufio.NewScanner(input)
	scanner.Buffer(make([]byte, 4096), 16*1024*1024)
	return &Reader{scanner: scanner}
}

// Read returns the next sentence, or io.EOF at the end of the file.
// Invalid IDs, incomplete or overlapping multiword tokens, and forms that do not
// match the text return an error. When # text is missing, Read rebuilds the text
// using SpaceAfter and SpacesAfter.
func (reader *Reader) Read() (Sentence, error) {
	var sentence Sentence
	var after []string
	hasText := false
	nextWord, coveredEnd := 1, 0
	for {
		ok := reader.scanner.Scan()
		if !ok {
			if err := reader.scanner.Err(); err != nil {
				return Sentence{}, fmt.Errorf("CoNLL-U line %d: %w", reader.line+1, err)
			}
			if len(sentence.Tokens) == 0 {
				return Sentence{}, io.EOF
			}
		} else {
			reader.line++
		}
		line := strings.TrimSuffix(reader.scanner.Text(), "\r")
		if !ok || line == "" {
			if len(sentence.Tokens) == 0 {
				sentence = Sentence{}
				hasText = false
				continue
			}
			if coveredEnd >= nextWord {
				return Sentence{}, reader.failure(sentence.ID, "incomplete multiword token")
			}
			if !hasText {
				var text strings.Builder
				for index, token := range sentence.Tokens {
					text.WriteString(token.Form)
					if index < len(sentence.Tokens)-1 {
						text.WriteString(after[index])
					}
				}
				sentence.Text = text.String()
			}
			if err := align(&sentence); err != nil {
				return Sentence{}, reader.failure(sentence.ID, err.Error())
			}
			return sentence, nil
		}
		if !utf8.ValidString(line) {
			return Sentence{}, reader.failure(sentence.ID, "invalid UTF-8")
		}
		if strings.HasPrefix(line, "#") {
			comment := strings.TrimLeft(strings.TrimPrefix(line, "#"), " \t")
			key, value, found := strings.Cut(comment, " = ")
			if found && key == "sent_id" {
				sentence.ID = value
			} else if found && key == "text" {
				if hasText {
					return Sentence{}, reader.failure(sentence.ID, "duplicate # text")
				}
				sentence.Text, hasText = value, true
			}
			continue
		}
		fields := strings.Split(line, "\t")
		if len(fields) != 10 {
			return Sentence{}, reader.failure(sentence.ID, "expected 10 tab-separated columns")
		}
		for _, field := range fields {
			if field == "" {
				return Sentence{}, reader.failure(sentence.ID, "empty field; use _ for missing values")
			}
		}
		id := fields[0]
		if base, decimal, found := strings.Cut(id, "."); found {
			anchor, e1 := integer(base, true)
			_, e2 := integer(decimal, false)
			if e1 != nil || e2 != nil || anchor != nextWord-1 {
				return Sentence{}, reader.failure(sentence.ID, "invalid empty-node ID "+id)
			}
			continue
		}
		if first, last, found := strings.Cut(id, "-"); found {
			start, e1 := integer(first, false)
			end, e2 := integer(last, false)
			if e1 != nil || e2 != nil || start != nextWord || end <= start || coveredEnd >= start {
				return Sentence{}, reader.failure(sentence.ID, "invalid or overlapping multiword-token ID "+id)
			}
			coveredEnd = end
		} else {
			number, err := integer(id, false)
			if err != nil || number != nextWord {
				return Sentence{}, reader.failure(sentence.ID, "expected consecutive word ID "+strconv.Itoa(nextWord)+", got "+id)
			}
			nextWord++
			sentence.Words++
			if number <= coveredEnd {
				continue
			}
		}
		spacing, err := spaceAfter(fields[9])
		if err != nil {
			return Sentence{}, reader.failure(sentence.ID, err.Error())
		}
		sentence.Tokens = append(sentence.Tokens, Token{ID: id, Form: fields[1]})
		after = append(after, spacing)
	}
}

func (reader *Reader) failure(id, message string) error {
	return fmt.Errorf("CoNLL-U line %d, sentence %q: %s", reader.line, id, message)
}

func integer(value string, zero bool) (int, error) {
	number, err := strconv.Atoi(value)
	if err != nil || number < 0 || (!zero && number == 0) || strconv.Itoa(number) != value {
		return 0, fmt.Errorf("invalid integer %q", value)
	}
	return number, nil
}

func spaceAfter(misc string) (string, error) {
	space := " "
	for _, item := range strings.Split(misc, "|") {
		if item == "SpaceAfter=No" {
			space = ""
		}
	}
	for _, item := range strings.Split(misc, "|") {
		value, found := strings.CutPrefix(item, "SpacesAfter=")
		if !found {
			continue
		}
		var decoded strings.Builder
		for offset := 0; offset < len(value); offset++ {
			if value[offset] != '\\' || offset+1 == len(value) {
				return "", fmt.Errorf("invalid SpacesAfter escape %q", value)
			}
			offset++
			switch value[offset] {
			case 's':
				decoded.WriteByte(' ')
			case 't':
				decoded.WriteByte('\t')
			case 'n':
				decoded.WriteByte('\n')
			case 'r':
				decoded.WriteByte('\r')
			default:
				return "", fmt.Errorf("unsupported SpacesAfter escape \\%c", value[offset])
			}
		}
		space = decoded.String()
	}
	return space, nil
}

func align(sentence *Sentence) error {
	if !utf8.ValidString(sentence.Text) {
		return fmt.Errorf("invalid UTF-8 text")
	}
	position := 0
	for index := range sentence.Tokens {
		token := &sentence.Tokens[index]
		position = skipSpace(sentence.Text, position)
		if !strings.HasPrefix(sentence.Text[position:], token.Form) {
			return fmt.Errorf("FORM %q (ID %s) does not match # text at byte %d", token.Form, token.ID, position)
		}
		token.Start = position
		position += len(token.Form)
		token.End = position
	}
	if skipSpace(sentence.Text, position) != len(sentence.Text) {
		return fmt.Errorf("unconsumed # text after byte %d", position)
	}
	return nil
}

func skipSpace(text string, position int) int {
	for position < len(text) {
		r, size := utf8.DecodeRuneInString(text[position:])
		if !unicode.IsSpace(r) {
			break
		}
		position += size
	}
	return position
}

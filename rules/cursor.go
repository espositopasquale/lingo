package rules

import (
	"unicode"
	"unicode/utf8"
)

type Cursor struct {
	text string
	pos  int
}

func NewCursor(text string, start int) Cursor {
	return Cursor{
		text: text,
		pos:  start,
	}
}

func (c Cursor) Pos() int {
	return c.pos
}

func (c Cursor) End() bool {
	return c.pos >= len(c.text)
}

func (c Cursor) Rune() (rune, int) {
	if c.End() {
		return utf8.RuneError, 0
	}

	return utf8.DecodeRuneInString(c.text[c.pos:])
}

func (c Cursor) Advance() Cursor {
	if c.End() {
		return c
	}

	_, size := c.Rune()

	return Cursor{
		text: c.text,
		pos:  c.pos + size,
	}
}

func (c Cursor) IsLetter() bool {
	if c.End() {
		return false
	}

	r, _ := c.Rune()
	return unicode.IsLetter(r)
}

// Recognizes combining marks, such as an accent stored after its letter.
func (c Cursor) IsMark() bool {
	if c.End() {
		return false
	}
	r, _ := c.Rune()
	return unicode.IsMark(r)
}

// Keeps dictionary matches from cutting a word short.
// Apostrophes in words like l'amico and combining accents also count.
func wordContinues(text string, end int) bool {
	c := NewCursor(text, end)
	if c.IsLetter() || c.IsDigit() || c.IsMark() {
		return true
	}
	if isApostrophe(c) {
		return c.Advance().IsLetter()
	}
	return false
}

func (c Cursor) IsDigit() bool {
	if c.End() {
		return false
	}

	r, _ := c.Rune()
	return unicode.IsDigit(r)
}

func (c Cursor) IsPunct() bool {
	if c.End() {
		return false
	}

	r, _ := c.Rune()
	return unicode.IsPunct(r)
}

func (c Cursor) MatchRune(expected rune) bool {
	if c.End() {
		return false
	}

	r, _ := c.Rune()
	return r == expected
}

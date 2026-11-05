package rules

import (
	"strings"
	"unicode"
	"unicode/utf8"
)

type UnitRule struct {
	units map[string]struct{}
}

func NewUnitRule(units []string) *UnitRule {
	m := make(map[string]struct{}, len(units))

	for _, unit := range units {
		unit = strings.TrimSpace(unit)

		if unit == "" {
			continue
		}

		m[strings.ToLower(unit)] = struct{}{}
	}

	return &UnitRule{
		units: m,
	}
}

func (r *UnitRule) Name() string {
	return "unit"
}

func (r *UnitRule) Priority() int {
	return 300
}

func (r *UnitRule) Match(
	text string,
	start int,
) (Match, bool) {
	if start < 0 || start >= len(text) {
		return Match{}, false
	}

	// Treat a unit as such only after a number: "12 cm", not "a casa".
	// Match the full symbol so units like m² and °C work too.
	previous := strings.TrimRightFunc(text[:start], unicode.IsSpace)
	last, _ := utf8.DecodeLastRuneInString(previous)
	if !unicode.IsDigit(last) {
		return Match{}, false
	}
	end := -1
	for unit := range r.units {
		candidateEnd := start + len(unit)
		if candidateEnd > len(text) || !strings.EqualFold(text[start:candidateEnd], unit) {
			continue
		}
		if !wordContinues(text, candidateEnd) && candidateEnd > end {
			end = candidateEnd
		}
	}
	if end < 0 {
		return Match{}, false
	}

	return Match{
		Start: start,
		End:   end,
		Type:  Unit,
	}, true
}

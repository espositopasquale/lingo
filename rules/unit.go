package rules

import "strings"

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

	c := NewCursor(text, start)

	if !c.IsLetter() {
		return Match{}, false
	}

	for c.IsLetter() {
		c = c.Advance()
	}

	end := c.Pos()
	value := text[start:end]

	if _, ok := r.units[strings.ToLower(value)]; !ok {
		return Match{}, false
	}

	if end < len(text) {
		next := NewCursor(text, end)

		if next.IsLetter() {
			return Match{}, false
		}
	}

	return Match{
		Start: start,
		End:   end,
		Type:  Unit,
	}, true
}

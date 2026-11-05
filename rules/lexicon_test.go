package rules

import "testing"

func TestIndexedAbbreviations(t *testing.T) {
	r := NewAbbreviationRuleWithNumericOnly([]string{"sig.", "sig.ra.", "k.", "Σ.", "art."}, []string{"No.", "art."})
	for _, tc := range []struct {
		text       string
		start, end int
		ok         bool
	}{
		{"è SIG.RA. Rossi", 3, 10, true},
		{"sig.ramo", 0, 0, false},
		{"sig.ra.mente", 0, 0, false},
		{"sig.ra.\u0301", 0, 0, false},
		{"sig.ra.'amica", 0, 0, false},
		{"K.", 0, 4, true},
		{"ς.", 0, 3, true},
		{"No.\t٢", 0, 3, true},
		{"No. grazie", 0, 0, false},
		{"No.", 0, 0, false},
		{"art. citato", 0, 4, true},
		{"test", -1, 0, false},
		{"test", 4, 0, false},
	} {
		match, ok := r.Match(tc.text, tc.start)
		if ok != tc.ok || ok && (match.Start != tc.start || match.End != tc.end || match.Type != Abbreviation) {
			t.Errorf("%q at %d: got %+v, %v", tc.text, tc.start, match, ok)
		}
	}
}

func TestAcronymBoundariesAndCase(t *testing.T) {
	r := NewAcronymRule([]string{"CNR", "AD", "CdA", "M5S"})
	for _, tc := range []struct {
		text string
		end  int
	}{
		{"C.N.R.", 6}, {"C.N.R", 5}, {"R.E.M.", 6}, {"É.È.", 6},
		{"CNR", 3}, {"CdA", 3}, {"CDA", 3}, {"M5S", 3}, {"AD", 2},
		{"Ad", -1}, {"ad", -1}, {"Cnr", -1}, {"Cda", -1}, {"SCONOSCIUTA", -1},
		{"CNR123", -1}, {"CNR’amico", -1}, {"C.N.R.foo", -1}, {"C.N.Rx", -1},
		{"A.", -1}, {"a.b.", -1}, {"", -1},
	} {
		match, ok := r.Match(tc.text, 0)
		if ok != (tc.end >= 0) || ok && (match.End != tc.end || match.Type != Acronym) {
			t.Errorf("%q: got %+v, %v; want end %d", tc.text, match, ok, tc.end)
		}
	}
}

package tokenizer

import "github.com/espositopasquale/lin-go/rules"

func italianUnits() *rules.UnitRule {
	return rules.NewUnitRule([]string{
		"mm",
		"cm",
		"dm",
		"m",
		"dam",
		"hm",
		"km",

		"mm²",
		"cm²",
		"m²",
		"km²",

		"ml",
		"cl",
		"dl",
		"l",
		"hl",

		"mg",
		"cg",
		"dg",
		"g",
		"dag",
		"hg",
		"kg",
		"t",

		"ms",
		"s",
		"min",
		"h",

		"°c",
		"°f",
		"k",

		"v",
		"mv",
		"kv",
		"a",
		"ma",

		"w",
		"kw",
		"mw",

		"j",
		"kj",
		"mj",
		"kwh",

		"hz",
		"khz",
		"mhz",
		"ghz",

		"b",
		"kb",
		"mb",
		"gb",
		"tb",
	})
}

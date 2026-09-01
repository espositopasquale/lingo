package tokenizer

import "github.com/espositopasquale/lin-go/rules"

func italianAcronyms() *rules.AcronymRule {
	return rules.NewAcronymRule([]string{
		"UE",
		"ONU",
		"NATO",
		"USA",
		"UK",

		"RAI",
		"FIAT",

		"INPS",
		"INAIL",
		"ISTAT",
		"ASL",
		"SSN",
		"SPID",
		"PEC",

		"IVA",
		"IRPEF",
		"IMU",
		"ISEE",
		"PIL",
		"PMI",

		"BCE",
		"ABI",
		"ACI",
		"AGCOM",
		"CONSOB",

		"ONG",
		"ONLUS",
		"ODV",

		"PD",
		"FI",
		"FDI",
		"M5S",

		"TG",
		"TV",
		"PC",
		"FAQ",
		"API",
		"URL",

		"CdA",
		"AD",
		"CEO",
	})
}

package tokenizer

import "github.com/espositopasquale/lingo/rules"

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

		// Institutions, international organizations and common technical acronyms.
		"CNR", "ENEA", "INFN", "ISS", "ANSA", "ANAC", "ARPA",
		"MEF", "MIUR", "MUR", "MISE", "MIMIT", "MAECI",
		"INPDAP", "ENI", "ENEL", "CUP", "TAR", "CNEL",
		"CONI", "FIGC", "CGIL", "CISL", "UIL",
		"OMS", "FAO", "UNESCO", "UNICEF", "OCSE", "FMI", "WTO",
		"PDF", "HTML", "XML", "JSON", "HTTP", "HTTPS", "SQL", "GPS",
		"IBAN", "BIC", "SWIFT", "PIN", "OTP", "GDPR", "DPO",
		"S.p.A.", "S.r.l.", "S.n.c.", "S.a.s.",
	})
}

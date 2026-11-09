package tokenizer

import "github.com/espositopasquale/lingo/rules"

func italianAbbreviations() *rules.AbbreviationRule {
	words := []string{

		// Mesi

		"gen.",
		"genn.",
		"feb.",
		"febbr.",
		"mar.",
		"apr.",
		"mag.",
		"magg.",
		"giu.",
		"lug.",
		"ago.",
		"ag.",
		"set.",
		"sett.",
		"ott.",
		"nov.",
		"dic.",
		"distr.",
		"l.",

		// Giorni della settimana

		"lun.",
		"mar.",
		"mer.",
		"gio.",
		"ven.",
		"sab.",
		"dom.",

		// Titoli e professioni

		"sig.",
		"sig.ra",
		"sig.ra.",
		"sig.na",
		"sig.na.",
		"sigg.",
		"sigg.re",

		"dr.",
		"dott.",
		"dott.ssa",
		"dott.sse",

		"avv.",
		"ing.",
		"arch.",
		"geom.",
		"rag.",
		"prof.",
		"prof.ssa",
		"prof.sse",

		"pres.",
		"on.",
		"sen.",
		"sen.ri",
		"sen.re",

		// Formule di cortesia

		"egr.",
		"egr.ia",
		"egr.i",

		"gent.mo",
		"gent.ma",
		"gent.mi",
		"gent.me",

		"ill.mo",
		"ill.ma",
		"ill.mi",
		"ill.me",

		"spett.",
		"spett.le",

		"ch.mo",
		"ch.ma",
		"ch.mi",
		"ch.me",

		"rev.",
		"rev.mo",

		// Abbreviazioni comuni

		"ecc.",
		"etc.",
		"es.",
		"e.g.",
		"cfr.",
		"v.",
		"ved.",
		"vedi",
		"n.",
		"nn.",
		"num.",
		"n.ro",
		"pag.",
		"pagg.",
		"cap.",
		"capp.",
		"par.",
		"art.",
		"artt.",
		"sez.",
		"vol.",
		"voll.",
		"fig.",
		"figg.",
		"tab.",
		"tav.",
		"pp.",

		// Azienda / società

		"s.p.a.",
		"s.p.a",
		"s.r.l.",
		"s.r.l",
		"s.n.c.",
		"s.a.s.",
		"s.a.",
		"soc.",
		"az.",
		"az.le",
		"uff.",
		"amm.",
		"ammin.",
		"dir.",
		"resp.",
		"tel.",
		"tel",
		"fax",
		"cell.",
	}
	words = append(words, italianImportedAbbreviations...)
	return rules.NewAbbreviationRuleWithNumericOnly(words, italianNumericAbbreviations)
}

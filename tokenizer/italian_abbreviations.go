package tokenizer

import "github.com/espositopasquale/lin-go/rules"

func italianAbbreviations() *rules.AbbreviationRule {
	return rules.NewAbbreviationRule([]string{

		// --------------------------------------------------
		// Mesi
		// --------------------------------------------------

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

		// --------------------------------------------------
		// Giorni della settimana
		// --------------------------------------------------

		"lun.",
		"mar.",
		"mer.",
		"gio.",
		"ven.",
		"sab.",
		"dom.",

		// --------------------------------------------------
		// Titoli e professioni
		// --------------------------------------------------

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

		// --------------------------------------------------
		// Formule di cortesia
		// --------------------------------------------------

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

		// --------------------------------------------------
		// Abbreviazioni comuni
		// --------------------------------------------------

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

		// --------------------------------------------------
		// Azienda / società
		// --------------------------------------------------

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

		// --------------------------------------------------
		// Istituzioni e sigle
		// --------------------------------------------------

		"ue",
		"unione europea",
		"onu",
		"nato",
		"usa",
		"uk",
		"ue",
		"rai",
		"inps",
		"inail",
		"iva",
		"irpef",
		"imu",
		"spid",
		"pec",
		"asl",
		"ssn",
		"istat",
		"agcom",
		"consob",
		"bce",
		"abi",
		"aci",

		// --------------------------------------------------
		// Politica / organizzazioni
		// --------------------------------------------------

		"pd",
		"fi",
		"fdi",
		"m5s",
		"ong",
		"onlus",
		"odv",

		// --------------------------------------------------
		// Economia / amministrazione
		// --------------------------------------------------

		"pil",
		"pil.",
		"cda",
		"ad",
		"ceo",
		"pmi",
		"iva",
		"irpef",
		"isee",
		"cciaa",

		// --------------------------------------------------
		// Unità e abbreviazioni tecniche
		// --------------------------------------------------

		"km",
		"m",
		"cm",
		"mm",
		"dm",
		"kg",
		"g",
		"mg",
		"l",
		"ml",
		"cl",
		"dl",

		// --------------------------------------------------
		// Informatica / comunicazione
		// --------------------------------------------------

		"pc",
		"tv",
		"tg",
		"msg",
		"app",
		"web",
		"faq",
		"api",
		"url",

		// --------------------------------------------------
		// Scrittura informale / chat
		// --------------------------------------------------

		"nn",
		"sn",
		"qnd",
		"qlcs",
		"cmq",
		"grz",
		"ke",
		"ki",
		"xkè",
		"xke",
		"xò",
		"x",
		"tvb",
	})
}

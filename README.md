**English** | [Italiano](README.it.md)

# LinGo

<picture>
  <source media="(prefers-color-scheme: dark)" srcset="assets/brand/wordmark-on-dark.svg">
  <img src="assets/brand/wordmark-on-light.svg" alt="LinGo — Lin in the Italian tricolour and go in Go blue" width="420">
</picture>

**Rule-based Italian tokenization in pure Go, with token types, UTF-8 offsets
and an evaluation suite based on Universal Dependencies.**

LinGo segments Italian text and recognizes words, numbers, dates, times,
abbreviations, acronyms, units of measure, emails, URLs, punctuation and symbols.
Each token keeps its original value and its position in the source text.

The project includes two profiles: one that keeps elided words whole, such as
`l'amico`, and one that follows the Italian surface-token conventions of
Universal Dependencies, which produces `l'` and `amico`.

> **Status: active development.** The package has tests, fuzzing and baselines on
> a real Italian corpus. APIs, rules and segmentation conventions may still
> change; suitability for production use must be verified on your own texts and
> requirements.

## Contents

- [Current status](#current-status)
- [Requirements and installation](#requirements-and-installation)
- [First use](#first-use)
- [API and token model](#api-and-token-model)
- [Italian profiles](#italian-profiles)
- [Recognition and conventions](#recognition-and-conventions)
- [Command-line interface](#command-line-interface)
- [Universal Dependencies and CoNLL-U](#universal-dependencies-and-conll-u)
- [Results and reproducibility](#results-and-reproducibility)
- [Architecture](#architecture)
- [Testing and development](#testing-and-development)
- [Known limitations](#known-limitations)
- [Roadmap](#roadmap)
- [Visual identity](#visual-identity)
- [License](#license)

## Current status

| Area | Available today |
| --- | --- |
| Language | Profiles dedicated to Italian. |
| Implementation | Go and the standard library, with no external Go dependencies. |
| Execution | Local, with no models to download and no external NLP services. |
| Tokenization | Deterministic rules and lexicons of abbreviations, acronyms and units. |
| Output | Original value, token type, start and end offsets in bytes. |
| Profiles | Standard Italian and UD Italian for surface tokens. |
| Tools | `lingo` CLI, CoNLL-U reader and `lingo-eval` evaluation CLI. |
| Quality | Offline suite, linguistic regressions, fuzzing and race detector checks. |
| External evaluation | Dev and test baselines on UD Italian-ISDT r2.17. |
| Performance | Repeatable throughput and allocation benchmarks; optimization still open. |
| Other European languages | Future work; no dedicated profile is available at the moment. |

Possible uses include preparing text for NLP pipelines, indexing, extracting
textual elements, and tools that must map tokens back to the original document.
The choice of profile depends on the conventions required by the application or
by the reference corpus.

## Requirements and installation

The module is:

```text
github.com/espositopasquale/lingo
```

The minimum version declared in [`go.mod`](go.mod) is **Go 1.27.0**.
The Go toolchain is enough to build and use the tokenizer, the CoNLL-U reader and
the evaluator. `make` is optional and provides shortcuts for project commands.

From a checkout of the repository:

```sh
go build ./...
go install ./cmd/lingo ./cmd/lingo-eval
```

Binaries are installed in the destination configured by Go. The `go run` examples
in the following sections work directly from the repository root and use the code
of that checkout.

UD corpora are downloaded separately and manually for external evaluation.
Tokenization and the offline tests work without those files.

## First use

```go
package main

import (
	"fmt"

	"github.com/espositopasquale/lingo/tokenizer"
)

func main() {
	text := "L'amico è nella sala."
	tk := tokenizer.NewItalianUD()

	for _, token := range tk.Tokenize(text) {
		fmt.Printf("%q\t%s\t[%d:%d]\n",
			token.Value, token.Type, token.Start, token.End)
	}
}
```

Output:

```text
"L'"    word         [0:2]
"amico" word         [2:7]
"è"     word         [8:10]
"nella" word         [11:16]
"sala"  word         [17:21]
"."     punctuation  [21:22]
```

The tab alignment of the output depends on the terminal. To keep `L'amico` as a
single word, use `tokenizer.NewItalian()`.

## API and token model

### Available functions

| API | Use |
| --- | --- |
| `tokenizer.Tokenize(text)` | Tokenizes with the standard Italian profile and keeps punctuation. |
| `tokenizer.TokenizeWithoutPunctuation(text)` | Uses the standard profile and excludes `Punctuation` tokens. |
| `tokenizer.NewItalian(options...)` | Builds a reusable standard Italian tokenizer. |
| `tokenizer.NewItalianUD(options...)` | Builds an Italian tokenizer following UD conventions. |
| `tk.Tokenize(text)` | Returns `[]tokenizer.Token` using the profile and options of the instance. |

To process many texts with the same profile, build the tokenizer once and reuse
it. The convenience functions build a new instance on every call.

### Structure and offsets

```go
type Token struct {
	Value string
	Type  TokenType
	Start int
	End   int
}
```

`Start` is inclusive and `End` is exclusive. Both are **UTF-8 byte offsets**,
measured on the original string. For every token:

```go
text[token.Start:token.End] == token.Value
```

In the previous example `è` takes two bytes and has the range `[8:10]`. Offsets
are not rune indexes or visual columns.

The tokenizer preserves capitalization, accents and original spelling. It does
not normalize Unicode. Combining accents are kept in the word; the sequence
`e\u0300` remains distinct from the precomposed form `è`.

Spaces and line breaks do not become tokens. To faithfully reconstruct the
document, keep the original text and use the offsets: concatenating only the
`Value` fields removes the spaces between tokens.

### Token types

| Go type | `String()` | Example or meaning |
| --- | --- | --- |
| `Word` | `word` | `amico`, `perché`. |
| `Number` | `number` | `1.250,50`; in the UD profile also forms such as `2,5%`. |
| `Date` | `date` | `31/08/2026`. |
| `Time` | `time` | `14:30`. |
| `Abbreviation` | `abbreviation` | `dott.`, `ecc.`. |
| `Acronym` | `acronym` | Known acronyms (`UE`, `M5S`, `CdA`) and dotted initials (`C.N.R.`). |
| `Unit` | `unit` | `kg`, `m²`, `°C` in a recognized numeric context. |
| `Email` | `email` | `nome@mail.esempio.it`. |
| `URL` | `url` | `https://esempio.it`. |
| `Punctuation` | `punctuation` | Commas, periods, quotation marks and other marks. |
| `Symbol` | `symbol` | Symbols and characters not classified by other rules. |
| `Unknown` | `unknown` | Initial value of the enumeration. |

These categories describe the lexical recognition performed by LinGo. UD
grammatical categories, such as `NOUN`, `VERB` and `DET`, belong to a different
annotation.

### Options

The current public option is `KeepPunctuation`, enabled by
`tokenizer.DefaultOptions()`:

```go
options := tokenizer.DefaultOptions()
options.KeepPunctuation = false

tk := tokenizer.NewItalianUD(options)
tokens := tk.Tokenize("L'amico?!")
// Values: ["L'", "amico"]
```

The options you provide replace the defaults. In particular, the zero value
`tokenizer.Options{}` has `KeepPunctuation: false`.

The filter removes only tokens classified as `Punctuation`: the apostrophe inside
a word or an elided prefix stays in its token. Symbols and percentages embedded
in a UD numeric token follow their own category.

## Italian profiles

| Text | Standard | UD |
| --- | --- | --- |
| `l'amico` | `l'amico` | `l'`, `amico` |
| `com'è` | `com'è` | `com'`, `è` |
| `po'` | `po'` | `po'` |
| `nella` | `nella` | `nella` |
| `dammelo` | `dammelo` | `dammelo` |
| `e-mail` | `e`, `-`, `mail` | `e-mail` |
| `2,5%` | `2,5`, `%` | `2,5%` |
| `!!!` | `!!!` | `!`, `!`, `!` |

The UD profile uses a list of elided prefixes and additional conventions for
compounds, English possessives found in Italian text, initials and some numeric
forms. Elision splitting is rule-based: it does not involve a grammatical
analysis of the context.

Choose the standard profile if you want to keep elided words whole in your
application. Choose the UD profile if you need to compare token boundaries with
annotations that follow those conventions. The differences between the profiles
are part of the intended behavior.

## Recognition and conventions

The engine tries the rules that apply at the current position. The rule with the
highest priority wins and, with equal priority, the longest match. If there is
still a tie, the first match in registration order is kept.

The specialized rules for URLs, emails, dates, times and abbreviations take
priority over the general word rule. This makes it possible, for example, to keep
`dott.` and an email address as single tokens.

Some useful conventions to know:

- Standard numbers follow Italian notation with a decimal comma and thousands
  groups separated by periods. The UD profile includes further forms such as
  ranges, percentages, `60x60` dimensions and abbreviated years like `'90`.
- Recognized dates have a one- or two-digit day and month, a four-digit year and
  a consistent separator: `/`, `-` or `.`. Times follow the hour:minutes form.
  Form recognition does not validate the calendar or the range of the values and
  does not convert the result to `time.Time`.
- The lexicons integrate the Italian prefixes from Moses and the Italian
  exceptions from spaCy, adapted to the conventions of the two profiles, without
  external dependencies. Abbreviation recognition uses an index by Unicode prefix
  and chooses the longest valid form: `sig.ra.` takes precedence over `sig.`.
- An all-uppercase word needs a lexicon entry to become an acronym. Explicit
  mixed-case spellings are allowed (`CdA`, `S.p.A.`); `ad`, `Ad`, `usa` and `Usa`
  remain words. Two or more dotted uppercase initials are recognized even without
  a dedicated entry, with an optional final period: `C.N.R.`, `C.N.R`, `R.E.M.`.
  This heuristic recognizes the form, not the meaning.
- `No.` is an abbreviation only before a digit, even after spaces; `No. grazie`
  keeps the period separate. The already unconditional entries `art.` and `pp.`
  stay that way. Single initials (`M. Rossi`) require a capital letter.
- The exceptions `e-mail`, `by-pass`, `c/c`, `e/o` and `n°` are `Word` tokens in
  both profiles. `C.so` is an abbreviation. The exceptions do not change the
  punctuation policy (`sett..` → `sett.`, `.`) or the elision policy (`all'art.`
  stays whole in the standard profile and is split in UD).
- Units are recognized in a numeric context, even with spaces between number and
  unit. Support includes forms such as `3 kg`, `3 m²` and `20 °C`.
- The email rules also handle subdomains. The URL rules keep normal trailing
  punctuation separate. These are tokenization rules, not complete validators of
  addresses or URLs.
- Unrecognized characters are still covered by the `Symbol` fallback. This does
  not imply complete grapheme or emoji-sequence segmentation.

## Command-line interface

### Tokenizing text

```sh
go run ./cmd/lingo -text "L'amico arriva alle 14:30."
go run ./cmd/lingo -ud -text "L'amico è nella sala."
go run ./cmd/lingo -ud -no-punct -text "L'amico è qui?!"
```

The CLI prints value, type and byte range. It also accepts positional arguments
and text from standard input:

```sh
go run ./cmd/lingo -ud "Un testo italiano."
go run ./cmd/lingo -ud < documento.txt
```

Without text, arguments or redirected input, it opens an interactive session.
`exit` and `quit` end the session. Positional arguments are joined with a single
space; to preserve the spacing of a document, use standard input or `-text` with
a single string.

| Flag | Default | Purpose |
| --- | --- | --- |
| `-text` | Empty string | Text to tokenize; if not empty it takes precedence over the other sources. |
| `-ud` | `false` | Selects the Italian UD profile. |
| `-no-punct` | `false` | Excludes punctuation tokens. |

### Evaluating a corpus

```sh
go run ./cmd/lingo-eval -corpus evaluation/testdata/italian.conllu -min-f1 1
go run ./cmd/lingo-eval -corpus data/ud/italian-isdt/r2.17/it_isdt-ud-dev.conllu -profile ud -max-errors 10
```

| Flag | Default | Purpose |
| --- | --- | --- |
| `-corpus` | Required | Path to a local CoNLL-U file. |
| `-profile` | `ud` | Profile `ud` or `default`. |
| `-format` | `text` | Human-readable report or `json`. |
| `-max-errors` | `5` | Maximum number of diagnostic examples; all sentences are evaluated. |
| `-min-f1` | `0` | F1 threshold in `[0,1]`; lower values make the command fail. |
| `-expect-sha256` | Disabled | Requires a specific corpus checksum. |

The exit code is `0` on success, `1` for read, evaluation, checksum or F1
threshold errors, and `2` for invalid arguments. A JSON report may be emitted even
when the F1 threshold is not reached: in automation, check the exit code of the
command.

## Universal Dependencies and CoNLL-U

[Universal Dependencies](https://universaldependencies.org/) is used as the
annotation reference and as the source of the evaluation corpora. LinGo applies
local rules: it does not train a model on UD files and does not call any external
service during tokenization or evaluation.

The integration concerns **surface tokens**, that is, the units that can be
aligned to the original text. In CoNLL-U, a line with ID `2-3` can represent
`nella`, while lines `2` and `3` describe the syntactic words `in` and `la`.

The [`conllu`](conllu/reader.go) reader:

- Uses the original form of the multiword line as a single surface token.
- Avoids counting the components of that multiword token again.
- Excludes empty nodes, with IDs such as `3.1`, from the token count.
- Aligns forms to `# text`, preserving byte offsets.
- If `# text` is missing, reconstructs the text using `SpaceAfter=No` and
  `SpacesAfter` for the supported `\s`, `\t`, `\n` and `\r` escapes.
- Reports format, structure or alignment errors instead of silently discarding
  problematic sentences.

The same distinction applies to `dammelo`: it remains a surface token, with no
automatic expansion into syntactic words. The package does not provide POS
tagging, lemmatization, dependency parsing or a complete validator of UD graphs.

Detailed documentation, with format sources and corpus attributions, is in the
[evaluation guide](docs/evaluation.md).

## Results and reproducibility

The current results, after the integration of the lexicons, refer to
**UD Italian-ISDT r2.17**:

| Set | Sentences | Gold surface tokens | F1 standard | F1 UD | Exact sentences UD |
| --- | ---: | ---: | ---: | ---: | ---: |
| Dev | 564 | 11,132 | 96.34% | **99.82%** | 555 / 564, i.e. 98.40% |
| Test | 482 | 9,680 | 95.89% | **99.71%** | 468 / 482, i.e. 97.10% |

The full reports with precision, recall, counts and SHA-256 are in
[`evaluation/baselines`](evaluation/baselines/README.md), with the
`-lexicons.json` suffix. The initial baselines are kept for comparison.

A token is correct when both ends of its range coincide with those of the gold
token. The metric includes punctuation and is aggregated over the whole set:

```text
precision = correct tokens / predicted tokens
recall    = correct tokens / gold tokens
F1        = 2 × correct tokens / (predicted tokens + gold tokens)
```

The exact-sentence percentage requires the whole token sequence of the sentence to
match the gold. Sentence boundaries are provided by the corpus: these results do
not measure sentence splitting, classification of LinGo token types, POS tagging,
syntactic expansions, UAS or LAS.

The rules were refined on dev; test was evaluated after those refinements.
Subsequent runs on the same files verify reproducibility. The lower F1 of the
standard profile also reflects its deliberate choice to keep elisions: the
comparison measures adherence to UD conventions, not the absolute usefulness of
the two profiles.

The coverage of a single treebank does not demonstrate the same accuracy on
social media text, OCR, dialects, specialized domains or other languages. The
results are a baseline for this corpus and this implementation.

### Repeating the evaluation

Manually download the dev and test files of release r2.17, together with
`LICENSE.txt` and the corpus README, as described in the
[guide](docs/evaluation.md#external-corpus-ud-italian-isdt-r217). Place them in:

```text
data/ud/italian-isdt/r2.17/
```

Then run:

```sh
make evaluate-dev
make evaluate-test
```

The targets verify the checksums pinned in the [Makefile](Makefile), produce JSON
and require an F1 of at least `0.99`. A file with different content makes the
checksum check fail. The train set is not needed for this rule-based evaluation.

For a different configuration:

```sh
make evaluate-dev UD_DIR=/path/to/corpus MIN_F1=0.995
```

To keep a report or work with another local CoNLL-U corpus:

```sh
go run ./cmd/lingo-eval -corpus path/corpus.conllu -format json -max-errors 0
```

The JSON includes the hash of the corpus bytes and excludes timestamps and
execution times. With the same data, code, profile and diagnostic limit, the
result is repeatable. `-expect-sha256` lets you enforce the same input in
subsequent runs. External corpora and local reports under `evaluation/results/`
are excluded from version control through [`.gitignore`](.gitignore).

## Architecture

| Path | Responsibility |
| --- | --- |
| [`tokenizer/`](tokenizer/) | API, Italian profiles, options, types and lexicons. |
| [`rules/`](rules/) | UTF-8 cursor, recognition rules and match selection. |
| [`conllu/`](conllu/) | Reading and alignment of CoNLL-U surface tokens. |
| [`evaluation/`](evaluation/) | Metrics, diagnostics, original fixtures and baselines. |
| [`cmd/lingo/`](cmd/lingo/) | CLI tokenization and interactive session. |
| [`cmd/lingo-eval/`](cmd/lingo-eval/) | Evaluation of a local corpus and quality thresholds. |
| [`docs/evaluation.md`](docs/evaluation.md) | Evaluation protocol, data and references. |
| [`assets/brand/`](assets/brand/) | Logos, icons, exports and visual guide. |

The tokenizer advances over the original string. It skips spaces, queries the
engine, emits the selected match and restarts from its end. The engine rejects
empty matches, out-of-range matches, or matches whose start differs from the
requested position.

The contract of a rule is:

```go
type Rule interface {
	Name() string
	Priority() int
	Match(text string, start int) (Match, bool)
}
```

`rules.NewEngine` lets you assemble rules at the engine level. The public
`Tokenizer` API currently exposes the two Italian constructors; it does not yet
offer a generic constructor to plug in a custom engine. To modify a profile of the
package you need to change its construction and add regressions for the desired
behavior.

## Testing and development

```sh
go test ./...
make check
make evaluate-fixtures
make benchmark
go test ./tokenizer -run '^$' -fuzz FuzzTokenizerDoesNotPanic -fuzztime 10s -parallel 2
```

| Command | Check |
| --- | --- |
| `go test ./...` or `make test` | Unit tests, regressions and fuzzing seeds. |
| `make check` | `go vet` and tests with the race detector; the latter requires a compatible Go/CGO environment. |
| `make evaluate-fixtures` | F1 equal to 1 on the 22 original examples of the repository. |
| `make evaluate` | Offline fixtures and evaluation of the external dev set. |
| `make evaluate-dev` | Dev corpus with checksum and F1 threshold. |
| `make evaluate-test` | Test corpus with checksum and F1 threshold. |
| `make benchmark` | Throughput and allocations of the two profiles, reusing the tokenizer; three runs. |

The fixture [`evaluation/testdata/italian.conllu`](evaluation/testdata/italian.conllu)
contains **22 original sentences and 164 gold tokens**. The POS and dependency
columns are structural placeholders: the fixture verifies segmentation, not
grammatical annotation. Its perfect score is a regression check on curated
examples and must be distinguished from the results on external corpora.

The tests also cover UTF-8 offsets, combining accents, classification of words
and acronyms, units, URLs, emails, multiword tokens, empty nodes, invalid CoNLL-U
input, metrics computed on known examples, checksums and report determinism.
Fuzzing exercises both profiles and verifies progress, text coverage and
consistency between values and ranges.

Speed measurements depend on the machine, the Go version and the texts used. The
repository does not publish a demonstrated performance advantage over other
libraries. The benchmarks are meant to observe and compare revisions under the
same conditions; engine optimization and allocations remain areas of work.

If the default Go cache is read-only, point to a writable directory:

```sh
GOCACHE=/tmp/lingo-build-cache make check
```

## Known limitations

- Dedicated language support covers Italian. There is no automatic language
  detection or catalog of European profiles.
- The rules and lexicons do not cover every elision, abbreviation, proper noun or
  domain convention. Some ambiguities require linguistic context.
- Percentages, ranges, hyphens and other forms may be segmented differently in
  individual treebanks. The UD profile does not imply perfect conformity with all
  Italian UD annotations.
- Unicode support preserves the text and handles letters and combining accents,
  but is not a complete implementation of UAX #29 or of grapheme clusters.
- Dates, times, numbers, emails and URLs are recognized with lexical heuristics;
  they do not replace validation and parsing of values in the application.
- The tokenization API processes a complete string and returns a slice. There is
  no incremental tokenizer over `io.Reader`; the CLI reads standard input
  entirely before tokenizing it.
- The evaluator aligns `FORM` and the original text exactly and stops the
  evaluation on misalignment. It does not correct or normalize the corpus
  automatically. The reader has a limit of 16 MiB per line.
- The package does not include sentence splitting, lemmatization, tagging,
  dependency parsing or BPE/WordPiece tokenization for language models.

## Roadmap

The following are future directions, not features that are already available:

1. Extend the Italian evaluation to different genres and domains, recording
   results and error cases per corpus.
2. Strengthen the conventions for apostrophes, abbreviations, compounds and
   numbers, keeping the differences between profiles explicit.
3. Optimize lexicon recognition, scanning and allocations based on the
   benchmarks.
4. Separate common rules, linguistic data and regional conventions to make it
   easier to build new profiles.
5. Gradually introduce other European languages, accompanying each language with
   fixtures, a reference corpus and reproducible results.
6. Consolidate the public API and evaluate interfaces for customization and
   incremental tokenization.

To contribute, include a minimal example, the profile concerned, the expected
tokens and the rationale for the convention adopted. A fix should have a
regression test and be compared on dev; the external test set should be reserved
for final checks, avoiding its use to choose rules.

## Visual identity

The logo keeps `Lin` in the diagonal Italian tricolour and `go` in Go blue.
Variants are available for light and dark backgrounds, icons, monogram, favicon,
monochrome printing and social sharing.

- [Asset guide](assets/brand/README.md)
- [Variants preview](assets/brand/preview.png)
- [Complete downloadable kit](assets/lingo-brand-kit.zip)

Image regeneration uses tools separate from the Go package; details and
requirements are in the asset guide.

## License

The project code is distributed under the **GNU GPL version 3** contained in
[`LICENSE`](LICENSE). External corpora keep their own licenses and attributions,
which must be kept together with the downloaded data. Notes on the adapted
lexicons are in [`THIRD_PARTY_NOTICES.md`](THIRD_PARTY_NOTICES.md).

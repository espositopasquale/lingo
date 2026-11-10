# Italian evaluation with Universal Dependencies

LinGo evaluates **surface tokenization** against CoNLL-U gold annotations.
The tokenizer, reader and evaluator are implemented in Go with no external Go
dependencies. Evaluation never downloads data or calls an external service.

## Profiles

| API / CLI | Behavior |
| --- | --- |
| `tokenizer.NewItalian()` / `lingo` | Preserves elisions such as `l'amico` as one word. |
| `tokenizer.NewItalianUD()` / `lingo -ud` | Separates known Italian elisions into `l'` and `amico`, keeping the apostrophe with the elided prefix. |

Both profiles preserve original spelling and return UTF-8 **byte offsets**.
Combining accents are kept with their word, without normalizing the input or
changing the offsets. The UD profile also handles borrowed English possessives,
hyphenated compounds, numeric ranges and percentages, abbreviated years, and
initials. These are deterministic conventions; individual treebanks can make
different decisions, especially for percentages, ranges, and compounds.

`NewItalianUD` accepts the same `Options` as `NewItalian`. Keep punctuation enabled
when evaluating: punctuation is part of the gold token stream. `-no-punct` remains
available for application output, including when using `lingo -ud`.

### Tokens versus syntactic words

For a CoNLL-U range such as `2-3 nella`, the evaluator takes the range's original
FORM (`nella`) as **one gold surface token**. Its component rows (`in`, `la`) are
not counted a second time. The same applies when component forms do not concatenate
to the surface spelling. Empty nodes such as `3.1` have no surface realization
and are ignored for token scoring.

LinGo currently does not expand `nella` into `in` + `la`, split enclitic forms
such as `dammelo` into syntactic words, predict UPOS tags, or parse dependency
trees. The CoNLL-U reader only consumes the fields needed for surface-token
evaluation; it is not a complete UD graph validator.

References:

- [UD CoNLL-U format](https://universaldependencies.org/format.html)
- [Italian tokenization conventions](https://universaldependencies.org/it/index.html)
- [Tokens and syntactic words](https://universaldependencies.org/u/overview/tokenization.html)

## Offline regression suite

```sh
go test ./...
make evaluate-fixtures
```

`evaluation/testdata/italian.conllu` contains **22 original LinGo examples**, not
an excerpt from an external UD corpus. Their POS/dependency columns are structural
placeholders; only the surface-token boundaries are gold for this suite. They
exercise elisions, apostrophe variants, MWTs, punctuation, dates, numbers, units,
URLs, emails, and UTF-8 offsets. Additional Go regressions cover Unicode combining
marks and generalized cases found during development.

The fixture target requires token F1 = 1. This is a regression check on curated
examples, not an estimate of accuracy on arbitrary Italian text. Unit tests also
check hand-calculated precision/recall, invalid predictions, malformed CoNLL-U,
corpus checksums, and deterministic JSON output.

## External corpus: UD Italian-ISDT r2.17

Use the fixed release **r2.17** from
[UniversalDependencies/UD_Italian-ISDT](https://github.com/UniversalDependencies/UD_Italian-ISDT/tree/r2.17).
Download these files manually, retaining the original names:

- [it_isdt-ud-dev.conllu](https://raw.githubusercontent.com/UniversalDependencies/UD_Italian-ISDT/r2.17/it_isdt-ud-dev.conllu)
- [it_isdt-ud-test.conllu](https://raw.githubusercontent.com/UniversalDependencies/UD_Italian-ISDT/r2.17/it_isdt-ud-test.conllu)
- [LICENSE.txt](https://raw.githubusercontent.com/UniversalDependencies/UD_Italian-ISDT/r2.17/LICENSE.txt)
- [README.md](https://raw.githubusercontent.com/UniversalDependencies/UD_Italian-ISDT/r2.17/README.md)

Place them in `data/ud/italian-isdt/r2.17/`. This directory is git-ignored.
The corpus retains its upstream license and attribution independently of LinGo's
code license; keep its LICENSE and README alongside the data. The train split is
not needed for this rule-based evaluation.

ISDT contributors include Cristina Bosco, Alessandro Lenci, Simonetta Montemagni
and Maria Simi. As requested by the corpus license, publications reporting
results obtained from these data should cite:

> Bosco, C., Montemagni, S., & Simi, M. (2013). Converting Italian Treebanks:
> Towards an Italian Stanford Dependency Treebank. *Proceedings of the 7th
> Linguistic Annotation Workshop & Interoperability with Discourse*, pp. 61–69.

Electronic publications should also link to the original
[ISDT project page](http://medialab.di.unipi.it/wiki/ISDT/), as requested by the
license. This historical URL may no longer be available, but is retained here to
preserve the specified attribution. UD Italian-ISDT r2.17 declares the corpus
under [CC BY-NC-SA 3.0](https://github.com/UniversalDependencies/UD_Italian-ISDT/blob/r2.17/LICENSE.txt).

### Repeatable commands

```sh
make evaluate-dev
make evaluate-test
```

Both commands verify the exact file SHA-256 recorded in the Makefile, emit JSON,
and fail if token F1 falls below `MIN_F1` (default `0.99`). A checksum mismatch
fails without emitting a score. `UD_DIR` can point to another local directory
containing the same release files. Use **dev** to investigate and improve rules;
reserve **test** for a final evaluation after implementation changes.

Use the CLI directly to inspect errors or compare profiles:

```sh
go run ./cmd/lingo-eval -corpus data/ud/italian-isdt/r2.17/it_isdt-ud-dev.conllu -profile ud -max-errors 10
go run ./cmd/lingo-eval -corpus data/ud/italian-isdt/r2.17/it_isdt-ud-dev.conllu -profile default -max-errors 0 -format json
```

For another local CoNLL-U corpus, supply `-corpus PATH`. Add
`-expect-sha256 DIGEST` to require the same bytes on subsequent runs. The CLI
returns exit status 0 on success, 1 for corpus/evaluation failures or a missed
F1 threshold, and 2 for invalid arguments.

The JSON reports in [`evaluation/baselines`](../evaluation/baselines/) record the
initial repeatable ISDT scores and the subsequent `-lexicons.json` reports for
both profiles and splits. They contain only
aggregate metrics and checksums, with no corpus text. Changing rules can change
these scores; retain a baseline for comparison and document intentional changes.

### What is measured

- A predicted token is correct when its start and end **byte offsets both match**
  a gold surface token in the same original sentence.
- Token precision = correct / predicted; recall = correct / gold;
  F1 = 2 × correct / (predicted + gold). Scores are micro-averaged and reported
  in `[0, 1]`. A zero denominator produces zero.
- Exact-sentence accuracy is the fraction of gold sentences whose complete
  surface-token sequence matches. Gold sentence boundaries are supplied by the
  corpus; this does **not** measure sentence-boundary detection.
- Token classification (`Word`, `Unit`, etc.) is covered by regression tests,
  not compared with UD UPOS tags, which describe a different annotation task.
- `-max-errors` limits diagnostic examples only. Every sentence is evaluated.

When `# text` exists, FORM values must align exactly to it, allowing whitespace
between tokens. Otherwise the reader reconstructs text from surface FORMs with
`SpaceAfter=No` and `SpacesAfter` (`\s`, `\t`, `\n`, `\r`). It never silently
skips unaligned sentences. Unknown spacing escapes and malformed token structure
produce explicit errors. Different spelling or normalization in FORM and text
must be corrected in a separately documented corpus copy before evaluation.

Predictions are also checked for nonempty, nonoverlapping, valid UTF-8 spans,
unchanged values, and coverage of all non-whitespace input. Reports contain the
raw corpus SHA-256 and exclude timestamps and timings, making JSON reproducible
for the same data, tokenizer code, profile and diagnostic limit.

This is a surface-token evaluation tool, not the complete official CoNLL shared
task scorer. It does not calculate word-alignment, tagging, UAS, or LAS scores.

## Performance and robustness

```sh
make check
make benchmark
go test ./tokenizer -run '^$' -fuzz FuzzTokenizerDoesNotPanic -fuzztime 10s
```

`make check` runs `go vet` and the race-enabled test suite. The benchmark reuses
each tokenizer on the same fixed sample and reports throughput and allocations;
timings depend on the machine and are deliberately kept outside quality reports.

In an environment with a read-only default Go build cache, select a writable
cache, for example `GOCACHE=/tmp/LinGo-build-cache make check`.

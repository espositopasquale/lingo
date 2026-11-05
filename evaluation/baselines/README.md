# UD Italian-ISDT r2.17 baseline

## Initial baseline

These reports contain aggregate surface-token scores and the SHA-256 of each
complete input file. They do not redistribute corpus sentences. The dataset is
the manually supplied [UD Italian-ISDT r2.17 release](https://github.com/UniversalDependencies/UD_Italian-ISDT/tree/r2.17).

| Split | Sentences | Gold surface tokens | Default F1 | UD F1 | UD exact sentences |
| --- | ---: | ---: | ---: | ---: | ---: |
| Dev | 564 | 11,132 | 96.24% | 99.82% | 555 / 564 (98.40%) |
| Test | 482 | 9,680 | 95.83% | 99.68% | 466 / 482 (96.68%) |

The default profile intentionally preserves elisions, so it is expected to
disagree with UD gold boundaries. The comparison measures compatibility with UD
surface-token conventions, not an overall ranking of application usefulness.
Rules were refined on dev; test was evaluated after those refinements.

- [Dev, default profile](isdt-r2.17-dev-default.json)
- [Dev, UD profile](isdt-r2.17-dev-ud.json)
- [Test, default profile](isdt-r2.17-test-default.json)
- [Test, UD profile](isdt-r2.17-test-ud.json)

## Integrated lexicons (2026-11-05)

The initial reports above are retained unchanged. Reports with the `-lexicons`
suffix record the Moses/spaCy integration, indexed matching, expanded acronym
dictionary, and dotted initialisms. The measurement exception `cm.` is adapted
to preserve a separate sentence-final period after the unit `cm`. Rules were
checked on dev before evaluating test; test diagnostics were not used to tune
this revision. Classification changes are checked separately in Go tests.

| Split | Default F1 | UD F1 | UD exact sentences |
| --- | ---: | ---: | ---: |
| Dev | 96.34% | 99.82% | 555 / 564 (98.40%) |
| Test | 95.89% | 99.71% | 468 / 482 (97.10%) |

- [Dev, default profile](isdt-r2.17-dev-default-lexicons.json)
- [Dev, UD profile](isdt-r2.17-dev-ud-lexicons.json)
- [Test, default profile](isdt-r2.17-test-default-lexicons.json)
- [Test, UD profile](isdt-r2.17-test-ud-lexicons.json)

To reproduce the current UD reports, run `make -s evaluate-dev` and
`make -s evaluate-test` from the repository root. Both commands verify the pinned
corpus checksum and use `-max-errors 0 -format json`. To reproduce the default
reports, use `go run ./cmd/lingo-eval -corpus PATH -profile default -max-errors 0
-format json`, using the corresponding local file and the checksum in its report.

The scores apply to token spans within gold sentence boundaries. They exclude
sentence splitting, POS tagging, dependency parsing and syntactic MWT expansion.
See the [evaluation guide](../../docs/evaluation.md) for formulas and limitations.

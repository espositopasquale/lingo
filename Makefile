GO ?= go
UD_DIR ?= data/ud/italian-isdt/r2.17
MIN_F1 ?= 0.99
VERSION ?= dev
BIN ?= release/lingo

# Check the local ISDT r2.17 files before evaluating, so repeated runs use
# exactly the same corpus.
UD_DEV_SHA256 := 1bc7dda04c58a4d87e2e6607a0e612462b87e7da6dde56be610d3749521af744
UD_TEST_SHA256 := cea8ca0e8d1ddf7ce509f10afff426d00e058982f49473e0744cd0ee477129e3

.PHONY: build test check evaluate evaluate-fixtures evaluate-dev evaluate-test benchmark

build:
	mkdir -p $(dir $(BIN))
	$(GO) build -trimpath -ldflags "-s -w -X main.version=$(VERSION)" -o $(BIN) ./cmd/lingo

test:
	$(GO) test ./...

check:
	$(GO) vet ./...
	$(GO) test -race ./...

evaluate: evaluate-fixtures evaluate-dev

evaluate-fixtures:
	$(GO) run ./cmd/lingo-eval -corpus evaluation/testdata/italian.conllu -profile ud -max-errors 0 -min-f1 1

evaluate-dev:
	$(GO) run ./cmd/lingo-eval -corpus $(UD_DIR)/it_isdt-ud-dev.conllu -profile ud -expect-sha256 $(UD_DEV_SHA256) -max-errors 0 -format json -min-f1 $(MIN_F1)

evaluate-test:
	$(GO) run ./cmd/lingo-eval -corpus $(UD_DIR)/it_isdt-ud-test.conllu -profile ud -expect-sha256 $(UD_TEST_SHA256) -max-errors 0 -format json -min-f1 $(MIN_F1)

benchmark:
	$(GO) test ./evaluation -run '^$$' -bench BenchmarkItalianProfiles -benchmem -count 3

bootstrap:
	# Check CONTRIBUTING.md for the required dependencies
	curl -fsSL https://get.jetpack.io/devbox | bash
	curl -sfL https://direnv.net/install.sh | bash

GO_SRC := $(shell find pkg -type f -name '*.go')

structlayout: cmd/structlayout/structlayout.go $(filter-out *_test.go,$(GO_SRC))
	go build -o $@ $<

mergelayout: cmd/mergelayout/mergelayout.go $(filter-out *_test.go,$(GO_SRC))
	go build -o $@ $<

debdownload: cmd/debdownload/debdownload.go $(filter-out *_test.go,$(GO_SRC))
	go build -o $@ $<

debuginfofind: cmd/debuginfofind/debuginfofind.go $(filter-out *_test.go,$(GO_SRC))
	go build -o $@ $<

.PHONY: build
build: structlayout mergelayout debdownload debuginfofind
	go build ./...

.PHONY: generate
generate: build generate/python generate/ruby generate/glibc

.PHONY: generate/python
generate/python:
	./scripts/download/python.sh
	./scripts/structlayout/python.sh
	./scripts/mergelayout/python.sh

.PHONY: generate/ruby
generate/ruby:
	./scripts/download/ruby.sh
	./scripts/structlayout/ruby.sh
	./scripts/mergelayout/ruby.sh

.PHONY: generate/glibc
generate/glibc:
	./scripts/download/glibc.sh
	./scripts/structlayout/glibc.sh
	./scripts/mergelayout/glibc.sh

.PHONY: clean
clean:
	rm -rf target

.PHONY: check
check: vet lint

.PHONY: vet
vet:
	go vet ./...

.PHONY: lint
lint:
	go run github.com/golangci/golangci-lint/cmd/golangci-lint@latest run

.PHONY: fix
fix:
	go run github.com/golangci/golangci-lint/cmd/golangci-lint@latest run --fix

.PHONY: format
format:
	go run mvdan.cc/gofumpt@latest -l -w .

.PHONY: tagalign
tagalign:
	go run github.com/4meepo/tagalign/cmd/tagalign@latest -fix -sort ./...

.PHONY: test
test: build
	go test ./...

.PHONY: test/integration
test/integration: build
	@echo "Make sure to run 'generate' if any generation code has changed"
	go test -tags=integration ./tests/integration/...


.PHONY: test/integration/update
test/integration/update:
	go test -count=1 -race -tags=integration ./tests/integration/... -update

TMPDIR := ./tmp
$(TMPDIR):
	mkdir -p $(TMPDIR)

$(TMPDIR)/structlayout-help.txt: $(TMPDIR) ./cmd/structlayout/structlayout.go
	mkdir -p ./tmp
	go run ./cmd/structlayout/structlayout.go -h > $@ 2>&1

$(TMPDIR)/mergelayout-help.txt: $(TMPDIR) ./cmd/mergelayout/mergelayout.go
	go run ./cmd/mergelayout/mergelayout.go -h > $@ 2>&1

$(TMPDIR)/debdownload-help.txt: $(TMPDIR) ./cmd/debdownload/debdownload.go
	go run ./cmd/debdownload/debdownload.go -h > $@ 2>&1

$(TMPDIR)/debuginfofind-help.txt: $(TMPDIR) ./cmd/debuginfofind/debuginfofind.go
	go run ./cmd/debuginfofind/debuginfofind.go -h > $@ 2>&1

.PHONY: README.md
README.md: $(TMPDIR)/structlayout-help.txt $(TMPDIR)/mergelayout-help.txt $(TMPDIR)/debdownload-help.txt $(TMPDIR)/debuginfofind-help.txt
	go run github.com/campoy/embedmd/v2@latest -w README.md
	devbox generate readme CONTRIBUTING.md

# TODO(kakkoyun): DRY this file.

target/static-libs:
	@echo "Building static-libs"
	./scripts/download_and_build_static_libs.sh

target/debug/deps/libelf.a: target/static-libs
	mkdir -p target/debug/deps
	cp target/static-libs/libelf.a target/debug/deps/libelf.a

target/release/deps/libelf.a: target/static-libs
	mkdir -p target/release/deps
	cp target/static-libs/libelf.a target/release/deps/libelf.a

target/debug/deps/libz.a: target/static-libs
	mkdir -p target/debug/deps
	cp target/static-libs/libz.a target/debug/deps/libz.a

target/release/deps/libz.a: target/static-libs
	mkdir -p target/release/deps
	cp target/static-libs/libz.a target/release/deps/libz.a

deps: target/debug/deps/libelf.a target/debug/deps/libz.a target/release/deps/libelf.a target/release/deps/libz.a

.PHONY: dirs
dirs:
	mkdir -p pkg/ruby/versions
	mkdir -p pkg/python/versions

.PHONY: build
build: target/debug/deps/libelf.a target/debug/deps/libz.a dirs
	@echo "Building runtime-data"
	cargo build
	cargo run
	go build ./...

.PHONY: release-build
release-build: target/release/deps/libelf.a target/release/deps/libz.a
	@echo "Building runtime-data"
	cargo build --release

.PHONY: clean
clean:
	rm -rf target

.PHONY: test
test:
	cargo test
	go test ./...

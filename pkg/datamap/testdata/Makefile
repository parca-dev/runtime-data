CC = zig cc
TARGETS ?= aarch64-linux-musl x86_64-linux-musl

.PHONY: debug/build
debug/build: test

.PHONY: build/cross
build/cross:
	@for target in $(TARGETS); do \
		dir=$$(echo $$target | sed 's/-linux-musl//'); \
		mkdir -p $$dir; \
		$(CC) -target $$target -o $$dir/test test.c; \
	done

# Build a test program with debug information.
test: test.c
	$(CC) -g -o $@ $<

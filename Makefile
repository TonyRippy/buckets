# Keep this list in sync as new implementation language directories are added.
IMPLEMENTATIONS := go

.PHONY: all build test clean release

all: build test

build:
	@for impl in $(IMPLEMENTATIONS); do \
		$(MAKE) -C $$impl build; \
	done

test:
	@for impl in $(IMPLEMENTATIONS); do \
		$(MAKE) -C $$impl test; \
	done

clean:
	@for impl in $(IMPLEMENTATIONS); do \
		$(MAKE) -C $$impl clean; \
	done

# See docs/runbooks/release.md.
release: export RELEASE_VERSION := $(VERSION)
release: export RELEASE_MAKE := $(MAKE)
release:
	@sh scripts/go-release.sh

# Keep this list in sync as new implementation language directories are added.
IMPLEMENTATIONS := go

.PHONY: all build test clean

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

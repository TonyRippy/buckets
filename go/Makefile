GO ?= go
GOCACHE ?= $(CURDIR)/.gocache
GOMODCACHE ?= $(CURDIR)/.gomodcache

GO_ENV := GOCACHE=$(GOCACHE) GOMODCACHE=$(GOMODCACHE)
# Discover workspace main module directories from go.work and expand to package patterns.
MODULE_PACKAGES := $(addsuffix /...,$(shell $(GO) list -work -m all))

.PHONY: all build test clean

all: build test

build:
	$(GO_ENV) $(GO) build $(MODULE_PACKAGES)

test:
	$(GO_ENV) $(GO) test $(MODULE_PACKAGES)

clean:
	rm -rf .gocache .gomodcache

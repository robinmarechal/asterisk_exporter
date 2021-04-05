GO := GO15VENDOREXPERIMENT=1 go
PROMU := $(GOPATH)/bin/promu
PKGS := $(shell $(GO) list ./... | grep -v /vendor/)

PREFIX ?= $(shell pwd)
BIN_DIR ?= $(shell pwd)

pkgs          = ./...
ifeq ($(GOHOSTARCH),amd64)
        ifeq ($(GOHOSTOS),$(filter $(GOHOSTOS),linux freebsd darwin windows))
                # Only supported on amd64
                test-flags := -race
        endif
endif

all: vet format tarball

style:
	@echo ">> checking code style"
	@! gofmt -d $(shell find . -path ./vendor -prune -o -name '*.go' -print) | grep '^'

format:
	@echo ">> formatting code"
	@$(GO) fmt $(PKGS)

vet:
	@echo ">> vetting code"
	@$(GO) vet $(PKGS)

build: test promu
	@echo ">> building binaries"
	@$(PROMU) build --prefix $(PREFIX)

tarball: build
	@echo ">> building release tarball"
	@$(PROMU) tarball $(BIN_DIR) --prefix $(PREFIX)

promu:
	@which promu > /dev/null; if [ $$? -ne 0 ]; then \
		$(GO) get -u github.com/prometheus/promu; \
	fi

test:
	@echo ">> running tests"
	$(GO) test -short $(test-flags) $(pkgs)

.PHONY: all style format build vet tarball promu test

GO ?= $(shell which go)
GOFMT := $(shell which gofmt) "-s"
SWAGGER ?= ${GOPATH}/bin/swag
GO_TEST_UNIT ?= ${GOPATH}/bin/go-junit-report
GOLINT ?= ${GOPATH}/bin/golint
PACKAGES ?= $(shell $(GO) list ./... | grep -v /api)
GOFILES := $(shell find . -name "*.go" -type f)
BUILD_DIR ?= ./build

all: swag test build_linux build_windows

.PHONY: install
install: swag
	$(GO) install ./cmd/london-server
	$(GO) install ./cmd/manchester-server

.PHONY: build_linux
build_linux: deps
	env GOOS=linux GOARCH=amd64 $(GO) build -o $(BUILD_DIR)/linux64/london-server --tags "linux" -a -tags netgo -ldflags '-w -extldflags "-static"' ./cmd/london-server
	env GOOS=linux GOARCH=amd64 $(GO) build -o $(BUILD_DIR)/linux64/manchester-server --tags "linux" -a -tags netgo -ldflags '-w -extldflags "-static"' ./cmd/manchester-server

.PHONY: build_windows
build_windows: deps
	@hash gcc-multilib > /dev/null 2>&1; if [ $$? -ne 0 ]; then \
		echo "Install gcc-multilib before running this!"; \
	fi
	@hash gcc-mingw-w64 > /dev/null 2>&1; if [ $$? -ne 0 ]; then \
		echo "Install gcc-mingw-w64 before running this!"; \
	fi
	env CGO_ENABLED=1 CC=x86_64-w64-mingw32-gcc CXX=x86_64-w64-mingw32-g++ GOOS=windows GOARCH=amd64 $(GO) build -o $(BUILD_DIR)/win64/london-server.exe ./cmd/london-server
	env CGO_ENABLED=1 CC=x86_64-w64-mingw32-gcc CXX=x86_64-w64-mingw32-g++ GOOS=windows GOARCH=amd64 $(GO) build -o $(BUILD_DIR)/win64/manchester-server.exe ./cmd/manchester-server

.PHONY: swag
swag: deps
	if [ ! -f $(SWAGGER) ]; then \
		$(GO) install github.com/swaggo/swag/cmd/swag@latest; \
	fi
	$(SWAGGER) init -g ../../cmd/london-server/main.go -o api/london -d internal/london
	$(SWAGGER) init -g ../../cmd/manchester-server/main.go -o api/manchester -d internal/manchester

.PHONY: test
test:
	$(GO) install github.com/jstemmer/go-junit-report@latest
	$(GO) test -v -covermode=atomic -coverpkg=$(PACKAGES) -coverprofile coverage.out $(PACKAGES) 2>&1 | $(GO_TEST_UNIT) > $(BUILD_DIR)/go-test-report.xml
	mv coverage.out $(BUILD_DIR)

.PHONY: fmt
fmt:
	$(GOFMT) -w $(GOFILES)

.PHONY: fmt-check
fmt-check:
	@diff=$$($(GOFMT) -d $(GOFILES)); \
	if [ -n "$$diff" ]; then \
		echo "Please run 'make fmt' and commit the result:"; \
		echo "$${diff}"; \
		exit 1; \
	fi;

.PHONY: vet
vet: swag
	$(GO) vet $(PACKAGES)

.PHONY: lint
lint:
	@hash golint > /dev/null 2>&1; if [ $$? -ne 0 ]; then \
		$(GO) install golang.org/x/lint/golint@latest; \
	fi
	for PKG in $(PACKAGES); do $(GOLINT) -set_exit_status $$PKG || exit 1; done;

.PHONY: deps
deps:
	mkdir -p build
	@hash go > /dev/null 2>&1; if [ $$? -ne 0 ]; then \
		echo "Install Go language before running this!"; \
	fi
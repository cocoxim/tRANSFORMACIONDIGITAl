GOCMD:=$(shell which go)
GOLINT:=$(shell which golint)
GOIMPORT:=$(shell which goimports)
GOFMT:=$(shell which gofmt)
GOBUILD:=$(GOCMD) build
GOINSTALL:=$(GOCMD) install
GOCLEAN:=$(GOCMD) clean
GOTEST:=$(GOCMD) test
GOGET:=$(GOCMD) get
GOLIST:=$(GOCMD) list
GOVET:=$(GOCMD) vet
GOPATH:=$(shell $(GOCMD) env GOPATH)
u := $(if $(update),-u)

BINARY_NAME:=swag
PACKAGES:=$(shell $(GOLIST) github.com/swaggo/swag github.com/swaggo/swag/cmd/swag github.com/swaggo/swag/gen github.com/swaggo/swag/format)
GOFILES:=$(shell find . -name "*.go" -type f)

export GO111MODULE := on

all: test build

.PHONY: build
build: deps
	$(GOBUILD) -o $(BINARY_NAME) ./cmd/swag

.PHONY: install
install: deps
	$(GOINSTALL) ./cmd/swag

.PHONY: test
test:
	echo "mode: count" > coverage.out
	for PKG in $(PACKAGES); do \
		$(GOCMD) test -v -covermode=count -coverprofile=profile.out $$PKG > tmp.out; \
		cat tmp.out; \
		if grep -q "^--- FAIL" tmp.out; then
VERSION:=0.1
TAG:=v$(VERSION)

COVEROUT = cover.out
GOFMTCHECK = test -z `gofmt -l -s -w *.go | tee /dev/stderr`
COVER = cd plugin/metadata_edns0 && go test -v -coverprofile=$(COVEROUT) -covermode=atomic -race
GOPATH?=$(HOME)/go
GITCOMMIT:=$(shell git describe --dirty --always)
BINARY:=coredns
SYSTEM:=
VERBOSE:=-v

all: fmt test
coredns: build

.PHONY: fmt
fmt:
	@echo "Checking format..."
	@$(GOFMTCHECK)

.PHONY: test
test:
	@echo "Running tests..."
	@$(COVER)

.PHONY: build
build:
	GO111MODULE=on CGO_ENABLED=0 $(SYSTEM) go build $(VERBOSE) -ldflags="-s -w -X github.com/coredns/coredns/coremain.GitCommit=$(GITCOMMIT)" -o $(BINARY)


# Use the 'release' target to start a release
.PHONY: release
release: commit push
	@echo Released $(VERSION)

.PHONY: commit
commit:
	@echo Committing release $(VERSION)
	git commit -am"Release $(VERSION)"
	git tag $(TAG)

.PHONY: push
push:
	@echo Pushing release $(VERSION) to master
	git push --tags
	git push

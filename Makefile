
.PHONY: all
all: deps build

.PHONY: deps
deps:
	go mod download

.PHONY: test
test:
	go test -v ./...

.PHONY: lint
lint:
	gometalinter --vendor ./...

.PHONY: build
build:
	go build -ldflags=-s disgoirc.go

.PHONY: run
run:
	go run disgoirc.go


GOMETALINTER_VERSION = 2.0.11

.PHONY: setup-ci
setup-ci:
	curl -Lo ci/gometalinter.tar.gz "https://github.com/alecthomas/gometalinter/releases/download/v${GOMETALINTER_VERSION}/gometalinter-${GOMETALINTER_VERSION}-linux-amd64.tar.gz"
	tar zxf ci/gometalinter.tar.gz -C linter

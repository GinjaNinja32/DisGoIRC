
.PHONY: deps
deps:
	glide up

.PHONY: setup
setup:
	go get -u github.com/Masterminds/glide
	go get -u github.com/alecthomas/gometalinter
	gometalinter -i

.PHONY: test
test:
	go test -v $(shell glide novendor)

.PHONY: lint
lint:
	gometalinter $(shell glide novendor)

.PHONY: build
build:
	go build -ldflags=-s disgoirc.go

.PHONY: run
run:
	go run disgoirc.go

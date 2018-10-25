
.PHONY: all
all: setup deps build

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
	gometalinter $(shell glide novendor) --deadline 120s --cyclo-over 15

.PHONY: build
build:
	go build -ldflags=-s disgoirc.go

.PHONY: run
run:
	go run disgoirc.go

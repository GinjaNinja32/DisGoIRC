
.PHONY: setup
setup:
	go get -u github.com/alecthomas/gometalinter
	gometalinter -i

.PHONY: lint
lint:
	gometalinter ./...

.PHONY: run
run:
	go run disgoirc.go

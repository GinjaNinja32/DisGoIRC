
.PHONY: deps
deps:
	glide up

.PHONY: setup
setup:
	go get -u github.com/Masterminds/glide
	go get -u github.com/alecthomas/gometalinter
	gometalinter -i

.PHONY: lint
lint: deps
	gometalinter ./...

.PHONY: run
run:
	go run disgoirc.go

GOMETALINTER_VERSION = 2.0.11

.PHONY: install
install:
	curl -Lo gometalinter.tar.gz "https://github.com/alecthomas/gometalinter/releases/download/v${GOMETALINTER_VERSION}/gometalinter-${GOMETALINTER_VERSION}-linux-amd64.tar.gz"
	rm -r gometalinter || true
	mkdir gometalinter
	tar zxf gometalinter.tar.gz -C gometalinter --strip-components=1 # remove leading gometalinter-$VERSION-linux-amd64 from path

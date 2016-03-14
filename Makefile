all: install

deps:
	glide install

install:
	go install -v ./cmd/...

test:
	go test $(glide novendor)

testrace:
	go test -x -race $(glide novendor)

.PHONY: all deps install test testrace

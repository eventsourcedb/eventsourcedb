all: install

deps:
	glide install

install:
	go install -ldflags="-s -w" -v ./cmd/...

test:
	go test $(glide novendor)

testrace:
	go test -x -race $(glide novendor)

gen:
	go generate $(glide novendor)

.PHONY: all deps install test testrace

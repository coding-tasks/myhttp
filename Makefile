.PHONY: test build

test:
	go test -v -race ./... -count=1

build:
	go build -o myhttp ./...

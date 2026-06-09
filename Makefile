.PHONY: test run version build clean

test:
	go test ./...

run:
	go run ./cmd/shipcheck audit .

version:
	go run ./cmd/shipcheck version

build:
	go build -o bin/shipcheck ./cmd/shipcheck

clean:
	rm -rf bin

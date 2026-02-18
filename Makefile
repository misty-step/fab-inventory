.PHONY: build test lint clean

BINARY_NAME=fab-inventory

build:
	go build -o $(BINARY_NAME) .

test:
	go test ./...

lint:
	golangci-lint run ./...

clean:
	rm -f $(BINARY_NAME)

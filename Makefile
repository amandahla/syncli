# Makefile for syncli

BINARY=syncli

.PHONY: all build clean test fmt lint

all: build

build:
	go build -o $(BINARY) .

clean:
	rm -f $(BINARY)

fmt:
	gofmt -w .

lint:
	golangci-lint run

# Run tests

test:
	go test ./...

# Run syncli
run:
	./$(BINARY)

# Install dependencies
deps:
	go mod tidy

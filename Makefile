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
	@if ! command -v golangci-lint >/dev/null 2>&1; then \
		echo "WARNING: golangci-lint not installed, skipping lint."; \
	else \
		golangci-lint run; \
	fi

# Run tests
test:
	go test -v -race -buildvcs ./...

# Run syncli
run:
	./$(BINARY)

# Install dependencies
deps:
	go mod tidy
	go mod download

# Run all checks
.PHONY: audit
audit:
	$(MAKE) test
	$(MAKE) lint
	go mod tidy -diff
	go mod verify
	test -z "$(shell gofmt -l .)"
	go vet ./...
	go run honnef.co/go/tools/cmd/staticcheck@latest -checks=all,-ST1000,-U1000 ./...
	go run golang.org/x/vuln/cmd/govulncheck@latest ./...

# Coverage report
.PHONY: coverage
coverage:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out
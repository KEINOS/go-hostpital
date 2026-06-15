.PHONY: check test lint gopls-analyzers

test:
	go test -race ./...

lint:
	golangci-lint run ./...

gopls-analyzers:
	go run ./cmd/scannercheck ./...

check: test lint gopls-analyzers

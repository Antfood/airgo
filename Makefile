.PHONY: test test-verbose lint

# Run all tests
test:
	go test ./airtable/... ./retry/... ./utils/...

# Verbose test output
test-verbose:
	go test -v ./airtable/... ./retry/... ./utils/...

# Run linter (install: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest)
lint:
	golangci-lint run ./...


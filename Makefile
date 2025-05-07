.PHONY: test test-verbose test-unit test-usability test-functional coverage

# Run all tests
test:
	go clean -testcache
	go test ./...

# Run all tests with verbose output
test-verbose:
	go test ./... -v

# Run only unit tests
test-unit:
	go clean -testcache
	go test ./internal/...

# Run only functional tests
test-functional:
	go clean -testcache
	go test ./tests/functional/...

# Custom formatted test output
test-docs:
	go clean -testcache
	@echo "======================= API TEST DOCUMENTATION ======================="
	@echo "Starting test run"
	@echo "=================================================================="
	@go test ./tests/functional/... -v | grep -v "^--- PASS" | grep -v "^PASS"
	@echo "=================================================================="
	@echo "Test run completed"
	@echo "=================================================================="

# Run tests with coverage
coverage:
	go test ./... -coverprofile=coverage.out
	go tool cover -html=coverage.out

# Build the application
build:
	go build -o bin/app ./cmd/main.go

# Run the application
run:
	go run ./cmd/main.go ./cmd/server.go

# Run docker compose
run-docker:
	docker compose up -d

# Clean build artifacts
clean:
	rm -rf bin/
	rm -f coverage.out

# Default target
default: test 
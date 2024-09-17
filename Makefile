.PHONY: all build-server build-client run-server clean reset-db

# Assuming the database file is named hvpm.db and is in the root directory
DB_FILE := ./hvpm.db

all: build-server build-client

build-server:
	@echo "Building server..."
	@go build -o hvr-server ./cmd/server

build-client:
	@echo "Building client..."
	@go build -o hvr ./pkg/client

run-server: build-server
	@echo "Running server..."
	@./hvr-server

clean:
	@echo "Cleaning up..."
	@rm -f hvr-server hvr

reset-db:
	@echo "Resetting database..."
	@if [ -f $(DB_FILE) ]; then rm $(DB_FILE) && echo "Database deleted."; else echo "Database file not found."; fi

.PHONY: test

test:
	@echo "Running tests..."
	@go test ./... -v

test-coverage:
	@echo "Running tests with coverage..."
	@go test ./... -coverprofile=coverage.out
	@go tool cover -html=coverage.out

.PHONY: all build-server build-client run-server

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

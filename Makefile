BINARY_NAME=openai
BINARY_PATH=bin/$(BINARY_NAME)
GO=go

# Get the current git commit hash
GIT_COMMIT=$(shell git rev-parse --short HEAD)
# Get the current git tag if any, else use 'dev'
GIT_TAG=$(shell git describe --tags 2>/dev/null || echo "dev")

# Build flags
LDFLAGS=-ldflags "-X github.com/hmm01i/openai/pkg/version.Current=$(GIT_TAG) -X github.com/hmm01i/openai/pkg/version.Commit=$(GIT_COMMIT)"

.PHONY: all build clean test

all: clean build

build:
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p bin
	$(GO) build $(LDFLAGS) -o $(BINARY_PATH) ./cmd/...

clean:
	@echo "Cleaning..."
	@rm -rf bin/
	@go clean

test:
	$(GO) test -v ./...

# Install binary to $GOPATH/bin
install: build
	@echo "Installing $(BINARY_NAME)..."
	@cp $(BINARY_PATH) $(GOPATH)/bin/

# Run go mod tidy
tidy:
	$(GO) mod tidy

# Development target that builds and runs the application
dev: build
	./$(BINARY_PATH)
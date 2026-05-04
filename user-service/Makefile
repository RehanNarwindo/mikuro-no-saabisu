hello:
	echo "Hello"

build:
	@echo "Building user-service..."
	go build -o bin/main main.go
	@echo "Build complete: bin/main"

run:
	go run main.go

fmt:
	@echo "Formatting code..."
	go fmt ./...
	@echo "Format complete"

vet:
	@echo "Running go vet..."
	go vet ./src/...
	@echo "Vet complete"

lint:
	@echo "Running golangci-lint..."
	golangci-lint run
	@echo "Linting complete"

check: fmt vet lint 
	@echo "All checks passed"

all: hello build

clean:
	rm -rf bin/main bin/main-*
	@echo "Cleaned binaries"

test:
	go test -v ./src/...

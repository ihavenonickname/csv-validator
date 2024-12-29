.SILENT:

.PHONY: all
all: build

.PHONY: build
build: clean
	go build -o ./bin/csv-validator ./cmd/csv-validator/main.go

# Clean build directory
.PHONY: clean
clean:
	rm -rf ./bin

# Format Go code
.PHONY: fmt
fmt:
	go fmt ./...

# Test the application
.PHONY: test
test: build
	bash ./test.bash

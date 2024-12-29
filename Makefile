.SILENT:

.PHONY: all
all: build

.PHONY: build
build: clean fmt
	go build -o ./bin/csv-validator ./cmd/csv-validator/main.go

.PHONY: clean
clean:
	rm -rf ./bin

.PHONY: fmt
fmt:
	go fmt ./... > /dev/null

.PHONY: test
test: build
	bash ./test.bash

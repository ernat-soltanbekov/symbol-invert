.PHONY: build run suggest analyze audit test test-verbose golden help

build:
	go build -o symbol-invert .

# make run FILE="Path to the file containing the ASCII art"
run:
	go run . --reverse="$(FILE)"

confidence:
	go run . --reverse="$(FILE)" --confidence

tolerant:
	go run . --reverse="$(FILE)" --tolerant

test:
	go test ./...

test-verbose:
	go test -v ./...

help:
	@echo "Доступные команды:"
	@echo '  make build'
	@echo '  make run FILE="Path to the file containing the ASCII art" '
	@echo '  make confidence FILE="Path to the file containing the ASCII art" '
	@echo '  make tolerant FILE="Path to the file containing the ASCII art" '
	@echo "  make test"
	@echo "  make test-verbose"

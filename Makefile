.PHONY: test run install-cli

test:
	go test ./...

run:
	go run ./cmd/scenaria --help

install-cli:
	go install ./cmd/scenaria

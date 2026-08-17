.PHONY: build
.PHONY: test

help:
	@echo "make build   Builds the project"
	@echo "make test    Runs the tests"

build:
	@go build \
		-ldflags="-s -w" \
		-o out/welp \
		cmd/welp.go

test:
	@go test ./...


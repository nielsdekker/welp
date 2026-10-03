#!make
-include .env

.PHONY: build
.PHONY: test

help:
	@echo "make build   Builds the project"
	@echo "make test    Runs the tests"

build:
	@echo "make debug $$GOOS $$GOARCH"
	@go build \
		-ldflags="-s -w" \
		-o out/welp \
		cmd/welp.go

test:
	@go test ./...

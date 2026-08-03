.PHONY: build
.PHONY: test

help:
	@echo "help"

build:
	@go build \
		-ldflags="-s -w" \
		-o out/welp \
		main.go

test:
	@go test ./tests/...


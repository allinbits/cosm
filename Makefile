all: install

mod:
	@go mod tidy

build: mod
	@go build -mod=readonly -o build/cosmos main.go

install: build
	@go install -mod=readonly ./...

.PHONY: all build install
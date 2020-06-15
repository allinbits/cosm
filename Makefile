all: install

install: build ui
	@go install -mod=readonly ./...

mod:
	@go mod tidy

build: mod
	@go build -mod=readonly -o build/cosmos main.go

ui:
	@rm -rf ui/dist
	@which npm 1>/dev/null && cd ui && npm install 1>/dev/null && npm run build 1>/dev/null

cli: build

.PHONY: all mod build ui install cli
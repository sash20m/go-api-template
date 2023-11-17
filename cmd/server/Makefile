.DEFAULT_GOAL := build

fmt:
		go fmt ./...
.PHONY:fmt

lint: fmt
		golint ./...
.PHONY:lint

vet: lint
		go vet ./...
		shadow ./...
.PHONY:vet

swag: vet
		swag init -d ./,../../internal
.PHONY: swag

build: swag
		go build main.go
.PHONY:build
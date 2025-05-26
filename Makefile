.PHONY: build run test clean fmt vet test-api

build:
	go build -o bin/qkrn ./cmd/qkrn

run:
	go run ./cmd/qkrn

test:
	go test -v ./...

test-api:
	@echo "Make sure the server is running on port 8080"
	@./scripts/test-api.sh

clean:
	rm -rf bin/

fmt:
	go fmt ./...

vet:
	go vet ./...

check: fmt vet test

deps:
	go mod tidy

release:
	CGO_ENABLED=0 GOOS=linux go build -a -ldflags '-extldflags "-static"' -o bin/qkrn-linux ./cmd/qkrn
	CGO_ENABLED=0 GOOS=darwin go build -a -ldflags '-extldflags "-static"' -o bin/qkrn-darwin ./cmd/qkrn

dev: run

<<<<<<< HEAD
.PHONY: clean lint test build

default: clean lint test build

clean:
	rm -rf coverage.txt

build:
	go build -ldflags "-s -w" -trimpath ./cmd/recvcheck/

test: clean
=======
all: build test lint

download:
	go mod download

build: download
	go build -C cmd/recvcheck

test:
>>>>>>> 70e0318b1 ([WIP] add archivista storage backend)
	go test -race -coverprofile=coverage.txt .

lint:
	golangci-lint run


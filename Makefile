# ex : shiftwidth=2 tabstop=2 softtabstop=2 :
SHELL := /bin/sh
GOPROCS := 4
SRC := $(wildcard *.go)

.PHONY: all
all: get-deps vet cov build

.PHONY: build
build: $(SRC)
	go build ./...

.PHONY: get-deps
get-deps:
	go get -d -v ./...

.PHONY: clean
clean:
	go clean -i ./...

.PHONY: format
format:
	go fmt ./...

.PHONY: cov
cov: test
	go tool cover -func=coverage.out

.PHONY: htmlcov
htmlcov: coverage.out
	go tool cover -html=coverage.out

coverage.out:
	go test -v -covermode=count -coverprofile=coverage.out

.PHONY: test
test:
	go test ./...

.PHONY: race
race:
	go test -race ./...

.PHONY: vet
vet:
	go vet -x ./...

.PHONY: run
run: all
	./mig31 -env=dev -offline

.PHONY: install
install: cov
	go install

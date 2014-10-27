# ex : shiftwidth=2 tabstop=2 softtabstop=2 :
SHELL := /bin/sh
GOPROCS := 4

.PHONY: all
all: get-deps test
	go build

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
htmlcov: test
	go tool cover -html=coverage.out

.PHONY: test
test:
	go test -coverprofile=coverage.out ./...

.PHONY: run
run: all
	mig31 -env=dev -offline

.PHONY: install
install: cov
	go install

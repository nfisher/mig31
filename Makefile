# ex : shiftwidth=2 tabstop=2 softtabstop=2 :
SHELL  = /bin/sh
PROJECT = github.com/hailocab/mig31
EXE  = $(GOPATH)/bin/mig31
SRC  = $(wildcard *.go)
TEST = $(wildcard *_test.go)
GOPROCS=4

.PHONY: all
all: $(SRC) $(EXE) 

.PHONY: clean
clean:
	go clean -i ./...

.PHONY: format
format: $(SRC)
	go fmt ./...

.PHONY: cov
cov:
	go test -coverprofile=coverage.out
	go tool cover -func=coverage.out

.PHONY: htmlcov
htmlcov:
	go test -coverprofile=coverage.out
	go tool cover -html=coverage.out

.PHONY: test
test:
	go test ./...

.PHONY: run
run: $(EXE)
	$(EXE) -env=dev -offline

$(EXE): test
	go install $(PROJECT)

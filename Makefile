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
	go test -coverprofile=coverage.out $(TEST) $(PROJECT)

.PHONY: test
test: test-main test-config

.PHONY: test-main
test-main:
	go test $(PROJECT)

.PHONY: test-config
test-config:
	go test $(PROJECT)/config

.PHONY: run
run: $(EXE)
	$(EXE)

$(EXE): test
	go install $(PROJECT)

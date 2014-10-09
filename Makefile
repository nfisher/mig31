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
test:
	go test ./...

.PHONY: test-main
test-main:
	go test $(PROJECT)

.PHONY: test-config
test-config:
	go test $(PROJECT)/config

.PHONY: test-migration
test-migration:
	go test $(PROJECT)/migration

.PHONY: run
run: $(EXE)
	$(EXE) -environment=dev -offline

$(EXE): test
	go install $(PROJECT)

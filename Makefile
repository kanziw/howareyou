GOPATH:=$(shell go env GOPATH)
APP?=howareyou

.PHONY: init
## init: initialize the application
init:
	go mod download

.PHONY: build
## build: build the application
build:
	go build -o build/howareyou cmd/main.go

.PHONY: run
## run: run the application
run:
	go run -v -race cmd/main.go

.PHONY: format
## format: format files
format:
	@go get -d golang.org/x/tools/cmd/goimports
	goimports -local github.com/kanziw -w .
	gofmt -s -w .
	go mod tidy

.PHONY: test
## test: run tests
test:
	@go install github.com/rakyll/gotest
	gotest -p 1 -race -cover -v ./...

.PHONY: coverage
## coverage: run tests with coverage
coverage:
	@go install github.com/rakyll/gotest
	gotest -p 1 -race -coverprofile=coverage.txt -covermode=atomic -v ./...

.PHONY: lint
## lint: check everything's okay
lint:
	golangci-lint run ./...
	go mod verify

.PHONY: generate
## generate: generate source code for mocking
generate:
	@go get -d golang.org/x/tools/cmd/stringer
	@go get -d github.com/golang/mock/gomock
	@go install github.com/golang/mock/mockgen
	go generate ./...
	$(MAKE) format

.PHONY: help
## help: prints this help message
help:
	@echo "Usage: \n"
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':'

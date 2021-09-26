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
	@go install golang.org/x/tools/cmd/goimports@v0.1.6
	@go install github.com/aristanetworks/goarista/cmd/importsort@latest
	goimports -local github.com/kanziw -w .
	importsort -s github.com/kanziw -w $$(find . -name "*.go")
	gofmt -s -w .
	go mod tidy

.PHONY: test
## test: run tests
test:
	@go install github.com/rakyll/gotest@v0.0.6
	gotest -p 1 -race -cover -v ./...

.PHONY: coverage
## coverage: run tests with coverage
coverage:
	@go install github.com/rakyll/gotest@v0.0.6
	gotest -p 1 -race -coverprofile=coverage.txt -covermode=atomic -v ./...

.PHONY: lint
## lint: check everything's okay
lint:
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.42.1
	golangci-lint run ./...

.PHONY: generate
## generate: generate source code for mocking and DB models
generate:
	@go install golang.org/x/tools/cmd/stringer@v0.1.6
	@go install github.com/golang/mock/mockgen@v1.6.0
	go generate ./...
	@go install github.com/volatiletech/sqlboiler/v4@v4.6.0
	@go install github.com/volatiletech/sqlboiler/v4/drivers/sqlboiler-mysql@v4.6.0
	sqlboiler --wipe --no-tests -p model -o ./model mysql
	$(MAKE) format

.PHONY: help
## help: prints this help message
help:
	@echo "Usage: \n"
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':'

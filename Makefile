.PHONY: all fmt test build-mock build
all: help

fmt:
	@echo "==> Fixing source code with gofmt..."
	gofmt -s -w ./

test: build-mock fmt
	go test ./... -v -count 1 -parallel 1 -race -coverprofile=coverage.txt -covermode=atomic $(TESTARGS) -timeout 600s

build-mock:
	go install github.com/golang/mock/mockgen@v1.6.0
	mockgen -destination=harbor/mocks/mock_api.go -package=mocks github.com/disaster37/harbor-cli/harbor/api API,ArtifactAPI

build: fmt
ifeq ($(OS),Windows_NT)
	go build -o  harbor-cli.exe
else
	go build -o harbor-cli
endif
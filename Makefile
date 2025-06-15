PKG := "github.com/infraboard/modules"
MOD_DIR := $(shell go env GOPATH)/pkg/mod
PKG_LIST := $(shell go list ${PKG}/... | grep -v /vendor/ | grep -v redis)
GO_FILES := $(shell find . -name '*.go' | grep -v /vendor/ | grep -v _test.go)
GIT_TAG := $(shell git describe --tags --exact-match 2>/dev/null || echo "UnTag")

.PHONY: all dep lint vet test test-coverage build clean

all: build

dep: ## Get the dependencies
	@go mod tidy

lint: ## Lint Golang files
	@golint -set_exit_status ${PKG_LIST}

vet: ## Run go vet
	@go vet ${PKG_LIST}

test: ## Run unittests
	@go test -short ${PKG_LIST}
	
test-coverage: ## Run tests with coverage
	@go test -short -coverprofile cover.out -covermode=atomic ${PKG_LIST} 
	@cat cover.out >> coverage.txt

install: ## Install depence go package
	@go install github.com/infraboard/mcube/v2/cmd@latest
	@go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	@go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	@go install github.com/favadi/protoc-go-inject-tag@latest

gen: ## Init Service
	@protoc -I=.. --go_out=. --go_opt=module=${PKG} --go-grpc_out=. --go-grpc_opt=module=${PKG} ../modules/*/apps/*/pb/*.proto
	@go fmt ./...

	@protoc-go-inject-tag -input=../modules/*/apps/*/*.pb.go
	@mcube generate enum -p -m ../modules/*/apps/*/*.pb.go


help: ## Display this help screen
	@grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
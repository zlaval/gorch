.DEFAULT_GOAL := help
.PHONY:help build clean lint flint tools init deps new-project

help: ## Display help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

build: ## Build the project
	go fmt ./...
	go build

clean: ## Clean-up dependencies
	go mod tidy

lint: ## Run static code analysis
	golangci-lint run
	govulncheck ./...

flint: ## Fix linter errors
	golangci-lint run --fix

tools: ## Install command line tools
	go install golang.org/x/vuln/cmd/govulncheck@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
    # Such go install/go get or "tools pattern" installations aren't guaranteed to work. Use the command below if there is any problem
    # curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.62.2

init: ## Init project
	go mod init gorch

deps: ## install all dependencies
	go get -u github.com/go-chi/chi/v5
	go get github.com/docker/docker/client
	go get github.com/shirou/gopsutil
	go get github.com/spf13/cobra@latest
	go get github.com/spf13/viper
	go get github.com/renstrom/shortuuid
	go get gopkg.in/yaml.v2
	go get go.etcd.io/bbolt@latest
	go get github.com/yusufpapurcu/wmi

new-project: init deps tools ## Create new project and install all dependencies
APP = api-gateway
GOBASE = $(shell pwd)
GOBIN = $(GOBASE)/build/bin
LINT_PATH = $(GOBASE)/build/lint
MAIN_APP = $(GOBASE)/cmd

###  mockgen -source=producer/producer.go -destination=producer/mocks/mock_producer.go -package=mocks

####  mockgen -source=consumer/consumer.go -destination=consumer/mocks/mock_consumer.go -package=mocks

help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

deps: ## Fetch required dependencies
	go mod tidy -compat=1.22
	go mod download

run: build ## Build and run program
	cd $(MAIN_APP) && go run .

lint: install-golangci ## Linter for developers
	$(LINT_PATH)/golangci-lint run --timeout=5m -c .golangci.yml

lint-fix:
	$(LINT_PATH)/golangci-lint run --timeout=5m -c .golangci.yml --fix

install-golangci: ## Install the correct version of lint
	@GOBIN=$(LINT_PATH) go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.58.1
	
# run-tests: ## Run tests 
# 	cd $(TEST_PATH) && go test .
# test-cover: ## Run tests with coverage
# 	cd $(TEST_PATH) && go test -cover

# test-coverage: ## Run tests and generate coverage profile
# 	cd $(TEST_PATH) && go test -coverprofile=coverage.out

# test-coverage-browser: ## Check the test coverage in the browser
# 	cd $(TEST_PATH) && go tool cover -html=coverage.out -o /tmp/coverage.html && wslview /tmp/coverage.html

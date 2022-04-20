build:
	go build

test:
	go fmt $(go list ./... | grep -v /vendor/)
	go vet $(go list ./... | grep -v /vendor/)
	go test -race $(go list ./... | grep -v /vendor/)

lint:
	golangci-lint run

lint.fix: ## Lint
	golangci-lint run --fix


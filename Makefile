.PHONY: dep, generate

## Ensure all dependencies are up to date
deps:
	@go mod tidy && go mod download

## Generate code
generate:
	GO_CMD_MOCKGEN=$(GO_CMD_MOCKGEN)
	@go generate ./...

## Lint
lint:
	@echo "Running golangci-lint..."
	@golangci-lint run ./...

unit_tests:
	@echo "Running unit tests..."
	@go test -race ./...

ci:
	@$(MAKE) lint
	@$(MAKE) unit_tests
.PHONY: dep generate lint unit_tests ci

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

service_tests:
	ginkgo -p \
		-coverpkg=./... \
		-coverprofile=service_coverage.out \
 		-tags="service_test" ./test/service/...

ci:
	@$(MAKE) lint
	@$(MAKE) unit_tests
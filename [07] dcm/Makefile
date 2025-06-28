START_LOG = @echo "======================= START OF LOG ======================="
END_LOG = @echo "======================== END OF LOG ======================="

.PHONY: build
build: ## Build the application RISC-V image with cartesi cli
	$(START_LOG)
	@cartesi build
	$(END_LOG)
	
.PHONY: generate
generate: ## Generate the application code
	$(START_LOG)
	@go generate ./...
	$(END_LOG)

.PHONY: test
test: ## Run the application tests
	$(START_LOG)
	@go test -p=1 ./... -coverprofile=./coverage.md -v
	$(END_LOG)

.PHONY: coverage
coverage: test ## Generate the application code coverage report
	$(START_LOG)
	@go tool cover -html=./coverage.md
	$(END_LOG)

.PHONY: state
state: ## Run the application state for devnet (demo)
	@chmod +x ./tools/state.sh
	@./tools/state.sh

.PHONY: state-dev
state-dev: ## Run the application state for devnet (demo)
	@chmod +x ./tools/state-dev.sh
	@./tools/state-dev.sh

.PHONY: state-mcp
state-mcp: ## Run the application state for devnet (demo)
	@chmod +x ./tools/state-mcp.sh
	@./tools/state-mcp.sh $(DAPP_ADDRESS)

.PHONY: contracts
contracts: ## Deploy the contracts
	@forge script ./contracts/script/Deploy.s.sol \
			--root ./contracts \
			--rpc-url http://localhost:8080/anvil \
			--private-key 0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80 \
			--broadcast \
			-vvv

.PHONY: help
help: ## Show help for each of the Makefile recipes
	@grep "##" $(MAKEFILE_LIST) | grep -v grep | sed -e 's/:.*##/:\t/'
.PHONY: run
run:
	@go run ./cmd/binance/

.PHONY: lint
lint:
	@golangci-lint run --new --fix

.PHONY: lint_%
lint_%:
	@golangci-lint run --new-from-rev=HEAD~$*

.PHONY: rmf
rmf:
	@echo "Removing files..."
	@rm -rf internal
	@rm -rf pkg
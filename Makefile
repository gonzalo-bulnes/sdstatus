.DEFAULT_GOAL := help

.PHONY: build
build: ## Compiles Golang binary.
	go build

.PHONY: test
test: ## Runs the test suite.
	go test -cover ./...

.PHONY: help
help: ## Prints this message and exits.
	@printf "Makefile for developing and testing FPF infrastructure.\n"
	@printf "Subcommands:\n\n"
	@perl -F':.*##\s+' -lanE '$$F[1] and say "\033[36m$$F[0]\033[0m : $$F[1]"' $(MAKEFILE_LIST) \
		| sort \
		| column -s ':' -t

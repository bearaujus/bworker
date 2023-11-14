.PHONY: test-no-race
test-no-race:
	@go clean -testcache
	go test ./... --cover

.PHONY: test-race
test-race:
	@go clean -testcache
	go test ./... --cover --race

.PHONY: test
test:
	@go mod tidy
	@make test-no-race
	@make test-race

.PHONY: test-sync
test-sync:
	@go clean -testcache
	go test ./... --cover

.PHONY: test-async
test-async:
	@go clean -testcache
	go test ./... --cover --race

.PHONY: test
test:
	@make test-sync
	@make test-async

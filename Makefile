.PHONY: build
build:
	@go build -o dependabot-templater main.go

.PHONY: lint
lint:
	@golangci-lint run ./...

.PHONY: test
test:
	@go clean -testcache
	@go test -v -cover -coverprofile coverage.txt -race ./...
	@echo
	@go tool cover -func coverage.txt

.PHONY: install
install: build
	go install

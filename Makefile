.PHONY: install lint test coverage test-report

install:
	go mod tidy

lint:
	@echo "== 🙆 linter =="
	golangci-lint run -v ./... --fix

test:
	go test ./... -race

test-report:
	go test ./... -coverprofile coverage.out
	go tool cover -html=coverage.out -o coverage.html
	open coverage.html

coverage:
	go tool cover -html=coverage.out

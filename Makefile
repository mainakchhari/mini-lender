build:
	@echo "Building go binary for mini-lender..."
	go build -o bin/app -v cmd/app/main.go
	@echo "Done."

run:
	@echo "Running mini-lender..."
	@GIN_MODE=release ./bin/app

test:
	@echo "Running tests..."
	@go test -v ./...
	@echo "Done."

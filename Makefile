generate-docs:
	swag init -g cmd/currency-rates/main.go

run:
	go run cmd/currency-rates/main.go

linters:
	golangci-lint run
	go mod tidy
	goimports -w .


generate-mocks:
	mockery --all -r
test:
	go test ./... -v -cover

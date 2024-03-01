generate-docs:
	swag init -g cmd/currency-rates/main.go

run:
	go run cmd/currency-rates/main.go

linters:
	golangci-lint run
	goimports -w .

generate-mocks:
	mockery --all -r
test:
	go test ./... -v -cover

start:
	docker-compose up -d
	docker exec -i -t currency-rates-app-1 bash -c "chmod +x wait-for-postgres.sh; sh wait-for-postgres.sh db app; GOOSE_DRIVER=postgres GOOSE_DBSTRING=\"user=postgres password=qwerty dbname=postgres sslmode=disable host=db port=5432\" goose -dir=./migrations up"

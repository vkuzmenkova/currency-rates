generate-docs:
	swag init -g cmd/currency-rates/main.go

run:
	export DEPLOYMENT=local
	go run cmd/currency-rates/main.go

work:
	workwebui -redis=":6379" -ns="currency_rates" -listen=":5040"


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

run-postgres:
	docker run --name postgres -p 5432:5432 -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=qwerty -e POSTGRES_DB=postgres -d postgres:13.3

apply-migrations:
	GOOSE_DRIVER=postgres GOOSE_DBSTRING="user=postgres password=qwerty dbname=postgres sslmode=disable host=localhost port=5432" goose -dir=./migrations up
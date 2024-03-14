FROM golang:1.21
WORKDIR /usr/src/app

ENV DEPLOYMENT=container

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY cmd/currency-rates/main.go ./cmd/currency-rates/
COPY configs ./configs
COPY docs ./docs
COPY internal ./internal
COPY middleware ./middleware
COPY migrations ./migrations
COPY models ./models
COPY wait-for-postgres.sh ./wait-for-postgres.sh

# install psql
RUN apt-get update
RUN apt-get -y install postgresql-client
RUN go install github.com/pressly/goose/v3/cmd/goose@latest

# make wait-for-postgres.sh executable
RUN chmod +x wait-for-postgres.sh

RUN go build -v -o /usr/local/bin/app ./cmd/currency-rates/

EXPOSE 8080

CMD ["app"]
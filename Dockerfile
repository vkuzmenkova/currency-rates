FROM golang:1.21
WORKDIR /usr/src/app

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY cmd/currency-rates/main.go ./cmd/currency-rates/
COPY configs ./configs
COPY docs ./docs
COPY internal ./internal
COPY middleware ./middleware
COPY migrations ./migrations
COPY models ./models
RUN go build -v -o /usr/local/bin/app ./cmd/currency-rates/

RUN go install github.com/pressly/goose/v3/cmd/goose@latest

EXPOSE 8080

CMD ["app"]
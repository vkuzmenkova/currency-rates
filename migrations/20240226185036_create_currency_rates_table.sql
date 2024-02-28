-- +goose Up
-- +goose StatementBegin
CREATE TABLE currency_rates (
    uuid UUID PRIMARY KEY,
    base SMALLINT,
    currency SMALLINT,
    rate REAL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS currency_rates;
-- +goose StatementEnd

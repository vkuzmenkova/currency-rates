-- +goose Up
-- +goose StatementBegin
CREATE TABLE currencies (
    id SERIAL PRIMARY KEY,
    currency_code VARCHAR(3) UNIQUE NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
INSERT INTO currencies (id, currency_code) VALUES
(1, 'USD'),
(2, 'EUR'),
(3, 'MXN');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS currencies;
-- +goose StatementEnd

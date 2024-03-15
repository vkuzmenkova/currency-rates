-- +goose Up
-- +goose StatementBegin
INSERT INTO currencies (id, currency_code) VALUES
(4, 'JPY'),
(5, 'GBP'),
(6, 'CHF');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- +goose StatementEnd

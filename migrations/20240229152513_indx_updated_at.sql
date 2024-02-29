-- +goose Up
-- +goose StatementBegin
CREATE INDEX idx_updated_at ON currency_rates(updated_at);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX idx_updated_at;
-- +goose StatementEnd

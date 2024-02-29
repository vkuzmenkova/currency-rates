-- +goose Up
-- +goose StatementBegin
CREATE INDEX idx_uuid ON currency_rates(uuid);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX idx_uuid;
-- +goose StatementEnd

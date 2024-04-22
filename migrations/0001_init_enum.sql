-- +goose Up
-- +goose StatementBegin
CREATE TYPE metric_type AS ENUM ('gauge', 'counter');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TYPE IF EXISTS metric_type;
-- +goose StatementEnd

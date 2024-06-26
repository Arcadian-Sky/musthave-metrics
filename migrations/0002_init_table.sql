-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS metrics (
			name varchar NOT NULL,
			counter numeric NULL,
			gauge double precision NULL,
			type metric_type  NOT NULL,
			CONSTRAINT constraint_name_type UNIQUE (name, type)
		);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE metrics;
-- +goose StatementEnd

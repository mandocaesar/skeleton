-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS event_store_data (
    event_id VARCHAR(255),
    aggregate_id VARCHAR(255),
    event_type VARCHAR(255),
    aggregate_type VARCHAR(255),
    version INT,
    metadata VARCHAR(255),
    status VARCHAR(255),
    timestamp TIMESTAMP,
    data TEXT,
    error VARCHAR(255)
);


-- +goose StatementEnd

-- +goose Down

-- +goose StatementBegin
DROP TABLE  IF EXISTS event_store_data
-- +goose StatementEnd
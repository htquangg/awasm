-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd
CREATE TABLE IF NOT EXISTS endpoints
(
    id                   VARCHAR(36) PRIMARY KEY             NOT NULL,
    created_at           TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at           TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    deleted_at           TIMESTAMP                           NULL,

    name                 VARCHAR(50)                         NOT NULL,
    runtime              VARCHAR(50)                         NOT NULL,
    active_deployment_id VARCHAR(36)                         NOT NULL
);

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
DROP TABLE if EXISTS endpoints;

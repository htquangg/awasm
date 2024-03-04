-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd
CREATE TABLE IF NOT EXISTS endpoints
(
    id                   VARCHAR(36) PRIMARY KEY NOT NULL,
    created_at           TIMESTAMP DEFAULT NOW() NOT NULL,
    updated_at           TIMESTAMP DEFAULT NOW() NOT NULL,
    deleted_at           TIMESTAMP               NULL,

    name                 VARCHAR(50)             NOT NULL,
    runtime              VARCHAR(50)             NOT NULL,
    active_deployment_id VARCHAR(36)             NOT NULL
);

CREATE TABLE IF NOT EXISTS deployments
(
    id          VARCHAR(36) PRIMARY KEY NOT NULL,
    created_at  TIMESTAMP DEFAULT NOW() NOT NULL,
    deleted_at  TIMESTAMP               NULL,

    endpoint_id VARCHAR(36)             NOT NULL,
    hash        CHAR(32)                NOT NULL,
    data        BYTEA                   NOT NULL
);

ALTER TABLE deployments
    ADD CONSTRAINT fk_deployment_endpoint_id FOREIGN KEY (endpoint_id) REFERENCES endpoints (id) ON DELETE NO ACTION ON UPDATE NO ACTION;

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
ALTER TABLE deployments
    DROP CONSTRAINT fk_deployment_endpoint_id;

DROP TABLE IF EXISTS deployments;

DROP TABLE IF EXISTS endpoints;

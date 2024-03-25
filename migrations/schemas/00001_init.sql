-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd
CREATE TABLE IF NOT EXISTS users
(
    id                     VARCHAR(36) PRIMARY KEY   NOT NULL,
    created_at             TIMESTAMPTZ DEFAULT NOW() NOT NULL,
    updated_at             TIMESTAMPTZ DEFAULT NOW() NOT NULL,
    deleted_at             TIMESTAMPTZ               NULL,

    name                   VARCHAR(256)              NOT NULL,
    email_hash             TEXT UNIQUE               NOT NULL,
    encrypted_email        BYTEA                     NOT NULL,
    email_decryption_nonce BYTEA                     NOT NULL,
    email_confirmed_at     TIMESTAMPTZ               NULL,

    last_login_at          TIMESTAMPTZ               NULL
);

CREATE TABLE IF NOT EXISTS endpoints
(
    id                   VARCHAR(36) PRIMARY KEY   NOT NULL,
    created_at           TIMESTAMPTZ DEFAULT NOW() NOT NULL,
    updated_at           TIMESTAMPTZ DEFAULT NOW() NOT NULL,
    deleted_at           TIMESTAMPTZ               NULL,

    name                 VARCHAR(256)              NOT NULL,
    runtime              VARCHAR(64)               NOT NULL,
    active_deployment_id VARCHAR(36)               NOT NULL
);

CREATE TABLE IF NOT EXISTS deployments
(
    id          VARCHAR(36) PRIMARY KEY   NOT NULL,
    created_at  TIMESTAMPTZ DEFAULT NOW() NOT NULL,
    deleted_at  TIMESTAMPTZ               NULL,

    endpoint_id VARCHAR(36)               NOT NULL,
    hash        CHAR(32)                  NOT NULL,
    data        BYTEA                     NOT NULL,
    CONSTRAINT fk_deployment_endpoint_id
        FOREIGN KEY (endpoint_id)
            REFERENCES endpoints (id)
            ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS srp_auth
(
    id          VARCHAR(36) PRIMARY KEY   NOT NULL,
    created_at  TIMESTAMPTZ DEFAULT NOW() NOT NULL,
    updated_at  TIMESTAMPTZ DEFAULT NOW() NOT NULL,

    user_id     VARCHAR(36)               NOT NULL,
    srp_user_id VARCHAR(36)               NOT NULL,
    salt        TEXT                      NOT NULL,
    verifier    TEXT                      NOT NULL,
    CONSTRAINT fk_srp_auth_user_id
        FOREIGN KEY (user_id)
            REFERENCES users (id)
            ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS srp_auth_temp
(
    id          VARCHAR(36) PRIMARY KEY   NOT NULL,
    created_at  TIMESTAMPTZ DEFAULT NOW() NOT NULL,

    user_id     VARCHAR(36)               NOT NULL,
    srp_user_id VARCHAR(36)               NOT NULL,
    salt        TEXT                      NOT NULL,
    verifier    TEXT                      NOT NULL,
    CONSTRAINT fk_srp_auth_temp_user_id
        FOREIGN KEY (user_id)
            REFERENCES users (id)
            ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS srp_challenges
(
    id               VARCHAR(36) PRIMARY KEY   NOT NULL,
    created_at       TIMESTAMPTZ DEFAULT NOW() NOT NULL,
    updated_at       TIMESTAMPTZ DEFAULT NOW() NOT NULL,

    srp_auth_temp_id VARCHAR(36)               NOT NULL,
    srp_user_id      VARCHAR(36)               NOT NULL,
    server_key       TEXT                      NOT NULL,
    srp_a            TEXT                      NOT NULL,
    verified_at      TIMESTAMPTZ               NULL,
    attempt_count    INT                       NOT NULL DEFAULT 0,
    CONSTRAINT fk_srp_challenges_srp_auth_temp_id
        FOREIGN KEY (srp_auth_temp_id)
            REFERENCES srp_auth_temp (id)
            ON DELETE CASCADE
);

CREATE TYPE aal_level_enum AS ENUM ('aal0','aal1', 'aal2', 'aal3');

CREATE TABLE IF NOT EXISTS tokens
(
    created_at   TIMESTAMPTZ DEFAULT NOW() NOT NULL,
    deleted_at   TIMESTAMPTZ               NULL,
    last_used_at TIMESTAMPTZ DEFAULT NOW() NOT NULL,

    user_id      VARCHAR(36)               NOT NULL,
    token        TEXT UNIQUE               NOT NULL,
    aal          aal_level_enum            NOT NULL DEFAULT 'aal0',
    ip           TEXT                      NULL,
    user_agent   TEXT                      NULL,
    CONSTRAINT fk_tokens_user_id
        FOREIGN KEY (user_id)
            REFERENCES users (id)
            ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_users_email_hash ON users (email_hash);

CREATE INDEX IF NOT EXISTS idx_deployments_endpoint_id ON deployments (endpoint_id);

CREATE INDEX IF NOT EXISTS idx_tokens_user_id ON tokens (user_id);

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
DROP TABLE IF EXISTS srp_challenges;

DROP INDEX IF EXISTS idx_tokens_user_id;

DROP TABLE IF EXISTS tokens;

DROP TYPE IF EXISTS aal_level_enum;

ALTER TABLE srp_auth_temp
    DROP CONSTRAINT fk_srp_auth_temp_user_id;

DROP TABLE IF EXISTS srp_auth_temp;

ALTER TABLE srp_auth
    DROP CONSTRAINT fk_srp_auth_user_id;

DROP TABLE IF EXISTS srp_auth;

DROP INDEX IF EXISTS idx_deployments_endpoint_id;

ALTER TABLE deployments
    DROP CONSTRAINT fk_deployment_endpoint_id;

DROP TABLE IF EXISTS deployments;

DROP TABLE IF EXISTS endpoints;

DROP INDEX idx_users_email_hash;

DROP TABLE IF EXISTS users;

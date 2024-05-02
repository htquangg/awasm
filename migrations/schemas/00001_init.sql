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
    email_accepted_at      TIMESTAMPTZ               NULL,

    last_login_at          TIMESTAMPTZ               NULL
);

CREATE INDEX IF NOT EXISTS idx_users_email_hash ON users (email_hash);

CREATE TABLE IF NOT EXISTS api_keys
(
    id            VARCHAR(36) PRIMARY KEY   NOT NULL,
    created_at    TIMESTAMPTZ DEFAULT NOW() NOT NULL,
    last_used_at  TIMESTAMPTZ               NULL,

    user_id       VARCHAR(36)               NOT NULL,
    key           VARCHAR(512)              NOT NULL,
    key_preview   VARCHAR(32)               NOT NULL,
    friendly_name TEXT                      NULL,
    CONSTRAINT fk_api_keys_user_id
        FOREIGN KEY (user_id)
            REFERENCES users (id)
            ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS endpoints
(
    id                   VARCHAR(36) PRIMARY KEY   NOT NULL,
    created_at           TIMESTAMPTZ DEFAULT NOW() NOT NULL,
    updated_at           TIMESTAMPTZ DEFAULT NOW() NOT NULL,
    deleted_at           TIMESTAMPTZ               NULL,

    user_id              VARCHAR(36)               NOT NULL,
    name                 VARCHAR(256)              NOT NULL,
    runtime              VARCHAR(64)               NOT NULL,
    active_deployment_id VARCHAR(36)               NOT NULL
);

CREATE TABLE IF NOT EXISTS deployments
(
    id          VARCHAR(36) PRIMARY KEY   NOT NULL,
    created_at  TIMESTAMPTZ DEFAULT NOW() NOT NULL,
    deleted_at  TIMESTAMPTZ               NULL,

    user_id     VARCHAR(36)               NOT NULL,
    endpoint_id VARCHAR(36)               NOT NULL,
    hash        CHAR(32)                  NOT NULL,
    data        BYTEA                     NOT NULL,
    CONSTRAINT fk_deployment_endpoint_id
        FOREIGN KEY (endpoint_id)
            REFERENCES endpoints (id)
            ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_deployments_endpoint_id ON deployments (endpoint_id);

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
    id               VARCHAR(36) PRIMARY KEY   NOT NULL,
    created_at       TIMESTAMPTZ DEFAULT NOW() NOT NULL,

    user_id          VARCHAR(36)               NOT NULL,
    srp_user_id      VARCHAR(36)               NOT NULL,
    salt             TEXT                      NOT NULL,
    verifier         TEXT                      NOT NULL,
    srp_challenge_id VARCHAR(36)               NOT NULL,
    CONSTRAINT fk_srp_auth_temp_user_id
        FOREIGN KEY (user_id)
            REFERENCES users (id)
            ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS srp_challenges
(
    id            VARCHAR(36) PRIMARY KEY NOT NULL,
    created_at    TIMESTAMPTZ                      DEFAULT NOW() NOT NULL,

    srp_user_id   VARCHAR(36)             NOT NULL,
    server_key    TEXT                    NOT NULL,
    srp_a         TEXT                    NOT NULL,
    verified_at   TIMESTAMPTZ             NULL,
    attempt_count INT                     NOT NULL DEFAULT 0
);

-- +goose StatementBegin
DO
$$
    BEGIN
        CREATE TYPE factor_type AS ENUM ('totp', 'webauthn');
        CREATE TYPE factor_status AS ENUM ('unverified', 'verified');
        CREATE TYPE aal_level_enum AS ENUM ('aal0','aal1', 'aal2', 'aal3');
    EXCEPTION
        WHEN duplicate_object then null;
    END
$$
language plpgsql;
-- +goose StatementEnd


CREATE TABLE IF NOT EXISTS mfa_factors
(
    id            VARCHAR(36) PRIMARY KEY   NOT NULL,
    created_at    TIMESTAMPTZ DEFAULT NOW() NOT NULL,
    last_used_at  TIMESTAMPTZ               NULL,

    user_id       VARCHAR(36)               NOT NULL,
    status        factor_status             NOT NULL,
    friendly_name TEXT                      NULL,
    factor_type   factor_type               NOT NULL,
    secret        TEXT                      NULL,
    CONSTRAINT fk_mfa_factors_user_id
        FOREIGN KEY (user_id)
            REFERENCES users (id)
            ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS mfa_challenges
(
    id          VARCHAR(36) PRIMARY KEY   NOT NULL,
    created_at  TIMESTAMPTZ DEFAULT NOW() NOT NULL,
    verified_at TIMESTAMPTZ               NULL,

    factor_id   VARCHAR(36)               NOT NULL,
    ip          TEXT                      NULL,
    CONSTRAINT fk_mfa_challenges_factor_id
        FOREIGN KEY (factor_id)
            REFERENCES mfa_factors (id)
            ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS sessions
(
    id           VARCHAR(36) PRIMARY KEY NOT NULL,
    created_at   TIMESTAMPTZ                      DEFAULT NOW() NOT NULL,
    deleted_at   TIMESTAMPTZ             NULL,
    last_used_at TIMESTAMPTZ             NULL,

    user_id      VARCHAR(36)             NOT NULL,
    factor_id    VARCHAR(36)             NULL,
    aal          aal_level_enum          NOT NULL DEFAULT 'aal0',
    ip           TEXT                    NULL,
    user_agent   TEXT                    NULL,
    not_after    TIMESTAMPTZ             NULL,
    CONSTRAINT fk_sessions_user_id
        FOREIGN KEY (user_id)
            REFERENCES users (id)
            ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_sessions_user_id ON sessions (user_id);


CREATE TABLE IF NOT EXISTS mfa_amr_claims
(
    created_at            TIMESTAMPTZ DEFAULT NOW() NOT NULL,
    updated_at            TIMESTAMPTZ DEFAULT NOW() NOT NULL,

    session_id            VARCHAR(36)               NOT NULL,
    authentication_method TEXT                      not null,
    CONSTRAINT pk_mfa_amr_claims_session_id_authentication_method UNIQUE (session_id, authentication_method),
    CONSTRAINT fk_mfa_amr_claims_session_id
        FOREIGN KEY (session_id)
            REFERENCES sessions (id)
            ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS refresh_tokens
(
    id         VARCHAR(36) PRIMARY KEY   NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW() NOT NULL,
    updated_at TIMESTAMPTZ DEFAULT NOW() NOT NULL,
    token      varchar(255)              NOT NULL,
    user_id    varchar(36)               NOT NULL,
    session_id varchar(36)               NOT NULL,
    revoked    bool                      NULL,
    CONSTRAINT fk_refresh_tokens_user_id
        FOREIGN KEY (user_id)
            REFERENCES users (id)
            ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_refresh_tokens_token ON refresh_tokens (token);

CREATE TABLE IF NOT EXISTS key_attributes
(
    user_id                                VARCHAR(36) PRIMARY KEY,
    kek_salt                               TEXT                 NOT NULL,
    encrypted_key                          TEXT                 NOT NULL,
    key_decryption_nonce                   TEXT                 NOT NULL,
    public_key                             TEXT                 NOT NULL,
    encrypted_secret_key                   TEXT                 NOT NULL,
    secret_key_decryption_nonce            TEXT                 NOT NULL,
    master_key_encrypted_with_recovery_key TEXT                 NOT NULL,
    master_key_decryption_nonce            TEXT                 NOT NULL,
    recovery_key_encrypted_with_master_key TEXT                 NOT NULL,
    recovery_key_decryption_nonce          TEXT                 NOT NULL,
    mem_limit                              INT DEFAULT 67108864 NOT NULL, -- crypto_pwhash_MEMLIMIT_INTERACTIVE
    ops_limit                              INT DEFAULT 2        NOT NULL, -- crypto_pwhash_OPSLIMIT_INTERACTIVE
    CONSTRAINT fk_key_attributes_user_id
        FOREIGN KEY (user_id)
            REFERENCES users (id)
            ON DELETE CASCADE
);

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
DROP TABLE IF EXISTS key_attributes;

DROP INDEX IF EXISTS idx_refresh_tokens_token;

DROP TABLE IF EXISTS refresh_tokens;

DROP TABLE IF EXISTS mfa_amr_claims;

DROP INDEX IF EXISTS idx_sessions_user_id;

DROP TABLE IF EXISTS sessions;

DROP TABLE IF EXISTS mfa_challenges;

DROP TABLE IF EXISTS mfa_factors;

DROP TYPE IF EXISTS factor_type;
DROP TYPE IF EXISTS factor_status;
DROP TYPE IF EXISTS aal_level_enum;

DROP TABLE IF EXISTS srp_challenges;

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

DROP TABLE IF EXISTS api_keys;

DROP TABLE IF EXISTS users;

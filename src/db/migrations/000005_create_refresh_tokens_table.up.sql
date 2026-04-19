CREATE TABLE refresh_tokens (
    id         BIGSERIAL PRIMARY KEY,
    session_id UUID      NOT NULL REFERENCES sessions(id) ON DELETE CASCADE,
    token_hash TEXT      NOT NULL UNIQUE,
    used       BOOLEAN   NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    expires_at TIMESTAMP NOT NULL
);

CREATE INDEX idx_refresh_tokens_session_id ON refresh_tokens(session_id);
CREATE INDEX idx_refresh_tokens_token_hash ON refresh_tokens(token_hash);

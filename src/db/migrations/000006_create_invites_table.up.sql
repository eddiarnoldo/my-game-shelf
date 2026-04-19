CREATE TABLE invites (
    id         BIGSERIAL    PRIMARY KEY,
    code       VARCHAR(36)  NOT NULL UNIQUE,
    email      VARCHAR(255) NOT NULL,
    created_by BIGINT       NOT NULL REFERENCES users(id),
    used       BOOLEAN      NOT NULL DEFAULT FALSE,
    used_by    BIGINT       REFERENCES users(id),
    expires_at TIMESTAMP    NOT NULL,
    created_at TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_invites_code  ON invites(code);
CREATE INDEX idx_invites_email ON invites(email);

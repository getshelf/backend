CREATE TABLE sessions (
    token_hash BYTEA PRIMARY KEY,
    account_id TEXT NOT NULL,
    expires_at TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ NOT NULL,

    CONSTRAINT sessions_token_hash_length
        CHECK (octet_length(token_hash) = 32),

    CONSTRAINT sessions_account_id_fkey
        FOREIGN KEY (account_id)
        REFERENCES accounts (id)
        ON DELETE CASCADE
);

CREATE INDEX sessions_account_id_idx
    ON sessions (account_id);

CREATE INDEX sessions_expires_at_idx
    ON sessions (expires_at);

CREATE TABLE
  IF NOT EXISTS user_credential (
    id BIGSERIAL PRIMARY KEY,
    status VARCHAR(1) NOT NULL DEFAULT 'A',
    uuid UUID NOT NULL DEFAULT gen_random_uuid (),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW (),
    updated_at TIMESTAMPTZ DEFAULT NOW (),
    deleted_at TIMESTAMPTZ DEFAULT NULL,
    username VARCHAR(45) NOT NULL UNIQUE,
    email VARCHAR(45) NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    is_verified BOOLEAN NOT NULL DEFAULT FALSE
  );

-- Index for fast lookups by UUID
CREATE INDEX idx_user_credential_uuid ON user_credential (uuid);
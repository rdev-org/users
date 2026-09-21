CREATE TABLE
  IF NOT EXISTS pending_registration (
    id BIGSERIAL PRIMARY KEY,
    email VARCHAR(254) NOT NULL UNIQUE,
    token_hash CHAR(64) NOT NULL UNIQUE,
    expires_at TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW ()
  );

-- Fast lookup when user submits token in Create Profile handler
CREATE INDEX idx_pending_registrations_token_hash ON pending_registrations (token_hash);
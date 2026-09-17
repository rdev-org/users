CREATE TABLE
  IF NOT EXISTS user_preference (
    id BIGSERIAL PRIMARY KEY,
    status VARCHAR(1) NOT NULL DEFAULT 'A',
    uuid UUID NOT NULL DEFAULT gen_random_uuid (),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW (),
    updated_at TIMESTAMPTZ DEFAULT NOW (),
    deleted_at TIMESTAMPTZ DEFAULT NULL,
    user_id UUID NOT NULL,
    language VARCHAR(15),
    currency VARCHAR(15),
    email_opt_in BOOLEAN DEFAULT FALSE,
    sms_opt_in BOOLEAN DEFAULT FALSE
  );

-- Index for fast lookups by user_id (essential for performance)
CREATE INDEX idx_user_preference_user_id ON user_preference (user_id);

-- Index for fast lookups by UUID
CREATE INDEX idx_user_preference_uuid ON user_preference (uuid);

CREATE UNIQUE INDEX idx_user_preference_single_preference ON user_preference (user_id)
WHERE
  deleted_at IS NULL;
CREATE TABLE
  IF NOT EXISTS user_profile (
    id BIGSERIAL PRIMARY KEY,
    status VARCHAR(1) NOT NULL DEFAULT 'A',
    uuid UUID NOT NULL DEFAULT gen_random_uuid (),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW (),
    updated_at TIMESTAMPTZ DEFAULT NOW (),
    deleted_at TIMESTAMPTZ DEFAULT NULL,
    user_id UUID NOT NULL,
    first_name VARCHAR(45) NOT NULL,
    middle_name VARCHAR(45) NULL,
    last_name VARCHAR(45) NOT NULL,
    mobile_number VARCHAR(15) NULL,
    date_of_birth VARCHAR(10) NULL,
    gender VARCHAR(1) NULL,
    avatar VARCHAR(255) NULL
  );

-- Index for fast lookups by user_id (essential for performance)
CREATE INDEX idx_user_profile_user_id ON user_profile (user_id);

-- Index for fast lookups by UUID
CREATE INDEX idx_user_profile_uuid ON user_profile (uuid);

CREATE UNIQUE INDEX idx_user_profile_single_profile ON user_profile (user_id)
WHERE
  deleted_at IS NULL;
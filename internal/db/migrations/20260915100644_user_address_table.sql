CREATE TABLE
  IF NOT EXISTS user_address (
    id BIGSERIAL PRIMARY KEY,
    status VARCHAR(1) NOT NULL DEFAULT 'A',
    uuid UUID NOT NULL DEFAULT gen_random_uuid (),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW (),
    updated_at TIMESTAMPTZ DEFAULT NOW (),
    deleted_at TIMESTAMPTZ DEFAULT NULL,
    user_id UUID NOT NULL,
    address_line1 VARCHAR(255) NOT NULL,
    address_line2 VARCHAR(255),
    city VARCHAR(100) NOT NULL,
    state_province VARCHAR(100) NOT NULL,
    postal_code VARCHAR(20) NOT NULL,
    country VARCHAR(100) NOT NULL,
    is_default BOOLEAN NOT NULL DEFAULT FALSE,
    address_type VARCHAR(50) NOT NULL DEFAULT 'primary'
  );

-- Index for fast lookups by user_id (essential for performance)
CREATE INDEX idx_user_address_user_id ON user_address (user_id);

-- Index for fast lookups by uuid
CREATE INDEX idx_user_address_uuid ON user_address (uuid);

-- Ensures a user can only have ONE active default address at a time
CREATE UNIQUE INDEX idx_user_address_single_default ON user_address (user_id)
WHERE
  is_default = TRUE
  AND deleted_at IS NULL;
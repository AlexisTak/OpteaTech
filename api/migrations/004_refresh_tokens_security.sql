ALTER TABLE admin_refresh_tokens
  ADD COLUMN IF NOT EXISTS fingerprint_hash TEXT,
  ADD COLUMN IF NOT EXISTS user_agent TEXT,
  ADD COLUMN IF NOT EXISTS ip_address VARCHAR(45),
  ADD COLUMN IF NOT EXISTS revoked_at TIMESTAMPTZ,
  ADD COLUMN IF NOT EXISTS replaced_by_token TEXT;

CREATE INDEX IF NOT EXISTS idx_admin_refresh_tokens_revoked ON admin_refresh_tokens(revoked_at);
CREATE INDEX IF NOT EXISTS idx_admin_refresh_tokens_fingerprint ON admin_refresh_tokens(fingerprint_hash);

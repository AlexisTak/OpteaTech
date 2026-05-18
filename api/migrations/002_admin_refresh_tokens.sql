CREATE TABLE IF NOT EXISTS admin_refresh_tokens (
  token TEXT PRIMARY KEY,
  user_email VARCHAR(200) NOT NULL,
  expires_at TIMESTAMPTZ NOT NULL,
  created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_admin_refresh_tokens_expires ON admin_refresh_tokens(expires_at);
CREATE INDEX IF NOT EXISTS idx_admin_refresh_tokens_user_email ON admin_refresh_tokens(user_email);

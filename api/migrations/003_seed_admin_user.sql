INSERT INTO admin_users (email, password_hash, name, role)
VALUES (
  'admin@optea.tech',
  '$2a$12$LQv3c1yqBWVHxkd0LHAkCOYz6TtxMQJqhN8/LewY5GyYzS3MebAJu',
  'Admin',
  'admin'
)
ON CONFLICT (email) DO NOTHING;

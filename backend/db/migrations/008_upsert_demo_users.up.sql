-- Upsert usuários demo com hashes corretos
-- admin123, user123, demo123

INSERT INTO users (username, password_hash, email, full_name, is_active)
VALUES
  ('admin', '$2a$10$kR.pK7uclXtW7Qrt3UlLiONpGCukqRBkwOKLkR/iynitqqdwSUTdG', 'admin@example.com', 'Administrator', true),
  ('user',  '$2a$10$1TGCvNlUXWSQmVvDl/zZBO1qy.W6XRWi95gEtgZZ3qB45HIcgYHwS', 'user@example.com',  'Regular User',   true),
  ('demo',  '$2a$10$22F5d06lzO.LHTPQP4aTFu8PM7f6iQTMLdw/KwK7DKEGSciWzFBGG', 'demo@example.com',  'Demo User',      true)
ON CONFLICT (username) DO UPDATE SET
  password_hash = EXCLUDED.password_hash,
  email = EXCLUDED.email,
  full_name = EXCLUDED.full_name,
  is_active = EXCLUDED.is_active,
  updated_at = CURRENT_TIMESTAMP;



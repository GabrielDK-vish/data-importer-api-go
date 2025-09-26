-- Upsert usuários demo com hashes corretos
-- admin123, user123, demo123

INSERT INTO users (username, password_hash, email, full_name, is_active)
VALUES
  ('admin', '$2a$10$N9qo8uLOickgx2ZMRZoMye.IjdQvOQ5eqGStBUKx6XgKnrQvp.Fl6', 'admin@example.com', 'Administrator', true),
  ('user',  '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', 'user@example.com',  'Regular User',   true),
  ('demo',  '$2a$10$TKh8H1.PfQx37YgCzwiKb.KjNyWgaHb9cbcoQgdIVFlYg7B77UdFm', 'demo@example.com',  'Demo User',      true)
ON CONFLICT (username) DO UPDATE SET
  password_hash = EXCLUDED.password_hash,
  email = EXCLUDED.email,
  full_name = EXCLUDED.full_name,
  is_active = EXCLUDED.is_active,
  updated_at = CURRENT_TIMESTAMP;



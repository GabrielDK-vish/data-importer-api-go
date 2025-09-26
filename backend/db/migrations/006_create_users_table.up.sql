CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    email VARCHAR(255),
    full_name VARCHAR(255),
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_users_username ON users(username);
CREATE INDEX idx_users_email ON users(email);

-- Inserir usuários padrão (senhas são hash de bcrypt)
-- admin123 -> $2a$10$N9qo8uLOickgx2ZMRZoMye.IjdQvOQ5eqGStBUKx6XgKnrQvp.Fl6
-- user123 ->  $2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi
-- demo123 ->  $2a$10$TKh8H1.PfQx37YgCzwiKb.KjNyWgaHb9cbcoQgdIVFlYg7B77UdFm
INSERT INTO users (username, password_hash, email, full_name) VALUES
('admin', '$2a$10$N9qo8uLOickgx2ZMRZoMye.IjdQvOQ5eqGStBUKx6XgKnrQvp.Fl6', 'admin@example.com', 'Administrator'),
('user', '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', 'user@example.com', 'Regular User'),
('demo', '$2a$10$TKh8H1.PfQx37YgCzwiKb.KjNyWgaHb9cbcoQgdIVFlYg7B77UdFm', 'demo@example.com', 'Demo User')
ON CONFLICT (username) DO NOTHING;

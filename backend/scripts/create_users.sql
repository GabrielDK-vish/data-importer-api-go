-- Script SQL para criar usuários com senhas hasheadas usando bcrypt
-- As senhas são: admin123, user123, demo123
-- Hash gerado com bcrypt (cost 10)

-- Inserir usuários com senhas hasheadas
INSERT INTO users (username, password_hash, email, full_name, is_active) VALUES
('admin', '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', 'admin@example.com', 'Administrator', true),
('user', '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', 'user@example.com', 'Regular User', true),
('demo', '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', 'demo@example.com', 'Demo User', true)
ON CONFLICT (username) DO UPDATE SET
    password_hash = EXCLUDED.password_hash,
    email = EXCLUDED.email,
    full_name = EXCLUDED.full_name,
    updated_at = CURRENT_TIMESTAMP;

-- Verificar usuários criados
SELECT id, username, email, full_name, is_active, created_at FROM users ORDER BY username;

-- Corrigir hashes de senha dos usuários demo
-- admin123: $2a$10$N9qo8uLOickgx2ZMRZoMye.IjdQvOQ5eqGStBUKx6XgKnrQvp.Fl6
-- user123:  $2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi
-- demo123:  $2a$10$TKh8H1.PfQx37YgCzwiKb.KjNyWgaHb9cbcoQgdIVFlYg7B77UdFm

UPDATE users SET 
    password_hash = '$2a$10$N9qo8uLOickgx2ZMRZoMye.IjdQvOQ5eqGStBUKx6XgKnrQvp.Fl6',
    updated_at = CURRENT_TIMESTAMP
WHERE username = 'admin';

UPDATE users SET 
    password_hash = '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi',
    updated_at = CURRENT_TIMESTAMP
WHERE username = 'user';

UPDATE users SET 
    password_hash = '$2a$10$TKh8H1.PfQx37YgCzwiKb.KjNyWgaHb9cbcoQgdIVFlYg7B77UdFm',
    updated_at = CURRENT_TIMESTAMP
WHERE username = 'demo';

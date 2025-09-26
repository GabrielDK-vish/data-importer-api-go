-- Script para corrigir os hashes de senha dos usuários
-- Execute este script no banco de dados PostgreSQL

-- Atualizar hash do admin
UPDATE users SET 
    password_hash = '$2a$10$kR.pK7uclXtW7Qrt3UlLiONpGCukqRBkwOKLkR/iynitqqdwSUTdG',
    updated_at = CURRENT_TIMESTAMP
WHERE username = 'admin';

-- Atualizar hash do user
UPDATE users SET 
    password_hash = '$2a$10$1TGCvNlUXWSQmVvDl/zZBO1qy.W6XRWi95gEtgZZ3qB45HIcgYHwS',
    updated_at = CURRENT_TIMESTAMP
WHERE username = 'user';

-- Atualizar hash do demo
UPDATE users SET 
    password_hash = '$2a$10$22F5d06lzO.LHTPQP4aTFu8PM7f6iQTMLdw/KwK7DKEGSciWzFBGG',
    updated_at = CURRENT_TIMESTAMP
WHERE username = 'demo';

-- Verificar se as atualizações foram aplicadas
SELECT username, password_hash, updated_at FROM users ORDER BY username;

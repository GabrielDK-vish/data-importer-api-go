# Autenticação com PostgreSQL

## Migração de Autenticação Local para Banco de Dados

O sistema foi migrado de autenticação hardcoded para autenticação baseada em banco de dados PostgreSQL.

### Mudanças Implementadas

1. **Nova tabela `users`**:
   - Campos: `id`, `username`, `password_hash`, `email`, `full_name`, `is_active`, `created_at`, `updated_at`
   - Índices em `username` e `email`
   - Senhas armazenadas com hash bcrypt

2. **Migration 006**: `006_create_users_table.up.sql`
   - Cria a tabela `users`
   - Insere usuários padrão com senhas hasheadas

3. **Modelo `User`**: Adicionado em `internal/models/models.go`

4. **Repository**: Método `GetUserByUsername()` em `internal/repository/repository.go`

5. **Service**: Método `ValidateUserCredentials()` em `internal/service/service.go`
   - Valida credenciais usando bcrypt
   - Retorna dados do usuário se válido

6. **API**: Login atualizado em `api/routes.go`
   - Usa `service.ValidateUserCredentials()` em vez de `auth.ValidateCredentials()`

### Usuários Padrão

| Username | Password | Email | Full Name |
|----------|----------|-------|-----------|
| admin    | admin123 | admin@example.com | Administrator |
| user     | user123  | user@example.com | Regular User |
| demo     | demo123  | demo@example.com | Demo User |

### Como Executar

1. **Executar migrations**:
   ```bash
   # As migrations incluem a criação da tabela e inserção dos usuários
   go run ./cmd/main.go
   ```

2. **Ou executar script SQL manualmente**:
   ```bash
   psql -d data_importer -f backend/scripts/create_users.sql
   ```

3. **Ou usar o script Go**:
   ```bash
   go run ./cmd/create_users.go
   ```

### Segurança

- Senhas são hasheadas com bcrypt (cost 10)
- Campo `password_hash` não é exposto no JSON (tag `json:"-"`)
- Usuários inativos (`is_active = false`) não podem fazer login
- JWT continua sendo usado para sessões

### Compatibilidade

- A função `auth.ValidateCredentials()` foi marcada como DEPRECATED
- Mantida apenas para compatibilidade, mas não é mais usada
- Login agora sempre consulta o banco de dados

### Próximos Passos

1. **Adicionar endpoint para criar usuários** (admin only)
2. **Implementar reset de senha**
3. **Adicionar roles/permissões**
4. **Logs de auditoria de login**
5. **Rate limiting no login**

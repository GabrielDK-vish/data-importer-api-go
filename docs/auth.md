# Sistema de Autenticação

## Visão Geral

O sistema utiliza autenticação baseada em banco de dados PostgreSQL com JWT para sessões.

## Estrutura

### Tabela Users
```sql
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    email VARCHAR(255),
    full_name VARCHAR(255),
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

### Usuários Padrão

| Username | Password | Email | Full Name |
|----------|----------|-------|-----------|
| admin    | admin123 | admin@example.com | Administrator |
| user     | user123  | user@example.com | Regular User |
| demo     | demo123  | demo@example.com | Demo User |

## Implementação

### Validação de Credenciais
```go
func (s *Service) ValidateUserCredentials(ctx context.Context, username, password string) (*models.User, error) {
    user, err := s.repo.GetUserByUsername(ctx, username)
    if err != nil {
        return nil, fmt.Errorf("erro ao buscar usuário: %w", err)
    }
    
    if user == nil {
        return nil, fmt.Errorf("usuário não encontrado")
    }
    
    if err := s.comparePassword(password, user.PasswordHash); err != nil {
        return nil, fmt.Errorf("senha inválida")
    }
    
    return user, nil
}
```

### Geração de Token JWT
```go
func GenerateToken(username string) (string, error) {
    expirationTime := time.Now().Add(24 * time.Hour)
    
    claims := &Claims{
        Username: username,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(expirationTime),
            IssuedAt:  jwt.NewNumericDate(time.Now()),
            NotBefore: jwt.NewNumericDate(time.Now()),
        },
    }
    
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString(jwtSecret)
}
```

## Segurança

- Senhas hasheadas com bcrypt (cost 10)
- Campo password_hash não exposto no JSON
- Usuários inativos não podem fazer login
- Tokens JWT expiram em 24 horas
- Validação de entrada em todos os endpoints

## Migrations

### 006_create_users_table.up.sql
Cria tabela users e insere usuários padrão.

### 008_upsert_demo_users.up.sql
Atualiza usuários com hashes corretos.

### 009_update_password_hashes.up.sql
Corrige hashes de senha em produção.

## Uso

### Login via API
```bash
curl -X POST https://data-importer-api-go.onrender.com/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username": "admin", "password": "admin123"}'
```

### Usar Token
```bash
curl -X GET https://data-importer-api-go.onrender.com/api/customers \
  -H "Authorization: Bearer <token>"
```

## Configuração

### Variáveis de Ambiente
```bash
JWT_SECRET=sua-chave-secreta-super-segura-aqui
DATABASE_URL=postgres://user:pass@host:port/db?sslmode=disable
```

### Middleware de Autenticação
```go
func (h *Handler) AuthMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        authHeader := r.Header.Get("Authorization")
        if authHeader == "" {
            http.Error(w, "Token de autorização necessário", http.StatusUnauthorized)
            return
        }
        
        tokenString := authHeader
        if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
            tokenString = authHeader[7:]
        }
        
        claims, err := auth.ValidateToken(tokenString)
        if err != nil {
            http.Error(w, "Token inválido", http.StatusUnauthorized)
            return
        }
        
        ctx := context.WithValue(r.Context(), "username", claims.Username)
        next.ServeHTTP(w, r.WithContext(ctx))
    })
}
```
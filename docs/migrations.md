# Guia de Migrations

## Visão Geral

Migrations SQL para gerenciar estrutura do banco PostgreSQL usando golang-migrate.

## Estrutura

```
backend/db/migrations/
├── 001_create_partners_table.up.sql
├── 001_create_partners_table.down.sql
├── 002_create_customers_table.up.sql
├── 002_create_customers_table.down.sql
├── 003_create_products_table.up.sql
├── 003_create_products_table.down.sql
├── 004_create_usages_table.up.sql
├── 004_create_usages_table.down.sql
├── 005_insert_test_data.up.sql
├── 005_insert_test_data.down.sql
├── 006_create_users_table.up.sql
├── 006_create_users_table.down.sql
├── 007_fix_user_passwords.up.sql
├── 008_upsert_demo_users.up.sql
├── 008_upsert_demo_users.down.sql
├── 009_update_password_hashes.up.sql
└── 009_update_password_hashes.down.sql
```

## Tabelas

### Partners
```sql
CREATE TABLE partners (
    id SERIAL PRIMARY KEY,
    partner_id VARCHAR(255) UNIQUE NOT NULL,
    partner_name VARCHAR(255) NOT NULL,
    mpn_id VARCHAR(255),
    tier2_mpn_id VARCHAR(255),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

### Customers
```sql
CREATE TABLE customers (
    id SERIAL PRIMARY KEY,
    customer_id VARCHAR(255) UNIQUE NOT NULL,
    customer_name VARCHAR(255) NOT NULL,
    customer_domain_name VARCHAR(255),
    country VARCHAR(100),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

### Products
```sql
CREATE TABLE products (
    id SERIAL PRIMARY KEY,
    product_id VARCHAR(255) UNIQUE NOT NULL,
    sku_id VARCHAR(255) NOT NULL,
    sku_name VARCHAR(255) NOT NULL,
    product_name VARCHAR(255) NOT NULL,
    meter_type VARCHAR(100),
    category VARCHAR(100),
    sub_category VARCHAR(100),
    unit_type VARCHAR(50),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

### Usages
```sql
CREATE TABLE usages (
    id SERIAL PRIMARY KEY,
    invoice_number VARCHAR(255),
    charge_start_date DATE,
    usage_date DATE NOT NULL,
    quantity DECIMAL(15,6) NOT NULL,
    unit_price DECIMAL(15,6) NOT NULL,
    billing_pre_tax_total DECIMAL(15,2) NOT NULL,
    resource_location VARCHAR(255),
    tags TEXT,
    benefit_type VARCHAR(100),
    partner_id INTEGER REFERENCES partners(id) ON DELETE CASCADE,
    customer_id INTEGER REFERENCES customers(id) ON DELETE CASCADE,
    product_id INTEGER REFERENCES products(id) ON DELETE CASCADE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

### Users
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

## Relacionamentos

### Chaves Estrangeiras
- usages.partner_id → partners.id
- usages.customer_id → customers.id
- usages.product_id → products.id

### Índices
```sql
CREATE INDEX idx_partners_partner_id ON partners(partner_id);
CREATE INDEX idx_customers_customer_id ON customers(customer_id);
CREATE INDEX idx_products_product_id ON products(product_id);
CREATE INDEX idx_usages_usage_date ON usages(usage_date);
CREATE INDEX idx_users_username ON users(username);
```

## Execução

### Automática (Recomendado)
Migrations executam automaticamente na inicialização da API.

### Manual Local
```bash
# Instalar golang-migrate
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# Executar migrations
migrate -path ./db/migrations -database "postgres://postgres:password@localhost:5432/data_importer?sslmode=disable" up

# Reverter migrations
migrate -path ./db/migrations -database "postgres://postgres:password@localhost:5432/data_importer?sslmode=disable" down
```

### Via Docker
```bash
# Executar migrations
docker-compose exec api migrate -path /app/db/migrations -database "postgres://postgres:password@postgres:5432/data_importer?sslmode=disable" up

# Reverter migrations
docker-compose exec api migrate -path /app/db/migrations -database "postgres://postgres:password@postgres:5432/data_importer?sslmode=disable" down
```

## Verificação

### Conectar ao Banco
```bash
# Via Docker
docker-compose exec postgres psql -U postgres -d data_importer

# Localmente
psql -U postgres -d data_importer
```

### Verificar Tabelas
```sql
-- Listar tabelas
\dt

-- Verificar estrutura
\d partners
\d customers
\d products
\d usages
\d users
```

### Verificar Dados
```sql
-- Contar registros
SELECT 'partners' as tabela, COUNT(*) as registros FROM partners
UNION ALL
SELECT 'customers', COUNT(*) FROM customers
UNION ALL
SELECT 'products', COUNT(*) FROM products
UNION ALL
SELECT 'usages', COUNT(*) FROM usages
UNION ALL
SELECT 'users', COUNT(*) FROM users;
```

## Dados de Teste

### Usuários Padrão
```sql
INSERT INTO users (username, password_hash, email, full_name, is_active) VALUES
('admin', '$2a$10$kR.pK7uclXtW7Qrt3UlLiONpGCukqRBkwOKLkR/iynitqqdwSUTdG', 'admin@example.com', 'Administrator', true),
('user', '$2a$10$1TGCvNlUXWSQmVvDl/zZBO1qy.W6XRWi95gEtgZZ3qB45HIcgYHwS', 'user@example.com', 'Regular User', true),
('demo', '$2a$10$22F5d06lzO.LHTPQP4aTFu8PM7f6iQTMLdw/KwK7DKEGSciWzFBGG', 'demo@example.com', 'Demo User', true)
ON CONFLICT (username) DO UPDATE SET
  password_hash = EXCLUDED.password_hash,
  email = EXCLUDED.email,
  full_name = EXCLUDED.full_name,
  is_active = EXCLUDED.is_active,
  updated_at = CURRENT_TIMESTAMP;
```

## Troubleshooting

### Erro: "relation does not exist"
Migrations não foram executadas.
```bash
# Verificar logs
docker-compose logs api | grep "Migrations executadas"

# Executar manualmente
docker-compose exec api migrate -path /app/db/migrations -database "postgres://postgres:password@postgres:5432/data_importer?sslmode=disable" up
```

### Erro: "duplicate key value violates unique constraint"
Dados duplicados. Usar ON CONFLICT DO NOTHING ou ON CONFLICT DO UPDATE.

### Erro: "foreign key constraint fails"
Referência inexistente. Verificar se dados de referência existem.

### Erro: "migrate: command not found"
```bash
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
export PATH=$PATH:$(go env GOPATH)/bin
```

## Histórico de Migrations

### 001-004: Tabelas Principais
Criação das tabelas partners, customers, products, usages.

### 005: Dados de Teste
Inserção de dados de exemplo para desenvolvimento.

### 006: Tabela Users
Criação da tabela de usuários para autenticação.

### 007: Correção de Senhas
Correção inicial dos hashes de senha.

### 008: Upsert de Usuários
Atualização idempotente dos usuários demo.

### 009: Atualização de Hashes
Correção final dos hashes de senha para produção.
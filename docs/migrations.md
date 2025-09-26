# Guia de Migrations

## Visão Geral

As migrations são scripts SQL que criam e modificam a estrutura do banco de dados PostgreSQL. Utilizamos a biblioteca `golang-migrate` para gerenciar as migrations automaticamente.

## Estrutura das Migrations

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
└── 005_insert_test_data.down.sql
```

## Tabelas Criadas

### 1. Partners (Parceiros)

```sql
CREATE TABLE IF NOT EXISTS partners (
    id SERIAL PRIMARY KEY,
    partner_id VARCHAR(255) UNIQUE NOT NULL,
    partner_name VARCHAR(255) NOT NULL,
    mpn_id VARCHAR(255),
    tier2_mpn_id VARCHAR(255),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_partners_partner_id ON partners(partner_id);
```

**Campos:**
- `id` - Chave primária auto-incremento
- `partner_id` - ID único do parceiro
- `partner_name` - Nome do parceiro
- `mpn_id` - ID do MPN (opcional)
- `tier2_mpn_id` - ID do Tier2 MPN (opcional)
- `created_at` - Data de criação
- `updated_at` - Data de atualização

### 2. Customers (Clientes)

```sql
CREATE TABLE IF NOT EXISTS customers (
    id SERIAL PRIMARY KEY,
    customer_id VARCHAR(255) UNIQUE NOT NULL,
    customer_name VARCHAR(255) NOT NULL,
    customer_domain_name VARCHAR(255),
    country VARCHAR(100),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_customers_customer_id ON customers(customer_id);
CREATE INDEX idx_customers_country ON customers(country);
```

**Campos:**
- `id` - Chave primária auto-incremento
- `customer_id` - ID único do cliente
- `customer_name` - Nome do cliente
- `customer_domain_name` - Domínio do cliente (opcional)
- `country` - País do cliente
- `created_at` - Data de criação
- `updated_at` - Data de atualização

### 3. Products (Produtos)

```sql
CREATE TABLE IF NOT EXISTS products (
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

CREATE INDEX idx_products_product_id ON products(product_id);
CREATE INDEX idx_products_sku_id ON products(sku_id);
CREATE INDEX idx_products_category ON products(category);
```

**Campos:**
- `id` - Chave primária auto-incremento
- `product_id` - ID único do produto
- `sku_id` - ID do SKU
- `sku_name` - Nome do SKU
- `product_name` - Nome do produto
- `meter_type` - Tipo de medição (opcional)
- `category` - Categoria do produto
- `sub_category` - Subcategoria (opcional)
- `unit_type` - Tipo de unidade (opcional)
- `created_at` - Data de criação
- `updated_at` - Data de atualização

### 4. Usages (Usos)

```sql
CREATE TABLE IF NOT EXISTS usages (
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

CREATE INDEX idx_usages_partner_id ON usages(partner_id);
CREATE INDEX idx_usages_customer_id ON usages(customer_id);
CREATE INDEX idx_usages_product_id ON usages(product_id);
CREATE INDEX idx_usages_usage_date ON usages(usage_date);
CREATE INDEX idx_usages_charge_start_date ON usages(charge_start_date);
CREATE INDEX idx_usages_invoice_number ON usages(invoice_number);
```

**Campos:**
- `id` - Chave primária auto-incremento
- `invoice_number` - Número da fatura (opcional)
- `charge_start_date` - Data de início da cobrança (opcional)
- `usage_date` - Data do uso (obrigatório)
- `quantity` - Quantidade utilizada
- `unit_price` - Preço unitário
- `billing_pre_tax_total` - Total antes dos impostos
- `resource_location` - Localização do recurso (opcional)
- `tags` - Tags adicionais (opcional)
- `benefit_type` - Tipo de benefício (opcional)
- `partner_id` - Referência ao parceiro
- `customer_id` - Referência ao cliente
- `product_id` - Referência ao produto
- `created_at` - Data de criação
- `updated_at` - Data de atualização

## Relacionamentos

### Chaves Estrangeiras

```sql
-- Usages → Partners
partner_id INTEGER REFERENCES partners(id) ON DELETE CASCADE

-- Usages → Customers  
customer_id INTEGER REFERENCES customers(id) ON DELETE CASCADE

-- Usages → Products
product_id INTEGER REFERENCES products(id) ON DELETE CASCADE
```

### Índices para Performance

```sql
-- Índices para consultas frequentes
CREATE INDEX idx_usages_usage_date ON usages(usage_date);
CREATE INDEX idx_usages_partner_id ON usages(partner_id);
CREATE INDEX idx_usages_customer_id ON usages(customer_id);
CREATE INDEX idx_usages_product_id ON usages(product_id);
```

## Execução das Migrations

### 1. Automática (Recomendado)

As migrations são executadas automaticamente quando a API inicia:

```go
// Em cmd/main.go
func runMigrations(databaseURL string) error {
    m, err := migrate.New(
        "file://db/migrations",
        databaseURL,
    )
    if err != nil {
        return fmt.Errorf("erro ao criar migrator: %w", err)
    }
    defer m.Close()

    if err := m.Up(); err != nil && err != migrate.ErrNoChange {
        return fmt.Errorf("erro ao executar migrations: %w", err)
    }

    log.Println("Migrations executadas com sucesso")
    return nil
}
```

### 2. Manual com Docker

```bash
# Executar migrations manualmente
docker-compose exec api migrate -path /app/db/migrations -database "postgres://postgres:password@postgres:5432/data_importer?sslmode=disable" up

# Reverter migrations
docker-compose exec api migrate -path /app/db/migrations -database "postgres://postgres:password@postgres:5432/data_importer?sslmode=disable" down
```

### 3. Localmente

```bash
# Instalar golang-migrate
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# Executar migrations
migrate -path ./db/migrations -database "postgres://postgres:password@localhost:5432/data_importer?sslmode=disable" up

# Reverter migrations
migrate -path ./db/migrations -database "postgres://postgres:password@localhost:5432/data_importer?sslmode=disable" down
```

## Dados de Teste

### Migration 005: Insert Test Data

```sql
-- Inserir dados de teste para partners
INSERT INTO partners (partner_id, partner_name, mpn_id, tier2_mpn_id) VALUES
('PARTNER001', 'Microsoft Corporation', 'MPN001', 'T2MPN001'),
('PARTNER002', 'Amazon Web Services', 'MPN002', 'T2MPN002'),
('PARTNER003', 'Google Cloud Platform', 'MPN003', 'T2MPN003')
ON CONFLICT (partner_id) DO NOTHING;

-- Inserir dados de teste para customers
INSERT INTO customers (customer_id, customer_name, customer_domain_name, country) VALUES
('CUST001', 'TechCorp Solutions', 'techcorp.com', 'Brazil'),
('CUST002', 'DataFlow Inc', 'dataflow.com', 'United States'),
('CUST003', 'CloudTech Ltd', 'cloudtech.co.uk', 'United Kingdom'),
('CUST004', 'InnovateSoft', 'innovatesoft.com', 'Canada')
ON CONFLICT (customer_id) DO NOTHING;

-- Inserir dados de teste para products
INSERT INTO products (product_id, sku_id, sku_name, product_name, meter_type, category, sub_category, unit_type) VALUES
('PROD001', 'SKU001', 'Azure Virtual Machine', 'Azure Compute', 'Compute', 'Virtual Machines', 'Standard', 'Hours'),
('PROD002', 'SKU002', 'AWS EC2 Instance', 'AWS Compute', 'Compute', 'EC2', 'General Purpose', 'Hours'),
('PROD003', 'SKU003', 'Google Cloud Storage', 'GCP Storage', 'Storage', 'Cloud Storage', 'Standard', 'GB-Month'),
('PROD004', 'SKU004', 'Azure SQL Database', 'Azure Database', 'Database', 'SQL Database', 'Standard', 'DTU-Hours')
ON CONFLICT (product_id) DO NOTHING;
```

## Verificação das Migrations

### 1. Conectar ao Banco

```bash
# Via Docker
docker-compose exec postgres psql -U postgres -d data_importer

# Localmente
psql -U postgres -d data_importer
```

### 2. Verificar Tabelas

```sql
-- Listar todas as tabelas
\dt

-- Verificar estrutura de uma tabela
\d partners
\d customers
\d products
\d usages
```

### 3. Verificar Dados

```sql
-- Contar registros em cada tabela
SELECT 'partners' as tabela, COUNT(*) as registros FROM partners
UNION ALL
SELECT 'customers', COUNT(*) FROM customers
UNION ALL
SELECT 'products', COUNT(*) FROM products
UNION ALL
SELECT 'usages', COUNT(*) FROM usages;
```

### 4. Verificar Relacionamentos

```sql
-- Verificar integridade referencial
SELECT 
    p.partner_name,
    COUNT(u.id) as total_usages
FROM partners p
LEFT JOIN usages u ON p.id = u.partner_id
GROUP BY p.id, p.partner_name;
```

## Troubleshooting

### Erro: "relation does not exist"

**Causa**: Migrations não foram executadas.

**Solução**: 
```bash
# Verificar se migrations foram executadas
docker-compose logs api | grep "Migrations executadas"

# Executar migrations manualmente
docker-compose exec api migrate -path /app/db/migrations -database "postgres://postgres:password@postgres:5432/data_importer?sslmode=disable" up
```

### Erro: "duplicate key value violates unique constraint"

**Causa**: Tentativa de inserir dados duplicados.

**Solução**: Usar `ON CONFLICT DO NOTHING` ou `ON CONFLICT DO UPDATE`.

### Erro: "foreign key constraint fails"

**Causa**: Tentativa de inserir uso com partner/customer/product inexistente.

**Solução**: Verificar se os dados de referência existem antes de inserir usages.

## Próximos Passos

- [ ] Adicionar migrations para índices adicionais
- [ ] Implementar versionamento de schema
- [ ] Adicionar validações de integridade
- [ ] Criar migrations para dados de produção
- [ ] Implementar rollback automático

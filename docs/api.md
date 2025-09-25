# 📚 Documentação da API

## Visão Geral

A API Data Importer é uma API REST desenvolvida em Golang que permite importar dados de faturamento e fornece endpoints para consulta e análise desses dados.

## Base URL

```
http://localhost:8080
```

## Autenticação

A API utiliza JWT (JSON Web Token) para autenticação. Para acessar endpoints protegidos, inclua o token no header:

```
Authorization: Bearer <seu_token>
```

## Endpoints

### 1. Autenticação

#### POST /auth/login

Autentica um usuário e retorna um token JWT.

**Request:**
```json
{
  "username": "admin",
  "password": "admin123"
}
```

**Response (200):**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": "admin"
}
```

**Response (401):**
```json
{
  "error": "Credenciais inválidas"
}
```

**Credenciais de Teste:**
- `admin` / `admin123`
- `user` / `user123`
- `demo` / `demo123`

### 2. Clientes

#### GET /api/customers

Retorna lista de todos os clientes.

**Headers:**
```
Authorization: Bearer <token>
```

**Response (200):**
```json
[
  {
    "id": 1,
    "customer_id": "CUST001",
    "customer_name": "TechCorp Solutions",
    "customer_domain_name": "techcorp.com",
    "country": "Brazil",
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:00Z"
  }
]
```

#### GET /api/customers/{id}/usage

Retorna o histórico de uso de um cliente específico.

**Headers:**
```
Authorization: Bearer <token>
```

**Path Parameters:**
- `id` (integer): ID do cliente

**Response (200):**
```json
[
  {
    "id": 1,
    "invoice_number": "INV001",
    "charge_start_date": "2024-01-01T00:00:00Z",
    "usage_date": "2024-01-15T00:00:00Z",
    "quantity": 100.0,
    "unit_price": 0.05,
    "billing_pre_tax_total": 5.0,
    "resource_location": "East US",
    "tags": "env:prod",
    "benefit_type": "None",
    "partner_id": 1,
    "customer_id": 1,
    "product_id": 1,
    "partner": {
      "partner_id": "PARTNER001",
      "partner_name": "Microsoft Corporation"
    },
    "customer": {
      "customer_id": "CUST001",
      "customer_name": "TechCorp Solutions"
    },
    "product": {
      "product_id": "PROD001",
      "product_name": "Azure Virtual Machine",
      "category": "Compute"
    }
  }
]
```

### 3. Relatórios

#### GET /api/reports/billing/monthly

Retorna faturamento agrupado por mês.

**Headers:**
```
Authorization: Bearer <token>
```

**Response (200):**
```json
[
  {
    "month": "2024-01",
    "total": 150.0,
    "count": 5
  },
  {
    "month": "2024-02",
    "total": 200.0,
    "count": 3
  }
]
```

#### GET /api/reports/billing/by-product

Retorna faturamento agrupado por produto.

**Headers:**
```
Authorization: Bearer <token>
```

**Response (200):**
```json
[
  {
    "product_id": "PROD001",
    "product_name": "Azure Virtual Machine",
    "category": "Compute",
    "total": 75.0,
    "count": 3
  },
  {
    "product_id": "PROD002",
    "product_name": "AWS EC2 Instance",
    "category": "Compute",
    "total": 50.0,
    "count": 2
  }
]
```

#### GET /api/reports/billing/by-partner

Retorna faturamento agrupado por parceiro.

**Headers:**
```
Authorization: Bearer <token>
```

**Response (200):**
```json
[
  {
    "partner_id": "PARTNER001",
    "partner_name": "Microsoft Corporation",
    "total": 100.0,
    "count": 4
  },
  {
    "partner_id": "PARTNER002",
    "partner_name": "Amazon Web Services",
    "total": 80.0,
    "count": 3
  }
]
```

## Códigos de Status HTTP

- `200 OK` - Requisição bem-sucedida
- `400 Bad Request` - Dados inválidos na requisição
- `401 Unauthorized` - Token inválido ou ausente
- `500 Internal Server Error` - Erro interno do servidor

## Exemplos de Uso

### 1. Autenticação e Listagem de Clientes

```bash
# 1. Fazer login
curl -X POST http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username": "admin", "password": "admin123"}'

# 2. Usar o token retornado para listar clientes
curl -X GET http://localhost:8080/api/customers \
  -H "Authorization: Bearer <seu_token>"
```

### 2. Consultar Uso de um Cliente

```bash
curl -X GET http://localhost:8080/api/customers/1/usage \
  -H "Authorization: Bearer <seu_token>"
```

### 3. Obter Relatórios

```bash
# Faturamento mensal
curl -X GET http://localhost:8080/api/reports/billing/monthly \
  -H "Authorization: Bearer <seu_token>"

# Faturamento por produto
curl -X GET http://localhost:8080/api/reports/billing/by-product \
  -H "Authorization: Bearer <seu_token>"

# Faturamento por parceiro
curl -X GET http://localhost:8080/api/reports/billing/by-partner \
  -H "Authorization: Bearer <seu_token>"
```

## Estrutura do Banco de Dados

### Tabelas Principais

1. **partners** - Dados dos parceiros
2. **customers** - Informações dos clientes
3. **products** - Catálogo de produtos/serviços
4. **usages** - Registros de uso e faturamento

### Relacionamentos

- `usages.partner_id` → `partners.id`
- `usages.customer_id` → `customers.id`
- `usages.product_id` → `products.id`

## Performance

- **Conexões pool** do PostgreSQL para melhor performance
- **Índices otimizados** nas colunas de busca
- **Queries agregadas** para relatórios eficientes
- **Middleware de CORS** configurado
- **Graceful shutdown** implementado

## Segurança

- **Autenticação JWT** obrigatória para endpoints protegidos
- **Validação de entrada** em todos os endpoints
- **Tratamento de erros** padronizado
- **Headers de segurança** configurados

## Upload de Arquivos

### POST /api/upload

Faz upload e processa arquivos CSV ou Excel para importação de dados.

**Headers:**
```
Content-Type: multipart/form-data
Authorization: Bearer <token>
```

**Parâmetros:**
- `file` (arquivo): Arquivo CSV ou Excel (.xlsx) para importação

**Resposta de Sucesso (200):**
```json
{
  "success": true,
  "message": "Arquivo processado com sucesso",
  "data": {
    "partners": 5,
    "customers": 10,
    "products": 8,
    "usages": 150
  }
}
```

**Resposta de Erro (400):**
```json
{
  "error": "Tipo de arquivo não suportado. Use .csv ou .xlsx"
}
```

**Formatos Suportados:**
- CSV (.csv) - Arquivo separado por vírgulas
- Excel (.xlsx) - Planilha do Microsoft Excel

**Colunas Obrigatórias:**
- `partner_id` - ID do parceiro
- `customer_id` - ID do cliente  
- `product_id` - ID do produto
- `usage_date` - Data do uso
- `quantity` - Quantidade
- `unit_price` - Preço unitário

**Exemplo de Upload com cURL:**
```bash
curl -X POST http://localhost:8080/api/upload \
  -H "Authorization: Bearer <seu_token>" \
  -F "file=@dados.csv"
```

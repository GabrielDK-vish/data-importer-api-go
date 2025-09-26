# API Documentation

## Base URL
```
https://data-importer-api-go.onrender.com
```

## Autenticação
A API utiliza JWT. Inclua o token no header:
```
Authorization: Bearer <token>
```

## Endpoints

### Autenticação

#### POST /auth/login
Autentica usuário e retorna token JWT.

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

**Credenciais:**
- admin / admin123
- user / user123
- demo / demo123

### Clientes

#### GET /api/customers
Lista todos os clientes.

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
Histórico de uso de um cliente.

**Response (200):**
```json
[
  {
    "id": 1,
    "invoice_number": "INV001",
    "usage_date": "2024-01-15T00:00:00Z",
    "quantity": 100.0,
    "unit_price": 0.05,
    "billing_pre_tax_total": 5.0,
    "partner": {
      "partner_name": "Microsoft Corporation"
    },
    "product": {
      "product_name": "Azure Virtual Machine",
      "category": "Compute"
    }
  }
]
```

### Relatórios

#### GET /api/reports/billing/monthly
Faturamento agrupado por mês.

**Response (200):**
```json
[
  {
    "month": "2024-01",
    "total": 150.0,
    "count": 5
  }
]
```

#### GET /api/reports/billing/by-product
Faturamento agrupado por produto.

**Response (200):**
```json
[
  {
    "product_id": "PROD001",
    "product_name": "Azure Virtual Machine",
    "category": "Compute",
    "total": 75.0,
    "count": 3
  }
]
```

#### GET /api/reports/billing/by-partner
Faturamento agrupado por parceiro.

**Response (200):**
```json
[
  {
    "partner_id": "PARTNER001",
    "partner_name": "Microsoft Corporation",
    "total": 100.0,
    "count": 4
  }
]
```

### Upload

#### POST /api/upload
Upload e processamento de arquivos CSV/Excel.

**Headers:**
```
Content-Type: multipart/form-data
Authorization: Bearer <token>
```

**Parâmetros:**
- file: Arquivo CSV ou Excel (.xlsx)

**Response (200):**
```json
{
  "success": true,
  "message": "Arquivo processado e dados substituídos com sucesso",
  "data": {
    "partners": 5,
    "customers": 10,
    "products": 8,
    "usages": 150
  }
}
```

**Colunas Obrigatórias:**
- partner_id
- customer_id
- product_id
- usage_date
- quantity
- unit_price

## Códigos de Status

- 200 OK - Sucesso
- 400 Bad Request - Dados inválidos
- 401 Unauthorized - Token inválido
- 500 Internal Server Error - Erro interno

## Exemplos de Uso

### Autenticação
```bash
curl -X POST https://data-importer-api-go.onrender.com/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username": "admin", "password": "admin123"}'
```

### Listar Clientes
```bash
curl -X GET https://data-importer-api-go.onrender.com/api/customers \
  -H "Authorization: Bearer <token>"
```

### Upload de Arquivo
```bash
curl -X POST https://data-importer-api-go.onrender.com/api/upload \
  -H "Authorization: Bearer <token>" \
  -F "file=@dados.xlsx"
```
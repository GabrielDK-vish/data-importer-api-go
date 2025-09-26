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
Retorna faturamento mensal.

**Response (200):**
```json
[
  {
    "month": "2023-01",
    "total": 12500.75,
    "records": 150
  },
  {
    "month": "2023-02",
    "total": 13750.25,
    "records": 175
  }
]
```

#### GET /api/reports/billing/by-product
Retorna faturamento por produto.

**Response (200):**
```json
[
  {
    "product_id": "prod1",
    "product_name": "Product One",
    "category": "Category A",
    "total": 8500.50,
    "count": 120
  },
  {
    "product_id": "prod2",
    "product_name": "Product Two",
    "category": "Category B",
    "total": 6750.25,
    "count": 85
  }
]
```

#### GET /api/reports/billing/by-partner
Retorna faturamento por parceiro.

**Response (200):**
```json
[
  {
    "partner_id": "partner1",
    "partner_name": "Partner One",
    "total": 9500.75,
    "count": 130
  },
  {
    "partner_id": "partner2",
    "partner_name": "Partner Two",
    "total": 7750.50,
    "count": 95
  }
]
```

#### GET /api/reports/billing/by-category
Retorna faturamento por categoria.

**Response (200):**
```json
[
  {
    "category": "Category A",
    "total": 8500.50,
    "count": 120
  },
  {
    "category": "Category B",
    "total": 6750.25,
    "count": 85
  }
]
```

#### GET /api/reports/billing/by-resource
Retorna faturamento por recurso.

**Response (200):**
```json
[
  {
    "resource": "Resource A",
    "total": 7500.50,
    "count": 110
  },
  {
    "resource": "Resource B",
    "total": 5750.25,
    "count": 75
  }
]
```

#### GET /api/reports/billing/by-customer
Retorna faturamento por cliente.

**Response (200):**
```json
[
  {
    "customer_id": "customer1",
    "customer_name": "Customer One",
    "total": 9500.75,
    "count": 130
  },
  {
    "customer_id": "customer2",
    "customer_name": "Customer Two",
    "total": 7750.50,
    "count": 95
  }
]
```

#### GET /api/reports/kpi
Retorna indicadores de performance (KPIs).

**Response (200):**
```json
{
  "total_records": 325,
  "total_categories": 5,
  "total_resources": 8,
  "total_customers": 12,
  "avg_billing_per_month": 13125.50,
  "processing_time_ms": 1250,
  "last_updated": "2023-09-28T14:30:00Z"
}
```

### Upload

#### POST /api/upload
Upload de arquivo Excel/CSV.

**Request:**
- Multipart form com campo `file`

**Response (200):**
```json
{
  "status": "success",
  "records_processed": 325,
  "processing_time_ms": 1250
}
```

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
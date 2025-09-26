# Desafio Técnico Full stack (Golang)

## Desafio

criar um importador para uma base de dados (postgres) que deverá ser feito na linguagem golang para armazenar os dados do arquivo enviado. Será avaliado a normalização dos dados na base de dados, e a performance do importador.

criar uma API com a linguagem golang contendo endpoint de autenticação e outros endpoints de consultas dos dados fornecidos e que estão importados na base de dados Postgres.

Será um diferencial criar um front em React para esse projeto, mostrando indicadores totalizadores / agrupamentos de categorias / recursos / clientes / Meses de cobrança.

Necessário publicar em algum link para avaliação, com a documentação de execução da aplicação.

## Solução Proposta

Sistema de importação e análise de dados de faturamento desenvolvido em Golang com frontend React. A solução implementa:

1. Importador de alta performance para arquivos Excel/CSV
2. API RESTful com autenticação JWT
3. Banco de dados PostgreSQL com normalização completa
4. Frontend React com dashboard interativo
5. Indicadores de performance (KPIs) carregados dinamicamente do banco de dados
6. Visualizações por categoria, recurso, cliente e período

## Arquitetura

### Backend (Golang)
- API REST com framework Chi
- Autenticação JWT
- Banco PostgreSQL com normalização completa
- Processamento de arquivos Excel/CSV
- Relatórios agregados por mês, produto e parceiro

### Frontend (React)
- Dashboard com gráficos interativos e KPIs dinâmicos
- Visualizações por categoria, recurso e cliente
- Relatórios de faturamento
- Lista de clientes e detalhamento de uso

### Banco de Dados
Estrutura normalizada em 4 entidades:
- **partners**: dados de parceiros
- **customers**: dados de clientes  
- **products**: catálogo de produtos/serviços
- **usages**: registros de uso e faturamento

## Execução

### Pré-requisitos
- Docker e Docker Compose
- Go 1.21+ (desenvolvimento)
- Node.js 16+ (desenvolvimento)

### Execução com Docker
```bash
# Clonar repositório
git clone <repository-url>
cd data-importer-api-go

# Executar todos os serviços
docker-compose up --build

# Acessar:
# Frontend: http://localhost:3000
# API: http://localhost:8080
# PostgreSQL: localhost:5432
```

### Execução Local
```bash
# Backend
cd backend
go mod tidy
go run ./cmd/main.go

# Frontend
cd frontend
npm install
npm start
```

## Credenciais de Teste

| Usuário | Senha    |
|---------|----------|
| admin   | admin123 |
| user    | user123  |
| demo    | demo123  |

## Endpoints da API

### Públicos
- `GET /` - Página inicial
- `GET /health` - Status da aplicação
- `POST /auth/login` - Autenticação

### Protegidos (Requer JWT)
- `GET /api/customers` - Listar clientes
- `GET /api/customers/{id}/usage` - Uso por cliente
- `GET /api/reports/billing/monthly` - Faturamento mensal
- `GET /api/reports/billing/by-product` - Faturamento por produto
- `GET /api/reports/billing/by-partner` - Faturamento por parceiro
- `GET /api/reports/billing/by-category` - Faturamento por categoria
- `GET /api/reports/billing/by-resource` - Faturamento por recurso
- `GET /api/reports/billing/by-customer` - Faturamento por cliente
- `GET /api/reports/kpi` - Indicadores de performance

## Estrutura de Dados

### Entidades Principais
- **Customers**: Informações dos clientes
- **Products**: Catálogo de produtos com categorias e tipos de recursos
- **Usages**: Registros de uso e faturamento
- **Partners**: Dados dos parceiros

### Indicadores de Performance (KPIs)
- Total de registros processados
- Total de categorias
- Total de recursos
- Total de clientes
- Média de faturamento mensal
- Tempo de processamento

### Mapeamento Automático
O sistema possui mapeamento automático inteligente que reconhece variações dos nomes das colunas:
- **Partner ID**: `PartnerId`, `Partner_ID`, `partner-id`, `partner id`
- **Customer ID**: `CustomerId`, `Customer_ID`, `customer-id`, `customer id`
- **Product ID**: `ProductId`, `Product_ID`, `product-id`, `product id`
- **Usage Date**: `UsageDate`, `Usage_Date`, `usage-date`, `usage date`, `Date`
- **Quantity**: `Quantity`, `Qty`, `quantity`, `qty`
- **Unit Price**: `UnitPrice`, `Unit_Price`, `unit-price`, `unit price`, `Price`

### Tratamento Robusto de Erros
- **Dados inválidos**: Sistema usa valores padrão e continua processamento
- **Datas inválidas**: Usa data atual como fallback
- **Números inválidos**: Usa 0 como fallback
- **Linhas vazias**: Ignoradas automaticamente
- **Processamento paralelo**: Melhora performance e tolerância a erros

### Comportamento
- Upload substitui completamente os dados existentes
- Processo atômico (tudo ou nada)
- Carregamento automático na inicialização se banco vazio
- Logs detalhados para debug e monitoramento

## Deploy em Produção

### URLs
- **Frontend**: https://data-importer-api-go.vercel.app/
- **Backend**: https://data-importer-api-go.onrender.com/
- **Health Check**: https://data-importer-api-go.onrender.com/health

### Plataformas
- **Frontend**: Vercel (deploy automático)
- **Backend**: Render (deploy automático)
- **Banco**: PostgreSQL no Render

## Estrutura do Projeto

```
data-importer-api-go/
├── backend/                    # API Golang
│   ├── cmd/                   # Aplicações
│   ├── internal/              # Código interno
│   ├── api/                   # Handlers HTTP
│   ├── db/migrations/         # Migrations SQL
│   └── go.mod
├── frontend/                  # React App
│   ├── src/                   # Código fonte
│   ├── public/                # Arquivos públicos
│   └── package.json
├── docs/                      # Documentação
├── docker-compose.yml         # Orquestração
└── README.md
```

## Documentação

- [API](./docs/api.md) - Documentação completa da API
- [Importador](./docs/importer.md) - Guia do importador
- [Migrations](./docs/migrations.md) - Controle de schema
- [Execução Local](./docs/local_setup.md) - Setup local
- [Autenticação](./docs/auth.md) - Sistema de autenticação

## Tecnologias

### Backend
- Go 1.21+
- Chi Router
- PostgreSQL
- JWT
- golang-migrate
- excelize

### Frontend
- React
- Recharts
- Axios
- React Router

### Infraestrutura
- Docker
- Docker Compose
- Render (Backend)
- Vercel (Frontend)
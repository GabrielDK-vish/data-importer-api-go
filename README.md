# Data Importer API

Sistema de importação e análise de dados de faturamento desenvolvido em Golang com frontend React.

## Visão Geral

O projeto consiste em uma API REST em Golang que importa dados de faturamento de arquivos Excel/CSV para PostgreSQL, com interface web React para visualização de relatórios e gráficos.

## Arquitetura

### Backend (Golang)
- API REST com framework Chi
- Autenticação JWT
- Banco PostgreSQL com normalização completa
- Processamento de arquivos Excel/CSV
- Relatórios agregados por mês, produto e parceiro

### Frontend (React)
- Dashboard com gráficos interativos
- Upload de arquivos via interface web
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
- `POST /api/upload` - Upload de arquivos

## Upload de Arquivos

### Formatos Suportados
- CSV (.csv)
- Excel (.xlsx)

### Colunas Obrigatórias
- partner_id
- customer_id
- product_id
- usage_date
- quantity
- unit_price

### Comportamento
- Upload substitui completamente os dados existentes
- Processo atômico (tudo ou nada)
- Carregamento automático na inicialização se banco vazio

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
# Data Importer API

Sistema de importação e análise de dados de faturamento desenvolvido em Go com frontend React.

## Funcionalidades

- Importador de alta performance para arquivos Excel/CSV
- API RESTful com autenticação JWT
- Banco de dados PostgreSQL com normalização completa
- Frontend React com dashboard interativo
- Indicadores de performance (KPIs) dinâmicos
- Visualizações por categoria, recurso, cliente e período

## Tecnologias

### Backend
- Go 1.21+ com Chi Router
- PostgreSQL com normalização completa
- JWT para autenticação
- Excelize para processamento de arquivos

### Frontend
- React 18.2.0 com React Router
- Recharts para gráficos
- Axios para comunicação com API

### Infraestrutura
- Docker e Docker Compose
- Render (Backend) e Vercel (Frontend)
- PostgreSQL no Render

## Execução

### Com Docker (Recomendado)
```bash
git clone <repository-url>
cd data-importer-api-go
docker-compose up --build
```

**URLs de Acesso:**
- Frontend: http://localhost:3000
- API: http://localhost:8080
- PostgreSQL: localhost:5432

### Desenvolvimento Local
```bash
# Backend
cd backend && go mod tidy && go run ./cmd/main.go

# Frontend
cd frontend && npm install && npm start
```

## Credenciais de Teste

| Usuário | Senha    |
|---------|----------|
| admin   | admin123 |
| user    | user123  |
| demo    | demo123  |

## URLs de Produção

- **Frontend**: https://data-importer-api-go.vercel.app/
- **Backend**: https://data-importer-api-go.onrender.com/
- **Health Check**: https://data-importer-api-go.onrender.com/health

## Documentação

- [Arquitetura](./docs/architecture.md) - Arquitetura do sistema
- [API](./docs/api.md) - Documentação da API
- [Importador](./docs/importer.md) - Guia do importador
- [Migrations](./docs/migrations.md) - Controle de schema
- [Setup Local](./docs/local_setup.md) - Configuração local
- [Autenticação](./docs/auth.md) - Sistema de autenticação
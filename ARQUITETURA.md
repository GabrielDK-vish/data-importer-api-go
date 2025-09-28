# Arquitetura do Sistema Data Importer

## Visão Geral
Sistema full-stack para importação e análise de dados de faturamento desenvolvido em Go com frontend React.

## Diagrama de Arquitetura

```
┌─────────────────────────────────────────────────────────────────┐
│                        FRONTEND (React)                        │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐             │
│  │   Login     │  │  Dashboard  │  │  Reports    │             │
│  │   Page      │  │   Page     │  │   Page      │             │
│  └─────────────┘  └─────────────┘  └─────────────┘             │
│         │                 │                 │                  │
│         └─────────────────┼─────────────────┘                  │
│                           │                                    │
│  ┌─────────────────────────────────────────────────────────────┐ │
│  │                AuthContext & API Service                   │ │
│  │              (JWT Token Management)                        │ │
│  └─────────────────────────────────────────────────────────────┘ │
└─────────────────────────────────────────────────────────────────┘
                           │ HTTP/HTTPS
                           │ REST API
                           ▼
┌─────────────────────────────────────────────────────────────────┐
│                      BACKEND (Go/Golang)                       │
│  ┌─────────────────────────────────────────────────────────────┐ │
│  │                    API Layer                                │ │
│  │  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐        │ │
│  │  │   Routes    │  │   Upload    │  │    Auth     │        │ │
│  │  │  Handler    │  │  Handler    │  │  Handler    │        │ │
│  │  └─────────────┘  └─────────────┘  └─────────────┘        │ │
│  └─────────────────────────────────────────────────────────────┘ │
│                           │                                    │
│  ┌─────────────────────────────────────────────────────────────┐ │
│  │                  Service Layer                              │ │
│  │  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐        │ │
│  │  │   Business   │  │   Import    │  │   Report    │        │ │
│  │  │   Logic      │  │   Logic     │  │   Logic    │        │ │
│  │  └─────────────┘  └─────────────┘  └─────────────┘        │ │
│  └─────────────────────────────────────────────────────────────┘ │
│                           │                                    │
│  ┌─────────────────────────────────────────────────────────────┐ │
│  │                Repository Layer                              │ │
│  │  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐        │ │
│  │  │   Partner   │  │  Customer   │  │   Product   │        │ │
│  │  │   Repo      │  │    Repo     │  │    Repo     │        │ │
│  │  └─────────────┘  └─────────────┘  └─────────────┘        │ │
│  │  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐        │ │
│  │  │   Usage     │  │    User     │  │   Report    │        │ │
│  │  │    Repo     │  │    Repo     │  │    Repo     │        │ │
│  │  └─────────────┘  └─────────────┘  └─────────────┘        │ │
│  └─────────────────────────────────────────────────────────────┘ │
└─────────────────────────────────────────────────────────────────┘
                           │ SQL Queries
                           │ Database Connection
                           ▼
┌─────────────────────────────────────────────────────────────────┐
│                    DATABASE (PostgreSQL)                       │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐             │
│  │  Partners   │  │  Customers  │  │  Products   │             │
│  │   Table     │  │   Table     │  │   Table     │             │
│  └─────────────┘  └─────────────┘  └─────────────┘             │
│  ┌─────────────┐  ┌─────────────┐                              │
│  │   Usages    │  │    Users    │                              │
│  │   Table     │  │   Table     │                              │
│  └─────────────┘  └─────────────┘                              │
└─────────────────────────────────────────────────────────────────┘
```

## Componentes Principais

### Frontend (React)
- **Login Page**: Autenticação com JWT
- **Dashboard**: KPIs e gráficos interativos
- **Reports Page**: Relatórios detalhados
- **Customers Page**: Lista e detalhes de clientes
- **AuthContext**: Gerenciamento de estado de autenticação
- **API Service**: Comunicação com backend

### Backend (Go)
- **API Layer**: Handlers HTTP e rotas
- **Service Layer**: Lógica de negócio
- **Repository Layer**: Acesso a dados
- **Auth Module**: JWT e validação
- **Import Module**: Processamento de Excel/CSV

### Database (PostgreSQL)
- **partners**: Dados dos parceiros
- **customers**: Informações dos clientes
- **products**: Catálogo de produtos/serviços
- **usages**: Registros de uso e faturamento
- **users**: Usuários do sistema

## Fluxo de Dados

### 1. Autenticação
```
User → Login Form → API /auth/login → JWT Token → Frontend Storage
```

### 2. Importação de Dados
```
Excel/CSV File → Upload Handler → Excel Parser → Data Validation → Database Insert
```

### 3. Visualização de Dados
```
Frontend Request → API Endpoint → Service Layer → Repository → Database Query → Response
```

## Tecnologias Utilizadas

### Frontend
- React 18.2.0
- React Router DOM 6.3.0
- Recharts 2.5.0 (gráficos)
- Axios 1.4.0 (HTTP client)

### Backend
- Go 1.21+
- Chi Router 5.0.10
- PostgreSQL (driver: pgx/v5)
- JWT (golang-jwt/jwt/v5)
- Excelize 2.8.0 (processamento Excel)

### Infraestrutura
- Docker & Docker Compose
- Render (deploy backend)
- Vercel (deploy frontend)
- PostgreSQL (banco de dados)

## Padrões Arquiteturais

### Clean Architecture
- Separação clara entre camadas
- Inversão de dependências
- Testabilidade

### Repository Pattern
- Abstração do acesso a dados
- Facilita testes unitários
- Flexibilidade para mudanças de banco

### JWT Authentication
- Stateless authentication
- Segurança em APIs REST
- Token-based authorization

## Deploy e Infraestrutura

### Produção
- **Frontend**: Vercel (https://data-importer-api-go.vercel.app/)
- **Backend**: Render (https://data-importer-api-go.onrender.com/)
- **Database**: PostgreSQL no Render

### Desenvolvimento Local
- Docker Compose para orquestração
- Hot reload para desenvolvimento
- Banco PostgreSQL local

## Segurança

### Autenticação
- JWT tokens com expiração
- Validação de credenciais no banco
- Middleware de autenticação

### Validação de Dados
- Sanitização de inputs
- Validação de tipos
- Tratamento de erros robusto

### CORS
- Configuração adequada de CORS
- Headers de segurança
- Validação de origens

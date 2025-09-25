# Desafio Técnico — Full Stack (Golang)

##  Desafio proposto:
> Você deverá criar um **importador para uma base de dados (Postgres)** que deverá ser feito em **Golang** para armazenar os dados do arquivo enviado.  
> Será avaliado:  
> - **Normalização dos dados** na base de dados  
> - **Performance** do importador  
>  
> Além disso:  
> - Criar uma **API em Golang** com:
>   - Endpoint de **autenticação**
>   - Endpoints de **consulta dos dados importados**  
> - **Diferencial:** Criar um **frontend em React** mostrando indicadores totalizadores, agrupamentos por categorias, recursos, clientes e meses de cobrança.  
> - Necessário **publicar em algum link** para avaliação, com documentação de execução.

---

##  Solução proposta

### 1. **Importador/Input de Dados** 
   -  **CLI em Go** - Importador de linha de comando para CSV/Excel
   -  **Upload via Web** - Interface React para upload de arquivos
   -  **Processamento** - Conversão e normalização automática
   -  **Performance** - Inserção em lote com `pgx.CopyFrom`

### 2. **Banco de Dados (PostgreSQL)**
   Estrutura normalizada em 4 entidades:
   - **`partners`**: dados de parceiros (PartnerId, PartnerName, MpnId, Tier2MpnId)
   - **`customers`**: dados de clientes (CustomerId, CustomerName, CustomerDomainName, Country)
   - **`products`**: catálogo de serviços ou produtos (ProductId, SkuId, SkuName, ProductName, MeterType, Category, SubCategory, UnitType)
   - **`usages`**: KPIs vinculando a `partner_id`, `customer_id` e `product_id`
     (InvoiceNumber, ChargeStartDate, UsageDate, Quantity, UnitPrice, BillingPreTaxTotal, ResourceLocation, Tags, BenefitType)

   ➝ **Normalização completa** com relacionamentos e índices otimizados

### 3. **API REST (Golang)**
   -  **Framework**: `chi` com middleware
   -  **Autenticação**: JWT com tokens seguros
   -  **Endpoints implementados**:
     - `POST /auth/login` → autenticação de usuário
     - `POST /api/upload` → upload e processamento de arquivos
     - `GET /api/customers` → listar todos os clientes
     - `GET /api/customers/{id}/usage` → consumo detalhado por cliente
     - `GET /api/reports/billing/monthly` → faturamento por mês
     - `GET /api/reports/billing/by-product` → faturamento por produto
     - `GET /api/reports/billing/by-partner` → faturamento por parceiro

### 4. **Frontend (React)**
   -  **Dashboard** com indicadores:
     - Faturamento total por mês (gráficos)
     - Ranking de clientes por consumo
     - Distribuição por produtos/recursos
     - Métricas de performance
   -  **Páginas**:
     - Login com autenticação implementada
     - Dashboard principal
     - Lista de clientes com detalhes
     - Relatórios de faturamento
     - **Upload de arquivos** (CSV/Excel) via interface web
   -  **Tecnologias**: React, Recharts, Axios, React Router

### 5. **Infraestrutura e Deploy**
   -  **Docker Compose** - Postgres + API + Frontend
   -  **Migrations** - Controle de schema com `golang-migrate`
   -  **Scripts de execução** - Linux/Mac e Windows
   -  **Documentação completa** - Guias de execução local e Docker
   -  **Deploy em produção** - Render/Railway (API) + Vercel (Frontend)


---

## Estrutura do Projeto

```
data-importer-api-go/
├── backend/                    # API Golang
│   ├── cmd/
│   │   ├── main.go            # Servidor API
│   │   └── importer/          # CLI Importador
│   │       ├── main.go        # Importador CSV
│   │       └── excel_importer.go # Importador Excel
│   ├── internal/
│   │   ├── models/            # Estrutura de dados
│   │   ├── repository/        # Camadas de dados
│   │   ├── service/          # Lógica do negócio
│   │   ├── auth/             # Autenticação JWT
│   │   └── config/           # Configuração de ambiente
│   ├── api/
│   │   ├── routes.go         # Rotas da API
│   │   ├── upload.go         # Upload de arquivo
│   │   └── upload_processor.go # Processamento e normalização
│   ├── db/migrations/        # Migrações do banco
│   ├── Dockerfile           # Container da API
│   └── go.mod               # Dependências Go
│
├── frontend/                  # React Frontend
│   ├── src/
│   │   ├── pages/           # Páginas da aplicação
│   │   │   ├── Login.js     # Página de login
│   │   │   ├── Dashboard.js # Dashboard 
│   │   │   ├── Customers.js # Lista dos clientes
│   │   │   ├── Reports.js   # Relatórios
│   │   │   └── Upload.js    # Upload de arquivos
│   │   ├── services/        # Serviços
│   │   ├── App.js           # Componente principal
│   │   └── index.js         # Ponto de entrada
│   ├── public/              # Arquivos estáticos
│   ├── Dockerfile           # Container do Frontend
│   └── package.json         # Dependências Node.js
│
├── docs/                     # Documentação
│   ├── api.md               # Documentação da API
│   ├── importer.md          # Guia do importador
│   ├── migrations.md       # Guia de migrações
│   └── local_setup.md      # Execução local
│
├── scripts/                  # Scripts de automação
│   ├── run_docker.sh        # Execução Docker (Linux/Mac)
│   └── run_docker.ps1       # Execução Docker (Windows)
│
├── docker-compose.yml       # Orquestração de containers
├── README.md                # Documentação principal
├── QUICK_START.md           # Guia rápido
├── UPLOAD_FEATURES.md       # Funcionalidades de upload
└── Reconfile fornecedores.xlsx     # Dados de exemplo
```

---




## Como Executar

### Pré-requisitos do projeto
- Docker e Docker Compose
- Go 1.21+ (para desenvolvimento)
- Node.js 16+ (para desenvolvimento)

### 1. Inicialização com Docker

#### Opção A: Script Automático
```bash
# Linux/Mac
./run_docker.sh

# Windows PowerShell
.\run_docker.ps1
```

#### Opção B: Manual
```bash
# Clonar o repositório
git clone <https://github.com/GabrielDK-vish/data-importer-api-go.git>
cd data-importer-api-go

# Executar todos os serviços
docker-compose up --build

# Acessar:
# - Frontend: http://localhost:3000
# - API: http://localhost:8080
# - PostgreSQL: localhost:5432

# Importar dados (opcional)
# Copie o arquivo Excel para o container e execute:
docker-compose exec api go run ./cmd/importer/excel_importer.go /app/Reconfile\ fornecedores.xlsx
```



## Credenciais criadas para teste

| Usuário | Senha    |
|---------|----------|
| admin   | admin123 |
| user    | user123  |
| demo    | demo123  |


## Configuração

### Variáveis de Ambiente

```bash
# Banco de dados
DATABASE_URL=postgres://postgres:password@localhost:5432/data_importer?sslmode=disable

# JWT
JWT_SECRET=sua-chave-secreta-super-segura-aqui

# Servidor
PORT=8080
```

### Docker Compose

O `docker-compose.yml` configura:
- **PostgreSQL** com healthcheck
- **API Golang** com dependências
- **Volumes persistentes** para dados
- **Rede interna** para comunicação

## Dados de Exemplo

### Arquivo de Exemplo
`Reconfile fornecedores.csv`


### Arquivo Excel Fornecido
O projeto foi desenvolvido para trabalhar com o arquivo `Reconfile fornecedores.xlsx` fornecido no teste.


## Deploy
- **API:** _(N/A)_  
- **Frontend:** _(N/A)_  

---

## Documentação 
- [**API** (endpoints e exemplos)](./docs/api.md) - Documentação completa da API
- [**Importador** (CLI e Upload)](./docs/importer.md) - Guia do importador
- [**Migrations** (banco de dados)](./docs/migrations.md) - Controle de schema
- [**Execução Local** (desenvolvimento)](./docs/local_setup.md) - Setup local


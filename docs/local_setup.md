# 🖥️ Guia de Execução Local

## Pré-requisitos

### 1. Go 1.21+
```bash
# Verificar versão do Go
go version

# Se não tiver instalado, baixe em: https://golang.org/dl/
```

### 2. Node.js 16+
```bash
# Verificar versão do Node
node --version
npm --version

# Se não tiver instalado, baixe em: https://nodejs.org/
```

### 3. PostgreSQL 15+
```bash
# Verificar se PostgreSQL está instalado
psql --version

# Se não tiver instalado:
# Windows: https://www.postgresql.org/download/windows/
# macOS: brew install postgresql
# Ubuntu: sudo apt install postgresql postgresql-contrib
```

## Configuração do Banco de Dados

### 1. Criar Banco de Dados
```bash
# Conectar ao PostgreSQL
psql -U postgres

# Criar banco de dados
CREATE DATABASE data_importer;

# Criar usuário (opcional)
CREATE USER data_user WITH PASSWORD 'password';
GRANT ALL PRIVILEGES ON DATABASE data_importer TO data_user;

# Sair do psql
\q
```

### 2. Configurar Variáveis de Ambiente
```bash
# Windows (PowerShell)
$env:DATABASE_URL="postgres://postgres:password@localhost:5432/data_importer?sslmode=disable"
$env:JWT_SECRET="sua-chave-secreta-super-segura-aqui"
$env:PORT="8080"

# Windows (CMD)
set DATABASE_URL=postgres://postgres:password@localhost:5432/data_importer?sslmode=disable
set JWT_SECRET=sua-chave-secreta-super-segura-aqui
set PORT=8080

# Linux/Mac
export DATABASE_URL="postgres://postgres:password@localhost:5432/data_importer?sslmode=disable"
export JWT_SECRET="sua-chave-secreta-super-segura-aqui"
export PORT="8080"
```

## Execução do Backend

### 1. Instalar Dependências
```bash
cd backend

# Instalar dependências Go
go mod tidy

# Instalar golang-migrate (para migrations)
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
```

### 2. Executar Migrations
```bash
# Executar migrations
migrate -path ./db/migrations -database "postgres://postgres:password@localhost:5432/data_importer?sslmode=disable" up

# Verificar se as tabelas foram criadas
psql -U postgres -d data_importer -c "\dt"
```

### 3. Iniciar API
```bash
# Executar API
go run ./cmd/main.go
```

**Saída esperada:**
```
✅ Migrations executadas com sucesso
🚀 Servidor iniciado na porta 8080
📊 Endpoints disponíveis:
   POST /auth/login
   GET  /api/customers
   GET  /api/customers/{id}/usage
   GET  /api/reports/billing/monthly
   GET  /api/reports/billing/by-product
   GET  /api/reports/billing/by-partner
```

## Execução do Frontend

### 1. Instalar Dependências
```bash
cd frontend

# Instalar dependências Node.js
npm install
```

### 2. Iniciar Frontend
```bash
# Executar em modo desenvolvimento
npm start
```

**Saída esperada:**
```
Compiled successfully!

You can now view data-importer-frontend in the browser.

  Local:            http://localhost:3000
  On Your Network:  http://192.168.1.100:3000
```

## Importação de Dados

### 1. Importar Dados de Exemplo (CSV)
```bash
cd backend

# Importar dados CSV de exemplo
go run ./cmd/importer/main.go ../sample_data.csv
```

### 2. Importar Arquivo Excel
```bash
cd backend

# Importar arquivo Excel fornecido
go run ./cmd/importer/excel_importer.go ../Reconfile\ fornecedores.xlsx
```

## Verificação do Sistema

### 1. Testar API
```bash
# Testar autenticação
curl -X POST http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username": "admin", "password": "admin123"}'

# Testar endpoint protegido (usar token retornado)
curl -X GET http://localhost:8080/api/customers \
  -H "Authorization: Bearer <seu_token>"
```

### 2. Verificar Banco de Dados
```bash
# Conectar ao banco
psql -U postgres -d data_importer

# Verificar tabelas
\dt

# Contar registros
SELECT 'partners' as tabela, COUNT(*) as registros FROM partners
UNION ALL
SELECT 'customers', COUNT(*) FROM customers
UNION ALL
SELECT 'products', COUNT(*) FROM products
UNION ALL
SELECT 'usages', COUNT(*) FROM usages;

# Sair
\q
```

### 3. Acessar Frontend
- Abrir navegador em: http://localhost:3000
- Fazer login com: admin / admin123

## Troubleshooting

### Erro: "connection refused" (PostgreSQL)
```bash
# Verificar se PostgreSQL está rodando
# Windows
net start postgresql-x64-15

# macOS
brew services start postgresql

# Ubuntu
sudo systemctl start postgresql
```

### Erro: "database does not exist"
```bash
# Criar banco de dados
psql -U postgres -c "CREATE DATABASE data_importer;"
```

### Erro: "migrate: command not found"
```bash
# Instalar golang-migrate
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# Adicionar ao PATH (se necessário)
export PATH=$PATH:$(go env GOPATH)/bin
```

### Erro: "module not found" (Go)
```bash
# Limpar cache e reinstalar
go clean -modcache
go mod download
go mod tidy
```

### Erro: "npm install failed" (Node.js)
```bash
# Limpar cache e reinstalar
npm cache clean --force
rm -rf node_modules package-lock.json
npm install
```

### Erro: "port already in use"
```bash
# Verificar processos usando a porta
# Windows
netstat -ano | findstr :8080
netstat -ano | findstr :3000

# Linux/Mac
lsof -i :8080
lsof -i :3000

# Matar processo (substitua PID)
# Windows
taskkill /PID <PID> /F

# Linux/Mac
kill -9 <PID>
```

## Estrutura de Arquivos Local

```
data-importer-api-go/
├── backend/                    # API Golang
│   ├── cmd/
│   │   ├── main.go            # Servidor principal
│   │   └── importer/          # Importadores
│   ├── internal/              # Código interno
│   ├── api/                   # Handlers HTTP
│   ├── db/migrations/         # Migrations SQL
│   └── go.mod                 # Dependências Go
├── frontend/                   # React App
│   ├── src/                   # Código fonte
│   ├── public/                # Arquivos públicos
│   └── package.json           # Dependências Node.js
├── docs/                      # Documentação
├── sample_data.csv            # Dados de exemplo
├── Reconfile fornecedores.xlsx # Arquivo Excel fornecido
└── README.md
```

## Comandos Úteis

### Backend
```bash
# Executar API
go run ./cmd/main.go

# Executar importador CSV
go run ./cmd/importer/main.go ../sample_data.csv

# Executar importador Excel
go run ./cmd/importer/excel_importer.go ../Reconfile\ fornecedores.xlsx

# Executar testes
go test ./...

# Build para produção
go build -o api ./cmd/main.go
```

### Frontend
```bash
# Executar em desenvolvimento
npm start

# Build para produção
npm run build

# Executar testes
npm test

# Instalar dependências
npm install
```

### Banco de Dados
```bash
# Executar migrations
migrate -path ./db/migrations -database "postgres://postgres:password@localhost:5432/data_importer?sslmode=disable" up

# Reverter migrations
migrate -path ./db/migrations -database "postgres://postgres:password@localhost:5432/data_importer?sslmode=disable" down

# Verificar status das migrations
migrate -path ./db/migrations -database "postgres://postgres:password@localhost:5432/data_importer?sslmode=disable" version
```

## Próximos Passos

1. **Executar Backend**: `go run ./cmd/main.go`
2. **Executar Frontend**: `npm start`
3. **Importar Dados**: Usar importadores CSV ou Excel
4. **Acessar Sistema**: http://localhost:3000
5. **Testar API**: http://localhost:8080

---

**🎉 Sistema funcionando localmente!**

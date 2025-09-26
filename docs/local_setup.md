# Guia de Execução Local

## Pré-requisitos

### Go 1.21+
```bash
go version
# Download: https://golang.org/dl/
```

### Node.js 16+
```bash
node --version
npm --version
# Download: https://nodejs.org/
```

### PostgreSQL 15+
```bash
psql --version
# Windows: https://www.postgresql.org/download/windows/
# macOS: brew install postgresql
# Ubuntu: sudo apt install postgresql postgresql-contrib
```

## Configuração do Banco

### Criar Banco de Dados
```bash
psql -U postgres
CREATE DATABASE data_importer;
\q
```

### Variáveis de Ambiente
```bash
# Windows (PowerShell)
$env:DATABASE_URL="postgres://postgres:password@localhost:5432/data_importer?sslmode=disable"
$env:JWT_SECRET="sua-chave-secreta-super-segura-aqui"
$env:PORT="8080"

# Linux/Mac
export DATABASE_URL="postgres://postgres:password@localhost:5432/data_importer?sslmode=disable"
export JWT_SECRET="sua-chave-secreta-super-segura-aqui"
export PORT="8080"
```

## Execução do Backend

### Instalar Dependências
```bash
cd backend
go mod tidy
```

### Executar API
```bash
go run ./cmd/main.go
```

**Saída esperada:**
```
Migrations executadas com sucesso
Servidor iniciado na porta 8080
Endpoints disponíveis:
   POST /auth/login
   GET  /api/customers
   GET  /api/customers/{id}/usage
   GET  /api/reports/billing/monthly
   GET  /api/reports/billing/by-product
   GET  /api/reports/billing/by-partner
```

## Execução do Frontend

### Instalar Dependências
```bash
cd frontend
npm install
```

### Iniciar Frontend
```bash
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

### Via Interface Web
1. Acesse http://localhost:3000
2. Faça login (admin/admin123)
3. Vá para Upload
4. Selecione arquivo Excel/CSV
5. Aguarde processamento

### Via CLI
```bash
cd backend

# CSV
go run ./cmd/importer/main.go ../dados.csv

# Excel
go run ./cmd/importer/excel_importer.go ../dados.xlsx
```

## Verificação

### Testar API
```bash
# Login
curl -X POST http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username": "admin", "password": "admin123"}'

# Usar token retornado
curl -X GET http://localhost:8080/api/customers \
  -H "Authorization: Bearer <token>"
```

### Verificar Banco
```bash
psql -U postgres -d data_importer

# Contar registros
SELECT 'partners' as tabela, COUNT(*) as registros FROM partners
UNION ALL
SELECT 'customers', COUNT(*) FROM customers
UNION ALL
SELECT 'products', COUNT(*) FROM products
UNION ALL
SELECT 'usages', COUNT(*) FROM usages;

\q
```

### Acessar Frontend
- URL: http://localhost:3000
- Login: admin / admin123

## Troubleshooting

### PostgreSQL não inicia
```bash
# Windows
net start postgresql-x64-15

# macOS
brew services start postgresql

# Ubuntu
sudo systemctl start postgresql
```

### Banco não existe
```bash
psql -U postgres -c "CREATE DATABASE data_importer;"
```

### Porta em uso
```bash
# Windows
netstat -ano | findstr :8080
taskkill /PID <PID> /F

# Linux/Mac
lsof -i :8080
kill -9 <PID>
```

### Dependências Go
```bash
go clean -modcache
go mod download
go mod tidy
```

### Dependências Node
```bash
npm cache clean --force
rm -rf node_modules package-lock.json
npm install
```

## Estrutura Local

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
└── README.md
```

## Comandos Úteis

### Backend
```bash
go run ./cmd/main.go                    # Executar API
go run ./cmd/importer/main.go dados.csv # Importar CSV
go test ./...                           # Executar testes
go build -o api ./cmd/main.go          # Build produção
```

### Frontend
```bash
npm start          # Desenvolvimento
npm run build      # Build produção
npm test           # Executar testes
npm install        # Instalar dependências
```

### Banco
```bash
# Executar migrations
migrate -path ./db/migrations -database "postgres://postgres:password@localhost:5432/data_importer?sslmode=disable" up

# Reverter migrations
migrate -path ./db/migrations -database "postgres://postgres:password@localhost:5432/data_importer?sslmode=disable" down
```
# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Instalar dependências necessárias
RUN apk add --no-cache git

# Copiar go mod e sum do backend
COPY backend/go.mod backend/go.sum ./

# Baixar dependências
RUN go mod download

# Copiar código fonte do backend
COPY backend/ ./

# Copiar arquivo Excel da raiz do projeto
COPY Reconfile\ fornecedores.xlsx ./

# Build da aplicação
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/main.go

# Runtime stage
FROM alpine:latest

# Instalar ca-certificates para HTTPS
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copiar o binário
COPY --from=builder /app/main ./

# Copiar migrations
COPY --from=builder /app/db ./db

# Copiar arquivo Excel
COPY --from=builder /app/Reconfile\ fornecedores.xlsx ./

# Expor porta
EXPOSE 8080

# Comando para executar
CMD ["./main"]

#!/bin/bash

# Script para executar o projeto com Docker
echo "Iniciando Docker..."

# Verificar se Docker está rodando
if ! docker info > /dev/null 2>&1; then
    echo " Docker não está rodando. Por favor, inicie o Docker Desktop."
    exit 1
fi

# Verificar se docker-compose está disponível
if ! command -v docker-compose &> /dev/null; then
    echo " docker-compose não encontrado. Instale o Docker Compose."
    exit 1
fi

echo "Construindo e iniciando containers..."
docker-compose up --build -d

echo "Aguardando serviços iniciarem..."
sleep 10

echo "Verificando status dos serviços..."
docker-compose ps

echo ""
echo "Serviços iniciados com sucesso!"
echo ""
echo "Acesse:"
echo "   Frontend: http://localhost:3000"
echo "   API:      http://localhost:8080"
echo "   Banco:    localhost:5432"
echo ""
echo "Para importar dados:"
echo "   docker-compose exec api go run ./cmd/importer/excel_importer.go /app/Reconfile\\ fornecedores.xlsx"
echo ""
echo "Para parar os serviços:"
echo "   docker-compose down"

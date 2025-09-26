# Script PowerShell para executar o projeto com Docker
Write-Host "Iniciando Data Importer com Docker..." -ForegroundColor Green

# Verificar se Docker está rodando
try {
    docker info | Out-Null
} catch {
    Write-Host "Docker não está rodando. Por favor, inicie o Docker Desktop." -ForegroundColor Red
    exit 1
}

# Verificar se docker-compose está disponível
if (-not (Get-Command docker-compose -ErrorAction SilentlyContinue)) {
    Write-Host "docker-compose não encontrado. Instale o Docker Compose." -ForegroundColor Red
    exit 1
}

Write-Host "Construindo e iniciando containers..." -ForegroundColor Yellow
docker-compose up --build -d

Write-Host "Aguardando serviços iniciarem..." -ForegroundColor Yellow
Start-Sleep -Seconds 10

Write-Host "Verificando status dos serviços..." -ForegroundColor Yellow
docker-compose ps

Write-Host ""
Write-Host "Serviços iniciados com sucesso!" -ForegroundColor Green
Write-Host ""
Write-Host "Acesse:" -ForegroundColor Cyan
Write-Host "   Frontend: http://localhost:3000" -ForegroundColor White
Write-Host "   API:      http://localhost:8080" -ForegroundColor White
Write-Host "   Banco:    localhost:5432" -ForegroundColor White
Write-Host ""
Write-Host "Para importar dados:" -ForegroundColor Cyan
Write-Host "   docker-compose exec api go run ./cmd/importer/excel_importer.go /app/Reconfile\ fornecedores.xlsx" -ForegroundColor White
Write-Host ""
Write-Host "Para parar os serviços:" -ForegroundColor Cyan
Write-Host "   docker-compose down" -ForegroundColor White

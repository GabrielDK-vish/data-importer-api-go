# Script simples para testar upload

Write-Host "🚀 Testando upload do arquivo Reconfile fornecedores.xlsx" -ForegroundColor Green

# Verificar se a API está rodando
Write-Host "🔍 Verificando API..." -ForegroundColor Yellow
try {
    $response = Invoke-WebRequest -Uri "http://localhost:8080/health" -Method GET -TimeoutSec 5
    Write-Host "✅ API está rodando" -ForegroundColor Green
} catch {
    Write-Host "❌ API não está rodando. Inicie o servidor primeiro." -ForegroundColor Red
    Write-Host "Execute: cd backend; go run cmd/main.go" -ForegroundColor Yellow
    exit 1
}

# Fazer login
Write-Host "🔐 Fazendo login..." -ForegroundColor Yellow
$loginBody = '{"username":"admin","password":"admin123"}'
$loginResponse = Invoke-RestMethod -Uri "http://localhost:8080/auth/login" -Method POST -Body $loginBody -ContentType "application/json"
$token = $loginResponse.token
Write-Host "✅ Login realizado" -ForegroundColor Green

# Fazer upload
Write-Host "📤 Fazendo upload..." -ForegroundColor Yellow
$filePath = "Reconfile fornecedores.xlsx"
$headers = @{"Authorization" = "Bearer $token"}

try {
    $uploadResponse = Invoke-RestMethod -Uri "http://localhost:8080/api/upload" -Method POST -InFile $filePath -ContentType "multipart/form-data" -Headers $headers
    Write-Host "✅ Upload realizado com sucesso!" -ForegroundColor Green
    Write-Host "📊 Resultado: $($uploadResponse | ConvertTo-Json -Depth 3)" -ForegroundColor Cyan
} catch {
    Write-Host "❌ Erro no upload: $($_.Exception.Message)" -ForegroundColor Red
    Write-Host "Detalhes: $($_.Exception.Response)" -ForegroundColor Red
}

Write-Host "🎉 Teste concluído!" -ForegroundColor Green

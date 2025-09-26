# Script para testar upload usando curl

Write-Host "🚀 Testando upload com curl..." -ForegroundColor Green

# Fazer login e obter token
Write-Host "🔐 Fazendo login..." -ForegroundColor Yellow
$loginResponse = Invoke-RestMethod -Uri "https://data-importer-api-go.onrender.com/auth/login" -Method POST -Body '{"username":"admin","password":"admin123"}' -ContentType "application/json"
$token = $loginResponse.token
Write-Host "✅ Login realizado com sucesso" -ForegroundColor Green

# Criar arquivo temporário com o token
$tokenFile = "token.txt"
$token | Out-File -FilePath $tokenFile -Encoding UTF8

# Criar comando curl
$curlCommand = @"
curl -X POST "https://data-importer-api-go.onrender.com/api/upload" \
  -H "Authorization: Bearer $token" \
  -F "file=@Reconfile fornecedores.xlsx"
"@

Write-Host "📤 Executando upload..." -ForegroundColor Yellow
Write-Host "Comando: $curlCommand" -ForegroundColor Cyan

# Executar curl
try {
    $result = Invoke-Expression $curlCommand
    Write-Host "✅ Upload realizado com sucesso!" -ForegroundColor Green
    Write-Host "Resultado: $result" -ForegroundColor White
} catch {
    Write-Host "❌ Erro no upload: $($_.Exception.Message)" -ForegroundColor Red
}

# Limpar arquivo temporário
Remove-Item $tokenFile -ErrorAction SilentlyContinue

Write-Host "🎉 Teste concluído!" -ForegroundColor Green

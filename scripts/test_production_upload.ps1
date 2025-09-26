# Script para testar upload no ambiente de produção

Write-Host "🚀 Testando upload no ambiente de produção..." -ForegroundColor Green

# URL da API de produção
$apiUrl = "https://data-importer-api-go.onrender.com"

# Verificar se a API está rodando
Write-Host "🔍 Verificando API de produção..." -ForegroundColor Yellow
try {
    $healthResponse = Invoke-RestMethod -Uri "$apiUrl/health" -Method GET -TimeoutSec 10
    Write-Host "✅ API de produção está rodando: $($healthResponse.status)" -ForegroundColor Green
} catch {
    Write-Host "❌ API de produção não está respondendo: $($_.Exception.Message)" -ForegroundColor Red
    exit 1
}

# Fazer login
Write-Host "🔐 Fazendo login..." -ForegroundColor Yellow
$loginData = @{
    username = "admin"
    password = "admin123"
} | ConvertTo-Json

try {
    $loginResponse = Invoke-RestMethod -Uri "$apiUrl/auth/login" -Method POST -Body $loginData -ContentType "application/json" -TimeoutSec 15
    $token = $loginResponse.token
    Write-Host "✅ Login realizado com sucesso" -ForegroundColor Green
} catch {
    Write-Host "❌ Erro no login: $($_.Exception.Message)" -ForegroundColor Red
    if ($_.Exception.Response) {
        $errorStream = $_.Exception.Response.GetResponseStream()
        $reader = New-Object System.IO.StreamReader($errorStream)
        $errorBody = $reader.ReadToEnd()
        Write-Host "Detalhes do erro: $errorBody" -ForegroundColor Red
    }
    exit 1
}

# Preparar headers para upload
$headers = @{
    "Authorization" = "Bearer $token"
}

# Caminho do arquivo
$filePath = "Reconfile fornecedores.xlsx"

# Verificar se arquivo existe
if (-not (Test-Path $filePath)) {
    Write-Host "❌ Arquivo não encontrado: $filePath" -ForegroundColor Red
    exit 1
}

Write-Host "📁 Arquivo encontrado: $filePath" -ForegroundColor Green
Write-Host "📏 Tamanho: $((Get-Item $filePath).Length) bytes" -ForegroundColor Cyan

# Fazer upload
Write-Host "📤 Iniciando upload para produção..." -ForegroundColor Yellow
try {
    $form = @{
        file = Get-Item $filePath
    }
    
    $uploadResponse = Invoke-RestMethod -Uri "$apiUrl/api/upload" -Method POST -Form $form -Headers $headers -TimeoutSec 120
    Write-Host "✅ Upload realizado com sucesso!" -ForegroundColor Green
    Write-Host "📊 Resultado:" -ForegroundColor Cyan
    Write-Host "  - Sucesso: $($uploadResponse.success)" -ForegroundColor White
    Write-Host "  - Mensagem: $($uploadResponse.message)" -ForegroundColor White
    Write-Host "  - Parceiros: $($uploadResponse.data.partners)" -ForegroundColor White
    Write-Host "  - Clientes: $($uploadResponse.data.customers)" -ForegroundColor White
    Write-Host "  - Produtos: $($uploadResponse.data.products)" -ForegroundColor White
    Write-Host "  - Registros de uso: $($uploadResponse.data.usages)" -ForegroundColor White
    
} catch {
    Write-Host "❌ Erro no upload: $($_.Exception.Message)" -ForegroundColor Red
    if ($_.Exception.Response) {
        $errorStream = $_.Exception.Response.GetResponseStream()
        $reader = New-Object System.IO.StreamReader($errorStream)
        $errorBody = $reader.ReadToEnd()
        Write-Host "Detalhes do erro: $errorBody" -ForegroundColor Red
    }
    exit 1
}

Write-Host "🎉 Teste de upload em produção concluído!" -ForegroundColor Green
Write-Host "🌐 Acesse o frontend em: https://data-importer-api-go.vercel.app/" -ForegroundColor Cyan

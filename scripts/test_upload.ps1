# Script para testar upload do arquivo Excel
Write-Host "🔐 Obtendo token de autenticação..."

$body = @{
    username = "admin"
    password = "admin123"
} | ConvertTo-Json

try {
    $res = Invoke-RestMethod -Uri "https://data-importer-api-go.onrender.com/auth/login" -Method Post -Body $body -ContentType "application/json"
    $token = $res.token
    Write-Host "✅ Token obtido: $($token.Substring(0,20))..."
    
    Write-Host "📤 Fazendo upload do arquivo Excel..."
    $result = curl.exe -X POST "https://data-importer-api-go.onrender.com/api/upload" -H "Authorization: Bearer $token" -F "file=@../Reconfile fornecedores.xlsx"
    
    if ($LASTEXITCODE -eq 0) {
        Write-Host "✅ Upload realizado com sucesso!"
        Write-Host "Resposta: $result"
    } else {
        Write-Host "❌ Erro no upload (código: $LASTEXITCODE)"
        Write-Host "Resposta: $result"
    }
} catch {
    Write-Host "❌ Erro: $($_.Exception.Message)"
    if ($_.Exception.Response) {
        $reader = New-Object System.IO.StreamReader($_.Exception.Response.GetResponseStream())
        $responseBody = $reader.ReadToEnd()
        Write-Host "Response Body: $responseBody"
    }
}

# Script simples para verificar arquivo Excel
Write-Host "Verificando arquivo Excel..."

$excelFile = "Reconfile fornecedores.xlsx"

if (Test-Path $excelFile) {
    Write-Host "Arquivo encontrado: $excelFile"
    $fileInfo = Get-Item $excelFile
    Write-Host "Tamanho: $([math]::Round($fileInfo.Length / 1MB, 2)) MB"
} else {
    Write-Host "Arquivo nao encontrado: $excelFile"
    exit 1
}

Write-Host "Arquivo existe e pode ser processado"

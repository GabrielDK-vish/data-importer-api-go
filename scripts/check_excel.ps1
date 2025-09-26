# Script para verificar arquivo Excel
Write-Host "🔍 Verificando arquivo Excel..."

$excelFile = "..\Reconfile fornecedores.xlsx"

if (Test-Path $excelFile) {
    Write-Host "✅ Arquivo encontrado: $excelFile"
    $fileInfo = Get-Item $excelFile
    Write-Host "📏 Tamanho: $([math]::Round($fileInfo.Length / 1MB, 2)) MB"
    Write-Host "📅 Última modificação: $($fileInfo.LastWriteTime)"
} else {
    Write-Host "❌ Arquivo não encontrado: $excelFile"
    exit 1
}

# Verificar se o arquivo pode ser aberto
try {
    $excel = New-Object -ComObject Excel.Application
    $excel.Visible = $false
    $workbook = $excel.Workbooks.Open((Resolve-Path $excelFile).Path)
    $worksheet = $workbook.Worksheets.Item(1)
    
    Write-Host "📊 Planilha: $($worksheet.Name)"
    Write-Host "📏 Linhas usadas: $($worksheet.UsedRange.Rows.Count)"
    Write-Host "📏 Colunas usadas: $($worksheet.UsedRange.Columns.Count)"
    
    # Ler cabeçalhos
    Write-Host "`nCabeçalhos encontrados:"
    for ($col = 1; $col -le $worksheet.UsedRange.Columns.Count; $col++) {
        $header = $worksheet.Cells.Item(1, $col).Value2
        if ($header) {
            Write-Host "  $($col.ToString().PadLeft(2)): $header"
        }
    }
    
    $workbook.Close($false)
    $excel.Quit()
    [System.Runtime.Interopservices.Marshal]::ReleaseComObject($excel) | Out-Null
    
    Write-Host "`nArquivo Excel é válido e pode ser processado"
} catch {
    Write-Host "❌ Erro ao abrir arquivo Excel: $($_.Exception.Message)"
    exit 1
}

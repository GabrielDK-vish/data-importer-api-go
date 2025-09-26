# Guia do Importador

## Visão Geral

Sistema de importação de dados CSV e Excel para PostgreSQL com alta performance usando inserção em lotes.

## Formatos Suportados

- CSV (.csv) - Separado por vírgula ou tab
- Excel (.xlsx) - Planilhas Microsoft Excel

## Colunas Obrigatórias

- partner_id
- customer_id
- product_id
- usage_date
- quantity
- unit_price

## Colunas Opcionais

- partner_name, mpn_id, tier2_mpn_id
- customer_name, customer_domain_name, country
- sku_id, sku_name, product_name, meter_type, category, sub_category, unit_type
- invoice_number, charge_start_date, billing_pre_tax_total
- resource_location, tags, benefit_type

## Execução

### Via API (Recomendado)
```bash
# Upload via interface web ou API
curl -X POST https://data-importer-api-go.onrender.com/api/upload \
  -H "Authorization: Bearer <token>" \
  -F "file=@dados.xlsx"
```

### Via CLI Local
```bash
cd backend

# CSV
go run ./cmd/importer/main.go dados.csv

# Excel
go run ./cmd/importer/excel_importer.go dados.xlsx
```

### Via Docker
```bash
# CSV
docker-compose exec api go run ./cmd/importer/main.go /app/dados.csv

# Excel
docker-compose exec api go run ./cmd/importer/excel_importer.go /app/dados.xlsx
```

## Processamento

### Normalização de Dados
Os dados são normalizados em 4 entidades:

**Partners**
```go
type Partner struct {
    PartnerID   string
    PartnerName string
    MpnID       string
    Tier2MpnID  string
}
```

**Customers**
```go
type Customer struct {
    CustomerID         string
    CustomerName       string
    CustomerDomainName string
    Country            string
}
```

**Products**
```go
type Product struct {
    ProductID   string
    SkuID       string
    SkuName     string
    ProductName string
    MeterType   string
    Category    string
    SubCategory string
    UnitType    string
}
```

**Usages**
```go
type Usage struct {
    InvoiceNumber      string
    ChargeStartDate    time.Time
    UsageDate          time.Time
    Quantity           float64
    UnitPrice          float64
    BillingPreTaxTotal float64
    ResourceLocation   string
    Tags               string
    BenefitType        string
}
```

### Mapeamento de Colunas
O sistema mapeia automaticamente colunas com aliases:

- PartnerId → partner_id
- CustomerCountry → country
- Unit → unit_type
- MeterCategory → category
- E muitas outras...

### Formatos de Data Suportados
- 2006-01-02 (ISO)
- 2006/01/02
- 02/01/2006
- 02-01-2006
- 1/2/2006
- Números seriais do Excel

## Performance

### Otimizações
- Inserção em lotes com pgx.CopyFrom
- Processamento em lotes de 1000 registros
- Validação otimizada de tipos
- Uso eficiente de memória

### Métricas Típicas
- Velocidade: ~10.000 registros/segundo
- Memória: ~50MB para arquivos de 100MB
- Tempo: ~30 segundos para 100.000 registros

## Comportamento do Upload

### Substituição Completa
- Upload substitui completamente dados existentes
- Processo atômico (tudo ou nada)
- Limpeza automática na ordem correta

### Carregamento Automático
- Sistema verifica dados na inicialização
- Carrega automaticamente se banco vazio
- Evita duplicação se dados existem

## Tratamento de Erros

### Validação de Colunas
```
Colunas obrigatórias não encontradas: [partner_id, customer_id]
Colunas disponíveis: [PartnerId, CustomerId, ProductId, ...]
```

### Parsing de Dados
```
Erro ao processar linha 15: coluna 'quantity' inválida
Erro ao parsear usage_date: formato de data inválido: 15/01/2024
```

### Logs de Processamento
```
Processado lote: 1000 registros
Processado último lote: 250 registros
Total processado: 1250 registros de 1250 linhas
```

## Troubleshooting

### Erro: "coluna obrigatória não encontrada"
Verificar se arquivo possui colunas obrigatórias.

### Erro: "formato de data inválido"
Usar formatos suportados ou verificar dados.

### Erro: "erro ao conectar ao banco"
Verificar se PostgreSQL está rodando.

### Erro: "arquivo Excel não possui planilhas"
Verificar se arquivo Excel é válido.

## Exemplos

### CSV de Exemplo
```csv
partner_id,partner_name,customer_id,customer_name,product_id,product_name,usage_date,quantity,unit_price,billing_pre_tax_total
PARTNER001,Microsoft Corporation,CUST001,TechCorp Solutions,PROD001,Azure Virtual Machine,2024-01-15,100.0,0.05,5.0
PARTNER002,Amazon Web Services,CUST002,DataFlow Inc,PROD002,AWS EC2 Instance,2024-01-16,200.0,0.08,16.0
```

### Verificação de Dados
```bash
# Conectar ao banco
psql -U postgres -d data_importer

# Contar registros
SELECT 'partners' as tabela, COUNT(*) as registros FROM partners
UNION ALL
SELECT 'customers', COUNT(*) FROM customers
UNION ALL
SELECT 'products', COUNT(*) FROM products
UNION ALL
SELECT 'usages', COUNT(*) FROM usages;
```
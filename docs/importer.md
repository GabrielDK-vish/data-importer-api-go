# Processamento de Dados

## Visão Geral

Sistema de processamento de dados Excel para PostgreSQL com alta performance usando inserção em lotes e normalização de dados.

## Estrutura de Dados

O sistema processa arquivos Excel e normaliza os dados nas seguintes entidades:

- **Partners**: Informações dos parceiros
- **Customers**: Dados dos clientes
- **Products**: Catálogo de produtos com categorias e tipos de recursos
- **Usages**: Registros de uso e faturamento

## Campos Processados

### Dados de Parceiros
- ID do parceiro
- Nome do parceiro
- Identificadores MPN

### Dados de Clientes
- ID do cliente
- Nome do cliente
- Domínio
- País

### Dados de Produtos
- ID do produto
- Nome do produto
- SKU
- Categoria
- Tipo de recurso
- Tipo de unidade

### Dados de Uso
- Data de uso
- Quantidade
- Preço unitário
- Total pré-impostos
- Localização do recurso

## Processamento Automático

O sistema processa automaticamente o arquivo Excel na inicialização, extraindo e normalizando os dados para o banco PostgreSQL.

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
O sistema possui mapeamento automático inteligente que reconhece variações dos nomes das colunas:

**Mapeamento Automático:**
- **Partner ID**: `PartnerId`, `Partner_ID`, `partner-id`, `partner id`
- **Customer ID**: `CustomerId`, `Customer_ID`, `customer-id`, `customer id`
- **Product ID**: `ProductId`, `Product_ID`, `product-id`, `product id`
- **Usage Date**: `UsageDate`, `Usage_Date`, `usage-date`, `usage date`, `Date`
- **Quantity**: `Quantity`, `Qty`, `quantity`, `qty`
- **Unit Price**: `UnitPrice`, `Unit_Price`, `unit-price`, `unit price`, `Price`

**Aliases Suportados:**
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
- 2006-01-02T15:04:05Z
- 2006-01-02T15:04:05.000Z
- Números seriais do Excel (convertidos automaticamente)

### Formatos de Números Suportados
- Vírgula como separador decimal (ex: 1,50)
- Ponto como separador decimal (ex: 1.50)
- Formato brasileiro com pontos para milhares (ex: 1.000,50)
- Caracteres não numéricos são removidos automaticamente

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


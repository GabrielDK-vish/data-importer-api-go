# Guia do Importador CLI

## Visão Geral

O Importador CLI é uma ferramenta desenvolvida em Golang para importar dados de arquivos CSV e Excel (.xlsx) para o banco de dados PostgreSQL. Foi projetado para alta performance usando inserção em lotes com `pgx.CopyFrom`.

## Tipos de Arquivo Suportados

- **CSV** - Arquivos separados por vírgula
- **Excel (.xlsx)** - Planilhas do Microsoft Excel

## Funcionalidades

- **Leitura de CSV e Excel** com validação de colunas
- **Normalização automática** dos dados
- **Inserção em lotes** para alta performance
- **Tratamento de erros** robusto
- **Validação de tipos** de dados
- **Suporte a múltiplos formatos** de data
- **Conversão automática** de números seriais do Excel
- **Carregamento automático** de dados na inicialização
- **Substituição inteligente** de dados existentes

## Estrutura do Arquivo CSV

### Colunas Obrigatórias

```csv
partner_id,customer_id,product_id,usage_date,quantity,unit_price
```

### Colunas Opcionais

```csv
partner_name,mpn_id,tier2_mpn_id,customer_name,customer_domain_name,country,
sku_id,sku_name,product_name,meter_type,category,sub_category,unit_type,
invoice_number,charge_start_date,billing_pre_tax_total,resource_location,tags,benefit_type
```

### Exemplo de CSV

```csv
partner_id,partner_name,customer_id,customer_name,product_id,product_name,usage_date,quantity,unit_price,billing_pre_tax_total
PARTNER001,Microsoft Corporation,CUST001,TechCorp Solutions,PROD001,Azure Virtual Machine,2024-01-15,100.0,0.05,5.0
PARTNER002,Amazon Web Services,CUST002,DataFlow Inc,PROD002,AWS EC2 Instance,2024-01-16,200.0,0.08,16.0
```

## Como Usar

### 1. Executar com Docker

```bash
# Importar arquivo CSV
docker-compose exec api go run ./cmd/importer/main.go /caminho/para/arquivo.csv

# Importar arquivo Excel
docker-compose exec api go run ./cmd/importer/excel_importer.go /caminho/para/arquivo.xlsx

# Exemplo com dados de exemplo
docker-compose exec api go run ./cmd/importer/main.go /app/sample_data.csv
```

### 2. Executar Localmente

```bash
cd backend

# Configurar variáveis de ambiente
export DATABASE_URL="postgres://postgres:password@localhost:5432/data_importer?sslmode=disable"

# Executar importador CSV
go run ./cmd/importer/main.go ../sample_data.csv

# Executar importador Excel
go run ./cmd/importer/excel_importer.go ../Reconfile\ fornecedores.xlsx
```

### 3. Sintaxe de Comando

```bash
# Para arquivos CSV
go run ./cmd/importer/main.go <arquivo.csv>

# Para arquivos Excel
go run ./cmd/importer/excel_importer.go <arquivo.xlsx>
```

## Processamento de Dados

### 1. Validação de Colunas

O importador verifica se as colunas obrigatórias estão presentes:
- `partner_id`
- `customer_id` 
- `product_id`
- `usage_date`
- `quantity`
- `unit_price`

### 2. Suporte a Excel (.xlsx)

O importador Excel possui funcionalidades específicas:

#### Características do Excel
- **Múltiplas planilhas**: Usa a primeira planilha automaticamente
- **Números seriais**: Converte automaticamente datas em formato serial do Excel
- **Formatação**: Preserva formatação de números e datas
- **Células vazias**: Trata células vazias como valores padrão

#### Formatos de Data Suportados no Excel
- Datas em formato texto: `2024-01-15`, `15/01/2024`
- Números seriais do Excel: `45285` (converte para 2024-01-15)
- Datas com hora: `2024-01-15 10:30:00`

#### Exemplo de Uso com Excel
```bash
# Copiar arquivo Excel para o container
docker cp "Reconfile fornecedores.xlsx" data-importer-api-go_api_1:/app/

# Importar arquivo Excel
docker-compose exec api go run ./cmd/importer/excel_importer.go /app/Reconfile\ fornecedores.xlsx
```

### 2. Normalização

Os dados são normalizados em 4 entidades:

#### Partners (Parceiros)
```go
type Partner struct {
    PartnerID   string
    PartnerName string
    MpnID       string
    Tier2MpnID  string
}
```

#### Customers (Clientes)
```go
type Customer struct {
    CustomerID         string
    CustomerName       string
    CustomerDomainName string
    Country            string
}
```

#### Products (Produtos)
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

#### Usages (Usos)
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

### 3. Inserção em Lotes

O importador processa os dados em lotes de 1000 registros para otimizar a performance:

```go
// Processar lote quando atingir o tamanho
if len(usages) >= batchSize {
    service.ProcessImportData(ctx, partners, customers, products, usages)
    // Limpar lotes
    partners = partners[:0]
    customers = customers[:0]
    products = products[:0]
    usages = usages[:0]
}
```

## Formatos de Data Suportados

O importador suporta múltiplos formatos de data:

- `2006-01-02` (ISO)
- `2006/01/02`
- `02/01/2006`
- `02-01-2006`

## Tratamento de Erros

### 1. Validação de Tipos

```go
// Converter string para float
parseFloat := func(value string) (float64, error) {
    if value == "" {
        return 0, nil
    }
    return strconv.ParseFloat(value, 64)
}
```

### 2. Logs de Erro

```
Erro ao ler linha 15: coluna 'quantity' inválida
Erro ao processar linha 23: formato de data inválido: 15/01/2024
```

### 3. Estatísticas de Processamento

```
Processado lote: 1000 registros
Processado último lote: 250 registros
Total processado: 1250 registros de 1250 linhas
```

## Performance

### Métricas Típicas

- **Velocidade**: ~10.000 registros/segundo
- **Memória**: ~50MB para arquivos de 100MB
- **Tempo**: ~30 segundos para 100.000 registros

### Otimizações Implementadas

1. **Inserção em lotes** usando `pgx.CopyFrom`
2. **Processamento paralelo** de dados
3. **Validação otimizada** de tipos
4. **Uso eficiente de memória**

## Exemplos de Uso

### 1. Importar Dados de Exemplo

```bash
# Usar dados de exemplo incluídos
docker-compose exec api go run ./cmd/importer/main.go /app/sample_data.csv
```

### 2. Importar Arquivo Personalizado

```bash
# Copiar arquivo para o container
docker cp meu_arquivo.csv data-importer-api-go_api_1:/app/

# Importar
docker-compose exec api go run ./cmd/importer/main.go /app/meu_arquivo.csv
```

### 3. Verificar Dados Importados

```bash
# Conectar ao banco
docker-compose exec postgres psql -U postgres -d data_importer

# Verificar tabelas
\dt

# Contar registros
SELECT COUNT(*) FROM partners;
SELECT COUNT(*) FROM customers;
SELECT COUNT(*) FROM products;
SELECT COUNT(*) FROM usages;
```

## Troubleshooting

### Erro: "coluna obrigatória não encontrada"

**Causa**: Arquivo CSV não possui colunas obrigatórias.

**Solução**: Verificar se o CSV possui as colunas:
- `partner_id`
- `customer_id`
- `product_id`
- `usage_date`
- `quantity`
- `unit_price`

### Erro: "formato de data inválido"

**Causa**: Data em formato não suportado.

**Solução**: Usar formatos suportados:
- `2024-01-15`
- `15/01/2024`
- `15-01-2024`

### Erro: "erro ao conectar ao banco"

**Causa**: Banco de dados não está rodando.

**Solução**: 
```bash
# Verificar se PostgreSQL está rodando
docker-compose ps

# Iniciar serviços
docker-compose up -d postgres
```

## Estrutura do Código

```
cmd/importer/
└── main.go              # Ponto de entrada do importador
    ├── processCSV()     # Processa arquivo CSV
    ├── parseRow()       # Parseia linha do CSV
    ├── parseFloat()     # Converte string para float
    └── parseDate()      # Converte string para data
```

## Carregamento Automático

O sistema possui carregamento automático de dados na inicialização:

### Funcionamento
- **Verificação Automática**: Na inicialização, o sistema verifica se existem dados no banco
- **Carregamento Inicial**: Se o banco estiver vazio, carrega automaticamente os dados do arquivo "Reconfile fornecedores.xlsx"
- **Evita Duplicação**: Se já existirem dados, não recarrega automaticamente

### Logs Esperados
```
Dados já existem no banco (X clientes encontrados)
```
ou
```
Carregando dados iniciais do arquivo: Reconfile fornecedores.xlsx
Dados iniciais carregados com sucesso
```

## Substituição de Dados

### Comportamento do Upload
- **Substituição Completa**: Upload de novos arquivos substitui completamente os dados existentes
- **Processo Atômico**: A operação é tudo ou nada (sem dados parciais)
- **Limpeza Automática**: Sistema limpa automaticamente as tabelas na ordem correta

### Logs de Substituição
```
Iniciando substituição de dados: X partners, Y customers, Z products, W usages
Dados substituídos com sucesso
```

## Próximos Passos

- [x] Suporte a arquivos Excel (.xlsx)
- [x] Carregamento automático de dados
- [x] Substituição inteligente de dados
- [ ] Modo dry-run para teste
- [ ] Relatório de importação detalhado
- [ ] Suporte a arquivos grandes (>1GB)
- [ ] Backup automático antes da substituição

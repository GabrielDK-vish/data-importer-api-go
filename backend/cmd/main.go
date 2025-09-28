package main

import (
	"context"
	"data-importer-api-go/api"
	"data-importer-api-go/internal/config"
	"data-importer-api-go/internal/models"
	"data-importer-api-go/internal/repository"
	"data-importer-api-go/internal/service"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/xuri/excelize/v2"
)

func main() {
	// Carregar configuração
	cfg := config.LoadConfig()

	// Conectar ao banco de dados
	db, err := pgxpool.New(context.Background(), cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Erro ao conectar ao banco de dados: %v", err)
	}
	defer db.Close()

	// Executar migrations
	if err := runMigrations(cfg.DatabaseURL); err != nil {
		log.Fatalf("Erro ao executar migrations: %v", err)
	}

	// Inicializar camadas
	repo := repository.NewRepository(db)
	svc := service.NewService(repo)

	// Carregar dados iniciais se não existirem (não deve encerrar a API em caso de erro)
	if err := loadInitialData(svc); err != nil {
		log.Printf("⚠️  Aviso: Não foi possível carregar dados iniciais: %v", err)
	}
	handler := api.NewHandler(svc)

	// Configurar rotas
	router := handler.SetupRoutes()

	// Configurar servidor
	server := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
	}

	// Iniciar servidor em goroutine
	go func() {
		log.Printf("Servidor iniciado na porta %s", cfg.Port)
		log.Printf("Endpoints disponíveis:")
		log.Printf("   POST /auth/login")
		log.Printf("   GET  /api/customers")
		log.Printf("   GET  /api/customers/{id}/usage")
		log.Printf("   GET  /api/reports/billing/monthly")
		log.Printf("   GET  /api/reports/billing/by-product")
		log.Printf("   GET  /api/reports/billing/by-partner")
		
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Erro ao iniciar servidor: %v", err)
		}
	}()

	// Aguardar sinal de interrupção
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println(" Parando servidor...")

	// Graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Erro ao parar servidor: %v", err)
	}

	log.Println("Servidor parado com sucesso")
}

func loadInitialData(svc *service.Service) (err error) {
	// Proteger contra panics para não derrubar o servidor
	defer func() {
		if r := recover(); r != nil {
			log.Printf("⚠️  Panic ao carregar dados iniciais: %v", r)
			err = fmt.Errorf("panic ao carregar dados iniciais")
		}
	}()

	// Verificar se já existem dados no banco
	ctx := context.Background()
	customers, err := svc.GetAllCustomers(ctx)
	if err != nil {
		return fmt.Errorf("erro ao verificar dados existentes: %w", err)
	}

	// Se já existem dados, não carregar novamente
	if len(customers) > 0 {
		log.Printf("Dados já existem no banco (%d clientes encontrados)", len(customers))
		return nil
	}

    // Tentar localizar o arquivo Excel inicial em múltiplos caminhos
    candidateFiles := []string{
        "Reconfile fornecedores.xlsx",                 // diretório atual (backend)
        "../Reconfile fornecedores.xlsx",              // raiz do repo
        "/app/Reconfile fornecedores.xlsx",            // caminho no container
        "/root/Reconfile fornecedores.xlsx",           // caminho no container runtime
        "./Reconfile fornecedores.xlsx",               // diretório atual explícito
    }
    var excelFile string
    for _, path := range candidateFiles {
        log.Printf("Tentando localizar arquivo em: %s", path)
        if _, err := os.Stat(path); err == nil {
            excelFile = path
            log.Printf("✅ Arquivo encontrado em: %s", path)
            break
        } else {
            log.Printf("❌ Arquivo não encontrado em: %s (erro: %v)", path, err)
        }
    }
    if excelFile == "" {
        log.Printf("⚠️  Arquivo 'Reconfile fornecedores.xlsx' não encontrado em nenhum caminho padrão, pulando carregamento inicial")
        log.Printf("📁 Caminhos testados: %v", candidateFiles)
        return nil
    }

    log.Printf("Carregando dados iniciais do arquivo: %s", excelFile)
	
	// Usar o importador Excel existente
	if err := processExcelFile(svc, excelFile); err != nil {
		return fmt.Errorf("erro ao processar arquivo inicial: %w", err)
	}

	log.Printf("Dados iniciais carregados com sucesso")
	return nil
}

func runMigrations(databaseURL string) error {
	m, err := migrate.New(
		"file://db/migrations",
		databaseURL,
	)
	if err != nil {
		return fmt.Errorf("erro ao criar migrator: %w", err)
	}
	defer m.Close()

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("erro ao executar migrations: %w", err)
	}
 
	log.Println("Migrations executadas com sucesso")
	return nil
}

func processExcelFile(svc *service.Service, filename string) error {
	// Abrir arquivo Excel
	f, err := excelize.OpenFile(filename)
	if err != nil {
		return fmt.Errorf("erro ao abrir arquivo Excel: %w", err)
	}
	defer f.Close()

	// Obter todas as planilhas
	sheetList := f.GetSheetList()
	if len(sheetList) == 0 {
		return fmt.Errorf("arquivo Excel não possui planilhas")
	}

	// Usar a primeira planilha
	sheetName := sheetList[0]
	log.Printf("Processando planilha: %s", sheetName)

	// Obter todas as linhas
	rows, err := f.GetRows(sheetName)
	if err != nil {
		return fmt.Errorf("erro ao ler linhas da planilha: %w", err)
	}

	if len(rows) < 2 {
		return fmt.Errorf("planilha deve ter pelo menos 2 linhas (cabeçalho + dados)")
	}

	// Processar cabeçalho
	header := rows[0]
	log.Printf("Cabeçalhos encontrados: %v", header)

	// Mapear índices das colunas com mapeamento inteligente (igual ao upload.go)
	columnMap := make(map[string]int)
	
	// Alias para mapear diferentes formatos de colunas
	alias := map[string]string{
		// Partner fields
		"partnerid":              "partner_id",
		"partnername":            "partner_name",
		"mpnid":                  "mpn_id",
		"tier2mpnid":             "tier2_mpn_id",
		"tier2mpn":               "tier2_mpn_id",
		
		// Customer fields
		"customerid":             "customer_id",
		"customername":           "customer_name",
		"customerdomainname":     "customer_domain_name",
		"customercountry":        "country",
		"customerdomain":         "customer_domain_name",
		
		// Product fields
		"productid":              "product_id",
		"skuid":                  "sku_id",
		"skuname":                "sku_name",
		"productname":            "product_name",
		"metertype":              "meter_type",
		"metercategory":          "category",
		"metersubcategory":       "sub_category",
		"unittype":               "unit_type",
		
		// Usage fields
		"invoicenumber":          "invoice_number",
		"chargestartdate":        "charge_start_date",
		"usagedate":              "usage_date",
		"quantity":               "quantity",
		"unitprice":              "unit_price",
		"billingpretaxtotal":     "billing_pre_tax_total",
		"resourcelocation":       "resource_location",
		"tags":                   "tags",
		"benefittype":            "benefit_type",
		
		// Novos campos identificados nos logs
		"availabilityid":         "availability_id",
		"publishername":          "publisher_name",
		"publisherid":            "publisher_id",
		"subscriptiondescription": "subscription_description",
		"subscriptionid":         "subscription_id",
		"chargeenddate":           "charge_end_date",
		"meterid":                 "meter_id",
		"metername":              "meter_name",
		"meterregion":            "meter_region",
		"unit":                   "unit",
		"consumedservice":        "consumed_service",
		"resourcegroup":          "resource_group",
		"resourceuri":            "resource_uri",
		"chargetype":             "charge_type",
		"billingcurrency":        "billing_currency",
		"pricingpretaxtotal":     "pricing_pre_tax_total",
		"pricingcurrency":        "pricing_currency",
		"serviceinfo1":           "service_info1",
		"serviceinfo2":           "service_info2",
		"additionalinfo":         "additional_info",
		"effectiveunitprice":     "effective_unit_price",
		"pctobcexchangerate":     "pc_to_bc_exchange_rate",
		"pctobcexchangeratedate": "pc_to_bc_exchange_rate_date",
		"entitlementid":          "entitlement_id",
		"entitlementdescription": "entitlement_description",
		"partnerearnedcreditpercentage": "partner_earned_credit_percentage",
		"creditpercentage":       "credit_percentage",
		"credittype":             "credit_type",
		"benefitorderid":         "benefit_order_id",
		"benefitid":              "benefit_id",
	}
	
	// Normalizar cabeçalhos e aplicar aliases
	normalize := func(s string) string {
		s = strings.ToLower(strings.TrimSpace(s))
		s = strings.ReplaceAll(s, " ", "")
		s = strings.ReplaceAll(s, "_", "")
		s = strings.ReplaceAll(s, "-", "")
		return s
	}
	
	for i, col := range header {
		n := normalize(col)
		key := n
		if mapped, ok := alias[n]; ok {
			key = mapped
		}
		columnMap[key] = i
		log.Printf("Mapeando coluna '%s' -> '%s' -> '%s' (índice %d)", col, n, key, i)
	}

	// Verificar colunas obrigatórias
	requiredColumns := []string{"partner_id", "customer_id", "product_id", "usage_date", "quantity", "unit_price"}
	for _, col := range requiredColumns {
		if _, exists := columnMap[col]; !exists {
			return fmt.Errorf("coluna obrigatória não encontrada: %s", col)
		}
	}

	// Processar dados em lotes
	batchSize := 1000
	var partners []models.Partner
	var customers []models.Customer
	var products []models.Product
	var usages []models.Usage

	partnerMap := make(map[string]*models.Partner)
	customerMap := make(map[string]*models.Customer)
	productMap := make(map[string]*models.Product)

	rowCount := 0
	processedCount := 0

	// Processar linhas de dados (pular cabeçalho)
	for i := 1; i < len(rows); i++ {
		record := rows[i]
		rowCount++

		// Pular linhas vazias
		if len(record) == 0 || allEmpty(record) {
			continue
		}

		// Processar linha
		partner, customer, product, usage, err := parseExcelRow(record, columnMap, rowCount)
		if err != nil {
			log.Printf("Erro ao processar linha %d: %v", rowCount, err)
			continue
		}

		// Adicionar partner se não existir
		if _, exists := partnerMap[partner.PartnerID]; !exists {
			partners = append(partners, *partner)
			partnerMap[partner.PartnerID] = partner
		}

		// Adicionar customer se não existir
		if _, exists := customerMap[customer.CustomerID]; !exists {
			customers = append(customers, *customer)
			customerMap[customer.CustomerID] = customer
		}

		// Adicionar product se não existir
		if _, exists := productMap[product.ProductID]; !exists {
			products = append(products, *product)
			productMap[product.ProductID] = product
		}

		// Adicionar usage
		usages = append(usages, *usage)

		// Processar lote quando atingir o tamanho
		if len(usages) >= batchSize {
			if err := svc.ProcessImportData(context.Background(), partners, customers, products, usages); err != nil {
				return fmt.Errorf("erro ao processar lote: %w", err)
			}

			processedCount += len(usages)
			log.Printf("Processado lote: %d registros", len(usages))

			// Limpar lotes
			partners = partners[:0]
			customers = customers[:0]
			products = products[:0]
			usages = usages[:0]
		}
	}

	// Processar último lote
	if len(usages) > 0 {
		if err := svc.ProcessImportData(context.Background(), partners, customers, products, usages); err != nil {
			return fmt.Errorf("erro ao processar último lote: %w", err)
		}
		processedCount += len(usages)
		log.Printf("Processado último lote: %d registros", len(usages))
	}

	log.Printf("Total processado: %d registros de %d linhas", processedCount, rowCount)
	return nil
}

func allEmpty(record []string) bool {
	for _, cell := range record {
		if strings.TrimSpace(cell) != "" {
			return false
		}
	}
	return true
}

func parseExcelRow(record []string, columnMap map[string]int, rowNum int) (*models.Partner, *models.Customer, *models.Product, *models.Usage, error) {
	// Função auxiliar para obter valor da coluna
	getValue := func(colName string) string {
		if idx, exists := columnMap[colName]; exists && idx < len(record) {
			return strings.TrimSpace(record[idx])
		}
		return ""
	}

	// Função auxiliar para converter string para float
	parseFloat := func(value string) (float64, error) {
		if value == "" {
			return 0, nil
		}
		// Remover vírgulas e converter
		value = strings.ReplaceAll(value, ",", ".")
		return strconv.ParseFloat(value, 64)
	}

	// Função auxiliar para converter string para data
	parseDate := func(value string) (time.Time, error) {
		if value == "" {
			return time.Time{}, nil
		}
		
		// Tentar diferentes formatos de data
		formats := []string{
			"2006-01-02",
			"2006/01/02",
			"02/01/2006",
			"02-01-2006",
			"02/01/06",   // ano curto
			"02-01-06",
			"1/2/2006",
			"1-2-2006",
			"2006-01-02 15:04:05",
			"2006/01/02 15:04:05",
			"2006-01-02T15:04:05Z",
			"2006-01-02T15:04:05.000Z",
		}
		
		
		for _, format := range formats {
			if t, err := time.Parse(format, value); err == nil {
				return t, nil
			}
		}
		
		// Tentar parsear como número serial do Excel
		if serial, err := strconv.ParseFloat(value, 64); err == nil {
			// Converter número serial do Excel para data
			// Excel usa 1900-01-01 como base, mas tem bug do ano bissexto
			baseDate := time.Date(1899, 12, 30, 0, 0, 0, 0, time.UTC)
			days := int(serial)
			return baseDate.AddDate(0, 0, days), nil
		}
		
		return time.Time{}, fmt.Errorf("formato de data inválido: %s", value)
	}

	// Criar Partner
	partner := &models.Partner{
		PartnerID:   getValue("partner_id"),
		PartnerName: getValue("partner_name"),
		MpnID:       getValue("mpn_id"),
		Tier2MpnID:  getValue("tier2_mpn_id"),
	}

	// Criar Customer
	customer := &models.Customer{
		CustomerID:         getValue("customer_id"),
		CustomerName:       getValue("customer_name"),
		CustomerDomainName: getValue("customer_domain_name"),
		Country:            getValue("country"),
	}

	// Criar Product
	product := &models.Product{
		ProductID:   getValue("product_id"),
		SkuID:       getValue("sku_id"),
		SkuName:     getValue("sku_name"),
		ProductName: getValue("product_name"),
		MeterType:   getValue("meter_type"),
		Category:    getValue("category"),
		SubCategory: getValue("sub_category"),
		UnitType:    getValue("unit_type"),
	}

	// Validar campos obrigatórios
	if partner.PartnerID == "" {
		return nil, nil, nil, nil, fmt.Errorf("partner_id é obrigatório")
	}
	if customer.CustomerID == "" {
		return nil, nil, nil, nil, fmt.Errorf("customer_id é obrigatório")
	}
	if product.ProductID == "" {
		return nil, nil, nil, nil, fmt.Errorf("product_id é obrigatório")
	}

	// Parsear datas
	usageDate, err := parseDate(getValue("usage_date"))
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("erro ao parsear usage_date: %w", err)
	}

	chargeStartDate, err := parseDate(getValue("charge_start_date"))
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("erro ao parsear charge_start_date: %w", err)
	}

	// Parsear valores numéricos
	quantity, err := parseFloat(getValue("quantity"))
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("erro ao parsear quantity: %w", err)
	}

	unitPrice, err := parseFloat(getValue("unit_price"))
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("erro ao parsear unit_price: %w", err)
	}

	billingPreTaxTotal, err := parseFloat(getValue("billing_pre_tax_total"))
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("erro ao parsear billing_pre_tax_total: %w", err)
	}

	// Criar Usage
	usage := &models.Usage{
		InvoiceNumber:      getValue("invoice_number"),
		ChargeStartDate:    timeToNullTime(chargeStartDate),
		UsageDate:          usageDate,
		Quantity:           quantity,
		UnitPrice:          unitPrice,
		BillingPreTaxTotal: billingPreTaxTotal,
		ResourceLocation:   getValue("resource_location"),
		Tags:               getValue("tags"),
		BenefitType:        getValue("benefit_type"),
		PartnerID:          0, // Será preenchido após inserção
		CustomerID:         0, // Será preenchido após inserção
		ProductID:          0, // Será preenchido após inserção
		// Campos temporários para resolução de IDs
		PartnerIDStr:       partner.PartnerID,
		CustomerIDStr:      customer.CustomerID,
		ProductIDStr:       product.ProductID,
	}

	return partner, customer, product, usage, nil
}

// timeToNullTime converte time.Time para sql.NullTime
func timeToNullTime(t time.Time) sql.NullTime {
	if t.IsZero() {
		return sql.NullTime{Valid: false}
	}
	return sql.NullTime{Time: t, Valid: true}
} 
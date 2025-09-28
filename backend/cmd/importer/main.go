package main

import (
	"context"
	"data-importer-api-go/internal/config"
	"data-importer-api-go/internal/models"
	"data-importer-api-go/internal/repository"
	"data-importer-api-go/internal/service"
	"database/sql"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Uso: go run ./cmd/importer/main.go <arquivo.csv>")
	}

	csvFile := os.Args[1]
	
	// Carregar configuração
	cfg := config.LoadConfig()

	// Conectar ao banco de dados
	db, err := pgxpool.New(context.Background(), cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Erro ao conectar ao banco de dados: %v", err)
	}
	defer db.Close()

	// Inicializar camadas
	repo := repository.NewRepository(db)
	svc := service.NewService(repo)

	// Processar arquivo CSV
	if err := processCSV(csvFile, svc); err != nil {
		log.Fatalf("Erro ao processar CSV: %v", err)
	}

	log.Println("✅ Importação concluída com sucesso!")
}

func processCSV(filename string, service *service.Service) error {
	// Abrir arquivo CSV
	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("erro ao abrir arquivo: %w", err)
	}
	defer file.Close()

	// Criar leitor CSV
	reader := csv.NewReader(file)
	reader.Comma = ','
	reader.LazyQuotes = true

	// Ler cabeçalho
	header, err := reader.Read()
	if err != nil {
		return fmt.Errorf("erro ao ler cabeçalho: %w", err)
	}

	log.Printf("📋 Cabeçalhos encontrados: %v", header)

	// Mapear índices das colunas
	columnMap := make(map[string]int)
	for i, col := range header {
		columnMap[strings.ToLower(strings.TrimSpace(col))] = i
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

	for {
		record, err := reader.Read()
		if err != nil {
			if err.Error() == "EOF" {
				break
			}
			log.Printf("⚠️  Erro ao ler linha %d: %v", rowCount+1, err)
			continue
		}

		rowCount++

		// Processar linha
		partner, customer, product, usage, err := parseRow(record, columnMap, rowCount)
		if err != nil {
			log.Printf("⚠️  Erro ao processar linha %d: %v", rowCount, err)
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
			if err := service.ProcessImportData(context.Background(), partners, customers, products, usages); err != nil {
				return fmt.Errorf("erro ao processar lote: %w", err)
			}

			processedCount += len(usages)
			log.Printf("📦 Processado lote: %d registros", len(usages))

			// Limpar lotes
			partners = partners[:0]
			customers = customers[:0]
			products = products[:0]
			usages = usages[:0]
		}
	}

	// Processar último lote
	if len(usages) > 0 {
		if err := service.ProcessImportData(context.Background(), partners, customers, products, usages); err != nil {
			return fmt.Errorf("erro ao processar último lote: %w", err)
		}
		processedCount += len(usages)
		log.Printf("📦 Processado último lote: %d registros", len(usages))
	}

	log.Printf("✅ Total processado: %d registros de %d linhas", processedCount, rowCount)
	return nil
}

func parseRow(record []string, columnMap map[string]int, rowNum int) (*models.Partner, *models.Customer, *models.Product, *models.Usage, error) {
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

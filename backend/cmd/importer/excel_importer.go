package main

import (
	"context"
	"data-importer-api-go/internal/config"
	"data-importer-api-go/internal/models"
	"data-importer-api-go/internal/repository"
	"data-importer-api-go/internal/service"
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/xuri/excelize/v2"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Uso: go run ./cmd/importer/excel_importer.go <arquivo.xlsx>")
	}

	excelFile := os.Args[1]
	
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

	// Processar arquivo Excel
	if err := processExcel(excelFile, svc); err != nil {
		log.Fatalf("Erro ao processar Excel: %v", err)
	}

	log.Println("Importação concluída com sucesso!")
}

func processExcel(filename string, service *service.Service) error {
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
			if err := service.ProcessImportData(context.Background(), partners, customers, products, usages); err != nil {
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
		if err := service.ProcessImportData(context.Background(), partners, customers, products, usages); err != nil {
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
		PartnerIDStr:       partner.PartnerID,    // Adicionado para mapeamento
		CustomerIDStr:      customer.CustomerID,  // Adicionado para mapeamento
		ProductIDStr:       product.ProductID,    // Adicionado para mapeamento
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

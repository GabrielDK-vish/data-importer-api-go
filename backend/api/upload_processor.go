package api

import (
	"data-importer-api-go/internal/models"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/xuri/excelize/v2"
)

// processExcelFile processa arquivo Excel
func (h *Handler) processExcelFile(file io.Reader) ([]models.Partner, []models.Customer, []models.Product, []models.Usage, error) {
	// Ler arquivo Excel
	f, err := excelize.OpenReader(file)
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("erro ao abrir arquivo Excel: %w", err)
	}
	defer f.Close()

	// Obter primeira planilha
	sheetList := f.GetSheetList()
	if len(sheetList) == 0 {
		return nil, nil, nil, nil, fmt.Errorf("arquivo Excel não possui planilhas")
	}

	sheetName := sheetList[0]
	rows, err := f.GetRows(sheetName)
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("erro ao ler planilha: %w", err)
	}

	if len(rows) < 2 {
		return nil, nil, nil, nil, fmt.Errorf("planilha deve ter pelo menos 2 linhas")
	}

	// Processar dados
	return h.processRows(rows)
}

// processCSVFile processa arquivo CSV
func (h *Handler) processCSVFile(file io.Reader) ([]models.Partner, []models.Customer, []models.Product, []models.Usage, error) {
	// Ler arquivo CSV
	reader := csv.NewReader(file)
	reader.Comma = ','
	reader.LazyQuotes = true

	// Ler todas as linhas
	var rows [][]string
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, nil, nil, nil, fmt.Errorf("erro ao ler CSV: %w", err)
		}
		rows = append(rows, record)
	}

	if len(rows) < 2 {
		return nil, nil, nil, nil, fmt.Errorf("CSV deve ter pelo menos 2 linhas")
	}

	// Processar dados
	return h.processRows(rows)
}

// processRows processa linhas de dados
func (h *Handler) processRows(rows [][]string) ([]models.Partner, []models.Customer, []models.Product, []models.Usage, error) {
	// Processar cabeçalho
	header := rows[0]
	columnMap := make(map[string]int)
	for i, col := range header {
		columnMap[strings.ToLower(strings.TrimSpace(col))] = i
	}

	// Verificar colunas obrigatórias
	requiredColumns := []string{"partner_id", "customer_id", "product_id", "usage_date", "quantity", "unit_price"}
	for _, col := range requiredColumns {
		if _, exists := columnMap[col]; !exists {
			return nil, nil, nil, nil, fmt.Errorf("coluna obrigatória não encontrada: %s", col)
		}
	}

	var partners []models.Partner
	var customers []models.Customer
	var products []models.Product
	var usages []models.Usage

	partnerMap := make(map[string]*models.Partner)
	customerMap := make(map[string]*models.Customer)
	productMap := make(map[string]*models.Product)

	// Processar linhas de dados
	for i := 1; i < len(rows); i++ {
		record := rows[i]
		
		// Pular linhas vazias
		if len(record) == 0 || h.allEmpty(record) {
			continue
		}

		// Processar linha
		partner, customer, product, usage, err := h.parseRow(record, columnMap, i)
		if err != nil {
			log.Printf("⚠️  Erro ao processar linha %d: %v", i, err)
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
	}

	return partners, customers, products, usages, nil
}

// allEmpty verifica se todas as células estão vazias
func (h *Handler) allEmpty(record []string) bool {
	for _, cell := range record {
		if strings.TrimSpace(cell) != "" {
			return false
		}
	}
	return true
}

// parseRow processa uma linha de dados
func (h *Handler) parseRow(record []string, columnMap map[string]int, rowNum int) (*models.Partner, *models.Customer, *models.Product, *models.Usage, error) {
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
			"2006-01-02 15:04:05",
			"2006/01/02 15:04:05",
		}
		
		for _, format := range formats {
			if t, err := time.Parse(format, value); err == nil {
				return t, nil
			}
		}
		
		// Tentar parsear como número serial do Excel
		if serial, err := strconv.ParseFloat(value, 64); err == nil {
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
		ChargeStartDate:    chargeStartDate,
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

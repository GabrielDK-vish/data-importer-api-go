package api

import (
	"data-importer-api-go/internal/models"
	"data-importer-api-go/internal/service"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/xuri/excelize/v2"
)

// UploadHandler lida com upload de arquivos
type UploadHandler struct {
	service *service.Service
}

func NewUploadHandler(service *service.Service) *UploadHandler {
	return &UploadHandler{service: service}
}

// UploadFileHandler processa upload de arquivos CSV/Excel
func (h *UploadHandler) UploadFileHandler(w http.ResponseWriter, r *http.Request) {
	// Verificar método
	if r.Method != http.MethodPost {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}

	// Parsear multipart form
	err := r.ParseMultipartForm(32 << 20) // 32MB max
	if err != nil {
		http.Error(w, "Erro ao parsear formulário", http.StatusBadRequest)
		return
	}

	// Obter arquivo
	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Erro ao obter arquivo", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Verificar tipo de arquivo
	contentType := header.Header.Get("Content-Type")
	fileName := header.Filename

	log.Printf("📁 Arquivo recebido: %s (%s)", fileName, contentType)
	log.Printf("📏 Tamanho do arquivo: %d bytes", header.Size)

	var partners []models.Partner
	var customers []models.Customer
	var products []models.Product
	var usages []models.Usage

	// Processar arquivo baseado na extensão
	if strings.HasSuffix(strings.ToLower(fileName), ".xlsx") {
		partners, customers, products, usages, err = h.processExcelFile(file)
	} else if strings.HasSuffix(strings.ToLower(fileName), ".csv") {
		partners, customers, products, usages, err = h.processCSVFile(file)
	} else {
		http.Error(w, "Tipo de arquivo não suportado. Use .csv ou .xlsx", http.StatusBadRequest)
		return
	}

	if err != nil {
		log.Printf("❌ Erro ao processar arquivo: %v", err)
		http.Error(w, fmt.Sprintf("Erro ao processar arquivo: %v", err), http.StatusInternalServerError)
		return
	}

	// Inserir dados no banco
	err = h.service.ProcessImportData(r.Context(), partners, customers, products, usages)
	if err != nil {
		log.Printf("❌ Erro ao inserir dados: %v", err)
		http.Error(w, fmt.Sprintf("Erro ao inserir dados: %v", err), http.StatusInternalServerError)
		return
	}

	// Resposta de sucesso
	response := map[string]interface{}{
		"success": true,
		"message": "Arquivo processado com sucesso",
		"data": map[string]interface{}{
			"partners":  len(partners),
			"customers": len(customers),
			"products":  len(products),
			"usages":    len(usages),
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *UploadHandler) processExcelFile(file io.Reader) ([]models.Partner, []models.Customer, []models.Product, []models.Usage, error) {
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

func (h *UploadHandler) processCSVFile(file io.Reader) ([]models.Partner, []models.Customer, []models.Product, []models.Usage, error) {
	// Ler todo o conteúdo para detectar delimitador
	data, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("erro ao ler arquivo: %w", err)
	}

	content := string(data)
	delimiter := ','
	// Heurística simples: se houver mais tabs do que vírgulas na primeira linha, usar tab
	firstLine := content
	if idx := strings.IndexAny(content, "\r\n"); idx != -1 {
		firstLine = content[:idx]
	}
	if strings.Count(firstLine, "\t") > strings.Count(firstLine, ",") {
		delimiter = '\t'
	}

	parseWith := func(delim rune) ([][]string, error) {
		reader := csv.NewReader(strings.NewReader(content))
		reader.Comma = delim
		reader.LazyQuotes = true
		var rows [][]string
		for {
			record, err := reader.Read()
			if err == io.EOF {
				break
			}
			if err != nil {
				return nil, fmt.Errorf("erro ao ler CSV: %w", err)
			}
			rows = append(rows, record)
		}
		return rows, nil
	}

	rows, err := parseWith(delimiter)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	if len(rows) < 2 {
		return nil, nil, nil, nil, fmt.Errorf("CSV deve ter pelo menos 2 linhas")
	}

	return h.processRows(rows)
}

func (h *UploadHandler) processRows(rows [][]string) ([]models.Partner, []models.Customer, []models.Product, []models.Usage, error) {
	// Processar cabeçalho
	header := rows[0]
	columnMap := make(map[string]int)

	// Normalizar cabeçalhos e aplicar aliases
	normalize := func(s string) string {
		s = strings.ToLower(strings.TrimSpace(s))
		s = strings.ReplaceAll(s, " ", "")
		s = strings.ReplaceAll(s, "_", "")
		s = strings.ReplaceAll(s, "-", "")
		return s
	}
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
		"unit":                   "unit_type",
		"unittype":               "unit_type",
		"resourcelocation":       "resource_location",
		"category":               "category",
		"subcategory":            "sub_category",
		
		// Usage fields
		"invoicenumber":          "invoice_number",
		"usagedate":              "usage_date",
		"chargestartdate":        "charge_start_date",
		"unitprice":              "unit_price",
		"effectiveunitprice":     "unit_price",
		"quantity":               "quantity",
		"billingpretaxtotal":     "billing_pre_tax_total",
		"billingcurrency":        "billing_currency",
		"pricingpretaxtotal":     "pricing_pre_tax_total",
		"pricingcurrency":        "pricing_currency",
		"benefittype":            "benefit_type",
		"tags":                   "tags",
		"additionalinfo":         "additional_info",
		"serviceinfo1":           "service_info1",
		"serviceinfo2":           "service_info2",
		"pcbcexchangerate":       "pc_to_bc_exchange_rate",
		"pcbcexchangeratedate":   "pc_to_bc_exchange_rate_date",
		"entitlementid":          "entitlement_id",
		"entitlementdescription": "entitlement_description",
		"partnerearnedcreditpercentage": "partner_earned_credit_percentage",
		"creditpercentage":       "credit_percentage",
		"credittype":             "credit_type",
		"benefitorderid":         "benefit_order_id",
		"benefitid":              "benefit_id",
	}

	for i, col := range header {
		n := normalize(col)
		key := n
		if mapped, ok := alias[n]; ok {
			key = mapped
		}
		columnMap[key] = i
	}

	// Verificar colunas obrigatórias
	requiredColumns := []string{"partner_id", "customer_id", "product_id", "usage_date", "quantity", "unit_price"}
	missingColumns := []string{}
	for _, col := range requiredColumns {
		if _, exists := columnMap[col]; !exists {
			missingColumns = append(missingColumns, col)
		}
	}
	
	if len(missingColumns) > 0 {
    log.Printf("⚠️  Colunas obrigatórias não encontradas: %v", missingColumns)
		log.Printf("📋 Colunas disponíveis: %v", getAvailableColumns(header))
		return nil, nil, nil, nil, fmt.Errorf("colunas obrigatórias não encontradas: %v. Colunas disponíveis: %v", missingColumns, getAvailableColumns(header))
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

func (h *UploadHandler) allEmpty(record []string) bool {
	for _, cell := range record {
		if strings.TrimSpace(cell) != "" {
			return false
		}
	}
	return true
}

func getAvailableColumns(header []string) []string {
    columns := make([]string, len(header))
    for i, col := range header {
        columns[i] = strings.ToLower(strings.TrimSpace(col))
    }
	return columns
}

func (h *UploadHandler) parseRow(record []string, columnMap map[string]int, rowNum int) (*models.Partner, *models.Customer, *models.Product, *models.Usage, error) {
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
            "1/2/2006",
            "1-2-2006",
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
		PartnerID:          0, 
		CustomerID:         0, // Será preenchido após inserçã
		ProductID:          0, 
	}

	return partner, customer, product, usage, nil
}

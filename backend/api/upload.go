package api

import (
	"data-importer-api-go/internal/models"
	"data-importer-api-go/internal/service"
	"database/sql"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
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
	log.Printf("📊 Processando planilha: %s", sheetName)
	
	rows, err := f.GetRows(sheetName)
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("erro ao ler planilha: %w", err)
	}

	log.Printf("📏 Total de linhas encontradas: %d", len(rows))

	if len(rows) < 2 {
		return nil, nil, nil, nil, fmt.Errorf("planilha deve ter pelo menos 2 linhas")
	}

	// Processar dados
	log.Printf("📊 Iniciando processamento de %d linhas de dados", len(rows))
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
	log.Printf("📋 Cabeçalhos encontrados: %v", header)
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
		log.Printf("🔗 Coluna %d: '%s' -> '%s'", i, col, key)
	}

	// Verificar colunas obrigatórias com mapeamento flexível
	requiredColumns := []string{"partner_id", "customer_id", "product_id", "usage_date", "quantity", "unit_price"}
	missingColumns := []string{}
	
	// Mapear colunas disponíveis para colunas obrigatórias
	columnMapping := make(map[string]string)
	for _, required := range requiredColumns {
		found := false
		for available, _ := range columnMap {
			if strings.Contains(strings.ToLower(available), strings.ToLower(required)) {
				columnMapping[required] = available
				found = true
				break
			}
		}
		if !found {
			missingColumns = append(missingColumns, required)
		}
	}
	
	if len(missingColumns) > 0 {
		log.Printf("⚠️  Colunas obrigatórias não encontradas: %v", missingColumns)
		log.Printf("📋 Colunas disponíveis: %v", getAvailableColumns(header))
		log.Printf("🔍 Tentando mapeamento automático...")
		
		// Tentar mapeamento automático mais agressivo
		for _, missing := range missingColumns {
			for available, _ := range columnMap {
				availableLower := strings.ToLower(strings.ReplaceAll(available, " ", ""))
				missingLower := strings.ToLower(strings.ReplaceAll(missing, "_", ""))
				
				if strings.Contains(availableLower, missingLower) || 
				   strings.Contains(missingLower, availableLower) ||
				   (strings.Contains(availableLower, "partner") && strings.Contains(missingLower, "partner")) ||
				   (strings.Contains(availableLower, "customer") && strings.Contains(missingLower, "customer")) ||
				   (strings.Contains(availableLower, "product") && strings.Contains(missingLower, "product")) ||
				   (strings.Contains(availableLower, "usage") && strings.Contains(missingLower, "usage")) ||
				   (strings.Contains(availableLower, "date") && strings.Contains(missingLower, "date")) ||
				   (strings.Contains(availableLower, "quantity") && strings.Contains(missingLower, "quantity")) ||
				   (strings.Contains(availableLower, "price") && strings.Contains(missingLower, "price")) {
					columnMapping[missing] = available
					log.Printf("✅ Mapeamento automático: '%s' -> '%s'", available, missing)
					break
				}
			}
		}
		
		// Verificar se ainda há colunas faltando
		stillMissing := []string{}
		for _, required := range requiredColumns {
			if _, exists := columnMapping[required]; !exists {
				stillMissing = append(stillMissing, required)
			}
		}
		
		if len(stillMissing) > 0 {
			return nil, nil, nil, nil, fmt.Errorf("colunas obrigatórias não encontradas: %v. Colunas disponíveis: %v", stillMissing, getAvailableColumns(header))
		}
	}
	
	// Atualizar columnMap com mapeamentos encontrados
	for required, available := range columnMapping {
		columnMap[required] = columnMap[available]
	}

	// Estruturas para processamento paralelo
	type rowResult struct {
		partner  *models.Partner
		customer *models.Customer
		product  *models.Product
		usage    *models.Usage
		err      error
		rowNum   int
	}

	// Canal para resultados
	resultChan := make(chan rowResult, len(rows)-1)
	
	// Worker pool para processamento paralelo
	numWorkers := 10 // Ajustar conforme necessário
	if numWorkers > len(rows)-1 {
		numWorkers = len(rows) - 1
	}
	
	// Canal para distribuir trabalho
	workChan := make(chan int, len(rows)-1)
	
	// Iniciar workers
	var wg sync.WaitGroup
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for rowNum := range workChan {
				record := rows[rowNum]
				
				// Pular linhas vazias
				if len(record) == 0 || h.allEmpty(record) {
					continue
				}

				// Processar linha
				partner, customer, product, usage, err := h.parseRow(record, columnMap, rowNum)
				resultChan <- rowResult{
					partner:  partner,
					customer: customer,
					product:  product,
					usage:    usage,
					err:      err,
					rowNum:   rowNum,
				}
			}
		}()
	}
	
	// Distribuir trabalho
	go func() {
		for i := 1; i < len(rows); i++ {
			workChan <- i
		}
		close(workChan)
	}()
	
	// Aguardar workers terminarem
	go func() {
		wg.Wait()
		close(resultChan)
	}()

	// Coletar resultados
	var partners []models.Partner
	var customers []models.Customer
	var products []models.Product
	var usages []models.Usage

	partnerMap := make(map[string]*models.Partner)
	customerMap := make(map[string]*models.Customer)
	productMap := make(map[string]*models.Product)

	processedCount := 0
	errorCount := 0

	for result := range resultChan {
		if result.err != nil {
			log.Printf("⚠️  Erro ao processar linha %d: %v", result.rowNum+1, result.err)
			errorCount++
			// Continuar processamento mesmo com erros
			continue
		}

		// Adicionar partner se não existir
		if _, exists := partnerMap[result.partner.PartnerID]; !exists {
			partners = append(partners, *result.partner)
			partnerMap[result.partner.PartnerID] = result.partner
		}

		// Adicionar customer se não existir
		if _, exists := customerMap[result.customer.CustomerID]; !exists {
			customers = append(customers, *result.customer)
			customerMap[result.customer.CustomerID] = result.customer
		}

		// Adicionar product se não existir
		if _, exists := productMap[result.product.ProductID]; !exists {
			products = append(products, *result.product)
			productMap[result.product.ProductID] = result.product
		}

		// Adicionar usage
		usages = append(usages, *result.usage)
		processedCount++
	}

	log.Printf("📊 Processamento concluído: %d linhas processadas, %d erros", processedCount, errorCount)

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
		
		// Limpar valor
		value = strings.TrimSpace(value)
		
		// Remover caracteres não numéricos exceto ponto, vírgula e sinal
		cleaned := ""
		for _, char := range value {
			if char >= '0' && char <= '9' || char == '.' || char == ',' || char == '-' || char == '+' {
				cleaned += string(char)
			}
		}
		
		if cleaned == "" {
			return 0, nil
		}
		
		// Substituir vírgula por ponto
		cleaned = strings.ReplaceAll(cleaned, ",", ".")
		
		// Verificar se há múltiplos pontos (formato brasileiro)
		if strings.Count(cleaned, ".") > 1 {
			// Se há múltiplos pontos, o último é o decimal
			parts := strings.Split(cleaned, ".")
			if len(parts) > 2 {
				cleaned = strings.Join(parts[:len(parts)-1], "") + "." + parts[len(parts)-1]
			}
		}
		
		return strconv.ParseFloat(cleaned, 64)
	}

	// Função auxiliar para converter string para data
    parseDate := func(value string) (time.Time, error) {
		if value == "" {
			return time.Time{}, nil
		}
		
		// Limpar valor
		value = strings.TrimSpace(value)
		
		// Tentar diferentes formatos de data
        formats := []string{
			"2006-01-02",
			"2006/01/02",
			"02/01/2006",
			"02-01-2006",
			"02/01/06",   // ano curto com /
			"02-01-06",   // ano curto com -
			"1/2/2006",
			"1-2-2006",
			"1/2/06",     // ano curto, sem zero à esquerda
			"1-2-06",     // ano curto, sem zero à esquerda
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
			// Verificar se é um número serial do Excel (geralmente entre 1 e 100000)
			if serial > 1 && serial < 100000 {
				baseDate := time.Date(1899, 12, 30, 0, 0, 0, 0, time.UTC)
				days := int(serial)
				return baseDate.AddDate(0, 0, days), nil
			}
		}
		
		// Se não conseguir parsear, retornar data atual como fallback
		log.Printf("⚠️  Não foi possível parsear data '%s', usando data atual", value)
		return time.Now(), nil
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

	// Parsear datas com tratamento de erro mais flexível
	usageDate, err := parseDate(getValue("usage_date"))
	if err != nil {
		log.Printf("⚠️  Erro ao parsear usage_date na linha %d: %v, usando data atual", rowNum+1, err)
		usageDate = time.Now()
	}

	chargeStartDate, err := parseDate(getValue("charge_start_date"))
	if err != nil {
		log.Printf("⚠️  Erro ao parsear charge_start_date na linha %d: %v, usando data atual", rowNum+1, err)
		chargeStartDate = time.Now()
	}

	// Parsear valores numéricos com tratamento de erro mais flexível
	quantity, err := parseFloat(getValue("quantity"))
	if err != nil {
		log.Printf("⚠️  Erro ao parsear quantity na linha %d: %v, usando 0", rowNum+1, err)
		quantity = 0
	}

	unitPrice, err := parseFloat(getValue("unit_price"))
	if err != nil {
		log.Printf("⚠️  Erro ao parsear unit_price na linha %d: %v, usando 0", rowNum+1, err)
		unitPrice = 0
	}

	billingPreTaxTotal, err := parseFloat(getValue("billing_pre_tax_total"))
	if err != nil {
		log.Printf("⚠️  Erro ao parsear billing_pre_tax_total na linha %d: %v, usando 0", rowNum+1, err)
		billingPreTaxTotal = 0
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
		PartnerIDStr:       partner.PartnerID,    // Adicionado para mapeamento
		CustomerIDStr:      customer.CustomerID,  // Adicionado para mapeamento
		ProductIDStr:       product.ProductID,    // Adicionado para mapeamento
		PartnerID:          0, 
		CustomerID:         0, // Será preenchido após inserção
		ProductID:          0, 
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

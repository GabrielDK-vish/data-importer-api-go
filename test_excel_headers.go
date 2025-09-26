package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/xuri/excelize/v2"
)

func main() {
	// Abrir arquivo Excel
	f, err := excelize.OpenFile("Reconfile fornecedores.xlsx")
	if err != nil {
		log.Fatalf("Erro ao abrir arquivo: %v", err)
	}
	defer f.Close()

	// Obter primeira planilha
	sheetList := f.GetSheetList()
	if len(sheetList) == 0 {
		log.Fatal("Arquivo Excel não possui planilhas")
	}

	sheetName := sheetList[0]
	rows, err := f.GetRows(sheetName)
	if err != nil {
		log.Fatalf("Erro ao ler planilha: %v", err)
	}

	if len(rows) < 1 {
		log.Fatal("Planilha vazia")
	}

	// Mostrar cabeçalhos
	header := rows[0]
	fmt.Printf("Total de colunas: %d\n", len(header))
	fmt.Println("\nCabeçalhos encontrados:")
	for i, col := range header {
		fmt.Printf("%d: %s\n", i, col)
	}

	// Normalizar e verificar mapeamento
	normalize := func(s string) string {
		s = strings.ToLower(strings.TrimSpace(s))
		s = strings.ReplaceAll(s, " ", "")
		s = strings.ReplaceAll(s, "_", "")
		s = strings.ReplaceAll(s, "-", "")
		return s
	}

	alias := map[string]string{
		"partnerid":              "partner_id",
		"partnername":            "partner_name",
		"mpnid":                  "mpn_id",
		"tier2mpnid":             "tier2_mpn_id",
		"customerid":             "customer_id",
		"customername":           "customer_name",
		"customerdomainname":     "customer_domain_name",
		"customercountry":        "country",
		"productid":              "product_id",
		"skuid":                  "sku_id",
		"skuname":                "sku_name",
		"productname":            "product_name",
		"metertype":              "meter_type",
		"metercategory":          "category",
		"metersubcategory":       "sub_category",
		"unit":                   "unit_type",
		"resourcelocation":       "resource_location",
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

	fmt.Println("\nMapeamento de colunas:")
	columnMap := make(map[string]int)
	for i, col := range header {
		n := normalize(col)
		key := n
		if mapped, ok := alias[n]; ok {
			key = mapped
		}
		columnMap[key] = i
		fmt.Printf("%s -> %s (posição %d)\n", col, key, i)
	}

	// Verificar colunas obrigatórias
	requiredColumns := []string{"partner_id", "customer_id", "product_id", "usage_date", "quantity", "unit_price"}
	fmt.Println("\nVerificação de colunas obrigatórias:")
	missingColumns := []string{}
	for _, col := range requiredColumns {
		if _, exists := columnMap[col]; !exists {
			missingColumns = append(missingColumns, col)
			fmt.Printf("❌ FALTANDO: %s\n", col)
		} else {
			fmt.Printf("✅ ENCONTRADA: %s (posição %d)\n", col, columnMap[col])
		}
	}

	if len(missingColumns) > 0 {
		fmt.Printf("\n❌ Colunas obrigatórias não encontradas: %v\n", missingColumns)
	} else {
		fmt.Println("\n✅ Todas as colunas obrigatórias foram encontradas!")
	}
}

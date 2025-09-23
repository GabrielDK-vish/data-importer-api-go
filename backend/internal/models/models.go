package models
import "time"

type Partner struct {
	ID         int    `json:"id"`
	PartnerId  string `json:"partner_id"`
	Name       string `json:"name"`
	MpnId      string `json:"mpn_id"`
	Tier2MpnId string `json:"tier2_mpn_id"`
}

type Customer struct {
	ID       int    `json:"id"`
	CustomerId string `json:"customer_id"`
	Name     string `json:"name"`
	Domain   string `json:"domain"`
	Country  string `json:"country"`
}

type Product struct {
	ID        int    `json:"id"`
	ProductId string `json:"product_id"`
	SkuId     string `json:"sku_id"`
	SkuName   string `json:"sku_name"`
	Name      string `json:"name"`
	MeterType string `json:"meter_type"`
	Category  string `json:"category"`
	SubCategory string `json:"sub_category"`
	UnitType  string `json:"unit_type"`
}

type Usage struct {
	ID               int       `json:"id"`
	InvoiceNumber    string    `json:"invoice_number"`
	PartnerID        int       `json:"partner_id"`
	CustomerID       int       `json:"customer_id"`
	ProductID        int       `json:"product_id"`
	ChargeStartDate  time.Time `json:"charge_start_date"`
	UsageDate        time.Time `json:"usage_date"`
	Quantity         float64   `json:"quantity"`
	UnitPrice        float64   `json:"unit_price"`
	BillingPreTaxTotal float64 `json:"billing_pre_tax_total"`
	ResourceLocation string    `json:"resource_location"`
	Tags             string    `json:"tags"`
	BenefitType      string    `json:"benefit_type"`
}

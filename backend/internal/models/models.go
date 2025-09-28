package models

import (
	"database/sql"
	"time"
)

// KPIData representa os dados de KPI do sistema
type KPIData struct {
	TotalRecords      int       `json:"totalRecords"`
	TotalCategories   int       `json:"totalCategories"`
	TotalResources    int       `json:"totalResources"`
	TotalCustomers    int       `json:"totalCustomers"`
	AvgBillingPerMonth float64  `json:"avgBillingPerMonth"`
	ProcessingTimeMs  int64     `json:"processingTimeMs"`
	LastUpdated       *time.Time `json:"lastUpdated"`
}

// CategoryBillingReport representa o faturamento por categoria
type CategoryBillingReport struct {
	Category    string  `json:"category"`
	Total       float64 `json:"total"`
	Count       int     `json:"count"`
}

// ResourceBillingReport representa o faturamento por recurso
type ResourceBillingReport struct {
	Resource    string  `json:"resource"`
	Total       float64 `json:"total"`
	Count       int     `json:"count"`
}

// CustomerBillingReport representa o faturamento por cliente
type CustomerBillingReport struct {
	CustomerID   int     `json:"customer_id"`
	CustomerName string  `json:"customer_name"`
	Total        float64 `json:"total"`
	Count        int     `json:"count"`
}

// Partner representa um parceiro
type Partner struct {
	ID          int       `json:"id" db:"id"`
	PartnerID   string    `json:"partner_id" db:"partner_id"`
	PartnerName string    `json:"partner_name" db:"partner_name"`
	MpnID       string    `json:"mpn_id" db:"mpn_id"`
	Tier2MpnID  string    `json:"tier2_mpn_id" db:"tier2_mpn_id"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

// Customer representa um cliente
type Customer struct {
	ID                 int       `json:"id" db:"id"`
	CustomerID         string    `json:"customer_id" db:"customer_id"`
	CustomerName       string    `json:"customer_name" db:"customer_name"`
	CustomerDomainName string    `json:"customer_domain_name" db:"customer_domain_name"`
	Country            string    `json:"country" db:"country"`
	CreatedAt          time.Time `json:"created_at" db:"created_at"`
	UpdatedAt          time.Time `json:"updated_at" db:"updated_at"`
}

// Product representa um produto/serviço
type Product struct {
	ID           int       `json:"id" db:"id"`
	ProductID    string    `json:"product_id" db:"product_id"`
	SkuID        string    `json:"sku_id" db:"sku_id"`
	SkuName      string    `json:"sku_name" db:"sku_name"`
	ProductName  string    `json:"product_name" db:"product_name"`
	MeterType    string    `json:"meter_type" db:"meter_type"`
	Category     string    `json:"category" db:"category"`
	SubCategory  string    `json:"sub_category" db:"sub_category"`
	UnitType     string    `json:"unit_type" db:"unit_type"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}

// Usage representa um registro de uso/faturamento
type Usage struct {
	ID                   int            `json:"id" db:"id"`
	InvoiceNumber        string         `json:"invoice_number" db:"invoice_number"`
	ChargeStartDate      sql.NullTime   `json:"charge_start_date" db:"charge_start_date"`
	UsageDate            time.Time      `json:"usage_date" db:"usage_date"`
	Quantity             float64        `json:"quantity" db:"quantity"`
	UnitPrice            float64        `json:"unit_price" db:"unit_price"`
	BillingPreTaxTotal   float64        `json:"billing_pre_tax_total" db:"billing_pre_tax_total"`
	ResourceLocation     string         `json:"resource_location" db:"resource_location"`
	Tags                 string         `json:"tags" db:"tags"`
	BenefitType          string         `json:"benefit_type" db:"benefit_type"`
	PartnerID            int            `json:"partner_id" db:"partner_id"`
	CustomerID           int            `json:"customer_id" db:"customer_id"`
	ProductID            int            `json:"product_id" db:"product_id"`
	CreatedAt            time.Time      `json:"created_at" db:"created_at"`
	UpdatedAt            time.Time      `json:"updated_at" db:"updated_at"`
	
	// Campos temporários para processamento (não são persistidos)
	PartnerIDStr         string         `json:"-" db:"-"`
	CustomerIDStr        string         `json:"-" db:"-"`
	ProductIDStr         string         `json:"-" db:"-"`
	
	// Relacionamentos
	Partner  *Partner  `json:"partner,omitempty"`
	Customer *Customer `json:"customer,omitempty"`
	Product  *Product  `json:"product,omitempty"`
}

// BillingReport representa relatórios de faturamento
type BillingReport struct {
	Month     string  `json:"month" db:"month"`
	Total     float64 `json:"total" db:"total"`
	Count     int     `json:"count" db:"count"`
}

type BillingByProduct struct {
	ProductID   string  `json:"product_id" db:"product_id"`
	ProductName string  `json:"product_name" db:"product_name"`
	Category    string  `json:"category" db:"category"`
	Total       float64 `json:"total" db:"total"`
	Count       int     `json:"count" db:"count"`
}

type BillingByPartner struct {
	PartnerID   string  `json:"partner_id" db:"partner_id"`
	PartnerName string  `json:"partner_name" db:"partner_name"`
	Total       float64 `json:"total" db:"total"`
	Count       int     `json:"count" db:"count"`
}

// LoginRequest representa a requisição de login
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// User representa um usuário do sistema
type User struct {
	ID           int       `json:"id" db:"id"`
	Username     string    `json:"username" db:"username"`
	PasswordHash string    `json:"-" db:"password_hash"`
	Email        string    `json:"email" db:"email"`
	FullName     string    `json:"full_name" db:"full_name"`
	IsActive     bool      `json:"is_active" db:"is_active"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}

// LoginResponse representa a resposta de login
type LoginResponse struct {
	Token string `json:"token"`
	User  string `json:"user"`
    // Opcional: métricas de processamento de importação executada no login
    ProcessingTimeMs int64 `json:"processing_time_ms,omitempty"`
    Partners         int   `json:"partners,omitempty"`
    Customers        int   `json:"customers,omitempty"`
    Products         int   `json:"products,omitempty"`
    Usages           int   `json:"usages,omitempty"`
}
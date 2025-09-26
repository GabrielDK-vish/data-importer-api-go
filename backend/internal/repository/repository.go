package repository

import (
	"context"
	"data-importer-api-go/internal/models"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

// GetAllCustomers retorna todos os clientes
func (r *Repository) GetAllCustomers(ctx context.Context) ([]models.Customer, error) {
	query := `
		SELECT id, customer_id, customer_name, customer_domain_name, country, created_at, updated_at
		FROM customers
		ORDER BY customer_name
	`
	
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar clientes: %w", err)
	}
	defer rows.Close()

	var customers []models.Customer
	for rows.Next() {
		var customer models.Customer
		err := rows.Scan(
			&customer.ID,
			&customer.CustomerID,
			&customer.CustomerName,
			&customer.CustomerDomainName,
			&customer.Country,
			&customer.CreatedAt,
			&customer.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("erro ao escanear cliente: %w", err)
		}
		customers = append(customers, customer)
	}

	return customers, nil
}

// GetUsageByCustomer retorna o uso de um cliente específico
func (r *Repository) GetUsageByCustomer(ctx context.Context, customerID int) ([]models.Usage, error) {
	query := `
		SELECT u.id, u.invoice_number, u.charge_start_date, u.usage_date, u.quantity, 
		       u.unit_price, u.billing_pre_tax_total, u.resource_location, u.tags, 
		       u.benefit_type, u.partner_id, u.customer_id, u.product_id, u.created_at, u.updated_at,
		       p.partner_id, p.partner_name,
		       c.customer_id, c.customer_name,
		       pr.product_id, pr.product_name, pr.category
		FROM usages u
		LEFT JOIN partners p ON u.partner_id = p.id
		LEFT JOIN customers c ON u.customer_id = c.id
		LEFT JOIN products pr ON u.product_id = pr.id
		WHERE u.customer_id = $1
		ORDER BY u.usage_date DESC
	`
	
	rows, err := r.db.Query(ctx, query, customerID)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar uso do cliente: %w", err)
	}
	defer rows.Close()

	var usages []models.Usage
	for rows.Next() {
		var usage models.Usage
		var partner models.Partner
		var customer models.Customer
		var product models.Product
		
		err := rows.Scan(
			&usage.ID, &usage.InvoiceNumber, &usage.ChargeStartDate, &usage.UsageDate,
			&usage.Quantity, &usage.UnitPrice, &usage.BillingPreTaxTotal,
			&usage.ResourceLocation, &usage.Tags, &usage.BenefitType,
			&usage.PartnerID, &usage.CustomerID, &usage.ProductID,
			&usage.CreatedAt, &usage.UpdatedAt,
			&partner.PartnerID, &partner.PartnerName,
			&customer.CustomerID, &customer.CustomerName,
			&product.ProductID, &product.ProductName, &product.Category,
		)
		if err != nil {
			return nil, fmt.Errorf("erro ao escanear uso: %w", err)
		}
		
		usage.Partner = &partner
		usage.Customer = &customer
		usage.Product = &product
		usages = append(usages, usage)
	}

	return usages, nil
}

// GetBillingMonthly retorna faturamento por mês
func (r *Repository) GetBillingMonthly(ctx context.Context) ([]models.BillingReport, error) {
	query := `
		SELECT 
			TO_CHAR(usage_date, 'YYYY-MM') as month,
			SUM(billing_pre_tax_total) as total,
			COUNT(*) as count
		FROM usages
		GROUP BY TO_CHAR(usage_date, 'YYYY-MM')
		ORDER BY month DESC
	`
	
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar faturamento mensal: %w", err)
	}
	defer rows.Close()

	var reports []models.BillingReport
	for rows.Next() {
		var report models.BillingReport
		err := rows.Scan(&report.Month, &report.Total, &report.Count)
		if err != nil {
			return nil, fmt.Errorf("erro ao escanear relatório mensal: %w", err)
		}
		reports = append(reports, report)
	}

	return reports, nil
}

// GetBillingByProduct retorna faturamento por produto
func (r *Repository) GetBillingByProduct(ctx context.Context) ([]models.BillingByProduct, error) {
	query := `
		SELECT 
			pr.product_id,
			pr.product_name,
			pr.category,
			SUM(u.billing_pre_tax_total) as total,
			COUNT(*) as count
		FROM usages u
		JOIN products pr ON u.product_id = pr.id
		GROUP BY pr.product_id, pr.product_name, pr.category
		ORDER BY total DESC
	`
	
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar faturamento por produto: %w", err)
	}
	defer rows.Close()

	var reports []models.BillingByProduct
	for rows.Next() {
		var report models.BillingByProduct
		err := rows.Scan(&report.ProductID, &report.ProductName, &report.Category, &report.Total, &report.Count)
		if err != nil {
			return nil, fmt.Errorf("erro ao escanear relatório por produto: %w", err)
		}
		reports = append(reports, report)
	}

	return reports, nil
}

// GetBillingByPartner retorna faturamento por parceiro
func (r *Repository) GetBillingByPartner(ctx context.Context) ([]models.BillingByPartner, error) {
	query := `
		SELECT 
			p.partner_id,
			p.partner_name,
			SUM(u.billing_pre_tax_total) as total,
			COUNT(*) as count
		FROM usages u
		JOIN partners p ON u.partner_id = p.id
		GROUP BY p.partner_id, p.partner_name
		ORDER BY total DESC
	`
	
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar faturamento por parceiro: %w", err)
	}
	defer rows.Close()

	var reports []models.BillingByPartner
	for rows.Next() {
		var report models.BillingByPartner
		err := rows.Scan(&report.PartnerID, &report.PartnerName, &report.Total, &report.Count)
		if err != nil {
			return nil, fmt.Errorf("erro ao escanear relatório por parceiro: %w", err)
		}
		reports = append(reports, report)
	}

	return reports, nil
}

// InsertPartner insere um novo parceiro
func (r *Repository) InsertPartner(ctx context.Context, partner *models.Partner) error {
	query := `
		INSERT INTO partners (partner_id, partner_name, mpn_id, tier2_mpn_id)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (partner_id) DO UPDATE SET
			partner_name = EXCLUDED.partner_name,
			mpn_id = EXCLUDED.mpn_id,
			tier2_mpn_id = EXCLUDED.tier2_mpn_id,
			updated_at = CURRENT_TIMESTAMP
		RETURNING id
	`
	
	err := r.db.QueryRow(ctx, query, partner.PartnerID, partner.PartnerName, partner.MpnID, partner.Tier2MpnID).Scan(&partner.ID)
	if err != nil {
		return fmt.Errorf("erro ao inserir parceiro: %w", err)
	}
	
	return nil
}

// InsertCustomer insere um novo cliente
func (r *Repository) InsertCustomer(ctx context.Context, customer *models.Customer) error {
	query := `
		INSERT INTO customers (customer_id, customer_name, customer_domain_name, country)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (customer_id) DO UPDATE SET
			customer_name = EXCLUDED.customer_name,
			customer_domain_name = EXCLUDED.customer_domain_name,
			country = EXCLUDED.country,
			updated_at = CURRENT_TIMESTAMP
		RETURNING id
	`
	
	err := r.db.QueryRow(ctx, query, customer.CustomerID, customer.CustomerName, customer.CustomerDomainName, customer.Country).Scan(&customer.ID)
	if err != nil {
		return fmt.Errorf("erro ao inserir cliente: %w", err)
	}
	
	return nil
}

// InsertProduct insere um novo produto
func (r *Repository) InsertProduct(ctx context.Context, product *models.Product) error {
	query := `
		INSERT INTO products (product_id, sku_id, sku_name, product_name, meter_type, category, sub_category, unit_type)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		ON CONFLICT (product_id) DO UPDATE SET
			sku_id = EXCLUDED.sku_id,
			sku_name = EXCLUDED.sku_name,
			product_name = EXCLUDED.product_name,
			meter_type = EXCLUDED.meter_type,
			category = EXCLUDED.category,
			sub_category = EXCLUDED.sub_category,
			unit_type = EXCLUDED.unit_type,
			updated_at = CURRENT_TIMESTAMP
		RETURNING id
	`
	
	err := r.db.QueryRow(ctx, query, product.ProductID, product.SkuID, product.SkuName, product.ProductName, product.MeterType, product.Category, product.SubCategory, product.UnitType).Scan(&product.ID)
	if err != nil {
		return fmt.Errorf("erro ao inserir produto: %w", err)
	}
	
	return nil
}

// InsertUsage insere um novo registro de uso
func (r *Repository) InsertUsage(ctx context.Context, usage *models.Usage) error {
	query := `
		INSERT INTO usages (invoice_number, charge_start_date, usage_date, quantity, unit_price, 
		                   billing_pre_tax_total, resource_location, tags, benefit_type, 
		                   partner_id, customer_id, product_id)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
		RETURNING id
	`
	
	err := r.db.QueryRow(ctx, query, usage.InvoiceNumber, usage.ChargeStartDate, usage.UsageDate, 
		usage.Quantity, usage.UnitPrice, usage.BillingPreTaxTotal, usage.ResourceLocation, 
		usage.Tags, usage.BenefitType, usage.PartnerID, usage.CustomerID, usage.ProductID).Scan(&usage.ID)
	if err != nil {
		return fmt.Errorf("erro ao inserir uso: %w", err)
	}
	
	return nil
}

// BulkInsertUsages insere múltiplos registros de uso usando CopyFrom para performance
func (r *Repository) BulkInsertUsages(ctx context.Context, usages []models.Usage) error {
	if len(usages) == 0 {
		return nil
	}

	// Preparar dados para CopyFrom
	rows := make([][]interface{}, len(usages))
	for i, usage := range usages {
		rows[i] = []interface{}{
			usage.InvoiceNumber,
			usage.ChargeStartDate,
			usage.UsageDate,
			usage.Quantity,
			usage.UnitPrice,
			usage.BillingPreTaxTotal,
			usage.ResourceLocation,
			usage.Tags,
			usage.BenefitType,
			usage.PartnerID,
			usage.CustomerID,
			usage.ProductID,
		}
	}

	// Usar CopyFrom para inserção em lote
	_, err := r.db.CopyFrom(ctx, pgx.Identifier{"usages"}, 
		[]string{"invoice_number", "charge_start_date", "usage_date", "quantity", 
			"unit_price", "billing_pre_tax_total", "resource_location", "tags", 
			"benefit_type", "partner_id", "customer_id", "product_id"}, 
		pgx.CopyFromRows(rows))
	
	if err != nil {
		return fmt.Errorf("erro ao inserir usos em lote: %w", err)
	}

	return nil
}

// ClearAllData limpa todos os dados das tabelas
func (r *Repository) ClearAllData(ctx context.Context) error {
	// Limpar dados na ordem correta (respeitando foreign keys)
	queries := []string{
		"DELETE FROM usages",
		"DELETE FROM products", 
		"DELETE FROM customers",
		"DELETE FROM partners",
	}
	
	for _, query := range queries {
		_, err := r.db.Exec(ctx, query)
		if err != nil {
			return fmt.Errorf("erro ao executar query %s: %w", query, err)
		}
	}
	
	return nil
}
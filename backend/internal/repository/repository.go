package repository

import (
	"context"
	"data-importer-api-go/internal/models"
	"fmt"
	"time"

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

// parseUsageDate tenta converter uma string para time.Time
func parseUsageDate(s string) (time.Time, error) {
    layouts := []string{
        "02-01-06",    // DD-MM-YY
        "01-02-06",    // MM-DD-YY
        "02/01/06",    // DD/MM/YY
        "01/02/06",    // MM/DD/YY
        "2006-01-02",  // YYYY-MM-DD
        "01/02/2006",  // MM/DD/YYYY
        "02/01/2006",  // DD/MM/YYYY
    }

    for _, layout := range layouts {
        if t, err := time.Parse(layout, s); err == nil {
            return t, nil
        }
    }

    return time.Time{}, fmt.Errorf("formato de data inválido: %s", s)
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

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("erro após iteração de faturamento mensal: %w", err)
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

// GetBillingByCategory retorna o faturamento por categoria
func (r *Repository) GetBillingByCategory(ctx context.Context) ([]models.CategoryBillingReport, error) {
	query := `
		SELECT 
			pr.category,
			SUM(u.billing_pre_tax_total) as total,
			COUNT(*) as count
		FROM usages u
		JOIN products pr ON u.product_id = pr.id
		GROUP BY pr.category
		ORDER BY total DESC
	`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("erro ao consultar faturamento por categoria: %w", err)
	}
	defer rows.Close()

	var reports []models.CategoryBillingReport
	for rows.Next() {
		var report models.CategoryBillingReport
		if err := rows.Scan(&report.Category, &report.Total, &report.Count); err != nil {
			return nil, fmt.Errorf("erro ao escanear resultado de faturamento por categoria: %w", err)
		}
		reports = append(reports, report)
	}

	return reports, nil
}

// GetBillingByResource retorna o faturamento por recurso
func (r *Repository) GetBillingByResource(ctx context.Context) ([]models.ResourceBillingReport, error) {
	query := `
		SELECT 
			u.resource_location as resource,
			SUM(u.billing_pre_tax_total) as total,
			COUNT(*) as count
		FROM usages u
		GROUP BY u.resource_location
		ORDER BY total DESC
	`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("erro ao consultar faturamento por recurso: %w", err)
	}
	defer rows.Close()

	var reports []models.ResourceBillingReport
	for rows.Next() {
		var report models.ResourceBillingReport
		if err := rows.Scan(&report.Resource, &report.Total, &report.Count); err != nil {
			return nil, fmt.Errorf("erro ao escanear resultado de faturamento por recurso: %w", err)
		}
		reports = append(reports, report)
	}

	return reports, nil
}

// GetBillingByCustomer retorna o faturamento por cliente
func (r *Repository) GetBillingByCustomer(ctx context.Context) ([]models.CustomerBillingReport, error) {
	query := `
		SELECT 
			c.customer_id,
			c.customer_name,
			SUM(u.billing_pre_tax_total) as total,
			COUNT(*) as count
		FROM usages u
		JOIN customers c ON u.customer_id = c.id
		GROUP BY c.customer_id, c.customer_name
		ORDER BY total DESC
	`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("erro ao consultar faturamento por cliente: %w", err)
	}
	defer rows.Close()

	var reports []models.CustomerBillingReport
	for rows.Next() {
		var report models.CustomerBillingReport
		if err := rows.Scan(&report.CustomerID, &report.CustomerName, &report.Total, &report.Count); err != nil {
			return nil, fmt.Errorf("erro ao escanear resultado de faturamento por cliente: %w", err)
		}
		reports = append(reports, report)
	}

	return reports, nil
}

// GetKPIData retorna os dados de KPI do sistema
func (r *Repository) GetKPIData(ctx context.Context) (*models.KPIData, error) {
	query := `
		WITH stats AS (
			SELECT 
				COUNT(DISTINCT u.id) as total_records,
				COUNT(DISTINCT pr.category) as total_categories,
				COUNT(DISTINCT u.resource_location) as total_resources,
				COUNT(DISTINCT c.id) as total_customers,
				AVG(monthly.total) as avg_billing_per_month,
				MAX(u.updated_at) as last_updated
			FROM 
				usages u
			JOIN 
				products pr ON u.product_id = pr.id
			JOIN 
				customers c ON u.customer_id = c.id
			LEFT JOIN (
				SELECT 
					DATE_TRUNC('month', usage_date) as month,
					SUM(billing_pre_tax_total) as total
				FROM 
					usages
				GROUP BY 
					DATE_TRUNC('month', usage_date)
			) monthly ON TRUE
		)
		SELECT 
			total_records,
			total_categories,
			total_resources,
			total_customers,
			COALESCE(avg_billing_per_month, 0) as avg_billing_per_month,
			last_updated
		FROM 
			stats
	`

	var kpiData models.KPIData
	var lastUpdated time.Time

	err := r.db.QueryRow(ctx, query).Scan(
		&kpiData.TotalRecords,
		&kpiData.TotalCategories,
		&kpiData.TotalResources,
		&kpiData.TotalCustomers,
		&kpiData.AvgBillingPerMonth,
		&lastUpdated,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("erro ao consultar dados de KPI: %w", err)
	}

	kpiData.LastUpdated = &lastUpdated

	return &kpiData, nil
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

func (r *Repository) InsertUsage(ctx context.Context, usage *models.Usage) error {
    query := `
        INSERT INTO usages (invoice_number, charge_start_date, usage_date, quantity, unit_price, 
                            billing_pre_tax_total, resource_location, tags, benefit_type, 
                            partner_id, customer_id, product_id)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
        RETURNING id
    `
    
    err := r.db.QueryRow(ctx, query, 
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
    ).Scan(&usage.ID)
    
    if err != nil {
        return fmt.Errorf("erro ao inserir uso: %w", err)
    }
    
    return nil
}


// Bufunc (r *Repository) BulkInsertUsages(ctx context.Context, usages []models.Usage) error {
	func (r *Repository) BulkInsertUsages(ctx context.Context, usages []models.Usage) error {
		if len(usages) == 0 {
			return nil
		}
	
		rows := make([][]interface{}, len(usages))
	
		for i, usage := range usages {
			var chargeStartDate, usageDate interface{}
				if usage.ChargeStartDate.Valid {
					chargeStartDate = usage.ChargeStartDate.Time
				} else {
					chargeStartDate = nil
				}

				if !usage.UsageDate.IsZero() {  // usage.UsageDate é time.Time
					usageDate = usage.UsageDate
				} else {
					usageDate = nil
				}

	
			rows[i] = []interface{}{
				usage.InvoiceNumber,
				chargeStartDate,
				usageDate,
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
	
		_, err := r.db.CopyFrom(
			ctx,
			pgx.Identifier{"usages"},
			[]string{
				"invoice_number",
				"charge_start_date",
				"usage_date",
				"quantity",
				"unit_price",
				"billing_pre_tax_total",
				"resource_location",
				"tags",
				"benefit_type",
				"partner_id",
				"customer_id",
				"product_id",
			},
			pgx.CopyFromRows(rows),
		)
	
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

// GetUserByUsername busca um usuário pelo username
func (r *Repository) GetUserByUsername(ctx context.Context, username string) (*models.User, error) {
	query := `
		SELECT id, username, password_hash, email, full_name, is_active, created_at, updated_at
		FROM users
		WHERE username = $1 AND is_active = true
	`
	
	var user models.User
	err := r.db.QueryRow(ctx, query, username).Scan(
		&user.ID,
		&user.Username,
		&user.PasswordHash,
		&user.Email,
		&user.FullName,
		&user.IsActive,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil // Usuário não encontrado
		}
		return nil, fmt.Errorf("erro ao buscar usuário: %w", err)
	}
	
	return &user, nil
}

// BulkInsertPartners insere múltiplos registros de parceiros usando CopyFrom para performance
func (r *Repository) BulkInsertPartners(ctx context.Context, partners []models.Partner) error {
	if len(partners) == 0 {
		return nil
	}

	// Preparar dados para CopyFrom
	rows := make([][]interface{}, len(partners))
	for i, partner := range partners {
		rows[i] = []interface{}{
			partner.PartnerID,
			partner.PartnerName,
			partner.MpnID,
			partner.Tier2MpnID,
		}
	}

	// Usar CopyFrom para inserção em lote
	_, err := r.db.CopyFrom(ctx, pgx.Identifier{"partners"}, 
		[]string{"partner_id", "partner_name", "mpn_id", "tier2_mpn_id"}, 
		pgx.CopyFromRows(rows))
	
	if err != nil {
		return fmt.Errorf("erro ao inserir parceiros em lote: %w", err)
	}

	return nil
}

// BulkInsertCustomers insere múltiplos registros de clientes usando CopyFrom para performance
func (r *Repository) BulkInsertCustomers(ctx context.Context, customers []models.Customer) error {
	if len(customers) == 0 {
		return nil
	}

	// Preparar dados para CopyFrom
	rows := make([][]interface{}, len(customers))
	for i, customer := range customers {
		rows[i] = []interface{}{
			customer.CustomerID,
			customer.CustomerName,
			customer.CustomerDomainName,
			customer.Country,
		}
	}

	// Usar CopyFrom para inserção em lote
	_, err := r.db.CopyFrom(ctx, pgx.Identifier{"customers"}, 
		[]string{"customer_id", "customer_name", "customer_domain_name", "country"}, 
		pgx.CopyFromRows(rows))
	
	if err != nil {
		return fmt.Errorf("erro ao inserir clientes em lote: %w", err)
	}

	return nil
}

// BulkInsertProducts insere múltiplos registros de produtos usando CopyFrom para performance
func (r *Repository) BulkInsertProducts(ctx context.Context, products []models.Product) error {
	if len(products) == 0 {
		return nil
	}

	rows := make([][]interface{}, len(products))
	for i, product := range products {
		rows[i] = []interface{}{
			product.ProductID,
			product.SkuID,
			product.ProductName,
			product.SkuName,
			product.MeterType,
			product.Category,
			product.SubCategory,
			product.UnitType,
		}
	}

	_, err := r.db.CopyFrom(
		ctx,
		pgx.Identifier{"products"},
		[]string{"product_id", "sku_id", "product_name", "sku_name", "meter_type", 
		         "category", "sub_category", "unit_type"},
		pgx.CopyFromRows(rows),
	)

	if err != nil {
		return fmt.Errorf("erro ao inserir produtos em lote: %w", err)
	}

	return nil
}


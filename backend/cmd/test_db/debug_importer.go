package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

// Usage representa a estrutura da tabela usages
type Usage struct {
	ID                 int
	InvoiceNumber      string
	ChargeStartDate    pgtype.Date
	UsageDate          time.Time
	Quantity           float64
	UnitPrice          float64
	BillingPreTaxTotal float64
	ResourceLocation   string
	Tags               string
	BenefitType        string
	PartnerID          int
	CustomerID         int
	ProductID          int
}

func main() {
	// String de conexão com o banco de dados
	dbURL := "postgresql://data_importer_db_user:Gx4hgHpOpFxY60QCyIBAmY6BlfULuktb@dpg-d3ar2n7fte5s7398mj3g-a.oregon-postgres.render.com/data_importer_db?sslmode=require"

	// Conectar ao banco de dados
	db, err := pgxpool.New(context.Background(), dbURL)
	if err != nil {
		log.Fatalf("Erro ao conectar ao banco de dados: %v", err)
	}
	defer db.Close()

	// Testar a conexão
	if err := db.Ping(context.Background()); err != nil {
		log.Fatalf("Erro ao fazer ping no banco de dados: %v", err)
	}
	fmt.Println("✅ Conexão com o banco de dados estabelecida com sucesso!")

	// Criar dados de teste
	usages := []Usage{
		{
			InvoiceNumber:      "TEST-002",
			UsageDate:          time.Now(),
			Quantity:           2.0,
			UnitPrice:          15.0,
			BillingPreTaxTotal: 30.0,
			ResourceLocation:   "East US",
			Tags:               "test",
			BenefitType:        "None",
			PartnerID:          1, // Assumindo que existe um parceiro com ID 1
			CustomerID:         1, // Assumindo que existe um cliente com ID 1
			ProductID:          1, // Assumindo que existe um produto com ID 1
		},
	}

	// Verificar se os IDs existem
	var partnerExists, customerExists, productExists bool
	
	err = db.QueryRow(context.Background(), "SELECT EXISTS(SELECT 1 FROM partners WHERE id = $1)", usages[0].PartnerID).Scan(&partnerExists)
	if err != nil {
		log.Fatalf("Erro ao verificar parceiro: %v", err)
	}
	
	err = db.QueryRow(context.Background(), "SELECT EXISTS(SELECT 1 FROM customers WHERE id = $1)", usages[0].CustomerID).Scan(&customerExists)
	if err != nil {
		log.Fatalf("Erro ao verificar cliente: %v", err)
	}
	
	err = db.QueryRow(context.Background(), "SELECT EXISTS(SELECT 1 FROM products WHERE id = $1)", usages[0].ProductID).Scan(&productExists)
	if err != nil {
		log.Fatalf("Erro ao verificar produto: %v", err)
	}
	
	fmt.Printf("Verificação de IDs: Partner ID %d existe: %v, Customer ID %d existe: %v, Product ID %d existe: %v\n", 
		usages[0].PartnerID, partnerExists, 
		usages[0].CustomerID, customerExists, 
		usages[0].ProductID, productExists)
	
	if !partnerExists || !customerExists || !productExists {
		fmt.Println("⚠️ Criando registros de teste para partner, customer e product...")
		
		// Criar partner se não existir
		if !partnerExists {
			_, err = db.Exec(context.Background(), 
				"INSERT INTO partners (partner_id, partner_name) VALUES ($1, $2) ON CONFLICT DO NOTHING", 
				"TEST-PARTNER", "Test Partner")
			if err != nil {
				log.Fatalf("Erro ao criar partner: %v", err)
			}
		}
		
		// Criar customer se não existir
		if !customerExists {
			_, err = db.Exec(context.Background(), 
				"INSERT INTO customers (customer_id, customer_name) VALUES ($1, $2) ON CONFLICT DO NOTHING", 
				"TEST-CUSTOMER", "Test Customer")
			if err != nil {
				log.Fatalf("Erro ao criar customer: %v", err)
			}
		}
		
		// Criar product se não existir
		if !productExists {
			_, err = db.Exec(context.Background(), 
				"INSERT INTO products (product_id, sku_id, sku_name, product_name) VALUES ($1, $2, $3, $4) ON CONFLICT DO NOTHING", 
				"TEST-PRODUCT", "TEST-SKU", "Test SKU", "Test Product")
			if err != nil {
				log.Fatalf("Erro ao criar product: %v", err)
			}
		}
		
		// Obter os IDs corretos
		err = db.QueryRow(context.Background(), "SELECT id FROM partners WHERE partner_id = $1", "TEST-PARTNER").Scan(&usages[0].PartnerID)
		if err != nil {
			log.Fatalf("Erro ao obter partner ID: %v", err)
		}
		
		err = db.QueryRow(context.Background(), "SELECT id FROM customers WHERE customer_id = $1", "TEST-CUSTOMER").Scan(&usages[0].CustomerID)
		if err != nil {
			log.Fatalf("Erro ao obter customer ID: %v", err)
		}
		
		err = db.QueryRow(context.Background(), "SELECT id FROM products WHERE product_id = $1", "TEST-PRODUCT").Scan(&usages[0].ProductID)
		if err != nil {
			log.Fatalf("Erro ao obter product ID: %v", err)
		}
		
		fmt.Printf("IDs atualizados: Partner ID: %d, Customer ID: %d, Product ID: %d\n", 
			usages[0].PartnerID, usages[0].CustomerID, usages[0].ProductID)
	}

	// Testar inserção em lote usando CopyFrom (similar ao BulkInsertUsages)
	rows := make([][]interface{}, len(usages))
	for i, usage := range usages {
		var chargeStartDate interface{}
		if usage.ChargeStartDate.Valid {
			chargeStartDate = usage.ChargeStartDate.Time
		} else {
			chargeStartDate = nil
		}

		rows[i] = []interface{}{
			usage.InvoiceNumber,
			chargeStartDate,
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

	// Usar transação para garantir consistência
	tx, err := db.Begin(context.Background())
	if err != nil {
		log.Fatalf("Erro ao iniciar transação: %v", err)
	}
	defer tx.Rollback(context.Background())

	fmt.Println("🔄 Iniciando inserção em lote...")
	
	// Imprimir os valores que serão inseridos
	fmt.Printf("Valores a serem inseridos: %+v\n", rows[0])
	
	_, err = tx.CopyFrom(
		context.Background(),
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
		log.Fatalf("Erro ao inserir usos em lote: %v", err)
	}

	// Commit da transação
	if err := tx.Commit(context.Background()); err != nil {
		log.Fatalf("Erro ao finalizar transação: %v", err)
	}

	fmt.Println("✅ Inserção em lote concluída com sucesso!")

	// Verificar se os registros foram inseridos
	var count int
	err = db.QueryRow(context.Background(), "SELECT COUNT(*) FROM usages WHERE invoice_number = $1", usages[0].InvoiceNumber).Scan(&count)
	if err != nil {
		log.Fatalf("Erro ao contar registros: %v", err)
	}
	fmt.Printf("📊 Registros inseridos com invoice_number '%s': %d\n", usages[0].InvoiceNumber, count)
}
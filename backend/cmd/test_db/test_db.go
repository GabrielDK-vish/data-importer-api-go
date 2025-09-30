package main

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

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

	// Verificar se a tabela usages existe
	var tableExists bool
	err = db.QueryRow(context.Background(), 
		"SELECT EXISTS (SELECT FROM information_schema.tables WHERE table_schema = 'public' AND table_name = 'usages')").Scan(&tableExists)
	if err != nil {
		log.Fatalf("Erro ao verificar se a tabela existe: %v", err)
	}

	if tableExists {
		fmt.Println("✅ Tabela usages existe no banco de dados")
	} else {
		fmt.Println("❌ Tabela usages não existe no banco de dados")
		return
	}

	// Verificar se há registros na tabela usages
	var count int
	err = db.QueryRow(context.Background(), "SELECT COUNT(*) FROM usages").Scan(&count)
	if err != nil {
		log.Fatalf("Erro ao contar registros na tabela usages: %v", err)
	}
	fmt.Printf("📊 Número de registros na tabela usages: %d\n", count)

	// Testar inserção manual na tabela usages
	_, err = db.Exec(context.Background(), `
		INSERT INTO usages (invoice_number, usage_date, quantity, unit_price, billing_pre_tax_total, partner_id, customer_id, product_id)
		VALUES ('TEST-001', CURRENT_DATE, 1.0, 10.0, 10.0, 
			(SELECT id FROM partners LIMIT 1), 
			(SELECT id FROM customers LIMIT 1), 
			(SELECT id FROM products LIMIT 1))
	`)
	if err != nil {
		log.Fatalf("Erro ao inserir registro de teste na tabela usages: %v", err)
	}
	fmt.Println("✅ Registro de teste inserido com sucesso na tabela usages")

	// Verificar novamente o número de registros
	err = db.QueryRow(context.Background(), "SELECT COUNT(*) FROM usages").Scan(&count)
	if err != nil {
		log.Fatalf("Erro ao contar registros na tabela usages após inserção: %v", err)
	}
	fmt.Printf("📊 Número de registros na tabela usages após inserção: %d\n", count)
}
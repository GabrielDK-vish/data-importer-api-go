package main

import (
	"context"
	"data-importer-api-go/api"
	"data-importer-api-go/internal/config"
	"data-importer-api-go/internal/repository"
	"data-importer-api-go/internal/service"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	// Carregar configuração
	cfg := config.LoadConfig()

	// Conectar ao banco de dados
	db, err := pgxpool.New(context.Background(), cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Erro ao conectar ao banco de dados: %v", err)
	}
	defer db.Close()

	// Executar migrations
	if err := runMigrations(cfg.DatabaseURL); err != nil {
		log.Fatalf("Erro ao executar migrations: %v", err)
	}

	// Inicializar camadas
	repo := repository.NewRepository(db)
	svc := service.NewService(repo)
	handler := api.NewHandler(svc)

	// Configurar rotas
	router := handler.SetupRoutes()

	// Configurar servidor
	server := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
	}

	// Iniciar servidor em goroutine
	go func() {
		log.Printf("🚀 Servidor iniciado na porta %s", cfg.Port)
		log.Printf("📊 Endpoints disponíveis:")
		log.Printf("   POST /auth/login")
		log.Printf("   GET  /api/customers")
		log.Printf("   GET  /api/customers/{id}/usage")
		log.Printf("   GET  /api/reports/billing/monthly")
		log.Printf("   GET  /api/reports/billing/by-product")
		log.Printf("   GET  /api/reports/billing/by-partner")
		
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Erro ao iniciar servidor: %v", err)
		}
	}()

	// Aguardar sinal de interrupção
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("🛑 Parando servidor...")

	// Graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Erro ao parar servidor: %v", err)
	}

	log.Println("✅ Servidor parado com sucesso")
}

func runMigrations(databaseURL string) error {
	m, err := migrate.New(
		"file://db/migrations",
		databaseURL,
	)
	if err != nil {
		return fmt.Errorf("erro ao criar migrator: %w", err)
	}
	defer m.Close()

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("erro ao executar migrations: %w", err)
	}

	log.Println("✅ Migrations executadas com sucesso")
	return nil
}
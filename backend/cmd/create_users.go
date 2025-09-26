package main

import (
	"context"
	"data-importer-api-go/internal/config"
	"data-importer-api-go/internal/repository"
	"fmt"
	"log"

	"golang.org/x/crypto/bcrypt"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Script para criar usuários com senhas hasheadas
func main() {
	// Carregar configuração
	cfg := config.LoadConfig()

	// Conectar ao banco de dados
	db, err := pgxpool.New(context.Background(), cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Erro ao conectar ao banco de dados: %v", err)
	}
	defer db.Close()

	// Criar repositório
	repo := repository.NewRepository(db)

	// Criar usuários com senhas hasheadas
	users := []struct {
		username string
		password string
		email    string
		fullName string
	}{
		{"admin", "admin123", "admin@example.com", "Administrator"},
		{"user", "user123", "user@example.com", "Regular User"},
		{"demo", "demo123", "demo@example.com", "Demo User"},
	}

	for _, u := range users {
		// Hash da senha
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.password), bcrypt.DefaultCost)
		if err != nil {
			log.Printf("Erro ao gerar hash para %s: %v", u.username, err)
			continue
		}

		// Inserir usuário
		query := `
			INSERT INTO users (username, password_hash, email, full_name, is_active)
			VALUES ($1, $2, $3, $4, true)
			ON CONFLICT (username) DO UPDATE SET
				password_hash = EXCLUDED.password_hash,
				email = EXCLUDED.email,
				full_name = EXCLUDED.full_name,
				updated_at = CURRENT_TIMESTAMP
		`

		_, err = db.Exec(context.Background(), query, u.username, string(hashedPassword), u.email, u.fullName)
		if err != nil {
			log.Printf("Erro ao inserir usuário %s: %v", u.username, err)
		} else {
			log.Printf("uário %s criado/atualizado com sucesso", u.username)
		}
	}

	log.Println("Processo de criação de usuários concluído!")
}

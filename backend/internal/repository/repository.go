package repository
import (
	"context"
	"log"
	"time"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	DB *pgxpool.Pool
}

func NewRepository(dbUrl string) *Repository {
	pool, err := pgxpool.New(context.Background(), dbUrl)
	if err != nil {
		log.Fatalf("Erro de conexão com Postgres: %v", err)
	}
	
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := pool.Ping(ctx); err != nil {
		log.Fatalf("Conexão intermitente: %v", err)
	}

	return &Repository{DB: pool}
}

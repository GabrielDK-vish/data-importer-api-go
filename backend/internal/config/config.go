package config
import (
	"log"
	"os"
)

type Config struct {
	DBUrl   string
	JWTKey  string
}

func LoadConfig() *Config {
	dbUrl := os.Getenv("DATABASE_URL")
	if dbUrl == "" {
		log.Fatal("Erro DATABASE")
	}

	jwtKey := os.Getenv("JWT_SECRET")
	if jwtKey == "" {
		log.Fatal("Erro JWT")
	}
	
	return &Config{
		DBUrl:  dbUrl,
		JWTKey: jwtKey,
	}
}

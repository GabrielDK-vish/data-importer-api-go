package main
import (
	"log"
	"net/http"
	"os"
	"github.com/yourusername/data-importer-api-go/internal/config"
	"github.com/yourusername/data-importer-api-go/internal/auth"
	"github.com/yourusername/data-importer-api-go/api"
)

func main() {
	cfg := config.LoadConfig()
	auth.InitJWT(cfg.JWTKey)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8050"
	}

	r := api.Routes()
	log.Printf("Endereço de servidor http://localhost:%s", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}

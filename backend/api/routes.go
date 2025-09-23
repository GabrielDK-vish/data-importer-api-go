package api
import (
	"net/http"
	"github.com/go-chi/chi/v5"
	"github.com/yourusername/data-importer-api-go/internal/auth"
)

func Routes() *chi.Mux {
	r := chi.NewRouter()

	// autenticação
	r.Post("/auth/login", LoginHandler)

	// rotas estabelecidas com segurança
	r.Group(func(r chi.Router) {
		r.Use(auth.VerifyToken)

		r.Get("/customers", GetCustomersHandler)
		r.Get("/customers/{id}/usage", GetCustomerUsageHandler)
		r.Get("/reports/billing/monthly", MonthlyBillingHandler)
		r.Get("/reports/billing/by-product", BillingByProductHandler)
		r.Get("/reports/billing/by-partner", BillingByPartnerHandler)
	})

	return r
}

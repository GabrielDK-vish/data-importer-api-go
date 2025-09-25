package api

import (
	"context"
	"data-importer-api-go/internal/auth"
	"data-importer-api-go/internal/models"
	"data-importer-api-go/internal/service"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

type Handler struct {
	service *service.Service
}

func NewHandler(service *service.Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) SetupRoutes() *chi.Mux {
	r := chi.NewRouter()

	// Middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	// Rotas públicas
	r.Post("/auth/login", h.LoginHandler)
	r.Get("/health", h.HealthCheckHandler)

	// Rotas protegidas
	r.Route("/api", func(r chi.Router) {
		r.Use(h.AuthMiddleware)
		
		// Clientes
		r.Get("/customers", h.GetCustomersHandler)
		r.Get("/customers/{id}/usage", h.GetCustomerUsageHandler)
		
		// Relatórios
		r.Get("/reports/billing/monthly", h.MonthlyBillingHandler)
		r.Get("/reports/billing/by-product", h.BillingByProductHandler)
		r.Get("/reports/billing/by-partner", h.BillingByPartnerHandler)
		
		// Upload de arquivos
		r.Post("/upload", h.UploadFileHandler)
	})

	return r
}

// LoginHandler autentica o usuário e retorna um token JWT
func (h *Handler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	var loginReq models.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&loginReq); err != nil {
		http.Error(w, "Erro ao decodificar JSON", http.StatusBadRequest)
		return
	}

	// Validar credenciais
	if !auth.ValidateCredentials(loginReq.Username, loginReq.Password) {
		http.Error(w, "Credenciais inválidas", http.StatusUnauthorized)
		return
	}

	// Gerar token
	token, err := auth.GenerateToken(loginReq.Username)
	if err != nil {
		http.Error(w, "Erro ao gerar token", http.StatusInternalServerError)
		return
	}

	response := models.LoginResponse{
		Token: token,
		User:  loginReq.Username,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// AuthMiddleware valida o token JWT
func (h *Handler) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Token de autorização necessário", http.StatusUnauthorized)
			return
		}

		// Extrair token do header "Bearer <token>"
		tokenString := authHeader
		if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
			tokenString = authHeader[7:]
		}

		// Validar token
		claims, err := auth.ValidateToken(tokenString)
		if err != nil {
			http.Error(w, "Token inválido", http.StatusUnauthorized)
			return
		}

		// Adicionar username ao contexto
		ctx := context.WithValue(r.Context(), "username", claims.Username)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// GetCustomersHandler retorna todos os clientes
func (h *Handler) GetCustomersHandler(w http.ResponseWriter, r *http.Request) {
	customers, err := h.service.GetAllCustomers(r.Context())
	if err != nil {
		http.Error(w, fmt.Sprintf("Erro ao buscar clientes: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(customers)
}

// GetCustomerUsageHandler retorna o uso de um cliente específico
func (h *Handler) GetCustomerUsageHandler(w http.ResponseWriter, r *http.Request) {
	customerIDStr := chi.URLParam(r, "id")
	customerID, err := strconv.Atoi(customerIDStr)
	if err != nil {
		http.Error(w, "ID do cliente inválido", http.StatusBadRequest)
		return
	}

	usages, err := h.service.GetUsageByCustomer(r.Context(), customerID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Erro ao buscar uso do cliente: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(usages)
}

// MonthlyBillingHandler retorna faturamento por mês
func (h *Handler) MonthlyBillingHandler(w http.ResponseWriter, r *http.Request) {
	reports, err := h.service.GetBillingMonthly(r.Context())
	if err != nil {
		http.Error(w, fmt.Sprintf("Erro ao buscar faturamento mensal: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(reports)
}

// BillingByProductHandler retorna faturamento por produto
func (h *Handler) BillingByProductHandler(w http.ResponseWriter, r *http.Request) {
	reports, err := h.service.GetBillingByProduct(r.Context())
	if err != nil {
		http.Error(w, fmt.Sprintf("Erro ao buscar faturamento por produto: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(reports)
}

// BillingByPartnerHandler retorna faturamento por parceiro
func (h *Handler) BillingByPartnerHandler(w http.ResponseWriter, r *http.Request) {
	reports, err := h.service.GetBillingByPartner(r.Context())
	if err != nil {
		http.Error(w, fmt.Sprintf("Erro ao buscar faturamento por parceiro: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(reports)
}

// UploadFileHandler processa upload de arquivos CSV/Excel
func (h *Handler) UploadFileHandler(w http.ResponseWriter, r *http.Request) {
	// Verificar método
	if r.Method != http.MethodPost {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}

	// Parsear multipart form
	err := r.ParseMultipartForm(32 << 20) // 32MB max
	if err != nil {
		http.Error(w, "Erro ao parsear formulário", http.StatusBadRequest)
		return
	}

	// Obter arquivo
	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Erro ao obter arquivo", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Verificar tipo de arquivo
	fileName := header.Filename
	log.Printf("📁 Arquivo recebido: %s", fileName)

	var partners []models.Partner
	var customers []models.Customer
	var products []models.Product
	var usages []models.Usage

	// Processar arquivo baseado na extensão
	if strings.HasSuffix(strings.ToLower(fileName), ".xlsx") {
		partners, customers, products, usages, err = h.processExcelFile(file)
	} else if strings.HasSuffix(strings.ToLower(fileName), ".csv") {
		partners, customers, products, usages, err = h.processCSVFile(file)
	} else {
		http.Error(w, "Tipo de arquivo não suportado. Use .csv ou .xlsx", http.StatusBadRequest)
		return
	}

	if err != nil {
		log.Printf("❌ Erro ao processar arquivo: %v", err)
		http.Error(w, fmt.Sprintf("Erro ao processar arquivo: %v", err), http.StatusInternalServerError)
		return
	}

	// Inserir dados no banco
	err = h.service.ProcessImportData(r.Context(), partners, customers, products, usages)
	if err != nil {
		log.Printf("❌ Erro ao inserir dados: %v", err)
		http.Error(w, fmt.Sprintf("Erro ao inserir dados: %v", err), http.StatusInternalServerError)
		return
	}

	// Resposta de sucesso
	response := map[string]interface{}{
		"success": true,
		"message": "Arquivo processado com sucesso",
		"data": map[string]interface{}{
			"partners":  len(partners),
			"customers": len(customers),
			"products":  len(products),
			"usages":    len(usages),
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// HealthCheckHandler retorna status de saúde da aplicação
func (h *Handler) HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"ok","service":"data-importer-api"}`))
}
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
	"time"

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
	r.Get("/", h.RootHandler)
	r.Post("/auth/login", h.LoginHandler)
	r.Get("/health", h.HealthCheckHandler)
	r.Get("/test", h.TestHandler)

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
		r.Get("/reports/billing/by-category", h.BillingByCategoryHandler)
		r.Get("/reports/billing/by-resource", h.BillingByResourceHandler)
		r.Get("/reports/billing/by-customer", h.BillingByCustomerHandler)
		r.Get("/reports/kpi", h.KPIHandler)
		
		// Upload
		r.Post("/upload", h.UploadFileHandler)
	})

	return r
}

// KPIHandler retorna as métricas de KPI do sistema
func (h *Handler) KPIHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	
	// Obter KPIs do banco de dados
	kpiData, err := h.service.GetKPIData(ctx)
	if err != nil {
		http.Error(w, "Erro ao obter KPIs: "+err.Error(), http.StatusInternalServerError)
		return
	}
	
	// Retornar KPIs como JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(kpiData)
}

// BillingByCategoryHandler retorna o faturamento agrupado por categoria
func (h *Handler) BillingByCategoryHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	
	// Obter dados de faturamento por categoria
	billingData, err := h.service.GetBillingByCategory(ctx)
	if err != nil {
		http.Error(w, "Erro ao obter dados de faturamento por categoria: "+err.Error(), http.StatusInternalServerError)
		return
	}
	
	// Retornar dados como JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(billingData)
}

// BillingByResourceHandler retorna o faturamento agrupado por recurso
func (h *Handler) BillingByResourceHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	
	// Obter dados de faturamento por recurso
	billingData, err := h.service.GetBillingByResource(ctx)
	if err != nil {
		http.Error(w, "Erro ao obter dados de faturamento por recurso: "+err.Error(), http.StatusInternalServerError)
		return
	}
	
	// Retornar dados como JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(billingData)
}

// BillingByCustomerHandler retorna o faturamento agrupado por cliente
func (h *Handler) BillingByCustomerHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	
	// Obter dados de faturamento por cliente
	billingData, err := h.service.GetBillingByCustomer(ctx)
	if err != nil {
		http.Error(w, "Erro ao obter dados de faturamento por cliente: "+err.Error(), http.StatusInternalServerError)
		return
	}
	
	// Retornar dados como JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(billingData)
}

// LoginHandler autentica o usuário e retorna um token JWT
func (h *Handler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	var loginReq models.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&loginReq); err != nil {
		http.Error(w, "Erro ao decodificar JSON", http.StatusBadRequest)
		return
	}

	// Validar credenciais no banco de dados
	user, err := h.service.ValidateUserCredentials(r.Context(), loginReq.Username, loginReq.Password)
	if err != nil {
		http.Error(w, "Credenciais inválidas", http.StatusUnauthorized)
		return
	}

	// Gerar token
	token, err := auth.GenerateToken(loginReq.Username)
	if err != nil {
		http.Error(w, "Erro ao gerar token", http.StatusInternalServerError)
		return
	}

	// Resposta simples de login sem reprocessamento
	response := models.LoginResponse{
		Token: token,
		User:  user.Username,
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
	log.Printf("Arquivo recebido: %s", fileName)

	var partners []models.Partner
	var customers []models.Customer
	var products []models.Product
	var usages []models.Usage

	// Processar arquivo baseado na extensão
	uploadHandler := NewUploadHandler(h.service)
	if strings.HasSuffix(strings.ToLower(fileName), ".xlsx") {
		partners, customers, products, usages, err = uploadHandler.processExcelFile(file)
	} else if strings.HasSuffix(strings.ToLower(fileName), ".csv") {
		partners, customers, products, usages, err = uploadHandler.processCSVFile(file)
	} else {
		http.Error(w, "Tipo de arquivo não suportado. Use .csv ou .xlsx", http.StatusBadRequest)
		return
	}

	if err != nil {
		log.Printf("Erro ao processar arquivo: %v", err)
		http.Error(w, fmt.Sprintf("Erro ao processar arquivo: %v", err), http.StatusInternalServerError)
		return
	}

	// Inserir dados no banco (substituindo dados existentes)
	log.Printf("Iniciando substituição de dados: %d partners, %d customers, %d products, %d usages", 
		len(partners), len(customers), len(products), len(usages))
	
	err = h.service.ProcessImportDataWithReplace(r.Context(), partners, customers, products, usages)
	if err != nil {
		log.Printf("Erro ao inserir dados: %v", err)
		http.Error(w, fmt.Sprintf("Erro ao inserir dados no banco: %v", err), http.StatusInternalServerError)
		return
	}
	
	log.Printf("Dados substituídos com sucesso")

	// Resposta de sucesso
	response := map[string]interface{}{
		"success": true,
		"message": "Arquivo processado e dados substituídos com sucesso",
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

// RootHandler retorna página inicial da API
func (h *Handler) RootHandler(w http.ResponseWriter, r *http.Request) {
	// Verificar se é uma requisição JSON
	if r.Header.Get("Accept") == "application/json" {
    apiInfo := map[string]interface{}{
			"service": "Data Importer API",
			"version": "1.0.0",
			"status":  "running",
			"endpoints": map[string]interface{}{
				"public": []string{
					"GET  /",
					"GET  /health",
					"POST /auth/login",
				},
				"protected": []string{
					"GET  /api/customers",
					"GET  /api/customers/{id}/usage",
					"GET  /api/reports/billing/monthly",
					"GET  /api/reports/billing/by-product",
					"GET  /api/reports/billing/by-partner",
				},
			},
			"documentation": "https://github.com/GabrielDK-vish/data-importer-api-go",
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(apiInfo)
		return
	}

	// Retornar página HTML
	html := `<!DOCTYPE html>
<html lang="pt-BR">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Data Importer API - Desafio Técnico</title>
    <style>
        * {
            margin: 0;
            padding: 0;
            box-sizing: border-box;
        }
        
        body {
            font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
            line-height: 1.6;
            color: #333;
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            min-height: 100vh;
        }
        
        .container {
            max-width: 1200px;
            margin: 0 auto;
            padding: 20px;
        }
        
        .header {
            text-align: center;
            color: white;
            margin-bottom: 40px;
        }
        
        .header h1 {
            font-size: 3rem;
            margin-bottom: 10px;
            text-shadow: 2px 2px 4px rgba(0,0,0,0.3);
        }
        
        .header p {
            font-size: 1.2rem;
            opacity: 0.9;
        }
        
        .card {
            background: white;
            border-radius: 15px;
            padding: 30px;
            margin-bottom: 30px;
            box-shadow: 0 10px 30px rgba(0,0,0,0.1);
            transition: transform 0.3s ease;
        }
        
        .card:hover {
            transform: translateY(-5px);
        }
        
        .card h2 {
            color: #667eea;
            margin-bottom: 20px;
            font-size: 1.8rem;
            border-bottom: 3px solid #667eea;
            padding-bottom: 10px;
        }
        
        .status {
            display: inline-block;
            background: #4CAF50;
            color: white;
            padding: 5px 15px;
            border-radius: 20px;
            font-size: 0.9rem;
            margin-bottom: 20px;
        }
        
        .endpoint {
            background: #f8f9fa;
            border-left: 4px solid #667eea;
            padding: 15px;
            margin: 10px 0;
            border-radius: 5px;
            font-family: 'Courier New', monospace;
        }
        
        .method {
            display: inline-block;
            padding: 3px 8px;
            border-radius: 3px;
            font-size: 0.8rem;
            font-weight: bold;
            margin-right: 10px;
        }
        
        .get { background: #4CAF50; color: white; }
        .post { background: #2196F3; color: white; }
        .put { background: #FF9800; color: white; }
        .delete { background: #f44336; color: white; }
        
        .credentials {
            background: #fff3cd;
            border: 1px solid #ffeaa7;
            border-radius: 8px;
            padding: 20px;
            margin: 20px 0;
        }
        
        .credentials h3 {
            color: #856404;
            margin-bottom: 15px;
        }
        
        .cred-table {
            width: 100%;
            border-collapse: collapse;
            margin-top: 10px;
        }
        
        .cred-table th, .cred-table td {
            padding: 10px;
            text-align: left;
            border-bottom: 1px solid #ddd;
        }
        
        .cred-table th {
            background: #f8f9fa;
            font-weight: bold;
        }
        
        .features {
            display: grid;
            grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
            gap: 20px;
            margin: 30px 0;
        }
        
        .feature {
            background: #f8f9fa;
            padding: 20px;
            border-radius: 10px;
            border-left: 4px solid #667eea;
        }
        
        .feature h3 {
            color: #667eea;
            margin-bottom: 10px;
        }
        
        .footer {
            text-align: center;
            color: white;
            margin-top: 40px;
            opacity: 0.8;
        }
        
        .tech-stack {
            display: flex;
            flex-wrap: wrap;
            gap: 10px;
            margin: 20px 0;
        }
        
        .tech {
            background: #667eea;
            color: white;
            padding: 5px 15px;
            border-radius: 20px;
            font-size: 0.9rem;
        }
        
        @media (max-width: 768px) {
            .header h1 {
                font-size: 2rem;
            }
            
            .container {
                padding: 10px;
            }
            
            .card {
                padding: 20px;
            }
        }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>Data Importer API</h1>
            <p>Desafio Técnico - Full Stack Developer</p>
            <span class="status">✅ Online</span>
        </div>
        
        <div class="card">
            <h2>📋 Sobre o Projeto</h2>
            <p>Esta é uma API desenvolvida em <strong>Go (Golang)</strong> para importação e análise de dados de faturamento. O projeto inclui:</p>
            
            <div class="features">
                <div class="feature">
                    <h3>Backend</h3>
                    <p>API REST em Go com PostgreSQL, autenticação JWT, e processamento de arquivos Excel/CSV.</p>
                </div>
                <div class="feature">
                    <h3>Frontend</h3>
                    <p>Interface React com dashboard, relatórios e upload de arquivos.</p>
                </div>
                <div class="feature">
                    <h3>Relatórios</h3>
                    <p>Análise de faturamento por mês, produto e parceiro com visualizações interativas.</p>
                </div>
                <div class="feature">
                    <h3>Segurança</h3>
                    <p>Autenticação JWT, validação de dados e tratamento de erros robusto.</p>
                </div>
            </div>
            
            <h3>Stack Tecnológica</h3>
            <div class="tech-stack">
                <span class="tech">Go (Golang)</span>
                <span class="tech">PostgreSQL</span>
                <span class="tech">React</span>
                <span class="tech">JWT</span>
                <span class="tech">Docker</span>
                <span class="tech">Chi Router</span>
                <span class="tech">Excelize</span>
                <span class="tech">Render</span>
                <span class="tech">Vercel</span>
            </div>
        </div>
        
        <div class="card">
            <h2>Acesso Rápido</h2>
            <div style="text-align: center; margin: 30px 0;">
                <a href="https://data-importer-api-go.vercel.app/" target="_blank" style="display: inline-block; background: linear-gradient(135deg, #667eea 0%, #764ba2 100%); color: white; padding: 15px 30px; text-decoration: none; border-radius: 25px; font-weight: bold; font-size: 1.1rem; margin: 10px; box-shadow: 0 4px 15px rgba(0,0,0,0.2); transition: transform 0.3s ease;">
                     Acessar Frontend
                </a>
                <a href="https://data-importer-api-go.onrender.com/health" target="_blank" style="display: inline-block; background: linear-gradient(135deg, #4CAF50 0%, #45a049 100%); color: white; padding: 15px 30px; text-decoration: none; border-radius: 25px; font-weight: bold; font-size: 1.1rem; margin: 10px; box-shadow: 0 4px 15px rgba(0,0,0,0.2); transition: transform 0.3s ease;">
                    🔍 Health Check
                </a>
            </div>
        </div>
        
        <div class="card">
            <h2>Credenciais de Teste</h2>
            <div class="credentials">
                <h3>Usuários Disponíveis</h3>
                <table class="cred-table">
                    <thead>
                        <tr>
                            <th>Usuário</th>
                            <th>Senha</th>
                            <th>Descrição</th>
                        </tr>
                    </thead>
                    <tbody>
                        <tr>
                            <td><code>admin</code></td>
                            <td><code>admin123</code></td>
                            <td>Administrador</td>
                        </tr>
                        <tr>
                            <td><code>user</code></td>
                            <td><code>user123</code></td>
                            <td>Usuário padrão</td>
                        </tr>
                        <tr>
                            <td><code>demo</code></td>
                            <td><code>demo123</code></td>
                            <td>Demonstração</td>
                        </tr>
                    </tbody>
                </table>
            </div>
        </div>
        
        <div class="card">
            <h2>Endpoints da API</h2>
            
            <h3>Endpoints Públicos</h3>
            <div class="endpoint">
                <span class="method get">GET</span> <strong>/</strong> - Página inicial da API
            </div>
            <div class="endpoint">
                <span class="method get">GET</span> <strong>/health</strong> - Status de saúde da aplicação
            </div>
            <div class="endpoint">
                <span class="method post">POST</span> <strong>/auth/login</strong> - Autenticação de usuário
            </div>
            
            <h3>Endpoints Protegidos (Requer Autenticação)</h3>
            <div class="endpoint">
                <span class="method get">GET</span> <strong>/api/customers</strong> - Listar todos os clientes
            </div>
            <div class="endpoint">
                <span class="method get">GET</span> <strong>/api/customers/{id}/usage</strong> - Uso detalhado por cliente
            </div>
            <div class="endpoint">
                <span class="method get">GET</span> <strong>/api/reports/billing/monthly</strong> - Faturamento por mês
            </div>
            <div class="endpoint">
                <span class="method get">GET</span> <strong>/api/reports/billing/by-product</strong> - Faturamento por produto
            </div>
            <div class="endpoint">
                <span class="method get">GET</span> <strong>/api/reports/billing/by-partner</strong> - Faturamento por parceiro
            </div>
            <div class="endpoint">
                <span class="method post">POST</span> <strong>/api/upload</strong> - Upload de arquivos Excel/CSV
            </div>
        </div>
        
        <div class="card">
            <h2>URLs de Produção</h2>
            <div style="display: grid; grid-template-columns: repeat(auto-fit, minmax(300px, 1fr)); gap: 20px; margin: 20px 0;">
                <div style="background: #e3f2fd; padding: 20px; border-radius: 10px; border-left: 4px solid #2196F3;">
                    <h3 style="color: #1976D2; margin-bottom: 10px;">Frontend</h3>
                    <p><strong>URL:</strong> <a href="https://data-importer-api-go.vercel.app/" target="_blank" style="color: #1976D2; text-decoration: none;">https://data-importer-api-go.vercel.app/</a></p>
                    <p><strong>Plataforma:</strong> Vercel</p>
                    <p><strong>Status:</strong> <span style="color: #4CAF50; font-weight: bold;"> Online</span></p>
                </div>
                <div style="background: #f3e5f5; padding: 20px; border-radius: 10px; border-left: 4px solid #9C27B0;">
                    <h3 style="color: #7B1FA2; margin-bottom: 10px;">🔧 Backend API</h3>
                    <p><strong>URL:</strong> <a href="https://data-importer-api-go.onrender.com/" target="_blank" style="color: #7B1FA2; text-decoration: none;">https://data-importer-api-go.onrender.com/</a></p>
                    <p><strong>Plataforma:</strong> Render</p>
                    <p><strong>Status:</strong> <span style="color: #4CAF50; font-weight: bold;"> Online</span></p>
                </div>
            </div>
        </div>
        
        <div class="card">
            <h2> Documentação e Recursos</h2>
            <div style="display: grid; grid-template-columns: repeat(auto-fit, minmax(250px, 1fr)); gap: 15px; margin: 20px 0;">
                <div style="background: #f8f9fa; padding: 15px; border-radius: 8px; border-left: 3px solid #667eea;">
                    <h4 style="color: #667eea; margin-bottom: 8px;">📖 Repositório</h4>
                    <p><a href="https://github.com/GabrielDK-vish/data-importer-api-go" target="_blank" style="color: #667eea; text-decoration: none;">GitHub - Código Fonte</a></p>
                </div>
                <div style="background: #f8f9fa; padding: 15px; border-radius: 8px; border-left: 3px solid #667eea;">
                    <h4 style="color: #667eea; margin-bottom: 8px;">🗄️ Banco de Dados</h4>
                    <p>PostgreSQL no Render</p>
                </div>
                <div style="background: #f8f9fa; padding: 15px; border-radius: 8px; border-left: 3px solid #667eea;">
                    <h4 style="color: #667eea; margin-bottom: 8px;">🐳 Containerização</h4>
                    <p>Docker & Docker Compose</p>
                </div>
                <div style="background: #f8f9fa; padding: 15px; border-radius: 8px; border-left: 3px solid #667eea;">
                    <h4 style="color: #667eea; margin-bottom: 8px;"> Deploy</h4>
                    <p>CI/CD Automático</p>
                </div>
            </div>
        </div>
        
        <div class="footer">
            <p>Desenvolvido por Gabriel - Desafio Técnico Full Stack</p>
            <p>API Version 1.0.0 | Status: Online </p>
        </div>
    </div>
</body>
</html>`

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte(html))
}

// HealthCheckHandler retorna status de saúde da aplicação
func (h *Handler) HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"ok","service":"data-importer-api"}`))
}

// TestHandler retorna informações de teste da API
func (h *Handler) TestHandler(w http.ResponseWriter, r *http.Request) {
	testInfo := map[string]interface{}{
		"status": "ok",
		"message": "API funcionando corretamente",
		"timestamp": time.Now().Format(time.RFC3339),
        "endpoints": map[string]string{
			"health": "/health",
			"login": "/auth/login",
			"customers": "/api/customers",
			"reports": "/api/reports/billing/monthly",
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(testInfo)
}
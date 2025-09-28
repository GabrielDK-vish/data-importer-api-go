package service

import (
	"context"
	"data-importer-api-go/internal/models"
	"data-importer-api-go/internal/repository"
	"fmt"
	
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	repo *repository.Repository
}

func NewService(repo *repository.Repository) *Service {
	return &Service{repo: repo}
}

// GetAllCustomers retorna todos os clientes
func (s *Service) GetAllCustomers(ctx context.Context) ([]models.Customer, error) {
	customers, err := s.repo.GetAllCustomers(ctx)
	if err != nil {
		return nil, fmt.Errorf("erro no service ao buscar clientes: %w", err)
	}
	return customers, nil
}

// GetAllPartners retorna todos os parceiros
func (s *Service) GetAllPartners(ctx context.Context) ([]models.Partner, error) {
	partners, err := s.repo.GetAllPartners(ctx)
	if err != nil {
		return nil, fmt.Errorf("erro no service ao buscar parceiros: %w", err)
	}
	return partners, nil
}

// GetAllProducts retorna todos os produtos
func (s *Service) GetAllProducts(ctx context.Context) ([]models.Product, error) {
	products, err := s.repo.GetAllProducts(ctx)
	if err != nil {
		return nil, fmt.Errorf("erro no service ao buscar produtos: %w", err)
	}
	return products, nil
}

// GetAllUsages retorna todos os usages
func (s *Service) GetAllUsages(ctx context.Context) ([]models.Usage, error) {
	usages, err := s.repo.GetAllUsages(ctx)
	if err != nil {
		return nil, fmt.Errorf("erro no service ao buscar usages: %w", err)
	}
	return usages, nil
}

// GetUsageByCustomer retorna o uso de um cliente específico
func (s *Service) GetUsageByCustomer(ctx context.Context, customerID int) ([]models.Usage, error) {
	// Validar se o customerID é válido
	if customerID <= 0 {
		return nil, fmt.Errorf("ID do cliente inválido")
	}

	usages, err := s.repo.GetUsageByCustomer(ctx, customerID)
	if err != nil {
		return nil, fmt.Errorf("erro no service ao buscar uso do cliente: %w", err)
	}
	return usages, nil
}

// GetBillingMonthly retorna faturamento por mês
func (s *Service) GetBillingMonthly(ctx context.Context) ([]models.BillingReport, error) {
	reports, err := s.repo.GetBillingMonthly(ctx)
	if err != nil {
		return nil, fmt.Errorf("erro no service ao buscar faturamento mensal: %w", err)
	}
	return reports, nil
}

// GetBillingByCategory retorna faturamento por categoria
func (s *Service) GetBillingByCategory(ctx context.Context) ([]models.CategoryBillingReport, error) {
	reports, err := s.repo.GetBillingByCategory(ctx)
	if err != nil {
		return nil, fmt.Errorf("erro no service ao buscar faturamento por categoria: %w", err)
	}
	return reports, nil
}

// GetBillingByResource retorna faturamento por recurso
func (s *Service) GetBillingByResource(ctx context.Context) ([]models.ResourceBillingReport, error) {
	reports, err := s.repo.GetBillingByResource(ctx)
	if err != nil {
		return nil, fmt.Errorf("erro no service ao buscar faturamento por recurso: %w", err)
	}
	return reports, nil
}

// GetBillingByCustomer retorna faturamento por cliente
func (s *Service) GetBillingByCustomer(ctx context.Context) ([]models.CustomerBillingReport, error) {
	reports, err := s.repo.GetBillingByCustomer(ctx)
	if err != nil {
		return nil, fmt.Errorf("erro no service ao buscar faturamento por cliente: %w", err)
	}
	return reports, nil
}

// GetKPIData retorna os dados de KPI do sistema
func (s *Service) GetKPIData(ctx context.Context) (*models.KPIData, error) {
	// Obter dados de KPI do repositório
	kpiData, err := s.repo.GetKPIData(ctx)
	if err != nil {
		return nil, fmt.Errorf("erro no service ao buscar dados de KPI: %w", err)
	}
	
	// Se não houver dados de KPI, criar um objeto vazio
	if kpiData == nil {
		kpiData = &models.KPIData{
			TotalRecords: 0,
			TotalCategories: 0,
			TotalResources: 0,
			TotalCustomers: 0,
			AvgBillingPerMonth: 0,
			ProcessingTimeMs: 0,
			LastUpdated: nil,
		}
	}
	
	return kpiData, nil
}

// GetBillingByProduct retorna faturamento por produto
func (s *Service) GetBillingByProduct(ctx context.Context) ([]models.BillingByProduct, error) {
	reports, err := s.repo.GetBillingByProduct(ctx)
	if err != nil {
		return nil, fmt.Errorf("erro no service ao buscar faturamento por produto: %w", err)
	}
	return reports, nil
}

// GetBillingByPartner retorna faturamento por parceiro
func (s *Service) GetBillingByPartner(ctx context.Context) ([]models.BillingByPartner, error) {
	reports, err := s.repo.GetBillingByPartner(ctx)
	if err != nil {
		return nil, fmt.Errorf("erro no service ao buscar faturamento por parceiro: %w", err)
	}
	return reports, nil
}

// ProcessImportData processa dados de importação com inserção em lote
func (s *Service) ProcessImportData(ctx context.Context, partners []models.Partner, customers []models.Customer, products []models.Product, usages []models.Usage) error {
	// Inserir partners individualmente para obter IDs
	partnerIDMap := make(map[string]int)
	for i := range partners {
		if err := s.repo.InsertPartner(ctx, &partners[i]); err != nil {
			fmt.Printf("Aviso ao inserir parceiro %s: %v\n", partners[i].PartnerID, err)
		} else {
			partnerIDMap[partners[i].PartnerID] = partners[i].ID
		}
	}

	// Inserir customers individualmente para obter IDs
	customerIDMap := make(map[string]int)
	for i := range customers {
		if err := s.repo.InsertCustomer(ctx, &customers[i]); err != nil {
			fmt.Printf("Aviso ao inserir cliente %s: %v\n", customers[i].CustomerID, err)
		} else {
			customerIDMap[customers[i].CustomerID] = customers[i].ID
		}
	}

	// Inserir products individualmente para obter IDs
	productIDMap := make(map[string]int)
	for i := range products {
		if err := s.repo.InsertProduct(ctx, &products[i]); err != nil {
			fmt.Printf("Aviso ao inserir produto %s: %v\n", products[i].ProductID, err)
		} else {
			productIDMap[products[i].ProductID] = products[i].ID
		}
	}

	// Atualizar usages com os IDs corretos
	for i := range usages {
		// Buscar partner_id baseado no partner_id do usage
		if partnerID, exists := partnerIDMap[usages[i].PartnerIDStr]; exists {
			usages[i].PartnerID = partnerID
		}
		
		// Buscar customer_id baseado no customer_id do usage
		if customerID, exists := customerIDMap[usages[i].CustomerIDStr]; exists {
			usages[i].CustomerID = customerID
		}
		
		// Buscar product_id baseado no product_id do usage
		if productID, exists := productIDMap[usages[i].ProductIDStr]; exists {
			usages[i].ProductID = productID
		}
	}

	// Inserir usages em lote
	if len(usages) > 0 {
		if err := s.repo.BulkInsertUsages(ctx, usages); err != nil {
			return fmt.Errorf("erro ao inserir usos em lote: %w", err)
		}
	}

	return nil
}

// ProcessImportDataWithReplace processa dados de importação substituindo dados existentes
func (s *Service) ProcessImportDataWithReplace(ctx context.Context, partners []models.Partner, customers []models.Customer, products []models.Product, usages []models.Usage) error {
	// Limpar dados existentes antes de inserir novos
	if err := s.repo.ClearAllData(ctx); err != nil {
		return fmt.Errorf("erro ao limpar dados existentes: %w", err)
	}

	// Processar dados normalmente
	return s.ProcessImportData(ctx, partners, customers, products, usages)
}

// ValidateUserCredentials valida credenciais de usuário no banco de dados
func (s *Service) ValidateUserCredentials(ctx context.Context, username, password string) (*models.User, error) {
	user, err := s.repo.GetUserByUsername(ctx, username)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar usuário: %w", err)
	}
	
	if user == nil {
		return nil, fmt.Errorf("usuário não encontrado")
	}
	
	// Verificar senha usando bcrypt
	if err := s.comparePassword(password, user.PasswordHash); err != nil {
		return nil, fmt.Errorf("senha inválida")
	}
	
	return user, nil
}

func (s *Service) comparePassword(password, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}

// GetPartnerIDByPartnerID busca o ID numérico de um partner pelo partner_id string
func (s *Service) GetPartnerIDByPartnerID(ctx context.Context, partnerID string) (int, error) {
	return s.repo.GetPartnerIDByPartnerID(ctx, partnerID)
}


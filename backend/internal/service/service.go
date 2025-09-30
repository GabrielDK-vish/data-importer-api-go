package service

import (
	"context"
	"data-importer-api-go/internal/models"
	"data-importer-api-go/internal/repository"
	"fmt"
	"time"
	
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
	// Iniciar contagem de tempo de processamento
	startTime := time.Now()
	
	// Não limpar dados existentes para manter o frontend funcionando
	// Apenas registrar que estamos processando novos dados
	fmt.Printf("Processando importação: %d partners, %d customers, %d products, %d usages\n", 
		len(partners), len(customers), len(products), len(usages))

	// Inserir partners individualmente para obter IDs
	partnerIDMap := make(map[string]int)
	for i := range partners {
		if partners[i].PartnerID == "" {
			fmt.Printf("Ignorando parceiro com ID vazio na posição %d\n", i)
			continue
		}
		if err := s.repo.InsertPartner(ctx, &partners[i]); err != nil {
			fmt.Printf("Aviso ao inserir parceiro %s: %v\n", partners[i].PartnerID, err)
		} else {
			partnerIDMap[partners[i].PartnerID] = partners[i].ID
			fmt.Printf("Partner inserido: %s (ID: %d)\n", partners[i].PartnerID, partners[i].ID)
		}
	}

	// Inserir customers individualmente para obter IDs
	customerIDMap := make(map[string]int)
	for i := range customers {
		if customers[i].CustomerID == "" {
			fmt.Printf("Ignorando cliente com ID vazio na posição %d\n", i)
			continue
		}
		if err := s.repo.InsertCustomer(ctx, &customers[i]); err != nil {
			fmt.Printf("Aviso ao inserir cliente %s: %v\n", customers[i].CustomerID, err)
		} else {
			customerIDMap[customers[i].CustomerID] = customers[i].ID
			fmt.Printf("Customer inserido: %s (ID: %d)\n", customers[i].CustomerID, customers[i].ID)
		}
	}

	// Inserir products individualmente para obter IDs
	productIDMap := make(map[string]int)
	for i := range products {
		if products[i].ProductID == "" {
			fmt.Printf("Ignorando produto com ID vazio na posição %d\n", i)
			continue
		}
		if err := s.repo.InsertProduct(ctx, &products[i]); err != nil {
			fmt.Printf("Aviso ao inserir produto %s: %v\n", products[i].ProductID, err)
		} else {
			productIDMap[products[i].ProductID] = products[i].ID
			fmt.Printf("Product inserido: %s (ID: %d)\n", products[i].ProductID, products[i].ID)
		}
	}

	// Verificar se temos IDs mapeados
	fmt.Printf("Mapeamento de IDs: %d partners, %d customers, %d products\n", 
		len(partnerIDMap), len(customerIDMap), len(productIDMap))

	// Verificar se temos pelo menos alguns IDs mapeados
	if len(partnerIDMap) == 0 || len(customerIDMap) == 0 || len(productIDMap) == 0 {
		fmt.Printf("Aviso: Um ou mais mapas de ID estão vazios. Isso pode causar problemas na importação.\n")
	}

	// Atualizar usages com os IDs corretos
	validUsages := make([]models.Usage, 0, len(usages))
	for i := range usages {
		// Verificar se temos os campos temporários preenchidos
		if usages[i].PartnerIDStr == "" || usages[i].CustomerIDStr == "" || usages[i].ProductIDStr == "" {
			fmt.Printf("Usage ignorado: campos de ID temporários vazios (linha %d)\n", i+1)
			continue
		}

		// Buscar partner_id baseado no partner_id do usage
		partnerID, partnerExists := partnerIDMap[usages[i].PartnerIDStr]
		if !partnerExists {
			fmt.Printf("Partner ID não encontrado para: %s (linha %d)\n", usages[i].PartnerIDStr, i+1)
			continue
		}
		usages[i].PartnerID = partnerID
		
		// Buscar customer_id baseado no customer_id do usage
		customerID, customerExists := customerIDMap[usages[i].CustomerIDStr]
		if !customerExists {
			fmt.Printf("Customer ID não encontrado para: %s (linha %d)\n", usages[i].CustomerIDStr, i+1)
			continue
		}
		usages[i].CustomerID = customerID
		
		// Buscar product_id baseado no product_id do usage
		productID, productExists := productIDMap[usages[i].ProductIDStr]
		if !productExists {
			fmt.Printf("Product ID não encontrado para: %s (linha %d)\n", usages[i].ProductIDStr, i+1)
			continue
		}
		usages[i].ProductID = productID

		// Verificar se a quantidade é válida
		if usages[i].Quantity <= 0 {
			fmt.Printf("Usage ignorado: quantidade inválida %.2f (linha %d)\n", usages[i].Quantity, i+1)
			continue
		}

		// Adicionar à lista de usages válidos
		validUsages = append(validUsages, usages[i])
		fmt.Printf("Usage mapeado: Partner=%d, Customer=%d, Product=%d, Quantidade=%f\n", 
			usages[i].PartnerID, usages[i].CustomerID, usages[i].ProductID, usages[i].Quantity)
	}

	// Inserir usages em lote
	if len(validUsages) > 0 {
		fmt.Printf("Inserindo %d usages válidos em lote\n", len(validUsages))
		if err := s.repo.BulkInsertUsages(ctx, validUsages); err != nil {
			fmt.Printf("Erro ao inserir usos em lote: %v\n", err)
			return fmt.Errorf("erro ao inserir usos em lote: %w", err)
		}
		fmt.Printf("Inserção em lote concluída com sucesso!\n")
	} else {
		fmt.Printf("Nenhum usage válido para inserir\n")
	}

	// Calcular e salvar a métrica de tempo de processamento
	endTime := time.Now()
	processingTime := endTime.Sub(startTime)
	
	metric := &models.ProcessingMetric{
		FileName:     "Reconfile fornecedores.xlsx",
		StartTime:    startTime,
		EndTime:      endTime,
		DurationMs:   processingTime.Milliseconds(),
		RecordsCount: len(usages),
		CreatedAt:    time.Now(),
	}
	
	if err := s.repo.SaveProcessingMetric(ctx, metric); err != nil {
		fmt.Printf("Erro ao salvar métrica de processamento: %v\n", err)
		// Não retornamos erro aqui para não interromper o fluxo principal
	} else {
		fmt.Printf("Métrica de processamento salva: %d ms para processar %d registros\n", 
			metric.DurationMs, metric.RecordsCount)
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

// GetProcessingMetrics retorna todas as métricas de processamento
func (s *Service) GetProcessingMetrics(ctx context.Context) ([]models.ProcessingMetric, error) {
	metrics, err := s.repo.GetAllProcessingMetrics(ctx)
	if err != nil {
		return nil, fmt.Errorf("erro no service ao buscar métricas de processamento: %w", err)
	}
	return metrics, nil
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


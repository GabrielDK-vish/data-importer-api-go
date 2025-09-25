package service

import (
	"context"
	"data-importer-api-go/internal/models"
	"data-importer-api-go/internal/repository"
	"fmt"
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

// ProcessImportData processa dados de importação
func (s *Service) ProcessImportData(ctx context.Context, partners []models.Partner, customers []models.Customer, products []models.Product, usages []models.Usage) error {
	// Inserir partners
	for _, partner := range partners {
		if err := s.repo.InsertPartner(ctx, &partner); err != nil {
			return fmt.Errorf("erro ao inserir parceiro %s: %w", partner.PartnerID, err)
		}
	}

	// Inserir customers
	for _, customer := range customers {
		if err := s.repo.InsertCustomer(ctx, &customer); err != nil {
			return fmt.Errorf("erro ao inserir cliente %s: %w", customer.CustomerID, err)
		}
	}

	// Inserir products
	for _, product := range products {
		if err := s.repo.InsertProduct(ctx, &product); err != nil {
			return fmt.Errorf("erro ao inserir produto %s: %w", product.ProductID, err)
		}
	}

	// Inserir usages em lote para performance
	if len(usages) > 0 {
		if err := s.repo.BulkInsertUsages(ctx, usages); err != nil {
			return fmt.Errorf("erro ao inserir usos em lote: %w", err)
		}
	}

	return nil
}
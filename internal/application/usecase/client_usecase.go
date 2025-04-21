package usecase

import (
	"fmt"
	"time"

	"github.com/lugondev/erp-warehouse-simple/internal/domain/entity"
)

type ClientUseCase interface {
	CreateClient(client *entity.Client) error
	GetClientByID(id uint) (*entity.Client, error)
	GetClientByCode(code string) (*entity.Client, error)
	GetClientByEmail(email string) (*entity.Client, error)
	UpdateClient(client *entity.Client) error
	DeleteClient(id uint) error
	ListClients(filter entity.ClientFilter) ([]entity.Client, error)

	// Address methods
	CreateAddress(address *entity.ClientAddress) error
	UpdateAddress(address *entity.ClientAddress) error
	DeleteAddress(id uint) error
	GetAddressesByClientID(clientID uint) ([]entity.ClientAddress, error)

	// Order history methods
	GetOrderHistory(clientID uint) (*entity.ClientOrderHistory, error)

	// Debt management methods
	GetClientDebt(clientID uint) (*entity.ClientDebt, error)
	UpdateClientDebt(clientID uint, amount float64) error

	// Loyalty methods
	UpdateLoyaltyPoints(clientID uint, points int) error
	UpdateLoyaltyTier(clientID uint, tier entity.ClientLoyaltyTier) error
	CalculateLoyaltyTier(clientID uint) (entity.ClientLoyaltyTier, error)
}

type ClientUseCaseImpl struct {
	clientRepo entity.ClientRepository
}

func NewClientUseCase(clientRepo entity.ClientRepository) ClientUseCase {
	return &ClientUseCaseImpl{
		clientRepo: clientRepo,
	}
}

// CreateClient creates a new client
func (uc *ClientUseCaseImpl) CreateClient(client *entity.Client) error {
	// Generate client code if not provided
	if client.Code == "" {
		client.Code = fmt.Sprintf("CUST-%d", time.Now().Unix())
	}

	// Set default values
	if client.Type == "" {
		client.Type = entity.ClientTypeIndividual
	}
	if client.LoyaltyTier == "" {
		client.LoyaltyTier = entity.ClientLoyaltyTierStandard
	}

	if err := uc.clientRepo.Create(client); err != nil {
		return err
	}

	// Audit logging would be done here in a real implementation
	return nil
}

// GetClientByID gets a client by ID
func (uc *ClientUseCaseImpl) GetClientByID(id uint) (*entity.Client, error) {
	return uc.clientRepo.FindByID(id)
}

// GetClientByCode gets a client by code
func (uc *ClientUseCaseImpl) GetClientByCode(code string) (*entity.Client, error) {
	return uc.clientRepo.FindByCode(code)
}

// GetClientByEmail gets a client by email
func (uc *ClientUseCaseImpl) GetClientByEmail(email string) (*entity.Client, error) {
	return uc.clientRepo.FindByEmail(email)
}

// UpdateClient updates a client
func (uc *ClientUseCaseImpl) UpdateClient(client *entity.Client) error {
	// Check if client exists
	_, err := uc.clientRepo.FindByID(client.ID)
	if err != nil {
		return err
	}

	if err := uc.clientRepo.Update(client); err != nil {
		return err
	}

	// Audit logging would be done here in a real implementation
	return nil
}

// DeleteClient deletes a client
func (uc *ClientUseCaseImpl) DeleteClient(id uint) error {
	// Check if client exists
	_, err := uc.clientRepo.FindByID(id)
	if err != nil {
		return err
	}

	if err := uc.clientRepo.Delete(id); err != nil {
		return err
	}

	// Audit logging would be done here in a real implementation
	return nil
}

// ListClients lists clients with optional filtering
func (uc *ClientUseCaseImpl) ListClients(filter entity.ClientFilter) ([]entity.Client, error) {
	return uc.clientRepo.List(filter)
}

// CreateAddress creates a new client address
func (uc *ClientUseCaseImpl) CreateAddress(address *entity.ClientAddress) error {
	if err := uc.clientRepo.CreateAddress(address); err != nil {
		return err
	}

	// Audit logging would be done here in a real implementation
	return nil
}

// UpdateAddress updates a client address
func (uc *ClientUseCaseImpl) UpdateAddress(address *entity.ClientAddress) error {
	// Get existing address to compare changes
	addresses, err := uc.clientRepo.FindAddressesByClientID(address.ClientID)
	if err != nil {
		return err
	}

	var existingAddress *entity.ClientAddress
	for i := range addresses {
		if addresses[i].ID == address.ID {
			existingAddress = &addresses[i]
			break
		}
	}

	if existingAddress == nil {
		return fmt.Errorf("address not found")
	}

	if err := uc.clientRepo.UpdateAddress(address); err != nil {
		return err
	}

	// Audit logging would be done here in a real implementation
	return nil
}

// DeleteAddress deletes a client address
func (uc *ClientUseCaseImpl) DeleteAddress(id uint) error {
	// We need to find the address first to get its details for the audit log
	// But we don't have a direct method to find an address by ID
	// This is a limitation in the current design

	if err := uc.clientRepo.DeleteAddress(id); err != nil {
		return err
	}

	// Audit logging would be done here in a real implementation
	return nil
}

// GetAddressesByClientID gets all addresses for a client
func (uc *ClientUseCaseImpl) GetAddressesByClientID(clientID uint) ([]entity.ClientAddress, error) {
	return uc.clientRepo.FindAddressesByClientID(clientID)
}

// GetOrderHistory gets a client's order history
func (uc *ClientUseCaseImpl) GetOrderHistory(clientID uint) (*entity.ClientOrderHistory, error) {
	return uc.clientRepo.GetOrderHistory(clientID)
}

// GetClientDebt gets a client's debt information
func (uc *ClientUseCaseImpl) GetClientDebt(clientID uint) (*entity.ClientDebt, error) {
	return uc.clientRepo.GetClientDebt(clientID)
}

// UpdateClientDebt updates a client's debt amount
func (uc *ClientUseCaseImpl) UpdateClientDebt(clientID uint, amount float64) error {
	// Check if client exists
	_, err := uc.clientRepo.FindByID(clientID)
	if err != nil {
		return err
	}

	if err := uc.clientRepo.UpdateClientDebt(clientID, amount); err != nil {
		return err
	}

	// Audit logging would be done here in a real implementation
	return nil
}

// UpdateLoyaltyPoints updates a client's loyalty points
func (uc *ClientUseCaseImpl) UpdateLoyaltyPoints(clientID uint, points int) error {
	// Check if client exists
	_, err := uc.clientRepo.FindByID(clientID)
	if err != nil {
		return err
	}

	if err := uc.clientRepo.UpdateLoyaltyPoints(clientID, points); err != nil {
		return err
	}

	// Audit logging would be done here in a real implementation

	// Recalculate loyalty tier based on new points
	if _, err := uc.CalculateLoyaltyTier(clientID); err != nil {
		return err
	}

	return nil
}

// UpdateLoyaltyTier updates a client's loyalty tier
func (uc *ClientUseCaseImpl) UpdateLoyaltyTier(clientID uint, tier entity.ClientLoyaltyTier) error {
	// Check if client exists
	_, err := uc.clientRepo.FindByID(clientID)
	if err != nil {
		return err
	}

	if err := uc.clientRepo.UpdateLoyaltyTier(clientID, tier); err != nil {
		return err
	}

	// Audit logging would be done here in a real implementation
	return nil
}

// CalculateLoyaltyTier calculates and updates a client's loyalty tier based on points and purchase history
func (uc *ClientUseCaseImpl) CalculateLoyaltyTier(clientID uint) (entity.ClientLoyaltyTier, error) {
	// Get client
	client, err := uc.clientRepo.FindByID(clientID)
	if err != nil {
		return "", err
	}

	// Get order history
	history, err := uc.clientRepo.GetOrderHistory(clientID)
	if err != nil {
		return "", err
	}

	// Calculate tier based on points and purchase history
	var newTier entity.ClientLoyaltyTier

	// Simple tier calculation logic
	// This could be made more sophisticated based on business rules
	if client.LoyaltyPoints >= 10000 || history.TotalSpent >= 100000 {
		newTier = entity.ClientLoyaltyTierPlatinum
	} else if client.LoyaltyPoints >= 5000 || history.TotalSpent >= 50000 {
		newTier = entity.ClientLoyaltyTierGold
	} else if client.LoyaltyPoints >= 1000 || history.TotalSpent >= 10000 {
		newTier = entity.ClientLoyaltyTierSilver
	} else {
		newTier = entity.ClientLoyaltyTierStandard
	}

	// Only update if tier has changed
	if newTier != client.LoyaltyTier {
		if err := uc.UpdateLoyaltyTier(clientID, newTier); err != nil {
			return "", err
		}
	}

	return newTier, nil
}

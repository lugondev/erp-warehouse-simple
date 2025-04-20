package usecase

import (
	"fmt"
	"time"

	"github.com/lugondev/erp-warehouse-simple/internal/domain/entity"
)

type CustomerUseCase interface {
	CreateCustomer(customer *entity.Customer) error
	GetCustomerByID(id uint) (*entity.Customer, error)
	GetCustomerByCode(code string) (*entity.Customer, error)
	GetCustomerByEmail(email string) (*entity.Customer, error)
	UpdateCustomer(customer *entity.Customer) error
	DeleteCustomer(id uint) error
	ListCustomers(filter entity.CustomerFilter) ([]entity.Customer, error)

	// Address methods
	CreateAddress(address *entity.CustomerAddress) error
	UpdateAddress(address *entity.CustomerAddress) error
	DeleteAddress(id uint) error
	GetAddressesByCustomerID(customerID uint) ([]entity.CustomerAddress, error)

	// Order history methods
	GetOrderHistory(customerID uint) (*entity.CustomerOrderHistory, error)

	// Debt management methods
	GetCustomerDebt(customerID uint) (*entity.CustomerDebt, error)
	UpdateCustomerDebt(customerID uint, amount float64) error

	// Loyalty methods
	UpdateLoyaltyPoints(customerID uint, points int) error
	UpdateLoyaltyTier(customerID uint, tier entity.CustomerLoyaltyTier) error
	CalculateLoyaltyTier(customerID uint) (entity.CustomerLoyaltyTier, error)
}

type CustomerUseCaseImpl struct {
	customerRepo entity.CustomerRepository
}

func NewCustomerUseCase(customerRepo entity.CustomerRepository) CustomerUseCase {
	return &CustomerUseCaseImpl{
		customerRepo: customerRepo,
	}
}

// CreateCustomer creates a new customer
func (uc *CustomerUseCaseImpl) CreateCustomer(customer *entity.Customer) error {
	// Generate customer code if not provided
	if customer.Code == "" {
		customer.Code = fmt.Sprintf("CUST-%d", time.Now().Unix())
	}

	// Set default values
	if customer.Type == "" {
		customer.Type = entity.CustomerTypeIndividual
	}
	if customer.LoyaltyTier == "" {
		customer.LoyaltyTier = entity.CustomerLoyaltyTierStandard
	}

	if err := uc.customerRepo.Create(customer); err != nil {
		return err
	}

	// Audit logging would be done here in a real implementation
	return nil
}

// GetCustomerByID gets a customer by ID
func (uc *CustomerUseCaseImpl) GetCustomerByID(id uint) (*entity.Customer, error) {
	return uc.customerRepo.FindByID(id)
}

// GetCustomerByCode gets a customer by code
func (uc *CustomerUseCaseImpl) GetCustomerByCode(code string) (*entity.Customer, error) {
	return uc.customerRepo.FindByCode(code)
}

// GetCustomerByEmail gets a customer by email
func (uc *CustomerUseCaseImpl) GetCustomerByEmail(email string) (*entity.Customer, error) {
	return uc.customerRepo.FindByEmail(email)
}

// UpdateCustomer updates a customer
func (uc *CustomerUseCaseImpl) UpdateCustomer(customer *entity.Customer) error {
	// Check if customer exists
	_, err := uc.customerRepo.FindByID(customer.ID)
	if err != nil {
		return err
	}

	if err := uc.customerRepo.Update(customer); err != nil {
		return err
	}

	// Audit logging would be done here in a real implementation
	return nil
}

// DeleteCustomer deletes a customer
func (uc *CustomerUseCaseImpl) DeleteCustomer(id uint) error {
	// Check if customer exists
	_, err := uc.customerRepo.FindByID(id)
	if err != nil {
		return err
	}

	if err := uc.customerRepo.Delete(id); err != nil {
		return err
	}

	// Audit logging would be done here in a real implementation
	return nil
}

// ListCustomers lists customers with optional filtering
func (uc *CustomerUseCaseImpl) ListCustomers(filter entity.CustomerFilter) ([]entity.Customer, error) {
	return uc.customerRepo.List(filter)
}

// CreateAddress creates a new customer address
func (uc *CustomerUseCaseImpl) CreateAddress(address *entity.CustomerAddress) error {
	if err := uc.customerRepo.CreateAddress(address); err != nil {
		return err
	}

	// Audit logging would be done here in a real implementation
	return nil
}

// UpdateAddress updates a customer address
func (uc *CustomerUseCaseImpl) UpdateAddress(address *entity.CustomerAddress) error {
	// Get existing address to compare changes
	addresses, err := uc.customerRepo.FindAddressesByCustomerID(address.CustomerID)
	if err != nil {
		return err
	}

	var existingAddress *entity.CustomerAddress
	for i := range addresses {
		if addresses[i].ID == address.ID {
			existingAddress = &addresses[i]
			break
		}
	}

	if existingAddress == nil {
		return fmt.Errorf("address not found")
	}

	if err := uc.customerRepo.UpdateAddress(address); err != nil {
		return err
	}

	// Audit logging would be done here in a real implementation
	return nil
}

// DeleteAddress deletes a customer address
func (uc *CustomerUseCaseImpl) DeleteAddress(id uint) error {
	// We need to find the address first to get its details for the audit log
	// But we don't have a direct method to find an address by ID
	// This is a limitation in the current design

	if err := uc.customerRepo.DeleteAddress(id); err != nil {
		return err
	}

	// Audit logging would be done here in a real implementation
	return nil
}

// GetAddressesByCustomerID gets all addresses for a customer
func (uc *CustomerUseCaseImpl) GetAddressesByCustomerID(customerID uint) ([]entity.CustomerAddress, error) {
	return uc.customerRepo.FindAddressesByCustomerID(customerID)
}

// GetOrderHistory gets a customer's order history
func (uc *CustomerUseCaseImpl) GetOrderHistory(customerID uint) (*entity.CustomerOrderHistory, error) {
	return uc.customerRepo.GetOrderHistory(customerID)
}

// GetCustomerDebt gets a customer's debt information
func (uc *CustomerUseCaseImpl) GetCustomerDebt(customerID uint) (*entity.CustomerDebt, error) {
	return uc.customerRepo.GetCustomerDebt(customerID)
}

// UpdateCustomerDebt updates a customer's debt amount
func (uc *CustomerUseCaseImpl) UpdateCustomerDebt(customerID uint, amount float64) error {
	// Check if customer exists
	_, err := uc.customerRepo.FindByID(customerID)
	if err != nil {
		return err
	}

	if err := uc.customerRepo.UpdateCustomerDebt(customerID, amount); err != nil {
		return err
	}

	// Audit logging would be done here in a real implementation
	return nil
}

// UpdateLoyaltyPoints updates a customer's loyalty points
func (uc *CustomerUseCaseImpl) UpdateLoyaltyPoints(customerID uint, points int) error {
	// Check if customer exists
	_, err := uc.customerRepo.FindByID(customerID)
	if err != nil {
		return err
	}

	if err := uc.customerRepo.UpdateLoyaltyPoints(customerID, points); err != nil {
		return err
	}

	// Audit logging would be done here in a real implementation

	// Recalculate loyalty tier based on new points
	if _, err := uc.CalculateLoyaltyTier(customerID); err != nil {
		return err
	}

	return nil
}

// UpdateLoyaltyTier updates a customer's loyalty tier
func (uc *CustomerUseCaseImpl) UpdateLoyaltyTier(customerID uint, tier entity.CustomerLoyaltyTier) error {
	// Check if customer exists
	_, err := uc.customerRepo.FindByID(customerID)
	if err != nil {
		return err
	}

	if err := uc.customerRepo.UpdateLoyaltyTier(customerID, tier); err != nil {
		return err
	}

	// Audit logging would be done here in a real implementation
	return nil
}

// CalculateLoyaltyTier calculates and updates a customer's loyalty tier based on points and purchase history
func (uc *CustomerUseCaseImpl) CalculateLoyaltyTier(customerID uint) (entity.CustomerLoyaltyTier, error) {
	// Get customer
	customer, err := uc.customerRepo.FindByID(customerID)
	if err != nil {
		return "", err
	}

	// Get order history
	history, err := uc.customerRepo.GetOrderHistory(customerID)
	if err != nil {
		return "", err
	}

	// Calculate tier based on points and purchase history
	var newTier entity.CustomerLoyaltyTier

	// Simple tier calculation logic
	// This could be made more sophisticated based on business rules
	if customer.LoyaltyPoints >= 10000 || history.TotalSpent >= 100000 {
		newTier = entity.CustomerLoyaltyTierPlatinum
	} else if customer.LoyaltyPoints >= 5000 || history.TotalSpent >= 50000 {
		newTier = entity.CustomerLoyaltyTierGold
	} else if customer.LoyaltyPoints >= 1000 || history.TotalSpent >= 10000 {
		newTier = entity.CustomerLoyaltyTierSilver
	} else {
		newTier = entity.CustomerLoyaltyTierStandard
	}

	// Only update if tier has changed
	if newTier != customer.LoyaltyTier {
		if err := uc.UpdateLoyaltyTier(customerID, newTier); err != nil {
			return "", err
		}
	}

	return newTier, nil
}

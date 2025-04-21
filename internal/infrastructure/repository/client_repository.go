package repository

import (
	"fmt"
	"time"

	"github.com/lugondev/erp-warehouse-simple/internal/domain/entity"
	"gorm.io/gorm"
)

type ClientRepositoryImpl struct {
	db *gorm.DB
}

func NewClientRepository(db *gorm.DB) entity.ClientRepository {
	return &ClientRepositoryImpl{db: db}
}

// Create creates a new client
func (r *ClientRepositoryImpl) Create(client *entity.Client) error {
	if err := r.db.Create(client).Error; err != nil {
		return fmt.Errorf("failed to create client: %w", err)
	}
	return nil
}

// FindByID finds a client by ID
func (r *ClientRepositoryImpl) FindByID(id uint) (*entity.Client, error) {
	var client entity.Client
	if err := r.db.Preload("Addresses").First(&client, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrRecordNotFound
		}
		return nil, fmt.Errorf("failed to find client by ID: %w", err)
	}
	return &client, nil
}

// FindByCode finds a client by code
func (r *ClientRepositoryImpl) FindByCode(code string) (*entity.Client, error) {
	var client entity.Client
	if err := r.db.Preload("Addresses").Where("code = ?", code).First(&client).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrRecordNotFound
		}
		return nil, fmt.Errorf("failed to find client by code: %w", err)
	}
	return &client, nil
}

// FindByEmail finds a client by email
func (r *ClientRepositoryImpl) FindByEmail(email string) (*entity.Client, error) {
	var client entity.Client
	if err := r.db.Preload("Addresses").Where("email = ?", email).First(&client).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrRecordNotFound
		}
		return nil, fmt.Errorf("failed to find client by email: %w", err)
	}
	return &client, nil
}

// Update updates a client
func (r *ClientRepositoryImpl) Update(client *entity.Client) error {
	if err := r.db.Save(client).Error; err != nil {
		return fmt.Errorf("failed to update client: %w", err)
	}
	return nil
}

// Delete deletes a client
func (r *ClientRepositoryImpl) Delete(id uint) error {
	if err := r.db.Delete(&entity.Client{}, id).Error; err != nil {
		return fmt.Errorf("failed to delete client: %w", err)
	}
	return nil
}

// List lists clients with optional filtering
func (r *ClientRepositoryImpl) List(filter entity.ClientFilter) ([]entity.Client, error) {
	var clients []entity.Client
	query := r.db.Model(&entity.Client{})

	// Apply filters
	if filter.Code != "" {
		query = query.Where("code LIKE ?", "%"+filter.Code+"%")
	}
	if filter.Name != "" {
		query = query.Where("name LIKE ?", "%"+filter.Name+"%")
	}
	if filter.Type != nil {
		query = query.Where("type = ?", *filter.Type)
	}
	if filter.Email != "" {
		query = query.Where("email LIKE ?", "%"+filter.Email+"%")
	}
	if filter.PhoneNumber != "" {
		query = query.Where("phone_number LIKE ?", "%"+filter.PhoneNumber+"%")
	}
	if filter.LoyaltyTier != nil {
		query = query.Where("loyalty_tier = ?", *filter.LoyaltyTier)
	}

	// Address-related filters
	if filter.City != "" || filter.Country != "" {
		query = query.Joins("JOIN client_addresses ON clients.id = client_addresses.client_id")
		if filter.City != "" {
			query = query.Where("client_addresses.city LIKE ?", "%"+filter.City+"%")
		}
		if filter.Country != "" {
			query = query.Where("client_addresses.country LIKE ?", "%"+filter.Country+"%")
		}
		query = query.Group("clients.id") // Avoid duplicates
	}

	if err := query.Preload("Addresses").Find(&clients).Error; err != nil {
		return nil, fmt.Errorf("failed to list clients: %w", err)
	}
	return clients, nil
}

// CreateAddress creates a new client address
func (r *ClientRepositoryImpl) CreateAddress(address *entity.ClientAddress) error {
	// If this is set as default, unset any existing default addresses for this client
	if address.IsDefault {
		if err := r.db.Model(&entity.ClientAddress{}).
			Where("client_id = ? AND is_default = ?", address.ClientID, true).
			Update("is_default", false).Error; err != nil {
			return fmt.Errorf("failed to update existing default addresses: %w", err)
		}
	}

	if err := r.db.Create(address).Error; err != nil {
		return fmt.Errorf("failed to create client address: %w", err)
	}
	return nil
}

// UpdateAddress updates a client address
func (r *ClientRepositoryImpl) UpdateAddress(address *entity.ClientAddress) error {
	// If this is set as default, unset any existing default addresses for this client
	if address.IsDefault {
		if err := r.db.Model(&entity.ClientAddress{}).
			Where("client_id = ? AND id != ? AND is_default = ?", address.ClientID, address.ID, true).
			Update("is_default", false).Error; err != nil {
			return fmt.Errorf("failed to update existing default addresses: %w", err)
		}
	}

	if err := r.db.Save(address).Error; err != nil {
		return fmt.Errorf("failed to update client address: %w", err)
	}
	return nil
}

// DeleteAddress deletes a client address
func (r *ClientRepositoryImpl) DeleteAddress(id uint) error {
	if err := r.db.Delete(&entity.ClientAddress{}, id).Error; err != nil {
		return fmt.Errorf("failed to delete client address: %w", err)
	}
	return nil
}

// FindAddressesByClientID finds all addresses for a client
func (r *ClientRepositoryImpl) FindAddressesByClientID(clientID uint) ([]entity.ClientAddress, error) {
	var addresses []entity.ClientAddress
	if err := r.db.Where("client_id = ?", clientID).Find(&addresses).Error; err != nil {
		return nil, fmt.Errorf("failed to find client addresses: %w", err)
	}
	return addresses, nil
}

// GetOrderHistory gets a client's order history
func (r *ClientRepositoryImpl) GetOrderHistory(clientID uint) (*entity.ClientOrderHistory, error) {
	var history entity.ClientOrderHistory

	// Check if sales_orders table exists
	var tableExists bool
	if err := r.db.Raw("SELECT EXISTS (SELECT FROM information_schema.tables WHERE table_name = 'sales_orders')").Scan(&tableExists).Error; err != nil {
		return nil, fmt.Errorf("failed to check if sales_orders table exists: %w", err)
	}

	if !tableExists {
		// Return empty history if table doesn't exist
		return &history, nil
	}

	// Get total orders count
	var count int64
	if err := r.db.Model(&entity.SalesOrder{}).Where("client_id = ?", clientID).Count(&count).Error; err != nil {
		return nil, fmt.Errorf("failed to count client orders: %w", err)
	}

	// If no orders, return empty history
	history.TotalOrders = int(count)
	if count == 0 {
		return &history, nil
	}

	// Get total spent
	if err := r.db.Model(&entity.SalesOrder{}).Where("client_id = ?", clientID).
		Select("COALESCE(SUM(grand_total), 0)").Scan(&history.TotalSpent).Error; err != nil {
		return nil, fmt.Errorf("failed to calculate total spent: %w", err)
	}

	// Calculate average order value
	history.AverageOrderValue = 0
	if history.TotalOrders > 0 {
		history.AverageOrderValue = history.TotalSpent / float64(history.TotalOrders)
	}

	// Get first order date
	var firstOrder entity.SalesOrder
	if err := r.db.Where("client_id = ?", clientID).Order("order_date ASC").First(&firstOrder).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("failed to get first order date: %w", err)
		}
	} else {
		history.FirstOrderDate = firstOrder.OrderDate
	}

	// Get last order date
	var lastOrder entity.SalesOrder
	if err := r.db.Where("client_id = ?", clientID).Order("order_date DESC").First(&lastOrder).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("failed to get last order date: %w", err)
		}
	} else {
		history.LastOrderDate = lastOrder.OrderDate
	}

	// Get frequent items (top 5)
	// This is a complex query that would need to be adjusted based on how items are stored in orders

	// This is a complex query that would need to be adjusted based on how items are stored in orders
	// For now, we'll return an empty list
	history.FrequentItems = []string{}

	return &history, nil
}

// GetClientDebt gets a client's debt information
func (r *ClientRepositoryImpl) GetClientDebt(clientID uint) (*entity.ClientDebt, error) {
	var debt entity.ClientDebt
	var client entity.Client

	// Get client's current debt
	if err := r.db.Select("current_debt").First(&client, clientID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrRecordNotFound
		}
		return nil, fmt.Errorf("failed to get client debt: %w", err)
	}

	debt.TotalDebt = client.CurrentDebt

	// Check if invoices table exists
	var tableExists bool
	if err := r.db.Raw("SELECT EXISTS (SELECT FROM information_schema.tables WHERE table_name = 'invoices')").Scan(&tableExists).Error; err != nil {
		return nil, fmt.Errorf("failed to check if invoices table exists: %w", err)
	}

	if !tableExists {
		// Return basic debt info if table doesn't exist
		return &debt, nil
	}

	// Get overdue debt (invoices past due date with status not PAID)
	if err := r.db.Model(&entity.Invoice{}).
		Joins("JOIN sales_orders ON invoices.sales_order_id = sales_orders.id").
		Where("sales_orders.client_id = ? AND invoices.status != 'PAID' AND invoices.due_date < ?", clientID, time.Now()).
		Select("COALESCE(SUM(invoices.total_amount), 0)").
		Scan(&debt.OverdueDebt).Error; err != nil {
		return nil, fmt.Errorf("failed to calculate overdue debt: %w", err)
	}

	// Get upcoming payments (invoices not yet due with status not PAID)
	if err := r.db.Model(&entity.Invoice{}).
		Joins("JOIN sales_orders ON invoices.sales_order_id = sales_orders.id").
		Where("sales_orders.client_id = ? AND invoices.status != 'PAID' AND invoices.due_date >= ?", clientID, time.Now()).
		Select("COALESCE(SUM(invoices.total_amount), 0)").
		Scan(&debt.UpcomingPayments).Error; err != nil {
		return nil, fmt.Errorf("failed to calculate upcoming payments: %w", err)
	}

	// Get last payment info
	// This would require a payments table which isn't in the current schema
	// For now, we'll leave these fields with zero values

	return &debt, nil
}

// UpdateClientDebt updates a client's debt amount
func (r *ClientRepositoryImpl) UpdateClientDebt(clientID uint, amount float64) error {
	if err := r.db.Model(&entity.Client{}).Where("id = ?", clientID).
		Update("current_debt", amount).Error; err != nil {
		return fmt.Errorf("failed to update client debt: %w", err)
	}
	return nil
}

// UpdateLoyaltyPoints updates a client's loyalty points
func (r *ClientRepositoryImpl) UpdateLoyaltyPoints(clientID uint, points int) error {
	if err := r.db.Model(&entity.Client{}).Where("id = ?", clientID).
		Update("loyalty_points", points).Error; err != nil {
		return fmt.Errorf("failed to update loyalty points: %w", err)
	}
	return nil
}

// UpdateLoyaltyTier updates a client's loyalty tier
func (r *ClientRepositoryImpl) UpdateLoyaltyTier(clientID uint, tier entity.ClientLoyaltyTier) error {
	if err := r.db.Model(&entity.Client{}).Where("id = ?", clientID).
		Update("loyalty_tier", tier).Error; err != nil {
		return fmt.Errorf("failed to update loyalty tier: %w", err)
	}
	return nil
}

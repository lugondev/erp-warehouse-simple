package repository

import (
	"fmt"
	"time"

	"github.com/lugondev/erp-warehouse-simple/internal/domain/entity"
	"gorm.io/gorm"
)

type CustomerRepositoryImpl struct {
	db *gorm.DB
}

func NewCustomerRepository(db *gorm.DB) entity.CustomerRepository {
	return &CustomerRepositoryImpl{db: db}
}

// Create creates a new customer
func (r *CustomerRepositoryImpl) Create(customer *entity.Customer) error {
	if err := r.db.Create(customer).Error; err != nil {
		return fmt.Errorf("failed to create customer: %w", err)
	}
	return nil
}

// FindByID finds a customer by ID
func (r *CustomerRepositoryImpl) FindByID(id uint) (*entity.Customer, error) {
	var customer entity.Customer
	if err := r.db.Preload("Addresses").First(&customer, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrRecordNotFound
		}
		return nil, fmt.Errorf("failed to find customer by ID: %w", err)
	}
	return &customer, nil
}

// FindByCode finds a customer by code
func (r *CustomerRepositoryImpl) FindByCode(code string) (*entity.Customer, error) {
	var customer entity.Customer
	if err := r.db.Preload("Addresses").Where("code = ?", code).First(&customer).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrRecordNotFound
		}
		return nil, fmt.Errorf("failed to find customer by code: %w", err)
	}
	return &customer, nil
}

// FindByEmail finds a customer by email
func (r *CustomerRepositoryImpl) FindByEmail(email string) (*entity.Customer, error) {
	var customer entity.Customer
	if err := r.db.Preload("Addresses").Where("email = ?", email).First(&customer).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrRecordNotFound
		}
		return nil, fmt.Errorf("failed to find customer by email: %w", err)
	}
	return &customer, nil
}

// Update updates a customer
func (r *CustomerRepositoryImpl) Update(customer *entity.Customer) error {
	if err := r.db.Save(customer).Error; err != nil {
		return fmt.Errorf("failed to update customer: %w", err)
	}
	return nil
}

// Delete deletes a customer
func (r *CustomerRepositoryImpl) Delete(id uint) error {
	if err := r.db.Delete(&entity.Customer{}, id).Error; err != nil {
		return fmt.Errorf("failed to delete customer: %w", err)
	}
	return nil
}

// List lists customers with optional filtering
func (r *CustomerRepositoryImpl) List(filter entity.CustomerFilter) ([]entity.Customer, error) {
	var customers []entity.Customer
	query := r.db.Model(&entity.Customer{})

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
		query = query.Joins("JOIN customer_addresses ON customers.id = customer_addresses.customer_id")
		if filter.City != "" {
			query = query.Where("customer_addresses.city LIKE ?", "%"+filter.City+"%")
		}
		if filter.Country != "" {
			query = query.Where("customer_addresses.country LIKE ?", "%"+filter.Country+"%")
		}
		query = query.Group("customers.id") // Avoid duplicates
	}

	if err := query.Preload("Addresses").Find(&customers).Error; err != nil {
		return nil, fmt.Errorf("failed to list customers: %w", err)
	}
	return customers, nil
}

// CreateAddress creates a new customer address
func (r *CustomerRepositoryImpl) CreateAddress(address *entity.CustomerAddress) error {
	// If this is set as default, unset any existing default addresses for this customer
	if address.IsDefault {
		if err := r.db.Model(&entity.CustomerAddress{}).
			Where("customer_id = ? AND is_default = ?", address.CustomerID, true).
			Update("is_default", false).Error; err != nil {
			return fmt.Errorf("failed to update existing default addresses: %w", err)
		}
	}

	if err := r.db.Create(address).Error; err != nil {
		return fmt.Errorf("failed to create customer address: %w", err)
	}
	return nil
}

// UpdateAddress updates a customer address
func (r *CustomerRepositoryImpl) UpdateAddress(address *entity.CustomerAddress) error {
	// If this is set as default, unset any existing default addresses for this customer
	if address.IsDefault {
		if err := r.db.Model(&entity.CustomerAddress{}).
			Where("customer_id = ? AND id != ? AND is_default = ?", address.CustomerID, address.ID, true).
			Update("is_default", false).Error; err != nil {
			return fmt.Errorf("failed to update existing default addresses: %w", err)
		}
	}

	if err := r.db.Save(address).Error; err != nil {
		return fmt.Errorf("failed to update customer address: %w", err)
	}
	return nil
}

// DeleteAddress deletes a customer address
func (r *CustomerRepositoryImpl) DeleteAddress(id uint) error {
	if err := r.db.Delete(&entity.CustomerAddress{}, id).Error; err != nil {
		return fmt.Errorf("failed to delete customer address: %w", err)
	}
	return nil
}

// FindAddressesByCustomerID finds all addresses for a customer
func (r *CustomerRepositoryImpl) FindAddressesByCustomerID(customerID uint) ([]entity.CustomerAddress, error) {
	var addresses []entity.CustomerAddress
	if err := r.db.Where("customer_id = ?", customerID).Find(&addresses).Error; err != nil {
		return nil, fmt.Errorf("failed to find customer addresses: %w", err)
	}
	return addresses, nil
}

// GetOrderHistory gets a customer's order history
func (r *CustomerRepositoryImpl) GetOrderHistory(customerID uint) (*entity.CustomerOrderHistory, error) {
	var history entity.CustomerOrderHistory

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
	if err := r.db.Model(&entity.SalesOrder{}).Where("customer_id = ?", customerID).Count(&count).Error; err != nil {
		return nil, fmt.Errorf("failed to count customer orders: %w", err)
	}

	// If no orders, return empty history
	history.TotalOrders = int(count)
	if count == 0 {
		return &history, nil
	}

	// Get total spent
	if err := r.db.Model(&entity.SalesOrder{}).Where("customer_id = ?", customerID).
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
	if err := r.db.Where("customer_id = ?", customerID).Order("order_date ASC").First(&firstOrder).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("failed to get first order date: %w", err)
		}
	} else {
		history.FirstOrderDate = firstOrder.OrderDate
	}

	// Get last order date
	var lastOrder entity.SalesOrder
	if err := r.db.Where("customer_id = ?", customerID).Order("order_date DESC").First(&lastOrder).Error; err != nil {
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

// GetCustomerDebt gets a customer's debt information
func (r *CustomerRepositoryImpl) GetCustomerDebt(customerID uint) (*entity.CustomerDebt, error) {
	var debt entity.CustomerDebt
	var customer entity.Customer

	// Get customer's current debt
	if err := r.db.Select("current_debt").First(&customer, customerID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrRecordNotFound
		}
		return nil, fmt.Errorf("failed to get customer debt: %w", err)
	}

	debt.TotalDebt = customer.CurrentDebt

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
		Where("sales_orders.customer_id = ? AND invoices.status != 'PAID' AND invoices.due_date < ?", customerID, time.Now()).
		Select("COALESCE(SUM(invoices.total_amount), 0)").
		Scan(&debt.OverdueDebt).Error; err != nil {
		return nil, fmt.Errorf("failed to calculate overdue debt: %w", err)
	}

	// Get upcoming payments (invoices not yet due with status not PAID)
	if err := r.db.Model(&entity.Invoice{}).
		Joins("JOIN sales_orders ON invoices.sales_order_id = sales_orders.id").
		Where("sales_orders.customer_id = ? AND invoices.status != 'PAID' AND invoices.due_date >= ?", customerID, time.Now()).
		Select("COALESCE(SUM(invoices.total_amount), 0)").
		Scan(&debt.UpcomingPayments).Error; err != nil {
		return nil, fmt.Errorf("failed to calculate upcoming payments: %w", err)
	}

	// Get last payment info
	// This would require a payments table which isn't in the current schema
	// For now, we'll leave these fields with zero values

	return &debt, nil
}

// UpdateCustomerDebt updates a customer's debt amount
func (r *CustomerRepositoryImpl) UpdateCustomerDebt(customerID uint, amount float64) error {
	if err := r.db.Model(&entity.Customer{}).Where("id = ?", customerID).
		Update("current_debt", amount).Error; err != nil {
		return fmt.Errorf("failed to update customer debt: %w", err)
	}
	return nil
}

// UpdateLoyaltyPoints updates a customer's loyalty points
func (r *CustomerRepositoryImpl) UpdateLoyaltyPoints(customerID uint, points int) error {
	if err := r.db.Model(&entity.Customer{}).Where("id = ?", customerID).
		Update("loyalty_points", points).Error; err != nil {
		return fmt.Errorf("failed to update loyalty points: %w", err)
	}
	return nil
}

// UpdateLoyaltyTier updates a customer's loyalty tier
func (r *CustomerRepositoryImpl) UpdateLoyaltyTier(customerID uint, tier entity.CustomerLoyaltyTier) error {
	if err := r.db.Model(&entity.Customer{}).Where("id = ?", customerID).
		Update("loyalty_tier", tier).Error; err != nil {
		return fmt.Errorf("failed to update loyalty tier: %w", err)
	}
	return nil
}

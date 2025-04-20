package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/lugondev/erp-warehouse-simple/internal/domain/entity"
	"gorm.io/gorm"
)

var (
	// ErrInvalidOrderStatus is returned when an operation is not allowed for the current order status
	ErrInvalidOrderStatus = errors.New("invalid order status for this operation")
)

// OrderRepository handles database operations for sales orders and delivery orders
type OrderRepository struct {
	db                *gorm.DB
	inventoryRepo     *InventoryRepository
	sequenceGenerator *SequenceGenerator
}

// NewOrderRepository creates a new OrderRepository
func NewOrderRepository(db *gorm.DB, inventoryRepo *InventoryRepository) *OrderRepository {
	return &OrderRepository{
		db:                db,
		inventoryRepo:     inventoryRepo,
		sequenceGenerator: NewSequenceGenerator(db),
	}
}

// CreateSalesOrder creates a new sales order
func (r *OrderRepository) CreateSalesOrder(ctx context.Context, order *entity.SalesOrder) error {
	if order.ID == "" {
		order.ID = uuid.New().String()
	}

	// Generate order number if not provided
	if order.OrderNumber == "" {
		seq, err := r.sequenceGenerator.NextSequence(ctx, "sales_order")
		if err != nil {
			return err
		}
		order.OrderNumber = fmt.Sprintf("SO-%s-%06d", time.Now().Format("20060102"), seq)
	}

	return r.db.WithContext(ctx).Create(order).Error
}

// GetSalesOrderByID retrieves a sales order by ID
func (r *OrderRepository) GetSalesOrderByID(ctx context.Context, id string) (*entity.SalesOrder, error) {
	var order entity.SalesOrder
	if err := r.db.WithContext(ctx).
		Preload("DeliveryOrders").
		Preload("Invoices").
		First(&order, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrRecordNotFound
		}
		return nil, err
	}
	return &order, nil
}

// GetSalesOrderByNumber retrieves a sales order by order number
func (r *OrderRepository) GetSalesOrderByNumber(ctx context.Context, orderNumber string) (*entity.SalesOrder, error) {
	var order entity.SalesOrder
	if err := r.db.WithContext(ctx).
		Preload("DeliveryOrders").
		Preload("Invoices").
		First(&order, "order_number = ?", orderNumber).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrRecordNotFound
		}
		return nil, err
	}
	return &order, nil
}

// ListSalesOrders retrieves a list of sales orders based on filter
func (r *OrderRepository) ListSalesOrders(ctx context.Context, filter *entity.SalesOrderFilter) ([]entity.SalesOrder, error) {
	var orders []entity.SalesOrder
	query := r.db.WithContext(ctx)

	if filter != nil {
		if filter.OrderNumber != "" {
			query = query.Where("order_number LIKE ?", "%"+filter.OrderNumber+"%")
		}
		if filter.CustomerID != nil {
			query = query.Where("customer_id = ?", *filter.CustomerID)
		}
		if filter.Status != nil {
			query = query.Where("status = ?", *filter.Status)
		}
		if filter.PaymentStatus != nil {
			query = query.Where("payment_status = ?", *filter.PaymentStatus)
		}
		if filter.StartDate != nil {
			query = query.Where("order_date >= ?", *filter.StartDate)
		}
		if filter.EndDate != nil {
			query = query.Where("order_date <= ?", *filter.EndDate)
		}
		if filter.ItemID != "" {
			// This requires a more complex query to search in the JSONB items array
			query = query.Where("items @> ?", fmt.Sprintf(`[{"item_id": "%s"}]`, filter.ItemID))
		}
	}

	if err := query.Order("created_at DESC").Find(&orders).Error; err != nil {
		return nil, err
	}
	return orders, nil
}

// UpdateSalesOrder updates an existing sales order
func (r *OrderRepository) UpdateSalesOrder(ctx context.Context, order *entity.SalesOrder) error {
	return r.db.WithContext(ctx).Save(order).Error
}

// UpdateSalesOrderStatus updates the status of a sales order
func (r *OrderRepository) UpdateSalesOrderStatus(ctx context.Context, id string, status entity.SalesOrderStatus) error {
	return r.db.WithContext(ctx).
		Model(&entity.SalesOrder{}).
		Where("id = ?", id).
		Update("status", status).
		Error
}

// CreateDeliveryOrder creates a new delivery order
func (r *OrderRepository) CreateDeliveryOrder(ctx context.Context, delivery *entity.DeliveryOrder) error {
	if delivery.ID == "" {
		delivery.ID = uuid.New().String()
	}

	// Generate delivery number if not provided
	if delivery.DeliveryNumber == "" {
		seq, err := r.sequenceGenerator.NextSequence(ctx, "delivery_order")
		if err != nil {
			return err
		}
		delivery.DeliveryNumber = fmt.Sprintf("DO-%s-%06d", time.Now().Format("20060102"), seq)
	}

	return r.db.WithContext(ctx).Create(delivery).Error
}

// GetDeliveryOrderByID retrieves a delivery order by ID
func (r *OrderRepository) GetDeliveryOrderByID(ctx context.Context, id string) (*entity.DeliveryOrder, error) {
	var delivery entity.DeliveryOrder
	if err := r.db.WithContext(ctx).
		Preload("SalesOrder").
		First(&delivery, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrRecordNotFound
		}
		return nil, err
	}
	return &delivery, nil
}

// ListDeliveryOrders retrieves a list of delivery orders based on filter
func (r *OrderRepository) ListDeliveryOrders(ctx context.Context, filter *entity.DeliveryOrderFilter) ([]entity.DeliveryOrder, error) {
	var deliveries []entity.DeliveryOrder
	query := r.db.WithContext(ctx)

	if filter != nil {
		if filter.DeliveryNumber != "" {
			query = query.Where("delivery_number LIKE ?", "%"+filter.DeliveryNumber+"%")
		}
		if filter.SalesOrderID != "" {
			query = query.Where("sales_order_id = ?", filter.SalesOrderID)
		}
		if filter.Status != nil {
			query = query.Where("status = ?", *filter.Status)
		}
		if filter.StartDate != nil {
			query = query.Where("delivery_date >= ?", *filter.StartDate)
		}
		if filter.EndDate != nil {
			query = query.Where("delivery_date <= ?", *filter.EndDate)
		}
		if filter.WarehouseID != "" {
			query = query.Where("warehouse_id = ?", filter.WarehouseID)
		}
	}

	if err := query.Order("created_at DESC").Find(&deliveries).Error; err != nil {
		return nil, err
	}
	return deliveries, nil
}

// UpdateDeliveryOrderStatus updates the status of a delivery order
func (r *OrderRepository) UpdateDeliveryOrderStatus(ctx context.Context, id string, status entity.DeliveryOrderStatus) error {
	return r.db.WithContext(ctx).
		Model(&entity.DeliveryOrder{}).
		Where("id = ?", id).
		Update("status", status).
		Error
}

// ProcessDelivery processes a delivery by updating inventory
func (r *OrderRepository) ProcessDelivery(ctx context.Context, deliveryID string, userID string) error {
	tx := r.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return tx.Error
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Get delivery order
	var delivery entity.DeliveryOrder
	if err := tx.Preload("Items").First(&delivery, "id = ?", deliveryID).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Check if delivery is in a valid state
	if delivery.Status != entity.DeliveryOrderStatusPreparing {
		tx.Rollback()
		return ErrInvalidOrderStatus
	}

	// Process each item in the delivery
	for _, item := range delivery.Items {
		// Create stock entry for inventory reduction
		stockEntry := &entity.StockEntry{
			WarehouseID: delivery.WarehouseID,
			ProductID:   item.ItemID,
			Type:        "OUT",
			Quantity:    item.ShippedQuantity,
			Reference:   delivery.DeliveryNumber,
			Note:        fmt.Sprintf("Delivery for Sales Order %s", delivery.SalesOrderID),
			CreatedBy:   userID,
		}

		// Process the stock entry using the inventory repository
		if err := r.inventoryRepo.ProcessStockEntry(ctx, stockEntry, userID); err != nil {
			tx.Rollback()
			return err
		}
	}

	// Update delivery status to in transit
	if err := tx.Model(&entity.DeliveryOrder{}).
		Where("id = ?", deliveryID).
		Update("status", entity.DeliveryOrderStatusInTransit).
		Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

// CreateInvoice creates a new invoice for a sales order
func (r *OrderRepository) CreateInvoice(ctx context.Context, invoice *entity.Invoice) error {
	if invoice.ID == "" {
		invoice.ID = uuid.New().String()
	}

	// Generate invoice number if not provided
	if invoice.InvoiceNumber == "" {
		seq, err := r.sequenceGenerator.NextSequence(ctx, "invoice")
		if err != nil {
			return err
		}
		invoice.InvoiceNumber = fmt.Sprintf("INV-%s-%06d", time.Now().Format("20060102"), seq)
	}

	return r.db.WithContext(ctx).Create(invoice).Error
}

// GetInvoiceByID retrieves an invoice by ID
func (r *OrderRepository) GetInvoiceByID(ctx context.Context, id string) (*entity.Invoice, error) {
	var invoice entity.Invoice
	if err := r.db.WithContext(ctx).
		Preload("SalesOrder").
		First(&invoice, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrRecordNotFound
		}
		return nil, err
	}
	return &invoice, nil
}

// ListInvoices retrieves a list of invoices based on filter
func (r *OrderRepository) ListInvoices(ctx context.Context, filter *entity.InvoiceFilter) ([]entity.Invoice, error) {
	var invoices []entity.Invoice
	query := r.db.WithContext(ctx)

	if filter != nil {
		if filter.InvoiceNumber != "" {
			query = query.Where("invoice_number LIKE ?", "%"+filter.InvoiceNumber+"%")
		}
		if filter.SalesOrderID != "" {
			query = query.Where("sales_order_id = ?", filter.SalesOrderID)
		}
		if filter.Status != nil {
			query = query.Where("status = ?", *filter.Status)
		}
		if filter.StartDate != nil {
			query = query.Where("issue_date >= ?", *filter.StartDate)
		}
		if filter.EndDate != nil {
			query = query.Where("issue_date <= ?", *filter.EndDate)
		}
	}

	if err := query.Order("created_at DESC").Find(&invoices).Error; err != nil {
		return nil, err
	}
	return invoices, nil
}

// UpdateInvoiceStatus updates the status of an invoice
func (r *OrderRepository) UpdateInvoiceStatus(ctx context.Context, id string, status entity.InvoiceStatus) error {
	return r.db.WithContext(ctx).
		Model(&entity.Invoice{}).
		Where("id = ?", id).
		Update("status", status).
		Error
}

// CheckStockAvailability checks if there is enough stock for all items in an order
func (r *OrderRepository) CheckStockAvailability(ctx context.Context, warehouseID string, items []entity.SalesOrderItem) (bool, map[string]float64, error) {
	insufficientItems := make(map[string]float64)

	for _, item := range items {
		// Get current inventory for this product in the warehouse
		inventory, err := r.inventoryRepo.GetByProductAndWarehouse(ctx, item.ItemID, warehouseID)
		if err != nil {
			if err == ErrRecordNotFound {
				// No inventory record means zero quantity
				insufficientItems[item.ItemID] = 0
				continue
			}
			return false, nil, err
		}

		// Check if there's enough stock
		if inventory.Quantity < item.Quantity {
			insufficientItems[item.ItemID] = inventory.Quantity
		}
	}

	return len(insufficientItems) == 0, insufficientItems, nil
}

// SequenceGenerator generates sequential numbers for various document types
type SequenceGenerator struct {
	db *gorm.DB
}

// NewSequenceGenerator creates a new SequenceGenerator
func NewSequenceGenerator(db *gorm.DB) *SequenceGenerator {
	return &SequenceGenerator{db: db}
}

// NextSequence gets the next sequence number for a given sequence type
func (sg *SequenceGenerator) NextSequence(ctx context.Context, sequenceType string) (uint, error) {
	// This is a simplified implementation. In production, you might want to use
	// database-specific features for sequences or implement a more robust solution.
	var sequence struct {
		ID    string `gorm:"primaryKey"`
		Value uint   `gorm:"not null"`
	}

	err := sg.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Try to get existing sequence
		result := tx.Raw("SELECT id, value FROM sequences WHERE id = ? FOR UPDATE", sequenceType).Scan(&sequence)
		if result.Error != nil && result.Error != gorm.ErrRecordNotFound {
			return result.Error
		}

		if result.RowsAffected == 0 {
			// Create new sequence starting at 1
			sequence.ID = sequenceType
			sequence.Value = 1
			if err := tx.Exec("INSERT INTO sequences (id, value) VALUES (?, ?)", sequenceType, sequence.Value).Error; err != nil {
				return err
			}
		} else {
			// Increment existing sequence
			sequence.Value++
			if err := tx.Exec("UPDATE sequences SET value = ? WHERE id = ?", sequence.Value, sequenceType).Error; err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		return 0, err
	}

	return sequence.Value, nil
}

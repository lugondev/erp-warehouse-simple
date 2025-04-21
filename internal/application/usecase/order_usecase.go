package usecase

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/lugondev/erp-warehouse-simple/internal/domain/entity"
	"github.com/lugondev/erp-warehouse-simple/internal/infrastructure/repository"
)

var (
	ErrInvalidOrderStatus = errors.New("invalid order status for this operation")
	ErrInsufficientStock  = errors.New("insufficient stock for order items")
)

// OrderUseCase handles business logic for sales orders and delivery orders
type OrderUseCase struct {
	orderRepo  *repository.OrderRepository
	stocksRepo *repository.StocksRepository
}

// NewOrderUseCase creates a new OrderUseCase
func NewOrderUseCase(orderRepo *repository.OrderRepository, stocksRepo *repository.StocksRepository) *OrderUseCase {
	return &OrderUseCase{
		orderRepo:  orderRepo,
		stocksRepo: stocksRepo,
	}
}

// CreateSalesOrder creates a new sales order with stock validation
func (u *OrderUseCase) CreateSalesOrder(ctx context.Context, order *entity.SalesOrder, warehouseID string, userID string) error {
	// Validate order items
	if len(order.Items) == 0 {
		return repository.ErrInvalidData
	}

	// Check stock availability
	available, insufficientItems, err := u.orderRepo.CheckStockAvailability(ctx, warehouseID, order.Items)
	if err != nil {
		return err
	}

	if !available {
		// Format error message with insufficient items
		itemIDs := ""
		for id, qty := range insufficientItems {
			if itemIDs != "" {
				itemIDs += ", "
			}
			itemIDs += fmt.Sprintf("%s (available: %.2f)", id, qty)
		}
		return fmt.Errorf("%w: %s", ErrInsufficientStock, itemIDs)
	}

	// Set initial status and created by
	order.Status = entity.SalesOrderStatusDraft
	createdByID, _ := parseUserID(userID)
	order.CreatedByID = createdByID
	order.OrderDate = time.Now()

	// Calculate totals
	u.calculateOrderTotals(order)

	// Create the order
	return u.orderRepo.CreateSalesOrder(ctx, order)
}

// Rest of the methods remain unchanged as they don't directly use stocksRepo
// Omitted for brevity...

// ConfirmSalesOrder changes a sales order from draft to confirmed status
func (u *OrderUseCase) ConfirmSalesOrder(ctx context.Context, orderID string, userID string) error {
	// Get the order
	order, err := u.orderRepo.GetSalesOrderByID(ctx, orderID)
	if err != nil {
		return err
	}

	// Validate current status
	if order.Status != entity.SalesOrderStatusDraft {
		return ErrInvalidOrderStatus
	}

	// Update status
	return u.orderRepo.UpdateSalesOrderStatus(ctx, orderID, entity.SalesOrderStatusConfirmed)
}

// CreateDeliveryOrder creates a delivery order for a sales order
func (u *OrderUseCase) CreateDeliveryOrder(ctx context.Context, delivery *entity.DeliveryOrder, userID string) error {
	// Get the sales order
	order, err := u.orderRepo.GetSalesOrderByID(ctx, delivery.SalesOrderID)
	if err != nil {
		return err
	}

	// Validate order status
	if order.Status != entity.SalesOrderStatusConfirmed && order.Status != entity.SalesOrderStatusProcessing {
		return ErrInvalidOrderStatus
	}

	// Set initial status and created by
	delivery.Status = entity.DeliveryOrderStatusPending
	createdByID, _ := parseUserID(userID)
	delivery.CreatedByID = createdByID

	// If shipping address not provided, use the one from sales order
	if delivery.ShippingAddress == "" {
		delivery.ShippingAddress = order.ShippingAddress
	}

	// Create the delivery order
	if err := u.orderRepo.CreateDeliveryOrder(ctx, delivery); err != nil {
		return err
	}

	// Update sales order status to processing
	return u.orderRepo.UpdateSalesOrderStatus(ctx, order.ID, entity.SalesOrderStatusProcessing)
}

// PrepareDelivery updates a delivery order status to preparing
func (u *OrderUseCase) PrepareDelivery(ctx context.Context, deliveryID string) error {
	// Get the delivery order
	delivery, err := u.orderRepo.GetDeliveryOrderByID(ctx, deliveryID)
	if err != nil {
		return err
	}

	// Validate current status
	if delivery.Status != entity.DeliveryOrderStatusPending {
		return ErrInvalidOrderStatus
	}

	// Update status
	return u.orderRepo.UpdateDeliveryOrderStatus(ctx, deliveryID, entity.DeliveryOrderStatusPreparing)
}

// ShipDelivery processes a delivery by updating inventory and changing status
func (u *OrderUseCase) ShipDelivery(ctx context.Context, deliveryID string, userID string) error {
	// Process the delivery (this will update inventory)
	if err := u.orderRepo.ProcessDelivery(ctx, deliveryID, userID); err != nil {
		return err
	}

	// Get the delivery order to update the sales order
	delivery, err := u.orderRepo.GetDeliveryOrderByID(ctx, deliveryID)
	if err != nil {
		return err
	}

	// Update sales order status to shipped
	return u.orderRepo.UpdateSalesOrderStatus(ctx, delivery.SalesOrderID, entity.SalesOrderStatusShipped)
}

// CompleteDelivery marks a delivery as delivered
func (u *OrderUseCase) CompleteDelivery(ctx context.Context, deliveryID string) error {
	// Get the delivery order
	delivery, err := u.orderRepo.GetDeliveryOrderByID(ctx, deliveryID)
	if err != nil {
		return err
	}

	// Validate current status
	if delivery.Status != entity.DeliveryOrderStatusInTransit {
		return ErrInvalidOrderStatus
	}

	// Update delivery status
	if err := u.orderRepo.UpdateDeliveryOrderStatus(ctx, deliveryID, entity.DeliveryOrderStatusDelivered); err != nil {
		return err
	}

	// Update sales order status to delivered
	return u.orderRepo.UpdateSalesOrderStatus(ctx, delivery.SalesOrderID, entity.SalesOrderStatusDelivered)
}

// CompleteSalesOrder marks a sales order as completed
func (u *OrderUseCase) CompleteSalesOrder(ctx context.Context, orderID string) error {
	// Get the order
	order, err := u.orderRepo.GetSalesOrderByID(ctx, orderID)
	if err != nil {
		return err
	}

	// Validate current status
	if order.Status != entity.SalesOrderStatusDelivered {
		return ErrInvalidOrderStatus
	}

	// Update status
	return u.orderRepo.UpdateSalesOrderStatus(ctx, orderID, entity.SalesOrderStatusCompleted)
}

// CreateInvoice creates an invoice for a sales order
func (u *OrderUseCase) CreateInvoice(ctx context.Context, invoice *entity.Invoice, userID string) error {
	// Get the sales order
	order, err := u.orderRepo.GetSalesOrderByID(ctx, invoice.SalesOrderID)
	if err != nil {
		return err
	}

	// Validate order status (can create invoice after confirmation)
	if order.Status == entity.SalesOrderStatusDraft || order.Status == entity.SalesOrderStatusCancelled {
		return ErrInvalidOrderStatus
	}

	// Set initial status and created by
	invoice.Status = entity.InvoiceStatusDraft
	createdByID, _ := parseUserID(userID)
	invoice.CreatedByID = createdByID
	invoice.IssueDate = time.Now()

	// If amount not provided, use the one from sales order
	if invoice.Amount == 0 {
		invoice.Amount = order.SubTotal
	}
	if invoice.TaxAmount == 0 {
		invoice.TaxAmount = order.TaxTotal
	}
	if invoice.TotalAmount == 0 {
		invoice.TotalAmount = order.GrandTotal
	}

	// Create the invoice
	return u.orderRepo.CreateInvoice(ctx, invoice)
}

// IssueInvoice changes an invoice from draft to issued status
func (u *OrderUseCase) IssueInvoice(ctx context.Context, invoiceID string) error {
	// Get the invoice
	invoice, err := u.orderRepo.GetInvoiceByID(ctx, invoiceID)
	if err != nil {
		return err
	}

	// Validate current status
	if invoice.Status != entity.InvoiceStatusDraft {
		return ErrInvalidOrderStatus
	}

	// Update status
	return u.orderRepo.UpdateInvoiceStatus(ctx, invoiceID, entity.InvoiceStatusIssued)
}

// PayInvoice marks an invoice as paid
func (u *OrderUseCase) PayInvoice(ctx context.Context, invoiceID string) error {
	// Get the invoice
	invoice, err := u.orderRepo.GetInvoiceByID(ctx, invoiceID)
	if err != nil {
		return err
	}

	// Validate current status
	if invoice.Status != entity.InvoiceStatusIssued && invoice.Status != entity.InvoiceStatusPartial {
		return ErrInvalidOrderStatus
	}

	// Update invoice status
	if err := u.orderRepo.UpdateInvoiceStatus(ctx, invoiceID, entity.InvoiceStatusPaid); err != nil {
		return err
	}

	// Get the sales order to update its payment status
	order, err := u.orderRepo.GetSalesOrderByID(ctx, invoice.SalesOrderID)
	if err != nil {
		return err
	}

	// Update sales order payment status
	order.PaymentStatus = entity.PaymentStatusPaid
	return u.orderRepo.UpdateSalesOrder(ctx, order)
}

// CancelSalesOrder cancels a sales order
func (u *OrderUseCase) CancelSalesOrder(ctx context.Context, orderID string) error {
	// Get the order
	order, err := u.orderRepo.GetSalesOrderByID(ctx, orderID)
	if err != nil {
		return err
	}

	// Can only cancel orders that are not completed or already cancelled
	if order.Status == entity.SalesOrderStatusCompleted || order.Status == entity.SalesOrderStatusCancelled {
		return ErrInvalidOrderStatus
	}

	// If order has deliveries, check if they can be cancelled
	for _, delivery := range order.DeliveryOrders {
		if delivery.Status == entity.DeliveryOrderStatusDelivered ||
			delivery.Status == entity.DeliveryOrderStatusInTransit {
			return errors.New("cannot cancel order with deliveries in transit or delivered")
		}

		// Cancel any pending or preparing deliveries
		if delivery.Status == entity.DeliveryOrderStatusPending ||
			delivery.Status == entity.DeliveryOrderStatusPreparing {
			if err := u.orderRepo.UpdateDeliveryOrderStatus(ctx, delivery.ID, entity.DeliveryOrderStatusCancelled); err != nil {
				return err
			}
		}
	}

	// Cancel any draft invoices
	for _, invoice := range order.Invoices {
		if invoice.Status == entity.InvoiceStatusDraft {
			if err := u.orderRepo.UpdateInvoiceStatus(ctx, invoice.ID, entity.InvoiceStatusCancelled); err != nil {
				return err
			}
		}
	}

	// Update order status
	return u.orderRepo.UpdateSalesOrderStatus(ctx, orderID, entity.SalesOrderStatusCancelled)
}

// GetSalesOrder retrieves a sales order by ID
func (u *OrderUseCase) GetSalesOrder(ctx context.Context, id string) (*entity.SalesOrder, error) {
	return u.orderRepo.GetSalesOrderByID(ctx, id)
}

// ListSalesOrders retrieves a list of sales orders based on filter
func (u *OrderUseCase) ListSalesOrders(ctx context.Context, filter *entity.SalesOrderFilter) ([]entity.SalesOrder, error) {
	return u.orderRepo.ListSalesOrders(ctx, filter)
}

// GetDeliveryOrder retrieves a delivery order by ID
func (u *OrderUseCase) GetDeliveryOrder(ctx context.Context, id string) (*entity.DeliveryOrder, error) {
	return u.orderRepo.GetDeliveryOrderByID(ctx, id)
}

// ListDeliveryOrders retrieves a list of delivery orders based on filter
func (u *OrderUseCase) ListDeliveryOrders(ctx context.Context, filter *entity.DeliveryOrderFilter) ([]entity.DeliveryOrder, error) {
	return u.orderRepo.ListDeliveryOrders(ctx, filter)
}

// GetInvoice retrieves an invoice by ID
func (u *OrderUseCase) GetInvoice(ctx context.Context, id string) (*entity.Invoice, error) {
	return u.orderRepo.GetInvoiceByID(ctx, id)
}

// ListInvoices retrieves a list of invoices based on filter
func (u *OrderUseCase) ListInvoices(ctx context.Context, filter *entity.InvoiceFilter) ([]entity.Invoice, error) {
	return u.orderRepo.ListInvoices(ctx, filter)
}

// Helper function to parse user ID from string to uint
func parseUserID(userID string) (uint, error) {
	var id uint
	_, err := fmt.Sscanf(userID, "%d", &id)
	return id, err
}

// calculateOrderTotals calculates the totals for a sales order
func (u *OrderUseCase) calculateOrderTotals(order *entity.SalesOrder) {
	var subTotal, taxTotal, discountTotal float64

	for i, item := range order.Items {
		// Calculate item total
		lineTotal := item.Quantity * item.UnitPrice

		// Apply discount
		discountAmount := lineTotal * (item.Discount / 100)
		lineTotal -= discountAmount
		discountTotal += discountAmount

		// Calculate tax
		taxAmount := lineTotal * (item.TaxRate / 100)
		taxTotal += taxAmount

		// Update item total
		order.Items[i].TaxAmount = taxAmount
		order.Items[i].TotalPrice = lineTotal + taxAmount

		// Add to subtotal
		subTotal += lineTotal
	}

	// Set order totals
	order.SubTotal = subTotal
	order.TaxTotal = taxTotal
	order.DiscountTotal = discountTotal
	order.GrandTotal = subTotal + taxTotal
}

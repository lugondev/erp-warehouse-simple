package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/lugondev/erp-warehouse-simple/internal/domain/entity"
	"github.com/lugondev/erp-warehouse-simple/internal/infrastructure/repository"
)

var (
	ErrInvalidPurchaseRequest = errors.New("invalid purchase request")
	ErrInvalidPurchaseOrder   = errors.New("invalid purchase order")
	ErrInvalidPurchaseReceipt = errors.New("invalid purchase receipt")
	ErrInvalidPurchasePayment = errors.New("invalid purchase payment")
	ErrOrderAlreadyReceived   = errors.New("purchase order already fully received")
	ErrOrderNotApproved       = errors.New("purchase order not approved")
	ErrOrderNotReceived       = errors.New("purchase order not received")
)

type PurchaseUseCase struct {
	purchaseRepo *repository.PurchaseRepository
	stocksRepo   *repository.StocksRepository
	vendorRepo   *repository.VendorRepository
	skuRepo      *repository.SKURepository
}

func NewPurchaseUseCase(
	purchaseRepo *repository.PurchaseRepository,
	stocksRepo *repository.StocksRepository,
	vendorRepo *repository.VendorRepository,
	skuRepo *repository.SKURepository,
) *PurchaseUseCase {
	return &PurchaseUseCase{
		purchaseRepo: purchaseRepo,
		stocksRepo:   stocksRepo,
		vendorRepo:   vendorRepo,
		skuRepo:      skuRepo,
	}
}

// Purchase Request methods

// CreatePurchaseRequest creates a new purchase request
func (u *PurchaseUseCase) CreatePurchaseRequest(ctx context.Context, request *entity.PurchaseRequest) error {
	if err := u.validatePurchaseRequest(request); err != nil {
		return err
	}

	request.Status = entity.PurchaseRequestStatusDraft
	request.RequestDate = time.Now()

	return u.purchaseRepo.CreatePurchaseRequest(ctx, request)
}

// GetPurchaseRequest gets a purchase request by ID
func (u *PurchaseUseCase) GetPurchaseRequest(ctx context.Context, id string) (*entity.PurchaseRequest, error) {
	return u.purchaseRepo.GetPurchaseRequestByID(ctx, id)
}

// UpdatePurchaseRequest updates a purchase request
func (u *PurchaseUseCase) UpdatePurchaseRequest(ctx context.Context, request *entity.PurchaseRequest) error {
	existingRequest, err := u.purchaseRepo.GetPurchaseRequestByID(ctx, request.ID)
	if err != nil {
		return err
	}

	// Cannot update if already ordered
	if existingRequest.Status == entity.PurchaseRequestStatusOrdered {
		return errors.New("cannot update a purchase request that has been ordered")
	}

	if err := u.validatePurchaseRequest(request); err != nil {
		return err
	}

	return u.purchaseRepo.UpdatePurchaseRequest(ctx, request)
}

// DeletePurchaseRequest deletes a purchase request
func (u *PurchaseUseCase) DeletePurchaseRequest(ctx context.Context, id string) error {
	request, err := u.purchaseRepo.GetPurchaseRequestByID(ctx, id)
	if err != nil {
		return err
	}

	// Cannot delete if not in draft or rejected status
	if request.Status != entity.PurchaseRequestStatusDraft && request.Status != entity.PurchaseRequestStatusRejected {
		return errors.New("can only delete purchase requests in draft or rejected status")
	}

	return u.purchaseRepo.DeletePurchaseRequest(ctx, id)
}

// ListPurchaseRequests lists purchase requests with filters
func (u *PurchaseUseCase) ListPurchaseRequests(ctx context.Context, filter *entity.PurchaseRequestFilter, page, pageSize int) ([]entity.PurchaseRequest, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	return u.purchaseRepo.ListPurchaseRequests(ctx, filter, page, pageSize)
}

// SubmitPurchaseRequest submits a purchase request for approval
func (u *PurchaseUseCase) SubmitPurchaseRequest(ctx context.Context, id string) error {
	request, err := u.purchaseRepo.GetPurchaseRequestByID(ctx, id)
	if err != nil {
		return err
	}

	if request.Status != entity.PurchaseRequestStatusDraft {
		return errors.New("only draft purchase requests can be submitted")
	}

	request.Status = entity.PurchaseRequestStatusSubmitted

	return u.purchaseRepo.UpdatePurchaseRequest(ctx, request)
}

// ApprovePurchaseRequest approves a purchase request
func (u *PurchaseUseCase) ApprovePurchaseRequest(ctx context.Context, id string, approverID uint, notes string) error {
	request, err := u.purchaseRepo.GetPurchaseRequestByID(ctx, id)
	if err != nil {
		return err
	}

	if request.Status != entity.PurchaseRequestStatusSubmitted {
		return errors.New("only submitted purchase requests can be approved")
	}

	now := time.Now()
	request.Status = entity.PurchaseRequestStatusApproved
	request.ApproverID = &approverID
	request.ApprovalDate = &now
	request.ApprovalNotes = notes

	return u.purchaseRepo.UpdatePurchaseRequest(ctx, request)
}

// RejectPurchaseRequest rejects a purchase request
func (u *PurchaseUseCase) RejectPurchaseRequest(ctx context.Context, id string, approverID uint, notes string) error {
	request, err := u.purchaseRepo.GetPurchaseRequestByID(ctx, id)
	if err != nil {
		return err
	}

	if request.Status != entity.PurchaseRequestStatusSubmitted {
		return errors.New("only submitted purchase requests can be rejected")
	}

	now := time.Now()
	request.Status = entity.PurchaseRequestStatusRejected
	request.ApproverID = &approverID
	request.ApprovalDate = &now
	request.ApprovalNotes = notes

	return u.purchaseRepo.UpdatePurchaseRequest(ctx, request)
}

// Purchase Order methods

// CreatePurchaseOrder creates a new purchase order
func (u *PurchaseUseCase) CreatePurchaseOrder(ctx context.Context, order *entity.PurchaseOrder) error {
	if err := u.validatePurchaseOrder(order); err != nil {
		return err
	}

	// Verify vendor exists
	if _, err := u.vendorRepo.FindByID(ctx, order.VendorID); err != nil {
		return err
	}

	order.Status = entity.PurchaseOrderStatusDraft
	order.PaymentStatus = entity.PaymentStatusPending
	order.OrderDate = time.Now()

	return u.purchaseRepo.CreatePurchaseOrder(ctx, order)
}

// GetPurchaseOrder gets a purchase order by ID
func (u *PurchaseUseCase) GetPurchaseOrder(ctx context.Context, id string) (*entity.PurchaseOrder, error) {
	return u.purchaseRepo.GetPurchaseOrderByID(ctx, id)
}

// UpdatePurchaseOrder updates a purchase order
func (u *PurchaseUseCase) UpdatePurchaseOrder(ctx context.Context, order *entity.PurchaseOrder) error {
	existingOrder, err := u.purchaseRepo.GetPurchaseOrderByID(ctx, order.ID)
	if err != nil {
		return err
	}

	// Cannot update if not in draft or submitted status
	if existingOrder.Status != entity.PurchaseOrderStatusDraft && existingOrder.Status != entity.PurchaseOrderStatusSubmitted {
		return errors.New("can only update purchase orders in draft or submitted status")
	}

	if err := u.validatePurchaseOrder(order); err != nil {
		return err
	}

	return u.purchaseRepo.UpdatePurchaseOrder(ctx, order)
}

// DeletePurchaseOrder deletes a purchase order
func (u *PurchaseUseCase) DeletePurchaseOrder(ctx context.Context, id string) error {
	order, err := u.purchaseRepo.GetPurchaseOrderByID(ctx, id)
	if err != nil {
		return err
	}

	// Cannot delete if not in draft status
	if order.Status != entity.PurchaseOrderStatusDraft {
		return errors.New("can only delete purchase orders in draft status")
	}

	return u.purchaseRepo.DeletePurchaseOrder(ctx, id)
}

// ListPurchaseOrders lists purchase orders with filters
func (u *PurchaseUseCase) ListPurchaseOrders(ctx context.Context, filter *entity.PurchaseOrderFilter, page, pageSize int) ([]entity.PurchaseOrder, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	return u.purchaseRepo.ListPurchaseOrders(ctx, filter, page, pageSize)
}

// SubmitPurchaseOrder submits a purchase order for approval
func (u *PurchaseUseCase) SubmitPurchaseOrder(ctx context.Context, id string) error {
	order, err := u.purchaseRepo.GetPurchaseOrderByID(ctx, id)
	if err != nil {
		return err
	}

	if order.Status != entity.PurchaseOrderStatusDraft {
		return errors.New("only draft purchase orders can be submitted")
	}

	order.Status = entity.PurchaseOrderStatusSubmitted

	return u.purchaseRepo.UpdatePurchaseOrder(ctx, order)
}

// ApprovePurchaseOrder approves a purchase order
func (u *PurchaseUseCase) ApprovePurchaseOrder(ctx context.Context, id string, approverID uint) error {
	order, err := u.purchaseRepo.GetPurchaseOrderByID(ctx, id)
	if err != nil {
		return err
	}

	if order.Status != entity.PurchaseOrderStatusSubmitted {
		return errors.New("only submitted purchase orders can be approved")
	}

	now := time.Now()
	order.Status = entity.PurchaseOrderStatusApproved
	order.ApprovedByID = &approverID
	order.ApprovalDate = &now

	return u.purchaseRepo.UpdatePurchaseOrder(ctx, order)
}

// SendPurchaseOrder marks a purchase order as sent to vendor
func (u *PurchaseUseCase) SendPurchaseOrder(ctx context.Context, id string) error {
	order, err := u.purchaseRepo.GetPurchaseOrderByID(ctx, id)
	if err != nil {
		return err
	}

	if order.Status != entity.PurchaseOrderStatusApproved {
		return errors.New("only approved purchase orders can be sent")
	}

	order.Status = entity.PurchaseOrderStatusSent

	return u.purchaseRepo.UpdatePurchaseOrder(ctx, order)
}

// ConfirmPurchaseOrder marks a purchase order as confirmed by vendor
func (u *PurchaseUseCase) ConfirmPurchaseOrder(ctx context.Context, id string) error {
	order, err := u.purchaseRepo.GetPurchaseOrderByID(ctx, id)
	if err != nil {
		return err
	}

	if order.Status != entity.PurchaseOrderStatusSent {
		return errors.New("only sent purchase orders can be confirmed")
	}

	order.Status = entity.PurchaseOrderStatusConfirmed

	return u.purchaseRepo.UpdatePurchaseOrder(ctx, order)
}

// CancelPurchaseOrder cancels a purchase order
func (u *PurchaseUseCase) CancelPurchaseOrder(ctx context.Context, id string) error {
	order, err := u.purchaseRepo.GetPurchaseOrderByID(ctx, id)
	if err != nil {
		return err
	}

	// Cannot cancel if already received or closed
	if order.Status == entity.PurchaseOrderStatusReceived || order.Status == entity.PurchaseOrderStatusClosed {
		return errors.New("cannot cancel purchase orders that are received or closed")
	}

	order.Status = entity.PurchaseOrderStatusCancelled

	return u.purchaseRepo.UpdatePurchaseOrder(ctx, order)
}

// ClosePurchaseOrder closes a purchase order
func (u *PurchaseUseCase) ClosePurchaseOrder(ctx context.Context, id string) error {
	order, err := u.purchaseRepo.GetPurchaseOrderByID(ctx, id)
	if err != nil {
		return err
	}

	// Can only close if received and paid
	if order.Status != entity.PurchaseOrderStatusReceived {
		return errors.New("only received purchase orders can be closed")
	}

	if order.PaymentStatus != entity.PaymentStatusPaid {
		return errors.New("only fully paid purchase orders can be closed")
	}

	order.Status = entity.PurchaseOrderStatusClosed

	return u.purchaseRepo.UpdatePurchaseOrder(ctx, order)
}

// CreatePurchaseOrderFromRequest creates a purchase order from a purchase request
func (u *PurchaseUseCase) CreatePurchaseOrderFromRequest(ctx context.Context, requestID string, vendorID uint, createdByID uint) (*entity.PurchaseOrder, error) {
	request, err := u.purchaseRepo.GetPurchaseRequestByID(ctx, requestID)
	if err != nil {
		return nil, err
	}

	if request.Status != entity.PurchaseRequestStatusApproved {
		return nil, errors.New("can only create purchase orders from approved purchase requests")
	}

	// Verify vendor exists
	if _, err := u.vendorRepo.FindByID(ctx, vendorID); err != nil {
		return nil, err
	}

	// Create order items from request items
	orderItems := make(entity.PurchaseOrderItems, 0, len(request.Items))
	subTotal := 0.0
	taxTotal := 0.0

	for _, item := range request.Items {
		// Get SKU details
		sku, err := u.skuRepo.GetSKUByID(ctx, item.SKUID)
		if err != nil {
			return nil, err
		}

		// Calculate item totals
		unitPrice := sku.Price
		taxRate := 0.0 // Default tax rate
		taxAmount := unitPrice * item.Quantity * (taxRate / 100)
		totalPrice := (unitPrice * item.Quantity) + taxAmount

		orderItem := entity.PurchaseOrderItem{
			SKUID:       item.SKUID,
			Quantity:    item.Quantity,
			UnitPrice:   unitPrice,
			TaxRate:     taxRate,
			TaxAmount:   taxAmount,
			Discount:    0,
			TotalPrice:  totalPrice,
			Description: item.Description,
		}

		orderItems = append(orderItems, orderItem)
		subTotal += unitPrice * item.Quantity
		taxTotal += taxAmount
	}

	// Create purchase order
	order := &entity.PurchaseOrder{
		VendorID:      vendorID,
		OrderDate:     time.Now(),
		ExpectedDate:  request.RequiredDate,
		Items:         orderItems,
		SubTotal:      subTotal,
		TaxTotal:      taxTotal,
		DiscountTotal: 0,
		GrandTotal:    subTotal + taxTotal,
		CurrencyCode:  request.CurrencyCode,
		Status:        entity.PurchaseOrderStatusDraft,
		PaymentStatus: entity.PaymentStatusPending,
		CreatedByID:   createdByID,
	}

	// Create the order
	if err := u.purchaseRepo.CreatePurchaseOrder(ctx, order); err != nil {
		return nil, err
	}

	// Link the request to the order
	if err := u.purchaseRepo.LinkPurchaseRequestToOrder(ctx, requestID, order.ID); err != nil {
		return nil, err
	}

	return order, nil
}

// Purchase Receipt methods

// CreatePurchaseReceipt creates a new purchase receipt
func (u *PurchaseUseCase) CreatePurchaseReceipt(ctx context.Context, receipt *entity.PurchaseReceipt, userID string) error {
	if err := u.validatePurchaseReceipt(receipt); err != nil {
		return err
	}

	// Verify purchase order exists and is in correct status
	order, err := u.purchaseRepo.GetPurchaseOrderByID(ctx, receipt.PurchaseOrderID)
	if err != nil {
		return err
	}

	if order.Status != entity.PurchaseOrderStatusConfirmed &&
		order.Status != entity.PurchaseOrderStatusPartial &&
		order.Status != entity.PurchaseOrderStatusSent {
		return errors.New("purchase order must be sent, confirmed, or partially received to create a receipt")
	}

	// Create receipt
	receipt.ReceiptDate = time.Now()
	if err := u.purchaseRepo.CreatePurchaseReceipt(ctx, receipt); err != nil {
		return err
	}

	// Update inventory for each received item
	for _, item := range receipt.Items {
		if item.ReceivedQuantity <= 0 {
			continue
		}

		// Create stock entry
		stockEntry := &entity.StockEntry{
			StoreID:   receipt.StoreID,
			SKUID:     item.SKUID,
			Type:      "IN",
			Quantity:  item.ReceivedQuantity,
			Reference: receipt.ReceiptNumber,
			Note:      "Purchase receipt",
			CreatedBy: userID,
		}

		if err := u.stocksRepo.ProcessStockEntry(ctx, stockEntry, userID); err != nil {
			return err
		}
	}

	return nil
}

// GetPurchaseReceipt gets a purchase receipt by ID
func (u *PurchaseUseCase) GetPurchaseReceipt(ctx context.Context, id string) (*entity.PurchaseReceipt, error) {
	return u.purchaseRepo.GetPurchaseReceiptByID(ctx, id)
}

// ListPurchaseReceiptsByOrder lists purchase receipts for a purchase order
func (u *PurchaseUseCase) ListPurchaseReceiptsByOrder(ctx context.Context, orderID string) ([]entity.PurchaseReceipt, error) {
	return u.purchaseRepo.ListPurchaseReceiptsByOrderID(ctx, orderID)
}

// Purchase Payment methods

// CreatePurchasePayment creates a new purchase payment
func (u *PurchaseUseCase) CreatePurchasePayment(ctx context.Context, payment *entity.PurchasePayment) error {
	if err := u.validatePurchasePayment(payment); err != nil {
		return err
	}

	// Verify purchase order exists and is in correct status
	order, err := u.purchaseRepo.GetPurchaseOrderByID(ctx, payment.PurchaseOrderID)
	if err != nil {
		return err
	}

	if order.Status != entity.PurchaseOrderStatusReceived && order.Status != entity.PurchaseOrderStatusPartial {
		return errors.New("purchase order must be received or partially received to create a payment")
	}

	// Check if payment would exceed the total amount
	totalPaid, err := u.purchaseRepo.GetTotalPaymentsByOrderID(ctx, payment.PurchaseOrderID)
	if err != nil {
		return err
	}

	if totalPaid+payment.Amount > order.GrandTotal {
		return errors.New("payment amount would exceed the order total")
	}

	payment.PaymentDate = time.Now()
	return u.purchaseRepo.CreatePurchasePayment(ctx, payment)
}

// GetPurchasePayment gets a purchase payment by ID
func (u *PurchaseUseCase) GetPurchasePayment(ctx context.Context, id string) (*entity.PurchasePayment, error) {
	return u.purchaseRepo.GetPurchasePaymentByID(ctx, id)
}

// ListPurchasePaymentsByOrder lists purchase payments for a purchase order
func (u *PurchaseUseCase) ListPurchasePaymentsByOrder(ctx context.Context, orderID string) ([]entity.PurchasePayment, error) {
	return u.purchaseRepo.ListPurchasePaymentsByOrderID(ctx, orderID)
}

// GetPurchaseOrderPaymentSummary gets payment summary for a purchase order
func (u *PurchaseUseCase) GetPurchaseOrderPaymentSummary(ctx context.Context, orderID string) (map[string]interface{}, error) {
	order, err := u.purchaseRepo.GetPurchaseOrderByID(ctx, orderID)
	if err != nil {
		return nil, err
	}

	totalPaid, err := u.purchaseRepo.GetTotalPaymentsByOrderID(ctx, orderID)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"order_id":       order.ID,
		"order_number":   order.OrderNumber,
		"grand_total":    order.GrandTotal,
		"total_paid":     totalPaid,
		"balance_due":    order.GrandTotal - totalPaid,
		"payment_status": order.PaymentStatus,
	}, nil
}

// Validation methods

func (u *PurchaseUseCase) validatePurchaseRequest(request *entity.PurchaseRequest) error {
	if request.RequesterID == 0 {
		return errors.New("requester is required")
	}

	if len(request.Items) == 0 {
		return errors.New("at least one item is required")
	}

	for _, item := range request.Items {
		if item.SKUID == "" {
			return errors.New("SKU ID is required")
		}
		if item.Quantity <= 0 {
			return errors.New("item quantity must be greater than zero")
		}
	}

	return nil
}

func (u *PurchaseUseCase) validatePurchaseOrder(order *entity.PurchaseOrder) error {
	if order.VendorID == 0 {
		return errors.New("vendor is required")
	}

	if order.CreatedByID == 0 {
		return errors.New("created by is required")
	}

	if len(order.Items) == 0 {
		return errors.New("at least one item is required")
	}

	for _, item := range order.Items {
		if item.SKUID == "" {
			return errors.New("SKU ID is required")
		}
		if item.Quantity <= 0 {
			return errors.New("item quantity must be greater than zero")
		}
		if item.UnitPrice < 0 {
			return errors.New("item unit price cannot be negative")
		}
	}

	return nil
}

func (u *PurchaseUseCase) validatePurchaseReceipt(receipt *entity.PurchaseReceipt) error {
	if receipt.PurchaseOrderID == "" {
		return errors.New("purchase order ID is required")
	}

	if receipt.StoreID == "" {
		return errors.New("store ID is required")
	}

	if receipt.ReceivedByID == 0 {
		return errors.New("received by is required")
	}

	if len(receipt.Items) == 0 {
		return errors.New("at least one item is required")
	}

	for _, item := range receipt.Items {
		if item.SKUID == "" {
			return errors.New("SKU ID is required")
		}
		if item.OrderedQuantity <= 0 {
			return errors.New("ordered quantity must be greater than zero")
		}
		if item.ReceivedQuantity < 0 {
			return errors.New("received quantity cannot be negative")
		}
		if item.RejectedQuantity < 0 {
			return errors.New("rejected quantity cannot be negative")
		}
	}

	return nil
}

func (u *PurchaseUseCase) validatePurchasePayment(payment *entity.PurchasePayment) error {
	if payment.PurchaseOrderID == "" {
		return errors.New("purchase order ID is required")
	}

	if payment.Amount <= 0 {
		return errors.New("payment amount must be greater than zero")
	}

	if payment.PaymentMethod == "" {
		return errors.New("payment method is required")
	}

	if payment.CreatedByID == 0 {
		return errors.New("created by is required")
	}

	return nil
}

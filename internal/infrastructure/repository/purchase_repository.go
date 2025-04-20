package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/lugondev/erp-warehouse-simple/internal/domain/entity"
	"gorm.io/gorm"
)

type PurchaseRepository struct {
	db *gorm.DB
}

func NewPurchaseRepository(db *gorm.DB) *PurchaseRepository {
	return &PurchaseRepository{db: db}
}

// Purchase Request methods

// CreatePurchaseRequest creates a new purchase request
func (r *PurchaseRepository) CreatePurchaseRequest(ctx context.Context, request *entity.PurchaseRequest) error {
	if request.ID == "" {
		request.ID = uuid.New().String()
	}
	if request.RequestNumber == "" {
		request.RequestNumber = fmt.Sprintf("PR-%s-%d", time.Now().Format("20060102"), time.Now().UnixNano()%1000)
	}
	return r.db.WithContext(ctx).Create(request).Error
}

// GetPurchaseRequestByID retrieves a purchase request by ID
func (r *PurchaseRepository) GetPurchaseRequestByID(ctx context.Context, id string) (*entity.PurchaseRequest, error) {
	var request entity.PurchaseRequest
	if err := r.db.WithContext(ctx).
		Preload("Requester").
		Preload("Approver").
		Preload("PurchaseOrder").
		First(&request, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrRecordNotFound
		}
		return nil, err
	}
	return &request, nil
}

// UpdatePurchaseRequest updates an existing purchase request
func (r *PurchaseRepository) UpdatePurchaseRequest(ctx context.Context, request *entity.PurchaseRequest) error {
	return r.db.WithContext(ctx).Save(request).Error
}

// DeletePurchaseRequest deletes a purchase request
func (r *PurchaseRepository) DeletePurchaseRequest(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&entity.PurchaseRequest{}, "id = ?", id).Error
}

// ListPurchaseRequests retrieves purchase requests with filters
func (r *PurchaseRepository) ListPurchaseRequests(ctx context.Context, filter *entity.PurchaseRequestFilter, page, pageSize int) ([]entity.PurchaseRequest, int64, error) {
	var requests []entity.PurchaseRequest
	var total int64

	query := r.db.WithContext(ctx).Model(&entity.PurchaseRequest{})

	if filter != nil {
		if filter.RequestNumber != "" {
			query = query.Where("request_number LIKE ?", "%"+filter.RequestNumber+"%")
		}
		if filter.RequesterID != nil {
			query = query.Where("requester_id = ?", *filter.RequesterID)
		}
		if filter.Status != nil {
			query = query.Where("status = ?", *filter.Status)
		}
		if filter.StartDate != nil {
			query = query.Where("request_date >= ?", *filter.StartDate)
		}
		if filter.EndDate != nil {
			query = query.Where("request_date <= ?", *filter.EndDate)
		}
		if filter.ItemID != "" {
			query = query.Where("items @> ?", fmt.Sprintf(`[{"item_id": "%s"}]`, filter.ItemID))
		}
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Offset((page - 1) * pageSize).Limit(pageSize).
		Preload("Requester").
		Preload("Approver").
		Order("created_at DESC").
		Find(&requests).Error; err != nil {
		return nil, 0, err
	}

	return requests, total, nil
}

// Purchase Order methods

// CreatePurchaseOrder creates a new purchase order
func (r *PurchaseRepository) CreatePurchaseOrder(ctx context.Context, order *entity.PurchaseOrder) error {
	if order.ID == "" {
		order.ID = uuid.New().String()
	}
	if order.OrderNumber == "" {
		order.OrderNumber = fmt.Sprintf("PO-%s-%d", time.Now().Format("20060102"), time.Now().UnixNano()%1000)
	}
	return r.db.WithContext(ctx).Create(order).Error
}

// GetPurchaseOrderByID retrieves a purchase order by ID
func (r *PurchaseRepository) GetPurchaseOrderByID(ctx context.Context, id string) (*entity.PurchaseOrder, error) {
	var order entity.PurchaseOrder
	if err := r.db.WithContext(ctx).
		Preload("Supplier").
		Preload("CreatedBy").
		Preload("ApprovedBy").
		Preload("PurchaseRequests").
		First(&order, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrRecordNotFound
		}
		return nil, err
	}
	return &order, nil
}

// UpdatePurchaseOrder updates an existing purchase order
func (r *PurchaseRepository) UpdatePurchaseOrder(ctx context.Context, order *entity.PurchaseOrder) error {
	return r.db.WithContext(ctx).Save(order).Error
}

// DeletePurchaseOrder deletes a purchase order
func (r *PurchaseRepository) DeletePurchaseOrder(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&entity.PurchaseOrder{}, "id = ?", id).Error
}

// ListPurchaseOrders retrieves purchase orders with filters
func (r *PurchaseRepository) ListPurchaseOrders(ctx context.Context, filter *entity.PurchaseOrderFilter, page, pageSize int) ([]entity.PurchaseOrder, int64, error) {
	var orders []entity.PurchaseOrder
	var total int64

	query := r.db.WithContext(ctx).Model(&entity.PurchaseOrder{})

	if filter != nil {
		if filter.OrderNumber != "" {
			query = query.Where("order_number LIKE ?", "%"+filter.OrderNumber+"%")
		}
		if filter.SupplierID != nil {
			query = query.Where("supplier_id = ?", *filter.SupplierID)
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
			query = query.Where("items @> ?", fmt.Sprintf(`[{"item_id": "%s"}]`, filter.ItemID))
		}
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Offset((page - 1) * pageSize).Limit(pageSize).
		Preload("Supplier").
		Preload("CreatedBy").
		Order("created_at DESC").
		Find(&orders).Error; err != nil {
		return nil, 0, err
	}

	return orders, total, nil
}

// Purchase Receipt methods

// CreatePurchaseReceipt creates a new purchase receipt
func (r *PurchaseRepository) CreatePurchaseReceipt(ctx context.Context, receipt *entity.PurchaseReceipt) error {
	if receipt.ID == "" {
		receipt.ID = uuid.New().String()
	}
	if receipt.ReceiptNumber == "" {
		receipt.ReceiptNumber = fmt.Sprintf("GRN-%s-%d", time.Now().Format("20060102"), time.Now().UnixNano()%1000)
	}

	tx := r.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return tx.Error
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Create receipt
	if err := tx.Create(receipt).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Update purchase order status
	order, err := r.GetPurchaseOrderByID(ctx, receipt.PurchaseOrderID)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Check if all items are received
	allReceived := true
	for _, orderItem := range order.Items {
		totalReceived := 0.0
		for _, receiptItem := range receipt.Items {
			if receiptItem.ItemID == orderItem.ItemID {
				totalReceived += receiptItem.ReceivedQuantity
			}
		}
		if totalReceived < orderItem.Quantity {
			allReceived = false
			break
		}
	}

	// Update order status
	if allReceived {
		order.Status = entity.PurchaseOrderStatusReceived
	} else {
		order.Status = entity.PurchaseOrderStatusPartial
	}

	if err := tx.Save(order).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

// GetPurchaseReceiptByID retrieves a purchase receipt by ID
func (r *PurchaseRepository) GetPurchaseReceiptByID(ctx context.Context, id string) (*entity.PurchaseReceipt, error) {
	var receipt entity.PurchaseReceipt
	if err := r.db.WithContext(ctx).
		Preload("PurchaseOrder").
		Preload("ReceivedBy").
		First(&receipt, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrRecordNotFound
		}
		return nil, err
	}
	return &receipt, nil
}

// ListPurchaseReceiptsByOrderID retrieves purchase receipts for a purchase order
func (r *PurchaseRepository) ListPurchaseReceiptsByOrderID(ctx context.Context, orderID string) ([]entity.PurchaseReceipt, error) {
	var receipts []entity.PurchaseReceipt
	if err := r.db.WithContext(ctx).
		Where("purchase_order_id = ?", orderID).
		Preload("ReceivedBy").
		Order("receipt_date DESC").
		Find(&receipts).Error; err != nil {
		return nil, err
	}
	return receipts, nil
}

// Purchase Payment methods

// CreatePurchasePayment creates a new purchase payment
func (r *PurchaseRepository) CreatePurchasePayment(ctx context.Context, payment *entity.PurchasePayment) error {
	if payment.ID == "" {
		payment.ID = uuid.New().String()
	}
	if payment.PaymentNumber == "" {
		payment.PaymentNumber = fmt.Sprintf("PAY-%s-%d", time.Now().Format("20060102"), time.Now().UnixNano()%1000)
	}

	tx := r.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return tx.Error
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Create payment
	if err := tx.Create(payment).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Update purchase order payment status
	order, err := r.GetPurchaseOrderByID(ctx, payment.PurchaseOrderID)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Get total payments
	var totalPaid float64
	if err := tx.Model(&entity.PurchasePayment{}).
		Where("purchase_order_id = ?", payment.PurchaseOrderID).
		Select("COALESCE(SUM(amount), 0)").
		Scan(&totalPaid).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Update payment status
	if totalPaid >= order.GrandTotal {
		order.PaymentStatus = entity.PaymentStatusPaid
	} else if totalPaid > 0 {
		order.PaymentStatus = entity.PaymentStatusPartial
	}

	if err := tx.Save(order).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

// GetPurchasePaymentByID retrieves a purchase payment by ID
func (r *PurchaseRepository) GetPurchasePaymentByID(ctx context.Context, id string) (*entity.PurchasePayment, error) {
	var payment entity.PurchasePayment
	if err := r.db.WithContext(ctx).
		Preload("PurchaseOrder").
		Preload("CreatedBy").
		First(&payment, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrRecordNotFound
		}
		return nil, err
	}
	return &payment, nil
}

// ListPurchasePaymentsByOrderID retrieves purchase payments for a purchase order
func (r *PurchaseRepository) ListPurchasePaymentsByOrderID(ctx context.Context, orderID string) ([]entity.PurchasePayment, error) {
	var payments []entity.PurchasePayment
	if err := r.db.WithContext(ctx).
		Where("purchase_order_id = ?", orderID).
		Preload("CreatedBy").
		Order("payment_date DESC").
		Find(&payments).Error; err != nil {
		return nil, err
	}
	return payments, nil
}

// GetTotalPaymentsByOrderID gets the total amount paid for a purchase order
func (r *PurchaseRepository) GetTotalPaymentsByOrderID(ctx context.Context, orderID string) (float64, error) {
	var totalPaid float64
	if err := r.db.WithContext(ctx).Model(&entity.PurchasePayment{}).
		Where("purchase_order_id = ?", orderID).
		Select("COALESCE(SUM(amount), 0)").
		Scan(&totalPaid).Error; err != nil {
		return 0, err
	}
	return totalPaid, nil
}

// LinkPurchaseRequestToOrder links a purchase request to a purchase order
func (r *PurchaseRepository) LinkPurchaseRequestToOrder(ctx context.Context, requestID string, orderID string) error {
	return r.db.WithContext(ctx).
		Model(&entity.PurchaseRequest{}).
		Where("id = ?", requestID).
		Updates(map[string]interface{}{
			"purchase_order_id": orderID,
			"status":            entity.PurchaseRequestStatusOrdered,
		}).Error
}

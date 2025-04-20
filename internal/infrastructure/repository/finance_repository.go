package repository

import (
	"context"
	"time"

	"github.com/lugondev/erp-warehouse-simple/internal/domain/entity"
	"gorm.io/gorm"
)

// FinanceRepository handles database operations for finance invoices and payments
type FinanceRepository struct {
	db *gorm.DB
}

// NewFinanceRepository creates a new finance repository
func NewFinanceRepository(db *gorm.DB) *FinanceRepository {
	return &FinanceRepository{db: db}
}

// CreateInvoice creates a new finance invoice
func (r *FinanceRepository) CreateInvoice(ctx context.Context, invoice *entity.FinanceInvoice) error {
	// Generate invoice number if not provided
	if invoice.InvoiceNumber == "" {
		prefix := "INV"
		if invoice.Type == entity.FinanceSalesInvoice {
			prefix = "SINV"
		} else if invoice.Type == entity.FinancePurchaseInvoice {
			prefix = "PINV"
		}
		invoice.InvoiceNumber = prefix + "-" + time.Now().Format("20060102-150405")
	}

	// Set default status if not provided
	if invoice.Status == "" {
		invoice.Status = entity.FinanceInvoiceDraft
	}

	// Set timestamps
	now := time.Now()
	invoice.CreatedAt = now
	invoice.UpdatedAt = now

	// Calculate amount due
	invoice.AmountDue = invoice.Total - invoice.AmountPaid

	return r.db.WithContext(ctx).Create(invoice).Error
}

// GetInvoiceByID retrieves a finance invoice by ID
func (r *FinanceRepository) GetInvoiceByID(ctx context.Context, id int64) (*entity.FinanceInvoice, error) {
	var invoice entity.FinanceInvoice
	if err := r.db.WithContext(ctx).First(&invoice, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrRecordNotFound
		}
		return nil, err
	}
	return &invoice, nil
}

// GetInvoiceByNumber retrieves a finance invoice by invoice number
func (r *FinanceRepository) GetInvoiceByNumber(ctx context.Context, invoiceNumber string) (*entity.FinanceInvoice, error) {
	var invoice entity.FinanceInvoice
	if err := r.db.WithContext(ctx).First(&invoice, "invoice_number = ?", invoiceNumber).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrRecordNotFound
		}
		return nil, err
	}
	return &invoice, nil
}

// UpdateInvoice updates a finance invoice
func (r *FinanceRepository) UpdateInvoice(ctx context.Context, invoice *entity.FinanceInvoice) error {
	// Update timestamp
	invoice.UpdatedAt = time.Now()

	// Calculate amount due
	invoice.AmountDue = invoice.Total - invoice.AmountPaid

	result := r.db.WithContext(ctx).Save(invoice)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrRecordNotFound
	}
	return nil
}

// UpdateInvoiceStatus updates the status of a finance invoice
func (r *FinanceRepository) UpdateInvoiceStatus(ctx context.Context, id int64, status entity.FinanceInvoiceStatus) error {
	result := r.db.WithContext(ctx).Model(&entity.FinanceInvoice{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"status":     status,
			"updated_at": time.Now(),
		})

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrRecordNotFound
	}
	return nil
}

// UpdateInvoicePayment updates the payment information of a finance invoice
func (r *FinanceRepository) UpdateInvoicePayment(ctx context.Context, id int64, amountPaid float64) error {
	// First get the invoice to calculate the new status
	var invoice entity.FinanceInvoice
	if err := r.db.WithContext(ctx).First(&invoice, "id = ?", id).Error; err != nil {
		return err
	}

	// Determine new status
	var status entity.FinanceInvoiceStatus
	if amountPaid >= invoice.Total {
		status = entity.FinanceInvoicePaid
	} else if amountPaid > 0 {
		status = entity.FinanceInvoicePartiallyPaid
	} else {
		status = invoice.Status
	}

	// Update the invoice
	result := r.db.WithContext(ctx).Model(&entity.FinanceInvoice{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"amount_paid": amountPaid,
			"amount_due":  invoice.Total - amountPaid,
			"status":      status,
			"updated_at":  time.Now(),
		})

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrRecordNotFound
	}
	return nil
}

// ListInvoices lists finance invoices based on filter criteria
func (r *FinanceRepository) ListInvoices(ctx context.Context, filter *entity.FinanceInvoiceFilter) ([]entity.FinanceInvoice, int64, error) {
	var invoices []entity.FinanceInvoice
	var total int64

	query := r.db.WithContext(ctx).Model(&entity.FinanceInvoice{})

	// Apply filters
	if filter.InvoiceNumber != "" {
		query = query.Where("invoice_number LIKE ?", "%"+filter.InvoiceNumber+"%")
	}
	if filter.Type != "" {
		query = query.Where("type = ?", filter.Type)
	}
	if filter.ReferenceID != "" {
		query = query.Where("reference_id = ?", filter.ReferenceID)
	}
	if filter.EntityID != 0 {
		query = query.Where("entity_id = ?", filter.EntityID)
	}
	if filter.EntityType != "" {
		query = query.Where("entity_type = ?", filter.EntityType)
	}
	if filter.Status != "" {
		query = query.Where("status = ?", filter.Status)
	}
	if filter.StartDate != nil {
		query = query.Where("issue_date >= ?", filter.StartDate)
	}
	if filter.EndDate != nil {
		query = query.Where("issue_date <= ?", filter.EndDate)
	}

	// Get total count
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply pagination
	if filter.Page <= 0 {
		filter.Page = 1
	}
	if filter.PageSize <= 0 {
		filter.PageSize = 10
	}
	offset := (filter.Page - 1) * filter.PageSize

	// Get invoices
	if err := query.Order("created_at DESC").Limit(filter.PageSize).Offset(offset).Find(&invoices).Error; err != nil {
		return nil, 0, err
	}

	return invoices, total, nil
}

// CreatePayment creates a new finance payment
func (r *FinanceRepository) CreatePayment(ctx context.Context, payment *entity.FinancePayment) error {
	// Generate payment number if not provided
	if payment.PaymentNumber == "" {
		payment.PaymentNumber = "PAY-" + time.Now().Format("20060102-150405")
	}

	// Set default status if not provided
	if payment.Status == "" {
		payment.Status = entity.FinancePaymentPending
	}

	// Set timestamps
	now := time.Now()
	payment.CreatedAt = now
	payment.UpdatedAt = now

	// Create payment
	if err := r.db.WithContext(ctx).Create(payment).Error; err != nil {
		return err
	}

	// Update invoice payment amount if payment is completed
	if payment.Status == entity.FinancePaymentCompleted {
		// Get current invoice
		invoice, err := r.GetInvoiceByID(ctx, payment.InvoiceID)
		if err != nil {
			return err
		}

		// Update invoice payment amount
		newAmountPaid := invoice.AmountPaid + payment.Amount
		if err := r.UpdateInvoicePayment(ctx, payment.InvoiceID, newAmountPaid); err != nil {
			return err
		}
	}

	return nil
}

// GetPaymentByID retrieves a finance payment by ID
func (r *FinanceRepository) GetPaymentByID(ctx context.Context, id int64) (*entity.FinancePayment, error) {
	var payment entity.FinancePayment
	if err := r.db.WithContext(ctx).First(&payment, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrRecordNotFound
		}
		return nil, err
	}
	return &payment, nil
}

// UpdatePayment updates a finance payment
func (r *FinanceRepository) UpdatePayment(ctx context.Context, payment *entity.FinancePayment) error {
	// First get the current payment to check status change
	currentPayment, err := r.GetPaymentByID(ctx, payment.ID)
	if err != nil {
		return err
	}

	// Update timestamp
	payment.UpdatedAt = time.Now()

	// Update payment
	result := r.db.WithContext(ctx).Save(payment)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrRecordNotFound
	}

	// Handle status change that affects invoice
	if currentPayment.Status != payment.Status || currentPayment.Amount != payment.Amount {
		// Get current invoice
		invoice, err := r.GetInvoiceByID(ctx, payment.InvoiceID)
		if err != nil {
			return err
		}

		// Calculate new amount paid
		newAmountPaid := invoice.AmountPaid

		// If previous status was completed, subtract the old amount
		if currentPayment.Status == entity.FinancePaymentCompleted {
			newAmountPaid -= currentPayment.Amount
		}

		// If new status is completed, add the new amount
		if payment.Status == entity.FinancePaymentCompleted {
			newAmountPaid += payment.Amount
		}

		// Update invoice payment amount
		if err := r.UpdateInvoicePayment(ctx, payment.InvoiceID, newAmountPaid); err != nil {
			return err
		}
	}

	return nil
}

// UpdatePaymentStatus updates the status of a finance payment
func (r *FinanceRepository) UpdatePaymentStatus(ctx context.Context, id int64, status entity.FinancePaymentStatus) error {
	// First get the current payment to check status change
	currentPayment, err := r.GetPaymentByID(ctx, id)
	if err != nil {
		return err
	}

	// Update payment status
	result := r.db.WithContext(ctx).Model(&entity.FinancePayment{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"status":     status,
			"updated_at": time.Now(),
		})

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrRecordNotFound
	}

	// Handle status change that affects invoice
	if currentPayment.Status != status {
		// Get current invoice
		invoice, err := r.GetInvoiceByID(ctx, currentPayment.InvoiceID)
		if err != nil {
			return err
		}

		// Calculate new amount paid
		newAmountPaid := invoice.AmountPaid

		// If previous status was completed, subtract the amount
		if currentPayment.Status == entity.FinancePaymentCompleted {
			newAmountPaid -= currentPayment.Amount
		}

		// If new status is completed, add the amount
		if status == entity.FinancePaymentCompleted {
			newAmountPaid += currentPayment.Amount
		}

		// Update invoice payment amount
		if err := r.UpdateInvoicePayment(ctx, currentPayment.InvoiceID, newAmountPaid); err != nil {
			return err
		}
	}

	return nil
}

// ListPayments lists finance payments based on filter criteria
func (r *FinanceRepository) ListPayments(ctx context.Context, filter *entity.FinancePaymentFilter) ([]entity.FinancePayment, int64, error) {
	var payments []entity.FinancePayment
	var total int64

	query := r.db.WithContext(ctx).Model(&entity.FinancePayment{})

	// Apply filters
	if filter.PaymentNumber != "" {
		query = query.Where("payment_number LIKE ?", "%"+filter.PaymentNumber+"%")
	}
	if filter.InvoiceID != 0 {
		query = query.Where("invoice_id = ?", filter.InvoiceID)
	}
	if filter.InvoiceNumber != "" {
		query = query.Where("invoice_number LIKE ?", "%"+filter.InvoiceNumber+"%")
	}
	if filter.EntityID != 0 {
		query = query.Where("entity_id = ?", filter.EntityID)
	}
	if filter.EntityType != "" {
		query = query.Where("entity_type = ?", filter.EntityType)
	}
	if filter.Status != "" {
		query = query.Where("status = ?", filter.Status)
	}
	if filter.PaymentMethod != "" {
		query = query.Where("payment_method = ?", filter.PaymentMethod)
	}
	if filter.StartDate != nil {
		query = query.Where("payment_date >= ?", filter.StartDate)
	}
	if filter.EndDate != nil {
		query = query.Where("payment_date <= ?", filter.EndDate)
	}

	// Get total count
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply pagination
	if filter.Page <= 0 {
		filter.Page = 1
	}
	if filter.PageSize <= 0 {
		filter.PageSize = 10
	}
	offset := (filter.Page - 1) * filter.PageSize

	// Get payments
	if err := query.Order("created_at DESC").Limit(filter.PageSize).Offset(offset).Find(&payments).Error; err != nil {
		return nil, 0, err
	}

	return payments, total, nil
}

// GetAccountsReceivable gets accounts receivable data
func (r *FinanceRepository) GetAccountsReceivable(ctx context.Context, startDate, endDate *time.Time) ([]entity.FinanceAccountsReceivable, error) {
	var receivables []entity.FinanceAccountsReceivable

	query := `
		SELECT
			entity_id,
			entity_name,
			id as invoice_id,
			invoice_number,
			issue_date as invoice_date,
			due_date,
			total as total_amount,
			amount_paid,
			amount_due,
			CASE
				WHEN due_date < NOW() AND amount_due > 0 THEN EXTRACT(DAY FROM NOW() - due_date)::int
				ELSE 0
			END as days_overdue,
			status,
			(
				SELECT MAX(payment_date)
				FROM finance_payments
				WHERE invoice_id = finance_invoices.id AND status = 'COMPLETED'
			) as last_payment_date
		FROM finance_invoices
		WHERE type = 'SALES'
	`

	args := []interface{}{}
	if startDate != nil {
		query += " AND issue_date >= ?"
		args = append(args, startDate)
	}
	if endDate != nil {
		query += " AND issue_date <= ?"
		args = append(args, endDate)
	}

	query += " ORDER BY due_date ASC"

	if err := r.db.WithContext(ctx).Raw(query, args...).Scan(&receivables).Error; err != nil {
		return nil, err
	}

	return receivables, nil
}

// GetAccountsPayable gets accounts payable data
func (r *FinanceRepository) GetAccountsPayable(ctx context.Context, startDate, endDate *time.Time) ([]entity.FinanceAccountsPayable, error) {
	var payables []entity.FinanceAccountsPayable

	query := `
		SELECT
			entity_id,
			entity_name,
			id as invoice_id,
			invoice_number,
			issue_date as invoice_date,
			due_date,
			total as total_amount,
			amount_paid,
			amount_due,
			CASE
				WHEN due_date < NOW() AND amount_due > 0 THEN EXTRACT(DAY FROM NOW() - due_date)::int
				ELSE 0
			END as days_overdue,
			status,
			(
				SELECT MAX(payment_date)
				FROM finance_payments
				WHERE invoice_id = finance_invoices.id AND status = 'COMPLETED'
			) as last_payment_date
		FROM finance_invoices
		WHERE type = 'PURCHASE'
	`

	args := []interface{}{}
	if startDate != nil {
		query += " AND issue_date >= ?"
		args = append(args, startDate)
	}
	if endDate != nil {
		query += " AND issue_date <= ?"
		args = append(args, endDate)
	}

	query += " ORDER BY due_date ASC"

	if err := r.db.WithContext(ctx).Raw(query, args...).Scan(&payables).Error; err != nil {
		return nil, err
	}

	return payables, nil
}

// GetFinanceReport generates a finance report for the specified period
func (r *FinanceRepository) GetFinanceReport(ctx context.Context, startDate, endDate time.Time) (*entity.FinanceReport, error) {
	var report entity.FinanceReport
	report.StartDate = startDate
	report.EndDate = endDate

	query := `
		WITH sales_data AS (
			SELECT
				COALESCE(SUM(total), 0) as total_revenue,
				COALESCE(SUM(tax_total), 0) as sales_tax
			FROM finance_invoices
			WHERE type = 'SALES'
			AND issue_date BETWEEN ? AND ?
		),
		purchase_data AS (
			SELECT
				COALESCE(SUM(total), 0) as total_cost,
				COALESCE(SUM(tax_total), 0) as purchase_tax
			FROM finance_invoices
			WHERE type = 'PURCHASE'
			AND issue_date BETWEEN ? AND ?
		)
		SELECT
			sales_data.total_revenue,
			purchase_data.total_cost,
			sales_data.total_revenue - purchase_data.total_cost as gross_profit,
			sales_data.sales_tax + purchase_data.purchase_tax as total_tax,
			sales_data.total_revenue - purchase_data.total_cost - (sales_data.sales_tax + purchase_data.purchase_tax) as net_profit
		FROM sales_data, purchase_data
	`

	if err := r.db.WithContext(ctx).Raw(query, startDate, endDate, startDate, endDate).Scan(&report).Error; err != nil {
		return nil, err
	}

	return &report, nil
}

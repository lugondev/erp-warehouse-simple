package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/lugondev/erp-warehouse-simple/internal/domain/entity"
	"github.com/lugondev/erp-warehouse-simple/internal/infrastructure/repository"
)

// FinanceUseCase handles business logic for finance operations
type FinanceUseCase struct {
	financeRepo *repository.FinanceRepository
}

// NewFinanceUseCase creates a new finance use case
func NewFinanceUseCase(financeRepo *repository.FinanceRepository) *FinanceUseCase {
	return &FinanceUseCase{
		financeRepo: financeRepo,
	}
}

// CreateInvoice creates a new finance invoice
func (u *FinanceUseCase) CreateInvoice(ctx context.Context, req *entity.CreateFinanceInvoiceRequest, userID int64) (*entity.FinanceInvoice, error) {
	// Calculate totals
	var subtotal, taxTotal, total float64
	var items entity.FinanceInvoiceItems

	for _, item := range req.Items {
		// Calculate item totals
		item.Subtotal = item.Quantity * item.UnitPrice
		item.TaxAmount = item.Subtotal * (item.TaxRate / 100)
		item.Total = item.Subtotal + item.TaxAmount

		// Add to invoice totals
		subtotal += item.Subtotal
		taxTotal += item.TaxAmount
		total += item.Total

		items = append(items, item)
	}

	// Apply discount
	total -= req.DiscountAmount

	// Create invoice entity
	invoice := &entity.FinanceInvoice{
		Type:           req.Type,
		ReferenceID:    req.ReferenceID,
		EntityID:       req.EntityID,
		EntityType:     req.EntityType,
		EntityName:     req.EntityType, // This should be replaced with actual name from customer/supplier
		IssueDate:      req.IssueDate,
		DueDate:        req.DueDate,
		Items:          items,
		Subtotal:       subtotal,
		TaxTotal:       taxTotal,
		DiscountAmount: req.DiscountAmount,
		Total:          total,
		AmountPaid:     0,
		AmountDue:      total,
		Status:         entity.FinanceInvoiceDraft,
		Notes:          req.Notes,
		CreatedBy:      userID,
	}

	// Save invoice
	if err := u.financeRepo.CreateInvoice(ctx, invoice); err != nil {
		return nil, fmt.Errorf("error creating invoice: %w", err)
	}

	return invoice, nil
}

// GetInvoiceByID retrieves a finance invoice by ID
func (u *FinanceUseCase) GetInvoiceByID(ctx context.Context, id int64) (*entity.FinanceInvoice, error) {
	invoice, err := u.financeRepo.GetInvoiceByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("error getting invoice: %w", err)
	}
	return invoice, nil
}

// GetInvoiceByNumber retrieves a finance invoice by invoice number
func (u *FinanceUseCase) GetInvoiceByNumber(ctx context.Context, invoiceNumber string) (*entity.FinanceInvoice, error) {
	invoice, err := u.financeRepo.GetInvoiceByNumber(ctx, invoiceNumber)
	if err != nil {
		return nil, fmt.Errorf("error getting invoice: %w", err)
	}
	return invoice, nil
}

// UpdateInvoice updates a finance invoice
func (u *FinanceUseCase) UpdateInvoice(ctx context.Context, id int64, req *entity.UpdateFinanceInvoiceRequest) (*entity.FinanceInvoice, error) {
	// Get existing invoice
	invoice, err := u.financeRepo.GetInvoiceByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("error getting invoice: %w", err)
	}

	// Check if invoice can be updated
	if invoice.Status == entity.FinanceInvoicePaid || invoice.Status == entity.FinanceInvoiceCancelled {
		return nil, fmt.Errorf("cannot update invoice with status %s", invoice.Status)
	}

	// Update fields
	if req.ReferenceID != "" {
		invoice.ReferenceID = req.ReferenceID
	}
	if !req.DueDate.IsZero() {
		invoice.DueDate = req.DueDate
	}
	if req.Notes != "" {
		invoice.Notes = req.Notes
	}
	if req.Status != "" {
		invoice.Status = req.Status
	}

	// Update items and recalculate totals if provided
	if len(req.Items) > 0 {
		var subtotal, taxTotal, total float64
		var items entity.FinanceInvoiceItems

		for _, item := range req.Items {
			// Calculate item totals
			item.Subtotal = item.Quantity * item.UnitPrice
			item.TaxAmount = item.Subtotal * (item.TaxRate / 100)
			item.Total = item.Subtotal + item.TaxAmount

			// Add to invoice totals
			subtotal += item.Subtotal
			taxTotal += item.TaxAmount
			total += item.Total

			items = append(items, item)
		}

		// Apply discount
		total -= req.DiscountAmount

		// Update invoice
		invoice.Items = items
		invoice.Subtotal = subtotal
		invoice.TaxTotal = taxTotal
		invoice.DiscountAmount = req.DiscountAmount
		invoice.Total = total
		invoice.AmountDue = total - invoice.AmountPaid
	}

	// Save invoice
	if err := u.financeRepo.UpdateInvoice(ctx, invoice); err != nil {
		return nil, fmt.Errorf("error updating invoice: %w", err)
	}

	return invoice, nil
}

// UpdateInvoiceStatus updates the status of a finance invoice
func (u *FinanceUseCase) UpdateInvoiceStatus(ctx context.Context, id int64, status entity.FinanceInvoiceStatus) error {
	// Get existing invoice
	invoice, err := u.financeRepo.GetInvoiceByID(ctx, id)
	if err != nil {
		return fmt.Errorf("error getting invoice: %w", err)
	}

	// Check if status change is valid
	if invoice.Status == entity.FinanceInvoiceCancelled {
		return fmt.Errorf("cannot update status of cancelled invoice")
	}

	if invoice.Status == entity.FinanceInvoicePaid && status != entity.FinanceInvoiceCancelled {
		return fmt.Errorf("cannot change status of paid invoice except to cancelled")
	}

	// Update status
	if err := u.financeRepo.UpdateInvoiceStatus(ctx, id, status); err != nil {
		return fmt.Errorf("error updating invoice status: %w", err)
	}

	return nil
}

// CancelInvoice cancels a finance invoice
func (u *FinanceUseCase) CancelInvoice(ctx context.Context, id int64) error {
	// Get existing invoice
	invoice, err := u.financeRepo.GetInvoiceByID(ctx, id)
	if err != nil {
		return fmt.Errorf("error getting invoice: %w", err)
	}

	// Check if invoice can be cancelled
	if invoice.Status == entity.FinanceInvoiceCancelled {
		return fmt.Errorf("invoice is already cancelled")
	}

	// Update status
	if err := u.financeRepo.UpdateInvoiceStatus(ctx, id, entity.FinanceInvoiceCancelled); err != nil {
		return fmt.Errorf("error cancelling invoice: %w", err)
	}

	return nil
}

// ListInvoices lists finance invoices based on filter criteria
func (u *FinanceUseCase) ListInvoices(ctx context.Context, filter *entity.FinanceInvoiceFilter) ([]entity.FinanceInvoice, int64, error) {
	invoices, total, err := u.financeRepo.ListInvoices(ctx, filter)
	if err != nil {
		return nil, 0, fmt.Errorf("error listing invoices: %w", err)
	}
	return invoices, total, nil
}

// CreatePayment creates a new finance payment
func (u *FinanceUseCase) CreatePayment(ctx context.Context, req *entity.CreateFinancePaymentRequest, userID int64) (*entity.FinancePayment, error) {
	// Get invoice
	invoice, err := u.financeRepo.GetInvoiceByID(ctx, req.InvoiceID)
	if err != nil {
		return nil, fmt.Errorf("error getting invoice: %w", err)
	}

	// Check if invoice can be paid
	if invoice.Status == entity.FinanceInvoiceCancelled {
		return nil, fmt.Errorf("cannot create payment for cancelled invoice")
	}

	// Check if payment amount is valid
	if req.Amount <= 0 {
		return nil, fmt.Errorf("payment amount must be greater than zero")
	}

	// Create payment entity
	payment := &entity.FinancePayment{
		InvoiceID:       invoice.ID,
		InvoiceNumber:   invoice.InvoiceNumber,
		EntityID:        invoice.EntityID,
		EntityType:      invoice.EntityType,
		EntityName:      invoice.EntityName,
		PaymentDate:     req.PaymentDate,
		PaymentMethod:   req.PaymentMethod,
		Amount:          req.Amount,
		Status:          entity.FinancePaymentPending,
		Notes:           req.Notes,
		ReferenceNumber: req.ReferenceNumber,
		CreatedBy:       userID,
	}

	// Save payment
	if err := u.financeRepo.CreatePayment(ctx, payment); err != nil {
		return nil, fmt.Errorf("error creating payment: %w", err)
	}

	return payment, nil
}

// GetPaymentByID retrieves a finance payment by ID
func (u *FinanceUseCase) GetPaymentByID(ctx context.Context, id int64) (*entity.FinancePayment, error) {
	payment, err := u.financeRepo.GetPaymentByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("error getting payment: %w", err)
	}
	return payment, nil
}

// UpdatePayment updates a finance payment
func (u *FinanceUseCase) UpdatePayment(ctx context.Context, id int64, req *entity.UpdateFinancePaymentRequest) (*entity.FinancePayment, error) {
	// Get existing payment
	payment, err := u.financeRepo.GetPaymentByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("error getting payment: %w", err)
	}

	// Check if payment can be updated
	if payment.Status == entity.FinancePaymentCancelled || payment.Status == entity.FinancePaymentRefunded {
		return nil, fmt.Errorf("cannot update payment with status %s", payment.Status)
	}

	// Update fields
	if req.PaymentMethod != "" {
		payment.PaymentMethod = req.PaymentMethod
	}
	if req.Amount > 0 {
		payment.Amount = req.Amount
	}
	if req.Status != "" {
		payment.Status = req.Status
	}
	if req.ReferenceNumber != "" {
		payment.ReferenceNumber = req.ReferenceNumber
	}
	if req.Notes != "" {
		payment.Notes = req.Notes
	}

	// Save payment
	if err := u.financeRepo.UpdatePayment(ctx, payment); err != nil {
		return nil, fmt.Errorf("error updating payment: %w", err)
	}

	return payment, nil
}

// ConfirmPayment confirms a finance payment
func (u *FinanceUseCase) ConfirmPayment(ctx context.Context, id int64) error {
	// Get existing payment
	payment, err := u.financeRepo.GetPaymentByID(ctx, id)
	if err != nil {
		return fmt.Errorf("error getting payment: %w", err)
	}

	// Check if payment can be confirmed
	if payment.Status != entity.FinancePaymentPending {
		return fmt.Errorf("only pending payments can be confirmed")
	}

	// Update status
	if err := u.financeRepo.UpdatePaymentStatus(ctx, id, entity.FinancePaymentCompleted); err != nil {
		return fmt.Errorf("error confirming payment: %w", err)
	}

	return nil
}

// CancelPayment cancels a finance payment
func (u *FinanceUseCase) CancelPayment(ctx context.Context, id int64) error {
	// Get existing payment
	payment, err := u.financeRepo.GetPaymentByID(ctx, id)
	if err != nil {
		return fmt.Errorf("error getting payment: %w", err)
	}

	// Check if payment can be cancelled
	if payment.Status == entity.FinancePaymentCancelled {
		return fmt.Errorf("payment is already cancelled")
	}

	if payment.Status == entity.FinancePaymentRefunded {
		return fmt.Errorf("cannot cancel refunded payment")
	}

	// Update status
	if err := u.financeRepo.UpdatePaymentStatus(ctx, id, entity.FinancePaymentCancelled); err != nil {
		return fmt.Errorf("error cancelling payment: %w", err)
	}

	return nil
}

// RefundPayment refunds a finance payment
func (u *FinanceUseCase) RefundPayment(ctx context.Context, id int64) error {
	// Get existing payment
	payment, err := u.financeRepo.GetPaymentByID(ctx, id)
	if err != nil {
		return fmt.Errorf("error getting payment: %w", err)
	}

	// Check if payment can be refunded
	if payment.Status != entity.FinancePaymentCompleted {
		return fmt.Errorf("only completed payments can be refunded")
	}

	// Update status
	if err := u.financeRepo.UpdatePaymentStatus(ctx, id, entity.FinancePaymentRefunded); err != nil {
		return fmt.Errorf("error refunding payment: %w", err)
	}

	return nil
}

// ListPayments lists finance payments based on filter criteria
func (u *FinanceUseCase) ListPayments(ctx context.Context, filter *entity.FinancePaymentFilter) ([]entity.FinancePayment, int64, error) {
	payments, total, err := u.financeRepo.ListPayments(ctx, filter)
	if err != nil {
		return nil, 0, fmt.Errorf("error listing payments: %w", err)
	}
	return payments, total, nil
}

// GetAccountsReceivable gets accounts receivable data
func (u *FinanceUseCase) GetAccountsReceivable(ctx context.Context, startDate, endDate *time.Time) ([]entity.FinanceAccountsReceivable, error) {
	receivables, err := u.financeRepo.GetAccountsReceivable(ctx, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("error getting accounts receivable: %w", err)
	}
	return receivables, nil
}

// GetAccountsPayable gets accounts payable data
func (u *FinanceUseCase) GetAccountsPayable(ctx context.Context, startDate, endDate *time.Time) ([]entity.FinanceAccountsPayable, error) {
	payables, err := u.financeRepo.GetAccountsPayable(ctx, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("error getting accounts payable: %w", err)
	}
	return payables, nil
}

// GetFinanceReport generates a finance report for the specified period
func (u *FinanceUseCase) GetFinanceReport(ctx context.Context, startDate, endDate time.Time) (*entity.FinanceReport, error) {
	report, err := u.financeRepo.GetFinanceReport(ctx, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("error generating finance report: %w", err)
	}
	return report, nil
}

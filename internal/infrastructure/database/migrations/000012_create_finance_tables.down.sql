-- Drop indexes
DROP INDEX IF EXISTS idx_finance_payments_payment_method;
DROP INDEX IF EXISTS idx_finance_payments_payment_date;
DROP INDEX IF EXISTS idx_finance_payments_status;
DROP INDEX IF EXISTS idx_finance_payments_entity_type;
DROP INDEX IF EXISTS idx_finance_payments_entity_id;
DROP INDEX IF EXISTS idx_finance_payments_invoice_id;
DROP INDEX IF EXISTS idx_finance_invoices_type;
DROP INDEX IF EXISTS idx_finance_invoices_due_date;
DROP INDEX IF EXISTS idx_finance_invoices_issue_date;
DROP INDEX IF EXISTS idx_finance_invoices_status;
DROP INDEX IF EXISTS idx_finance_invoices_entity_type;
DROP INDEX IF EXISTS idx_finance_invoices_entity_id;
-- Drop tables
DROP TABLE IF EXISTS finance_tax_rates;
DROP TABLE IF EXISTS finance_payments;
DROP TABLE IF EXISTS finance_invoices;
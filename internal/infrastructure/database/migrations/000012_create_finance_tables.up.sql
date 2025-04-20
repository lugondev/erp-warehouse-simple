-- Create finance_invoices table
CREATE TABLE IF NOT EXISTS finance_invoices (
	id SERIAL PRIMARY KEY,
	invoice_number VARCHAR(50) NOT NULL UNIQUE,
	type VARCHAR(20) NOT NULL,
	reference_id VARCHAR(100),
	entity_id BIGINT NOT NULL,
	entity_type VARCHAR(20) NOT NULL,
	entity_name VARCHAR(255) NOT NULL,
	issue_date TIMESTAMP NOT NULL,
	due_date TIMESTAMP NOT NULL,
	items JSONB NOT NULL,
	subtotal DECIMAL(15, 2) NOT NULL,
	tax_total DECIMAL(15, 2) NOT NULL,
	discount_amount DECIMAL(15, 2) NOT NULL DEFAULT 0,
	total DECIMAL(15, 2) NOT NULL,
	amount_paid DECIMAL(15, 2) NOT NULL DEFAULT 0,
	amount_due DECIMAL(15, 2) NOT NULL,
	status VARCHAR(20) NOT NULL,
	notes TEXT,
	created_by BIGINT NOT NULL,
	created_at TIMESTAMP NOT NULL DEFAULT NOW(),
	updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);
-- Create finance_payments table
CREATE TABLE IF NOT EXISTS finance_payments (
	id SERIAL PRIMARY KEY,
	payment_number VARCHAR(50) NOT NULL UNIQUE,
	invoice_id BIGINT NOT NULL,
	invoice_number VARCHAR(50) NOT NULL,
	entity_id BIGINT NOT NULL,
	entity_type VARCHAR(20) NOT NULL,
	entity_name VARCHAR(255) NOT NULL,
	payment_date TIMESTAMP NOT NULL,
	payment_method VARCHAR(20) NOT NULL,
	amount DECIMAL(15, 2) NOT NULL,
	status VARCHAR(20) NOT NULL,
	notes TEXT,
	reference_number VARCHAR(100),
	created_by BIGINT NOT NULL,
	created_at TIMESTAMP NOT NULL DEFAULT NOW(),
	updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);
-- Create indexes for better query performance
CREATE INDEX IF NOT EXISTS idx_finance_invoices_entity_id ON finance_invoices(entity_id);
CREATE INDEX IF NOT EXISTS idx_finance_invoices_entity_type ON finance_invoices(entity_type);
CREATE INDEX IF NOT EXISTS idx_finance_invoices_status ON finance_invoices(status);
CREATE INDEX IF NOT EXISTS idx_finance_invoices_issue_date ON finance_invoices(issue_date);
CREATE INDEX IF NOT EXISTS idx_finance_invoices_due_date ON finance_invoices(due_date);
CREATE INDEX IF NOT EXISTS idx_finance_invoices_type ON finance_invoices(type);
CREATE INDEX IF NOT EXISTS idx_finance_payments_invoice_id ON finance_payments(invoice_id);
CREATE INDEX IF NOT EXISTS idx_finance_payments_entity_id ON finance_payments(entity_id);
CREATE INDEX IF NOT EXISTS idx_finance_payments_entity_type ON finance_payments(entity_type);
CREATE INDEX IF NOT EXISTS idx_finance_payments_status ON finance_payments(status);
CREATE INDEX IF NOT EXISTS idx_finance_payments_payment_date ON finance_payments(payment_date);
CREATE INDEX IF NOT EXISTS idx_finance_payments_payment_method ON finance_payments(payment_method);
-- Create tax_rates table for managing different tax rates
CREATE TABLE IF NOT EXISTS finance_tax_rates (
	id SERIAL PRIMARY KEY,
	name VARCHAR(100) NOT NULL,
	rate DECIMAL(5, 2) NOT NULL,
	description TEXT,
	is_default BOOLEAN NOT NULL DEFAULT FALSE,
	is_active BOOLEAN NOT NULL DEFAULT TRUE,
	created_at TIMESTAMP NOT NULL DEFAULT NOW(),
	updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);
-- Insert default tax rates
INSERT INTO finance_tax_rates (name, rate, description, is_default, is_active)
VALUES (
		'Standard Rate',
		10.00,
		'Standard tax rate',
		TRUE,
		TRUE
	),
	(
		'Reduced Rate',
		5.00,
		'Reduced tax rate for certain goods',
		FALSE,
		TRUE
	),
	(
		'Zero Rate',
		0.00,
		'Zero-rated goods and services',
		FALSE,
		TRUE
	);
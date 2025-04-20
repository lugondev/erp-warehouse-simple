-- Create purchase_requests table
CREATE TABLE IF NOT EXISTS purchase_requests (
	id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	request_number VARCHAR(50) NOT NULL UNIQUE,
	requester_id INTEGER NOT NULL REFERENCES users(id),
	request_date TIMESTAMP NOT NULL,
	required_date TIMESTAMP,
	items JSONB NOT NULL,
	reason TEXT,
	status VARCHAR(20) NOT NULL DEFAULT 'DRAFT',
	approver_id INTEGER REFERENCES users(id),
	approval_date TIMESTAMP,
	approval_notes TEXT,
	department_id INTEGER,
	total_estimated DECIMAL(15, 2),
	currency_code VARCHAR(3) DEFAULT 'USD',
	attachment_urls TEXT [],
	purchase_order_id UUID,
	created_at TIMESTAMP NOT NULL DEFAULT NOW(),
	updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);
-- Create purchase_orders table
CREATE TABLE IF NOT EXISTS purchase_orders (
	id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	order_number VARCHAR(50) NOT NULL UNIQUE,
	supplier_id INTEGER NOT NULL REFERENCES suppliers(id),
	order_date TIMESTAMP NOT NULL,
	expected_date TIMESTAMP,
	items JSONB NOT NULL,
	sub_total DECIMAL(15, 2) NOT NULL,
	tax_total DECIMAL(15, 2) DEFAULT 0,
	discount_total DECIMAL(15, 2) DEFAULT 0,
	grand_total DECIMAL(15, 2) NOT NULL,
	currency_code VARCHAR(3) DEFAULT 'USD',
	payment_terms VARCHAR(100),
	status VARCHAR(20) NOT NULL DEFAULT 'DRAFT',
	payment_status VARCHAR(20) NOT NULL DEFAULT 'PENDING',
	shipping_address TEXT,
	shipping_method VARCHAR(100),
	notes TEXT,
	attachment_urls TEXT [],
	created_by_id INTEGER NOT NULL REFERENCES users(id),
	approved_by_id INTEGER REFERENCES users(id),
	approval_date TIMESTAMP,
	created_at TIMESTAMP NOT NULL DEFAULT NOW(),
	updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);
-- Add foreign key to purchase_requests
ALTER TABLE purchase_requests
ADD CONSTRAINT fk_purchase_requests_purchase_order FOREIGN KEY (purchase_order_id) REFERENCES purchase_orders(id);
-- Create purchase_receipts table
CREATE TABLE IF NOT EXISTS purchase_receipts (
	id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	receipt_number VARCHAR(50) NOT NULL UNIQUE,
	purchase_order_id UUID NOT NULL REFERENCES purchase_orders(id),
	receipt_date TIMESTAMP NOT NULL,
	items JSONB NOT NULL,
	warehouse_id UUID NOT NULL REFERENCES warehouses(id),
	received_by_id INTEGER NOT NULL REFERENCES users(id),
	notes TEXT,
	attachment_urls TEXT [],
	created_at TIMESTAMP NOT NULL DEFAULT NOW(),
	updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);
-- Create purchase_payments table
CREATE TABLE IF NOT EXISTS purchase_payments (
	id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	payment_number VARCHAR(50) NOT NULL UNIQUE,
	purchase_order_id UUID NOT NULL REFERENCES purchase_orders(id),
	payment_date TIMESTAMP NOT NULL,
	amount DECIMAL(15, 2) NOT NULL,
	payment_method VARCHAR(50) NOT NULL,
	reference_number VARCHAR(100),
	notes TEXT,
	created_by_id INTEGER NOT NULL REFERENCES users(id),
	created_at TIMESTAMP NOT NULL DEFAULT NOW(),
	updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);
-- Create indexes
CREATE INDEX idx_purchase_requests_requester ON purchase_requests(requester_id);
CREATE INDEX idx_purchase_requests_status ON purchase_requests(status);
CREATE INDEX idx_purchase_requests_purchase_order ON purchase_requests(purchase_order_id);
CREATE INDEX idx_purchase_orders_supplier ON purchase_orders(supplier_id);
CREATE INDEX idx_purchase_orders_status ON purchase_orders(status);
CREATE INDEX idx_purchase_orders_payment_status ON purchase_orders(payment_status);
CREATE INDEX idx_purchase_receipts_purchase_order ON purchase_receipts(purchase_order_id);
CREATE INDEX idx_purchase_receipts_warehouse ON purchase_receipts(warehouse_id);
CREATE INDEX idx_purchase_payments_purchase_order ON purchase_payments(purchase_order_id);
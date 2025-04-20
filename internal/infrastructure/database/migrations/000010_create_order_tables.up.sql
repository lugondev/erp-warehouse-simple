-- Create sequences table if not exists
CREATE TABLE IF NOT EXISTS sequences (
	id VARCHAR(50) PRIMARY KEY,
	value INTEGER NOT NULL
);
-- Create sales_orders table
CREATE TABLE IF NOT EXISTS sales_orders (
	id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	order_number VARCHAR(50) NOT NULL UNIQUE,
	customer_id INTEGER NOT NULL,
	order_date TIMESTAMP NOT NULL,
	items JSONB NOT NULL,
	sub_total DECIMAL(15, 2) NOT NULL,
	tax_total DECIMAL(15, 2) NOT NULL DEFAULT 0,
	discount_total DECIMAL(15, 2) NOT NULL DEFAULT 0,
	grand_total DECIMAL(15, 2) NOT NULL,
	status VARCHAR(20) NOT NULL DEFAULT 'DRAFT',
	payment_method VARCHAR(20),
	payment_status VARCHAR(20) NOT NULL DEFAULT 'PENDING',
	shipping_address TEXT,
	billing_address TEXT,
	notes TEXT,
	created_by_id INTEGER NOT NULL,
	created_at TIMESTAMP NOT NULL DEFAULT NOW(),
	updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
	FOREIGN KEY (customer_id) REFERENCES users(id),
	FOREIGN KEY (created_by_id) REFERENCES users(id)
);
-- Create delivery_orders table
CREATE TABLE IF NOT EXISTS delivery_orders (
	id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	delivery_number VARCHAR(50) NOT NULL UNIQUE,
	sales_order_id UUID NOT NULL,
	delivery_date TIMESTAMP NOT NULL,
	items JSONB NOT NULL,
	shipping_address TEXT NOT NULL,
	status VARCHAR(20) NOT NULL DEFAULT 'PENDING',
	tracking_number VARCHAR(100),
	shipping_method VARCHAR(50),
	warehouse_id UUID NOT NULL,
	notes TEXT,
	created_by_id INTEGER NOT NULL,
	created_at TIMESTAMP NOT NULL DEFAULT NOW(),
	updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
	FOREIGN KEY (sales_order_id) REFERENCES sales_orders(id),
	FOREIGN KEY (warehouse_id) REFERENCES warehouses(id),
	FOREIGN KEY (created_by_id) REFERENCES users(id)
);
-- Create invoices table
CREATE TABLE IF NOT EXISTS invoices (
	id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	invoice_number VARCHAR(50) NOT NULL UNIQUE,
	sales_order_id UUID NOT NULL,
	issue_date TIMESTAMP NOT NULL,
	due_date TIMESTAMP NOT NULL,
	amount DECIMAL(15, 2) NOT NULL,
	tax_amount DECIMAL(15, 2) NOT NULL DEFAULT 0,
	total_amount DECIMAL(15, 2) NOT NULL,
	status VARCHAR(20) NOT NULL DEFAULT 'DRAFT',
	notes TEXT,
	created_by_id INTEGER NOT NULL,
	created_at TIMESTAMP NOT NULL DEFAULT NOW(),
	updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
	FOREIGN KEY (sales_order_id) REFERENCES sales_orders(id),
	FOREIGN KEY (created_by_id) REFERENCES users(id)
);
-- Create indexes
CREATE INDEX idx_sales_orders_customer_id ON sales_orders(customer_id);
CREATE INDEX idx_sales_orders_status ON sales_orders(status);
CREATE INDEX idx_sales_orders_payment_status ON sales_orders(payment_status);
CREATE INDEX idx_sales_orders_order_date ON sales_orders(order_date);
CREATE INDEX idx_delivery_orders_sales_order_id ON delivery_orders(sales_order_id);
CREATE INDEX idx_delivery_orders_status ON delivery_orders(status);
CREATE INDEX idx_delivery_orders_warehouse_id ON delivery_orders(warehouse_id);
CREATE INDEX idx_delivery_orders_delivery_date ON delivery_orders(delivery_date);
CREATE INDEX idx_invoices_sales_order_id ON invoices(sales_order_id);
CREATE INDEX idx_invoices_status ON invoices(status);
CREATE INDEX idx_invoices_issue_date ON invoices(issue_date);
CREATE INDEX idx_invoices_due_date ON invoices(due_date);
-- Add initial sequence values
INSERT INTO sequences (id, value)
VALUES ('sales_order', 1) ON CONFLICT (id) DO NOTHING;
INSERT INTO sequences (id, value)
VALUES ('delivery_order', 1) ON CONFLICT (id) DO NOTHING;
INSERT INTO sequences (id, value)
VALUES ('invoice', 1) ON CONFLICT (id) DO NOTHING;
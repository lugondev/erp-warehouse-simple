CREATE TABLE IF NOT EXISTS vendors (
	id SERIAL PRIMARY KEY,
	code VARCHAR(50) UNIQUE NOT NULL,
	name VARCHAR(255) NOT NULL,
	type VARCHAR(50),
	address TEXT,
	country VARCHAR(100),
	email VARCHAR(255),
	phone VARCHAR(50),
	website VARCHAR(255),
	tax_id VARCHAR(100),
	payment_method VARCHAR(100),
	payment_days INTEGER,
	currency VARCHAR(10),
	rating DECIMAL(3, 2) DEFAULT 0,
	created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
CREATE TABLE IF NOT EXISTS vendor_skus (
	id SERIAL PRIMARY KEY,
	code VARCHAR(50) UNIQUE NOT NULL,
	name VARCHAR(255) NOT NULL,
	description TEXT,
	unit_price DECIMAL(15, 2),
	currency VARCHAR(10),
	created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
CREATE TABLE IF NOT EXISTS vendor_sku_mappings (
	vendor_id INTEGER REFERENCES vendors(id) ON DELETE CASCADE,
	sku_id INTEGER REFERENCES vendor_skus(id) ON DELETE CASCADE,
	PRIMARY KEY (vendor_id, sku_id)
);
CREATE TABLE IF NOT EXISTS contracts (
	id SERIAL PRIMARY KEY,
	vendor_id INTEGER REFERENCES vendors(id) ON DELETE CASCADE,
	contract_no VARCHAR(100) UNIQUE NOT NULL,
	start_date TIMESTAMP WITH TIME ZONE,
	end_date TIMESTAMP WITH TIME ZONE,
	terms TEXT,
	status VARCHAR(50),
	created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
CREATE TABLE IF NOT EXISTS vendor_ratings (
	id SERIAL PRIMARY KEY,
	vendor_id INTEGER REFERENCES vendors(id) ON DELETE CASCADE,
	score DECIMAL(3, 2) NOT NULL,
	category VARCHAR(100),
	comment TEXT,
	rated_by INTEGER REFERENCES users(id),
	created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
-- Create indexes for better query performance
CREATE INDEX idx_vendors_code ON vendors(code);
CREATE INDEX idx_vendors_name ON vendors(name);
CREATE INDEX idx_vendor_skus_code ON vendor_skus(code);
CREATE INDEX idx_contracts_vendor ON contracts(vendor_id);
CREATE INDEX idx_contracts_dates ON contracts(start_date, end_date);
CREATE INDEX idx_vendor_ratings ON vendor_ratings(vendor_id);
CREATE INDEX idx_vendor_sku_mappings_vendor ON vendor_sku_mappings(vendor_id);
CREATE INDEX idx_vendor_sku_mappings_sku ON vendor_sku_mappings(sku_id);
CREATE TYPE sku_status AS ENUM ('ACTIVE', 'INACTIVE', 'ARCHIVED');
CREATE TABLE IF NOT EXISTS sku_categories (
	id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	name VARCHAR(255) NOT NULL,
	description TEXT,
	parent_id UUID REFERENCES sku_categories(id),
	created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
	UNIQUE(name)
);
CREATE TABLE IF NOT EXISTS skus (
	id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	sku_code VARCHAR(100) UNIQUE NOT NULL,
	name VARCHAR(255) NOT NULL,
	description TEXT,
	unit_of_measure VARCHAR(50) NOT NULL,
	price DECIMAL(15, 2) DEFAULT 0,
	category_id UUID REFERENCES sku_categories(id),
	technical_specs JSONB,
	primary_vendor_id INTEGER REFERENCES vendors(id),
	backup_vendor_id INTEGER REFERENCES vendors(id),
	image_url VARCHAR(255),
	min_stock_level DECIMAL(15, 3) DEFAULT 0,
	max_stock_level DECIMAL(15, 3),
	reorder_point DECIMAL(15, 3) DEFAULT 0,
	lead_time_days INTEGER DEFAULT 0,
	status sku_status NOT NULL DEFAULT 'ACTIVE',
	is_purchasable BOOLEAN DEFAULT true,
	is_sellable BOOLEAN DEFAULT true,
	tax_rate DECIMAL(5, 2) DEFAULT 0,
	created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
-- Create indexes for better query performance
CREATE INDEX idx_skus_sku_code ON skus(sku_code);
CREATE INDEX idx_skus_name ON skus(name);
CREATE INDEX idx_skus_category_id ON skus(category_id);
CREATE INDEX idx_skus_primary_vendor_id ON skus(primary_vendor_id);
CREATE INDEX idx_skus_backup_vendor_id ON skus(backup_vendor_id);
CREATE INDEX idx_skus_status ON skus(status);
CREATE INDEX idx_sku_categories_parent_id ON sku_categories(parent_id);
-- Create table for SKU prices history
CREATE TABLE IF NOT EXISTS sku_price_history (
	id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	sku_id UUID REFERENCES skus(id) ON DELETE CASCADE,
	price DECIMAL(15, 2) NOT NULL,
	effective_from TIMESTAMP WITH TIME ZONE NOT NULL,
	effective_to TIMESTAMP WITH TIME ZONE,
	created_by UUID REFERENCES users(id),
	created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
-- Create table for SKU-Store relationships (additional attributes specific to SKU in a store)
CREATE TABLE IF NOT EXISTS sku_store_settings (
	sku_id UUID REFERENCES skus(id) ON DELETE CASCADE,
	store_id UUID REFERENCES stores(id) ON DELETE CASCADE,
	min_stock_level DECIMAL(15, 3),
	max_stock_level DECIMAL(15, 3),
	reorder_point DECIMAL(15, 3),
	bin_location VARCHAR(50),
	created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
	PRIMARY KEY (sku_id, store_id)
);
-- Create table for SKU barcodes
CREATE TABLE IF NOT EXISTS sku_barcodes (
	id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	sku_id UUID REFERENCES skus(id) ON DELETE CASCADE,
	barcode_type VARCHAR(50) NOT NULL,
	barcode_value VARCHAR(100) NOT NULL,
	is_primary BOOLEAN DEFAULT false,
	created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
	UNIQUE(barcode_type, barcode_value)
);
-- Create indexes for the new tables
CREATE INDEX idx_sku_price_history_sku_id ON sku_price_history(sku_id);
CREATE INDEX idx_sku_price_history_effective_from ON sku_price_history(effective_from);
CREATE INDEX idx_sku_store_settings_store_id ON sku_store_settings(store_id);
CREATE INDEX idx_sku_barcodes_sku_id ON sku_barcodes(sku_id);
CREATE INDEX idx_sku_barcodes_value ON sku_barcodes(barcode_value);
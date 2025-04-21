CREATE TYPE store_type AS ENUM ('RAW', 'FINISHED', 'GENERAL');
CREATE TYPE store_status AS ENUM ('ACTIVE', 'INACTIVE');
CREATE TABLE IF NOT EXISTS stores (
	id UUID PRIMARY KEY,
	name VARCHAR(255) NOT NULL,
	address TEXT,
	type store_type NOT NULL,
	manager_id UUID NOT NULL REFERENCES users(id),
	contact VARCHAR(255),
	status store_status NOT NULL DEFAULT 'ACTIVE',
	created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
	UNIQUE(name)
);
CREATE TABLE IF NOT EXISTS stocks (
	id UUID PRIMARY KEY,
	sku_id UUID NOT NULL,
	store_id UUID NOT NULL REFERENCES stores(id),
	quantity DECIMAL(15, 3) NOT NULL DEFAULT 0,
	bin_location VARCHAR(50),
	shelf_number VARCHAR(50),
	zone_code VARCHAR(50),
	batch_number VARCHAR(100),
	lot_number VARCHAR(100),
	manufacture_date TIMESTAMP WITH TIME ZONE,
	expiry_date TIMESTAMP WITH TIME ZONE,
	created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
	UNIQUE(sku_id, store_id)
);
CREATE TABLE IF NOT EXISTS stock_entries (
	id UUID PRIMARY KEY,
	store_id UUID NOT NULL REFERENCES stores(id),
	sku_id UUID NOT NULL,
	type VARCHAR(10) NOT NULL CHECK (type IN ('IN', 'OUT')),
	quantity DECIMAL(15, 3) NOT NULL,
	batch_number VARCHAR(100),
	lot_number VARCHAR(100),
	reference VARCHAR(100),
	note TEXT,
	created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
	created_by UUID NOT NULL REFERENCES users(id)
);
CREATE TABLE IF NOT EXISTS stock_history (
	id UUID PRIMARY KEY,
	stock_id UUID NOT NULL REFERENCES stocks(id),
	type VARCHAR(10) NOT NULL CHECK (type IN ('IN', 'OUT', 'ADJUST')),
	quantity DECIMAL(15, 3) NOT NULL,
	previous_qty DECIMAL(15, 3) NOT NULL,
	new_qty DECIMAL(15, 3) NOT NULL,
	reference UUID REFERENCES stock_entries(id),
	note TEXT,
	created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
	created_by UUID NOT NULL REFERENCES users(id)
);
CREATE INDEX idx_stocks_store_id ON stocks(store_id);
CREATE INDEX idx_stocks_sku_id ON stocks(sku_id);
CREATE INDEX idx_stock_entries_store_id ON stock_entries(store_id);
CREATE INDEX idx_stock_entries_sku_id ON stock_entries(sku_id);
CREATE INDEX idx_stock_history_stock_id ON stock_history(stock_id);
CREATE INDEX idx_stock_history_reference ON stock_history(reference);
CREATE TYPE warehouse_type AS ENUM ('RAW', 'FINISHED', 'GENERAL');
CREATE TYPE warehouse_status AS ENUM ('ACTIVE', 'INACTIVE');
CREATE TABLE IF NOT EXISTS warehouses (
	id UUID PRIMARY KEY,
	name VARCHAR(255) NOT NULL,
	address TEXT,
	type warehouse_type NOT NULL,
	manager_id UUID NOT NULL REFERENCES users(id),
	contact VARCHAR(255),
	status warehouse_status NOT NULL DEFAULT 'ACTIVE',
	created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
	UNIQUE(name)
);
CREATE TABLE IF NOT EXISTS inventories (
	id UUID PRIMARY KEY,
	product_id UUID NOT NULL,
	warehouse_id UUID NOT NULL REFERENCES warehouses(id),
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
	UNIQUE(product_id, warehouse_id)
);
CREATE TABLE IF NOT EXISTS stock_entries (
	id UUID PRIMARY KEY,
	warehouse_id UUID NOT NULL REFERENCES warehouses(id),
	product_id UUID NOT NULL,
	type VARCHAR(10) NOT NULL CHECK (type IN ('IN', 'OUT')),
	quantity DECIMAL(15, 3) NOT NULL,
	batch_number VARCHAR(100),
	lot_number VARCHAR(100),
	reference VARCHAR(100),
	note TEXT,
	created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
	created_by UUID NOT NULL REFERENCES users(id)
);
CREATE TABLE IF NOT EXISTS inventory_history (
	id UUID PRIMARY KEY,
	inventory_id UUID NOT NULL REFERENCES inventories(id),
	type VARCHAR(10) NOT NULL CHECK (type IN ('IN', 'OUT', 'ADJUST')),
	quantity DECIMAL(15, 3) NOT NULL,
	previous_qty DECIMAL(15, 3) NOT NULL,
	new_qty DECIMAL(15, 3) NOT NULL,
	reference UUID REFERENCES stock_entries(id),
	note TEXT,
	created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
	created_by UUID NOT NULL REFERENCES users(id)
);
CREATE INDEX idx_inventories_warehouse_id ON inventories(warehouse_id);
CREATE INDEX idx_inventories_product_id ON inventories(product_id);
CREATE INDEX idx_stock_entries_warehouse_id ON stock_entries(warehouse_id);
CREATE INDEX idx_stock_entries_product_id ON stock_entries(product_id);
CREATE INDEX idx_inventory_history_inventory_id ON inventory_history(inventory_id);
CREATE INDEX idx_inventory_history_reference ON inventory_history(reference);
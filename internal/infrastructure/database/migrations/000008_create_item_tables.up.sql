CREATE TYPE item_status AS ENUM ('ACTIVE', 'INACTIVE', 'ARCHIVED');
CREATE TABLE IF NOT EXISTS item_categories (
	id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	name VARCHAR(255) NOT NULL,
	description TEXT,
	parent_id UUID REFERENCES item_categories(id),
	created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
	UNIQUE(name)
);
CREATE TABLE IF NOT EXISTS items (
	id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	sku VARCHAR(100) UNIQUE NOT NULL,
	name VARCHAR(255) NOT NULL,
	description TEXT,
	unit_of_measure VARCHAR(50) NOT NULL,
	price DECIMAL(15, 2) DEFAULT 0,
	category VARCHAR(100),
	technical_specs JSONB,
	manufacturer_id INTEGER REFERENCES suppliers(id),
	supplier_id INTEGER REFERENCES suppliers(id),
	image_url VARCHAR(255),
	status item_status NOT NULL DEFAULT 'ACTIVE',
	created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
-- Create indexes for better query performance
CREATE INDEX idx_items_sku ON items(sku);
CREATE INDEX idx_items_name ON items(name);
CREATE INDEX idx_items_category ON items(category);
CREATE INDEX idx_items_manufacturer_id ON items(manufacturer_id);
CREATE INDEX idx_items_supplier_id ON items(supplier_id);
CREATE INDEX idx_items_status ON items(status);
CREATE INDEX idx_item_categories_parent_id ON item_categories(parent_id);
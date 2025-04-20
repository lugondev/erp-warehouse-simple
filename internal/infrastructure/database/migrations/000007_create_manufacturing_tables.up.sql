CREATE TABLE manufacturing_facilities (
	id SERIAL PRIMARY KEY,
	name VARCHAR(255) NOT NULL,
	address TEXT NOT NULL,
	type VARCHAR(100) NOT NULL,
	capacity INTEGER NOT NULL,
	manager VARCHAR(255),
	created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
CREATE TABLE production_orders (
	id SERIAL PRIMARY KEY,
	product_id INTEGER NOT NULL,
	quantity INTEGER NOT NULL,
	start_date TIMESTAMP WITH TIME ZONE,
	deadline TIMESTAMP WITH TIME ZONE NOT NULL,
	status VARCHAR(50) NOT NULL DEFAULT 'pending',
	facility_id INTEGER NOT NULL,
	completed_qty INTEGER DEFAULT 0,
	defect_qty INTEGER DEFAULT 0,
	notes TEXT,
	created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
	FOREIGN KEY (facility_id) REFERENCES manufacturing_facilities(id),
	FOREIGN KEY (product_id) REFERENCES inventory_items(id)
);
CREATE TABLE bill_of_materials (
	id SERIAL PRIMARY KEY,
	product_id INTEGER NOT NULL,
	name VARCHAR(255) NOT NULL,
	version VARCHAR(50) NOT NULL,
	created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
	FOREIGN KEY (product_id) REFERENCES inventory_items(id)
);
CREATE TABLE bom_items (
	id SERIAL PRIMARY KEY,
	bom_id INTEGER NOT NULL,
	material_id INTEGER NOT NULL,
	quantity_needed DECIMAL(10, 2) NOT NULL,
	unit_of_measure VARCHAR(50) NOT NULL,
	notes TEXT,
	created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
	FOREIGN KEY (bom_id) REFERENCES bill_of_materials(id),
	FOREIGN KEY (material_id) REFERENCES inventory_items(id)
);
CREATE TABLE mrp_calculations (
	id SERIAL PRIMARY KEY,
	production_id INTEGER NOT NULL,
	material_id INTEGER NOT NULL,
	required_qty DECIMAL(10, 2) NOT NULL,
	available_qty DECIMAL(10, 2),
	shortage_qty DECIMAL(10, 2),
	unit_of_measure VARCHAR(50),
	calculated_at TIMESTAMP WITH TIME ZONE NOT NULL,
	created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
	FOREIGN KEY (production_id) REFERENCES production_orders(id),
	FOREIGN KEY (material_id) REFERENCES inventory_items(id)
);
-- Add indexes for better query performance
CREATE INDEX idx_production_orders_facility ON production_orders(facility_id);
CREATE INDEX idx_production_orders_product ON production_orders(product_id);
CREATE INDEX idx_bom_items_material ON bom_items(material_id);
CREATE INDEX idx_mrp_calculations_production ON mrp_calculations(production_id);
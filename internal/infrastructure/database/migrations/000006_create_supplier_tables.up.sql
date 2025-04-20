CREATE TABLE IF NOT EXISTS suppliers (
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
CREATE TABLE IF NOT EXISTS products (
	id SERIAL PRIMARY KEY,
	code VARCHAR(50) UNIQUE NOT NULL,
	name VARCHAR(255) NOT NULL,
	description TEXT,
	unit_price DECIMAL(15, 2),
	currency VARCHAR(10),
	created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
CREATE TABLE IF NOT EXISTS supplier_products (
	supplier_id INTEGER REFERENCES suppliers(id) ON DELETE CASCADE,
	product_id INTEGER REFERENCES products(id) ON DELETE CASCADE,
	PRIMARY KEY (supplier_id, product_id)
);
CREATE TABLE IF NOT EXISTS contracts (
	id SERIAL PRIMARY KEY,
	supplier_id INTEGER REFERENCES suppliers(id) ON DELETE CASCADE,
	contract_no VARCHAR(100) UNIQUE NOT NULL,
	start_date TIMESTAMP WITH TIME ZONE,
	end_date TIMESTAMP WITH TIME ZONE,
	terms TEXT,
	status VARCHAR(50),
	created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
CREATE TABLE IF NOT EXISTS supplier_ratings (
	id SERIAL PRIMARY KEY,
	supplier_id INTEGER REFERENCES suppliers(id) ON DELETE CASCADE,
	score DECIMAL(3, 2) NOT NULL,
	category VARCHAR(100),
	comment TEXT,
	rated_by INTEGER REFERENCES users(id),
	created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX idx_suppliers_code ON suppliers(code);
CREATE INDEX idx_suppliers_name ON suppliers(name);
CREATE INDEX idx_products_code ON products(code);
CREATE INDEX idx_contracts_supplier ON contracts(supplier_id);
CREATE INDEX idx_contracts_dates ON contracts(start_date, end_date);
CREATE INDEX idx_supplier_ratings ON supplier_ratings(supplier_id);
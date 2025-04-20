-- Create customers table
CREATE TABLE IF NOT EXISTS customers (
	id SERIAL PRIMARY KEY,
	code VARCHAR(50) NOT NULL UNIQUE,
	name VARCHAR(255) NOT NULL,
	type VARCHAR(50) NOT NULL DEFAULT 'INDIVIDUAL',
	email VARCHAR(255) UNIQUE,
	phone_number VARCHAR(50),
	tax_id VARCHAR(50),
	contacts JSONB,
	credit_limit DECIMAL(15, 2) DEFAULT 0,
	current_debt DECIMAL(15, 2) DEFAULT 0,
	loyalty_tier VARCHAR(50) NOT NULL DEFAULT 'STANDARD',
	loyalty_points INTEGER DEFAULT 0,
	notes TEXT,
	created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
	updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);
-- Create customer addresses table
CREATE TABLE IF NOT EXISTS customer_addresses (
	id SERIAL PRIMARY KEY,
	customer_id INTEGER NOT NULL REFERENCES customers(id) ON DELETE CASCADE,
	type VARCHAR(50) NOT NULL DEFAULT 'BOTH',
	street VARCHAR(255) NOT NULL,
	city VARCHAR(100) NOT NULL,
	state VARCHAR(100),
	postal_code VARCHAR(20),
	country VARCHAR(100) NOT NULL,
	is_default BOOLEAN DEFAULT FALSE,
	created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
	updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);
-- Create index on customer_id for faster lookups
CREATE INDEX IF NOT EXISTS idx_customer_addresses_customer_id ON customer_addresses(customer_id);
-- Create index on customer name for faster searches
CREATE INDEX IF NOT EXISTS idx_customers_name ON customers(name);
-- Create index on customer email for faster searches
CREATE INDEX IF NOT EXISTS idx_customers_email ON customers(email);
-- Create index on customer type for filtering
CREATE INDEX IF NOT EXISTS idx_customers_type ON customers(type);
-- Create index on customer loyalty tier for filtering
CREATE INDEX IF NOT EXISTS idx_customers_loyalty_tier ON customers(loyalty_tier);
-- Add foreign key to sales_orders table if it exists
DO $$ BEGIN IF EXISTS (
	SELECT
	FROM information_schema.tables
	WHERE table_name = 'sales_orders'
) THEN -- Check if the column already exists
IF NOT EXISTS (
	SELECT
	FROM information_schema.columns
	WHERE table_name = 'sales_orders'
		AND column_name = 'customer_id'
) THEN
ALTER TABLE sales_orders
ADD COLUMN customer_id INTEGER REFERENCES customers(id);
END IF;
END IF;
END $$;
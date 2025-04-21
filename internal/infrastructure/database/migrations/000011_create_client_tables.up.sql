-- Create clients table
CREATE TABLE IF NOT EXISTS clients (
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
-- Create client addresses table
CREATE TABLE IF NOT EXISTS client_addresses (
	id SERIAL PRIMARY KEY,
	client_id INTEGER NOT NULL REFERENCES clients(id) ON DELETE CASCADE,
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
-- Create index on client_id for faster lookups
CREATE INDEX IF NOT EXISTS idx_client_addresses_client_id ON client_addresses(client_id);
-- Create index on client name for faster searches
CREATE INDEX IF NOT EXISTS idx_clients_name ON clients(name);
-- Create index on client email for faster searches
CREATE INDEX IF NOT EXISTS idx_clients_email ON clients(email);
-- Create index on client type for filtering
CREATE INDEX IF NOT EXISTS idx_clients_type ON clients(type);
-- Create index on client loyalty tier for filtering
CREATE INDEX IF NOT EXISTS idx_clients_loyalty_tier ON clients(loyalty_tier);
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
		AND column_name = 'client_id'
) THEN
ALTER TABLE sales_orders
ADD COLUMN client_id INTEGER REFERENCES clients(id);
END IF;
END IF;
END $$;
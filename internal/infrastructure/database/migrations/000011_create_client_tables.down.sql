-- Remove foreign key from sales_orders table if it exists
DO $$ BEGIN IF EXISTS (
	SELECT
	FROM information_schema.tables
	WHERE table_name = 'sales_orders'
) THEN -- Check if the column exists
IF EXISTS (
	SELECT
	FROM information_schema.columns
	WHERE table_name = 'sales_orders'
		AND column_name = 'client_id'
) THEN
ALTER TABLE sales_orders DROP COLUMN client_id;
END IF;
END IF;
END $$;
-- Drop client addresses table
DROP TABLE IF EXISTS client_addresses;
-- Drop clients table
DROP TABLE IF EXISTS clients;
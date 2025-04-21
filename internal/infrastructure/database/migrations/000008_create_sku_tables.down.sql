-- Drop indexes
DROP INDEX IF EXISTS idx_sku_barcodes_value;
DROP INDEX IF EXISTS idx_sku_barcodes_sku_id;
DROP INDEX IF EXISTS idx_sku_store_settings_store_id;
DROP INDEX IF EXISTS idx_sku_price_history_effective_from;
DROP INDEX IF EXISTS idx_sku_price_history_sku_id;
DROP INDEX IF EXISTS idx_sku_categories_parent_id;
DROP INDEX IF EXISTS idx_skus_status;
DROP INDEX IF EXISTS idx_skus_backup_vendor_id;
DROP INDEX IF EXISTS idx_skus_primary_vendor_id;
DROP INDEX IF EXISTS idx_skus_category_id;
DROP INDEX IF EXISTS idx_skus_name;
DROP INDEX IF EXISTS idx_skus_sku_code;
-- Drop tables
DROP TABLE IF EXISTS sku_barcodes;
DROP TABLE IF EXISTS sku_store_settings;
DROP TABLE IF EXISTS sku_price_history;
DROP TABLE IF EXISTS skus;
DROP TABLE IF EXISTS sku_categories;
-- Drop types
DROP TYPE IF EXISTS sku_status;
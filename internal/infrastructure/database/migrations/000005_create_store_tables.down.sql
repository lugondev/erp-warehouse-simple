DROP INDEX IF EXISTS idx_stock_history_reference;
DROP INDEX IF EXISTS idx_stock_history_stock_id;
DROP INDEX IF EXISTS idx_stock_entries_sku_id;
DROP INDEX IF EXISTS idx_stock_entries_store_id;
DROP INDEX IF EXISTS idx_stocks_sku_id;
DROP INDEX IF EXISTS idx_stocks_store_id;
DROP TABLE IF EXISTS stock_history;
DROP TABLE IF EXISTS stock_entries;
DROP TABLE IF EXISTS stocks;
DROP TABLE IF EXISTS stores;
DROP TYPE IF EXISTS store_status;
DROP TYPE IF EXISTS store_type;
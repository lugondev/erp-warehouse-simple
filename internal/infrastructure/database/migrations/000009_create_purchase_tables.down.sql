-- Drop indexes
DROP INDEX IF EXISTS idx_purchase_payments_purchase_order;
DROP INDEX IF EXISTS idx_purchase_receipts_warehouse;
DROP INDEX IF EXISTS idx_purchase_receipts_purchase_order;
DROP INDEX IF EXISTS idx_purchase_orders_payment_status;
DROP INDEX IF EXISTS idx_purchase_orders_status;
DROP INDEX IF EXISTS idx_purchase_orders_supplier;
DROP INDEX IF EXISTS idx_purchase_requests_purchase_order;
DROP INDEX IF EXISTS idx_purchase_requests_status;
DROP INDEX IF EXISTS idx_purchase_requests_requester;
-- Drop foreign key constraint
ALTER TABLE IF EXISTS purchase_requests DROP CONSTRAINT IF EXISTS fk_purchase_requests_purchase_order;
-- Drop tables
DROP TABLE IF EXISTS purchase_payments;
DROP TABLE IF EXISTS purchase_receipts;
DROP TABLE IF EXISTS purchase_requests;
DROP TABLE IF EXISTS purchase_orders;
-- Drop indexes
DROP INDEX IF EXISTS idx_invoices_due_date;
DROP INDEX IF EXISTS idx_invoices_issue_date;
DROP INDEX IF EXISTS idx_invoices_status;
DROP INDEX IF EXISTS idx_invoices_sales_order_id;
DROP INDEX IF EXISTS idx_delivery_orders_delivery_date;
DROP INDEX IF EXISTS idx_delivery_orders_warehouse_id;
DROP INDEX IF EXISTS idx_delivery_orders_status;
DROP INDEX IF EXISTS idx_delivery_orders_sales_order_id;
DROP INDEX IF EXISTS idx_sales_orders_order_date;
DROP INDEX IF EXISTS idx_sales_orders_payment_status;
DROP INDEX IF EXISTS idx_sales_orders_status;
DROP INDEX IF EXISTS idx_sales_orders_customer_id;
-- Drop tables
DROP TABLE IF EXISTS invoices;
DROP TABLE IF EXISTS delivery_orders;
DROP TABLE IF EXISTS sales_orders;
-- Remove sequence values
DELETE FROM sequences
WHERE id IN ('sales_order', 'delivery_order', 'invoice');
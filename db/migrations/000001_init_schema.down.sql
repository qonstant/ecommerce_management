-- Drop foreign key constraints
ALTER TABLE "payments" DROP CONSTRAINT IF EXISTS payments_order_id_fkey;
ALTER TABLE "payments" DROP CONSTRAINT IF EXISTS payments_user_id_fkey;
ALTER TABLE "order_items" DROP CONSTRAINT IF EXISTS order_items_product_id_fkey;
ALTER TABLE "order_items" DROP CONSTRAINT IF EXISTS order_items_order_id_fkey;
ALTER TABLE "orders" DROP CONSTRAINT IF EXISTS orders_user_id_fkey;

-- Drop tables
DROP TABLE IF EXISTS "payments";
DROP TABLE IF EXISTS "order_items";
DROP TABLE IF EXISTS "orders";
DROP TABLE IF EXISTS "products";
DROP TABLE IF EXISTS "users";

-- Drop types
DROP TYPE IF EXISTS "payment_status";
DROP TYPE IF EXISTS "order_status";

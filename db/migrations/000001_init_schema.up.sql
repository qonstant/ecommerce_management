CREATE TYPE "order_status" AS ENUM (
  'new',
  'processing',
  'completed'
);

CREATE TYPE "payment_status" AS ENUM (
  'successful',
  'unsuccessful'
);

CREATE TABLE "users" (
  "id" BIGSERIAL PRIMARY KEY,
  "full_name" varchar(255) NOT NULL,
  "email" varchar(320) UNIQUE NOT NULL,
  "address" text NOT NULL,
  "registration_date" timestamp NOT NULL DEFAULT NOW(),
  "role" varchar(50) NOT NULL
);

CREATE TABLE "products" (
  "id" BIGSERIAL PRIMARY KEY,
  "name" varchar(255) NOT NULL,
  "description" text NOT NULL,
  "price" numeric(10,2) NOT NULL,
  "category" varchar(255) NOT NULL,
  "stock_quantity" int NOT NULL,
  "addition_date" timestamp NOT NULL DEFAULT NOW()
);

CREATE TABLE "orders" (
  "id" BIGSERIAL PRIMARY KEY,
  "user_id" BIGINT NOT NULL,
  "total_amount" numeric(10,2) NOT NULL,
  "order_date" timestamp NOT NULL DEFAULT NOW(),
  "status" order_status NOT NULL
);

CREATE TABLE "order_items" (
  "id" BIGSERIAL PRIMARY KEY,
  "order_id" BIGINT NOT NULL,
  "product_id" BIGINT NOT NULL,
  "quantity" int NOT NULL,
  "price" numeric(10,2) NOT NULL
);

CREATE TABLE "payments" (
  "id" BIGSERIAL PRIMARY KEY,
  "user_id" BIGINT NOT NULL,
  "order_id" BIGINT NOT NULL,
  "amount" numeric(10,2) NOT NULL,
  "payment_date" timestamp NOT NULL DEFAULT NOW(),
  "status" payment_status NOT NULL
);

ALTER TABLE "orders" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");

ALTER TABLE "order_items" ADD FOREIGN KEY ("order_id") REFERENCES "orders" ("id");

ALTER TABLE "order_items" ADD FOREIGN KEY ("product_id") REFERENCES "products" ("id");

ALTER TABLE "payments" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");

ALTER TABLE "payments" ADD FOREIGN KEY ("order_id") REFERENCES "orders" ("id");

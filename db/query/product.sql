-- name: GetProduct :one
SELECT * FROM products WHERE id = $1 LIMIT 1;

-- name: ListProducts :many
SELECT * FROM products ORDER BY addition_date ASC;

-- name: CreateProduct :one
INSERT INTO products (name, description, price, category, stock_quantity, addition_date) 
VALUES ($1, $2, $3, $4, $5, NOW()) 
RETURNING *;

-- name: UpdateProduct :one
UPDATE products SET 
    name = $2,
    description = $3,
    price = $4,
    category = $5,
    stock_quantity = $6
WHERE id = $1 
RETURNING *;

-- name: DeleteProduct :exec
DELETE FROM products WHERE id = $1;

-- name: SearchProductsByName :many
SELECT * FROM products WHERE name ILIKE '%' || $1 || '%' ORDER BY addition_date ASC;

-- name: SearchProductsByCategory :many
SELECT * FROM products WHERE category = $1 ORDER BY addition_date ASC;

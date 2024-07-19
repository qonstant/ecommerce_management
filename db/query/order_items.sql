-- name: GetOrderItem :one
SELECT * FROM order_items WHERE id = $1 LIMIT 1;

-- name: ListOrderItems :many
SELECT * FROM order_items ORDER BY id ASC;

-- name: CreateOrderItem :one
INSERT INTO order_items (order_id, product_id, quantity, price) 
VALUES ($1, $2, $3, $4) 
RETURNING *;

-- name: UpdateOrderItem :one
UPDATE order_items SET 
    order_id = $2,
    product_id = $3,
    quantity = $4,
    price = $5
WHERE id = $1 
RETURNING *;

-- name: DeleteOrderItem :exec
DELETE FROM order_items WHERE id = $1;

-- name: ListOrderItemsByOrder :many
SELECT * FROM order_items WHERE order_id = $1 ORDER BY id ASC;

-- name: ListOrderItemsByProduct :many
SELECT * FROM order_items WHERE product_id = $1 ORDER BY id ASC;

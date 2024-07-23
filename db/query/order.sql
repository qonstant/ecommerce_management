-- name: GetOrder :one
SELECT * FROM orders WHERE id = $1 LIMIT 1;

-- name: ListOrders :many
SELECT * FROM orders ORDER BY order_date ASC;

-- name: CreateOrder :one
INSERT INTO orders (user_id, total_amount, order_date) 
VALUES ($1, $2, NOW()) 
RETURNING id, user_id, total_amount, order_date, status;

-- name: UpdateOrder :one
UPDATE orders SET 
    user_id = $2,
    total_amount = $3,
    status = $4
WHERE id = $1 
RETURNING *;

-- name: DeleteOrder :exec
DELETE FROM orders WHERE id = $1;

-- name: SearchOrdersByUser :many
SELECT * FROM orders WHERE user_id = $1 ORDER BY order_date ASC;

-- name: SearchOrdersByStatus :many
SELECT * FROM orders WHERE status = $1 ORDER BY order_date ASC;

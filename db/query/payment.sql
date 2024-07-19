-- name: GetPayment :one
SELECT * FROM payments WHERE id = $1 LIMIT 1;

-- name: ListPayments :many
SELECT * FROM payments ORDER BY payment_date ASC;

-- name: CreatePayment :one
INSERT INTO payments (user_id, order_id, amount, payment_date, status) 
VALUES ($1, $2, $3, NOW(), $4) 
RETURNING *;

-- name: UpdatePayment :one
UPDATE payments SET 
    user_id = $2,
    order_id = $3,
    amount = $4,
    status = $5
WHERE id = $1 
RETURNING *;

-- name: DeletePayment :exec
DELETE FROM payments WHERE id = $1;

-- name: SearchPaymentsByUser :many
SELECT * FROM payments WHERE user_id = $1 ORDER BY payment_date ASC;

-- name: SearchPaymentsByOrder :many
SELECT * FROM payments WHERE order_id = $1 ORDER BY payment_date ASC;

-- name: SearchPaymentsByStatus :many
SELECT * FROM payments WHERE status = $1 ORDER BY payment_date ASC;

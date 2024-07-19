-- name: GetUser :one
SELECT * FROM users WHERE id = $1 LIMIT 1;

-- name: ListUsers :many
SELECT * FROM users ORDER BY registration_date ASC;

-- name: CreateUser :one
INSERT INTO users (full_name, email, address, registration_date, role) 
VALUES ($1, $2, $3, NOW(), $4) 
RETURNING *;

-- name: UpdateUser :one
UPDATE users SET 
    full_name = $2,
    email = $3,
    address = $4,
    role = $5
WHERE id = $1 
RETURNING *;

-- name: DeleteUser :exec
DELETE FROM users WHERE id = $1;

-- name: SearchUsersByName :many
SELECT * FROM users WHERE full_name ILIKE '%' || $1 || '%' ORDER BY registration_date ASC;

-- name: SearchUsersByEmail :many
SELECT * FROM users WHERE email = $1 ORDER BY registration_date ASC;

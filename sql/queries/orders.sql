-- name: CreateOrder :one
INSERT INTO orders (id, user_id, total_amount, delivery_address)
VALUES (
    gen_random_uuid(),
    $1,
    $2,
    $3
) RETURNING *;

-- name: GetOrderHistory :many
SELECT * FROM orders 
WHERE user_id = $1;

-- name: UpdateOrderStatus :one
UPDATE orders 
SET 
    status = $1,
    updated_at = NOW()
WHERE id = $2
RETURNING *;
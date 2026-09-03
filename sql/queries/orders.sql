-- name: CreateOrder :one
INSERT INTO orders (id, user_id, total_amount, delivery_address)
VALUES (
    gen_random_uuid(),
    $1,
    $2,
    $3
) RETURNING *;
-- name: CreateOrderItem :one
INSERT INTO order_items (id, order_id, product_id, quantity, price_at_purchase)
VALUES (
    gen_random_uuid(),
    $1, 
    $2,
    $3, 
    $4
) RETURNING *;
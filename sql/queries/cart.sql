-- name: AddItemToCart :one
INSERT INTO cart_items (id, user_id, product_id)
VALUES (
    gen_random_uuid(),
    $1,
    $2
)
ON CONFLICT (user_id, product_id)
DO UPDATE SET 
    quantity = cart_items.quantity + 1,
    updated_at = NOW()
RETURNING *;
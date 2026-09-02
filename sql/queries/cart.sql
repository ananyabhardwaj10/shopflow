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

-- name: GetAllCartItems :many
SELECT cart_items.id, cart_items.quantity, products.name, products.price, products.id as product_id
FROM cart_items
INNER JOIN products
ON cart_items.product_id = products.id 
WHERE user_id = $1;

-- name: UpdateCartItemQuantity :one
UPDATE cart_items 
SET 
    quantity = $1,
    updated_at = NOW()
WHERE id = $2 AND user_id = $3
RETURNING *;
-- name: CreateOrderItem :one
INSERT INTO order_items (id, order_id, product_id, quantity, price_at_purchase)
VALUES (
    gen_random_uuid(),
    $1, 
    $2,
    $3, 
    $4
) RETURNING *;

-- name: UpdateOrderItemStatus :one 
UPDATE order_items
SET 
    status = $1, 
    updated_at = NOW()
WHERE order_items.id = $2 AND order_items.order_id = $3
AND EXISTS (
    SELECT 1 FROM products p
    JOIN sellers s ON p.seller_id = s.id 
    WHERE p.id = order_items.product_id
    AND s.user_id = $4
)
RETURNING *;

-- name: CheckAllOrderItemsConfirmed :one
SELECT COUNT(*) FROM order_items
WHERE order_id = $1
AND status != 'confirmed';
-- name: CreateProduct :one
INSERT INTO products (id, name, description, price, stock_quantity, seller_id, category_id) 
VALUES (
    gen_random_uuid(),
    $1,
    $2,
    $3,
    $4,
    $5,
    $6
) RETURNING *;

-- name: GetAllProductsBySeller :many
SELECT * FROM products
WHERE seller_id = $1
LIMIT $2 OFFSET $3;

-- name: UpdateProductDetails :one
UPDATE products
SET 
    name = COALESCE(sqlc.narg('name'), name),
    description = COALESCE(sqlc.narg('description'), description),
    price = COALESCE(sqlc.narg('price'), price),
    stock_quantity = COALESCE(sqlc.narg('stock_quantity'), stock_quantity),
    updated_at = NOW()
WHERE id = $1 AND seller_id = $2
RETURNING *;

-- name: DeleteProduct :exec
DELETE FROM products
WHERE id = $1 AND seller_id = $2;

-- name: GetAllProductsByCategory :many
SELECT * FROM products 
WHERE category_id = $1
LIMIT $2 OFFSET $3;

-- name: GetProductByID :one
SELECT * FROM products 
WHERE id = $1;
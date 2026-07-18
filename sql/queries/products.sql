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
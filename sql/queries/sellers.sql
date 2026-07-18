-- name: CreateSeller :one
INSERT INTO sellers (id, business_name, user_id) 
VALUES (
    gen_random_uuid(),
    $1,
    $2
) RETURNING *;

-- name: GetSellerByUserID :one
SELECT * FROM sellers 
WHERE user_id = $1;
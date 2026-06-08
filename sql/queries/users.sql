-- name: CreateUser :one
INSERT INTO users (id, first_name, last_name, contact_number, email, hashed_password)
VALUES(
    gen_random_uuid(),
    $1,
    $2,
    $3,
    $4,
    $5
) RETURNING *;

-- name: GetUserByEmail :one
SELECT * FROM users 
WHERE email = $1;
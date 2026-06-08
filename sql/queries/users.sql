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

-- name: UpdateUserDetails :one
UPDATE users 
SET 
    first_name = COALESCE(sqlc.narg('first_name'), first_name),
    last_name = COALESCE(sqlc.narg('last_name'), last_name),
    contact_number = COALESCE(sqlc.narg('contact_number'), contact_number),
    address = COALESCE(sqlc.narg('address'), address),
    email = COALESCE(sqlc.narg('email'), email),
    updated_at = NOW()
WHERE id = $1
RETURNING *;
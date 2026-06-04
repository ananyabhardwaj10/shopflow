-- +goose Up
CREATE TABLE users (
    id UUID PRIMARY KEY,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    first_name TEXT NOT NULL,
    last_name TEXT NOT NULL,
    contact_number VARCHAR(15) NOT NULL,
    address TEXT DEFAULT NULL,
    email TEXT UNIQUE NOT NULL,
    hashed_password TEXT NOT NULL,
    role TEXT NOT NULL DEFAULT 'customer'
);

-- +goose Down
DROP TABLE users;
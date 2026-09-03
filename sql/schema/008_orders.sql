-- +goose Up
CREATE TABLE orders (
    id UUID PRIMARY KEY,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    user_id UUID NOT NULL, 
    total_amount NUMERIC(10, 2) NOT NULL,
    status TEXT NOT NULL DEFAULT 'pending', 
    delivery_address TEXT NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- +goose Down
DROP TABLE orders;
-- +goose Up
CREATE TABLE order_items (
    id UUID PRIMARY KEY,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    order_id UUID NOT NULL, 
    product_id UUID NOT NULL,
    quantity INTEGER NOT NULL,
    price_at_purchase NUMERIC(10, 2) NOT NULL, 
    FOREIGN KEY (order_id) REFERENCES orders(id) ON DELETE CASCADE
);

-- +goose Down
DROP TABLE order_items;

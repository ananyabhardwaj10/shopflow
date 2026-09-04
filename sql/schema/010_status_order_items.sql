-- +goose Up
ALTER TABLE order_items 
ADD status TEXT NOT NULL DEFAULT 'pending';

-- +goose Down
ALTER TABLE order_items 
DROP COLUMN status;
-- +goose Up
INSERT INTO categories (id, name) 
VALUES 
    (gen_random_uuid(), 'Electronics'),
    (gen_random_uuid(), 'Clothing'),
    (gen_random_uuid(), 'Books'),
    (gen_random_uuid(), 'Footwear'),
    (gen_random_uuid(), 'Sports'),
    (gen_random_uuid(), 'Home'),
    (gen_random_uuid(), 'Beauty'),
    (gen_random_uuid(), 'Toys');

-- +goose Down
DELETE FROM categories;
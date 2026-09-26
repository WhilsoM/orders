-- +goose Up

CREATE TABLE inventory (
    product_id BIGINT PRIMARY KEY,
    quantity INTEGER NOT NULL,
    unit_price INTEGER NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);

-- +goose Down

DROP TABLE inventory;

CREATE TABLE IF NOT EXISTS order_items (
    id BIGSERIAL PRIMARY KEY,
    order_id BIGINT NOT NULL REFERENCES orders(id),
    book_id BIGINT NOT NULL REFERENCES books(id),
    quantity INT NOT NULL,
    price DECIMAL(10,2) NOT NULL,
    UNIQUE (order_id , book_id)
)
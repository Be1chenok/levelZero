CREATE TABLE IF NOT EXISTS orders(
    uid VARCHAR(64) PRIMARY KEY,
    track_number VARCHAR(64) NOT NULL,
    entry VARCHAR(64) NOT NULL,
    locale VARCHAR(6) NOT NULL,
    internal_signature VARCHAR(64) NOT NULL,
    customer_id VARCHAR(64) NOT NULL,
    delivery_service VARCHAR(64) NOT NULL,
    shardkey VARCHAR(64) NOT NULL,
    sm_id INT NOT NULL,
    date_created TIMESTAMP NOT NULL,
    oof_shard VARCHAR(64) NOT NULL
);

CREATE TABLE IF NOT EXISTS deliveries(
    id BIGSERIAL PRIMARY KEY,
    order_uid VARCHAR(64) UNIQUE,
    name VARCHAR(64) NOT NULL,
    phone VARCHAR(16) NOT NULL,
    zip VARCHAR(255) NOT NULL,
    city VARCHAR(255) NOT NULL,
    address VARCHAR(255) NOT NULL,
    region VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL,
    FOREIGN KEY (order_uid) REFERENCES orders (uid) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS payments(
    id BIGSERIAL PRIMARY KEY,
    order_uid VARCHAR(64) UNIQUE,
    transaction VARCHAR(64) NOT NULL,
    request_id VARCHAR(64) NOT NULL,
    currency VARCHAR(6) NOT NULL,
    provider VARCHAR(64) NOT NULL,
    amount INT NOT NULL,
    payment_dt INT NOT NULL,
    bank VARCHAR(64) NOT NULL,
    delivery_cost INT NOT NULL,
    goods_total INT NOT NULL,
    custom_fee INT NOT NULL,
    FOREIGN KEY (order_uid) REFERENCES orders (uid) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS items (
    id BIGSERIAL PRIMARY KEY,
    order_uid VARCHAR(64) NOT NULL,
    chrt_id INT NOT NULL,
    track_number VARCHAR(64) NOT NULL,
    price INT NOT NULL,
    rid VARCHAR(64) NOT NULL,
    name VARCHAR(64) NOT NULL,
    sale INT NOT NULL,
    size VARCHAR(64) NOT NULL,
    total_price INT NOT NULL,
    nm_id INT NOT NULL,
    brand VARCHAR(64) NOT NULL,
    status INT NOT NULL,
    FOREIGN KEY (order_uid) REFERENCES orders (uid) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_deliveries ON deliveries (order_uid);
CREATE INDEX IF NOT EXISTS idx_payments ON payments (order_uid);
CREATE INDEX IF NOT EXISTS idx_items ON items (order_uid);
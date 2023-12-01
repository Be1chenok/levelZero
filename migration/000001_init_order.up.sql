CREATE TABLE IF NOT EXISTS orders(
    order_uid VARCHAR(64) PRIMARY KEY,
    track_number VARCHAR(64),
    entry VARCHAR(64),
    locale VARCHAR(6),
    internal_signature VARCHAR(64),
    customer_id VARCHAR(64),
    delivery_service VARCHAR(64),
    shardkey VARCHAR(64),
    sm_id INT,
    date_created VARCHAR(64),
    oof_shard VARCHAR(64)
);

CREATE TABLE IF NOT EXISTS delivery(
    order_uid VARCHAR(64) PRIMARY KEY,
    name VARCHAR(64),
    phone VARCHAR(16),
    zip VARCHAR(255),
    city VARCHAR(255),
    address VARCHAR(255),
    region VARCHAR(255),
    email VARCHAR(255),
    FOREIGN KEY (order_uid) REFERENCES orders (order_uid)
);

CREATE TABLE IF NOT EXISTS payment(
    order_uid VARCHAR(64) PRIMARY KEY,
    transaction VARCHAR(64),
    request_id VARCHAR(64),
    currency VARCHAR(6),
    provider VARCHAR(64),
    amount INT,
    payment_dt INT,
    bank VARCHAR(64),
    delivery_cost INT,
    goods_total INT,
    custom_fee INT,
    FOREIGN KEY (order_uid) REFERENCES orders (order_uid)
);

CREATE TABLE IF NOT EXISTS items (
    chrt_id INT PRIMARY KEY,
    track_number VARCHAR(64),
    price INT,
    rid VARCHAR(64),
    name VARCHAR(64),
    sale INT,
    size VARCHAR(64),
    total_price INT,
    nm_id INT,
    brand VARCHAR(64),
    status INT,
    order_uid VARCHAR(64),
    FOREIGN KEY (order_uid) REFERENCES orders (order_uid)
);
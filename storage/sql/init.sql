CREATE TABLE IF NOT EXISTS orders(
    order_uid UUID NOT NULL PRIMARY KEY,
    track_number VARCHAR(255) NOT NULL UNIQUE,
    "entry" VARCHAR(255) NOT NULL,
    locale VARCHAR(255) NOT NULL,
    internal_signature VARCHAR(255),
    customer_id VARCHAR(255) NOT NULL,
    delivery_service VARCHAR(255) NOT NULL,
    shardkey VARCHAR(255) NOT NULL,
    sm_id INTEGER NOT NULL,
    date_created TIMESTAMP NOT NULL,
    oof_shard VARCHAR(255) NOT NULL
);

CREATE TABLE IF NOT EXISTS items(
    order_uid UUID NOT NULL REFERENCES orders(order_uid),
    chrt_id INTEGER NOT NULL,
    track_number VARCHAR(255) NOT NULL REFERENCES orders(track_number),
    price DECIMAL NOT NULL CHECK(price >= 0),
    rid VARCHAR(255) NOT NULL,
    "name" VARCHAR(255) NOT NULL,
    sale DECIMAL CHECK(sale >= 0),
    "size" VARCHAR(255) NOT NULL,
    total_price DECIMAL CHECK(total_price >= 0),
    nm_id INTEGER NOT NULL,
    brand VARCHAR(255) NOT NULL,
    "status" INTEGER
);

CREATE TABLE IF NOT EXISTS payments(
    order_uid UUID NOT NULL REFERENCES orders(order_uid),
    "transaction" UUID NOT NULL PRIMARY KEY,
    request_id VARCHAR(255),
    currency VARCHAR(10) NOT NULL,
    "provider" VARCHAR(255) NOT NULL,
    amount INTEGER NOT NULL,
    payment_dt BIGINT NOT NULL,
    bank VARCHAR(255) NOT NULL,
    delivery_cost DECIMAL NOT NULL CHECK(delivery_cost >= 0),
    goods_total INTEGER NOT NULL CHECK(goods_total >= 0),
    custom_fee INTEGER NOT NULL
);

CREATE TABLE IF NOT EXISTS deliveries(
    order_uid UUID NOT NULL REFERENCES orders(order_uid),
    "name" VARCHAR(255) NOT NULL,
    phone VARCHAR(20),
    zip VARCHAR(255),
    city VARCHAR(255),
    "address" VARCHAR(255),
    region VARCHAR(255),
    email VARCHAR(255)
);
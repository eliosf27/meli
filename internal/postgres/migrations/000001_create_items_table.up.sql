CREATE TABLE IF NOT EXISTS order_history
(
    id         SERIAL  PRIMARY KEY,
    bundle_id  TEXT    NULL,
    courier_id INTEGER NULL,
    order_id   INTEGER NOT NULL,
    state      TEXT    NOT NULL,
    event      TEXT    NOT NULL,
    data       JSONB   NULL,
    created_at TIMESTAMP(0) DEFAULT NOW(),
    updated_at TIMESTAMP(0) DEFAULT NOW()
);

CREATE INDEX ix_order_history_order_id
    ON order_history (order_id);

CREATE INDEX ix_order_history_courier_id
    ON order_history (courier_id);

CREATE INDEX ix_order_history_bundle_id
    ON order_history (bundle_id);

CREATE INDEX ix_order_history_created_at
    ON order_history (created_at DESC);
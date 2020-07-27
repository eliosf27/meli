CREATE TABLE IF NOT EXISTS item
(
    id          SERIAL PRIMARY KEY,
    item_id     TEXT    NULL,
    title       TEXT    NULL,
    category_id TEXT    NULL,
    price       DECIMAL NULL,
    start_time  TEXT    NOT NULL,
    stop_time   TEXT    NOT NULL,
    created_at  TIMESTAMP(0) DEFAULT NOW(),
    updated_at  TIMESTAMP(0) DEFAULT NOW()
);

CREATE INDEX ix_items_item_id
    ON item (item_id);

CREATE INDEX ix_items_category_id
    ON item (category_id);

CREATE INDEX ix_items_created_at
    ON item (created_at DESC);
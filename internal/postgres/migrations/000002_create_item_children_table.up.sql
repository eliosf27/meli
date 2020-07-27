CREATE TABLE IF NOT EXISTS item_children
(
    id         SERIAL PRIMARY KEY,
    item_id    TEXT NULL,
    stop_time  TEXT NOT NULL,
    created_at TIMESTAMP(0) DEFAULT NOW(),
    updated_at TIMESTAMP(0) DEFAULT NOW()
);

CREATE INDEX ix_item_children_item_id
    ON item_children (item_id);

CREATE INDEX ix_item_children_created_at
    ON item_children (created_at DESC);
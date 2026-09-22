CREATE TABLE collections (
    id TEXT PRIMARY KEY,
    title TEXT NOT NULL,
    icon TEXT,
    parent_id TEXT NULL,
    owner_id TEXT NOT NULL,
    sort_order INTEGER NOT NULL,
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL,

    CONSTRAINT collections_sort_order_nonnegative
        CHECK (sort_order >= 0),


    FOREIGN KEY (parent_id) REFERENCES collections(id),
    FOREIGN KEY (owner_id) REFERENCES accounts(id) ON DELETE CASCADE,

    CONSTRAINT collections_account_sort_order_unique
        UNIQUE (owner_id, sort_order)
        DEFERRABLE INITIALLY DEFERRED
);

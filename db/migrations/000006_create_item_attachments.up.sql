CREATE TABLE item_attachments (
    id TEXT PRIMARY KEY,

    owner_id TEXT NOT NULL,
    item_id TEXT NOT NULL,
    asset_id TEXT NOT NULL,

    role TEXT NOT NULL DEFAULT 'attachment',
    position INTEGER NOT NULL,

    caption TEXT,
    alt_text TEXT,
    content_ref TEXT,

    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL,

    CONSTRAINT item_attchments_role_valid CHECK (
        role IN ('primary', 'cover', 'inline', 'attachment')
    ),

    CONSTRAINT item_attachments_position_nonnegative CHECK (
        position >= 0
    ),

    CONSTRAINT item_attachments_item_fkey
        FOREIGN KEY (owner_id, item_id)
        REFERENCES items (owner_id, id)
        ON DELETE CASCADE,

    CONSTRAINT item_attachments_asset_fkey
        FOREIGN KEY (owner_id, asset_id)
        REFERENCES assets (owner_id, id)
        ON DELETE NO ACTION
        DEFERRABLE INITIALLY DEFERRED,

        CONSTRAINT item_attachments_position_unique
            UNIQUE (item_id, position)
            DEFERRABLE INITIALLY DEFERRED
);

CREATE INDEX item_attachments_owner_asset_idx
    ON item_attachments (owner_id, asset_id);

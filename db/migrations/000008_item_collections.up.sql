ALTER TABLE collections
    ADD CONSTRAINT collections_owner_id_id_unique
    UNIQUE (owner_id, id);

ALTER TABLE items
    ADD COLUMN collection_id TEXT;

ALTER TABLE items
    ADD CONSTRAINT items_collection_owner_fk
    FOREIGN KEY (owner_id, collection_id)
    REFERENCES collections(owner_id, id);

CREATE INDEX items_owner_collection_created_at_idx
    ON items (
        owner_id,
        collection_id,
        created_at DESC,
        id DESC
    );

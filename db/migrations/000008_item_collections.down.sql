DROP INDEX items_owner_collection_created_at_idx;

ALTER TABLE items
    DROP CONSTRAINT items_collection_owner_fk,
    DROP COLUMN collection_id;

ALTER TABLE collections
    DROP CONSTRAINT collections_owner_id_id_unique;

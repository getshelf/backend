CREATE TABLE items (
    id TEXT PRIMARY KEY,
    owner_id TEXT NOT NULL REFERENCES accounts(id) ON DELETE CASCADE,

    kind TEXT NOT NULL,
    title TEXT,
    description TEXT,

    -- Captured content
    content_text TEXT,
    content_html TEXT,
    content_structured JSONB,

    -- Original source
    source_url TEXT,
    canonical_url TEXT,
    provider TEXT,
    external_id TEXT,

    author_name TEXT,
    author_handle TEXT,
    author_url TEXT,
    published_at TIMESTAMPTZ,

    -- User annotation
    user_note TEXT,
    is_favorite BOOLEAN NOT NULL DEFAULT FALSE,

    metadata JSONB NOT NULL DEFAULT '{}'::jsonb,

    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL,

    CONSTRAINT items_kind_valid CHECK (
        kind IN (
            'link', 'article', 'post', 'image',
            'video', 'audio', 'document', 'note'
        )
    ),

    CONSTRAINT items_metadata_object CHECK (
        jsonb_typeof(metadata) = 'object'
    ),

    CONSTRAINT items_owner_id_id_unique UNIQUE (owner_id, id)
);

CREATE INDEX items_owner_created_at_idx
    ON items (owner_id, created_at DESC, id DESC)

CREATE TABLE assets (
    id TEXT PRIMARY KEY,
    owner_id TEXT NOT NULL REFERENCES accounts(id) ON DELETE CASCADE,

    kind TEXT NOT NULL,
    original_url TEXT,

    storage_key TEXT,

    file_name TEXT,
    mime_type TEXT,
    byte_size BIGINT,
    chechjsum_sha256 TEXT,

    width INTEGER,
    height INTEGER,
    duration_ms BIGINT,

    download_status TEXT NOT NULL DEFAULT 'pending',
    download_error TEXT,

    metadata JSONB NOT NULL DEFAULT '{}'::jsonb,

    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL,

    CONSTRAINT assets_kind_valid CHECK (
        kind IN ('image', 'video', 'audio', 'document', 'other')
    ),

    CONSTRAINT assets_download_status_valid CHECK (
        download_status IN ('pending', 'downloading', 'ready', 'failed')
    ),

    CONSTRAINT assets_ready_has_storage CHECK (
        donwload_status <> 'ready'
        OR (
            storage_key IS NOT NULL
            AND length(storage_key) > 0
        )
    ),

    CONSTRAINT assets_byte_size_nonnegative CHECK (byte_size >= 0),
    CONSTRAINT assets_width_positive CHECK (width > 0),
    CONSTRAINT assets_height_positive CHECK (height > 0),
    CONSTRAINT assets_duration_nonnegative CHECK (duration_ms >= 0),

    CONSTRAINT assets_checksum_sha256_valid CHECK (
        checksum_sha256 ~ '^[0-9a-f]{64}$'
    ),

    CONSTRAINT assets_metadata_object CHECK (
        jsonb_typeof(metadata) = 'object'
    ),

    CONSTRAINT assets_owner_id_id_unique UNIQUE (owner_id, id)
)

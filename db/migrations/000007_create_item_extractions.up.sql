CREATE TABLE item_extractions (
    item_id TEXT PRIMARY KEY REFERENCES items(id) ON DELETE CASCADE,

    status TEXT NOT NULL DEFAULT 'pending',

    extractor TEXT,
    extractor_version TEXT,

    attempt_count INTEGER NOT NULL DEFAULT 0,
    last_error TEXT,
    warnings TEXT[] NOT NULL DEFAULT ARRAY[]::text[],

    raw_response_key TEXT,

    started_at TIMESTAMPTZ,
    extracted_at TIMESTAMPTZ,

    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL,

    CONSTRAINT item_extractions_status_valid CHECK (
        status IN ('pending', 'processing', 'ready', 'partial', 'failed')
    ),

    CONSTRAINT item_extractions_attempt_count_nonnegative CHECK (
        attempt_count >= 0
    )
)

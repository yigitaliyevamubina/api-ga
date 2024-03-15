CREATE TABLE IF NOT EXISTS comments (
    id uuid NOT NULL,
    content TEXT NOT NULL,
    post_id uuid NOT NULL,
    owner_id uuid NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP
);
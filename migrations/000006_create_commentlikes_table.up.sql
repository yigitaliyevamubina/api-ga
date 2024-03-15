CREATE TABLE IF NOT EXISTS commentlikes (
    id SERIAL PRIMARY KEY,
    user_id UUID NOT NULL UNIQUE,
    comment_id UUID NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,
    CONSTRAINT unique_comment_user_combination UNIQUE(user_id, comment_id)
);

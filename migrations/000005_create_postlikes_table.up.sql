CREATE TABLE IF NOT EXISTS postlikes (
    id SERIAL PRIMARY KEY,
    user_id UUID NOT NULL,
    post_id UUID NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,
    CONSTRAINT unique_post_user_combination UNIQUE(user_id, post_id)
);
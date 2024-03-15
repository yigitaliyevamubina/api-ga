CREATE TABLE IF NOT EXISTS posts (
    id uuid NOT NULL,
    title TEXT NOT NULL,
    image_url TEXT NOT NULL,
    owner_id uuid NOT NULL
);
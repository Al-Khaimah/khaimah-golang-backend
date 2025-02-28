CREATE TABLE podcast (
    id UUID PRIMARY KEY,
    title VARCHAR NOT NULL,
    description VARCHAR,
    image_url TEXT,
    content_url TEXT,
    likes_count INT DEFAULT 0
);
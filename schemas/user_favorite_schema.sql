CREATE TABLE user_favorite (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL,
    podcast_id UUID NOT NULL,
    CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES "user"(id),
    CONSTRAINT fk_podcast FOREIGN KEY (podcast_id) REFERENCES podcast(id),
    UNIQUE (user_id, podcast_id)
);
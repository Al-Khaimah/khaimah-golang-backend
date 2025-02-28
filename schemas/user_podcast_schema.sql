CREATE TABLE user_podcast (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL,
    podcast_id UUID NOT NULL,
    category_id UUID NOT NULL,
    resume_position INT DEFAULT 0,
    is_completed BOOLEAN DEFAULT false,
    CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES "user"(id),
    CONSTRAINT fk_podcast FOREIGN KEY (podcast_id) REFERENCES podcast(id),
    CONSTRAINT fk_category FOREIGN KEY (category_id) REFERENCES category(id),
    UNIQUE (user_id, podcast_id)
);
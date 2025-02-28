CREATE TABLE user_category (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL,
    category_id UUID NOT NULL,
    CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES "user"(id),
    CONSTRAINT fk_category FOREIGN KEY (category_id) REFERENCES category(id),
    UNIQUE (user_id, category_id)
);
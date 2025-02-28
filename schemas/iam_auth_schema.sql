CREATE TABLE iam_auth (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES "user"(id),
    password VARCHAR NOT NULL,
    is_active BOOLEAN DEFAULT true,
    last_login TIMESTAMP,
    CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES "user"(id)
);
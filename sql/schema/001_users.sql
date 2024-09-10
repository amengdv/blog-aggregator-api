-- +goose Up
CREATE TABLE users (
    id UUID PRIMARY KEY,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    email TEXT UNIQUE NOT NULL,
    name TEXT NOT NULL,
    password TEXT NOT NULL,
    refresh_token VARCHAR(64) UNIQUE,
    tkn_expires_at TIMESTAMP
);

-- +goose Down
DROP TABLE users;

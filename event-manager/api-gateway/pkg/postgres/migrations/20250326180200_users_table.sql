-- +goose Up
CREATE TABLE users
(
    user_id    SERIAL PRIMARY KEY,
    username   VARCHAR(30) NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    email VARCHAR NOT NULL UNIQUE,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    is_deleted BOOLEAN NOT NULL DEFAULT FALSE,
    role VARCHAR NOT NULL,
    created_at TIMESTAMP   NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- +goose Down
DROP TABLE users;

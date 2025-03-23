-- +goose Up
CREATE TABLE comments (
    comment_id SERIAL PRIMARY KEY,
    event_id INT NOT NULL,
    sender_id INT NOT NULL,
    content TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    is_deleted BOOLEAN NOT NULL DEFAULT FALSE,
    is_read BOOLEAN NOT NULL DEFAULT FALSE
--     FOREIGN KEY (event_id) REFERENCES event(event_id) ON DELETE CASCADE,
--     FOREIGN KEY (sender_id) REFERENCES "user"(user_id) ON DELETE CASCADE
);

-- +goose Down
DROP TABLE comments;

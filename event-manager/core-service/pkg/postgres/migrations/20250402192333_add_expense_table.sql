-- +goose Up
CREATE TABLE expense
(
    expense_id   SERIAL PRIMARY KEY,
    event_id     INT,
    created_by   INT,
    description  TEXT       NOT NULL,
    amount       FLOAT      NOT NULL,
    currency     VARCHAR(3) NOT NULL DEFAULT 'RUB',
    split_method TEXT       NOT NULL,
    created_at   TIMESTAMP  NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (event_id) REFERENCES events(event_id) ON DELETE CASCADE,
    FOREIGN KEY (created_by) REFERENCES users(user_id) ON DELETE CASCADE
);

-- +goose Down
DROP TABLE expense;

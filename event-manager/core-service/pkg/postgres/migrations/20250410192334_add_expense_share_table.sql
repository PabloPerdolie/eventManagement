-- +goose Up
CREATE TABLE expense_share
(
    share_id    SERIAL PRIMARY KEY,
    expense_id  INT          NOT NULL,
    user_id     INT          NOT NULL,
    amount      FLOAT        NOT NULL,
    is_paid     BOOLEAN      NOT NULL DEFAULT FALSE,
    paid_at     TIMESTAMP,
    FOREIGN KEY (expense_id) REFERENCES expense(expense_id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE
);

-- +goose Down
DROP TABLE expense_share; 
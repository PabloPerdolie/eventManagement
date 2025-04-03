-- +goose Up
CREATE TABLE tasks
(
    task_id      SERIAL PRIMARY KEY,
    event_id     INT,
    parent_id    INT REFERENCES tasks(task_id) DEFAULT NULL, -- ?
    title        TEXT        NOT NULL,
    description  TEXT        NOT NULL,
    story_points INT,
    priority     VARCHAR(20),
    status       VARCHAR(20) NOT NULL,
    created_at   TIMESTAMP   NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (event_id) REFERENCES events (event_id) ON DELETE CASCADE
);

CREATE TABLE task_assignment
(
    task_assignment_id SERIAL PRIMARY KEY,
    task_id            INT,
    user_id            INT,
    assigned_at        TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    completed_at       TIMESTAMP,
    FOREIGN KEY (task_id) REFERENCES tasks (task_id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users (user_id) ON DELETE CASCADE
);

-- +goose Down
DROP TABLE task_assignment;

DROP TABLE tasks;

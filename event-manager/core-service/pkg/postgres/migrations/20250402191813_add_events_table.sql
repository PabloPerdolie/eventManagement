-- +goose Up
CREATE TABLE events
(
    event_id     SERIAL PRIMARY KEY,
    organizer_id INT,
    title        TEXT        NOT NULL,
    description  TEXT        NOT NULL,
    start_date   TIMESTAMP   NOT NULL,
    end_date     TIMESTAMP   NOT NULL,
    location     TEXT,
    status       VARCHAR(10) NOT NULL DEFAULT 'planned',
    created_at   TIMESTAMP   NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (organizer_id) REFERENCES users (user_id) ON DELETE CASCADE
);

CREATE TABLE event_participant
(
    event_participant_id SERIAL PRIMARY KEY,
    event_id INT,
    user_id INT,
    role VARCHAR(15) NOT NULL,
    joined_at TIMESTAMP, --?
    is_confirmed bool, --?
    FOREIGN KEY (event_id) REFERENCES events (event_id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users (user_id) ON DELETE CASCADE
);

-- +goose Down
DROP TABLE event_participant;

DROP TABLE events;

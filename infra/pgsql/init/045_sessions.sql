CREATE TABLE sessions
(
    id            TEXT PRIMARY KEY,
    user_id       TEXT        NOT NULL,
    user_agent    TEXT        NOT NULL,
    status        TEXT        NOT NULL,
    expiry        TIMESTAMPTZ NOT NULL,
    refresh_token TEXT        NOT NULL
);
CREATE TABLE users (
    id    BIGSERIAL PRIMARY KEY,
    email TEXT      NOT NULL CHECK (char_length(email) <= 255)
);

CREATE TABLE user_actions (
    id          BIGSERIAL   PRIMARY KEY,
    user_id     BIGINT      NOT NULL REFERENCES users (id),
    action      TEXT        NOT NULL CHECK (char_length(action) < 100),
    occurred_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
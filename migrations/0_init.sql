-- +migrate Up
CREATE TABLE users
(
    id         BIGSERIAL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
    deleted_at TIMESTAMP WITH TIME ZONE NULL,

    email      TEXT                     NOT NULL,
    passhash   TEXT                     NOT NULL,

    PRIMARY KEY (id)
);
CREATE UNIQUE INDEX users_email_idx ON users (email);

comment
    ON COLUMN users.email IS 'email';
comment
    ON COLUMN users.passhash IS 'password hash';

CREATE TABLE sessions
(
    id         BIGSERIAL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
    deleted_at TIMESTAMP WITH TIME ZONE NULL,

    user_id    BIGINT                   NOT NULL,
    token      TEXT                     NOT NULL,
    extra      JSONB                    NOT NULL,
    PRIMARY KEY (id)
);
CREATE UNIQUE INDEX sessions_access_token_idx ON sessions (token);
CREATE INDEX sessions_user_id_idx ON sessions (user_id);

comment
    ON COLUMN sessions.user_id IS 'user id';
comment
    ON COLUMN sessions.token IS 'access token';
comment
    ON COLUMN sessions.extra IS 'extra data: ip, user-agent etc';

-- +migrate Down
DROP TABLE IF EXISTS sessions;
DROP TABLE IF EXISTS users;

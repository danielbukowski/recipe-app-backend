-- +goose Up
-- +goose StatementBegin
CREATE TABLE users(
    user_id UUID PRIMARY KEY,
    email TEXT NOT NULL UNIQUE,
    password TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX users_email_idx ON users(email) INCLUDE (password);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX idx_email;
DROP table users;
-- +goose StatementEnd

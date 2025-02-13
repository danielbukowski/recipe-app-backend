-- +goose Up
-- +goose StatementBegin
CREATE TABLE email_verifications(
    user_id UUID REFERENCES users(user_id) ON DELETE CASCADE,
    code TEXT NOT NULL,
    expires_at TIMESTAMP NOT NULL DEFAULT NOW(),
    PRIMARY KEY (user_id, code)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE email_verifications;
-- +goose StatementEnd

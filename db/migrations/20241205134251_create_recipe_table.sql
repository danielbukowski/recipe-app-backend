-- +goose Up
-- +goose StatementBegin
CREATE TABLE recipes(
    recipe_id uuid PRIMARY KEY,
    title text NOT NULL,
    content text NOT NULL,
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE recipes;
-- +goose StatementEnd 

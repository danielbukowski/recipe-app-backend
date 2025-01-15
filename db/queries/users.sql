-- name: CreateUser :exec
INSERT INTO users (
    user_id,
    email,
    password
) VALUES ($1, $2, $3);

-- name: GetUserByEmail :one
SELECT email, password FROM users 
    WHERE email = $1
    LIMIT 1;
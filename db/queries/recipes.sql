-- name: CreateRecipe :one
INSERT INTO recipes (
    recipe_id, 
    title, 
    content
) VALUES ($1, $2, $3)
RETURNING recipe_id;

-- name: GetRecipeById :one
SELECT * FROM recipes
    WHERE recipe_id = $1 LIMIT 1;

-- name: UpdateRecipeById :exec
UPDATE recipes
    SET title = $3, content = $4, updated_at = sqlc.arg(new_updated_at)
    WHERE recipe_id = $1 AND updated_at = $2;


-- name: DeleteRecipeById :exec
DELETE FROM recipes 
    WHERE recipe_id = $1;


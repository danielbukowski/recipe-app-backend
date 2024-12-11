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
    SET title = $2, content = $3
    WHERE recipe_id = $1;


-- name: DeleteRecipeById :exec
DELETE FROM recipes 
    WHERE recipe_id = $1;


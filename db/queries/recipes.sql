-- name: CreateRecipe :exec
INSERT INTO recipes (
    recipe_id, 
    title, 
    content, 
    updated_at
)
    VALUES ($1, $2, $3, $4);

-- name: GetRecipeById :one
SELECT * FROM recipes
    WHERE recipe_id = $1 LIMIT 1;

-- name: UpdateRecipeById :exec
UPDATE recipes
    SET title = $1, content = $2
    WHERE recipe_id = $3;


-- name: DeleteRecipeById :exec
DELETE FROM recipes 
    WHERE recipe_id = $1;


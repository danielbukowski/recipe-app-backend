// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0

package sqlc

import (
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type Recipe struct {
	RecipeID  uuid.UUID
	Title     string
	Content   string
	CreatedAt pgtype.Timestamp
	UpdatedAt pgtype.Timestamp
}

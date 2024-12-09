package recipe

import (
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

const databaseConnectionTimeout = 4 * time.Second

type service struct {
	Logger *zap.Logger
	dbpool *pgxpool.Pool
}

func NewService(logger *zap.Logger, dbpool *pgxpool.Pool) *service {
	return &service{
		Logger: logger,
		dbpool: dbpool,
	}
}

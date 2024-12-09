package recipe

import (
	"go.uber.org/zap"
)

type controller struct {
	Logger        *zap.Logger
	recipeService *service
}

func NewController(logger *zap.Logger, recipeService *service) *controller {
	return &controller{
		Logger:        logger,
		recipeService: recipeService,
	}
}


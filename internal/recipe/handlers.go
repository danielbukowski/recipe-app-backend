package recipe

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/danielbukowski/recipe-app-backend/internal/shared"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

// TODO: test out the new handlers

type handler struct {
	logger        *zap.Logger
	recipeService recipeService
}

type recipeService interface {
	GetRecipeById(context.Context, uuid.UUID) (RecipeResponse, error)
	DeleteRecipeById(context.Context, uuid.UUID) error
	CreateNewRecipe(context.Context, NewRecipeRequest) (uuid.UUID, error)
	UpdateRecipeById(context.Context, uuid.UUID, time.Time, UpdateRecipeRequest) error
}

func NewHandler(logger *zap.Logger, recipeService recipeService) *handler {
	return &handler{
		logger:        logger,
		recipeService: recipeService,
	}
}

// @Summary		Create a recipe
// @Description	Insert a new recipe by providing a request body with a title and a content for the recipe.
// @Tags			recipes
//
// @Accept			json
// @Produce		json
// @Param			NewRecipeRequest	body		recipe.NewRecipeRequest	true	"Request body with title and content"
//
// @Success		201					{object}	shared.CommonResponse
// @Failure		400					{object}	shared.CommonResponse
// @Failure		408					{object}	shared.CommonResponse
// @Failure		500					{object}	shared.CommonResponse
//
// @Router			/api/v1/recipes [POST]
func (h *handler) createRecipe(c echo.Context) error {
	var requestBody = NewRecipeRequest{}

	if err := c.Bind(&requestBody); err != nil {
		return c.JSON(http.StatusBadRequest, shared.CommonResponse{Message: "missing a valid JSON request body"})
	}

	if err := c.Validate(&requestBody); err != nil {
		return err
	}

	recipeId, err := h.recipeService.CreateNewRecipe(c.Request().Context(), requestBody)
	if err != nil {
		return err
	}

	h.logger.Info("saved a new recipe to database")

	c.Response().Header().Add("Location", fmt.Sprintf("http://localhost:8080/api/v1/recipes/%v", recipeId.String()))
	return c.JSON(http.StatusCreated, shared.CommonResponse{Message: "successfully saved a recipe"})
}

// @Summary		Update a recipe
// @Description	Update a title or a content of a recipe by ID.
// @Tags			recipes
//
// @Accept			json
// @Produce		json
// @Param			id					path	string						true	"UUID for a recipe resource"
// @Param			UpdateRecipeRequest	body	recipe.UpdateRecipeRequest	true	"Request body for updating title and content fields of a recipe"
//
// @Success		204
// @Failure		400	{object}	shared.CommonResponse
// @Failure		408	{object}	shared.CommonResponse
// @Failure		500	{object}	shared.CommonResponse
//
// @Router			/api/v1/recipes/{id} [PUT]
func (h *handler) updateRecipeById(c echo.Context) error {
	recipeIdParam := c.Param("id")

	if recipeIdParam == "" {
		return c.JSON(http.StatusBadRequest, shared.CommonResponse{Message: "missing ID param for recipe"})
	}

	recipeId, err := uuid.Parse(recipeIdParam)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, shared.CommonResponse{Message: "the received ID is not a valid UUID"})
	}

	var requestBody = UpdateRecipeRequest{}

	if err := c.Bind(&requestBody); err != nil {
		return c.JSON(http.StatusBadRequest, shared.CommonResponse{Message: "missing a valid JSON request body"})
	}

	recipeFromDb, err := h.recipeService.GetRecipeById(c.Request().Context(), recipeId)
	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			return echo.NewHTTPError(http.StatusNotFound, shared.CommonResponse{Message: "could not find a recipe with this id"})
		default:
			// TODO: check what is the response of it
			return err
		}
	}

	if requestBody.Title == "" {
		requestBody.Title = recipeFromDb.Title
	}

	if requestBody.Content == "" {
		requestBody.Content = recipeFromDb.Content
	}

	if err := c.Validate(requestBody); err != nil {
		// TODO: customize this error to return a better descriptive message of wrong fields
		return err
	}

	err = h.recipeService.UpdateRecipeById(c.Request().Context(), recipeId, recipeFromDb.UpdatedAt, requestBody)
	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			return echo.NewHTTPError(http.StatusConflict, shared.CommonResponse{Message: "conflict occurred when trying to update a recipe"})
		default:
			// TODO: think about returning here a http internal error
			return err
		}
	}

	h.logger.Info("updated a recipe",
		zap.String("recipeId", recipeId.String()),
	)

	c.NoContent(http.StatusNoContent)
	return nil
}

// @Summary		Delete a recipe
// @Description	Delete a recipe by ID.
// @Tags			recipes
//
// @Accept			json
// @Produce		json
// @Param			id	path	string	true	"UUID for a recipe"
//
// @Success		204
// @Failure		400	{object}	shared.CommonResponse
// @Failure		408	{object}	shared.CommonResponse
// @Failure		500	{object}	shared.CommonResponse
//
// @Router			/api/v1/recipes/{id} [DELETE]
func (h *handler) deleteRecipeById(c echo.Context) error {
	recipeIdParam := c.Param("id")

	if recipeIdParam == "" {
		return c.JSON(http.StatusBadRequest, shared.CommonResponse{Message: "missing ID param for recipe"})
	}

	recipeId, err := uuid.Parse(recipeIdParam)
	if err != nil {
		return c.JSON(http.StatusBadRequest, shared.CommonResponse{Message: "the received ID is not a valid UUID"})
	}

	err = h.recipeService.DeleteRecipeById(c.Request().Context(), recipeId)
	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			return echo.NewHTTPError(http.StatusNoContent)
		default:
			// TODO: think about returning here a http internal error or think how to handle this in the global exception handler
			return err
		}

	}

	h.logger.Info("deleted a recipe from database", zap.String("recipeId", recipeId.String()))

	return c.NoContent(http.StatusNoContent)
}

// @Summary		Get a recipe
// @Description	Get a recipe by ID.
// @Tags			recipes
//
// @Accept			json
// @Produce		json
// @Param			id	path		string	true	"UUID for a recipe"
//
// @Success		200	{object}	shared.DataResponse[recipe.RecipeResponse]
// @Failure		400	{object}	shared.CommonResponse
// @Failure		408	{object}	shared.CommonResponse
// @Failure		500	{object}	shared.CommonResponse
//
// @Router			/api/v1/recipes/{id} [GET]
func (h *handler) getRecipeById(c echo.Context) error {
	recipeIdParam := c.Param("id")

	if recipeIdParam == "" {
		return c.JSON(http.StatusBadRequest, shared.CommonResponse{Message: "missing ID param for recipe"})
	}

	recipeId, err := uuid.Parse(recipeIdParam)
	if err != nil {
		return c.JSON(http.StatusBadRequest, shared.CommonResponse{Message: "the received ID is not a valid UUID"})
	}

	recipe, err := h.recipeService.GetRecipeById(c.Request().Context(), recipeId)
	// TO REMEMBER: if know type of a error, then return JSON
	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			return c.JSON(http.StatusNotFound, shared.CommonResponse{Message: "could not find recipe with this UUID"})
		default:
			// TODO: again think about how to handle this error better
			return err
		}
	}

	return c.JSON(http.StatusOK, shared.DataResponse[RecipeResponse]{Data: recipe})
}

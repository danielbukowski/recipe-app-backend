package recipe

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/danielbukowski/recipe-app-backend/internal/shared"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

var cachedRecipeKeyPrefix = "recipe_"

type handler struct {
	logger        *zap.Logger
	cache         cacheStorage
	recipeService recipeService
}

type recipeService interface {
	GetRecipeById(context.Context, uuid.UUID) (RecipeResponse, error)
	DeleteRecipeById(context.Context, uuid.UUID) error
	CreateNewRecipe(context.Context, NewRecipeRequest) (uuid.UUID, error)
	UpdateRecipeById(context.Context, uuid.UUID, time.Time, UpdateRecipeRequest) error
}

type cacheStorage interface {
	InsertItem(key string, value []byte, expiration int32) error
	GetItem(key string) ([]byte, error)
	DeleteItem(key string) error
}

func NewHandler(logger *zap.Logger, cacheStorage cacheStorage, recipeService recipeService) *handler {
	return &handler{
		logger:        logger,
		cache:         cacheStorage,
		recipeService: recipeService,
	}
}

// CreateRecipe godoc
//
//	@Summary		Create a new recipe
//	@Description	Insert a new recipe by providing a request body with title and content for the recipe you want to save.
//	@Tags			recipes
//
//	@Accept			json
//	@Produce		json
//	@Param			NewRecipeRequest	body		recipe.NewRecipeRequest				true	"Request body with title and content."
//
//	@Success		201					{object}	shared.CommonResponse				"Recipe saved successfully."
//	@Failure		400					{object}	validator.ValidationErrorResponse	"Invalid data provided."
//	@Failure		404					{object}	shared.CommonResponse				"Recipe not found."
//
//	@Router			/api/v1/recipes [POST]
func (h *handler) CreateRecipe(c echo.Context) error {
	if err := shared.ValidateJSONContentType(c); err != nil {
		return err
	}

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

// UpdateRecipeByID godoc
//
//	@Summary		Update a recipe
//	@Description	Update title or content of a recipe by UUID.
//	@Tags			recipes
//
//	@Accept			json
//
//	@Produce		json
//	@Param			id					path		string						true	"UUID of a recipe."
//	@Param			UpdateRecipeRequest	body		recipe.UpdateRecipeRequest	true	"Request body with title and content for updating a recipe."
//
//	@Success		204					"Recipe  	updated successfully."
//	@Failure		400					{object}	validator.ValidationErrorResponse	"Invalid data provided."
//	@Failure		409					{object}	shared.CommonResponse				"Database conflict occurred when trying to saving a recipe."
//
//	@Router			/api/v1/recipes/{id} [PUT]
func (h *handler) UpdateRecipeById(c echo.Context) error {
	if err := shared.ValidateJSONContentType(c); err != nil {
		return err
	}

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
		return nil
	}

	if requestBody.Title == "" {
		requestBody.Title = recipeFromDb.Title
	}

	if requestBody.Content == "" {
		requestBody.Content = recipeFromDb.Content
	}

	if err := c.Validate(requestBody); err != nil {
		return err
	}

	err = h.recipeService.UpdateRecipeById(c.Request().Context(), recipeId, recipeFromDb.UpdatedAt, requestBody)
	if err != nil {
		return err
	}

	if err = h.cache.DeleteItem(cachedRecipeKeyPrefix + recipeId.String()); err != nil {
		h.logger.Error("failed to delete a recipe from the cache", zap.String("recipe_id", recipeId.String()), zap.Error(err))
	}

	h.logger.Info("successfully updated a recipe", zap.String("recipeId", recipeId.String()))

	return c.NoContent(http.StatusNoContent)
}

// DeleteRecipeByID godoc
//
//	@Summary		Delete a recipe
//	@Description	Delete a recipe by ID.
//	@Tags			recipes
//
//	@Produce		json
//	@Param			id	path	string	true	"UUID for a recipe"
//
//	@Success		204	"Recipe deleted successfully."
//	@Failure		400	{object}	validator.ValidationErrorResponse	"Invalid data provided."
//
//	@Router			/api/v1/recipes/{id} [DELETE]
func (h *handler) DeleteRecipeById(c echo.Context) error {
	recipeIdParam := c.Param("id")

	if recipeIdParam == "" {
		return c.JSON(http.StatusBadRequest, shared.CommonResponse{Message: "missing ID param for recipe"})
	}

	recipeId, err := uuid.Parse(recipeIdParam)
	if err != nil {
		return c.JSON(http.StatusBadRequest, shared.CommonResponse{Message: "the received ID is not a valid UUID"})
	}

	if err := h.recipeService.DeleteRecipeById(c.Request().Context(), recipeId); err != nil {
		return err
	}

	if err = h.cache.DeleteItem(cachedRecipeKeyPrefix + recipeId.String()); err != nil {
		h.logger.Error("failed to delete a recipe from the cache", zap.String("recipe_id", recipeId.String()), zap.Error(err))
	}

	h.logger.Info("successfully deleted a recipe from database", zap.String("recipeId", recipeId.String()))

	return c.NoContent(http.StatusNoContent)
}

// GetRecipeByID godoc
//
//	@Summary		Get a recipe
//	@Description	Get a recipe by ID.
//	@Tags			recipes
//
//	@Produce		json
//
//	@Param			id	path		string										true	"UUID for a recipe"
//
//	@Success		200	{object}	shared.DataResponse[recipe.RecipeResponse]	"Recipe fetched successfully."
//	@Failure		400	{object}	validator.ValidationErrorResponse			"Invalid data provided."
//	@Failure		404	{object}	shared.CommonResponse						"Recipe is not found."
//
//	@Router			/api/v1/recipes/{id} [GET]
func (h *handler) GetRecipeById(c echo.Context) error {
	recipeIdParam := c.Param("id")

	if recipeIdParam == "" {
		return c.JSON(http.StatusBadRequest, shared.CommonResponse{Message: "missing ID param for recipe"})
	}

	recipeId, err := uuid.Parse(recipeIdParam)
	if err != nil {
		return c.JSON(http.StatusBadRequest, shared.CommonResponse{Message: "the received ID is not a valid UUID"})
	}

	// Check for a recipe in the cache.
	if cachedRecipe, err := h.cache.GetItem(cachedRecipeKeyPrefix + recipeId.String()); err == nil {
		var recipe RecipeResponse

		if err := json.Unmarshal(cachedRecipe, &recipe); err == nil {
			return c.JSON(http.StatusOK, shared.DataResponse[RecipeResponse]{Data: recipe})
		}
	}

	recipe, err := h.recipeService.GetRecipeById(c.Request().Context(), recipeId)
	if err != nil {
		return err
	}

	// Save the recipe to the cache.
	encodedRecipe, err := json.Marshal(recipe)
	if err != nil {
		h.logger.Error("failed to encoded a recipe to the cache", zap.Error(err))

		return c.JSON(http.StatusOK, shared.DataResponse[RecipeResponse]{Data: recipe})
	}

	if err = h.cache.InsertItem(cachedRecipeKeyPrefix+recipeId.String(), encodedRecipe, 60*15); err != nil {
		h.logger.Error("failed to insert a recipe to the cache", zap.Error(err))
	}

	return c.JSON(http.StatusOK, shared.DataResponse[RecipeResponse]{Data: recipe})
}

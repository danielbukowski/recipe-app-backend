package auth

import (
	"context"
	"net/http"

	"github.com/danielbukowski/recipe-app-backend/internal/shared"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

type handler struct {
	userService    userService
	logger         *zap.Logger
	sessionStorage sessionStorage
}

type userService interface {
	CreateUser(context.Context, SignUpRequest) error

type sessionStorage interface {
	CreateNew(value []byte) (string, error)
}

func NewHandler(logger *zap.Logger, userService userService, sessionStorage sessionStorage) *handler {
	return &handler{
		userService:    userService,
		logger:         logger,
		sessionStorage: sessionStorage,
	}
}

func (h *handler) SignUp(c echo.Context) error {
	var requestBody = SignUpRequest{}

	if err := c.Bind(&requestBody); err != nil {
		return c.JSON(http.StatusBadRequest, shared.CommonResponse{Message: "missing a valid JSON request body"})
	}

	if err := c.Validate(&requestBody); err != nil {
		return err
	}

	err := h.userService.CreateUser(c.Request().Context(), requestBody)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, shared.CommonResponse{Message: "successfully create a user account"})
}

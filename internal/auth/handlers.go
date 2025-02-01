package auth

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/danielbukowski/recipe-app-backend/internal/shared"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

type handler struct {
	userService    userService
	logger         *zap.Logger
	sessionStorage sessionStorage
	domainName     string
	isDev          bool
}

type userService interface {
	CreateUser(context.Context, SignUpRequest) error
	SignIn(ctx context.Context, signInRequest SignInRequest) (SignInResponse, error)
}

type sessionStorage interface {
	CreateNew(value []byte) (string, error)
	Delete(c echo.Context)
	AttachSessionCookieToClient(sessionID string, c echo.Context)
}

func NewHandler(logger *zap.Logger, userService userService, sessionStorage sessionStorage, domainName string, isDev bool) *handler {
	return &handler{
		userService:    userService,
		logger:         logger,
		sessionStorage: sessionStorage,
		domainName:     domainName,
		isDev:          isDev,
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

func (h *handler) signIn(c echo.Context) error {
	var requestBody = SignInRequest{}

	if err := c.Bind(&requestBody); err != nil {
		return c.JSON(http.StatusBadRequest, shared.CommonResponse{Message: "missing a valid JSON request body"})
	}

	if err := c.Validate(&requestBody); err != nil {
		return err
	}

	signInResponse, err := h.userService.SignIn(c.Request().Context(), requestBody)
	if err != nil {
		return err
	}

	jsonEncodedSession, err := json.Marshal(signInResponse)
	if err != nil {
		return err
	}

	sessionID, err := h.sessionStorage.CreateNew(jsonEncodedSession)
	if err != nil {
		return err
	}

	h.sessionStorage.AttachSessionCookieToClient(sessionID, c)

	return c.JSON(http.StatusOK, shared.CommonResponse{Message: "successfully sign in"})
}

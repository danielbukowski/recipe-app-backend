package auth

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/danielbukowski/recipe-app-backend/internal/shared"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

const (
	defaultSessionExpirationTime = 86400 * 14
)

type handler struct {
	userService       userService
	logger            *zap.Logger
	sessionStorage    sessionStorage
	isDev             bool
	sessionCookieName string
}

type userService interface {
	CreateUser(context.Context, SignUpRequest) error
	SignIn(ctx context.Context, signInRequest SignInRequest) (SignInResponse, error)
}

type sessionStorage interface {
	CreateNew(value []byte, expiration int32) (string, error)
	Delete(key string)
}

func NewHandler(logger *zap.Logger, userService userService, sessionStorage sessionStorage, isDev bool, sessionCookieName string) *handler {
	return &handler{
		userService:       userService,
		logger:            logger,
		sessionStorage:    sessionStorage,
		isDev:             isDev,
		sessionCookieName: sessionCookieName,
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

	sessionID, err := h.sessionStorage.CreateNew(jsonEncodedSession, defaultSessionExpirationTime)
	if err != nil {
		return err
	}

	cookie := http.Cookie{
		Name:     h.sessionCookieName,
		Value:    sessionID,
		Path:     "/",
		MaxAge:   defaultSessionExpirationTime,
		Secure:   !h.isDev,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	}

	c.SetCookie(&cookie)

	return c.JSON(http.StatusOK, shared.CommonResponse{Message: "successfully sign in"})
}

func (h *handler) signOut(c echo.Context) error {
	sessionID, err := c.Cookie(h.sessionCookieName)
	if err != nil {
		return err
	}

	h.sessionStorage.Delete(sessionID.Value)

	// Delete a session cookie from a client's browser
	cookie := http.Cookie{
		Name:     h.sessionCookieName,
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		Secure:   !h.isDev,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	}

	c.SetCookie(&cookie)

	return c.NoContent(http.StatusNoContent)
}

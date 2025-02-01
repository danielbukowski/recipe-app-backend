package auth

import (
	"github.com/labstack/echo/v4"
)

func (h *handler) RegisterRoutes(e *echo.Echo) {
	e.POST("api/v1/auth/signup", h.SignUp)
	e.POST("api/v1/auth/signin", h.signIn)
	e.POST("api/v1/auth/signout", h.signOut)
}

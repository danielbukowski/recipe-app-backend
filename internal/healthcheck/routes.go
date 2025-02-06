package healthcheck

import "github.com/labstack/echo/v4"

func (h *handler) RegisterRoutes(e *echo.Echo) {
	e.GET("api/v1/health", h.CheckHealth)
}

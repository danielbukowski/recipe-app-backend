package healthcheck

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type handler struct {
}

func NewHandler() *handler {
	return &handler{}
}

func (h *handler) checkHealth(c echo.Context) error {
	return c.NoContent(http.StatusOK)
}

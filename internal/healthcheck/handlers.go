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

// CheckHealth godoc
//
//	@Summary		Check health
//	@Description	Check the status of the recipe API.
//	@Tags			health
//
//	@Success		200	"The API is healthy."
//
//	@Router			/api/v1/health [GET]
func (h *handler) CheckHealth(c echo.Context) error {
	return c.NoContent(http.StatusOK)
}

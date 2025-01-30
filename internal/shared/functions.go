package shared

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func ValidateJSONContentType(c echo.Context) error {
	contentType := c.Request().Header.Get(echo.HeaderContentType)

	if contentType == "" || contentType != echo.MIMEApplicationJSON {
		return echo.NewHTTPError(http.StatusUnsupportedMediaType, CommonResponse{
			Message: "Only 'application/json' content type is allowed",
		})
	}

	return nil
}

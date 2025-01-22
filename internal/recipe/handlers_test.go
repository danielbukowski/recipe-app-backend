package recipe_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	mock_recipe "github.com/danielbukowski/recipe-app-backend/gen/_mocks/recipe"
	"github.com/danielbukowski/recipe-app-backend/internal/recipe"
	"github.com/danielbukowski/recipe-app-backend/internal/validator"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
)

func TestCreateRecipeHandler(t *testing.T) {
	testCases := []struct {
		name           string
		requestBody    string
		wantStatusCode int
	}{
		{
			name:           "no request body",
			requestBody:    "",
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name:           "request body is empty json",
			requestBody:    "{}",
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name: "content length of title and content in request body are too short",
			requestBody: `{
							"title": "cake",
							"content": "cook it"
						}`,
			wantStatusCode: http.StatusBadRequest,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			// given
			e := echo.New()
			e.Validator = validator.New()
			server := &http.Server{Handler: e}

			req := httptest.NewRequest(http.MethodPost, "/api/v1/recipes", nil)
			rec := httptest.NewRecorder()

			logger := zap.NewNop()
			recipeService := mock_recipe.NewMockRecipeService(gomock.NewController(t))

			handler := recipe.NewHandler(logger, recipeService)
			handler.RegisterRoutes(e)

			// when
			server.Handler.ServeHTTP(rec, req)

			// then
			assert.Equal(t, tc.wantStatusCode, rec.Code)
		})
	}
}

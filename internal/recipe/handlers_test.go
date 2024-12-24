package recipe_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/danielbukowski/recipe-app-backend/internal/recipe"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
)

func TestCreateRecipeHandler(t *testing.T) {
	testCases := []struct {
		name        string
		requestBody string
		statusCode  int
	}{
		{
			name:        "when there is no request body, then return UnsupportedMediaType status code",
			requestBody: "",
			statusCode:  http.StatusUnsupportedMediaType,
		},
		{
			name:        "when request body is empty json, then return BadRequest status code",
			requestBody: "{}",
			statusCode:  http.StatusBadRequest,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			// given
			l := zap.NewNop()

			recipeService := NewMockRecipeService(gomock.NewController(t))
			router := gin.New()
			handler := recipe.NewHandler(l, recipeService)
			resp := httptest.NewRecorder()

			// when
			req, err := http.NewRequest(http.MethodPost, "/api/v1/recipes", strings.NewReader(string(tc.requestBody)))
			assert.NoError(t, err)

			handler.RegisterRoutes(router)
			router.ServeHTTP(resp, req)

			// then
			assert.Equal(t, tc.statusCode, resp.Code)
		})
	}
}

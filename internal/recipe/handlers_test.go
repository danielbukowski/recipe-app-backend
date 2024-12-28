package recipe_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	mock_recipe "github.com/danielbukowski/recipe-app-backend/gen/_mocks/recipe"
	"github.com/danielbukowski/recipe-app-backend/internal/recipe"
	"github.com/gin-gonic/gin"
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
			wantStatusCode: http.StatusUnsupportedMediaType,
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
			logger := zap.NewNop()
			recipeService := mock_recipe.NewMockRecipeService(gomock.NewController(t))
			router := gin.New()
			handler := recipe.NewHandler(logger, recipeService)
			handler.RegisterRoutes(router)
			recorder := httptest.NewRecorder()

			// when
			req, err := http.NewRequest(http.MethodPost, "/api/v1/recipes", strings.NewReader(string(tc.requestBody)))
			assert.NoError(t, err)

			router.ServeHTTP(recorder, req)

			// then
			assert.Equal(t, tc.wantStatusCode, recorder.Code)
		})
	}
}

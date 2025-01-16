package auth

import (
	"context"
	"errors"
	"net/http"

	"github.com/danielbukowski/recipe-app-backend/internal/shared"
	"github.com/danielbukowski/recipe-app-backend/internal/validator"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgconn"
	"go.uber.org/zap"
)

type handler struct {
	userService userService
	logger      *zap.Logger
}

type userService interface {
	CreateUser(context.Context, SignUpRequest) error
}

func NewHandler(logger *zap.Logger, userService userService) *handler {
	return &handler{
		userService: userService,
		logger:      logger,
	}
}

func (h *handler) SignUp(c *gin.Context) {
	var requestBody = SignUpRequest{}

	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusUnsupportedMediaType, gin.H{
			"message": "missing JSON request body",
		})
		return
	}

	v := validator.New()

	if validateSignUpRequestBody(v, requestBody); !v.Valid() {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "request body did not pass the validation",
			"fields":  v.Errors,
		})
		return
	}

	err := h.userService.CreateUser(c, requestBody)
	if err != nil {
		var pgErr *pgconn.PgError

		switch {
		case errors.As(err, &pgErr):
			switch pgErr.Code {
			// unique_violation code
			case "23505":
				c.JSON(http.StatusBadRequest, shared.CommonResponse{
					Message: "user with this email already exists",
				})
			default:
				c.JSON(http.StatusInternalServerError, gin.H{
					"message": http.StatusText(http.StatusInternalServerError),
				})
				h.logger.Error(errors.Join(errors.New("got unexpected error code from database"), err).Error(),
					zap.String("method", c.Request.Method),
					zap.String("path", c.FullPath()))
			}
		case errors.Is(err, context.DeadlineExceeded):
			c.JSON(http.StatusRequestTimeout, gin.H{
				"message": "failed to create a user account in time",
			})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": http.StatusText(http.StatusInternalServerError),
			})

			h.logger.Error(err.Error(),
				zap.Stack("stackError"),
				zap.String("method", c.Request.Method),
				zap.String("path", c.FullPath()),
			)
		}
		return

	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "successfully create a user account",
	})
}


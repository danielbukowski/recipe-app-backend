package auth

import (
	"github.com/gin-gonic/gin"
)

func (h *handler) RegisterRoutes(gin *gin.Engine) {
	gin.POST("api/v1/auth/signup", h.SignUp)
}

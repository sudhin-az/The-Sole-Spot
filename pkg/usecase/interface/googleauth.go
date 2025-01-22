package interfaces

import (
	"ecommerce_clean_architecture/pkg/utils/models"

	"github.com/gin-gonic/gin"
)

type AuthUseCase interface {
	HandleGoogleCallback(c *gin.Context, code string) (models.TempUser, string, error)
}

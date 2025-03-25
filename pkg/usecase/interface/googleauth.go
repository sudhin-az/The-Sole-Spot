package interfaces

import (
	"ecommerce_clean_arch/pkg/utils/models"

	"github.com/gin-gonic/gin"
)

type AuthUseCase interface {
	HandleGoogleCallback(c *gin.Context, code string) (models.TempUser, string, error)
}

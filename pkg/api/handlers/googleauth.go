package handlers

import (
	"ecommerce_clean_arch/pkg/usecase"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authUseCase *usecase.AuthUseCase
}

func NewAuthHandler(authUseCase *usecase.AuthUseCase) *AuthHandler {
	return &AuthHandler{authUseCase: authUseCase}
}

// GoogleLogin godoc
// @Summary Initiates Google login
// @Description Redirects the user to Google for authentication
// @Tags Auth
// @Produce json
// @Success 302 {object} gin.H{"url": "string"} "Redirects to Google login URL"
// @Router /auth/google/login [get]
func (h *AuthHandler) GoogleLogin(c *gin.Context) {
	url := h.authUseCase.HandleGoogleLogin()
	c.Redirect(http.StatusTemporaryRedirect, url)
}

// GoogleCallback godoc
// @Summary Handles Google login callback
// @Description Receives the authorization code from Google and retrieves user information and token
// @Tags Auth
// @Accept json
// @Produce json
// @Param code query string true "Authorization code from Google"
// @Success 200 {object} gin.H{"user": "User ", "token": "string"} "User  information and token"
// @Failure 400 {object} gin.H{"error": "string"} "Missing 'code' query parameter"
// @Failure 500 {object} gin.H{"error": "string"} "Internal server error"
// @Router /auth/google/callback [get]
func (h *AuthHandler) GoogleCallback(c *gin.Context) {
	code := c.Query("code")
	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing 'code' query parameter"})
		return
	}

	user, token, err := h.authUseCase.HandleGoogleCallback(c, code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user":  user,
		"token": token,
	})
}

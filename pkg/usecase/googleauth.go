package usecase

import (
	"ecommerce_clean_architecture/pkg/domain"
	"ecommerce_clean_architecture/pkg/helper"
	"ecommerce_clean_architecture/pkg/repository"
	"ecommerce_clean_architecture/pkg/utils"
	"ecommerce_clean_architecture/pkg/utils/models"
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
	"gorm.io/gorm"
)

type AuthUseCase struct {
	userRepo    repository.UserRepository
	OAuthConfig *oauth2.Config
}

func NewAuthUseCase(userRepo repository.UserRepository, oauthConfig *oauth2.Config) *AuthUseCase {
	return &AuthUseCase{
		userRepo:    userRepo,
		OAuthConfig: oauthConfig,
	}
}

func (uc *AuthUseCase) HandleGoogleLogin() string {
	return uc.OAuthConfig.AuthCodeURL("state")
}

func (uc *AuthUseCase) HandleGoogleCallback(c *gin.Context, code string) (models.User, string, error) {

	token, err := uc.OAuthConfig.Exchange(oauth2.NoContext, code)
	if err != nil {
		return models.User{}, "", errors.New("failed to exchange token")
	}

	resp, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
	if err != nil {
		return models.User{}, "", errors.New("failed to get user information")
	}
	defer resp.Body.Close()

	var googleUser domain.GoogleResponse
	if err := utils.ParseJSON(resp.Body, &googleUser); err != nil {
		return models.User{}, "", errors.New("failed to parse user information")
	}

	user := models.User{
		Email: googleUser.Email,
	}

	existingUser, err := uc.userRepo.GetUserByEmail(user.Email)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			if err := uc.userRepo.CreateUser(user); err != nil {
				return models.User{}, "", errors.New("failed to create user through Google SSO")
			}
			// existingUser = user
		} else {
			return models.User{}, "", errors.New("failed to fetch user information")
		}
	}

	if existingUser.Blocked {
		return models.User{}, "", errors.New("user is unauthorized to access")
	}

	tokenString, err := helper.GenerateTokenUsers(uint(existingUser.ID), existingUser.Email, time.Now().Add(24*time.Hour))
	if err != nil {
		return models.User{}, "", errors.New("failed to create authorization token")
	}

	return existingUser, tokenString, nil
}

package handlers

import (
	"ecommerce_clean_architecture/pkg/usecase"
	"ecommerce_clean_architecture/pkg/utils/models"
	"ecommerce_clean_architecture/pkg/utils/response"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type UserHandler struct {
	userUseCase usecase.UserUseCase
}

func NewUserHandler(u usecase.UserUseCase) *UserHandler {
	return &UserHandler{
		userUseCase: u,
	}
}

// UserSignUp handles the user signup process.
func (h *UserHandler) UserSignUp(c *gin.Context) {
	var user models.UserSignUp

	if err := c.ShouldBindJSON(&user); err != nil {
		errRes := response.ClientResponse(http.StatusBadRequest, "Invalid request data", nil, err.Error())
		c.JSON(http.StatusBadRequest, errRes)
		return
	}

	// Check if email or phone already exists
	if h.userUseCase.IsEmailExists(user.Email) || h.userUseCase.IsPhoneExists(user.Phone) {
		errRes := response.ClientResponse(http.StatusConflict, "Email or Phone already exists", nil, "duplicate email or phone")
		c.JSON(http.StatusConflict, errRes)
		return
	}

	// Save user data temporarily and generate OTP
	tokenUsers, err := h.userUseCase.SaveTempUserAndGenerateOTP(user)
	if err != nil {
		errRes := response.ClientResponse(http.StatusInternalServerError, "Signup failed", nil, err.Error())
		c.JSON(http.StatusInternalServerError, errRes)
		return
	}

	successRes := response.ClientResponse(http.StatusOK, "OTP sent successfully", tokenUsers, nil)
	c.JSON(http.StatusOK, successRes)
}

// VerifyOTP handles OTP verification and user creation.
func (h *UserHandler) VerifyOTP(c *gin.Context) {
	email := c.Param("email")
	email = strings.Trim(email, "\"")

	fmt.Println("hellooooooooo", email)
	if email == "" {
		errRes := response.ClientResponse(http.StatusBadRequest, "Email is required", nil, "missing email parameter")
		c.JSON(http.StatusBadRequest, errRes)
		return
	}

	var verifyUser models.VerifyOTP
	if err := c.ShouldBindJSON(&verifyUser); err != nil {
		errRes := response.ClientResponse(http.StatusBadRequest, "Invalid request data", nil, err.Error())
		c.JSON(http.StatusBadRequest, errRes)
		return
	}

	// Pass the email and OTP for verification
	tokenUsers, err := h.userUseCase.VerifyOTPAndRegisterUser(email, verifyUser.OTP)
	if err != nil {
		errRes := response.ClientResponse(http.StatusUnauthorized, "OTP verification failed", nil, err.Error())
		c.JSON(http.StatusUnauthorized, errRes)
		return
	}

	// If OTP verification is successful, you can return the token or any other relevant response
	successRes := response.ClientResponse(http.StatusOK, "OTP verified and user registered successfully", tokenUsers, nil)
	c.JSON(http.StatusOK, successRes)

}

// ResendOTP handles the OTP resend functionality.
func (h *UserHandler) ResendOTP(c *gin.Context) {
	email := c.Param("email")

	if err := h.userUseCase.ResendOTP(email); err != nil {
		errRes := response.ClientResponse(http.StatusInternalServerError, "Failed to resend OTP", nil, err.Error())
		c.JSON(http.StatusInternalServerError, errRes)
		return
	}

	successRes := response.ClientResponse(http.StatusOK, "OTP resent successfully", nil, nil)
	c.JSON(http.StatusOK, successRes)
}

// UserLogin handles the user login process.
func (h *UserHandler) UserLogin(c *gin.Context) {
	var user models.UserLogin
	if err := c.ShouldBindJSON(&user); err != nil {
		errRes := response.ClientResponse(http.StatusBadRequest, "fields provided are in wrong format", nil, err.Error())
		c.JSON(http.StatusBadRequest, errRes)
		return
	}

	err := validator.New().Struct(user)
	if err != nil {
		errRes := response.ClientResponse(http.StatusBadRequest, "constraints not satisfied", nil, err.Error())
		c.JSON(http.StatusBadRequest, errRes)
		return
	}

	userDetails, err := h.userUseCase.UserLogin(user)
	if err != nil {
		errRes := response.ClientResponse(http.StatusInternalServerError, "User could not be logged in", nil, err.Error())
		c.JSON(http.StatusInternalServerError, errRes)
		return
	}

	successRes := response.ClientResponse(http.StatusCreated, "User successfully logged in", userDetails, nil)
	c.JSON(http.StatusCreated, successRes)
}

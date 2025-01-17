package handlers

import (
	"ecommerce_clean_architecture/pkg/domain"
	"ecommerce_clean_architecture/pkg/helper"
	"ecommerce_clean_architecture/pkg/usecase"
	"ecommerce_clean_architecture/pkg/utils/models"
	"ecommerce_clean_architecture/pkg/utils/response"
	"strconv"

	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userUseCase usecase.UserUseCase
}

func NewUserHandler(u usecase.UserUseCase) *UserHandler {
	return &UserHandler{
		userUseCase: u,
	}
}

func (h *UserHandler) UserSignup(c *gin.Context) {
	var user models.User

	if err := c.ShouldBindJSON(&user); err != nil {
		errRes := response.ClientResponse(http.StatusBadRequest, "Invalid request data", nil, err.Error())
		c.JSON(http.StatusBadRequest, errRes)
		return
	}
	if err := usecase.ValidateUser(user); err != nil {
		errRes := response.ClientResponse(http.StatusBadRequest, "Validation failed", nil, err.Error())
		c.JSON(http.StatusBadRequest, errRes)
		return
	}

	tokenUsers, err := h.userUseCase.SaveTempUserAndGenerateOTP(user)
	if err != nil {
		errRes := response.ClientResponse(http.StatusInternalServerError, "Signup failed", nil, err.Error())
		c.JSON(http.StatusInternalServerError, errRes)
		return
	}

	successRes := response.ClientResponse(http.StatusOK, "OTP sent successfully", tokenUsers, nil)
	c.JSON(http.StatusOK, successRes)
}

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

	tokenUsers, err := h.userUseCase.VerifyOTPAndRegisterUser(email, verifyUser.OTP)
	if err != nil {
		errRes := response.ClientResponse(http.StatusUnauthorized, "OTP verification failed", nil, err.Error())
		c.JSON(http.StatusUnauthorized, errRes)
		return
	}

	successRes := response.ClientResponse(http.StatusOK, "OTP verified and user registered successfully", tokenUsers, nil)
	c.JSON(http.StatusOK, successRes)

}

func (h *UserHandler) ResendOTP(c *gin.Context) {
	email := c.Param("email")
	fmt.Println("Email:", email)
	if err := h.userUseCase.ResendOTP(email); err != nil {
		errRes := response.ClientResponse(http.StatusInternalServerError, "Failed to resend OTP", nil, err.Error())
		c.JSON(http.StatusInternalServerError, errRes)
		return
	}

	successRes := response.ClientResponse(http.StatusOK, "OTP resent successfully", nil, nil)
	c.JSON(http.StatusOK, successRes)
}

func (h *UserHandler) UserLogin(c *gin.Context) {
	var loginReq models.UserLogin

	// Bind JSON to LoginRequest struct
	if err := c.ShouldBindJSON(&loginReq); err != nil {
		errRes := response.ClientResponse(http.StatusBadRequest, "Fields provided are in the wrong format", nil, err.Error())
		c.JSON(http.StatusBadRequest, errRes)
		return
	}

	// Validate login fields
	if message, err := helper.ValidateAddress(loginReq); err != nil {
		errRes := response.ClientResponse(http.StatusBadRequest, "Validation failed", message, err.Error())
		c.JSON(http.StatusBadRequest, errRes)
		return
	}

	userDetails, err := h.userUseCase.UserLogin(loginReq)
	if err != nil {
		errRes := response.ClientResponse(http.StatusInternalServerError, "User could not be logged in", nil, err.Error())
		c.JSON(http.StatusInternalServerError, errRes)
		return
	}

	successRes := response.ClientResponse(http.StatusCreated, "User successfully logged in", userDetails, nil)
	c.JSON(http.StatusCreated, successRes)
}

func (h *UserHandler) GetProducts(c *gin.Context) {
	products, err := h.userUseCase.GetProducts()
	if err != nil {
		errorRes := response.ClientResponse(http.StatusInternalServerError, "could not retrieve records of products", nil, err.Error())
		c.JSON(http.StatusInternalServerError, errorRes)
		return
	}
	successRes := response.ClientResponse(http.StatusOK, "Successfully retrieved the products", products, nil)
	c.JSON(http.StatusOK, successRes)
}
func (cat *UserHandler) ListCategory(c *gin.Context) {
	category, err := cat.userUseCase.ListCategory()
	if err != nil {
		errorRes := response.ClientResponse(http.StatusInternalServerError, "could not retrieve records of categories", nil, err.Error())
		c.JSON(http.StatusInternalServerError, errorRes)
		return
	}
	successRes := response.ClientResponse(http.StatusOK, "Successfully retrieved the categories", category, nil)
	c.JSON(http.StatusOK, successRes)
}

func (u *UserHandler) UserProfile(c *gin.Context) {
	userID, ok := c.Get("id")
	if !ok {
		errRes := response.ClientResponse(http.StatusUnauthorized, "User ID not found in context", nil, nil)
		c.JSON(http.StatusUnauthorized, errRes)
		return
	}

	idInt, ok := userID.(int)
	if !ok {
		errRes := response.ClientResponse(http.StatusBadRequest, "Invalid user ID format", nil, nil)
		c.JSON(http.StatusBadRequest, errRes)
		return
	}

	idStr := idInt
	userprofile, err := u.userUseCase.UserProfile(idStr)
	if err != nil {
		errRes := response.ClientResponse(http.StatusNotFound, "Failed to retrieve user profile details", nil, err.Error())
		c.JSON(http.StatusNotFound, errRes)
		return
	}

	successRes := response.ClientResponse(http.StatusOK, "User profile details retrieved successfully", userprofile, nil)
	c.JSON(http.StatusOK, successRes)
}
func (u *UserHandler) UpdateProfile(c *gin.Context) {
	var profile models.User

	userID, ok := c.Get("id")
	if !ok {
		errRes := response.ClientResponse(http.StatusUnauthorized, "User ID not found in context", nil, nil)
		c.JSON(http.StatusUnauthorized, errRes)
		return
	}
	profile.ID = userID.(int)

	err := c.ShouldBindJSON(&profile)
	if err != nil {
		errRes := response.ClientResponse(http.StatusNotFound, "fields provided are in wrong format", nil, err.Error())
		c.JSON(http.StatusNotFound, errRes)
		return
	}
	if err := usecase.ValidateUser(profile); err != nil {
		errRes := response.ClientResponse(http.StatusBadRequest, "Validation failed", nil, err.Error())
		c.JSON(http.StatusBadRequest, errRes)
		return
	}

	userProfile, err := u.userUseCase.UpdateProfile(profile)
	if err != nil {
		errRes := response.ClientResponse(http.StatusNotFound, "Failed to update user profile details", nil, err.Error())
		c.JSON(http.StatusNotFound, errRes)
		return
	}
	successRes := response.ClientResponse(http.StatusOK, "User profile details retrieved successfully", userProfile, nil)
	c.JSON(http.StatusOK, successRes)
}

func (h *UserHandler) ForgotPassword(c *gin.Context) {
	userID, ok := c.Get("id")
	if !ok {
		errRes := response.ClientResponse(http.StatusUnauthorized, "User ID not found in context", nil, nil)
		c.JSON(http.StatusUnauthorized, errRes)
		return
	}
	ID := userID.(int)
	var input models.NewPassword
	err := c.ShouldBindJSON(&input)
	if err != nil {
		errRes := response.ClientResponse(http.StatusNotFound, "fields provided are in wrong format", nil, err.Error())
		c.JSON(http.StatusNotFound, errRes)
		return
	}
	if err := usecase.ValidatePassword(input); err != nil {
		errRes := response.ClientResponse(http.StatusBadRequest, "Validation failed", nil, err.Error())
		c.JSON(http.StatusBadRequest, errRes)
		return
	}
	userPassword, err := h.userUseCase.ForgotPassword(ID, input)
	if err != nil {
		errRes := response.ClientResponse(http.StatusNotFound, "Failed to forgot password", nil, err.Error())
		c.JSON(http.StatusNotFound, errRes)
		return
	}
	successRes := response.ClientResponse(http.StatusOK, "Password changed successfully", userPassword, nil)
	c.JSON(http.StatusOK, successRes)
}
func (u *UserHandler) AddAddress(c *gin.Context) {
	userID, ok := c.Get("id")
	if !ok {
		errRes := response.ClientResponse(http.StatusUnauthorized, "User ID not found in context", nil, nil)
		c.JSON(http.StatusUnauthorized, errRes)
		return
	}
	ID := userID.(int)

	var address models.AddAddress
	if err := c.BindJSON(&address); err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "Fields provided are in wrong format", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}

	if message, err := helper.ValidateAddress(address); err != nil {
		errRes := response.ClientResponse(http.StatusBadRequest, "Validation failed", message, err.Error())
		c.JSON(http.StatusBadRequest, errRes)
		return
	}

	address, err := u.userUseCase.AddAddress(ID, address)
	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "Could not add the address", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}

	successRes := response.ClientResponse(http.StatusOK, "Successfully added address", address, nil)
	c.JSON(http.StatusOK, successRes)
}
func (u *UserHandler) UpdateAddress(c *gin.Context) {
	userID, ok := c.Get("id")
	if !ok {
		errRes := response.ClientResponse(http.StatusUnauthorized, "User ID not found in context", nil, nil)
		c.JSON(http.StatusUnauthorized, errRes)
		return
	}
	ID := userID.(int)

	addressIDStr := c.Query("id")
	if addressIDStr == "" {
		errRes := response.ClientResponse(http.StatusBadRequest, "Address ID is required", nil, nil)
		c.JSON(http.StatusBadRequest, errRes)
		return
	}
	addressID, err := strconv.Atoi(addressIDStr)
	if err != nil {
		errRes := response.ClientResponse(http.StatusBadRequest, "Invalid Address ID", nil, err.Error())
		c.JSON(http.StatusBadRequest, errRes)
		return
	}

	var address domain.Address
	if err := c.BindJSON(&address); err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "Fields provided are in wrong format", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}

	address.ID = addressID

	if message, err := helper.ValidateAddress(address); err != nil {
		errRes := response.ClientResponse(http.StatusBadRequest, "Validation failed", message, err.Error())
		c.JSON(http.StatusBadRequest, errRes)
		return
	}

	address, err = u.userUseCase.UpdateAddress(ID, address)
	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "Could not edit the address", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}

	successRes := response.ClientResponse(http.StatusOK, "Successfully updated address", address, nil)
	c.JSON(http.StatusOK, successRes)
}

func (u *UserHandler) DeleteAddress(c *gin.Context) {
	idParam := c.Query("id")
	if idParam == "" {
		errRes := response.ClientResponse(http.StatusBadRequest, "id parameter is required", nil, "missing id")
		c.JSON(http.StatusBadRequest, errRes)
		return
	}
	userID, err := strconv.Atoi(idParam)
	if err != nil {
		errRes := response.ClientResponse(http.StatusBadRequest, "the id provided is invalid", nil, err.Error())
		c.JSON(http.StatusBadRequest, errRes)
		return
	}
	err = u.userUseCase.DeleteAddress(userID)
	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "Could not delete the address", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}
	successRes := response.ClientResponse(http.StatusOK, "Successfully deleted address", nil, nil)
	c.JSON(http.StatusOK, successRes)
}

func (u *UserHandler) GetAllAddresses(c *gin.Context) {
	userID, ok := c.Get("id")
	if !ok {
		errRes := response.ClientResponse(http.StatusUnauthorized, "User ID not found in context", nil, nil)
		c.JSON(http.StatusUnauthorized, errRes)
		return
	}
	ID := userID.(int)
	address, err := u.userUseCase.GetAllAddresses(ID)
	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "could not open checkout", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}
	successRes := response.ClientResponse(http.StatusOK, "Successfully got all records", address, nil)
	c.JSON(http.StatusOK, successRes)
}

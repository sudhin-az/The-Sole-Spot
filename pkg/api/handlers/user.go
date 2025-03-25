package handlers

import (
	"ecommerce_clean_arch/pkg/helper"
	"ecommerce_clean_arch/pkg/usecase"
	"ecommerce_clean_arch/pkg/utils/models"
	"ecommerce_clean_arch/pkg/utils/response"
	"log"
	"strconv"

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

// UserSignup godoc
// @Summary User signup
// @Description Registers a new user and sends an OTP for verification
// @Tags Users
// @Accept json
// @Produce json
// @Param user body models.User true "User  details"
// @Success 200 {object} response.ClientResponse
// @Failure 400 {object} response.ClientResponse
// @Failure 500 {object} response.ClientResponse
// @Router /signup [post]
func (h *UserHandler) UserSignup(c *gin.Context) {
	var user models.User

	if err := c.ShouldBindJSON(&user); err != nil {
		errRes := response.UserResponse("Invalid request data")
		c.JSON(http.StatusBadRequest, errRes)
		return
	}
	if err := usecase.ValidateUser(user); err != nil {
		errRes := response.UserResponse("Validation failed")
		c.JSON(http.StatusBadRequest, errRes)
		return
	}

	_, err := h.userUseCase.SaveTempUserAndGenerateOTP(user)
	if err != nil {
		errRes := response.UserResponse("Signup failed")
		c.JSON(http.StatusInternalServerError, errRes)
		return
	}

	successRes := response.UserResponse("OTP sent successfully")
	c.JSON(http.StatusOK, successRes)
}

// VerifyOTP godoc
// @Summary Verify OTP
// @Description Verifies the OTP sent to the user's email and registers the user
// @Tags Users
// @Param email path string true "User  email"
// @Param verifyUser  body models.VerifyOTP true "OTP details"
// @Produce json
// @Success 200 {object} response.ClientResponse
// @Failure 400 {object} response.ClientResponse
// @Failure 401 {object} response.ClientResponse
// @Router /verify/{email} [post]
func (h *UserHandler) VerifyOTP(c *gin.Context) {

	email := c.Param("email")
	email = strings.Trim(email, "\"")

	log.Println("hellooooooooo", email)
	if email == "" {
		errRes := response.UserResponse("Email is required")
		c.JSON(http.StatusBadRequest, errRes)
		return
	}

	var verifyUser models.VerifyOTP
	if err := c.ShouldBindJSON(&verifyUser); err != nil {
		errRes := response.UserResponse("Invalid request data")
		c.JSON(http.StatusBadRequest, errRes)
		return
	}

	_, err := h.userUseCase.VerifyOTPAndRegisterUser(email, verifyUser.OTP)
	if err != nil {
		errRes := response.UserResponse("OTP verification failed")
		c.JSON(http.StatusUnauthorized, errRes)
		return
	}

	successRes := response.UserResponse("OTP verified and user registered successfully")
	c.JSON(http.StatusOK, successRes)

}

// ResendOTP godoc
// @Summary Resend OTP
// @Description Resends the OTP to the user's email
// @Tags Users
// @Param email path string true "User  email"
// @Produce json
// @Success 200 {object} response.ClientResponse
// @Failure 500 {object} response.ClientResponse
// @Router /resend-otp/{email} [post]
func (h *UserHandler) ResendOTP(c *gin.Context) {
	email := c.Param("email")
	log.Println("Email:", email)
	if err := h.userUseCase.ResendOTP(email); err != nil {
		errRes := response.ClientResponse(http.StatusInternalServerError, "Failed to resend OTP", nil, err.Error())
		c.JSON(http.StatusInternalServerError, errRes)
		return
	}

	successRes := response.ClientResponse(http.StatusOK, "OTP resent successfully", nil, nil)
	c.JSON(http.StatusOK, successRes)
}

// UserLogin godoc
// @Summary User login
// @Description Logs in a user and returns user details
// @Tags Users
// @Accept json
// @Produce json
// @Param loginReq body models.UserLogin true "Login details"
// @Success 201 {object} response.ClientResponse
// @Failure 400 {object} response.ClientResponse
// @Failure 500 {object} response.ClientResponse
// @Router /login [post]
func (h *UserHandler) UserLogin(c *gin.Context) {
	var loginReq models.UserLogin

	if err := c.ShouldBindJSON(&loginReq); err != nil {
		errRes := response.ClientResponse(http.StatusBadRequest, "Fields provided are in the wrong format", nil, err.Error())
		c.JSON(http.StatusBadRequest, errRes)
		return
	}

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

// GetProducts godoc
// @Summary Get all products
// @Description Retrieves a list of all products
// @Tags Users
// @Produce json
// @Success 200 {object} response.ClientResponse
// @Failure 500 {object} response.ClientResponse
// @Router /products [get]
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

// ListCategory godoc
// @Summary Get all categories
// @Description Retrieves a list of all product categories
// @Tags Users
// @Produce json
// @Success 200 {object} response.ClientResponse
// @Failure 500 {object} response.ClientResponse
// @Router /categories [get]
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

// UserProfile godoc
// @Summary Get user profile
// @Description Retrieves the profile details of the authenticated user
// @Tags Users
// @Produce json
// @Success 200 {object} response.ClientResponse
// @Failure 401 {object} response.ClientResponse
// @Failure 400 {object} response.ClientResponse
// @Failure 404 {object} response.ClientResponse
// @Router /profile [get]
func (u *UserHandler) UserProfile(c *gin.Context) {
	userID, ok := c.Get("id")
	if !ok {
		errRes := response.ClientResponse(http.StatusUnauthorized, "User ID not found in context", nil, nil)
		c.JSON(http.StatusUnauthorized, errRes)
		return
	}

	idInt, ok := userID.(int)
	if !ok {
		errRes := response.ClientResponse(http.StatusBadRequest, "Invalid user ID format in context", nil, nil)
		c.JSON(http.StatusBadRequest, errRes)
		return
	}
	userProfile, err := u.userUseCase.UserProfile(idInt)
	if err != nil {
		errRes := response.ClientResponse(http.StatusNotFound, "Failed to retrieve user profile details", nil, err.Error())
		c.JSON(http.StatusNotFound, errRes)
		return
	}

	successRes := response.ClientResponse(http.StatusOK, "User profile details retrieved successfully", userProfile, nil)
	c.JSON(http.StatusOK, successRes)
}

// UpdateProfile godoc
// @Summary Update user profile
// @Description Updates the profile details of the authenticated user
// @Tags Users
// @Accept json
// @Produce json
// @Param profile body models.User true "User  profile details"
// @Success 200 {object} response.ClientResponse
// @Failure 401 {object} response.ClientResponse
// @Failure 400 {object} response.ClientResponse
// @Failure 404 {object} response.ClientResponse
// @Router /profile [put]
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
		errRes := response.ClientResponse(http.StatusBadRequest, "fields provided are in wrong format", nil, err.Error())
		c.JSON(http.StatusBadRequest, errRes)
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

	successRes := response.ClientResponse(http.StatusOK, "User profile details updated successfully", userProfile, nil)
	c.JSON(http.StatusOK, successRes)
}

// SendOTP godoc
// @Summary Send OTP
// @Description Sends an OTP to the user's email for verification
// @Tags Users
// @Accept json
// @Produce json
// @Param sendOTP body models.SendOTP true "Email to send OTP"
// @Success 200 {object} response.ClientResponse
// @Failure 400 {object} response.ClientResponse
// @Failure 500 {object} response.ClientResponse
// @Router /send-otp [post]
func (h *UserHandler) SendOTP(c *gin.Context) {
	var sendOTP models.SendOTP
	if err := c.ShouldBindJSON(&sendOTP); err != nil {
		c.JSON(http.StatusBadRequest, response.ClientResponse(http.StatusBadRequest, "Invalid input format", nil, err.Error()))
		return
	}

	token, err := h.userUseCase.GenerateAndSendOTP(sendOTP.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.ClientResponse(http.StatusInternalServerError, "Failed to send OTP", nil, err.Error()))
		return
	}
	responseData := map[string]interface{}{
		"token": token,
	}
	c.JSON(http.StatusOK, response.ClientResponse(http.StatusOK, "OTP sent successfully", responseData, nil))
}

// ForgotPassword godoc
// @Summary Forgot Password
// @Description Resets the user's password using the provided OTP and new password
// @Tags Users
// @Accept json
// @Produce json
// @Param input body models.ForgotPassword true "Forgot password details"
// @Success 200 {object} response.ClientResponse
// @Failure 400 {object} response.ClientResponse
// @Failure 401 {object} response.ClientResponse
// @Failure 500 {object} response.ClientResponse
// @Router /forgot-password [post]
func (h *UserHandler) ForgotPassword(c *gin.Context) {
	var input models.ForgotPassword
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, response.ClientResponse(http.StatusBadRequest, "Invalid input format", nil, err.Error()))
		return
	}

	tokenString := c.GetHeader("Authorization")

	email, err := helper.VerifyTemporaryToken(tokenString)
	if err != nil {
		c.JSON(http.StatusUnauthorized, response.ClientResponse(http.StatusUnauthorized, "Invalid or expired token", nil, err.Error()))
		return
	}

	if input.Password != input.ConfirmPassword {
		c.JSON(http.StatusBadRequest, response.ClientResponse(http.StatusBadRequest, "Passwords do not match", nil, nil))
		return
	}

	user, err := h.userUseCase.ResetPassword(email, input.Otp, input.Password, input.ConfirmPassword)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.ClientResponse(http.StatusInternalServerError, "Failed to reset password", nil, err.Error()))
		return
	}
	c.JSON(http.StatusOK, response.ClientResponse(http.StatusOK, "Password reset successfully", user, nil))
}

// ChangePassword godoc
// @Summary Change Password
// @Description Changes the user's password
// @Tags Users
// @Accept json
// @Produce json
// @Param input body models.NewPassword true "New password details"
// @Success 200 {object} response.ClientResponse
// @Failure 401 {object} response.ClientResponse
// @Failure 400 {object} response.ClientResponse
// @Failure 404 {object} response.ClientResponse
// @Router /change-password [post]
func (h *UserHandler) ChangePassword(c *gin.Context) {
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
	userPassword, err := h.userUseCase.ChangePassword(ID, input)
	if err != nil {
		errRes := response.ClientResponse(http.StatusNotFound, "Failed to change password", nil, err.Error())
		c.JSON(http.StatusNotFound, errRes)
		return
	}
	successRes := response.ClientResponse(http.StatusOK, "Password changed successfully", userPassword, nil)
	c.JSON(http.StatusOK, successRes)
}

// AddAddress godoc
// @Summary Add a new address
// @Description Adds a new address for the authenticated user
// @Tags Users
// @Accept json
// @Produce json
// @Param address body models.AddAddress true "Address details"
// @Success 200 {object} response.ClientResponse
// @Failure 401 {object} response.ClientResponse
// @Failure 400 {object} response.ClientResponse
// @Router /addresses [post]
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

// UpdateAddress godoc
// @Summary Update an existing address
// @Description Updates the details of an existing address for the authenticated user
// @Tags Users
// @Accept json
// @Produce json
// @Param id query int true "Address ID"
// @Param address body models.AddAddress true "Updated address details"
// @Success 200 {object} response.ClientResponse
// @Failure 401 {object} response.ClientResponse
// @Failure 400 {object} response.ClientResponse
// @Router /addresses [put]
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

	var address models.AddAddress
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

// DeleteAddress godoc
// @Summary Delete an address
// @Description Deletes an existing address for the authenticated user
// @Tags Users
// @Param id query int true "Address ID"
// @Success 200 {object} response.ClientResponse
// @Failure 400 {object} response.ClientResponse
// @Router /addresses [delete]
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

// GetAllAddresses godoc
// @Summary Get all addresses
// @Description Retrieves all addresses for the authenticated user
// @Tags Users
// @Produce json
// @Success 200 {object} response.ClientResponse
// @Failure 401 {object} response.ClientResponse
// @Failure 400 {object} response.ClientResponse
// @Router /addresses [get]
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

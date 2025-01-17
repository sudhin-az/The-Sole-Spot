package handlers

import (
	"ecommerce_clean_architecture/pkg/usecase"
	"ecommerce_clean_architecture/pkg/utils/models"
	"ecommerce_clean_architecture/pkg/utils/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AdminHandler struct {
	adminUseCase usecase.AdminUseCase
}

func NewAdminHandler(usecase usecase.AdminUseCase) *AdminHandler {
	return &AdminHandler{
		adminUseCase: usecase,
	}
}

func (ad *AdminHandler) SignUpHandler(c *gin.Context) {
	var admin models.AdminSignUp

	if err := c.ShouldBindJSON(&admin); err != nil {
		errRes := response.ClientResponse(http.StatusBadRequest, "fields provided are wrong", nil, err.Error())
		c.JSON(http.StatusBadRequest, errRes)
		return
	}

	adminDetails, err := ad.adminUseCase.SignUpHandler(admin)
	if err != nil {
		errRes := response.ClientResponse(http.StatusInternalServerError, "cannot authenticate Admin", nil, err.Error())
		c.JSON(http.StatusInternalServerError, errRes)
		return
	}

	successRes := response.ClientResponse(http.StatusCreated, "Successfully signed up the user", adminDetails, nil)
	c.JSON(http.StatusCreated, successRes)
}

func (ad *AdminHandler) LoginHandler(c *gin.Context) {
	var admin models.AdminLogin

	if err := c.ShouldBindJSON(&admin); err != nil {
		errRes := response.ClientResponse(http.StatusBadRequest, "details not in correct format", nil, err.Error())
		c.JSON(http.StatusBadRequest, errRes)
		return
	}

	adminDetails, err := ad.adminUseCase.LoginHandler(admin)
	if err != nil {
		errRes := response.ClientResponse(http.StatusInternalServerError, "cannot authenticate Admin", nil, err.Error())
		c.JSON(http.StatusInternalServerError, errRes)
		return
	}

	successRes := response.ClientResponse(http.StatusCreated, "Successfully login the user", adminDetails, nil)
	c.JSON(http.StatusCreated, successRes)
}

func (ad *AdminHandler) GetUsers(c *gin.Context) {
	var listusers models.UserSignUp

	if err := c.ShouldBindJSON(&listusers); err != nil {
		errRes := response.ClientResponse(http.StatusBadRequest, "Constraints are not in correct format", nil, err.Error())
		c.JSON(http.StatusBadRequest, errRes)
		return
	}

	users, err := ad.adminUseCase.GetUsers(listusers)
	if err != nil {
		errorRes := response.ClientResponse(http.StatusInternalServerError, "could not retrieve records of users", nil, err.Error())
		c.JSON(http.StatusInternalServerError, errorRes)
		return
	}
	successRes := response.ClientResponse(http.StatusOK, "Successfully retrieved the users", users, nil)
	c.JSON(http.StatusOK, successRes)
}

package handlers

import (
	"ecommerce_clean_architecture/pkg/helper"
	"ecommerce_clean_architecture/pkg/usecase"
	"ecommerce_clean_architecture/pkg/utils/models"
	"ecommerce_clean_architecture/pkg/utils/response"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type CouponHandler struct {
	usecase usecase.CouponUseCase
}

func NewCouponHandler(usecase usecase.CouponUseCase) *CouponHandler {
	return &CouponHandler{
		usecase: usecase,
	}
}

var validate = validator.New()

func (coup *CouponHandler) CreateNewCoupon(c *gin.Context) {
	var newCoupon models.Coupon
	if err := c.ShouldBindJSON(&newCoupon); err != nil {
		errRes := response.ClientResponse(http.StatusBadRequest, "fields provided are in wrong format", nil, err.Error())
		c.JSON(http.StatusBadRequest, errRes)
		return
	}

	if newCoupon.ExpireDateStr != "" {
		parsedData, err := time.Parse("2006-01-02", newCoupon.ExpireDateStr)
		if err != nil {
			errRes := response.ClientResponse(http.StatusBadRequest, "Invalid date format, use YYYY-MM-DD", nil, err.Error())
			c.JSON(http.StatusBadRequest, errRes)
			return
		}
		newCoupon.ExpireDate = parsedData
	}
	if err := validate.Struct(newCoupon); err != nil {
		errorMessages := helper.ValidationErrorToText(err)
		errRes := response.ClientResponse(http.StatusBadRequest, strings.Join(errorMessages, ", "), nil, nil)
		c.JSON(http.StatusBadRequest, errRes)
		return
	}
	coupons, err := coup.usecase.AddCoupon(newCoupon)
	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "Could not add the Coupon", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}
	successRes := response.ClientResponse(http.StatusOK, "Successfully added Coupon", coupons, nil)
	c.JSON(http.StatusOK, successRes)
}

func (coup *CouponHandler) MakeCouponInvalid(c *gin.Context) {
	id, err := strconv.Atoi(c.Query("id"))
	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "fields provided are in wrong format", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}

	if err := coup.usecase.MakeCouponInvalid(id); err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "Coupon cannot be made invalid", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}

	successRes := response.ClientResponse(http.StatusOK, "Successfully made Coupon as invalid", nil, nil)
	c.JSON(http.StatusOK, successRes)
}

func (coup *CouponHandler) GetAllCoupons(c *gin.Context) {
	coupons, err := coup.usecase.GetAllCoupons()
	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "error in getting coupons", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}

	successRes := response.ClientResponse(http.StatusOK, "Successfully got all coupons", coupons, nil)
	c.JSON(http.StatusOK, successRes)

}

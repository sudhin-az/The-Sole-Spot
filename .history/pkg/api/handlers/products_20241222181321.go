package handlers

import (
	"ecommerce_clean_architecture/pkg/usecase"
	"ecommerce_clean_architecture/pkg/utils/models"
	"ecommerce_clean_architecture/pkg/utils/response"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ProductHandler struct {
	ProductUseCase usecase.ProductUseCase
}

func NewProductHandler(usecase usecase.ProductUseCase) *ProductHandler {
	return &ProductHandler{
		ProductUseCase: usecase,
	}
}

func (p *ProductHandler) AddProduct(c *gin.Context) {
	var addproduct models.AddProduct

	if err := c.ShouldBindJSON(&addproduct); err != nil {
		errRes := response.ClientResponse(http.StatusBadRequest, "the constraints are given wrong", nil, err.Error())
		c.JSON(http.StatusBadRequest, errRes)
		return
	}

	products, err := p.ProductUseCase.AddProduct(addproduct)
	if err != nil {
		errRes := response.ClientResponse(http.StatusBadRequest, "the products cannot be added", nil, err.Error())
		c.JSON(http.StatusBadRequest, errRes)
		return
	}

	successRes := response.ClientResponse(http.StatusOK, "the product added successfully", products, nil)
	c.JSON(http.StatusOK, successRes)

}

func (p *ProductHandler) UpdateProduct(c *gin.Context) {
	productID, err := strconv.Atoi(c.Query("id"))
	if err != nil {
		errRes := response.ClientResponse(http.StatusBadRequest, "check the parameter", nil, err.Error())
		c.JSON(http.StatusBadRequest, errRes)
		return
	}

	var products models.ProductResponse
	if err := c.ShouldBindJSON(&products); err != nil {
		errRes := response.ClientResponse(http.StatusBadRequest, "the constraints are given wrong", nil, err.Error())
		c.JSON(http.StatusBadRequest, errRes)
		return
	}

	updateProduct, err := p.ProductUseCase.UpdateProduct(products, productID)
	if err != nil {
		errRes := response.ClientResponse(http.StatusBadRequest, "the product cannot be updated", nil, err.Error())
		c.JSON(http.StatusBadRequest, errRes)
		return
	}
	successRes := response.ClientResponse(http.StatusOK, "the product is updated", updateProduct, nil)
	c.JSON(http.StatusOK, successRes)
}

func (p *ProductHandler) DeleteProduct(c *gin.Context) {
	productID, err := strconv.Atoi(c.Query("id"))
	if err != nil {
		errRes := response.ClientResponse(http.StatusBadRequest, "check the parameter", nil, err.Error())
		c.JSON(http.StatusBadRequest, errRes)
		return
	}
	err = p.ProductUseCase.DeleteProduct(productID)
	if err != nil {
		errRes := response.ClientResponse(http.StatusBadRequest, "the product cannot be deleted", nil, err.Error())
		c.JSON(http.StatusBadRequest, errRes)
		return
	}
	successRes := response.ClientResponse(http.StatusOK, "the product is deleted", nil, nil)
	c.JSON(http.StatusOK, successRes)
}

package handlers

import (
	"ecommerce_clean_arch/pkg/usecase"
	"ecommerce_clean_arch/pkg/utils/models"
	"ecommerce_clean_arch/pkg/utils/response"
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

// AddProduct godoc
// @Summary Add a new product
// @Description Adds a new product to the inventory
// @Tags Products
// @Accept json
// @Produce json
// @Param product body models.AddProduct true "Product details"
// @Success 200 {object} response.ClientResponse
// @Failure 400 {object} response.ClientResponse
// @Router /products [post]
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

// UpdateProduct godoc
// @Summary Update an existing product
// @Description Updates the details of an existing product by its ID
// @Tags Products
// @Accept json
// @Produce json
// @Param id query int true "Product ID"
// @Param product body models.ProductResponse true "Updated product details"
// @Success 200 {object} response.ClientResponse
// @Failure 400 {object} response.ClientResponse
// @Router /products [put]
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

// DeleteProduct godoc
// @Summary Delete a product
// @Description Deletes a product from the inventory by its ID
// @Tags Products
// @Param id query int true "Product ID"
// @Produce json
// @Success 200 {object} response.ClientResponse
// @Failure 400 {object} response.ClientResponse
// @Router /products [delete]
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

// SearchProduct godoc
// @Summary Search for products
// @Description Searches for products based on category and sorting criteria
// @Tags Products
// @Param category_id query string false "Category ID"
// @Param sort_by query string false "Sort by criteria"
// @Produce json
// @Success 200 {object} response.ClientResponse
// @Failure 500 {object} response.ClientResponse
// @Router /products/search [get]
func (p *ProductHandler) SearchProduct(c *gin.Context) {

	categoryID := c.Query("category_id")
	sortBy := c.Query("sort_by")

	products, err := p.ProductUseCase.SearchProduct(categoryID, sortBy)
	if err != nil {
		errRes := response.ClientResponse(http.StatusInternalServerError, "Failed to fetch products", nil, err.Error())
		c.JSON(http.StatusInternalServerError, errRes)
		return
	}

	successRes := response.ClientResponse(http.StatusOK, "Products fetched successfully", products, nil)
	c.JSON(http.StatusOK, successRes)
}

package handlers

import (
	"ecommerce_clean_arch/pkg/domain"
	"ecommerce_clean_arch/pkg/usecase"
	"ecommerce_clean_arch/pkg/utils/response"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CategoryHandler struct {
	CategoryUseCase usecase.CategoryUseCase
}

func NewCategoryHandler(usecase usecase.CategoryUseCase) *CategoryHandler {
	return &CategoryHandler{
		CategoryUseCase: usecase,
	}
}

// @Summary Add a new category
// @Description Creates a new category
// @Tags Category
// @Accept json
// @Produce json
// @Param category body domain.Category true "Category details"
// @Success 200 {object} response.Response{}
// @Failure 400 {object} response.Response{}
// @Router /category [post]
func (cat *CategoryHandler) AddCategory(c *gin.Context) {
	var category domain.Category

	if err := c.ShouldBindJSON(&category); err != nil {
		errRes := response.ClientResponse(http.StatusBadRequest, "the parameters are wrong", nil, err.Error())
		c.JSON(http.StatusBadRequest, errRes)
		return
	}

	CategoryResponse, err := cat.CategoryUseCase.AddCategory(category)
	if err != nil {
		errRes := response.ClientResponse(http.StatusBadRequest, "the category cannot be added", nil, err.Error())
		c.JSON(http.StatusBadRequest, errRes)
		return
	}
	successRes := response.ClientResponse(http.StatusOK, "the category is added successfully", CategoryResponse, nil)
	c.JSON(http.StatusOK, successRes)
}

// @Summary Update an existing category
// @Description Updates an existing category
// @Tags Category
// @Accept json
// @Produce json
// @Param id query int true "Category ID to update"
// @Param category body domain.Category true "Updated category details"
// @Success 200 {object} response.Response{}
// @Failure 400 {object} response.Response{}
// @Router /category [put]
func (cat *CategoryHandler) UpdateCategory(c *gin.Context) {
	categoryID, err := strconv.Atoi(c.Query("id"))
	if err != nil {
		errRes := response.ClientResponse(http.StatusBadRequest, "check the parameter", nil, err.Error())
		c.JSON(http.StatusBadRequest, errRes)
		return
	}

	var category domain.Category
	if err := c.ShouldBindJSON(&category); err != nil {
		errRes := response.ClientResponse(http.StatusBadRequest, "the constraints are given wrong", nil, err.Error())
		c.JSON(http.StatusBadRequest, errRes)
		return
	}

	updateCategory, err := cat.CategoryUseCase.UpdateCategory(category, categoryID)
	if err != nil {
		errRes := response.ClientResponse(http.StatusBadRequest, "the category cannot be updated", nil, err.Error())
		c.JSON(http.StatusBadRequest, errRes)
		return
	}

	successRes := response.ClientResponse(http.StatusOK, "the category is updated", updateCategory, nil)
	c.JSON(http.StatusOK, successRes)
}

// @Summary Delete a category
// @Description Deletes a category by ID
// @Tags Category
// @Produce json
// @Param id query int true "Category ID to delete"
// @Success 200 {object} response.Response{}
// @Failure 400 {object} response.Response{}
// @Router /category [delete]
func (cat *CategoryHandler) DeleteCategory(c *gin.Context) {
	categoryID, err := strconv.Atoi(c.Query("id"))
	if err != nil {
		errRes := response.ClientResponse(http.StatusBadRequest, "check the parameter", nil, err.Error())
		c.JSON(http.StatusBadRequest, errRes)
		return
	}

	err = cat.CategoryUseCase.DeleteCategory(categoryID)
	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "fields provided are in wrong format", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}

	successRes := response.ClientResponse(http.StatusOK, "Successfully deleted the Category", nil, nil)
	c.JSON(http.StatusOK, successRes)
}

package handlers

import (
	"ecommerce_clean_architecture/pkg/usecase"
	"ecommerce_clean_architecture/pkg/utils/models"
	"ecommerce_clean_architecture/pkg/utils/response"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type WishlistHandler struct {
	wishlistUseCase usecase.WishlistUseCase
}

func NewWishlistHandler(usecase usecase.WishlistUseCase) *WishlistHandler {
	return &WishlistHandler{
		wishlistUseCase: usecase,
	}
}

// AddToWishList godoc
// @Summary Add a product to the wishlist
// @Description Adds a specified product to the authenticated user's wishlist
// @Tags Wishlist
// @Accept json
// @Produce json
// @Param wishlist body models.WishlistRequest true "Wishlist request details"
// @Success 200 {object} response.ClientResponse
// @Failure 401 {object} response.ClientResponse
// @Failure 400 {object} response.ClientResponse
// @Router /wishlist [post]

func (w *WishlistHandler) AddToWishList(c *gin.Context) {
	userID, ok := c.Get("id")
	if !ok {
		errRes := response.ClientResponse(http.StatusUnauthorized, "User ID not found in context", nil, nil)
		c.JSON(http.StatusUnauthorized, errRes)
		return
	}
	UserID := userID.(int)
	var wishlist models.WishlistRequest
	if err := c.ShouldBindJSON(&wishlist); err != nil {
		errRes := response.ClientResponse(http.StatusBadRequest, "Invalid input data", nil, err.Error())
		c.JSON(http.StatusBadRequest, errRes)
		return
	}

	wishlist.UserID = UserID

	if wishlist.ProductID == 0 {
		errRes := response.ClientResponse(http.StatusBadRequest, "Product ID is required", nil, nil)
		c.JSON(http.StatusBadRequest, errRes)
		return
	}
	err := w.wishlistUseCase.AddToWishList(wishlist.ProductID, wishlist.UserID)
	if err != nil {
		errRes := response.ClientResponse(http.StatusBadRequest, "failed to item add to the wishlist", nil, err.Error())
		c.JSON(http.StatusBadRequest, errRes)
		return
	}
	successRes := response.ClientResponse(http.StatusOK, "SuccessFully added product to the wishlist", wishlist, nil)
	c.JSON(http.StatusOK, successRes)
}

// RemoveFromWishList godoc
// @Summary Remove a product from the wishlist
// @Description Removes a specified product from the authenticated user's wishlist
// @Tags Wishlist
// @Param product_id query int true "Product ID"
// @Success 200 {object} response.ClientResponse
// @Failure 401 {object} response.ClientResponse
// @Failure 400 {object} response.ClientResponse
// @Router /wishlist [delete]

func (w *WishlistHandler) RemoveFromWishList(c *gin.Context) {
	userID, ok := c.Get("id")
	if !ok {
		errRes := response.ClientResponse(http.StatusUnauthorized, "User ID not found in context", nil, nil)
		c.JSON(http.StatusUnauthorized, errRes)
		return
	}
	UserID := userID.(int)
	id := c.Query("product_id")
	productID, err := strconv.Atoi(id)
	err = w.wishlistUseCase.RemoveFromWishList(productID, UserID)
	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "failed to remove item from wishlist", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}
	successRes := response.ClientResponse(http.StatusOK, "SuccessFully deleted product from wishlist", nil, nil)
	c.JSON(http.StatusOK, successRes)

}

// GetWishList godoc
// @Summary Get the user's wishlist
// @Description Retrieves the wishlist for the authenticated user
// @Tags Wishlist
// @Success 200 {object} response.ClientResponse
// @Failure 401 {object} response.ClientResponse
// @Failure 400 {object} response.ClientResponse
// @Router /wishlist [get]

func (w *WishlistHandler) GetWishList(c *gin.Context) {
	userID, ok := c.Get("id")
	if !ok {
		errRes := response.ClientResponse(http.StatusUnauthorized, "User ID Not found in context", nil, nil)
		c.JSON(http.StatusUnauthorized, errRes)
		return
	}
	UserID := userID.(int)

	wishlist, err := w.wishlistUseCase.GetWishList(UserID)
	if err != nil {
		errRes := response.ClientResponse(http.StatusBadRequest, "failed to retrieve wishlist detailss", nil, err.Error())
		c.JSON(http.StatusBadRequest, errRes)
		return
	}
	successRes := response.ClientResponse(http.StatusOK, "SuccessFully retrieved wishlist", nil, wishlist)
	c.JSON(http.StatusOK, successRes)
}

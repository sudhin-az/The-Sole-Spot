package handlers

import (
	"ecommerce_clean_architecture/pkg/usecase"
	"ecommerce_clean_architecture/pkg/utils/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

type WalletHandler struct {
	usecase usecase.WalletUseCase
}

func NewWalletHandler(usecase usecase.WalletUseCase) *WalletHandler {
	return &WalletHandler{usecase: usecase}
}

// ViewWallet godoc
// @Summary View user wallet
// @Description Retrieves the wallet details for the authenticated user
// @Tags Wallet
// @Produce json
// @Success 200 {object} response.ClientResponse
// @Failure 401 {object} response.ClientResponse
// @Failure 400 {object} response.ClientResponse
// @Router /wallet [get]

func (wal *WalletHandler) ViewWallet(c *gin.Context) {
	userID, ok := c.Get("id")
	if !ok {
		errRes := response.ClientResponse(http.StatusUnauthorized, "User ID not found in context", nil, nil)
		c.JSON(http.StatusUnauthorized, errRes)
		return
	}
	UserID := userID.(int)
	wallet, err := wal.usecase.GetUserWallet(UserID)
	if err != nil {
		errRes := response.ClientResponse(http.StatusBadRequest, "fields provided are wrong", nil, err.Error())
		c.JSON(http.StatusBadRequest, errRes)
		return
	}
	successRes := response.ClientResponse(http.StatusOK, "Walllet is successfully shown", wallet, nil)
	c.JSON(http.StatusOK, successRes)
}

// GetWalletTransaction godoc
// @Summary Get wallet transaction history
// @Description Retrieves the transaction history for the authenticated user's wallet
// @Tags Wallet
// @Produce json
// @Success 200 {object} response.ClientResponse
// @Failure 401 {object} response.ClientResponse
// @Failure 400 {object} response.ClientResponse
// @Router /wallet/transactions [get]

func (wal *WalletHandler) GetWalletTransaction(c *gin.Context) {
	userID, ok := c.Get("id")
	if !ok {
		errRes := response.ClientResponse(http.StatusUnauthorized, "User ID not found in context", nil, nil)
		c.JSON(http.StatusUnauthorized, errRes)
		return
	}
	UserID := userID.(int)
	walletHistory, err := wal.usecase.GetWalletTransaction(UserID)
	if err != nil {
		errRes := response.ClientResponse(http.StatusBadRequest, "wallet history cannot be retrived", nil, err.Error())
		c.JSON(http.StatusBadRequest, errRes)
		return
	}
	successRes := response.ClientResponse(http.StatusOK, "WallletHistory is successfully shown", walletHistory, nil)
	c.JSON(http.StatusOK, successRes)
}

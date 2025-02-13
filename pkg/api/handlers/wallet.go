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

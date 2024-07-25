package http_handler

import (
	"banking-app/services/wallet"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type WalletHandler struct {
	walletService wallet.WalletService
}

func NewWalletHandler(walletService wallet.WalletService) *WalletHandler {
	return &WalletHandler{walletService: walletService}
}

type CrDrRequestBody struct {
	AccountId primitive.ObjectID `json:"account_id" binding:"required"`
	Amount    float64            `json:"amount" binding:"required"`
	Reference string             `json:"reference" binding:"required"`
}

type CrDrResponseBody struct {
	Balance float64 `json:"balance"`
}

// creditWallet godoc
// @Summary Credit Wallet
// @Description Credit Wallet
// @Tags credit
// @Accept json
// @Produce json
// @Param payment body CrDrRequestBody true "Credit wallet request"
// @Success 201 {object} CrDrResponseBody
// @Router /bank-api/credit-wallet [post]
func (h *WalletHandler) CreditWallet(c *gin.Context) {
	var req CrDrRequestBody
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	balance, err := h.walletService.CreditWallet(req.AccountId, req.Amount, req.Reference)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, CrDrResponseBody{Balance: balance})
}

// debitWallet godoc
// @Summary Debit Wallet
// @Description Debit Wallet
// @Tags debit
// @Accept json
// @Produce json
// @Param payment body CrDrRequestBody true "Debit wallet request"
// @Success 201 {object} CrDrResponseBody
// @Router /bank-api/debit-wallet [post]
func (h *WalletHandler) DebitWallet(c *gin.Context) {
	var req CrDrRequestBody
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	balance, err := h.walletService.DebitWallet(req.AccountId, req.Amount, req.Reference)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, CrDrResponseBody{Balance: balance})
}

package http_handler

import (
	"net/http"
	"third-party/services/transaction"

	"github.com/gin-gonic/gin"
)

type TransactionHandler struct {
	transactionService transaction.TransactionService
}

func NewWTransactionHandler(transactionService transaction.TransactionService) *TransactionHandler {
	return &TransactionHandler{transactionService: transactionService}
}

type CreatePayment struct {
	AccountId string  `json:"account_id" binding:"required"`
	Reference string  `json:"reference" binding:"required"`
	Amount    float64 `json:"amount" binding:"required"`
}

// createTransaction godoc
// @Summary Create a Transaction
// @Description Create a new Transaction
// @Tags payments
// @Accept json
// @Produce json
// @Param payment body CreatePayment true "Payment info"
// @Success 201 {object} models.Transaction
// @Router /third-party/payments [post]
func (h *TransactionHandler) CreatePaymentTransaction(c *gin.Context) {
	req := CreatePayment{}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	transaction, err := h.transactionService.CreatePaymentTransaction(req.AccountId, req.Reference, req.Amount)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, transaction)

}

// getTransaction godoc
// @Summary Get a transaction
// @Description Get Transaction by payment reference
// @Tags payments
// @Produce json
// @Param reference path string true "Payment Reference"
// @Success 200 {object} models.Transaction
// @Router /third-party/payments/{reference} [get]
func (h *TransactionHandler) GetPaymentTransaction(c *gin.Context) {
	reference := c.Param("reference")

	transaction, err := h.transactionService.GetPaymentTransaction(reference)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, transaction)
}

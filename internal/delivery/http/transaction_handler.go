package http

import (
	"net/http"
	"strconv"

	"xyz-multifinance/internal/model"
	"xyz-multifinance/internal/usecase"

	"github.com/gin-gonic/gin"
)

type TransactionHandler struct {
	transactionUsecase usecase.TransactionUsecase
}

func NewTransactionHandler(uc usecase.TransactionUsecase) *TransactionHandler {
	return &TransactionHandler{transactionUsecase: uc}
}

func (h *TransactionHandler) CreateTransaction(c *gin.Context) {
	var tx model.Transaction
	if err := c.ShouldBindJSON(&tx); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	err := h.transactionUsecase.CreateTransaction(&tx)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Transaction created"})
}

func (h *TransactionHandler) GetTransactionsByCustomer(c *gin.Context) {
	customerIDStr := c.Param("customer_id")
	customerID, err := strconv.ParseUint(customerIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid customer_id"})
		return
	}

	txs, err := h.transactionUsecase.GetTransactionsByCustomer(uint(customerID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Transactions not found"})
		return
	}

	c.JSON(http.StatusOK, txs)
}

func (h *TransactionHandler) UpdateTransaction(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid transaction id"})
		return
	}

	var updatedTx model.Transaction
	if err := c.ShouldBindJSON(&updatedTx); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	err = h.transactionUsecase.UpdateTransaction(uint(id), &updatedTx)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Transaction updated"})
}

func (h *TransactionHandler) DeleteTransaction(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid transaction id"})
		return
	}

	err = h.transactionUsecase.DeleteTransaction(uint(id))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Transaction deleted"})
}

package http

import (
	"net/http"
	"strconv"

	"xyz-multifinance/internal/model"
	"xyz-multifinance/internal/usecase"

	"github.com/gin-gonic/gin"
)

type LimitHandler struct {
	limitUsecase usecase.LimitUsecase
}

func NewLimitHandler(uc usecase.LimitUsecase) *LimitHandler {
	return &LimitHandler{limitUsecase: uc}
}

func (h *LimitHandler) CreateLimit(c *gin.Context) {
	var limit model.Limit
	if err := c.ShouldBindJSON(&limit); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	if err := h.limitUsecase.CreateLimit(&limit); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Limit created successfully"})
}

func (h *LimitHandler) DeleteLimit(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid limit ID"})
		return
	}

	if err := h.limitUsecase.DeleteLimit(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Limit deleted successfully"})
}

func (h *LimitHandler) GetLimitByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	limitWithRemaining, err := h.limitUsecase.GetLimitByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Limit not found"})
		return
	}
	if limitWithRemaining == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Limit not found"})
		return
	}

	c.JSON(http.StatusOK, limitWithRemaining)
}

func (h *LimitHandler) GetLimitsByCustomerID(c *gin.Context) {
	customerIDStr := c.Param("customer_id")
	customerID, err := strconv.ParseUint(customerIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid customer_id"})
		return
	}

	limits, err := h.limitUsecase.GetLimitsByCustomer(uint(customerID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get limits"})
		return
	}

	c.JSON(http.StatusOK, limits)
}

func (h *LimitHandler) GetLimitByCustomerAndTenor(c *gin.Context) {
	customerIDStr := c.Param("customer_id")
	tenorStr := c.Param("tenor")

	customerID, err := strconv.ParseUint(customerIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid customer_id"})
		return
	}

	tenor, err := strconv.Atoi(tenorStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tenor"})
		return
	}

	limits, err := h.limitUsecase.GetLimitsByCustomer(uint(customerID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Limit not found"})
		return
	}

	var limitFound *usecase.LimitWithRemaining
	for i, l := range limits {
		if l.Tenor == tenor {
			limitFound = &limits[i]
			break
		}
	}

	if limitFound == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Limit for tenor not found"})
		return
	}

	c.JSON(http.StatusOK, limitFound)
}

func (h *LimitHandler) UpdateLimit(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var updateData map[string]interface{}
	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	allowedFields := map[string]bool{"tenor_month": true, "limit_amount": true}
	for key := range updateData {
		if !allowedFields[key] {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid field: " + key})
			return
		}
	}

	err = h.limitUsecase.UpdateLimit(uint(id), updateData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Limit updated successfully"})
}

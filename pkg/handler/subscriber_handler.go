package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"hirelly-backend/pkg/domain"
	"hirelly-backend/pkg/service"
)

// SubscriberHandler handles subscriber API operations.
type SubscriberHandler struct {
	svc service.SubscriberService
}

// NewSubscriberHandler constructs a SubscriberHandler.
func NewSubscriberHandler(svc service.SubscriberService) *SubscriberHandler {
	return &SubscriberHandler{svc: svc}
}

// Subscribe handles registration of talent / corporate subscribers.
func (h *SubscriberHandler) Subscribe(c *gin.Context) {
	var req domain.CreateSubscriberRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body", "details": err.Error()})
		return
	}

	sub, err := h.svc.Subscribe(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "subscribed successfully",
		"record":  sub,
	})
}

// List handles listing of subscribers.
func (h *SubscriberHandler) List(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	subs, err := h.svc.ListSubscribers(c.Request.Context(), limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed retrieving subscribers"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success":     true,
		"subscribers": subs,
		"total":       len(subs),
	})
}

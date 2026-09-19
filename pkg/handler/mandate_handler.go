package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"hirelly-backend/pkg/domain"
	"hirelly-backend/pkg/service"
)

// MandateHandler handles executive mandate API operations.
type MandateHandler struct {
	svc service.AdvisoryService
}

// NewMandateHandler constructs a MandateHandler.
func NewMandateHandler(svc service.AdvisoryService) *MandateHandler {
	return &MandateHandler{svc: svc}
}

// Create handles incoming mandate/application submissions.
func (h *MandateHandler) Create(c *gin.Context) {
	var req domain.CreateMandateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body", "details": err.Error()})
		return
	}

	mandate, emailStatus, err := h.svc.CreateMandate(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success":         true,
		"savedToDatabase": true,
		"emailStatus":     emailStatus,
		"record":          mandate,
	})
}

// List handles retrieval of mandates with pagination.
func (h *MandateHandler) List(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	mandates, err := h.svc.ListMandates(c.Request.Context(), limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed retrieving mandates", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success":      true,
		"applications": mandates,
		"total":        len(mandates),
	})
}

// GetByID retrieves a single mandate by its identifier.
func (h *MandateHandler) GetByID(c *gin.Context) {
	id := c.Param("id")
	mandate, err := h.svc.GetMandate(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "mandate not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"record":  mandate,
	})
}

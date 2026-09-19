package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"hirelly-backend/internal/repository"
)

// JobHandler handles job opportunity retrieval.
type JobHandler struct {
	repo repository.JobRepository
}

// NewJobHandler constructs a JobHandler.
func NewJobHandler(repo repository.JobRepository) *JobHandler {
	return &JobHandler{repo: repo}
}

// List returns active executive search opportunities.
func (h *JobHandler) List(c *gin.Context) {
	jobs, err := h.repo.List(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed retrieving jobs"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"jobs":    jobs,
		"source":  "hirelly_go_core",
	})
}

package handler

import (
	"context"
	"net/http"
	"runtime"
	"time"

	"github.com/gin-gonic/gin"
	"hirelly-backend/pkg/config"
	"hirelly-backend/pkg/repository/postgres"
)

var startTime = time.Now()

// HealthHandler provides health check endpoints.
type HealthHandler struct {
	cfg *config.Config
	db  *postgres.Database
}

// NewHealthHandler constructs a HealthHandler.
func NewHealthHandler(cfg *config.Config, db *postgres.Database) *HealthHandler {
	return &HealthHandler{cfg: cfg, db: db}
}

// Check responds with detailed service and infrastructure diagnostics.
func (h *HealthHandler) Check(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
	defer cancel()

	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	dbStatus := "connected"
	if h.db.Pool != nil {
		if err := h.db.Pool.Ping(ctx); err != nil {
			dbStatus = "degraded (ping failed)"
		}
	} else if h.db.SupabaseURL != "" {
		dbStatus = "supabase_rest_active"
	} else {
		dbStatus = "not_configured"
	}

	c.JSON(http.StatusOK, gin.H{
		"service":     "Hirelly Executive Search Backend Core",
		"status":      "healthy",
		"environment": h.cfg.Env,
		"uptime":      time.Since(startTime).String(),
		"database": gin.H{
			"status":   dbStatus,
			"has_pool": h.db.Pool != nil,
			"supabase": h.db.SupabaseURL != "",
		},
		"system": gin.H{
			"go_version": runtime.Version(),
			"goroutines": runtime.NumGoroutine(),
			"alloc_mb":   m.Alloc / 1024 / 1024,
		},
	})
}

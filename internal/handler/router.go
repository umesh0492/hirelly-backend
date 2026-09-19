package handler

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/umesh0492/go-libs/ginmw"
	"github.com/umesh0492/go-libs/health"

	"hirelly-backend/internal/config"
	"hirelly-backend/internal/repository"
	"hirelly-backend/internal/repository/postgres"
	"hirelly-backend/internal/service"
)

var appStartTime = time.Now()

// RouterDeps encapsulates handler dependencies.
type RouterDeps struct {
	Config     *config.Config
	DB         *postgres.Database
	MandateSvc service.AdvisoryService
	SubSvc     service.SubscriberService
	JobRepo    repository.JobRepository
}

// SetupRouter initializes the Gin engine with Abeta-dev/go-libs middleware and API routes.
func SetupRouter(deps *RouterDeps) *gin.Engine {
	if deps.Config.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()

	// 1. Enterprise Middlewares from github.com/Abeta-dev/go-libs/ginmw
	r.Use(ginmw.RequestID())
	r.Use(ginmw.SecurityHeaders("hirelly-api"))
	r.Use(ginmw.Recovery())
	r.Use(ginmw.Logger())
	r.Use(func(c *gin.Context) {
		if c.Request.Body != nil && c.Request.ContentLength > 0 {
			ginmw.LimitBodyDefault()(c)
		} else {
			c.Next()
		}
	})
	r.Use(CORSMiddleware(deps.Config.CORSOrigins))

	// 2. Health & Readiness Engine from github.com/Abeta-dev/go-libs/health
	healthChecker := health.New(
		health.WithStartTime(appStartTime),
		health.WithChecker("postgres", func(ctx context.Context) error {
			if deps.DB.Pool != nil {
				return deps.DB.Pool.Ping(ctx)
			}
			if deps.DB.SupabaseURL != "" {
				return nil // Supabase REST active
			}
			return fmt.Errorf("no database configured")
		}),
	)

	// Mount health endpoints
	r.GET("/health", gin.WrapF(healthChecker.Handler))
	r.GET("/ready", gin.WrapF(healthChecker.Handler))

	// Root welcome endpoint
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"platform":   "Hirelly Executive Search Engine",
			"status":     "operational",
			"framework":  "Go + Abeta-dev/go-libs + Abeta-dev/go-app-kit",
			"health":     "/health",
			"request_id": c.GetString(ginmw.RequestIDHeader),
		})
	})

	mandateHandler := NewMandateHandler(deps.MandateSvc)
	subHandler := NewSubscriberHandler(deps.SubSvc)
	jobHandler := NewJobHandler(deps.JobRepo)

	// Versioned API routes (/api/v1/...)
	v1 := r.Group("/api/v1")
	{
		v1.POST("/mandates", mandateHandler.Create)
		v1.GET("/mandates", mandateHandler.List)
		v1.GET("/mandates/:id", mandateHandler.GetByID)

		v1.POST("/subscribers", subHandler.Subscribe)
		v1.GET("/subscribers", subHandler.List)

		v1.GET("/jobs", jobHandler.List)
	}

	// 1:1 Backward-compatible Vercel frontend routes (/api/...)
	legacy := r.Group("/api")
	{
		legacy.POST("/applications", mandateHandler.Create)
		legacy.GET("/applications", mandateHandler.List)
		legacy.GET("/jobs", jobHandler.List)
	}

	return r
}

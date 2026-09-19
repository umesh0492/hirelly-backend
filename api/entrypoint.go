package handler

import (
	"context"
	"log"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"hirelly-backend/pkg/config"
	"hirelly-backend/pkg/handler"
	"hirelly-backend/pkg/repository/postgres"
	"hirelly-backend/pkg/service"
)

var (
	router *gin.Engine
	once   sync.Once
)

func initApp() {
	cfg := config.Load()
	ctx := context.Background()

	db, err := postgres.NewDatabase(ctx, cfg)
	if err != nil {
		log.Printf("⚠️ Vercel Go Serverless Database warning: %v", err)
	}

	mandateRepo := postgres.NewMandateRepository(db)
	subRepo := postgres.NewSubscriberRepository(db)
	jobRepo := postgres.NewJobRepository(db)

	emailSvc := service.NewEmailService(cfg)
	mandateSvc := service.NewAdvisoryService(mandateRepo, emailSvc)
	subSvc := service.NewSubscriberService(subRepo)

	router = handler.SetupRouter(&handler.RouterDeps{
		Config:     cfg,
		DB:         db,
		MandateSvc: mandateSvc,
		SubSvc:     subSvc,
		JobRepo:    jobRepo,
	})
}

// Handler is the Vercel serverless function entrypoint.
func Handler(w http.ResponseWriter, r *http.Request) {
	once.Do(initApp)
	router.ServeHTTP(w, r)
}

package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"hirelly-backend/pkg/config"
	"hirelly-backend/pkg/handler"
	"hirelly-backend/pkg/repository/postgres"
	"hirelly-backend/pkg/service"
)

func main() {
	log.Println("⚡ Starting Hirelly Executive Backend Core Engine...")

	// 1. Load configuration
	cfg := config.Load()

	// 2. Initialize Database layer (PostgreSQL connection pool + Supabase REST)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	db, err := postgres.NewDatabase(ctx, cfg)
	if err != nil {
		log.Fatalf("❌ Database initialization error: %v", err)
	}
	defer db.Close()

	// 3. Initialize Repositories
	mandateRepo := postgres.NewMandateRepository(db)
	subRepo := postgres.NewSubscriberRepository(db)
	jobRepo := postgres.NewJobRepository(db)

	// 4. Initialize Services
	emailSvc := service.NewEmailService(cfg)
	mandateSvc := service.NewAdvisoryService(mandateRepo, emailSvc)
	subSvc := service.NewSubscriberService(subRepo)

	// 5. Initialize Router
	router := handler.SetupRouter(&handler.RouterDeps{
		Config:     cfg,
		DB:         db,
		MandateSvc: mandateSvc,
		SubSvc:     subSvc,
		JobRepo:    jobRepo,
	})

	// 6. Setup HTTP Server
	server := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// 7. Start server in goroutine
	go func() {
		fmt.Printf("🚀 Hirelly Backend Service listening on http://localhost:%s (env: %s)\n", cfg.Port, cfg.Env)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("❌ HTTP server error: %v", err)
		}
	}()

	// 8. Graceful shutdown handler
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("🛑 Shutting down Hirelly Backend gracefully...")
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Printf("⚠️ Server shutdown error: %v", err)
	}

	log.Println("✅ Hirelly Backend stopped cleanly.")
}

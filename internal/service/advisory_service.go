package service

import (
	"context"
	"fmt"
	"log"
	"time"

	"hirelly-backend/internal/domain"
	"hirelly-backend/internal/repository"
	"hirelly-backend/pkg/uuid"
)

// AdvisoryService coordinates executive mandate operations.
type AdvisoryService interface {
	CreateMandate(ctx context.Context, req *domain.CreateMandateRequest) (*domain.Mandate, string, error)
	ListMandates(ctx context.Context, limit, offset int) ([]*domain.Mandate, error)
	GetMandate(ctx context.Context, id string) (*domain.Mandate, error)
}

type advisoryService struct {
	repo  repository.MandateRepository
	email EmailService
}

// NewAdvisoryService constructs an AdvisoryService instance.
func NewAdvisoryService(repo repository.MandateRepository, email EmailService) AdvisoryService {
	return &advisoryService{
		repo:  repo,
		email: email,
	}
}

func (s *advisoryService) CreateMandate(ctx context.Context, req *domain.CreateMandateRequest) (*domain.Mandate, string, error) {
	if err := req.Validate(); err != nil {
		return nil, "", err
	}

	mandate := &domain.Mandate{
		ID:        uuid.New(),
		Name:      req.Name,
		Email:     req.Email,
		Org:       req.Org,
		Type:      req.Type,
		Message:   req.Message,
		Status:    domain.StatusNew,
		CreatedAt: time.Now().UTC(),
	}

	// 1. Persist mandate into database (PostgreSQL pool / Supabase)
	if err := s.repo.Create(ctx, mandate); err != nil {
		return nil, "", fmt.Errorf("failed persisting mandate: %w", err)
	}

	// 2. Dispatch transactional confirmation email
	var emailStatus string
	msgID, err := s.email.SendMandateConfirmation(ctx, mandate)
	if err != nil {
		log.Printf("⚠️ Mandate created (%s) but email delivery failed: %v", mandate.ID, err)
		emailStatus = "failed"
	} else {
		log.Printf("📧 Mandate confirmation sent. Message ID: %s", msgID)
		emailStatus = "sent"
	}

	return mandate, emailStatus, nil
}

func (s *advisoryService) ListMandates(ctx context.Context, limit, offset int) ([]*domain.Mandate, error) {
	return s.repo.List(ctx, limit, offset)
}

func (s *advisoryService) GetMandate(ctx context.Context, id string) (*domain.Mandate, error) {
	return s.repo.GetByID(ctx, id)
}

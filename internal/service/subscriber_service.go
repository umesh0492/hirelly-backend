package service

import (
	"context"
	"fmt"
	"time"

	"hirelly-backend/internal/domain"
	"hirelly-backend/internal/repository"
	"hirelly-backend/pkg/uuid"
)

// SubscriberService coordinates subscriber registration and listing.
type SubscriberService interface {
	Subscribe(ctx context.Context, req *domain.CreateSubscriberRequest) (*domain.Subscriber, error)
	ListSubscribers(ctx context.Context, limit int) ([]*domain.Subscriber, error)
}

type subscriberService struct {
	repo repository.SubscriberRepository
}

// NewSubscriberService constructs a SubscriberService instance.
func NewSubscriberService(repo repository.SubscriberRepository) SubscriberService {
	return &subscriberService{repo: repo}
}

func (s *subscriberService) Subscribe(ctx context.Context, req *domain.CreateSubscriberRequest) (*domain.Subscriber, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	sub := &domain.Subscriber{
		ID:        uuid.New(),
		Email:     req.Email,
		Role:      req.Role,
		CreatedAt: time.Now().UTC(),
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	}

	if err := s.repo.Create(ctx, sub); err != nil {
		return nil, fmt.Errorf("failed creating subscriber: %w", err)
	}

	return sub, nil
}

func (s *subscriberService) ListSubscribers(ctx context.Context, limit int) ([]*domain.Subscriber, error) {
	return s.repo.List(ctx, limit)
}

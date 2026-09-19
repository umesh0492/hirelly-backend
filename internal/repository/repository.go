package repository

import (
	"context"
	"hirelly-backend/internal/domain"
)

// MandateRepository abstracts persistence operations for executive mandates.
type MandateRepository interface {
	Create(ctx context.Context, mandate *domain.Mandate) error
	GetByID(ctx context.Context, id string) (*domain.Mandate, error)
	List(ctx context.Context, limit, offset int) ([]*domain.Mandate, error)
}

// SubscriberRepository abstracts persistence for candidate/client subscribers.
type SubscriberRepository interface {
	Create(ctx context.Context, subscriber *domain.Subscriber) error
	List(ctx context.Context, limit int) ([]*domain.Subscriber, error)
}

// JobRepository abstracts retrieval of active executive roles.
type JobRepository interface {
	List(ctx context.Context) ([]*domain.Job, error)
}

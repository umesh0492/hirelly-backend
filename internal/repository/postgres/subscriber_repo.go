package postgres

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"hirelly-backend/internal/domain"
	"hirelly-backend/internal/repository"
)

type subscriberRepository struct {
	db *Database
}

// NewSubscriberRepository constructs a new subscriber repository.
func NewSubscriberRepository(db *Database) repository.SubscriberRepository {
	return &subscriberRepository{db: db}
}

func (r *subscriberRepository) Create(ctx context.Context, s *domain.Subscriber) error {
	if s.CreatedAt.IsZero() {
		s.CreatedAt = time.Now().UTC()
	}
	if s.Timestamp == "" {
		s.Timestamp = s.CreatedAt.Format(time.RFC3339)
	}

	// 1. Direct PostgreSQL pool
	if r.db.Pool != nil {
		query := `
			INSERT INTO subscribers (id, email, role, created_at, timestamp)
			VALUES ($1, $2, $3, $4, $5)
			ON CONFLICT (id) DO NOTHING
		`
		_, err := r.db.Pool.Exec(ctx, query, s.ID, s.Email, s.Role, s.CreatedAt, s.Timestamp)
		if err == nil {
			return nil
		}
	}

	// 2. Supabase REST API
	payload := map[string]any{
		"id":         s.ID,
		"email":      s.Email,
		"role":       s.Role,
		"created_at": s.CreatedAt.Format(time.RFC3339),
		"timestamp":  s.Timestamp,
	}

	_, err := r.db.QuerySupabase(ctx, "POST", "subscribers", payload)
	if err != nil {
		return fmt.Errorf("failed saving subscriber: %w", err)
	}

	return nil
}

func (r *subscriberRepository) List(ctx context.Context, limit int) ([]*domain.Subscriber, error) {
	if limit <= 0 {
		limit = 20
	}

	if r.db.Pool != nil {
		query := `SELECT id, email, role, created_at, COALESCE(timestamp, '') FROM subscribers ORDER BY created_at DESC LIMIT $1`
		rows, err := r.db.Pool.Query(ctx, query, limit)
		if err == nil {
			defer rows.Close()
			var subs []*domain.Subscriber
			for rows.Next() {
				var s domain.Subscriber
				if err := rows.Scan(&s.ID, &s.Email, &s.Role, &s.CreatedAt, &s.Timestamp); err == nil {
					subs = append(subs, &s)
				}
			}
			return subs, nil
		}
	}

	// Supabase REST fallback
	resBytes, err := r.db.QuerySupabase(ctx, "GET", fmt.Sprintf("subscribers?select=*&order=created_at.desc&limit=%d", limit), nil)
	if err != nil {
		return []*domain.Subscriber{}, nil
	}

	var subs []*domain.Subscriber
	if err := json.Unmarshal(resBytes, &subs); err != nil {
		return []*domain.Subscriber{}, nil
	}
	return subs, nil
}

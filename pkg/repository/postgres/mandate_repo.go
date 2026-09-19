package postgres

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"hirelly-backend/pkg/domain"
	"hirelly-backend/pkg/repository"
)

type mandateRepository struct {
	db *Database
}

// NewMandateRepository constructs a new mandate persistence repository.
func NewMandateRepository(db *Database) repository.MandateRepository {
	return &mandateRepository{db: db}
}

func (r *mandateRepository) Create(ctx context.Context, m *domain.Mandate) error {
	if m.CreatedAt.IsZero() {
		m.CreatedAt = time.Now().UTC()
	}
	if m.Status == "" {
		m.Status = domain.StatusNew
	}

	// 1. Attempt direct PostgreSQL Pool query if pool is active
	if r.db.Pool != nil {
		query := `
			INSERT INTO applications (id, name, email, org, type, message, status, created_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		`
		_, err := r.db.Pool.Exec(ctx, query, m.ID, m.Name, m.Email, m.Org, m.Type, m.Message, string(m.Status), m.CreatedAt)
		if err == nil {
			return nil
		}
		log.Printf("⚠️ Direct postgres applications insert error: %v. Trying subscribers table fallback...", err)

		// Fallback to subscribers table in direct Postgres
		subQuery := `
			INSERT INTO subscribers (id, email, role, created_at, timestamp)
			VALUES ($1, $2, $3, $4, $5)
		`
		_, subErr := r.db.Pool.Exec(ctx, subQuery, m.ID, m.Email, m.Type, m.CreatedAt, m.CreatedAt.Format(time.RFC3339))
		if subErr == nil {
			return nil
		}
		log.Printf("⚠️ Direct postgres subscribers fallback error: %v. Trying Supabase REST API...", subErr)
	}

	// 2. Supabase REST API fallback
	appPayload := map[string]any{
		"id":         m.ID,
		"name":       m.Name,
		"email":      m.Email,
		"org":        m.Org,
		"type":       m.Type,
		"message":    m.Message,
		"created_at": m.CreatedAt.Format(time.RFC3339),
	}

	_, err := r.db.QuerySupabase(ctx, "POST", "applications", appPayload)
	if err == nil {
		return nil
	}

	// If 'applications' table does not exist in schema cache, fallback to 'subscribers' table
	subPayload := map[string]any{
		"id":         m.ID,
		"email":      m.Email,
		"role":       m.Type,
		"created_at": m.CreatedAt.Format(time.RFC3339),
		"timestamp":  m.CreatedAt.Format(time.RFC3339),
	}

	_, subErr := r.db.QuerySupabase(ctx, "POST", "subscribers", subPayload)
	if subErr != nil {
		return fmt.Errorf("failed saving to both applications and subscribers: %w", subErr)
	}

	return nil
}

func (r *mandateRepository) GetByID(ctx context.Context, id string) (*domain.Mandate, error) {
	if r.db.Pool != nil {
		query := `SELECT id, name, email, org, type, message, status, created_at FROM applications WHERE id = $1`
		row := r.db.Pool.QueryRow(ctx, query, id)
		var m domain.Mandate
		var statusStr string
		err := row.Scan(&m.ID, &m.Name, &m.Email, &m.Org, &m.Type, &m.Message, &statusStr, &m.CreatedAt)
		if err == nil {
			m.Status = domain.MandateStatus(statusStr)
			return &m, nil
		}
	}

	// Supabase REST fallback
	resBytes, err := r.db.QuerySupabase(ctx, "GET", fmt.Sprintf("applications?id=eq.%s&select=*", id), nil)
	if err != nil {
		return nil, domain.ErrNotFound
	}

	var results []domain.Mandate
	if err := json.Unmarshal(resBytes, &results); err == nil && len(results) > 0 {
		return &results[0], nil
	}

	return nil, domain.ErrNotFound
}

func (r *mandateRepository) List(ctx context.Context, limit, offset int) ([]*domain.Mandate, error) {
	if limit <= 0 {
		limit = 20
	}

	if r.db.Pool != nil {
		query := `SELECT id, name, email, org, type, message, status, created_at FROM applications ORDER BY created_at DESC LIMIT $1 OFFSET $2`
		rows, err := r.db.Pool.Query(ctx, query, limit, offset)
		if err == nil {
			defer rows.Close()
			var mandates []*domain.Mandate
			for rows.Next() {
				var m domain.Mandate
				var statusStr string
				if err := rows.Scan(&m.ID, &m.Name, &m.Email, &m.Org, &m.Type, &m.Message, &statusStr, &m.CreatedAt); err == nil {
					m.Status = domain.MandateStatus(statusStr)
					mandates = append(mandates, &m)
				}
			}
			return mandates, nil
		}
	}

	// Supabase REST fallback
	resBytes, err := r.db.QuerySupabase(ctx, "GET", fmt.Sprintf("applications?select=*&order=created_at.desc&limit=%d&offset=%d", limit, offset), nil)
	if err != nil {
		// If applications table is missing, try subscribers table
		subBytes, subErr := r.db.QuerySupabase(ctx, "GET", fmt.Sprintf("subscribers?select=*&order=created_at.desc&limit=%d&offset=%d", limit, offset), nil)
		if subErr != nil {
			return []*domain.Mandate{}, nil
		}
		var subs []domain.Subscriber
		if err := json.Unmarshal(subBytes, &subs); err == nil {
			var list []*domain.Mandate
			for _, s := range subs {
				list = append(list, &domain.Mandate{
					ID:        s.ID,
					Name:      "Subscriber",
					Email:     s.Email,
					Type:      s.Role,
					CreatedAt: s.CreatedAt,
					Status:    domain.StatusNew,
				})
			}
			return list, nil
		}
		return []*domain.Mandate{}, nil
	}

	var mandates []*domain.Mandate
	if err := json.Unmarshal(resBytes, &mandates); err != nil {
		return []*domain.Mandate{}, nil
	}
	return mandates, nil
}

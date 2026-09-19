package domain

import (
	"strings"
	"time"
)

// Subscriber represents a registered executive candidate, recruiter, or company.
type Subscriber struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at"`
	Timestamp string    `json:"timestamp,omitempty"`
}

// CreateSubscriberRequest is the payload for subscribing to the platform waitlist.
type CreateSubscriberRequest struct {
	Email string `json:"email"`
	Role  string `json:"role"`
}

// Validate validates the subscription request payload.
func (r *CreateSubscriberRequest) Validate() error {
	r.Email = strings.TrimSpace(r.Email)
	if r.Email == "" {
		return ErrEmailRequired
	}
	if !strings.Contains(r.Email, "@") || !strings.Contains(r.Email, ".") {
		return ErrInvalidEmail
	}
	if r.Role == "" {
		r.Role = "Executive"
	}
	return nil
}

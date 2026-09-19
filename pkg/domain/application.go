package domain

import (
	"strings"
	"time"
)

// MandateStatus represents the lifecycle state of an executive mandate inquiry.
type MandateStatus string

const (
	StatusNew        MandateStatus = "NEW"
	StatusInReview   MandateStatus = "IN_REVIEW"
	StatusContacted  MandateStatus = "CONTACTED"
	StatusClosed     MandateStatus = "CLOSED"
)

// Mandate represents an executive search mandate or contact submission.
type Mandate struct {
	ID        string        `json:"id"`
	Name      string        `json:"name"`
	Email     string        `json:"email"`
	Org       string        `json:"org"`
	Type      string        `json:"type"`
	Message   string        `json:"message"`
	Status    MandateStatus `json:"status"`
	CreatedAt time.Time     `json:"created_at"`
	UpdatedAt time.Time     `json:"updated_at,omitempty"`
}

// CreateMandateRequest defines the incoming payload from web forms.
type CreateMandateRequest struct {
	Name    string `json:"name"`
	Email   string `json:"email"`
	Org     string `json:"org"`
	Type    string `json:"type"`
	Message string `json:"msg"`
	Role    string `json:"role"`
}

// Validate checks for required fields in the mandate request.
func (r *CreateMandateRequest) Validate() error {
	r.Email = strings.TrimSpace(r.Email)
	if r.Email == "" {
		return ErrEmailRequired
	}
	if !strings.Contains(r.Email, "@") || !strings.Contains(r.Email, ".") {
		return ErrInvalidEmail
	}
	if r.Name == "" {
		r.Name = "Executive Leader"
	}
	if r.Org == "" {
		r.Org = "Confidential Organization"
	}
	if r.Type == "" {
		if r.Role != "" {
			r.Type = r.Role
		} else {
			r.Type = "Executive Search Request"
		}
	}
	return nil
}

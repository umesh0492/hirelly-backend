package postgres

import (
	"context"
	"encoding/json"
	"time"

	"hirelly-backend/pkg/domain"
	"hirelly-backend/pkg/repository"
)

type jobRepository struct {
	db *Database
}

// NewJobRepository constructs a job repository.
func NewJobRepository(db *Database) repository.JobRepository {
	return &jobRepository{db: db}
}

var defaultJobs = []*domain.Job{
	{
		ID:        "job-exec-1",
		Title:     "Chief Technology Officer (AI & Distributed Systems)",
		Company:   "Tier-1 AI Unicorn",
		Icon:      "fa-brain",
		Location:  "Bangalore / Singapore (Hybrid)",
		Type:      "Full-time / Executive",
		Salary:    "$220,000 - $350,000 + Equity",
		SalaryVal: 220000,
		Tags:      []string{"GenAI", "Golang", "Kubernetes", "C-Suite"},
		Posted:    "1 day ago",
		CreatedAt: time.Now().Add(-24 * time.Hour),
	},
	{
		ID:        "job-exec-2",
		Title:     "Managing Director — Private Equity Tech Advisory",
		Company:   "Global Capital Partners",
		Icon:      "fa-chart-line",
		Location:  "London / Dubai (Global Mobility)",
		Type:      "Partner Track",
		Salary:    "£250,000 - £400,000 + Carry",
		SalaryVal: 320000,
		Tags:      []string{"Private Equity", "M&A", "Board Advisory"},
		Posted:    "2 days ago",
		CreatedAt: time.Now().Add(-48 * time.Hour),
	},
	{
		ID:        "job-exec-3",
		Title:     "Head of Global Engineering",
		Company:   "Fintech Infrastructure Scaleup",
		Icon:      "fa-credit-card",
		Location:  "San Francisco / Gurgaon",
		Type:      "Full-time / Leadership",
		Salary:    "$200,000 - $310,000",
		SalaryVal: 200000,
		Tags:      []string{"High Frequency", "Rust", "Go", "Payments"},
		Posted:    "3 days ago",
		CreatedAt: time.Now().Add(-72 * time.Hour),
	},
}

func (r *jobRepository) List(ctx context.Context) ([]*domain.Job, error) {
	// Try Supabase REST query
	resBytes, err := r.db.QuerySupabase(ctx, "GET", "jobs?select=*&order=created_at.desc", nil)
	if err == nil {
		var liveJobs []*domain.Job
		if err := json.Unmarshal(resBytes, &liveJobs); err == nil && len(liveJobs) > 0 {
			return liveJobs, nil
		}
	}

	return defaultJobs, nil
}

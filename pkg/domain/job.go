package domain

import "time"

// Job represents an executive opportunity posted on Hirelly.
type Job struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	Company   string    `json:"company"`
	Icon      string    `json:"icon"`
	Location  string    `json:"location"`
	Type      string    `json:"type"`
	Salary    string    `json:"salary"`
	SalaryVal int       `json:"salary_val,omitempty"`
	Tags      []string  `json:"tags"`
	Posted    string    `json:"posted"`
	CreatedAt time.Time `json:"created_at"`
}

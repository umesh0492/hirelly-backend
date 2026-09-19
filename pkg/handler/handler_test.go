package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"hirelly-backend/pkg/config"
	"hirelly-backend/pkg/domain"
	"hirelly-backend/pkg/repository/postgres"
	"hirelly-backend/pkg/service"
)

// MockMandateRepo implements repository.MandateRepository in-memory for testing.
type mockMandateRepo struct {
	mandates []*domain.Mandate
}

func (m *mockMandateRepo) Create(ctx context.Context, mandate *domain.Mandate) error {
	m.mandates = append(m.mandates, mandate)
	return nil
}

func (m *mockMandateRepo) GetByID(ctx context.Context, id string) (*domain.Mandate, error) {
	for _, mandate := range m.mandates {
		if mandate.ID == id {
			return mandate, nil
		}
	}
	return nil, domain.ErrNotFound
}

func (m *mockMandateRepo) List(ctx context.Context, limit, offset int) ([]*domain.Mandate, error) {
	return m.mandates, nil
}

// MockSubscriberRepo implements repository.SubscriberRepository.
type mockSubscriberRepo struct {
	subs []*domain.Subscriber
}

func (m *mockSubscriberRepo) Create(ctx context.Context, s *domain.Subscriber) error {
	m.subs = append(m.subs, s)
	return nil
}

func (m *mockSubscriberRepo) List(ctx context.Context, limit int) ([]*domain.Subscriber, error) {
	return m.subs, nil
}

// MockEmailService implements service.EmailService.
type mockEmailService struct{}

func (m *mockEmailService) SendMandateConfirmation(ctx context.Context, mandate *domain.Mandate) (string, error) {
	return "mock_msg_12345", nil
}

func setupTestRouter() (*mockMandateRepo, *mockSubscriberRepo, http.Handler) {
	cfg := &config.Config{
		Env:         "test",
		CORSOrigins: []string{"*"},
	}

	mandateRepo := &mockMandateRepo{}
	subRepo := &mockSubscriberRepo{}
	jobRepo := postgres.NewJobRepository(&postgres.Database{})
	emailSvc := &mockEmailService{}

	mandateSvc := service.NewAdvisoryService(mandateRepo, emailSvc)
	subSvc := service.NewSubscriberService(subRepo)

	router := SetupRouter(&RouterDeps{
		Config:     cfg,
		DB:         &postgres.Database{SupabaseURL: "https://test.supabase.co"},
		MandateSvc: mandateSvc,
		SubSvc:     subSvc,
		JobRepo:    jobRepo,
	})

	return mandateRepo, subRepo, router
}

func TestHealthCheck(t *testing.T) {
	_, _, router := setupTestRouter()

	req, _ := http.NewRequest("GET", "/health", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp map[string]any
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "healthy", resp["status"])
	assert.NotNil(t, w.Header().Get("X-Request-ID"))
}

func TestCreateMandate_Success(t *testing.T) {
	repo, _, router := setupTestRouter()

	payload := domain.CreateMandateRequest{
		Name:    "Vikram Malhotra",
		Email:   "vikram@apexholdings.com",
		Org:     "Apex Holdings",
		Type:    "Board Advisory",
		Message: "Seeking independent non-executive directors.",
	}
	body, _ := json.Marshal(payload)

	req, _ := http.NewRequest("POST", "/api/v1/mandates", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp map[string]any
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.True(t, resp["success"].(bool))
	assert.Equal(t, "sent", resp["emailStatus"])
	assert.Len(t, repo.mandates, 1)
	assert.Equal(t, "Vikram Malhotra", repo.mandates[0].Name)
}

func TestCreateMandate_InvalidEmail(t *testing.T) {
	_, _, router := setupTestRouter()

	payload := domain.CreateMandateRequest{
		Name:  "Vikram Malhotra",
		Email: "not-an-email",
	}
	body, _ := json.Marshal(payload)

	req, _ := http.NewRequest("POST", "/api/v1/mandates", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCreateSubscriber_Success(t *testing.T) {
	_, subRepo, router := setupTestRouter()

	payload := domain.CreateSubscriberRequest{
		Email: "candidate.leader@hirelly.in",
		Role:  "Candidate",
	}
	body, _ := json.Marshal(payload)

	req, _ := http.NewRequest("POST", "/api/v1/subscribers", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Len(t, subRepo.subs, 1)
	assert.Equal(t, "candidate.leader@hirelly.in", subRepo.subs[0].Email)
}

func TestListJobs(t *testing.T) {
	_, _, router := setupTestRouter()

	req, _ := http.NewRequest("GET", "/api/v1/jobs", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp map[string]any
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.True(t, resp["success"].(bool))
	jobs := resp["jobs"].([]any)
	assert.NotEmpty(t, jobs)
}

func TestLegacyApplicationsEndpoint(t *testing.T) {
	repo, _, router := setupTestRouter()

	payload := map[string]string{
		"name":  "Legacy Client",
		"email": "client@legacy.com",
		"org":   "Legacy Corp",
		"type":  "Executive Search Request",
		"msg":   "Mandate request from legacy Vercel form",
	}
	body, _ := json.Marshal(payload)

	req, _ := http.NewRequest("POST", "/api/applications", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Len(t, repo.mandates, 1)
	assert.Equal(t, "Legacy Client", repo.mandates[0].Name)
}

package httpapi

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"project-budget-service/internal/domain"
)

func TestUnauthorizedRequestUsesUnifiedError(t *testing.T) {
	handler := NewHandler(&fakeService{})
	mux := http.NewServeMux()
	handler.Register(mux)

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/projects", nil)
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", rec.Code)
	}
	var body apiResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatal(err)
	}
	if body.Error == nil || body.Error.Code != string(domain.CodeUnauthorized) {
		t.Fatalf("unexpected error body: %+v", body.Error)
	}
}

func TestCreateProjectRoute(t *testing.T) {
	handler := NewHandler(&fakeService{})
	mux := http.NewServeMux()
	handler.Register(mux)

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/v1/projects", strings.NewReader(`{"name":"ERP"}`))
	req.Header.Set("Authorization", "Bearer 11111111-1111-1111-1111-111111111111")
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d: %s", rec.Code, rec.Body.String())
	}
}

func TestGetResourcesRoute(t *testing.T) {
	handler := NewHandler(&fakeService{})
	mux := http.NewServeMux()
	handler.Register(mux)

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/resources", nil)
	req.Header.Set("Authorization", "Bearer 11111111-1111-1111-1111-111111111111")
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
}

type fakeService struct{}

func (f *fakeService) CreateProject(ctx context.Context, userID domain.UserID, in domain.CreateProject) (*domain.Project, error) {
	return &domain.Project{ID: "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa", Name: in.Name, Status: domain.StatusPlanning}, nil
}

func (f *fakeService) ListProjects(ctx context.Context, userID domain.UserID) ([]domain.Project, error) {
	return []domain.Project{{ID: "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa", Name: "ERP"}}, nil
}

func (f *fakeService) GetProject(ctx context.Context, userID domain.UserID, projectID domain.ProjectID) (*domain.Project, error) {
	return &domain.Project{ID: projectID, Name: "ERP"}, nil
}

func (f *fakeService) UpdateProject(ctx context.Context, userID domain.UserID, projectID domain.ProjectID, in domain.UpdateProject) (*domain.Project, error) {
	return &domain.Project{ID: projectID, Name: "ERP"}, nil
}

func (f *fakeService) DeleteProject(ctx context.Context, userID domain.UserID, projectID domain.ProjectID) error {
	return nil
}

func (f *fakeService) CreateStage(ctx context.Context, userID domain.UserID, in domain.CreateStage) (*domain.Stage, error) {
	return &domain.Stage{ID: "bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb", ProjectID: in.ProjectID, Name: in.Name}, nil
}

func (f *fakeService) UpdateStage(ctx context.Context, userID domain.UserID, in domain.UpdateStage) (*domain.Stage, error) {
	return &domain.Stage{ID: in.StageID, ProjectID: in.ProjectID}, nil
}

func (f *fakeService) DeleteStage(ctx context.Context, userID domain.UserID, projectID domain.ProjectID, stageID domain.StageID) error {
	return nil
}

func (f *fakeService) ListResources(ctx context.Context) ([]domain.Resource, error) {
	return []domain.Resource{{ID: "cccccccc-cccc-cccc-cccc-cccccccccccc", Name: "Anna"}}, nil
}

func (f *fakeService) AddProjectResource(ctx context.Context, userID domain.UserID, in domain.AddProjectResource) (*domain.ProjectResource, error) {
	return &domain.ProjectResource{ProjectID: in.ProjectID, ResourceID: in.ResourceID}, nil
}

func (f *fakeService) ListProjectResources(ctx context.Context, userID domain.UserID, projectID domain.ProjectID) ([]domain.ProjectResource, error) {
	return []domain.ProjectResource{}, nil
}

func (f *fakeService) UpdateProjectResourceWorkload(ctx context.Context, userID domain.UserID, projectID domain.ProjectID, resourceID domain.ResourceID, workload float64) (*domain.ProjectResource, error) {
	return &domain.ProjectResource{ProjectID: projectID, ResourceID: resourceID, WorkloadPercentage: workload}, nil
}

func (f *fakeService) DeleteProjectResource(ctx context.Context, userID domain.UserID, projectID domain.ProjectID, resourceID domain.ResourceID) error {
	return nil
}

func (f *fakeService) GrantAccess(ctx context.Context, userID domain.UserID, projectID domain.ProjectID, access domain.ProjectAccess) error {
	return nil
}

func (f *fakeService) RevokeAccess(ctx context.Context, userID domain.UserID, projectID domain.ProjectID, target domain.UserID) error {
	return nil
}

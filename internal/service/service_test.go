package service

import (
	"context"
	"errors"
	"testing"

	"go.uber.org/mock/gomock"

	"project-budget-service/internal/domain"
)

func TestCreateProjectSetsOwnerAndCleansName(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := NewMockProjectRepository(ctrl)
	budgets := NewMockBudgetPublisher(ctrl)
	svc := NewProjectService(repo, budgets)
	userID := domain.UserID("11111111-1111-1111-1111-111111111111")
	expected := &domain.Project{ID: "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa", Name: "ERP"}

	repo.EXPECT().
		CreateProject(gomock.Any(), gomock.Any()).
		DoAndReturn(func(ctx context.Context, in domain.CreateProject) (*domain.Project, error) {
			if in.OwnerID != userID {
				t.Fatalf("owner id mismatch: %s", in.OwnerID)
			}
			if in.Name != "ERP" {
				t.Fatalf("name was not cleaned: %q", in.Name)
			}
			return expected, nil
		})

	got, err := svc.CreateProject(context.Background(), userID, domain.CreateProject{Name: "  ERP  "})
	if err != nil {
		t.Fatal(err)
	}
	if got != expected {
		t.Fatal("unexpected project pointer")
	}
}

func TestUpdateProjectRequiresEditorOrOwner(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := NewMockProjectRepository(ctrl)
	budgets := NewMockBudgetPublisher(ctrl)
	svc := NewProjectService(repo, budgets)

	repo.EXPECT().
		GetAccessRole(gomock.Any(), domain.ProjectID("aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa"), domain.UserID("11111111-1111-1111-1111-111111111111")).
		Return(domain.RoleReader, nil)

	_, err := svc.UpdateProject(
		context.Background(),
		"11111111-1111-1111-1111-111111111111",
		"aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa",
		domain.UpdateProject{},
	)
	var appErr *domain.AppError
	if !errors.As(err, &appErr) || appErr.Code != domain.CodeForbidden {
		t.Fatalf("expected forbidden error, got %v", err)
	}
}

func TestCreateStagePublishesBudgetJob(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := NewMockProjectRepository(ctrl)
	budgets := NewMockBudgetPublisher(ctrl)
	svc := NewProjectService(repo, budgets)
	start, _ := domain.ParseDate("2026-02-01")
	end, _ := domain.ParseDate("2026-02-03")
	projectID := domain.ProjectID("aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa")
	userID := domain.UserID("11111111-1111-1111-1111-111111111111")
	stage := &domain.Stage{ProjectID: projectID, Name: "Discovery", DurationDays: 3}

	repo.EXPECT().GetAccessRole(gomock.Any(), projectID, userID).Return(domain.RoleEditor, nil)
	repo.EXPECT().CreateStage(gomock.Any(), gomock.Any()).Return(stage, nil)
	budgets.EXPECT().Publish(gomock.Any(), projectID).Return(nil)

	got, err := svc.CreateStage(context.Background(), userID, domain.CreateStage{
		ProjectID: projectID,
		Name:      "Discovery",
		StartDate: start,
		EndDate:   end,
	})
	if err != nil {
		t.Fatal(err)
	}
	if got.DurationDays != 3 {
		t.Fatalf("unexpected duration: %d", got.DurationDays)
	}
}

func TestAddProjectResourceValidatesWorkload(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := NewMockProjectRepository(ctrl)
	budgets := NewMockBudgetPublisher(ctrl)
	svc := NewProjectService(repo, budgets)
	projectID := domain.ProjectID("aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa")
	userID := domain.UserID("11111111-1111-1111-1111-111111111111")

	repo.EXPECT().GetAccessRole(gomock.Any(), projectID, userID).Return(domain.RoleOwner, nil)

	_, err := svc.AddProjectResource(context.Background(), userID, domain.AddProjectResource{
		ProjectID:          projectID,
		ResourceID:         "bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb",
		WorkloadPercentage: 101,
	})
	var appErr *domain.AppError
	if !errors.As(err, &appErr) || appErr.Code != domain.CodeValidation {
		t.Fatalf("expected validation error, got %v", err)
	}
}

func TestRevokeOwnAccessIsRejected(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := NewMockProjectRepository(ctrl)
	budgets := NewMockBudgetPublisher(ctrl)
	svc := NewProjectService(repo, budgets)
	projectID := domain.ProjectID("aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa")
	userID := domain.UserID("11111111-1111-1111-1111-111111111111")

	repo.EXPECT().GetAccessRole(gomock.Any(), projectID, userID).Return(domain.RoleOwner, nil)

	err := svc.RevokeAccess(context.Background(), userID, projectID, userID)
	var appErr *domain.AppError
	if !errors.As(err, &appErr) || appErr.Code != domain.CodeValidation {
		t.Fatalf("expected validation error, got %v", err)
	}
}

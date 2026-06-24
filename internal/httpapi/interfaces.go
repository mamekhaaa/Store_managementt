package httpapi

import (
	"context"

	"project-budget-service/internal/domain"
)

type ProjectUseCase interface {
	CreateProject(context.Context, domain.UserID, domain.CreateProject) (*domain.Project, error)
	ListProjects(context.Context, domain.UserID) ([]domain.Project, error)
	GetProject(context.Context, domain.UserID, domain.ProjectID) (*domain.Project, error)
	UpdateProject(context.Context, domain.UserID, domain.ProjectID, domain.UpdateProject) (*domain.Project, error)
	DeleteProject(context.Context, domain.UserID, domain.ProjectID) error
	CreateStage(context.Context, domain.UserID, domain.CreateStage) (*domain.Stage, error)
	UpdateStage(context.Context, domain.UserID, domain.UpdateStage) (*domain.Stage, error)
	DeleteStage(context.Context, domain.UserID, domain.ProjectID, domain.StageID) error
	ListResources(context.Context) ([]domain.Resource, error)
	AddProjectResource(context.Context, domain.UserID, domain.AddProjectResource) (*domain.ProjectResource, error)
	ListProjectResources(context.Context, domain.UserID, domain.ProjectID) ([]domain.ProjectResource, error)
	UpdateProjectResourceWorkload(context.Context, domain.UserID, domain.ProjectID, domain.ResourceID, float64) (*domain.ProjectResource, error)
	DeleteProjectResource(context.Context, domain.UserID, domain.ProjectID, domain.ResourceID) error
	GrantAccess(context.Context, domain.UserID, domain.ProjectID, domain.ProjectAccess) error
	RevokeAccess(context.Context, domain.UserID, domain.ProjectID, domain.UserID) error
}

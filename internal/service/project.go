package service

import (
	"context"

	"project-budget-service/internal/domain"
)

type ProjectRepository interface {
	CreateProject(context.Context, domain.CreateProject) (*domain.Project, error)
	ListProjects(context.Context, domain.UserID) ([]domain.Project, error)
	GetProject(context.Context, domain.ProjectID) (*domain.Project, error)
	UpdateProject(context.Context, domain.ProjectID, domain.UpdateProject) (*domain.Project, error)
	DeleteProject(context.Context, domain.ProjectID) error
	CreateStage(context.Context, domain.CreateStage) (*domain.Stage, error)
	UpdateStage(context.Context, domain.UpdateStage) (*domain.Stage, error)
	DeleteStage(context.Context, domain.ProjectID, domain.StageID) error
	ListResources(context.Context) ([]domain.Resource, error)
	AddProjectResource(context.Context, domain.AddProjectResource) (*domain.ProjectResource, error)
	ListProjectResources(context.Context, domain.ProjectID) ([]domain.ProjectResource, error)
	UpdateProjectResourceWorkload(context.Context, domain.ProjectID, domain.ResourceID, float64) (*domain.ProjectResource, error)
	DeleteProjectResource(context.Context, domain.ProjectID, domain.ResourceID) error
	GrantAccess(context.Context, domain.ProjectID, domain.ProjectAccess) error
	RevokeAccess(context.Context, domain.ProjectID, domain.UserID) error
	GetAccessRole(context.Context, domain.ProjectID, domain.UserID) (domain.AccessRole, error)
}

type BudgetPublisher interface {
	Publish(context.Context, domain.ProjectID) error
}

type ProjectService struct {
	repo    ProjectRepository
	budgets BudgetPublisher
}

func NewProjectService(repo ProjectRepository, budgets BudgetPublisher) *ProjectService {
	return &ProjectService{repo: repo, budgets: budgets}
}

func (s *ProjectService) CreateProject(ctx context.Context, userID domain.UserID, in domain.CreateProject) (*domain.Project, error) {
	name, err := domain.CleanName(in.Name)
	if err != nil {
		return nil, err
	}
	in.Name = name
	in.OwnerID = userID
	return s.repo.CreateProject(ctx, in)
}

func (s *ProjectService) ListProjects(ctx context.Context, userID domain.UserID) ([]domain.Project, error) {
	return s.repo.ListProjects(ctx, userID)
}

func (s *ProjectService) GetProject(ctx context.Context, userID domain.UserID, projectID domain.ProjectID) (*domain.Project, error) {
	if err := s.requireRole(ctx, userID, projectID, domain.RoleReader, domain.RoleEditor, domain.RoleOwner); err != nil {
		return nil, err
	}
	return s.repo.GetProject(ctx, projectID)
}

func (s *ProjectService) UpdateProject(ctx context.Context, userID domain.UserID, projectID domain.ProjectID, in domain.UpdateProject) (*domain.Project, error) {
	if err := s.requireRole(ctx, userID, projectID, domain.RoleEditor, domain.RoleOwner); err != nil {
		return nil, err
	}
	if in.Name != nil {
		name, err := domain.CleanName(*in.Name)
		if err != nil {
			return nil, err
		}
		in.Name = &name
	}
	if in.Status != nil {
		if err := domain.ValidateStatus(*in.Status); err != nil {
			return nil, err
		}
	}
	return s.repo.UpdateProject(ctx, projectID, in)
}

func (s *ProjectService) DeleteProject(ctx context.Context, userID domain.UserID, projectID domain.ProjectID) error {
	if err := s.requireRole(ctx, userID, projectID, domain.RoleOwner); err != nil {
		return err
	}
	return s.repo.DeleteProject(ctx, projectID)
}

func (s *ProjectService) CreateStage(ctx context.Context, userID domain.UserID, in domain.CreateStage) (*domain.Stage, error) {
	if err := s.requireRole(ctx, userID, in.ProjectID, domain.RoleEditor, domain.RoleOwner); err != nil {
		return nil, err
	}
	name, err := domain.CleanName(in.Name)
	if err != nil {
		return nil, err
	}
	if _, err := domain.DurationDays(in.StartDate, in.EndDate); err != nil {
		return nil, err
	}
	in.Name = name
	stage, err := s.repo.CreateStage(ctx, in)
	if err != nil {
		return nil, err
	}
	return stage, s.budgets.Publish(ctx, in.ProjectID)
}

func (s *ProjectService) UpdateStage(ctx context.Context, userID domain.UserID, in domain.UpdateStage) (*domain.Stage, error) {
	if err := s.requireRole(ctx, userID, in.ProjectID, domain.RoleEditor, domain.RoleOwner); err != nil {
		return nil, err
	}
	if _, err := domain.DurationDays(in.StartDate, in.EndDate); err != nil {
		return nil, err
	}
	stage, err := s.repo.UpdateStage(ctx, in)
	if err != nil {
		return nil, err
	}
	return stage, s.budgets.Publish(ctx, in.ProjectID)
}

func (s *ProjectService) DeleteStage(ctx context.Context, userID domain.UserID, projectID domain.ProjectID, stageID domain.StageID) error {
	if err := s.requireRole(ctx, userID, projectID, domain.RoleEditor, domain.RoleOwner); err != nil {
		return err
	}
	if err := s.repo.DeleteStage(ctx, projectID, stageID); err != nil {
		return err
	}
	return s.budgets.Publish(ctx, projectID)
}

func (s *ProjectService) ListResources(ctx context.Context) ([]domain.Resource, error) {
	return s.repo.ListResources(ctx)
}

func (s *ProjectService) AddProjectResource(ctx context.Context, userID domain.UserID, in domain.AddProjectResource) (*domain.ProjectResource, error) {
	if err := s.requireRole(ctx, userID, in.ProjectID, domain.RoleEditor, domain.RoleOwner); err != nil {
		return nil, err
	}
	if err := domain.ValidateWorkload(in.WorkloadPercentage); err != nil {
		return nil, err
	}
	res, err := s.repo.AddProjectResource(ctx, in)
	if err != nil {
		return nil, err
	}
	return res, s.budgets.Publish(ctx, in.ProjectID)
}

func (s *ProjectService) ListProjectResources(ctx context.Context, userID domain.UserID, projectID domain.ProjectID) ([]domain.ProjectResource, error) {
	if err := s.requireRole(ctx, userID, projectID, domain.RoleReader, domain.RoleEditor, domain.RoleOwner); err != nil {
		return nil, err
	}
	return s.repo.ListProjectResources(ctx, projectID)
}

func (s *ProjectService) UpdateProjectResourceWorkload(ctx context.Context, userID domain.UserID, projectID domain.ProjectID, resourceID domain.ResourceID, workload float64) (*domain.ProjectResource, error) {
	if err := s.requireRole(ctx, userID, projectID, domain.RoleEditor, domain.RoleOwner); err != nil {
		return nil, err
	}
	if err := domain.ValidateWorkload(workload); err != nil {
		return nil, err
	}
	res, err := s.repo.UpdateProjectResourceWorkload(ctx, projectID, resourceID, workload)
	if err != nil {
		return nil, err
	}
	return res, s.budgets.Publish(ctx, projectID)
}

func (s *ProjectService) DeleteProjectResource(ctx context.Context, userID domain.UserID, projectID domain.ProjectID, resourceID domain.ResourceID) error {
	if err := s.requireRole(ctx, userID, projectID, domain.RoleEditor, domain.RoleOwner); err != nil {
		return err
	}
	if err := s.repo.DeleteProjectResource(ctx, projectID, resourceID); err != nil {
		return err
	}
	return s.budgets.Publish(ctx, projectID)
}

func (s *ProjectService) GrantAccess(ctx context.Context, userID domain.UserID, projectID domain.ProjectID, access domain.ProjectAccess) error {
	if err := s.requireRole(ctx, userID, projectID, domain.RoleOwner); err != nil {
		return err
	}
	if err := domain.ValidateRole(access.Role); err != nil {
		return err
	}
	return s.repo.GrantAccess(ctx, projectID, access)
}

func (s *ProjectService) RevokeAccess(ctx context.Context, userID domain.UserID, projectID domain.ProjectID, target domain.UserID) error {
	if err := s.requireRole(ctx, userID, projectID, domain.RoleOwner); err != nil {
		return err
	}
	if target == userID {
		return domain.NewValidationError("owner cannot revoke own access")
	}
	return s.repo.RevokeAccess(ctx, projectID, target)
}

func (s *ProjectService) requireRole(ctx context.Context, userID domain.UserID, projectID domain.ProjectID, allowed ...domain.AccessRole) error {
	role, err := s.repo.GetAccessRole(ctx, projectID, userID)
	if err != nil {
		return err
	}
	for _, item := range allowed {
		if role == item {
			return nil
		}
	}
	return domain.NewForbiddenError("insufficient project role")
}

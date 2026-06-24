package httpapi

import "project-budget-service/internal/domain"

type createProjectRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type updateProjectRequest struct {
	Name        *string `json:"name"`
	Description *string `json:"description"`
	Status      *string `json:"status"`
}

func (r updateProjectRequest) toDomain() domain.UpdateProject {
	var status *domain.ProjectStatus
	if r.Status != nil {
		value := domain.ProjectStatus(*r.Status)
		status = &value
	}
	return domain.UpdateProject{Name: r.Name, Description: r.Description, Status: status}
}

type stageRequest struct {
	Name      string      `json:"name"`
	StartDate domain.Date `json:"start_date"`
	EndDate   domain.Date `json:"end_date"`
}

type projectResourceRequest struct {
	ResourceID         string  `json:"resource_id"`
	WorkloadPercentage float64 `json:"workload_percentage"`
}

type workloadRequest struct {
	WorkloadPercentage float64 `json:"workload_percentage"`
}

type accessRequest struct {
	UserID string `json:"user_id"`
	Role   string `json:"role"`
}

type revokeAccessRequest struct {
	UserID string `json:"user_id"`
}

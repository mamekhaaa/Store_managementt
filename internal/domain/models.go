package domain

type Project struct {
	ID           ProjectID     `json:"id"`
	Name         string        `json:"name"`
	Description  string        `json:"description"`
	Status       ProjectStatus `json:"status"`
	TotalCost    float64       `json:"total_cost"`
	TotalRevenue float64       `json:"total_revenue"`
	Margin       float64       `json:"margin"`
}

type Stage struct {
	ID           StageID   `json:"id"`
	ProjectID    ProjectID `json:"project_id"`
	Name         string    `json:"name"`
	StartDate    Date      `json:"start_date"`
	EndDate      Date      `json:"end_date"`
	DurationDays int       `json:"duration_days"`
}

type Resource struct {
	ID                ResourceID `json:"id"`
	Name              string     `json:"name"`
	Role              string     `json:"role"`
	CostPerDay        float64    `json:"cost_per_day"`
	BillingRatePerDay float64    `json:"billing_rate_per_day"`
}

type ProjectResource struct {
	ID                         ProjectResourceID `json:"id"`
	ProjectID                  ProjectID         `json:"project_id"`
	ResourceID                 ResourceID        `json:"resource_id"`
	ResourceName               string            `json:"resource_name"`
	WorkloadPercentage         float64           `json:"workload_percentage"`
	EffectiveCostPerDay        float64           `json:"effective_cost_per_day"`
	EffectiveBillingRatePerDay float64           `json:"effective_billing_rate_per_day"`
}

type ProjectAccess struct {
	UserID UserID     `json:"user_id"`
	Role   AccessRole `json:"role"`
}

type CreateProject struct {
	Name        string
	Description string
	OwnerID     UserID
}

type UpdateProject struct {
	Name        *string
	Description *string
	Status      *ProjectStatus
}

type CreateStage struct {
	ProjectID ProjectID
	Name      string
	StartDate Date
	EndDate   Date
}

type UpdateStage struct {
	ProjectID ProjectID
	StageID   StageID
	StartDate Date
	EndDate   Date
}

type AddProjectResource struct {
	ProjectID          ProjectID
	ResourceID         ResourceID
	WorkloadPercentage float64
}

func ValidateWorkload(v float64) error {
	if v < 0 || v > 100 {
		return NewValidationError("workload_percentage must be between 0 and 100")
	}
	return nil
}

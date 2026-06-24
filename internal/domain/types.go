package domain

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/google/uuid"
)

type ProjectID string
type StageID string
type ResourceID string
type ProjectResourceID string
type UserID string

type ProjectStatus string
type AccessRole string

const (
	StatusPlanning  ProjectStatus = "planning"
	StatusActive    ProjectStatus = "active"
	StatusCompleted ProjectStatus = "completed"
	StatusArchived  ProjectStatus = "archived"

	RoleReader AccessRole = "reader"
	RoleEditor AccessRole = "editor"
	RoleOwner  AccessRole = "owner"
)

type Date time.Time

func ParseProjectID(v string) (ProjectID, error) {
	if _, err := uuid.Parse(v); err != nil {
		return "", NewValidationError("invalid project id")
	}
	return ProjectID(v), nil
}

func ParseStageID(v string) (StageID, error) {
	if _, err := uuid.Parse(v); err != nil {
		return "", NewValidationError("invalid stage id")
	}
	return StageID(v), nil
}

func ParseResourceID(v string) (ResourceID, error) {
	if _, err := uuid.Parse(v); err != nil {
		return "", NewValidationError("invalid resource id")
	}
	return ResourceID(v), nil
}

func ParseUserID(v string) (UserID, error) {
	if _, err := uuid.Parse(v); err != nil {
		return "", NewValidationError("invalid user id")
	}
	return UserID(v), nil
}

func NewProjectID() ProjectID { return ProjectID(uuid.NewString()) }
func NewStageID() StageID     { return StageID(uuid.NewString()) }
func NewProjectResourceID() ProjectResourceID {
	return ProjectResourceID(uuid.NewString())
}

func NewDate(t time.Time) Date {
	y, m, d := t.Date()
	return Date(time.Date(y, m, d, 0, 0, 0, 0, time.UTC))
}

func ParseDate(v string) (Date, error) {
	t, err := time.Parse("2006-01-02", v)
	if err != nil {
		return Date{}, NewValidationError("date must use YYYY-MM-DD format")
	}
	return NewDate(t), nil
}

func (d Date) Time() time.Time { return time.Time(d) }

func (d Date) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.Time().Format("2006-01-02"))
}

func (d *Date) UnmarshalJSON(b []byte) error {
	var v string
	if err := json.Unmarshal(b, &v); err != nil {
		return err
	}
	parsed, err := ParseDate(v)
	if err != nil {
		return err
	}
	*d = parsed
	return nil
}

func DurationDays(start, end Date) (int, error) {
	if end.Time().Before(start.Time()) {
		return 0, NewValidationError("end_date must be greater than or equal to start_date")
	}
	return int(end.Time().Sub(start.Time()).Hours()/24) + 1, nil
}

func ValidateStatus(v ProjectStatus) error {
	switch v {
	case StatusPlanning, StatusActive, StatusCompleted, StatusArchived:
		return nil
	default:
		return NewValidationError("unknown project status")
	}
}

func ValidateRole(v AccessRole) error {
	switch v {
	case RoleReader, RoleEditor, RoleOwner:
		return nil
	default:
		return NewValidationError("unknown access role")
	}
}

func CleanName(v string) (string, error) {
	v = strings.TrimSpace(v)
	if len(v) < 2 || len(v) > 120 {
		return "", NewValidationError("name length must be between 2 and 120")
	}
	return v, nil
}

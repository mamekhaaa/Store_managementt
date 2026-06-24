package background

import (
	"context"
	"log/slog"

	"golang.org/x/sync/errgroup"

	"project-budget-service/internal/domain"
)

type budgetRepository interface {
	RecalculateProjectBudget(context.Context, domain.ProjectID) error
}

type BudgetQueue struct {
	repo   budgetRepository
	log    *slog.Logger
	jobs   chan domain.ProjectID
	cancel context.CancelFunc
	group  *errgroup.Group
}

func NewBudgetQueue(parent context.Context, repo budgetRepository, log *slog.Logger) *BudgetQueue {
	ctx, cancel := context.WithCancel(parent)
	group, ctx := errgroup.WithContext(ctx)
	q := &BudgetQueue{
		repo:   repo,
		log:    log,
		jobs:   make(chan domain.ProjectID, 64),
		cancel: cancel,
		group:  group,
	}
	group.Go(func() error {
		for {
			select {
			case <-ctx.Done():
				return nil
			case projectID := <-q.jobs:
				if err := q.repo.RecalculateProjectBudget(ctx, projectID); err != nil {
					log.Error("budget recalculation failed", "project_id", projectID, "error", err)
				}
			}
		}
	})
	return q
}

func (q *BudgetQueue) Publish(ctx context.Context, projectID domain.ProjectID) error {
	select {
	case q.jobs <- projectID:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (q *BudgetQueue) Shutdown(ctx context.Context) error {
	q.cancel()
	done := make(chan error, 1)
	go func() { done <- q.group.Wait() }()
	select {
	case err := <-done:
		return err
	case <-ctx.Done():
		return ctx.Err()
	}
}

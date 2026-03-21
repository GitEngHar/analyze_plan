package summary

import (
	"context"
	"terraform-sammary/internal/domain/repository"
)

type Confirm interface {
	Execute(ctx context.Context, filePath string) (error, *ResultResourceAddress)
}

type ConfirmImpl struct {
	repo repository.PlanRepository
}

func NewConfirm(repo repository.PlanRepository) Confirm {
	return ConfirmImpl{
		repo: repo,
	}
}

func (c ConfirmImpl) Execute(ctx context.Context, filePath string) (error, *ResultResourceAddress) {
	return nil, nil
}

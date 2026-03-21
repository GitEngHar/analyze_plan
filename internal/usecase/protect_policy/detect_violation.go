package protect_policy

import (
	"context"
	"terraform-sammary/internal/domain/repository"
)

type DetectViolation interface {
	Execute(ctx context.Context, filePath string) error
}

type DetectViolationImpl struct {
	repo repository.PolicyRepository
}

func NewDetectViolation(repo repository.PolicyRepository) DetectViolation {
	return &DetectViolationImpl{
		repo: repo,
	}
}

func (v *DetectViolationImpl) Execute(ctx context.Context, filePath string) error {
	return nil
}

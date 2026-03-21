package protect_replace

import (
	"context"
	"terraform-sammary/internal/usecase/summary"
)

type Detect interface {
	Execute(ctx context.Context, plan summary.ResultResourceAddress) error
}

type DetectImpl struct {
}

func NewDetect() Detect {
	return &DetectImpl{}
}

func (d *DetectImpl) Execute(ctx context.Context, plan summary.ResultResourceAddress) error {
	return nil
}

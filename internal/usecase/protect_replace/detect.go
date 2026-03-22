package protect_replace

import (
	"context"
	"fmt"
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

func (d *DetectImpl) Execute(_ context.Context, plan summary.ResultResourceAddress) error {
	var isReplaceWarning = len(plan.Replace) > 0
	fmt.Println("\nReplace Detected\n----------------")
	// replace warn
	if isReplaceWarning {
		for _, address := range plan.Replace {
			fmt.Println("+/- ", address)
		}
	}
	return nil
}

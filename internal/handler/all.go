package handler

import (
	"context"
	"fmt"
	"os"
	"terraform-sammary/internal/infra/repository/json"
	"terraform-sammary/internal/infra/repository/yaml"
	"terraform-sammary/internal/usecase/protect_policy"
	"terraform-sammary/internal/usecase/protect_replace"
	"terraform-sammary/internal/usecase/summary"
)

// As tha project expand,responsibilities more distributed
var isUsePolicy bool

func NewActionsHandler() ActionsHandler {
	return ActionsHandler{}
}

type ActionsHandler struct {
}

func (h ActionsHandler) Handle(args []string) error {
	var (
		planFilePath  string
		policyDirPath string
	)
	switch len(args) {
	case 2:
		fmt.Printf("planFilePath: %s \n", os.Args[1])
		break
	case 3:
		fmt.Printf("planFilePath: %s \n, policyFilePath: %s\n", os.Args[1], os.Args[2])
		isUsePolicy = true
		break
	default:
		return fmt.Errorf("usage: go run analyze_plan.go plan.json policy_path")
	}
	ctx := context.Background()
	planRepo := json.NewPlanJSONRepository(planFilePath)
	summaryUc := summary.NewConfirm(planRepo)
	err, resultResourceAddress := summaryUc.Execute(ctx, planFilePath)
	if err != nil {
		return err
	}
	replaceDetectUc := protect_replace.NewDetect()
	err = replaceDetectUc.Execute(ctx, *resultResourceAddress)
	if err != nil {
		return err
	}
	if isUsePolicy {
		policyRepo := yaml.NewPolicyYAMLRepository(policyDirPath)
		detectViolationUc := protect_policy.NewDetectViolation(policyRepo)
		err = detectViolationUc.Execute(ctx, policyDirPath)
		if err != nil {
			return err
		}
	}
	return nil
}

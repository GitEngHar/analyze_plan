package protect_policy

import (
	"context"
	"fmt"
	"strings"
	"terraform-sammary/internal/domain/repository"
	"terraform-sammary/internal/usecase/summary"
	"time"
)

type DetectViolation interface {
	Execute(ctx context.Context, resourceTypeToNames summary.ResultResourceTypeToNames, dirPath string) error
}

type DetectViolationImpl struct {
	repo repository.PolicyRepository
}

func NewDetectViolation(repo repository.PolicyRepository) DetectViolation {
	return &DetectViolationImpl{
		repo: repo,
	}
}

func (v *DetectViolationImpl) Execute(ctx context.Context, resourceTypeToNames summary.ResultResourceTypeToNames, dirPath string) error {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	policies, err := v.repo.ReadPolicies(dirPath)
	if err != nil {
		return err
	}
	for _, deleteProtectRule := range policies.DeleteProtectPolicy.ProtectRules {
		deleteResourceNames := resourceTypeToNames.Delete[deleteProtectRule.Resource]
		if deleteResourceNames != nil {
			for _, deleteResourceName := range deleteResourceNames {
				outputPolicyViolation(deleteProtectRule.Severity, fmt.Sprintf("%s %s", "+/-", deleteProtectRule.Resource), deleteResourceName)
			}
		}
	}
	for _, replaceProtectRule := range policies.ReplaceProtectPolicy.ProtectRules {
		replaceResourceNames := resourceTypeToNames.Replace[replaceProtectRule.Resource]
		if replaceResourceNames != nil {
			for _, name := range replaceResourceNames {
				outputPolicyViolation(replaceProtectRule.Severity, fmt.Sprintf("%s %s", "+/-", replaceProtectRule.Resource), name)
			}
		}
	}
	return nil
}

func outputPolicyViolation(severity string, resourceType, resourceName string) {
	switch strings.ToLower(severity) {
	case "critical":
		fmt.Printf("🚨 %s.%s \n", resourceType, resourceName)
		break
	case "warning", "warn":
		fmt.Printf("⚠️ %s.%s \n", resourceType, resourceName)
		break
	}
}

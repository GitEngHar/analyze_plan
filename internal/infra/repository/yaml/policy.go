package yaml

import (
	"sync"
	"terraform-sammary/internal/domain/repository"
	dPolicy "terraform-sammary/internal/domain/value-object/policy"
)

type PolicyYAMLRepositoryImpl struct {
	filePath string
	mu       sync.Mutex
}

func NewPolicyYAMLRepository(filePath string) repository.PolicyRepository {
	return &PolicyYAMLRepositoryImpl{
		filePath: filePath,
	}
}

func (y *PolicyYAMLRepositoryImpl) ReadPolicy() (dPolicy.ResourceProtectPolicies, error) {
	var policy dPolicy.ResourceProtectPolicies
	// TODO add details
	return policy, nil
}

package yaml

import (
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
	"sync"
	"terraform-sammary/internal/domain/repository"
	dPolicy "terraform-sammary/internal/domain/value-object/policy"
)

const (
	deletePolicyYamlName  = "delete.yaml"
	replacePolicyYamlName = "replace.yaml"
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

func (y *PolicyYAMLRepositoryImpl) ReadPolicies(dirPath string) (*dPolicy.ResourceProtectPolicies, error) {
	var policies dPolicy.ResourceProtectPolicies
	files, err := os.ReadDir(dirPath)
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}
		path := filepath.Join(dirPath, file.Name())
		switch file.Name() {
		case deletePolicyYamlName:
			protectPolicy, er := loadProtectPolicy(path)
			if er != nil {
				return nil, er
			}
			policies.DeleteProtectPolicy = *protectPolicy
		case replacePolicyYamlName:
			protectPolicy, er := loadProtectPolicy(path)
			if er != nil {
				return nil, er
			}
			policies.ReplaceProtectPolicy = *protectPolicy
		default:
			continue
		}
	}
	return &policies, nil
}

func loadProtectPolicy(path string) (*dPolicy.ResourceProtectPolicy, error) {
	data, er := os.ReadFile(path)
	if er != nil {
		return nil, er
	}
	var policy *dPolicy.ResourceProtectPolicy
	err := yaml.Unmarshal(data, &policy)
	if err != nil {
		return nil, err
	}
	return policy, nil
}

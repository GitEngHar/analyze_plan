package json

import (
	"encoding/json"
	"os"
	"sync"
	"terraform-sammary/internal/domain/repository"
	dPlan "terraform-sammary/internal/domain/value-object/plan"
)

type PlanJSONRepositoryImpl struct {
	filePath string
	mu       sync.Mutex
}

func NewPlanJSONRepository(filePath string) repository.PlanRepository {
	return &PlanJSONRepositoryImpl{
		filePath: filePath,
	}
}

func (p *PlanJSONRepositoryImpl) ReadPlan() (dPlan.Plan, error) {
	var plan dPlan.Plan
	p.mu.Lock()
	defer p.mu.Unlock()
	data, err := os.ReadFile(p.filePath)
	if err != nil {
		return dPlan.Plan{}, err
	}
	err = json.Unmarshal(data, &plan)
	if err != nil {
		return dPlan.Plan{}, err
	}
	return plan, nil
}

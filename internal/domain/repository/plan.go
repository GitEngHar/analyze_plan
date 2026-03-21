package repository

import dPlan "terraform-sammary/internal/domain/value-object/plan"

type PlanRepository interface {
	ReadPlan() (dPlan.Plan, error)
}

package repository

import dPlan "terraform-sammary/internal/domain/value-object/plan"

type PlanJSONRepository interface {
	ReadPlan() (dPlan.Plan, error)
}

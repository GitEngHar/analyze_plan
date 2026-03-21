package repository

import (
	dPolicy "terraform-sammary/internal/domain/value-object/policy"
)

type PolicyRepository interface {
	ReadPolicy() (dPolicy.ResourceProtectPolicies, error)
}

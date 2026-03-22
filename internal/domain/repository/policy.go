package repository

import (
	dPolicy "terraform-sammary/internal/domain/value-object/policy"
)

type PolicyRepository interface {
	ReadPolicies() (*dPolicy.ResourceProtectPolicies, error)
}

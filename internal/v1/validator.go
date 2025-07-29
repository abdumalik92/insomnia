package v1

import (
	"github.com/alifcapital/keycloak_module/src/iam"
	"github.com/go-playground/validator/v10"
)

type IRuleContainer interface {
	Rules() map[string]string
}

func NewValidator() *validator.Validate {
	v := validator.New()

	registerRules(
		v,
		iam.CreateUserRequest{},
		iam.UpdateUserRequest{},
	)

	return v
}

func registerRules(v *validator.Validate, requests ...IRuleContainer) {
	for _, r := range requests {
		v.RegisterStructValidationMapRules(r.Rules(), r)
	}
}

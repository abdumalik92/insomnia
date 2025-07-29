package controller

import (
	"github.com/alifcapital/keycloak_module/conf"
	v1 "github.com/alifcapital/keycloak_module/internal/v1"
	"github.com/alifcapital/keycloak_module/src/iam"
	"go.uber.org/fx"
)

type Controller struct {
	fx.In

	Cfg            *conf.Config
	IAMService     *iam.Service
	KeycloakClient *v1.KeycloakHTTPClient
	Publisher      *v1.Publisher
}

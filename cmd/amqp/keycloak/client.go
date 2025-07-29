package keycloak

import (
	"github.com/alifcapital/keycloak_module/cmd/amqp/controller"
	"github.com/alifcapital/keycloak_module/conf"
	"github.com/alifcapital/keycloak_module/internal/utils"
	"github.com/alifcapital/rabbitmq/mqutils"
	"go.uber.org/fx"
)

type Client struct {
	fx.In

	Cfg                   *conf.Config
	Pool                  *mqutils.Pool
	Controller            controller.Controller
	ErrorHandler          mqutils.ErrorHandler
	PanicRecoveryCallback mqutils.PanicRecoveryCallback
}

func (c Client) Start() error {
	// 1. create rabbitmq.Client
	clientConfig := utils.NewLocalRabbitMQClientConfig(c.Cfg)
	client, err := c.Pool.Register(clientConfig)
	if err != nil {
		return err
	}

	// 2. register all consumers
	deadLetterConsumer := c.deadLetterConsumer()
	if err := client.Consume(deadLetterConsumer); err != nil {
		return err
	}
	keycloakEventsConsumer := c.keycloakEvents()
	return client.Consume(keycloakEventsConsumer)
}

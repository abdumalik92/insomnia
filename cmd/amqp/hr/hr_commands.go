package hr

import (
	"github.com/alifcapital/keycloak_module/internal/utils"
	"github.com/alifcapital/rabbitmq"
	"github.com/alifcapital/rabbitmq/mqutils"
)

const (
	eventCreateUser = "create_user"
	eventUpdateUser = "update_user"
)

func (c Client) hrCommands() rabbitmq.AMQPConsumer {
	// 1. create router
	mids := []mqutils.Middleware{
		mqutils.NewPanicRecoveryMiddleware(c.PanicRecoveryCallback),
		mqutils.NewTracerMiddleware(),
	}
	router := mqutils.NewRouter(c.ErrorHandler, mids...)

	// 2. register event routing
	ctrl := c.Controller
	router.RegisterEventHandler(eventCreateUser, ctrl.CreatUser)
	router.RegisterEventHandler(eventUpdateUser, ctrl.UpdateUser)

	queueName := utils.RabbitMQueueName(c.Cfg, "hr_commands")
	return utils.NewAMQPConsumer(
		c.Cfg.NewRabbitMQPrimaryExchange,
		queueName,
		c.Cfg.RabbitMQDeadLetterExchange,
		true,
		router.GetEventNames(),
		router,
	)
}

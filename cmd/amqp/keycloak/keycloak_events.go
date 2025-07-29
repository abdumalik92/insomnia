package keycloak

import (
	"strings"

	"github.com/alifcapital/keycloak_module/internal/utils"
	"github.com/alifcapital/rabbitmq"
	"github.com/alifcapital/rabbitmq/mqutils"
)

const (
	eventUserActionsInKK    = "KK.EVENT.ADMIN.{REALM_ID}.SUCCESS.USER.*"
	eventProfileActionsInKK = "KK.EVENT.CLIENT.{REALM_ID}.SUCCESS.account.UPDATE_PROFILE"
)

func (c Client) keycloakEvents() rabbitmq.AMQPConsumer {
	// 1. create router
	mids := []mqutils.Middleware{
		mqutils.NewPanicRecoveryMiddleware(c.PanicRecoveryCallback),
		mqutils.NewTracerMiddleware(),
		utils.NewConsumerTraceLoggerMid(),
	}
	router := mqutils.NewRouter(c.ErrorHandler, mids...)

	// 2. register event routing
	keycloakUserTopics := strings.ReplaceAll(eventUserActionsInKK, "{REALM_ID}", c.Cfg.KeycloakRealm)
	keycloakProfileTopics := strings.ReplaceAll(eventProfileActionsInKK, "{REALM_ID}", c.Cfg.KeycloakRealm)

	keycloakUserCreated := strings.ReplaceAll(keycloakUserTopics, "*", "CREATE")
	keycloakUserUpdated := strings.ReplaceAll(keycloakUserTopics, "*", "UPDATE")
	//keycloakUserDeleted := strings.ReplaceAll(keycloakUserTopics, "*", "DELETE")

	ctrl := c.Controller
	router.RegisterEventHandler(keycloakUserCreated, ctrl.BroadcastUserCreated)
	router.RegisterEventHandler(keycloakUserUpdated, ctrl.BroadcastUserUpdated)
	//router.RegisterEventHandler(keycloakUserDeleted, ctrl.BroadcastUserDeleted)
	router.RegisterEventHandler(keycloakProfileTopics, ctrl.BroadcastUserChangedThroughProfile)

	queueName := utils.RabbitMQueueName(c.Cfg, "keycloak_events")
	return utils.NewAMQPConsumer(
		c.Cfg.KeycloakRabbitMQPrimaryExchange,
		queueName,
		c.Cfg.RabbitMQDeadLetterExchange,
		true,
		router.GetEventNames(),
		router,
	)
}

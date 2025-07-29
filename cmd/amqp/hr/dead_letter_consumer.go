package hr

import (
	"github.com/alifcapital/keycloak_module/internal/utils"
	"github.com/alifcapital/rabbitmq"
)

func (c Client) deadLetterConsumer() rabbitmq.AMQPConsumer {
	dlq := utils.RabbitMQueueName(c.Cfg, "poisoned_msg")
	return utils.NewAMQPConsumer(
		c.Cfg.RabbitMQDeadLetterExchange,
		dlq,
		"",
		true,
		[]string{"#"},
		nil,
	)
}

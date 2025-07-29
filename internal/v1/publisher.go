package v1

import (
	"context"
	"encoding/json"
	"time"

	"github.com/alifcapital/keycloak_module/conf"
	"github.com/alifcapital/keycloak_module/internal/utils"
	"github.com/alifcapital/rabbitmq"
	"github.com/alifcapital/rabbitmq/mqutils"
)

type Publisher struct {
	cfg           *conf.Config
	client        *rabbitmq.Client
	testEnvClient *rabbitmq.Client
}

func NewPublisher(cfg *conf.Config, pool *mqutils.Pool) (*Publisher, error) {
	clientCfg, err := utils.NewRabbitMQClientConfig(cfg, cfg.NewRabbitMQPrimaryVHost)
	if err != nil {
		return nil, err
	}
	client, err := pool.Register(clientCfg)
	if err != nil {
		return nil, err
	}

	testEnvCfg := utils.NewTestEnvRabbitMQClientConfig(cfg)
	testEnvClient, err := pool.Register(testEnvCfg)
	if err != nil {
		return nil, err
	}

	return &Publisher{
		cfg:           cfg,
		client:        client,
		testEnvClient: testEnvClient,
	}, nil
}

func (p *Publisher) Publish(ctx context.Context, eventName string, obj any) error {
	body, err := json.Marshal(obj)
	if err != nil {
		return err
	}

	timeOutCtx, cancel := context.WithTimeout(ctx, time.Second*15)
	defer cancel()

	msg := mqutils.NewMessage("", body)
	if err := mqutils.Publish(timeOutCtx, p.cfg.NewRabbitMQPrimaryExchange, eventName, msg, p.client); err != nil {
		return err
	}

	return mqutils.Publish(timeOutCtx, p.cfg.TestRabbitMQPrimaryExchange, eventName, msg, p.testEnvClient)
}

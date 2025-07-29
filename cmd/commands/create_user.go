package main

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"time"

	"github.com/alifcapital/keycloak_module/conf"
	"github.com/alifcapital/keycloak_module/internal/utils"
	"github.com/alifcapital/keycloak_module/src/iam"
	"github.com/alifcapital/rabbitmq"
	"github.com/alifcapital/rabbitmq/mqutils"
)

func CreateUser() {
	v, err := conf.NewViper()
	if err != nil {
		panic(err)
	}
	cfg, err := conf.NewConfig(v)
	if err != nil {
		panic(err)
	}

	clientCfg, err := utils.NewRabbitMQClientConfig(cfg, cfg.NewRabbitMQPrimaryVHost)
	if err != nil {
		panic(err)
	}
	client, err := rabbitmq.NewClient(clientCfg)
	if err != nil {
		panic(err)
	}
	defer client.Close()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	msg := mqutils.NewMessage("", randUser())
	err = client.Publish(
		ctx,
		cfg.NewRabbitMQPrimaryExchange,
		"create_user",
		false,
		false,
		msg,
	)
	if err != nil {
		panic(err)
	}
}

func randUser() []byte {
	req := iam.CreateUserRequest{
		Username:      fmt.Sprintf("john_doe_%d", rand.Int()),
		Email:         fmt.Sprintf("john%d@mail.com", rand.Int()),
		FirstName:     "John",
		LastName:      "Doe",
		Enabled:       true,
		EmailVerified: true,
		Attributes:    nil,
	}
	body, err := json.Marshal(&req)
	if err != nil {
		panic(err)
	}
	return body
}

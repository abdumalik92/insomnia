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

func UpdateUser() {
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

	msg := mqutils.NewMessage("", userUpd("6927e3d8-b7e6-4e0f-92e3-ff18e20359b5"))
	err = client.Publish(
		ctx,
		cfg.NewRabbitMQPrimaryExchange,
		"update_user",
		false,
		false,
		msg,
	)
	if err != nil {
		panic(err)
	}
}

func userUpd(id string) []byte {
	req := iam.UpdateUserRequest{
		ID:            id,
		Username:      fmt.Sprintf("john_doe_%d", rand.Int()),
		Email:         fmt.Sprintf("MINN%d@mail.com", rand.Int()),
		FirstName:     "BONN",
		LastName:      "SOMM",
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

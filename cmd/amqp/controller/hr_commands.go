package controller

import (
	"context"
	"errors"
	"fmt"

	"github.com/alifcapital/keycloak_module/src/iam"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	amqp "github.com/rabbitmq/amqp091-go"
)

func (c Controller) CreatUser(ctx context.Context, msg amqp.Delivery) error {
	span, newCtx := opentracing.StartSpanFromContext(ctx, "ctrl.CreateUser")
	defer span.Finish()

	req := new(iam.CreateUserRequest)
	if err := unmarshal(msg.Body, req); err != nil {
		return errors.Join(fmt.Errorf("marshal failed for: %s", string(msg.Body)), err)
	}

	user, err := c.IAMService.CreateUser(newCtx, req)
	if err != nil {
		return err
	}

	span.LogFields(log.Bool("user_created", true))

	buff := map[string]any{}
	if err := unmarshal(msg.Body, &buff); err != nil {
		return errors.Join(fmt.Errorf("marshal failed for: %s", string(msg.Body)), err)
	}
	initialPass, _ := buff["initial_password"].(string)
	return c.IAMService.SetTemporaryPassword(newCtx, &iam.SetTemporaryPasswordRequest{
		UserID:            user.ID,
		TemporaryPassword: initialPass,
	})
}

func (c Controller) UpdateUser(ctx context.Context, msg amqp.Delivery) error {
	span, newCtx := opentracing.StartSpanFromContext(ctx, "ctrl.UpdateUser")
	defer span.Finish()

	req := new(iam.UpdateUserRequest)
	if err := unmarshal(msg.Body, req); err != nil {
		return errors.Join(fmt.Errorf("marshal failed for: %s", string(msg.Body)), err)
	}
	return c.IAMService.UpdateUser(newCtx, req)
}

package controller

import (
	"context"

	"github.com/opentracing/opentracing-go"
	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	// application publishes
	eventUserCreated = "keycloak.user_created"
	eventUserUpdated = "keycloak.user_updated"
	//eventUserDeleted = "keycloak.user_deleted"
)

type resourceUpdatedInKeycloak struct {
	ResourceType  string `json:"resourceType"`
	OperationType string `json:"operationType"`
	ResourcePath  string `json:"resourcePath"`
}

func (c Controller) BroadcastUserCreated(ctx context.Context, msg amqp.Delivery) error {
	span, newCtx := opentracing.StartSpanFromContext(ctx, "ctrl.BroadcastUserCreated")
	defer span.Finish()

	payload := new(resourceUpdatedInKeycloak)
	if err := unmarshal(msg.Body, payload); err != nil {
		return err
	}

	userID, err := extractUserId(payload.ResourcePath)
	if err != nil {
		return err
	}

	user, err := c.KeycloakClient.Get(newCtx, userID)
	if err != nil {
		return err
	}

	return c.Publisher.Publish(newCtx, eventUserCreated, user)
}

func (c Controller) BroadcastUserUpdated(ctx context.Context, msg amqp.Delivery) error {
	span, newCtx := opentracing.StartSpanFromContext(ctx, "ctrl.BroadcastUserUpdated")
	defer span.Finish()

	payload := new(resourceUpdatedInKeycloak)
	if err := unmarshal(msg.Body, payload); err != nil {
		return err
	}

	userID, err := extractUserId(payload.ResourcePath)
	if err != nil {
		return err
	}

	user, err := c.KeycloakClient.Get(newCtx, userID)
	if err != nil {
		return err
	}

	return c.Publisher.Publish(newCtx, eventUserUpdated, user)
}

//func (c Controller) BroadcastUserDeleted(ctx context.Context, msg amqp.Delivery) error {
//	span, newCtx := opentracing.StartSpanFromContext(ctx, "ctrl.BroadcastUserDeleted")
//	defer span.Finish()
//
//	payload := new(resourceUpdatedInKeycloak)
//	if err := unmarshal(msg.Body, payload); err != nil {
//		return err
//	}
//
//	userID, err := extractUserId(payload.ResourcePath)
//	if err != nil {
//		return err
//	}
//
//	return c.Publisher.Publish(newCtx, eventUserDeleted, userID)
//}

type profileChangedInKeycloak struct {
	UserID string `json:"userId"`
}

func (c Controller) BroadcastUserChangedThroughProfile(ctx context.Context, msg amqp.Delivery) error {
	span, newCtx := opentracing.StartSpanFromContext(ctx, "ctrl.BroadcastUserChangedThroughProfile")
	defer span.Finish()

	payload := new(profileChangedInKeycloak)
	if err := unmarshal(msg.Body, payload); err != nil {
		return err
	}

	user, err := c.KeycloakClient.Get(newCtx, payload.UserID)
	if err != nil {
		return err
	}
	return c.Publisher.Publish(newCtx, eventUserUpdated, user)
}

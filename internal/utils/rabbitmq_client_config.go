package utils

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/alifcapital/keycloak_module/conf"
	"github.com/alifcapital/rabbitmq"
	"github.com/alifcapital/rabbitmq/mqutils"
	"github.com/getsentry/sentry-go"
	"github.com/google/uuid"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	amqp "github.com/rabbitmq/amqp091-go"
)

func NewRabbitMQClientConfig(cfg *conf.Config, vhost string) (rabbitmq.ClientConfig, error) {
	tlsConfig, err := getTlsConfig(cfg)
	if err != nil {
		return rabbitmq.ClientConfig{}, err
	}
	return rabbitmq.ClientConfig{
		NetworkErrCallback: func(a *amqp.Error) {
			err := fmt.Errorf("network issues with rabbitmq, error: %s, vhost: %s", a.Error(), vhost)
			sentry.CaptureException(err)
		},

		AutoRecoveryInterval: time.Second * 10,

		AutoRecoveryErrCallback: func(err error) bool {
			joinedErr := errors.Join(err, fmt.Errorf("connection autorecovery failed"))
			sentry.CaptureException(joinedErr)
			return true // always try to reconnect
		},

		ConsumerAutoRecoveryErrCallback: func(consumer rabbitmq.AMQPConsumer, err error) {
			joinedErr := errors.Join(err, fmt.Errorf("consumer autorecovery failed"))
			sentry.CaptureException(joinedErr)
		},

		DialConfig: rabbitmq.DialConfig{
			User:     cfg.NewRabbitMQUser,
			Password: cfg.NewRabbitMQPassword,
			Host:     cfg.NewRabbitMQHost,
			Port:     cfg.NewRabbitMQPort,
			AMQPConfig: amqp.Config{
				Vhost:           vhost,
				TLSClientConfig: tlsConfig,
			},
		},

		PublisherConfirmEnabled: true,
		PublisherConfirmNowait:  false,
		ConsumerQos:             0,
		ConsumerPrefetchSize:    0,
		ConsumerGlobal:          false,
	}, nil
}

func NewLocalRabbitMQClientConfig(cfg *conf.Config) rabbitmq.ClientConfig {
	return rabbitmq.ClientConfig{
		NetworkErrCallback: func(a *amqp.Error) {
			err := fmt.Errorf("network issues with rabbitmq, error: %s, vhost: %s", a.Error(), cfg.KeycloakRabbitMQVHost)
			sentry.CaptureException(err)
		},

		AutoRecoveryInterval: time.Second * 10,

		AutoRecoveryErrCallback: func(err error) bool {
			joinedErr := errors.Join(err, fmt.Errorf("connection autorecovery failed"))
			sentry.CaptureException(joinedErr)
			return true // always try to reconnect
		},

		ConsumerAutoRecoveryErrCallback: func(consumer rabbitmq.AMQPConsumer, err error) {
			joinedErr := errors.Join(err, fmt.Errorf("consumer autorecovery failed"))
			sentry.CaptureException(joinedErr)
		},

		DialConfig: rabbitmq.DialConfig{
			User:     cfg.KeycloakRabbitMQUser,
			Password: cfg.KeycloakRabbitMQPassword,
			Host:     cfg.KeycloakRabbitMQHost,
			Port:     cfg.KeycloakRabbitMQPort,
			AMQPConfig: amqp.Config{
				Vhost: cfg.KeycloakRabbitMQVHost,
			},
		},

		PublisherConfirmEnabled: true,
		PublisherConfirmNowait:  false,
		ConsumerQos:             0,
		ConsumerPrefetchSize:    0,
		ConsumerGlobal:          false,
	}
}

func NewTestEnvRabbitMQClientConfig(cfg *conf.Config) rabbitmq.ClientConfig {
	return rabbitmq.ClientConfig{
		NetworkErrCallback: func(a *amqp.Error) {
			err := fmt.Errorf("network issues with rabbitmq, error: %s, vhost: %s", a.Error(), cfg.KeycloakRabbitMQVHost)
			sentry.CaptureException(err)
		},

		AutoRecoveryInterval: time.Second * 10,

		AutoRecoveryErrCallback: func(err error) bool {
			joinedErr := errors.Join(err, fmt.Errorf("connection autorecovery failed"))
			sentry.CaptureException(joinedErr)
			return true // always try to reconnect
		},

		ConsumerAutoRecoveryErrCallback: func(consumer rabbitmq.AMQPConsumer, err error) {
			joinedErr := errors.Join(err, fmt.Errorf("consumer autorecovery failed"))
			sentry.CaptureException(joinedErr)
		},

		DialConfig: rabbitmq.DialConfig{
			User:     cfg.TestRabbitMQUser,
			Password: cfg.TestRabbitMQPassword,
			Host:     cfg.TestRabbitMQHost,
			Port:     cfg.TestRabbitMQPort,
			AMQPConfig: amqp.Config{
				Vhost: cfg.TestRabbitMQVHost,
			},
		},

		PublisherConfirmEnabled: true,
		PublisherConfirmNowait:  false,
		ConsumerQos:             0,
		ConsumerPrefetchSize:    0,
		ConsumerGlobal:          false,
	}
}

func NewAMQPConsumer(
	exchange, queueName, queueDle string,
	declareExchange bool,
	routingKeys []string,
	consumer rabbitmq.IConsumer,
) rabbitmq.AMQPConsumer {
	queueArgs := amqp.Table{
		"x-queue-type":     "quorum",
		"x-delivery-limit": 3,
	}
	if queueDle != "" {
		queueArgs["x-dead-letter-exchange"] = queueDle
	}

	return rabbitmq.AMQPConsumer{
		ExchangeParams: rabbitmq.ExchangeParams{
			Name:       exchange,
			Type:       amqp.ExchangeTopic,
			Durable:    true,
			AutoDelete: false,
			Internal:   false,
			Nowait:     false,
			Args:       nil,
			// custom flag
			DeclareExchange: declareExchange,
		},
		QueueParams: rabbitmq.QueueParams{
			Name:       queueName,
			Durable:    true,
			AutoDelete: false,
			Exclusive:  false,
			Nowait:     false,
			Args:       queueArgs,
		},
		QueueBindParams: rabbitmq.QueueBindParams{
			Nowait: false,
			Args:   nil,
		},
		ConsumerParams: rabbitmq.ConsumerParams{
			RoutingKeys: routingKeys,
			ConsumerID:  uuid.NewString(),
			AutoAck:     false,
			Exclusive:   false,
			NoLocal:     false,
			Nowait:      false,
			Args:        nil,
		},
		IConsumer: consumer,
	}
}

func getTlsConfig(cfg *conf.Config) (*tls.Config, error) {
	var tlsConfig *tls.Config

	if cfg.NewRabbitMQUseTLS {
		tlsConfig = new(tls.Config)

		tlsConfig.RootCAs = x509.NewCertPool()

		ca, err := os.ReadFile(cfg.NewRabbitMQCACert)
		if err != nil {
			return nil, err
		}
		tlsConfig.RootCAs.AppendCertsFromPEM(ca)

		cert, err := tls.LoadX509KeyPair(cfg.NewRabbitMQClientCert, cfg.NewRabbitMQClientKey)
		if err != nil {
			return nil, err
		}
		tlsConfig.Certificates = append(tlsConfig.Certificates, cert)
	}

	return tlsConfig, nil
}

func RabbitMQueueName(cfg *conf.Config, queueName string) string {
	return mqutils.NewQueueName(cfg.AppName, cfg.AppENV, queueName)
}

func NewConsumerTraceLoggerMid() mqutils.Middleware {
	return func(next rabbitmq.IConsumer) rabbitmq.IConsumer {
		return rabbitmq.ConsumerFunc(func(ctx context.Context, msg amqp.Delivery) {
			span, newCtx := opentracing.StartSpanFromContext(ctx, "LOG_MESSAGE")
			defer span.Finish()

			span.LogFields(log.String("body", string(msg.Body)))

			next.Consume(newCtx, msg)
		})
	}
}

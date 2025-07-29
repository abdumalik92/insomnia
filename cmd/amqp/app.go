package amqp

import (
	"context"
	"errors"

	"github.com/alifcapital/components/jaeger"
	"github.com/alifcapital/keycloak_module/cmd/amqp/hr"
	"github.com/alifcapital/keycloak_module/cmd/amqp/keycloak"
	"github.com/alifcapital/keycloak_module/conf"
	"github.com/alifcapital/keycloak_module/internal/v1"
	"github.com/alifcapital/keycloak_module/src/iam"
	"github.com/alifcapital/rabbitmq/mqutils"
	"github.com/opentracing/opentracing-go"
	"go.uber.org/fx"
)

func Run() {
	app := fx.New(
		provide(),

		fx.Invoke(invoke),
	)

	app.Run()
}

func provide() fx.Option {
	return fx.Provide(
		conf.NewViper,
		conf.NewConfig,

		v1.NewValidator,
		v1.NewKeycloakHTTPClient,
		v1.NewZapLogger,
		v1.NewPublisher,

		iam.NewService,
		iam.NewUserAttributesFactory,
		fx.Annotate(v1.NewKeycloakHTTPClient,
			fx.As(new(iam.IUserRepository)),
			fx.As(new(iam.IPasswordBroker)),
		),

		mqutils.NewPool,

		newSentry,
		newJaegerConf,
		newAMQPErrorHandler,
		newAMQPPanicRecoveryCallback,
	)
}

func invoke(
	lc fx.Lifecycle,
	pool *mqutils.Pool,
	tracer opentracing.Tracer,
	jaegerCloser jaeger.Closer,
	keycloakClient keycloak.Client,
	hrClient hr.Client,
	_ sentryInst,
) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			opentracing.SetGlobalTracer(tracer)
			keycloakClientErr := keycloakClient.Start()
			hrClientErr := hrClient.Start()
			return errors.Join(keycloakClientErr, hrClientErr)
		},
		OnStop: func(ctx context.Context) error {
			errCloseErr := pool.Close()
			tracerCloseErr := jaegerCloser.Close()
			return errors.Join(errCloseErr, tracerCloseErr)
		},
	})
}

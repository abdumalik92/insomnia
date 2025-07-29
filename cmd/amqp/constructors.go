package amqp

import (
	"github.com/alifcapital/components/jaeger"
	"github.com/alifcapital/keycloak_module/conf"
	"github.com/getsentry/sentry-go"
	"github.com/opentracing/opentracing-go"
)

type sentryInst struct{}

func newSentry(cfg *conf.Config) (sentryInst, error) {
	if err := sentry.Init(sentry.ClientOptions{Dsn: cfg.SentryDSN}); err != nil {
		return sentryInst{}, err
	}
	return sentryInst{}, nil
}

func newJaegerConf(cfg *conf.Config) (opentracing.Tracer, jaeger.Closer, error) {
	jaegerCfg := jaeger.Config{
		AgentHostPort: cfg.JaegerAgentHostPort,
		ServiceName:   cfg.JaegerServiceName,
		SamplerType:   cfg.JaegerSamplerType,
		SamplerParam:  cfg.JaegerSamplerParam,
		LogsEnables:   cfg.JaegerLogsEnabled,
		LogSpans:      cfg.JaegerReporterLogSpans,
	}
	return jaeger.NewJaegerTracer(jaegerCfg)
}

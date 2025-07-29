package amqp

import (
	"context"
	"fmt"
	"github.com/alifcapital/rabbitmq/mqutils"
	"github.com/getsentry/sentry-go"
	amqp "github.com/rabbitmq/amqp091-go"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

// newAMQPErrorHandler handles errors with trace + logs
func newAMQPErrorHandler(logger *zap.Logger, tracer trace.Tracer) mqutils.ErrorHandler {
	return func(ctx context.Context, msg amqp.Delivery, err error) bool {
		_, span := tracer.Start(ctx, "amqp.ErrorHandler")
		defer span.End()

		span.RecordError(err)
		span.SetAttributes(
			attribute.String("routing_key", msg.RoutingKey),
			attribute.String("exchange", msg.Exchange),
			attribute.String("body", string(msg.Body)),
		)

		// sentry
		sentry.CaptureException(err)

		// zap
		logger.Error("error",
			zap.Error(err),
			zap.String("msg.MessageId", msg.MessageId),
			zap.String("msg.Exchange", msg.Exchange),
			zap.String("msg.Body", string(msg.Body)),
		)

		return false // no requeue
	}
}

//func newAMQPErrorHandler(logger *zap.Logger) mqutils.ErrorHandler {
//	return func(ctx context.Context, msg amqp.Delivery, err error) bool {
//		span, _ := opentracing.StartSpanFromContext(ctx, "amqp.ErrorHandler")
//		defer span.Finish()
//
//		// opentracing
//		ext.LogError(span, err)
//		span.LogFields(tracelog.String("reflect.TypeOf(err)", reflect.TypeOf(err).String()))
//		span.LogFields(tracelog.String("routing_key", msg.RoutingKey))
//		span.LogFields(tracelog.String("exchange", msg.Exchange))
//		span.LogFields(tracelog.String("body", string(msg.Body)))
//
//		// sentry
//		sentry.CaptureException(err)
//
//		// zap
//		logger.Error("error",
//			zap.Error(err),
//			zap.String("msg.MessageId", msg.MessageId),
//			zap.String("msg.Exchange", msg.Exchange),
//			zap.String("msg.Body", string(msg.Body)),
//		)
//
//		// no requeue for now
//		return false
//	}
//}

// newAMQPPanicRecoveryCallback handles panics with trace + logs
func newAMQPPanicRecoveryCallback(logger *zap.Logger, tracer trace.Tracer) mqutils.PanicRecoveryCallback {
	return func(ctx context.Context, msg amqp.Delivery, recErr any) {
		_, span := tracer.Start(ctx, "amqp.PanicRecoveryCallback")
		defer span.End()

		err := fmt.Errorf("panic: %+v", recErr)
		span.RecordError(err)
		span.SetAttributes(
			attribute.String("recErr", err.Error()),
			attribute.String("panic", "true"),
			attribute.String("routing_key", msg.RoutingKey),
			attribute.String("exchange", msg.Exchange),
			attribute.String("body", string(msg.Body)),
		)

		// sentry
		sentry.CaptureException(err)

		// zap
		logger.Error("panic",
			zap.Error(err),
			zap.String("msg.MessageId", msg.MessageId),
			zap.String("msg.Exchange", msg.Exchange),
			zap.String("msg.Body", string(msg.Body)),
		)

		// NACK msg
		if err := msg.Nack(false, false); err != nil {
			span.RecordError(err)
		}
	}
}

//func newAMQPPanicRecoveryCallback(logger *zap.Logger) mqutils.PanicRecoveryCallback {
//	return func(ctx context.Context, msg amqp.Delivery, recErr any) {
//		span, _ := opentracing.StartSpanFromContext(ctx, "amqp.PanicRecoveryCallback")
//		defer span.Finish()
//
//		err := fmt.Errorf("panic: %+v", recErr)
//
//		// opentracing
//		ext.LogError(span, err)
//		span.SetTag("panic", "true")
//		span.LogFields(tracelog.Object("recErr", recErr))
//		span.LogFields(tracelog.String("routing_key", msg.RoutingKey))
//		span.LogFields(tracelog.String("exchange", msg.Exchange))
//		span.LogFields(tracelog.String("body", string(msg.Body)))
//
//		// sentry
//		sentry.CaptureException(err)
//
//		// zap
//		logger.Error("panic",
//			zap.Error(err),
//			zap.String("msg.MessageId", msg.MessageId),
//			zap.String("msg.Exchange", msg.Exchange),
//			zap.String("msg.Body", string(msg.Body)),
//		)
//
//		// NACK msg
//		if err := msg.Nack(false, false); err != nil {
//			span.LogFields(tracelog.String("nack_failed", err.Error()))
//		}
//	}
//}

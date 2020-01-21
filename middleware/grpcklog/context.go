package grpcklog

import (
	"context"

	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
)

type ctxLoggerMarker struct{}

type ctxLogger struct {
	fields map[string]interface{}
}

var (
	ctxLoggerKey = &ctxLoggerMarker{}
)

// AddFields adds fields to the logger.
func AddFields(ctx context.Context, fields map[string]interface{}) {
	l, ok := ctx.Value(ctxLoggerKey).(*ctxLogger)
	if !ok || l == nil {
		return
	}
	for k, v := range fields {
		l.fields[k] = v
	}
}

// Extract takes the call-scoped logrus.Entry from ctx_logrus middleware.
//
// If the ctx_logrus middleware wasn't used, a no-op `logrus.Entry` is returned. This makes it safe to
// use regardless.
func Extract(ctx context.Context) map[string]interface{} {
	l, ok := ctx.Value(ctxLoggerKey).(*ctxLogger)
	if !ok || l == nil {
		return nil
	}

	fields := make(map[string]interface{})

	// Add grpc_ctxtags tags metadata until now.
	tags := grpc_ctxtags.Extract(ctx)
	for k, v := range tags.Values() {
		fields[k] = v
	}

	// Add logrus fields added until now.
	for k, v := range l.fields {
		fields[k] = v
	}
	return fields
}

// ToContext adds the logrus.Entry to the context for extraction later.
// Returning the new context that has been created.
func ToContext(ctx context.Context, callLog map[string]interface{}) context.Context {
	l := &ctxLogger{
		fields: callLog,
	}
	return context.WithValue(ctx, ctxLoggerKey, l)
}

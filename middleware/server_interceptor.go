package grpcklog

import (
	"context"
	"path"
	"time"

	"github.com/go-logr/logr"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"google.golang.org/grpc"
)

// UnaryServerInterceptor returns a new unary server interceptors that adds fields to the context.
func UnaryServerInterceptor(log logr.Logger, opts ...Option) grpc.UnaryServerInterceptor {
	o := evaluateServerOpt(opts)
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		startTime := time.Now()
		newCtx := newLoggerForCall(ctx, info.FullMethod, startTime)

		resp, err := handler(newCtx, req)

		if !o.shouldLog(info.FullMethod, err) {
			return resp, err
		}
		code := o.codeFunc(err)
		level := o.levelFunc(code)
		durField, durVal := o.durationFunc(time.Since(startTime))
		fields := Extract(newCtx)
		fields[durField] = durVal
		fields["grpc.code"] = code.String()

		levelLogf(log, level, "finished streaming call with code "+code.String(), fields, err)

		return resp, err
	}
}

// StreamServerInterceptor returns a new streaming server interceptor that adds fields to the context.
func StreamServerInterceptor(log logr.Logger, opts ...Option) grpc.StreamServerInterceptor {
	o := evaluateServerOpt(opts)
	return func(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		startTime := time.Now()
		newCtx := newLoggerForCall(stream.Context(), info.FullMethod, startTime)
		wrapped := grpc_middleware.WrapServerStream(stream)
		wrapped.WrappedContext = newCtx

		err := handler(srv, wrapped)

		if !o.shouldLog(info.FullMethod, err) {
			return err
		}
		code := o.codeFunc(err)
		level := o.levelFunc(code)
		durField, durVal := o.durationFunc(time.Since(startTime))
		fields := Extract(newCtx)
		fields[durField] = durVal
		fields["grpc.code"] = code.String()
		if err != nil {
			fields["error"] = err
		}

		levelLogf(log, level, "finished streaming call with code "+code.String(), fields, err)

		return err
	}
}

func newLoggerForCall(ctx context.Context, fullMethodString string, start time.Time) context.Context {
	service := path.Dir(fullMethodString)[1:]
	method := path.Base(fullMethodString)
	callLog := map[string]interface{}{
		"SystemField":     "grpc server",
		"grpc.service":    service,
		"grpc.method":     method,
		"grpc.start_time": start.Format(time.RFC3339),
	}

	if d, ok := ctx.Deadline(); ok {
		callLog = map[string]interface{}{
			"grpc.request.deadline": d.Format(time.RFC3339),
		}
	}

	fromcontext := Extract(ctx)
	for k, v := range fromcontext {
		callLog[k] = v
	}
	return ToContext(ctx, callLog)
}

func levelLogf(log logr.Logger, level KlogLevel, format string, fields map[string]interface{}, err error) {
	switch level {
	case InfoLog:
		log.Info("Info - "+format, "details", fields)
	case WarningLog:
		log.Info("Warning - "+format, "details", fields)
	case ErrorLog:
		log.Error(err, format, "details", fields)
	}
}

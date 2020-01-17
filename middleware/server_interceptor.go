package grpcklog

import (
	"context"
	"fmt"
	"path"
	"time"

	"google.golang.org/grpc"
	"k8s.io/klog/klogr"
)

// UnaryServerInterceptor returns a new unary server interceptors that adds logrus.Entry to the context.
func UnaryServerInterceptor(opts ...Option) grpc.UnaryServerInterceptor {
	o := evaluateServerOpt(opts)
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		startTime := time.Now()
		newCtx := newLoggerForCall(ctx, info.FullMethod, startTime)

		resp, err := handler(newCtx, req)

		if !o.shouldLog(info.FullMethod, err) {
			return resp, err
		}
		durField, durVal := o.durationFunc(time.Since(startTime))

		fields := Extract(newCtx)
		fields[durField] = durVal

		var values string
		for k, v := range fields {
			values += fmt.Sprintln(k, v)
		}

		log := klogr.New().WithName("MyName").WithValues("user", "Dan")

		//log.Info("finished unary call with code", "val1", err.Error(), "val2", map[string]int{"k": 1})
		log.Info("finished unary call with code", "val1", err.Error(), "val2", values)

		return resp, err
	}
}

func newLoggerForCall(ctx context.Context, fullMethodString string, start time.Time) context.Context {
	service := path.Dir(fullMethodString)[1:]
	method := path.Base(fullMethodString)
	callLog := map[string]interface{}{
		"SystemField":     "grpc",
		"KindField":       "server",
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

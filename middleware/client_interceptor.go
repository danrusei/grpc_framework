package grpcklog

import (
	"context"
	"fmt"
	"path"
	"time"

	"google.golang.org/grpc"
	"k8s.io/klog/klogr"
)

// UnaryClientInterceptor returns a new unary client interceptor that optionally logs the execution of external gRPC calls.
func UnaryClientInterceptor(opts ...Option) grpc.UnaryClientInterceptor {
	o := evaluateClientOpt(opts)
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		fields := newClientLoggerFields(ctx, method)
		startTime := time.Now()
		err := invoker(ctx, method, req, reply, cc, opts...)
		logFinalClientLine(o, fields, startTime, err, "finished client unary call")
		return err
	}
}

func newClientLoggerFields(ctx context.Context, fullMethodString string) map[string]interface{} {
	service := path.Dir(fullMethodString)[1:]
	method := path.Base(fullMethodString)
	return map[string]interface{}{
		"SystemField":  "grpc",
		"KindField":    "client",
		"grpc.service": service,
		"grpc.method":  method,
	}
}

func logFinalClientLine(o *options, fields map[string]interface{}, startTime time.Time, err error, msg string) {
	durField, durVal := o.durationFunc(time.Now().Sub(startTime))

	fields[durField] = durVal

	var values string
	for k, v := range fields {
		values += fmt.Sprintln(k, v)
	}

	if err != nil {
		fields["error"] = err
	}
	log := klogr.New().WithName("MyName").WithValues("user", "Dan")

	//log.Info("finished unary call with code", "val1", err.Error(), "val2", map[string]int{"k": 1})
	log.Info("finished unary call with code", "val1", err.Error(), "val2", values)
}

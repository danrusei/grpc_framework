package grpcklog

import (
	"context"
	"path"
	"time"

	api "github.com/Danr17/grpc_framework/proto"
	"github.com/go-logr/logr"
	"google.golang.org/grpc"
)

// UnaryClientInterceptor returns a new unary client interceptor that optionally logs the execution of external gRPC calls.
func UnaryClientInterceptor(log logr.Logger, opts ...Option) grpc.UnaryClientInterceptor {
	o := evaluateClientOpt(opts)
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		request := req.(*api.ClientRequestType)
		log.Info("requesting all product types from vendor: " + request.Vendor)
		fields := newClientLoggerFields(ctx, method)
		startTime := time.Now()
		err := invoker(ctx, method, req, reply, cc, opts...)
		logFinalClientLine(o, log, fields, startTime, err, "finished client unary call")
		return err
	}
}

func newClientLoggerFields(ctx context.Context, fullMethodString string) map[string]interface{} {
	service := path.Dir(fullMethodString)[1:]
	method := path.Base(fullMethodString)
	return map[string]interface{}{
		"SystemField":  "grpc client",
		"grpc.service": service,
		"grpc.method":  method,
	}
}

func logFinalClientLine(o *options, log logr.Logger, fields map[string]interface{}, startTime time.Time, err error, msg string) {
	code := o.codeFunc(err)
	level := o.levelFunc(code)
	durField, durVal := o.durationFunc(time.Now().Sub(startTime))
	fields[durField] = durVal

	if err != nil {
		fields["error"] = err
	}

	var name string = "Unknown User"
	if val, ok := fields["Name"]; ok {
		name = val.(string)
	}

	log.WithName(name).WithValues("user", "Dan")

	switch level {
	case InfoLog:
		log.WithName(name).Info("Info - The call finished with code "+code.String(), "details", fields)
	case WarningLog:
		log.Info("Warning - The call finished with code "+code.String(), "details", fields)
	case ErrorLog:
		log.Error(err, "The call finished with code "+code.String(), "details", fields)
	}
}

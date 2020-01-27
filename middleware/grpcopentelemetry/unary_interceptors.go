package grpcopentelemetry

// gRPC tracing middleware
// https://github.com/open-telemetry/opentelemetry-specification/blob/master/specification/data-rpc.md
import (
	"context"

	"go.opentelemetry.io/otel/plugin/grpctrace"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"go.opentelemetry.io/otel/api/core"
	"go.opentelemetry.io/otel/api/distributedcontext"
	"go.opentelemetry.io/otel/api/global"
	"go.opentelemetry.io/otel/api/key"
	"go.opentelemetry.io/otel/api/trace"
)

// UnaryServerInterceptor intercepts and extracts incoming trace data
func UnaryServerInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	requestMetadata, _ := metadata.FromIncomingContext(ctx)
	metadataCopy := requestMetadata.Copy()

	entries, spanCtx := grpctrace.Extract(ctx, &metadataCopy)
	ctx = distributedcontext.WithMap(ctx, distributedcontext.NewMap(distributedcontext.MapUpdate{
		MultiKV: entries,
	}))

	grpcServerKey := key.New("grpc.server")
	serverSpanAttrs := []core.KeyValue{
		grpcServerKey.String("api-server"),
	}

	tr := global.TraceProvider().Tracer("grpc_tracer")
	ctx, span := tr.Start(
		ctx,
		"Cloud-Products-types",
		trace.WithAttributes(serverSpanAttrs...),
		trace.ChildOf(spanCtx),
		trace.WithSpanKind(trace.SpanKindServer),
	)
	defer span.End()

	return handler(ctx, req)
}

// UnaryClientInterceptor intercepts and injects outgoing trace
func UnaryClientInterceptor(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	requestMetadata, _ := metadata.FromOutgoingContext(ctx)
	metadataCopy := requestMetadata.Copy()

	tr := global.TraceProvider().Tracer("grpc_tracer")
	err := tr.WithSpan(ctx, "Cloud-Products-types",
		func(ctx context.Context) error {
			grpctrace.Inject(ctx, &metadataCopy)
			ctx = metadata.NewOutgoingContext(ctx, metadataCopy)

			err := invoker(ctx, method, req, reply, cc, opts...)
			setTraceStatus(ctx, err)
			return err
		})
	return err
}

func setTraceStatus(ctx context.Context, err error) {
	if err != nil {
		s, _ := status.FromError(err)
		trace.SpanFromContext(ctx).SetStatus(s.Code())
	} else {
		trace.SpanFromContext(ctx).SetStatus(codes.OK)
	}
}

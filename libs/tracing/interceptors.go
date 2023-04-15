package tracing

import (
	"context"
)

import (
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/uber/jaeger-client-go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func Tracer(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (resp interface{}, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, info.FullMethod)
	defer span.Finish()

	span.SetTag("method", info.FullMethod)

	if spanContext, ok := span.Context().(jaeger.SpanContext); ok {
		ctx = metadata.AppendToOutgoingContext(ctx, "x-trace-id", spanContext.TraceID().String())
	}

	res, err := handler(ctx, req)

	if status.Code(err) != codes.OK {
		ext.Error.Set(span, true)
	}

	span.SetTag("status_code", status.Code(err).String())

	return res, err
}

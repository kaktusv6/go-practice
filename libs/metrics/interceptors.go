package metrics

import (
	"context"
	"time"
)

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

func Metrics(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (resp interface{}, err error) {
	RequestsCounter.Inc()

	timeStart := time.Now()

	res, err := handler(ctx, req)

	elapsed := time.Since(timeStart)

	statusStr := status.Code(err).String()

	HistogramResponseTime.WithLabelValues(statusStr).Observe(elapsed.Seconds())
	ResponseCounter.WithLabelValues(statusStr).Inc()
	return res, err
}

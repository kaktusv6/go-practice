package tracing

import (
	"github.com/uber/jaeger-client-go/config"
	"go.uber.org/zap"
	"route256/libs/logger"
)

func Init(serviceName string) {
	config := &config.Configuration{
		ServiceName: serviceName,
		Sampler: &config.SamplerConfig{
			Type:  "const",
			Param: 1,
		},
	}
	config, _ = config.FromEnv()

	_, err := config.InitGlobalTracer(serviceName)

	if err != nil {
		logger.Fatal("Cannot init tracing", zap.Error(err))
	}
}

package v1

import (
	"io"
	"os"

	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/config"
	jaegercfg "github.com/uber/jaeger-client-go/config"
	"github.com/uber/jaeger-lib/metrics/prometheus"
)

type JaegerInfo struct {
	JAEGER_SERVICE_NAME string
	JAEGER_AGENT_HOST   string
	JAEGER_AGENT_PORT   string
	JAEGER_ENDPOINT     string
}

const (
	operationName = "test-gogowang"
)

func (jaegerInfo *JaegerInfo) InitFromEnvironment() *JaegerInfo {
	os.Setenv("JAEGER_AGENT_HOST", jaegerInfo.JAEGER_AGENT_HOST)
	os.Setenv("JAEGER_AGENT_PORT", jaegerInfo.JAEGER_AGENT_PORT)
	os.Setenv("JAEGER_SERVICE_NAME", jaegerInfo.JAEGER_SERVICE_NAME)
	os.Setenv("JAEGER_REPORTER_FLUSH_INTERVAL", "5ns")
	return jaegerInfo
}

func (jaegerInfo *JaegerInfo) NewClient() (opentracing.Tracer, io.Closer, error) {
	metricsFactory := prometheus.New()
	tracer, closer, err := config.Configuration{
		ServiceName: "your-service-name",
		Sampler: &jaegercfg.SamplerConfig{
			Type:  jaeger.SamplerTypeRateLimiting,
			Param: 0.001,
		},
		Reporter: &jaegercfg.ReporterConfig{
			LogSpans: true,
		},
	}.NewTracer(
		config.Metrics(metricsFactory),
	)
	opentracing.SetGlobalTracer(tracer)
	return tracer, closer, err
}

func (jaegerInfo *JaegerInfo) InitStartSpan(tracer opentracing.Tracer, opts ...opentracing.StartSpanOption) opentracing.Span {
	return tracer.StartSpan(operationName, opts...)
}

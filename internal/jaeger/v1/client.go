package v1

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
)

func InitClients() error {
	info := &JaegerInfo{
		JAEGER_AGENT_PORT:   "16831",
		JAEGER_AGENT_HOST:   "9.135.225.72",
		JAEGER_SERVICE_NAME: "gogo-wang-test-service",
		// JAEGER_ENDPOINT:     "http://9.135.225.72:6831/api/traces",
	}

	tracer, closer, err := info.InitFromEnvironment().NewClient()
	defer closer.Close()
	if err != nil {
		panic(err)
	}
	fmt.Println(tracer, err)
	span := info.InitStartSpan(tracer)
	span = span.SetTag("test-gogo-tag", "test-gogo-tag--value")
	defer span.Finish()
	Recursion(context.Background(), span, 0)
	return nil
}

func Recursion(ctxparam context.Context, span opentracing.Span, flag int) {

	if flag > 1<<3 {
		return
	}
	ctx := opentracing.ContextWithSpan(ctxparam, span)
	reqSpan, _ := opentracing.StartSpanFromContext(ctx, "test-"+strconv.Itoa(flag))
	defer reqSpan.Finish()
	span = span.SetTag("test-gogo-tag", "test-gogo-tag--value")
	reqSpan.LogKV("info-key-1"+strconv.Itoa(flag), "info-value-1"+strconv.Itoa(flag))
	reqSpan.LogFields(
		log.String("event", "event-value-1"+strconv.Itoa(flag)),
		log.String("value", "value-1"+strconv.Itoa(flag)),
	)
	flag++
	time.Sleep(time.Millisecond * 10)
	ctxR := opentracing.ContextWithSpan(context.Background(), span)
	Recursion(ctxR, reqSpan, flag)
}

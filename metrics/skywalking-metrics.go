package metrics

import (
	"context"
	"github.com/SkyAPM/go2sky"
	http2 "github.com/SkyAPM/go2sky/plugins/http"
	"github.com/SkyAPM/go2sky/reporter"
	"github.com/julioisaac/daxxer-api/internal/logs"
	"go.uber.org/zap"
	"net/http"
	"sync"
)

var once sync.Once
var Instance *Metrics
var instanceTracer *go2sky.Tracer

func Setup() *Metrics {
	once.Do(func() {
		Instance = &Metrics{
			Tracer:     getTracer(),
			Middleware: getMiddleware(),
			Client:     getClient(),
		}
	})
	return Instance
}

func getTracer() interface{} {
	if instanceTracer == nil {
		re, err := reporter.NewGRPCReporter("skywalking-oap:11800")

		if err != nil {
			logs.Instance.Log.Error(context.Background(), "error creating reporter", zap.Error(err))
			return nil
		}

		tracer, err := go2sky.NewTracer("daxxer-api", go2sky.WithReporter(re))
		if err != nil {
			logs.Instance.Log.Error(context.Background(), "error creating tracer", zap.Error(err))
			return nil
		}
		logs.Instance.Log.Debug(context.Background(), "reporter created")
		instanceTracer = tracer
	}
	return instanceTracer
}

func getMiddleware() interface{} {
	tracer := getTracer().(*go2sky.Tracer)
	middleware, err := http2.NewServerMiddleware(tracer)
	if err != nil {
		logs.Instance.Log.Error(context.Background(), "error creating server middleware", zap.Error(err))
		return nil
	}
	return middleware
}

func getClient() *http.Client   {
	tracer := getTracer().(*go2sky.Tracer)
	if tracer == nil {
		logs.Instance.Log.Warn(context.Background(), "error trying to recover tracer")
		return nil
	}
	client, err := http2.NewClient(tracer)
	if err != nil {
		logs.Instance.Log.Debug(context.Background(), "error trying to create client", zap.Error(err))
		return nil
	}
	return client
}
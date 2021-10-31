package metrics

import (
	"github.com/SkyAPM/go2sky"
	http2 "github.com/SkyAPM/go2sky/plugins/http"
	"github.com/SkyAPM/go2sky/reporter"
	"log"
	"net/http"
	"os"
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
		// log
		logger := log.New(os.Stderr, "WithLogger", log.LstdFlags)
		options := reporter.WithLogger(logger)
		re, err := reporter.NewGRPCReporter("skywalking-oap:11800", options)

		if err != nil {
			logger.Fatalf("new reporter error %v \n", err)
			return nil
		}

		tracer, err := go2sky.NewTracer("daxxer-api", go2sky.WithReporter(re))
		if err != nil {
			logger.Fatalf("create tracer error %v \n", err)
			return nil
		}
		logger.Print("reporter created...")
		instanceTracer = tracer
	}
	return instanceTracer
}

func getMiddleware() interface{} {
	tracer := getTracer().(*go2sky.Tracer)
	middleware, err := http2.NewServerMiddleware(tracer)
	if err != nil {
		log.Fatalf("create server middleware error %v \n", err)
		return nil
	}
	return middleware
}

func getClient() *http.Client   {
	tracer := getTracer().(*go2sky.Tracer)
	if tracer == nil {
		log.Fatalf("recovering tracer error\n")
		return nil
	}
	client, err := http2.NewClient(tracer)
	if err != nil {
		log.Fatalf("create client error %v \n", err)
		return nil
	}
	return client
}
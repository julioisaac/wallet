package metrics

import (
	"fmt"
	"github.com/SkyAPM/go2sky"
	http2 "github.com/SkyAPM/go2sky/plugins/http"
	"github.com/SkyAPM/go2sky/reporter"
	"log"
	"net/http"
	"os"
)

var (
	instanceTracer *go2sky.Tracer
	instanceClient *http.Client
)

type metrics struct{}

func Metric() *metrics {
	return &metrics{}
}

func (m *metrics) GetTracer() *go2sky.Tracer {
	if instanceTracer == nil {
		logger := log.New(os.Stderr, "daxxer-wallet-api", log.LstdFlags)
		options := reporter.WithLogger(logger)
		re, err := reporter.NewGRPCReporter("oap:11800", options)

		if err != nil {
			fmt.Printf("error creating reporter")
			return nil
		}

		tracer, err := go2sky.NewTracer("daxxer-api", go2sky.WithReporter(re), go2sky.WithInstance("daxxer-service-1"))
		if err != nil {
			fmt.Printf("error creating tracer")
			return nil
		}
		fmt.Printf("reporter created")
		instanceTracer = tracer
	}
	return instanceTracer
}

func (m *metrics) GetClient() *http.Client   {
	if instanceClient == nil {
		tracer := m.GetTracer()
		if tracer == nil {
			fmt.Printf("error trying to recover tracer")
			return nil
		}
		client, err := http2.NewClient(tracer)
		if err != nil {
			fmt.Printf("error trying to create client")
			return nil
		}
		instanceClient = client
	}
	return instanceClient
}
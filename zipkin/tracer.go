package zipkin

import (
	"github.com/openzipkin/zipkin-go"
	httpreporter "github.com/openzipkin/zipkin-go/reporter/http"
	"go.uber.org/zap"
	"time"
)

func NewHttpTracer(conf HttpTracerConf, logger *zap.Logger) *zipkin.Tracer {
	// create a reporter to be used by the tracer
	if len(conf.ReporterAddress) == 0 {
		conf.ReporterAddress = "http://localhost:9411/api/v2/spans"
	}
	zipkinReporter := httpreporter.NewReporter(conf.ReporterAddress,
		httpreporter.Serializer(JSONSerializer{}))
	//
	if len(conf.ServiceName) == 0 || len(conf.ServiceHostPort) == 0 {
		logger.Fatal("service not config")
	}
	// set-up the local endpoint for our service
	endpoint, err := zipkin.NewEndpoint(conf.ServiceName, conf.ServiceHostPort)
	if err != nil {
		logger.Fatal("unable to create local endpoint", zap.Error(err))
	}
	// set-up our sampling strategy
	sampler, err := zipkin.NewBoundarySampler(0.01, time.Now().UnixNano())
	if err != nil {
		logger.Fatal("unable to create sampler", zap.Error(err))
	}
	// initialize the tracer
	tracer, err := zipkin.NewTracer(
		zipkinReporter,
		zipkin.WithLocalEndpoint(endpoint),
		zipkin.WithSampler(sampler),
	)
	if err != nil {
		logger.Fatal("unable to create tracer", zap.Error(err))
	}
	logger.Info("created zipkin tracer", zap.Any("conf", conf))
	return tracer
}

package httpclient

import (
	"github.com/openzipkin/zipkin-go"
	zipkinhttp "github.com/openzipkin/zipkin-go/middleware/http"
	"go.uber.org/zap"
)

func NewZipkinClient(tracer *zipkin.Tracer, log *zap.Logger) *zipkinhttp.Client {
	// create global zipkin traced http client
	client, err := zipkinhttp.NewClient(tracer,
		zipkinhttp.ClientTrace(true),
		zipkinhttp.ClientTags(map[string]string{
			"clientType": "golang",
		}),
		zipkinhttp.TransportOptions(zipkinhttp.RoundTripper(NewFastHttpTransport())))
	if err != nil {
		log.Fatal("unable to create client", zap.Error(err))
	}
	return client
}

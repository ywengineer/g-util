package es

import (
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esapi"
	jsoniter "github.com/json-iterator/go"
	zipkinhttp "github.com/openzipkin/zipkin-go/middleware/http"
	"github.com/ywengineer/g-util/client"
	"go.uber.org/zap"
	"net/http"
)

var defaultTransport = client.NewFastHttpTransport()

func NewESClientWithZipkin(address []string, client *zipkinhttp.Client, log *zap.Logger) *esapi.API {
	cfg := elasticsearch.Config{
		Addresses: address,
		Transport: &innerTransport{
			client: client,
		},
	}
	es, err := elasticsearch.NewClient(cfg)
	if err != nil {
		log.Fatal("Error creating es client", zap.Error(err))
	} else {
		infoJson := make(map[string]interface{})
		if info, e := es.Info(func(request *esapi.InfoRequest) {
			request.Human = true
			request.Pretty = true
		}); e != nil {
			log.Fatal("Error creating es client through info method", zap.Error(e))
		} else if e := jsoniter.NewDecoder(info.Body).Decode(&infoJson); e != nil {
			log.Fatal("Error creating es client through info method", zap.Error(e))
		} else {
			log.Info("es client created.", zap.Any("info", infoJson))
			return esapi.New(es)
		}
	}
	return nil
}

func NewESClient(address []string, log *zap.Logger) *esapi.API {
	return NewESClientWithZipkin(address, nil, log)
}

// Transport implements the estransport interface with
// the github.com/valyala/fasthttp HTTP client.
//
type innerTransport struct {
	client *zipkinhttp.Client
}

// RoundTrip performs the request and returns a response or error
//
func (t *innerTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	if t.client == nil {
		return defaultTransport.RoundTrip(req)
	}
	return t.client.DoWithAppSpan(req, "elastic")
}

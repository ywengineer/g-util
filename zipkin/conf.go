package zipkin

import (
	jsoniter "github.com/json-iterator/go"
	"github.com/openzipkin/zipkin-go/model"
)

type HttpTracerConf struct {
	ReporterAddress string `json:"reporter_address"`
	ServiceName     string `json:"service_name"`
	ServiceHostPort string `json:"service_host_port"`
}

// JSONSerializer implements the default JSON encoding SpanSerializer.
type JSONSerializer struct{}

// Serialize takes an array of Zipkin SpanModel objects and returns a JSON
// encoding of it.
func (JSONSerializer) Serialize(spans []*model.SpanModel) ([]byte, error) {
	return jsoniter.Marshal(spans)
}

// ContentType returns the ContentType needed for this encoding.
func (JSONSerializer) ContentType() string {
	return "application/json;charset=UTF-8"
}

package es

import (
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"go.uber.org/zap"
	"io/ioutil"
)

func NewESClient(address []string, log *zap.Logger) *elasticsearch.Client {
	cfg := elasticsearch.Config{
		Addresses: address,
		Transport: &transport{},
		//Transport: &http.Transport{
		//	MaxIdleConnsPerHost:   10,
		//	ResponseHeaderTimeout: time.Millisecond,
		//	DialContext:           (&net.Dialer{Timeout: time.Nanosecond}).DialContext,
		//	TLSClientConfig: &tls.Config{
		//		MinVersion: tls.VersionTLS11,
		//		// ...
		//	},
		//},
	}
	es, err := elasticsearch.NewClient(cfg)
	if err != nil {
		log.Fatal("Error creating the client", zap.Error(err))
	} else {
		if info, e := es.Info(func(request *esapi.InfoRequest) {
			request.Human = true
			request.Pretty = true
		}); e != nil {
			log.Fatal("Error creating the client through info method", zap.Error(e))
		} else if t, e := ioutil.ReadAll(info.Body); e != nil {
			log.Fatal("Error creating the client through info method", zap.Error(e))
		} else {
			log.Info("elastic client created.", zap.Any("info", string(t)))
		}
	}
	return es
}

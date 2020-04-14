package es

import (
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esapi"
	jsoniter "github.com/json-iterator/go"
	"go.uber.org/zap"
)

func NewESClient(address []string, log *zap.Logger) *esapi.API {
	cfg := elasticsearch.Config{
		Addresses: address,
		Transport: &FastHttpTransport{},
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
